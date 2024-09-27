package images

import (
	"os"
	"path/filepath"
	"receipt_uploader/internal/futils"
	"receipt_uploader/internal/test_utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateImages(t *testing.T) {
	inputPath := "./test.jpg"
	outputDir := "./output/"
	os.MkdirAll(outputDir, 0755)
	defer os.RemoveAll(outputDir)

	t.Run("succeed", func(t *testing.T) {
		createErr := test_utils.CreateTestImage(inputPath, 800, 1200)
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

		resizedImg, resizeErr := resizeImage(testFilePath, 0, height) // use 0 for width for keep the original ratio of image
		assert.Nil(t, resizeErr)

		bounds := resizedImg.Bounds()
		assert.Equal(t, height, bounds.Dy())
		assert.Equal(t, width, bounds.Dx())
	})

	t.Run("should fail, non existing file", func(t *testing.T) {
		testFilePath := "non-existing.jpg"

		resizedImg, resizeErr := resizeImage(testFilePath, 0, 100)
		assert.Nil(t, resizedImg)
		assert.NotNil(t, resizeErr)
		assert.Contains(t, resizeErr.Error(), "os.Open() failed:")
	})

	t.Run("should fail, invalid image data content", func(t *testing.T) {
		fPath := t.TempDir() + "/invalid.jpg"
		invalidContent := []byte("this is not an image")
		writeErr := os.WriteFile(fPath, invalidContent, 0644)
		assert.Nil(t, writeErr)

		resizedImg, err := resizeImage(fPath, 100, 100)
		assert.Nil(t, resizedImg)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "image.Decode() failed:")
	})
}

func TestSaveImage(t *testing.T) {
	t.Run("succeed", func(t *testing.T) {
		tempDir := t.TempDir()
		fPath := filepath.Join(tempDir, "test.jpg")
		destPath := filepath.Join(tempDir, "saved-test.jpg")

		createImgErr := test_utils.CreateTestImage(fPath, 100, 100)
		assert.Nil(t, createImgErr)

		resizedImg, resizeErr := resizeImage(fPath, 0, 100)
		assert.Nil(t, resizeErr)

		saveErr := saveImage(&resizedImg, destPath)
		assert.Nil(t, saveErr)

		assert.FileExists(t, destPath)

		info, _ := os.Stat(destPath)
		assert.Greater(t, info.Size(), int64(0))
	})

	t.Run("should fail, non-exising dir", func(t *testing.T) {
		tempDir := t.TempDir()
		fPath := filepath.Join(tempDir, "test.jpg")
		destPath := filepath.Join("non-existing-path/", "saved-test.jpg")

		createImgErr := test_utils.CreateTestImage(fPath, 100, 100)
		assert.Nil(t, createImgErr)

		resizedImg, resizeErr := resizeImage(fPath, 0, 100)
		assert.Nil(t, resizeErr)

		saveErr := saveImage(&resizedImg, destPath)
		assert.NotNil(t, saveErr)
		assert.Contains(t, saveErr.Error(), "os.Create() failed, err:")
	})
}
