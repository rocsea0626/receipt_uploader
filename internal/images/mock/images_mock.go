package images_mock

import (
	"errors"
	"log"
	"net/http"
)

type ServiceMock struct{}

func (s *ServiceMock) ParseImage(r *http.Request) ([]byte, error) {
	log.Println("images_mock.ParseImage()")
	return nil, nil
}

func (s *ServiceMock) GenerateResizedImages(srcPath, destDir string) error {
	log.Printf("images_mock.GenerateResizedImages(srcPath: %s)", srcPath)
	if destDir == "mock_generate_images_failed" {
		return errors.New("mock GenerateResizedImages() failed")
	}
	return nil
}

func (s *ServiceMock) SaveUpload(bytes *[]byte, destDir string) (string, error) {
	log.Println("images_mock.SaveUpload()")
	return "", nil
}

func (s *ServiceMock) GetImage(receiptId, size, srcDir string) ([]byte, string, error) {
	log.Printf("images_mock.GetImage(receiptId: %s, size: %s)", receiptId, size)
	if receiptId == "mockgetimagefailed" {
		return nil, "", errors.New("mock GetImage() failed")
	}
	return nil, "", nil
}
