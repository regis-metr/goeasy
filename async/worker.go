package async

import (
	"sync"
)

type workerStatus int

const (
	workerStatusPending = 0
	workerStatusWorking = 1
)

type worker struct {
	taskQueue chan Task
	status    workerStatus
	statusMu  sync.Mutex
}

func newWorker(taskQueue chan Task) *worker {
	return &worker{
		taskQueue: taskQueue,
		status:    workerStatusPending,
	}
}

func (w *worker) start() {
	w.statusMu.Lock()
	defer w.statusMu.Unlock()

	if w.status == workerStatusWorking {
		return
	}
	w.status = workerStatusWorking
	go w.doStart()

}

func (w *worker) doStart() {
	for w.status == workerStatusWorking {
		// TODO: stop chan or check if taskQueue is closed
		task := <-w.taskQueue

		// TODO: check context deadline
		task.Do()

		if easyWaiterTask, ok := task.(easyWaiter); ok {
			easyWaiterTask.done()
		}

	}
}

func (w *worker) stop() {
	w.statusMu.Lock()
	defer w.statusMu.Unlock()
	w.status = workerStatusPending
}
