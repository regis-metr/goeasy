package async

import "context"

type TaskStatus int

const (
	TaskStatusPending    TaskStatus = 0
	TaskStatusSuccessful TaskStatus = 1
	TaskStatusFailed     TaskStatus = 2
)

type TaskWaiter interface {
	Wait()
}

type Task interface {
	Do()
	Context() context.Context
	Status() TaskStatus
}
