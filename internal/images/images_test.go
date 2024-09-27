package images

import (
	"os"
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
}
