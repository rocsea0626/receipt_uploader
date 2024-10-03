package resize_queue

import (
	"context"
	"fmt"
	"receipt_uploader/internal/constants"
	"receipt_uploader/internal/images"
	"receipt_uploader/internal/logging"
	"receipt_uploader/internal/models/tasks"
	"sync"
	"time"
)

type TaskFunc func() error

type ResizeQueue struct {
	DoneChan      chan struct{}
	tasks         chan tasks.ResizeTask
	wg            sync.WaitGroup
	imagesService images.ServiceType
	mu            sync.Mutex
}

func NewService(capacity int, service images.ServiceType) *ResizeQueue {
	return &ResizeQueue{
		DoneChan:      make(chan struct{}),
		tasks:         make(chan tasks.ResizeTask, capacity),
		imagesService: service,
	}
}

func (q *ResizeQueue) Start(stopChan <-chan struct{}) {
	fmt.Println("starting task queue...")
	logging.Infof("queue size: %d, capacity: %d", len(q.tasks), cap(q.tasks))

	go q.process()

	<-stopChan
	fmt.Println("Stopping task queue...")

	q.close()
	q.wait()
	fmt.Println("Task queue stopped")
	close(q.DoneChan)
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

func (q *ResizeQueue) process() {
	fmt.Println("task queue starts running...")

	q.wg.Add(1)
	for task := range q.tasks {
		err := q.withTimeout(task, constants.RESIZE_TIMEOUT*time.Second)
		if err != nil {
			logging.Errorf("WithTimeout() failed, path: '%s', err: %s", task.ImageMeta.Path, err)
		}
	}
	q.wg.Done()
}

func (q *ResizeQueue) wait() {
	fmt.Printf("processing remaining %d enqueued tasks\n", len(q.tasks))
	q.wg.Wait()
}

func (q *ResizeQueue) close() {
	fmt.Println("closing task queue...")
	close(q.tasks)
}

func (q *ResizeQueue) withTimeout(task tasks.ResizeTask, timeout time.Duration) error {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	errChan := make(chan error, 1)

	go func() {
		defer close(errChan)

		startTime := time.Now()
		err := q.imagesService.GenerateResizedImages(&task.ImageMeta, task.DestDir)
		if err != nil {
			errChan <- fmt.Errorf("GenerateResizedImages() failed, err: %w", err)
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
