package task_queue

import (
	"fmt"
	"receipt_uploader/internal/images"
	"receipt_uploader/internal/logging"
	"receipt_uploader/internal/models/tasks"
	"sync"
)

type TaskQueue struct {
	tasks         chan tasks.ResizeTask
	wg            sync.WaitGroup
	imagesService images.ServiceType
}

func NewService(capacity int, service images.ServiceType) *TaskQueue {
	return &TaskQueue{
		tasks:         make(chan tasks.ResizeTask, capacity),
		imagesService: service,
	}
}

func (tq *TaskQueue) Start(stopChan <-chan struct{}) {
	fmt.Println("starting task queue...")
	logging.Infof("queue size: %d, capacity: %d", len(tq.tasks), cap(tq.tasks))

	go tq.Process()

	<-stopChan
	fmt.Println("Stopping task queue...")

	tq.Close()
	tq.Wait()
	fmt.Println("Task queue stopped")
}

func (tq *TaskQueue) Enqueue(task tasks.ResizeTask) bool {
	logging.Debugf("queue size: %d, capacity: %d", len(tq.tasks), cap(tq.tasks))

	select {
	case tq.tasks <- task:
		return true
	default:
		return false
	}
}

func (tq *TaskQueue) Process() {
	fmt.Println("task queue starts running...")

	for task := range tq.tasks {
		tq.wg.Add(1)
		err := tq.imagesService.GenerateResizedImages(&task.ImageMeta, task.DestDir)
		if err != nil {
			logging.Errorf("GenerateResizedImages() failed, path: '%s', err: %s", task.ImageMeta.Path, err)
		}
		tq.wg.Done()
	}
}

func (tq *TaskQueue) Wait() {
	tq.wg.Wait()
}

func (tq *TaskQueue) Close() {
	fmt.Println("closing task queue...")
	close(tq.tasks)
}
