package state

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"
	"time"

	"github.com/iquidus/blockspider/common"
)

type State struct {
	Syncing bool    `json:"syncing"`
	Config  *Config `json:"config"`
}

type StateData struct {
	ChainId   *uint64         `json:"chainId"`
	Head      common.RawBlock `json:"head"`
	Timestamp int64           `json:"updated"`
}

type Config struct {
	Path string `json:"path"`
}

var state *StateData = nil
var lock sync.Mutex

var Marshal = func(v interface{}) (io.Reader, error) {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

var Unmarshal = func(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

// create new state instance
func Init(cfg *Config, chainId *uint64, startBlock common.RawBlock) (*State, error) {
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

func (s *State) Update(block common.RawBlock) error {
	state.Head = block
	state.Timestamp = time.Now().Unix()
	return save(s.Config.Path)
}

func load(path string) error {
	lock.Lock()
	defer lock.Unlock()
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return Unmarshal(f, &state)
}

func save(path string) error {
	lock.Lock()
	defer lock.Unlock()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	r, err := Marshal(state)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, r)
	return err
}
