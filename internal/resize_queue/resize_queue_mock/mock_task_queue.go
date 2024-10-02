package resize_queue_mock

import (
	"receipt_uploader/internal/logging"
	"receipt_uploader/internal/models/tasks"
)

type ServiceMock struct{}

func (tq *ServiceMock) Start(stopChan <-chan struct{}) {
	logging.Debugf("resize_queue_mock.Start()")
}

func (tq *ServiceMock) Enqueue(task tasks.ResizeTask) bool {
	logging.Debugf("resize_queue_mock.Enqueue(taks: #%v)", task)

	if task.DestDir == "./test_image_enqueue_failed" {
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
