package images_mock

import (
	"errors"
	"log"
	"net/http"
)

type ServiceMock struct{}

// DecodeImage implements images.ServiceType.
func (s *ServiceMock) DecodeImage(r *http.Request) ([]byte, error) {
	log.Println("images_mock.DecodeImage()")
	return nil, nil
}

func (s *ServiceMock) GenerateImages(srcPath string, destDir string) error {
	log.Printf("images_mock.GenerateImages(srcPath: %s, destDir: %s)", srcPath, destDir)
	if destDir == "mock_generate_images_failed" {
		return errors.New("mock GenerateImages() failed")
	}
	return nil
}

// SaveUpload implements images.ServiceType.
func (s *ServiceMock) SaveUpload(bytes []byte, tmpDir string) (string, error) {
	log.Println("images_mock.SaveUpload()")
	return "", nil
}
