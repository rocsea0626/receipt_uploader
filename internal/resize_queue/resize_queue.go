package resize_queue

import (
	"context"
	"fmt"
	"receipt_uploader/internal/images"
	"receipt_uploader/internal/logging"
	"receipt_uploader/internal/models/tasks"
	"sync"
	"time"
)

type TaskFunc func() error

type ResizeQueue struct {
	tasks         chan tasks.ResizeTask
	wg            sync.WaitGroup
	imagesService images.ServiceType
	mu            sync.Mutex
}

func NewService(capacity int, service images.ServiceType) *ResizeQueue {
	return &ResizeQueue{
		tasks:         make(chan tasks.ResizeTask, capacity),
		imagesService: service,
	}
}

func (q *ResizeQueue) Start(stopChan <-chan struct{}) {
	fmt.Println("starting task queue...")
	logging.Infof("queue size: %d, capacity: %d", len(q.tasks), cap(q.tasks))

	go q.Process()

	<-stopChan
	fmt.Println("Stopping task queue...")

	q.Close()
	q.Wait()
	fmt.Println("Task queue stopped")
}

func (q *ResizeQueue) Enqueue(task tasks.ResizeTask) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	select {
	case q.tasks <- task:
		return true
	default:
		return false
	}
}

func (q *ResizeQueue) Process() {
	fmt.Println("task queue starts running...")

	for task := range q.tasks {
		q.wg.Add(1)
		err := q.WithTimeout(task, 2*time.Second)
		if err != nil {
			logging.Errorf("withTimeout() failed, path: '%s', err: %s", task.ImageMeta.Path, err)
		}
		q.wg.Done()
	}
}

func (q *ResizeQueue) Wait() {
	q.wg.Wait()
}

func (q *ResizeQueue) Close() {
	fmt.Println("closing task queue...")
	close(q.tasks)
}

func (q *ResizeQueue) WithTimeout(task tasks.ResizeTask, timeout time.Duration) error {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	errChan := make(chan error, 1)

	go func() {
		defer func() {
			close(errChan)
		}()

		startTime := time.Now()
		err := q.imagesService.GenerateResizedImages(&task.ImageMeta, task.DestDir)
		if err != nil {
			errChan <- fmt.Errorf("s.ImageService.GenerateResizedImages() failed, err: %w", err)
			return
		}
		elapsedTime := time.Since(startTime)
		logging.Infof("resizeImages() completes with %d ms", elapsedTime.Milliseconds())
	}()

	select {
	case genErr := <-errChan:
		return genErr
	case <-ctx.Done():
		return fmt.Errorf("resizeImages() timed out")
	}
}
