package worker_pool_mock

import (
	"receipt_uploader/internal/logging"
	"receipt_uploader/internal/worker_pool"
	"strings"
)

type ServiceMock struct {
	Capacity int
}

func (tq *ServiceMock) Start(stopChan <-chan struct{}) {
	logging.Debugf("resize_queue_mock.Start()")
}

func (tq *ServiceMock) Submit(task worker_pool.Task) bool {
	logging.Debugf("resize_queue_mock.Submit(taks: #%v)", task)
	if strings.Contains(task.Name, "mock_submit_failed") {
		logging.Debugf("mock submit failed")
		return false
	}
	return true
}
