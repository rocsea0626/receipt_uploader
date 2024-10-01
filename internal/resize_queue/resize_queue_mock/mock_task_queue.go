package resize_queue_mock

import (
	"fmt"
	"receipt_uploader/internal/models/tasks"
)

type ServiceMock struct{}

func (tq *ServiceMock) Start(stopChan <-chan struct{}) {
	fmt.Println("resize_queue_mock.Start()")
}

func (tq *ServiceMock) Enqueue(task tasks.ResizeTask) bool {
	fmt.Println("resize_queue_mock.Enqueue()")
	return true
}

func (q *ServiceMock) Process() {
	fmt.Println("resize_queue_mock.Enqueue()")
}

func (q *ServiceMock) Wait() {
	fmt.Println("resize_queue_mock.Wait()")

}

func (q *ServiceMock) Close() {
	fmt.Println("resize_queue_mock.Close()")
}
