package images_mock

import (
	"errors"
	"log"
	"net/http"
	"receipt_uploader/internal/constants"
	"receipt_uploader/internal/models/image_meta"
	"time"
)

type ServiceMock struct{}

func (s *ServiceMock) ParseImage(r *http.Request) ([]byte, error) {
	log.Println("images_mock.ParseImage()")
	return nil, nil
}

func (s *ServiceMock) GenerateResizedImages(imageMeta *image_meta.ImageMeta, destDir string) error {
	log.Printf("images_mock.GenerateResizedImages(srcPath: %s)", imageMeta.Path)

	if destDir == "mock_generate_images_failed" {
		return errors.New("mock GenerateResizedImages() failed")
	}

	if destDir == "mock_generate_images_timeout" {
		time.Sleep(constants.IMAGE_WORKER_TIMEOUT)
		return nil
	}
	return nil
}

func (s *ServiceMock) SaveUpload(bytes *[]byte, username, destDir string) (*image_meta.ImageMeta, error) {
	log.Println("images_mock.SaveUpload()")
	return nil, nil
}

func (s *ServiceMock) GetImage(imageMeta *image_meta.ImageMeta) ([]byte, string, error) {
	log.Printf("images_mock.GetImage(receiptId: %s)", imageMeta.ReceiptID)
	if imageMeta.ReceiptID == "mockgetimagefailed" {
		return nil, "", errors.New("mock GetImage() failed")
	}
	return nil, "", nil
}
