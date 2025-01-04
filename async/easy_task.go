package async

var _ easyWaiter = NewEasyWait()

type EasyWait struct {
	isDone   bool
	doneChan chan struct{}
}

func NewEasyWait() *EasyWait {
	return &EasyWait{
		doneChan: make(chan struct{}, 1),
	}
}

// Wait for the task to be done. Note that this is not thread-safe, only one routine should be waiting
func (ew *EasyWait) Wait() {
	if ew.isDone {
		return
	}
	<-ew.doneChan
}

// Done marks the task as done. Note that this is not thread-safe and only one routine
// should be calling this. It is not recommended for this to be called outside the library.
func (ew *EasyWait) done() {
	ew.doneChan <- struct{}{}
	ew.isDone = true
}

type easyWaiter interface {
	Wait()
	done()
}
