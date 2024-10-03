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
