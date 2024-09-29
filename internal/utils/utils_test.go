package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateDestPath(t *testing.T) {
	outputDir := "output"

	t.Run("succeed", func(t *testing.T) {
		fPath := "/input/test_resize_image.jpg"
		newPath := GenerateDestPath(fPath, outputDir, "small")
		assert.Equal(t, "output/test_resize_image_small.jpg", newPath)
	})

	t.Run("succeed, no extension", func(t *testing.T) {
		fPath := "/input/test_resize_image"
		newPath := GenerateDestPath(fPath, outputDir, "medium")
		assert.Equal(t, "output/test_resize_image_medium", newPath)
	})

	t.Run("succeed, no path & extension", func(t *testing.T) {
		fPath := "test_resize_image"
		newPath := GenerateDestPath(fPath, outputDir, "large")
		assert.Equal(t, "output/test_resize_image_large", newPath)
	})
}
