package image_worker

import (
	"log"
	"os"
	"path/filepath"
	"receipt_uploader/internal/images"
	images_mock "receipt_uploader/internal/images/mock"
	"receipt_uploader/internal/logging"
	"receipt_uploader/internal/models/configs"
	"receipt_uploader/internal/models/image_meta"
	"receipt_uploader/internal/test_utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResizeImages(t *testing.T) {
	t.Run("succeed", func(t *testing.T) {
		imageService := images.NewService(&configs.AllowedDimensions)
		service := &Service{
			ImageService: imageService,
		}

		baseDir := "image_worker"
		uploadsDir := filepath.Join(baseDir, "uploads")
		destDir := filepath.Join(baseDir, "resized")
		os.MkdirAll(uploadsDir, 0755)
		os.MkdirAll(destDir, 0755)
		defer os.RemoveAll(baseDir)

		username := "user_1"
		extension := "jpg"
		testFilename := username + "#receiptupload123.jpg"
		testFilePath := filepath.Join(uploadsDir, testFilename)
		imageFile := image_meta.FromFormData(username, extension, uploadsDir)
		test_utils.CreateTestImage(imageFile.Path, 1000, 1200)

		resizeErr := service.ResizeImages(uploadsDir, destDir)
		assert.Nil(t, resizeErr)

		smallImagePath := image_meta.GetResizedPath(imageFile, filepath.Join(destDir, username), "small")
		mediumImagePath := image_meta.GetResizedPath(imageFile, filepath.Join(destDir, username), "medium")
		largeImagePath := image_meta.GetResizedPath(imageFile, filepath.Join(destDir, username), "large")

		_, smallErr := os.Stat(smallImagePath)
		assert.Nil(t, smallErr)
		log.Println("smallErr: ", smallErr)
		_, mediumErr := os.Stat(mediumImagePath)
		assert.Nil(t, mediumErr)
		_, largeErr := os.Stat(largeImagePath)
		assert.Nil(t, largeErr)

		_, err := os.Stat(testFilePath)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("should fail, GenerateResizedImages() failed", func(t *testing.T) {
		mockImagesService := &images_mock.ServiceMock{}
		service := &Service{
			ImageService: mockImagesService,
		}

		baseDir := "image_worker"
		uploadsDir := filepath.Join(baseDir, "uploads")
		destDir := "mock_generate_images_failed"
		os.MkdirAll(uploadsDir, 0755)
		os.MkdirAll(destDir, 0755)
		defer os.RemoveAll(baseDir)

		username := "user_1"
		extension := "jpg"
		imageFile := image_meta.FromFormData(username, extension, uploadsDir)
		test_utils.CreateTestImage(imageFile.Path, 1000, 1200)

		resizeErr := service.ResizeImages(uploadsDir, destDir)
		assert.NotNil(t, resizeErr)
		logging.Debugf("resizeErr: %s", resizeErr.Error())

		_, fileErr := os.Stat(imageFile.Path)
		assert.Nil(t, fileErr)
	})
}
