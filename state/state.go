package state

import (
	"errors"
	"sync"
	"time"

	"github.com/iquidus/blockspider/common"
	"github.com/iquidus/blockspider/disk"
)

type State struct {
	Syncing bool    `json:"syncing"`
	Config  *Config `json:"config"`
}

type StateData struct {
	ChainId   *uint64      `json:"chainId"`
	Head      common.Block `json:"head"`
	Timestamp int64        `json:"updated"`
}

type Config struct {
	Path string `json:"path"`
}

var state *StateData = nil
var lock sync.Mutex

// create new state instance
func Init(cfg *Config, chainId *uint64, startBlock common.Block) (*State, error) {
	s := &State{
		Syncing: false,
		Config:  cfg,
	}
	err := load(cfg.Path)
	if err != nil {
		// set singleton
		state = &StateData{
			ChainId:   chainId,
			Head:      startBlock,
			Timestamp: time.Now().Unix(),
		}
		// write to disc
		err = save(cfg.Path)
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
		err := load(s.Config.Path)
		if err != nil {
			return nil, errors.New("State has not be initialized. run Init() first")
		} else {
			return state, nil
		}
	}
}

func (s *State) Update(block common.Block) error {
	state.Head = block
	state.Timestamp = time.Now().Unix()
	return save(s.Config.Path)
}

func load(path string) error {
	lock.Lock()
	defer lock.Unlock()
	err := disk.ReadJsonFile[StateData](path, state)
	if err != nil {
		return err
	}
	return nil
}

func save(path string) error {
	lock.Lock()
	defer lock.Unlock()
	err := disk.WriteJsonFile[StateData](*state, path, 0644)
	if err != nil {
		return err
	}
	return nil
}
