package crawler

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/iquidus/blockspider/common"
	"github.com/iquidus/blockspider/kafka"
	"github.com/iquidus/blockspider/syncronizer"
)

func (c *Crawler) RunLoop() {
	// create log channel
	c.logChan = make(chan *logObject)
	// start crawling blocks
	c.crawlBlocks()
	// close log channel
	close(c.logChan)
	// update crawler state
	c.state.Syncing = false
}

func (c *Crawler) crawlBlocks() {
	// check if a sync is already in progress
	if c.state.Syncing {
		c.logger.Warn("Sync already in progress; quitting.")
		return
	}
	// set syncing true to block additional syncs
	c.state.Syncing = true

	// get local head
	state, err := c.state.Get()
	if err != nil || state == nil {
		c.logger.Error("couldn't get head from state", "err", err)
		return
	}

	// add head to cache
	c.cache.Push(state.Head)
	c.logger.Debug("fetched block from local state", "number", state.Head)

	// get remote head
	chainHead, err := c.rpc.LatestBlockNumber()
	if err != nil {
		c.logger.Error("couldn't get block number", "err", err)
	}
	c.logger.Debug("fetched block from node", "number", chainHead)

	// set current block to head + 1
	currentBlock := state.Head.Number + 1

	// create and start sync logger
	syncLogger := c.logger.New()
	startLogger(c.logChan, syncLogger)
	start := time.Now()
	syncLogger.Debug("started sync at", "t", start)

	// add new sync to task chain
	taskChain := syncronizer.NewSync(c.cfg.MaxRoutines)
	for ; currentBlock <= chainHead; currentBlock++ {
		// capture blockNumber
		b := currentBlock
		// add link to task chain
		taskChain.AddLink(func(r *syncronizer.Task) {
			// get remote block
			rawBlock, err := c.rpc.GetBlockByHeight(b)
			if err != nil {
				syncLogger.Error("failed getting block", "err", err)
				c.state.Syncing = false
				r.AbortSync()
				return
			}

			// convert remote block to common.Block
			block, err := rawBlock.Convert(*c.rpc)
			if err != nil {
				syncLogger.Error("failed converting block", "err", err)
				c.state.Syncing = false
				r.AbortSync()
				return
			}

			// check if sync should abort
			abort := r.Link()
			if abort {
				syncLogger.Debug("Aborting routine")
				return
			}
			// process
			c.syncBlock(block, r)
		})
	}

	abort := taskChain.Finish()

	if abort {
		syncLogger.Debug("Aborted sync")
	} else {
		syncLogger.Debug("terminated sync", "t", time.Since(start))
	}
}

// validates local head against remote block with same height
// returns the valid block, dropped block, isValid, error
func (c *Crawler) validateBlock() (*common.Block, *common.Block, bool, error) {
	// if there's no blocks in chain bail out
	if c.cache.Count() > 0 {
		// remove local block from cache
		local, _ := c.cache.Pop()
		// fetch remote block from node
		rawRemote, err := c.rpc.GetBlockByHeight(local.Number)
		if err != nil {
			return nil, nil, false, err
		}
		// compares local and remote block hash
		if local.Hash == rawRemote.Hash {
			return &local, nil, true, nil
		} else {
			// convert remote block to common.Block
			remote, err := rawRemote.Convert(*c.rpc)
			if err != nil {
				return nil, nil, false, err
			}
			return &remote, &local, false, nil
		}
	} else {
		return nil, nil, false, errors.New("No blocks in chain to validate")
	}
}

func (c *Crawler) reorg() error {
	var commonAncestor *common.Block
	sidechainmap := make(map[uint64]common.Block)
	sidechain := []common.Block{}
	dropped := []common.Block{}

	for {
		// loop until common ancestor is found
		if commonAncestor == nil {
			// compare local "head" against remote block
			b, d, ok, _ := c.validateBlock()
			if !ok && b != nil {
				// if compare fails check to make sure we are not already
				// handling this block
				_, e := sidechainmap[b.Number]
				if !e {
					// store in map to prevent multiple fires
					sidechainmap[b.Number] = *b
					// add to sidechain (map is not ordered so use this)
					sidechain = append(sidechain, *b)
					if d != nil {
						// if a block was dropped add to slice
						dropped = append(dropped, *d)
					}
				}
			} else {
				// first block match. set common ancestor then end loop.
				commonAncestor = b
				break
			}
		}
	}

	// log common ancenstor
	c.logger.Warn("Common ancestor found", "block", commonAncestor.Number, "hash", commonAncestor.Hash)
	// common ancestor was popped off the chain during above loop, push it back on
	c.cache.Push(*commonAncestor)

	// process old blocks
	for i := 0; i < len(dropped); i++ {
		c.logger.Warn("Dropping local block", "number", dropped[i].Number, "hash", dropped[i].Hash)
		err := c.sendReorgHooks(dropped[i])
		if err != nil {
			return errors.New("Failed to send reorg hook: " + err.Error())
		}
	}

	// process new blocks
	for i := len(sidechain) - 1; i >= 0; i-- {
		c.cache.Push(sidechain[i])
		c.logger.Info("Adding remote block", "number", sidechain[i].Number, "hash", sidechain[i].Hash)
		err := c.sendBlockMessage(dropped[i])
		if err != nil {
			return errors.New("Failed to send reorg hook: " + err.Error())
		}
	}

	return nil
}

