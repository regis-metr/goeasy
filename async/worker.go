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
	wg        *sync.WaitGroup
}

func newWorker(taskQueue chan Task, wg *sync.WaitGroup) *worker {
	return &worker{
		taskQueue: taskQueue,
		status:    workerStatusPending,
		wg:        wg,
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
	w.wg.Add(1)
	for w.status == workerStatusWorking {
		task := <-w.taskQueue
		if task == nil {
			continue
		}

		// TODO: check context deadline
		doTask(task)

	}
	w.wg.Done()
}

func (w *worker) stop() {
	w.statusMu.Lock()
	defer w.statusMu.Unlock()
	w.status = workerStatusPending
}
