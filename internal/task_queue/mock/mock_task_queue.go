package task_queue_mock

import (
	"fmt"
	"receipt_uploader/internal/models/tasks"
)

type ServiceMock struct{}

func (tq *ServiceMock) Start(stopChan <-chan struct{}) {
	fmt.Println("task_queue_mock.Start()")
}

func (tq *ServiceMock) Enqueue(task tasks.ResizeTask) bool {
	fmt.Println("task_queue_mock.Enqueue()")
	return true
}
