package state

import (
	"errors"
	"sync"
	"time"

	"github.com/iquidus/blockspider/cache"
	"github.com/iquidus/blockspider/common"
	"github.com/iquidus/blockspider/disk"
)

type State struct {
	Syncing bool    `json:"syncing"`
	Config  *Config `json:"config"`
	Cache   *cache.BlockStack[common.Block]
}

type StateData struct {
	ChainId   *uint64      `json:"chainId"`
	Head      common.Block `json:"head"`
	Timestamp int64        `json:"updated"`
}

type Config struct {
	Path       string `json:"path"`
	CacheLimit int    `json:"cache"`
}

type StateFile struct {
	ChainId   *uint64        `json:"chainId"`
	Head      common.Block   `json:"head"`
	Timestamp int64          `json:"updated"`
	Cache     []common.Block `json:"cache,omitempty"`
}

var state *StateData = nil
var lock sync.Mutex

// create new state instance
func Init(cfg *Config, chainId *uint64, startBlock common.Block) (*State, error) {
	s := &State{
		Syncing: false,
		Config:  cfg,
		Cache:   cache.New[common.Block](&cfg.CacheLimit),
	}
	err := s.load()
	if err != nil {
		// set singleton
		state = &StateData{
			ChainId:   chainId,
			Head:      startBlock,
			Timestamp: time.Now().Unix(),
		}
		// write to disc
		err = s.save()
		if err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (s *State) Get() (*StateData, error) {
	if state != nil {
		// state is set, return it
		return state, nil
	} else {
		// read from disc
		err := s.load()
		if err != nil {
			return nil, errors.New("State has not be initialized. run Init() first")
		} else {
			return state, nil
		}
	}
}

func (s *State) Update(block common.Block) error {
	// s.Cache.Push(state.Head)
	state.Head = block
	state.Timestamp = time.Now().Unix()
	return s.save()
}

func (s *State) load() error {
	lock.Lock()
	defer lock.Unlock()
	err := disk.ReadJsonFile[StateData](s.Config.Path, state)
	if err != nil {
		return err
	}
	return nil
}

func (s *State) save() error {
	lock.Lock()
	defer lock.Unlock()
	var payload = StateFile{
		ChainId:   state.ChainId,
		Head:      state.Head,
		Timestamp: state.Timestamp,
		// Cache: s.Cache.Items(),
	}
	err := disk.WriteJsonFile[StateFile](payload, s.Config.Path, 0644)
	if err != nil {
		return err
	}
	return nil
}
