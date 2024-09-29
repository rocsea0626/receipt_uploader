package utils

import (
	"os"
	"receipt_uploader/internal/test_utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveUploadImage(t *testing.T) {
	userToken := ""

	t.Run("succeed", func(t *testing.T) {
		fileName := "test_image_save_upload.jpg"
		uploadsDir := "uploads"
		err := os.MkdirAll(uploadsDir, 0755)
		defer os.RemoveAll(uploadsDir)
		assert.NoError(t, err)

		createErr := test_utils.CreateTestImage(fileName, 1200, 1200)
		assert.Nil(t, createErr)
		defer os.Remove(fileName)

		req, reqErr := test_utils.GenerateUploadRequest(t, "/upload", fileName, userToken)
		assert.Nil(t, reqErr)

		tmpPath, err := SaveUploadedImage(req, uploadsDir)
		assert.NoError(t, err)
		assert.NotEmpty(t, tmpPath)
	})

	t.Run("should fail, invalid image width, w=500, h=1000", func(t *testing.T) {
		fileName := "test_image_save_upload.jpg"
		tmpDir := "tmp"
		err := os.Mkdir(tmpDir, 0755)
		defer os.RemoveAll(tmpDir)
		assert.NoError(t, err)

		createErr := test_utils.CreateTestImage(fileName, 500, 1000)
		assert.Nil(t, createErr)
		defer os.Remove(fileName)

		req, reqErr := test_utils.GenerateUploadRequest(t, "/upload", fileName, userToken)
		assert.Nil(t, reqErr)

		tmpPath, err := SaveUploadedImage(req, tmpDir)
		assert.NotNil(t, err)
		assert.Empty(t, tmpPath)
	})

	t.Run("should fail, invalid image height, w=600, h=500", func(t *testing.T) {
		fileName := "test_image_save_upload.jpg"
		tmpDir := "tmp"
		err := os.Mkdir(tmpDir, 0755)
		defer os.RemoveAll(tmpDir)
		assert.NoError(t, err)

		createErr := test_utils.CreateTestImage(fileName, 500, 1000)
		assert.Nil(t, createErr)
		defer os.Remove(fileName)

		req, reqErr := test_utils.GenerateUploadRequest(t, "/upload", fileName, userToken)
		assert.Nil(t, reqErr)

		tmpPath, err := SaveUploadedImage(req, tmpDir)
		assert.NotNil(t, err)
		assert.Empty(t, tmpPath)
	})
}

func TestGetOutputPath(t *testing.T) {
	outputDir := "output"

	t.Run("succeed", func(t *testing.T) {
		fPath := "/input/test_resize_image.jpg"
		newPath := GetOutputPath(fPath, outputDir, "small")
		assert.Equal(t, "output/test_resize_image_small.jpg", newPath)
	})

	t.Run("succeed, no extension", func(t *testing.T) {
		fPath := "/input/test_resize_image"
		newPath := GetOutputPath(fPath, outputDir, "medium")
		assert.Equal(t, "output/test_resize_image_medium", newPath)
	})

	t.Run("succeed, no path & extension", func(t *testing.T) {
		fPath := "test_resize_image"
		newPath := GetOutputPath(fPath, outputDir, "large")
		assert.Equal(t, "output/test_resize_image_large", newPath)
	})
}
