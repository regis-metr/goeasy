package async

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testTask struct {
	status  TaskStatus
	context context.Context
	doFunc  func()
	t       *testing.T
}

func (t *testTask) Do() {
	if t.doFunc != nil {
		t.doFunc()
	}
}

func (t *testTask) Context() context.Context {
	return t.context
}

func (t *testTask) Status() TaskStatus {
	return t.status
}

type easyWaiterTask struct {
	status  TaskStatus
	context context.Context
	t       *testing.T
	*EasyWait
}

func (t *easyWaiterTask) Do() {
	t.status = TaskStatusSuccessful
}

func (t *easyWaiterTask) Context() context.Context {
	return t.context
}

func (t *easyWaiterTask) Status() TaskStatus {
	return t.status
}

func TestAddTaskFull(t *testing.T) {
	opts := WorkerPoolOptions{
		Workers:       1,
		MaxQueuedTask: 1,
	}
	wp := NewWorkerPool(opts)
	addErr1 := wp.AddTask(&testTask{})
	addErr2 := wp.AddTask(&testTask{})
	assert.Nil(t, addErr1, "First task returned an error")
	assert.Equal(t, addErr2, fmt.Errorf("queue is full"), addErr2, "Second task did not return expected error")
}

func TestStop(t *testing.T) {
	opts := WorkerPoolOptions{
		Workers: 10,
	}
	wp := NewWorkerPool(opts)
	wp.Start()
	t1Exec, t2Exec := false, false
	wp.AddTask(&testTask{
		doFunc: func() { t1Exec = true },
	})
	wp.AddTask(&testTask{doFunc: func() { t2Exec = true }})
	wp.Stop()

	assert.Equal(t, 0, len(wp.taskQueue), "taskQueue is not empty")
	assert.Equal(t, true, t1Exec, "task 1 not executed")
	assert.Equal(t, true, t2Exec, "task 2 not executed")
}
func TestStopFinishAllTasks(t *testing.T) {
	opts := WorkerPoolOptions{
		Workers: 1,
	}
	wp := NewWorkerPool(opts)
	wp.Start()
	t1Exec, t2Exec := false, false
	wp.AddTask(&testTask{
		doFunc: func() {
			time.Sleep(5 * time.Second)
			t1Exec = true
		},
	})
	wp.AddTask(&testTask{doFunc: func() { t2Exec = true }})
	wp.Stop()

	assert.Equal(t, 0, len(wp.taskQueue), "taskQueue is not empty")
	assert.Equal(t, true, t1Exec, "task 1 not executed")
	assert.Equal(t, true, t2Exec, "task 2 not executed")
}
