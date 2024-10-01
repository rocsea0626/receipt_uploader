package resize_queue

import "receipt_uploader/internal/models/tasks"

type ServiceType interface {
	Start(stopChan <-chan struct{})
	Enqueue(task tasks.ResizeTask) bool
}
