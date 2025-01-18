package async

import (
	"fmt"
	"sync"
	"time"
)

const DefaultMaxQueuedTask = 20
const DefaultWorkers = 1

// WorkerPoolOptions contains settings of WorkerPool
type WorkerPoolOptions struct {
	// Workers indicates the number of workers(routines) to be spawned
	Workers int

	// MaxQueuedTask indicates the maximum number of pending tasks to be executed by the worker pool.
	MaxQueuedTask int
}

// WorkerPool maintains a group of workers which limits the number of routines that are spawned.
type WorkerPool struct {
	mu        sync.Mutex
	workers   []*worker
	taskQueue chan Task
	wg        sync.WaitGroup
}

// NewWorkerPool creates a new instance of WorkerPool
func NewWorkerPool(options WorkerPoolOptions) *WorkerPool {
	workerPool := WorkerPool{}

	maxQueuedTask := DefaultMaxQueuedTask
	if options.MaxQueuedTask > 0 {
		maxQueuedTask = options.MaxQueuedTask
	}
	workerPool.taskQueue = make(chan Task, maxQueuedTask)

	workers := DefaultWorkers
	if options.Workers > 0 {
		workers = options.Workers
	}
	for i := 0; i < workers; i++ {
		workerPool.workers = append(workerPool.workers, newWorker(workerPool.taskQueue, &workerPool.wg))
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

	close(wp.taskQueue)

	for len(wp.taskQueue) > 0 {
		time.Sleep(500 * time.Millisecond)
	}

	for _, worker := range wp.workers {
		worker.stop()
	}

	wp.wg.Wait()
	return nil
}

// AddTask adds the task the pool of tasks
func (wp *WorkerPool) AddTask(task Task) error {
	select {
	case wp.taskQueue <- task:
		return nil
	default:
		return fmt.Errorf("queue is full")
	}
}
