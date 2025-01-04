package test

import (
	"context"

	"testing"

	"github.com/bryan-t/goeasy/async"
)

type easyWaiterTask struct {
	status  async.TaskStatus
	context context.Context
	t       *testing.T
	*async.EasyWait
}

func (t *easyWaiterTask) Do() {
	t.status = async.TaskStatusSuccessful
}

func (t *easyWaiterTask) Context() context.Context {
	return t.context
}

func (t *easyWaiterTask) Status() async.TaskStatus {
	return t.status
}

func TestStartEasyWait(t *testing.T) {
	opts := async.WorkerPoolOptions{
		Workers: 2,
	}
	workerPool := async.NewWorkerPool(opts)
	workerPool.Start()

	task := easyWaiterTask{
		EasyWait: async.NewEasyWait(),
		t:        t,
	}

	workerPool.AddTask(&task)
	(&task).Wait()

	workerPool.Stop()
}
