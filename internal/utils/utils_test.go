package utils

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"receipt_uploader/internal/futils"
	"receipt_uploader/internal/test_utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResizeImage(t *testing.T) {
	t.Run("succeed", func(t *testing.T) {
		testFilePath := "test_resize_image.jpg"
		createImgErr := test_utils.CreateTestImage(testFilePath)
		assert.Nil(t, createImgErr)
		defer os.Remove(testFilePath)

		width, height := uint(100), uint(100)

		resizedImg, resizeErr := resizeImage(testFilePath, width, height)
		assert.Nil(t, resizeErr)

		bounds := resizedImg.Bounds()
		if bounds.Dx() != int(width) || bounds.Dy() != int(height) {
			t.Errorf("Resized image has unexpected dimensions: got %dx%d, want %dx%d",
				bounds.Dx(), bounds.Dy(), width, height)
		}
	})
}

func TestGenerateImages(t *testing.T) {
	inputPath := "./test.jpg"
	outputDir := "./output/"
	defer os.RemoveAll(outputDir)
	os.MkdirAll(outputDir, 0755)

	t.Run("succeed", func(t *testing.T) {
		createErr := test_utils.CreateTestImage(inputPath)
		assert.Nil(t, createErr)

		genErr := GenerateImages(inputPath, outputDir)
		assert.Nil(t, genErr)

		smallImagePath := futils.GetOutputPath(inputPath, outputDir, "small")
		mediumImagePath := futils.GetOutputPath(inputPath, outputDir, "medium")
		largeImagePath := futils.GetOutputPath(inputPath, outputDir, "large")

		_, smallErr := os.Stat(smallImagePath)
		assert.Nil(t, smallErr)
		_, mediumErr := os.Stat(mediumImagePath)
		assert.Nil(t, mediumErr)
		_, largeErr := os.Stat(largeImagePath)
		assert.Nil(t, largeErr)

		os.Remove(inputPath)
	})
}

func TestSaveUploadImage(t *testing.T) {
	t.Run("succeed", func(t *testing.T) {
		fileName := "test_image_save_upload.jpg"
		tmpDir := "tmp"
		err := os.Mkdir(tmpDir, 0755)
		defer os.RemoveAll(tmpDir)
		assert.NoError(t, err)

		body, writer := test_utils.CreateImageForUpload(t, fileName, 300, 200)
		defer os.Remove(fileName)
		defer writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		tmpPath, err := SaveUploadedImage(req, tmpDir)
		assert.NoError(t, err)
		assert.NotEmpty(t, tmpPath)
		log.Printf("tmpPath: %s", tmpPath)
	})
}
