package worker_pool

import (
	"context"
	"fmt"
	"receipt_uploader/internal/images"
	"receipt_uploader/internal/logging"
	"sync"
	"time"
)

type TaskFunc func() error
type Task struct {
	Name string
	Func TaskFunc
}

type WorkerPool struct {
	DoneChan      chan struct{}
	tasks         chan Task
	wg            sync.WaitGroup
	imagesService images.ServiceType
	mu            sync.Mutex
	workerCount   int
}

func NewService(capacity, workerCount int, service images.ServiceType) *WorkerPool {
	return &WorkerPool{
		DoneChan:      make(chan struct{}),
		tasks:         make(chan Task, capacity),
		imagesService: service,
		workerCount:   workerCount,
	}
}

func (wp *WorkerPool) Start(stopChan <-chan struct{}) {
	fmt.Println("starting worker pool...")
	logging.Infof("queue capacity: %d, num of workers: %d", cap(wp.tasks), wp.workerCount)

	for i := 0; i < wp.workerCount; i++ {
		go wp.worker(i)
	}

	<-stopChan
	fmt.Println("stopping worker pool...")

	wp.close()

	wp.wait()
	fmt.Println("worker pool stopped")
	close(wp.DoneChan)
}

func (wp *WorkerPool) Submit(task Task) bool {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	select {
	case wp.tasks <- task:
		return true
	default:
		return false
	}
}

func (wp *WorkerPool) wait() {
	wp.wg.Wait()
}

func (wp *WorkerPool) close() {
	fmt.Printf("closing task queue, len(tasks): %d\n", len(wp.tasks))
	close(wp.tasks)
}

func (wp *WorkerPool) processTask(task Task) {
	wp.wg.Add(1)

	err := withTimeout(task, 2*time.Second)
	if err != nil {
		logging.Errorf("withTimeout() failed, err: %s", err)
	}
	logging.Infof("withTimeout() done")
	wp.wg.Done()
}

func (wp *WorkerPool) worker(workerId int) {
	fmt.Printf("worker_%d started...\n", workerId)

	for task := range wp.tasks {
		wp.processTask(task)
	}
}

func withTimeout(task Task, timeout time.Duration) error {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	errChan := make(chan error, 1)

	go func() {
		defer close(errChan)

		startTime := time.Now()
		err := task.Func()
		if err != nil {
			errChan <- fmt.Errorf("task failed, task: %v, err: %w", task, err)
			return
		}
		elapsedTime := time.Since(startTime)
		logging.Infof("task completes with %d ms, task: %s", elapsedTime.Milliseconds(), task.Name)
	}()

	select {
	case genErr := <-errChan:
		return genErr
	case <-ctx.Done():
		return fmt.Errorf("%s timed out", task.Name)
	}
}
