package syncronizer

import (
	"os"

	"github.com/ethereum/go-ethereum/log"
)

// Returns a new sync object with no routines
// Tasks should be linked inside the task function body via Task.Link()
// After syncing all tasks, it is necessary to call sync.Finish() to let
// the syncronizer know it can quit

func NewSync(maxRoutines int) *Synchronizer {
	if maxRoutines == 0 {
		log.Error("Error, cannot start sync with 0 maxroutines, should be atleast 1")
		os.Exit(1)
	}

	s := &Synchronizer{}

	s.routines = make(chan *Task, maxRoutines)

	// Buffered channels so sends on these don't block
	s.abortChan = make(chan *Task, maxRoutines)
	s.quitChan = make(chan int, 1)

	// Unbuffered so this blocks
	s.nextChannel = make(chan int)

	s.startTaskHandler()

	return s
}
