package images_mock

import (
	"errors"
	"log"
	"net/http"
	"receipt_uploader/internal/models/image_meta"
)

type ServiceMock struct{}

func (s *ServiceMock) ParseImage(r *http.Request) ([]byte, error) {
	log.Println("images_mock.ParseImage()")
	return nil, nil
}

func (s *ServiceMock) GenerateResizedImages(imageFile *image_meta.ImageMeta, destDir string) error {
	log.Printf("images_mock.GenerateResizedImages(srcPath: %s)", imageFile.Path)

	if destDir == "mock_generate_images_failed" {
		return errors.New("mock GenerateResizedImages() failed")
	}
	return nil
}

func (s *ServiceMock) SaveUpload(bytes *[]byte, username, destDir string) (*image_meta.ImageMeta, error) {
	log.Println("images_mock.SaveUpload()")
	return nil, nil
}

func (s *ServiceMock) GetImage(imageFile *image_meta.ImageMeta) ([]byte, string, error) {
	log.Printf("images_mock.GetImage(receiptId: %s)", imageFile.ReceiptID)
	if imageFile.ReceiptID == "mockgetimagefailed" {
		return nil, "", errors.New("mock GetImage() failed")
	}
	return nil, "", nil
}
