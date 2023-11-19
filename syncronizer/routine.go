package syncronizer

type Task struct {
	ranInit, hang, done chan int
	fn                  func()
	abort               bool
	abortFunc           func()
	activeSync          *Synchronizer
}

// Link should be called exactly once inside the fn of every task.
func (r *Task) Link() (close bool) {
	r.ranInit <- 0
	close = r.receive()
	return
}

func (r *Task) AbortSync() {
	r.activeSync.abortChan <- r
}

func (r *Task) stop() {
	close(r.hang)
}

func (r *Task) closeNext() {
	r.abortFunc()
}

func (r *Task) wait() {
	<-r.ranInit
}

func (r *Task) release() {
	r.hang <- 0
}

func (r *Task) finish() {
	<-r.done
}

func (r *Task) run() {
	r.fn()
}

func (r *Task) receive() (closed bool) {
	for {
		select {
		case _, more := <-r.hang:
			if more {
				return false
			} else {
				r.closeNext()
				return true
			}
		}
	}
}

func newTask(s *Synchronizer, fn func(*Task), hang chan int) *Task {
	// Buffered channels so hooks don't block when they're not supposed to
	r := &Task{make(chan int, 1), hang, make(chan int, 1), nil, false, nil, s}

	rFn := func() {
		fn(r)

		r.done <- 0
	}

	r.fn = rFn

	return r
}
