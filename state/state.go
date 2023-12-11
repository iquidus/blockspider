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
	ChainId   *uint64 `json:"chainId"`
	Timestamp int64   `json:"updated"`
}

type Config struct {
	Path       string `json:"path"`
	CacheLimit int    `json:"cache"`
}

type StateFile struct {
	ChainId   *uint64        `json:"chainId"`
	Timestamp int64          `json:"updated"`
	Cache     []common.Block `json:"cache"`
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
			Timestamp: time.Now().Unix(),
		}
		s.Cache.Push(startBlock)
		// write to disc
		err = s.save()
		if err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (s *State) Save() error {
	return s.save()
}

func (s *State) load() error {
	lock.Lock()
	defer lock.Unlock()
	var sf StateFile
	err := disk.ReadJsonFile[StateFile](s.Config.Path, &sf)
	if err != nil {
		return err
	}
	state = &StateData{
		ChainId:   sf.ChainId,
		Timestamp: sf.Timestamp,
	}
	if sf.Cache == nil {
		return errors.New("cache is nil")
	}
	for i := len(sf.Cache) - 1; i >= 0; i-- {
		s.Cache.Push(sf.Cache[i])
	}
	return nil
}

func (s *State) save() error {
	lock.Lock()
	defer lock.Unlock()
	var sf = StateFile{
		ChainId:   state.ChainId,
		Timestamp: time.Now().Unix(),
		Cache:     s.Cache.Items(),
	}
	err := disk.WriteJsonFile[StateFile](sf, s.Config.Path, 0644)
	if err != nil {
		return err
	}
	return nil
}
