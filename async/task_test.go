package async

import (
	"context"
	"testing"

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
	doFunc  func()
	*EasyWait
}

func (t *easyWaiterTask) Do() {
	if t.doFunc != nil {
		t.doFunc()
	}
}

func (t *easyWaiterTask) Context() context.Context {
	return t.context
}

func (t *easyWaiterTask) Status() TaskStatus {
	return t.status
}

func TestStart(t *testing.T) {
	taskExecuted := false
	task := easyWaiterTask{
		doFunc: func() {
			taskExecuted = true
		},
		EasyWait: NewEasyWait(),
	}
	Start(&task)
	task.Wait()

	assert.True(t, taskExecuted, "Task was not executed")

}
