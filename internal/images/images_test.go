package images

import (
	"bytes"
	"image"
	"os"
	"path/filepath"
	"receipt_uploader/internal/models/configs"
	"receipt_uploader/internal/models/image_meta"
	"receipt_uploader/internal/test_utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateImages(t *testing.T) {
	baseDir := "test-gen-images"
	uploadDir := filepath.Join(baseDir, "uploads")
	username := "user1"
	filaname := username + "#test.jpg"
	srcPath := filepath.Join(uploadDir, filaname)
	destDir := filepath.Join(baseDir, "resized")

	os.MkdirAll(uploadDir, 0755)
	os.MkdirAll(destDir, 0755)
	defer os.RemoveAll(baseDir)

	service := NewService(&configs.AllowedDimensions)

	t.Run("succeed", func(t *testing.T) {
		createErr := test_utils.CreateTestImageJPG(srcPath, 800, 1200)
		assert.Nil(t, createErr)

		imageMeta, imageErr := image_meta.FromUploadDir(srcPath)
		assert.Nil(t, imageErr)

		genErr := service.GenerateResizedImages(imageMeta, destDir)
		assert.Nil(t, genErr)
		imgFile, imgErr := image_meta.FromUploadDir(srcPath)
		assert.Nil(t, imgErr)

		smallImagePath := image_meta.GetResizedPath(imgFile, filepath.Join(destDir, username), "small")
		mediumImagePath := image_meta.GetResizedPath(imgFile, filepath.Join(destDir, username), "medium")
		largeImagePath := image_meta.GetResizedPath(imgFile, filepath.Join(destDir, username), "large")

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

		createImgErr := test_utils.CreateTestImageJPG(testFilePath, orgWidth, orgHeight)
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
	username := "test_user"
	baseDir := "mock-get-images"
	srcDir := filepath.Join(baseDir, username)

	os.MkdirAll(srcDir, 0755)
	defer os.RemoveAll(baseDir)

	service := NewService(&configs.AllowedDimensions)

	t.Run("succeed, no size", func(t *testing.T) {
		receiptId := "receiptId1"
		fileName := receiptId + ".jpg"
		fPath := filepath.Join(srcDir, fileName)

		createImgErr := test_utils.CreateTestImageJPG(fPath, 100, 100)
		assert.Nil(t, createImgErr)

		imageMeta := &image_meta.ImageMeta{
			Path:     fPath,
			FileName: fileName,
		}

		fileBytes, fName, getErr := service.GetImage(imageMeta)
		assert.Nil(t, getErr)
		assert.Greater(t, len(fileBytes), 0)
		assert.Equal(t, fileName, fName)

	})

	t.Run("succeed, size=small", func(t *testing.T) {
		receiptId := "receiptId1"
		size := "small"
		fileName := receiptId + "_" + size + ".jpg"
		fPath := filepath.Join(srcDir, fileName)

		createImgErr := test_utils.CreateTestImageJPG(fPath, 100, 100)
		assert.Nil(t, createImgErr)

		imageMeta := &image_meta.ImageMeta{
			Path:     fPath,
			FileName: fileName,
		}

		fileBytes, fName, getErr := service.GetImage(imageMeta)
		assert.Nil(t, getErr)
		assert.Greater(t, len(fileBytes), 0)
		assert.Equal(t, fileName, fName)

	})

	t.Run("succeed, size=large", func(t *testing.T) {
		receiptId := "receiptId1"
		size := "large"
		fileName := receiptId + "_" + size + ".jpg"
		fPath := filepath.Join(srcDir, fileName)

		createImgErr := test_utils.CreateTestImageJPG(fPath, 100, 100)
		assert.Nil(t, createImgErr)

		imageMeta := &image_meta.ImageMeta{
			Path:     fPath,
			FileName: fileName,
		}

		fileBytes, fName, getErr := service.GetImage(imageMeta)
		assert.Nil(t, getErr)
		assert.Greater(t, len(fileBytes), 0)
		assert.Equal(t, fileName, fName)

	})

	t.Run("should fail, non existing file", func(t *testing.T) {
		receiptId := "non-existing"
		size := "large"
		fileName := receiptId + "_" + size + ".jpg"
		fPath := filepath.Join(srcDir, fileName)

		imageMeta := &image_meta.ImageMeta{
			Path:     fPath,
			FileName: fileName,
		}

		fileBytes, fName, getErr := service.GetImage(imageMeta)
		assert.NotNil(t, getErr)
		assert.ErrorIs(t, getErr, os.ErrNotExist)
		assert.Nil(t, fileBytes)
		assert.Equal(t, "", fName)

	})
}
