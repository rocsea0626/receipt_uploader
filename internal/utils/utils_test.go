package utils

import (
	"log"
	"os"
	"receipt_uploader/internal/test_utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveUploadImage(t *testing.T) {
	t.Run("succeed", func(t *testing.T) {
		fileName := "test_image_save_upload.jpg"
		tmpDir := "tmp"
		err := os.Mkdir(tmpDir, 0755)
		defer os.RemoveAll(tmpDir)
		assert.NoError(t, err)

		createErr := test_utils.CreateTestImage(fileName, 1200, 1200)
		assert.Nil(t, createErr)
		defer os.Remove(fileName)

		req, reqErr := test_utils.GenerateUploadRequest(t, "/upload", fileName)
		assert.Nil(t, reqErr)

		tmpPath, err := SaveUploadedImage(req, tmpDir)
		assert.NoError(t, err)
		assert.NotEmpty(t, tmpPath)
		log.Printf("tmpPath: %s", tmpPath)
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

		req, reqErr := test_utils.GenerateUploadRequest(t, "/upload", fileName)
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

		req, reqErr := test_utils.GenerateUploadRequest(t, "/upload", fileName)
		assert.Nil(t, reqErr)

		tmpPath, err := SaveUploadedImage(req, tmpDir)
		assert.NotNil(t, err)
		assert.Empty(t, tmpPath)
	})
}
