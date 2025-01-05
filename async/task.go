package async

import "context"

// TaskStatus indicates the state of a task
type TaskStatus int

const (
	// TaskStatusPending indicates that the task have not yet been processed.
	TaskStatusPending TaskStatus = 0

	// TaskStatusSuccessful indicates that a task has been successfully processed.
	TaskStatusSuccessful TaskStatus = 1

	// TaskStatusSuccessful indicates that a task failed.
	TaskStatusFailed TaskStatus = 2
)

// TaskWaiter is an interface for tasks that exposes functionality to wait until
// the task is done.
type TaskWaiter interface {
	Wait()
}

// Task is a single unit of work that can be done.
type Task interface {
	// Do the task
	Do()

	// Context should return the context of the task
	Context() context.Context

	// Status should return the status of the task whether it was successful or not
	Status() TaskStatus
}

// Start doing the task
func Start(task Task) {
	go doTask(task)

}

func doTask(task Task) {
	task.Do()
	if easyWaiterTask, ok := task.(easyWaiter); ok {
		easyWaiterTask.done()
	}
}
