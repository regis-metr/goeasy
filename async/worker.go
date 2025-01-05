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
	stopChan  chan struct{}
}

func newWorker(taskQueue chan Task) *worker {
	return &worker{
		taskQueue: taskQueue,
		status:    workerStatusPending,
		stopChan:  make(chan struct{}),
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
		task := <-w.taskQueue
		if task == nil {
			continue
		}

		// TODO: check context deadline
		doTask(task)

	}
}

func (w *worker) stop() {
	w.statusMu.Lock()
	defer w.statusMu.Unlock()
	w.status = workerStatusPending
}
