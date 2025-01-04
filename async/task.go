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
