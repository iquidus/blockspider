package syncronizer

func (s *Synchronizer) startTaskHandler() {

	// As tasks are created with s.AddLink, and block at task.Link(), this goroutine will
	// receive them from the channel; it waits until the tasks calls t.Link() (if it didn't already),
	// resumes execution and waits for it to finish, doing this for each task until
	// the channel is empty or one of the tasks calls t.Abort()

	go func() {
		var abort bool
	loop:
		for {
			select {
			case task := <-s.routines:

				task.wait()

				abort = s.didAbort(task)

				if abort {
					s.aborted = true
					task.stop()
					break loop
				}

				task.release()

				task.finish()

				abort = s.didAbort(task)

				if abort {
					s.aborted = true
					task.closeNext()
					break loop
				}

			// Sometimes, if we add n tasks in a loop, the first one will be executed right away
			// (go scheduler quirks??) before others are inserted, so the number of items in the channel
			// can't be relied upon to quit the taskHandler;
			// s.shouldQuit() will return true only when s.Finish() has been called
			// and the default case will only run if there are no more task to receive (len(s.routines) == 0)

			default:
				if s.shouldQuit() {
					s.quit()
					break loop
				}
			}
		}
		if s.aborted {
			s.flushTasks()
		}
		return
	}()
}

type Synchronizer struct {
	routines, abortChan   chan *Task
	quitChan, nextChannel chan int
	aborted               bool
}

// AddLink creates a new task with the function body it's provided, sets up hooks and
// queues it for execution

func (s *Synchronizer) AddLink(body func(*Task)) {

	if s.aborted {
		return
	}

	nr := newTask(s, body, s.nextChannel)

	c := make(chan int)
	s.nextChannel = c

	nr.abortFunc = func() {
		close(c)
	}

	go nr.run()

	s.routines <- nr

	return
}

// Finish hangs until all tasks have completed executions, and there are no more tasks
// or if the sync was aborted
// when finish is called no new tasks should be added

func (s *Synchronizer) Finish() (aborted bool) {
	s.quitChan <- 0
	for {
		select {
		case _, more := <-s.nextChannel:
			if more {
				return false
			}
			return true
		}
	}
}

func (s *Synchronizer) quit() {
	s.nextChannel <- 0
}

func (s *Synchronizer) shouldQuit() bool {
	select {
	case <-s.quitChan:
		return true
	default:
		return false
	}
}

// Check abortChan if any task sent an abort signal
// returns true if it is equal to the one that called the method

func (s *Synchronizer) didAbort(t *Task) bool {
	select {
	case closedTask := <-s.abortChan:
		if closedTask == t {
			return true
		} else {
			s.abortChan <- closedTask
			return false
		}
	default:
		return false
	}
}

// Sometimes, if there is an abort when len(s.routines) == maxRoutines, and there's a call to
// AddLink stuck on inserting a task, everything blocks. To avoid that, after aborting a sync
// we flush tasks from the channel so AddLink can send task and return

// Todo: maybe it's quicker if we just remove one task as all possible calls to AddLink would return immediately

func (s *Synchronizer) flushTasks() {
loop:
	for {
		select {
		case <-s.routines:
			if len(s.routines) == 0 {
				break loop
			}
		}
	}
	return
}
