package async

import "context"

var _ easyWaiter = NewEasyWait()

// EasyWait implements [async.TaskWaiter] and can be embedded to tasks.
// Sample usage:
//
//	type PublishEventTask struct {
//		Event interface{}
//		*async.EasyWait
//	}
//	... // implement Task interface
//	func main() {
//		task := PublishEventTask {
//			Event: struct{}{},
//			EasyWait: async.NewEasyWait(),
//		}
//		async.Start()
//		task.Wait() // function provided by async.EasyWait
//	}
type EasyWait struct {
	isDone   bool
	doneChan chan struct{}
}

// NewEasyWait creates a new instance of EasyWait
func NewEasyWait() *EasyWait {
	return &EasyWait{
		doneChan: make(chan struct{}, 1),
	}
}

// Wait for the task to be done. Note that this is not thread-safe, only one routine should be waiting.
// Do not execute a task multiple times as well.
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

// EasyTask is an embeddable struct which provides some functions required by
// [async.Task]. By embedding EasyTask, you only need to implement [async.Task.Do]
// Sample usage:
//
//	type PublishEventTask struct {
//		Event interface{}
//		async.EasyTask
//		*async.EasyWait
//	}
//	func (p *PublishEventTask) Do() { // only this needs to be implemented
//		// do task...
//		p.TaskStatus = async.TaskStatusSuccessful // update the task status
//	}
//	func main() {
//		task := PublishEventTask {
//			Event: struct{}{},
//			EasyWait: async.NewEasyWait(),
//		}
//		task.TaskContext = context.Background()
//		async.Start()
//		task.Wait() // function provided by async.EasyWait
//	}
type EasyTask struct {
	TaskContext context.Context
	TaskStatus  TaskStatus
}

// Status returns the status of the task
func (et EasyTask) Status() TaskStatus {
	return et.TaskStatus
}

// Context returns the context of the task
func (et EasyTask) Context() context.Context {
	return et.TaskContext
}
