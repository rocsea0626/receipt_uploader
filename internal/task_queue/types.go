package task_queue

type ServiceType interface {
	Start(stopChan <-chan struct{})
	Enqueue(task Task) bool
	Process()
	Wait()
	Close()
}
