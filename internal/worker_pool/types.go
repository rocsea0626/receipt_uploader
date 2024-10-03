package worker_pool

type ServiceType interface {
	Start(stopChan <-chan struct{})
	Submit(task Task) bool
}
