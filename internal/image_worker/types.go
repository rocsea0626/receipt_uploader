package image_worker

type ServiceType interface {
	Start(stopChan <-chan struct{})
}
