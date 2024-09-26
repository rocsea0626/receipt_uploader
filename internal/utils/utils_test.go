package utils

import (
	"os"
	"receipt_uploader/test_utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResizeImage(t *testing.T) {
	t.Run("succeed", func(t *testing.T) {
		testFilePath := "test_resize_image.png"
		createImgErr := test_utils.CreateTestImage(testFilePath)
		assert.Nil(t, createImgErr)
		defer os.Remove(testFilePath)

		width, height := uint(100), uint(100)

		resizedImg, resizeErr := ResizeImage(testFilePath, width, height)
		assert.Nil(t, resizeErr)

		bounds := resizedImg.Bounds()
		if bounds.Dx() != int(width) || bounds.Dy() != int(height) {
			t.Errorf("Resized image has unexpected dimensions: got %dx%d, want %dx%d",
				bounds.Dx(), bounds.Dy(), width, height)
		}
	})
}

func TestGetFileNameSmall(t *testing.T) {
	t.Run("succeed", func(t *testing.T) {
		fPath := "/images/test_resize_image.png"
		newPath := AppendSuffix(fPath, "small")
		assert.Equal(t, "/images/test_resize_image_small.png", newPath)
	})

	t.Run("succeed, no extension", func(t *testing.T) {
		fPath := "/images/test_resize_image"
		newPath := AppendSuffix(fPath, "medium")
		assert.Equal(t, "/images/test_resize_image_medium", newPath)
	})
}
