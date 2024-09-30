package image_worker

type ServiceType interface {
	ResizeImages(srcDir, destDir string) error
}
