package images

import (
	"bytes"
	"image"
	"os"
	"path/filepath"
	"receipt_uploader/internal/models/configs"
	"receipt_uploader/internal/test_utils"
	"receipt_uploader/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateImages(t *testing.T) {
	srcPath := "./test.jpg"
	destDir := "./output/"

	os.MkdirAll(destDir, 0755)
	defer os.RemoveAll(destDir)

	service := NewService(&configs.AllowedDimensions)

	t.Run("succeed", func(t *testing.T) {
		createErr := test_utils.CreateTestImage(srcPath, 800, 1200)
		assert.Nil(t, createErr)

		fileBytes, readErr := os.ReadFile(srcPath)
		assert.Nil(t, readErr)

		genErr := service.GenerateResizedImages(&fileBytes, srcPath, destDir)
		assert.Nil(t, genErr)

		smallImagePath := utils.GenerateDestPath(srcPath, destDir, "small")
		mediumImagePath := utils.GenerateDestPath(srcPath, destDir, "medium")
		largeImagePath := utils.GenerateDestPath(srcPath, destDir, "large")

		_, smallErr := os.Stat(smallImagePath)
		assert.Nil(t, smallErr)
		_, mediumErr := os.Stat(mediumImagePath)
		assert.Nil(t, mediumErr)
		_, largeErr := os.Stat(largeImagePath)
		assert.Nil(t, largeErr)

		os.Remove(srcPath)
	})
}

func TestResizeImage(t *testing.T) {

	t.Run("succeed", func(t *testing.T) {
		testFilePath := "test_resize_image.jpg"
		orgWidth := 800
		orgHeight := 1200

		createImgErr := test_utils.CreateTestImage(testFilePath, orgWidth, orgHeight)
		assert.Nil(t, createImgErr)
		defer os.Remove(testFilePath)

		height := orgHeight / 5
		width := orgWidth / 5

		fileBytes, readErr := os.ReadFile(testFilePath)
		assert.Nil(t, readErr)

		img, _, decodeErr := image.Decode(bytes.NewReader(fileBytes))
		assert.Nil(t, decodeErr)

		resizedBytes, resizeErr := resizeImage(&img, 0, height) // use 0 for width for keep the original ratio of image
		assert.Nil(t, resizeErr)

		reader := bytes.NewReader(resizedBytes)
		img, _, err := image.Decode(reader)
		assert.Nil(t, err)
		bounds := img.Bounds()
		assert.Equal(t, height, bounds.Dy())
		assert.Equal(t, width, bounds.Dx())
	})
}

func TestGetImage(t *testing.T) {
	srcDir := "./mock-get-images"

	os.MkdirAll(srcDir, 0755)
	defer os.RemoveAll(srcDir)

	service := NewService(&configs.AllowedDimensions)

	t.Run("succeed, size=small", func(t *testing.T) {
		receiptId := "test1get2image"
		size := "small"
		fileName := receiptId + "_" + size + ".jpg"
		fPath := filepath.Join(srcDir, fileName)

		createImgErr := test_utils.CreateTestImage(fPath, 100, 100)
		assert.Nil(t, createImgErr)

		fileBytes, fName, getErr := service.GetImage(receiptId, size, srcDir)
		assert.Nil(t, getErr)
		assert.Greater(t, len(fileBytes), 0)
		assert.Equal(t, fileName, fName)

	})

	t.Run("succeed, size=large", func(t *testing.T) {
		receiptId := "test-get-image"
		size := "large"
		fileName := receiptId + "_" + size + ".jpg"
		fPath := filepath.Join(srcDir, fileName)

		createImgErr := test_utils.CreateTestImage(fPath, 100, 100)
		assert.Nil(t, createImgErr)

		fileBytes, fName, getErr := service.GetImage(receiptId, size, srcDir)
		assert.Nil(t, getErr)
		assert.Greater(t, len(fileBytes), 0)
		assert.Equal(t, fileName, fName)

	})

	t.Run("should fail, non existing file", func(t *testing.T) {
		receiptId := "test-get-image-non-existing"
		size := "large"

		fileBytes, fName, getErr := service.GetImage(receiptId, size, srcDir)
		assert.NotNil(t, getErr)
		assert.ErrorIs(t, getErr, os.ErrNotExist)
		assert.Nil(t, fileBytes)
		assert.Equal(t, "", fName)

	})
}
