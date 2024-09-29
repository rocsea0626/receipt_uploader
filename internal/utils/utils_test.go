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