func (c *Crawler) sendBlockMessage(block common.Block) error {
	var bp = kafka.BlocksPayload{
		Method: "PUSH",
		Block:  block,
	}

	payload, err := json.Marshal(bp)
	if err != nil {
		return err
	}

	err = c.blockWriter.WriteMessages(context.Background(), payload)

	if err != nil {
		c.logger.Error("failed to write messages", "err", err)
	}

	return nil
}

func (c *Crawler) sendReorgHooks(block common.Block) error {
	var bp = kafka.BlocksPayload{
		Method: "POP",
		Block:  block,
	}

	payload, err := json.Marshal(bp)
	if err != nil {
		return err
	}

	err = c.blockWriter.WriteMessages(context.Background(), payload)
	if err != nil {
		return err
	}

	return nil
}

// func (c *Crawler) sendEventsMessage(events []common.Log, topic string) error {
// 	var ep = kafka.EventsPayload{
// 		Method: "PUSH",
// 		Events: events,
// 	}

// 	payload, err := json.Marshal(ep)
// 	if err != nil {
// 		return err
// 	}

// 	err = c.eventWriter.WriteMessagesWithTopic(context.Background(), payload, "events")
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (c *Crawler) syncBlock(block common.Block, task *syncronizer.Task) {
	// get parent block from cache
	parent, err := c.cache.Peak()
	if err == nil {
		if parent.Hash != block.ParentHash {
			// A reorg has occurred
			c.logger.Warn("Chain reorg detected", "parent", parent.Number, "hash", parent.Hash, "block", block.Number, "hash", block.Hash, "parent", block.ParentHash)
			err := c.reorg()
			if err != nil {
				c.logger.Error("Failed to determine common ancestor", "err", err)
			}
			newHead, err := c.cache.Peak()
			if err != nil {
				c.logger.Error("Failed to peak head from cache", "err", err)
			}
			err = c.state.Update(newHead)
			if err != nil {
				c.logger.Error("Failed to update local state after reorg", "err", err)
				return
			}
			// abort sync
			task.AbortSync()
			return
		}
	} else {
		c.logger.Error("Failed to peak block cache", "err", err)
	}

	// handle block hook here
	err = c.sendBlockMessage(block)
	if err != nil {
		c.logger.Error("Failed to send block hook", "err", err)
	}

	// TODO(iquidus): Running a getlogs with each filter is expensive
	// 	              migrate to a single getlogs call (or derive from block itself), and filter locally

	// handle getlogs requests
	// logopts := c.cfg.Kafka.Events
	// for i := 0; i < len(logopts); i++ {
	// 	logs, err := c.rpc.GetLogs(logopts[i].Addresses, block.Hash, logopts[i].Topics)
	// 	if err != nil {
	// 		c.logger.Error("Failed to get logs", "err", err)
	// 		return
	// 	}
	// 	if len(logs) > 0 {
	// 		// handle webhook here
	// 		err = c.sendEventsMessage(logs, logopts[i].Topic)
	// 		if err != nil {
	// 			c.logger.Error("Failed to send getlogs hook", "err", err)
	// 		}
	// 	}
	// }

	// write block to state
	err = c.state.Update(block)
	if err != nil {
		c.logger.Error("err", err)
		return
	}

	// add block to cache for next iteration
	c.cache.Push(block)

	// log
	c.log(block.Number, len(block.Logs))
}

func (c *Crawler) log(blockNo uint64, txns int) {
	c.logChan <- &logObject{
		blockNo: blockNo,
		txns:    txns,
	}
}
