package task_queue

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

type taskQueue struct {
	tasks         chan Task
	wg            sync.WaitGroup
	imagesService images.ServiceType
	mu            sync.Mutex
}

func NewService(capacity int, service images.ServiceType) *taskQueue {
	return &taskQueue{
		tasks:         make(chan Task, capacity),
		imagesService: service,
	}
}

func (q *taskQueue) Start(stopChan <-chan struct{}) {
	fmt.Println("starting task queue...")
	logging.Infof("queue size: %d, capacity: %d", len(q.tasks), cap(q.tasks))

	go q.Process()

	<-stopChan
	fmt.Println("Stopping task queue...")

	q.Close()
	q.Wait()
	fmt.Println("Task queue stopped")
}

func (q *taskQueue) Enqueue(task Task) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	select {
	case q.tasks <- task:
		return true
	default:
		return false
	}
}

func (q *taskQueue) Process() {
	fmt.Println("task queue starts running...")

	for task := range q.tasks {
		q.wg.Add(1)
		err := withTimeout(task, 2*time.Second)
		if err != nil {
			logging.Errorf("WithTimeout() failed, err: %s", err)
		}
		q.wg.Done()
	}
}

func (q *taskQueue) Wait() {
	q.wg.Wait()
}

func (q *taskQueue) Close() {
	fmt.Println("closing task queue...")
	close(q.tasks)
}

func withTimeout(task Task, timeout time.Duration) error {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	errChan := make(chan error, 1)

	go func() {
		defer func() {
			close(errChan)
		}()

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
		return fmt.Errorf("resizeImages() timed out")
	}
}
