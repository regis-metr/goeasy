package async

import "sync"

// WorkerPoolOptions contains settings of WorkerPool
type WorkerPoolOptions struct {
	// Workers indicates the number of workers(routines) to be spawned
	Workers int
}

// WorkerPool maintains a group of workers which limits the number of routines that are spawned.
type WorkerPool struct {
	mu        sync.Mutex
	workers   []*worker
	taskQueue chan Task
}

// NewWorkerPool creates a new instance of WorkerPool
func NewWorkerPool(options WorkerPoolOptions) *WorkerPool {
	workerPool := WorkerPool{}
	for i := 0; i < options.Workers; i++ {
		workerPool.workers = append(workerPool.workers, newWorker(workerPool.taskQueue))
	}

	return &workerPool
}

// Start the workers
func (wp *WorkerPool) Start() error {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	for _, worker := range wp.workers {
		worker.start()
	}
	return nil
}

// Stop the workers
func (wp *WorkerPool) Stop() error {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	for _, worker := range wp.workers {
		worker.stop()
	}
	return nil
}

// AddTask adds the task the pool of tasks
func (wp *WorkerPool) AddTask(task Task) error {
	wp.taskQueue <- task
	return nil
}
