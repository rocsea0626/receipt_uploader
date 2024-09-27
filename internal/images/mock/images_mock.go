package images_mock

import (
	"errors"
	"log"
	"net/http"
	"receipt_uploader/internal/models/configs"
)

type ServiceMock struct {
	Config *configs.Config
}

func (s *ServiceMock) DecodeImage(r *http.Request) ([]byte, error) {
	log.Println("images_mock.DecodeImage()")
	return nil, nil
}

func (s *ServiceMock) GenerateImages(srcPath string) error {
	log.Printf("images_mock.GenerateImages(srcPath: %s)", srcPath)
	if s.Config.DIR_IMAGES == "mock_generate_images_failed" {
		return errors.New("mock GenerateImages() failed")
	}
	return nil
}

func (s *ServiceMock) SaveUpload(bytes []byte) (string, error) {
	log.Println("images_mock.SaveUpload()")
	return "", nil
}
