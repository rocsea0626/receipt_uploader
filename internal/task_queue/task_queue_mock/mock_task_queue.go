package task_queue_mock

import (
	"receipt_uploader/internal/logging"
	"receipt_uploader/internal/task_queue"
	"strings"
)

type ServiceMock struct {
	Capacity int
}

func (tq *ServiceMock) Start(stopChan <-chan struct{}) {
	logging.Debugf("resize_queue_mock.Start()")
}

func (tq *ServiceMock) Enqueue(task task_queue.Task) bool {
	logging.Debugf("resize_queue_mock.Enqueue(taks: #%v)", task)
	if strings.Contains(task.Name, "mock_enqueue_timeout") {
		logging.Debugf("mock enqueue timeout")
		return false
	}
	return true
}

func (q *ServiceMock) Process() {
	logging.Debugf("resize_queue_mock.Enqueue()")
}

func (q *ServiceMock) Wait() {
	logging.Debugf("resize_queue_mock.Wait()")

}

func (q *ServiceMock) Close() {
	logging.Debugf("resize_queue_mock.Close()")
}
