package image_meta

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromPath(t *testing.T) {
	dir := "test-from-path"

	t.Run("Valid Input", func(t *testing.T) {
		fileName := "user1#123456.jpg"
		path := filepath.Join(dir, fileName)
		expected := &ImageMeta{
			Path:      path,
			Extension: ".jpg",
			Username:  "user1",
			Dir:       dir,
			ReceiptID: "123456",
			FileName:  fileName,
		}

		imgFile, err := FromUploadDir(path)

		assert.Nil(t, err)
		assert.Equal(t, expected.Path, path)
		assert.Equal(t, expected.Dir, dir)
		assert.Equal(t, expected.Username, imgFile.Username)
		assert.Equal(t, expected.Extension, imgFile.Extension)
		assert.Equal(t, expected.ReceiptID, imgFile.ReceiptID)
		assert.Equal(t, expected.FileName, imgFile.FileName)
	})

	t.Run("Invalid Input (missing username)", func(t *testing.T) {
		path := filepath.Join(dir, "123456.jpg")

		_, err := FromUploadDir(path)
		assert.NotNil(t, err)
	})
}

func TestFromUpload(t *testing.T) {
	t.Run("Valid Upload", func(t *testing.T) {
		username := "user2"
		extension := "png"
		uploadDir := "test-image-files/uploads"

		imgFile := FromFormData(username, extension, uploadDir)

		expectedPath := filepath.Join(uploadDir, username+"#"+imgFile.ReceiptID+"."+extension)

		assert.NotNil(t, imgFile)
		assert.Equal(t, username, imgFile.Username)
		assert.Equal(t, expectedPath, imgFile.Path)
		assert.True(t, strings.HasSuffix(imgFile.Path, extension))
	})
}

func TestGetResizedPath(t *testing.T) {
	baseDir := "test-get-resized-path"
	uploadDir := filepath.Join(baseDir, "uploads")
	path := filepath.Join(uploadDir, "user1#123456.jpg")
	destDir := filepath.Join(baseDir, "/resized")

	imgFile := &ImageMeta{
		Path:      path,
		Dir:       uploadDir,
		Username:  "user1",
		ReceiptID: "123456",
		Extension: ".jpg",
	}

	t.Run("With Size", func(t *testing.T) {
		size := "medium"

		expectedPath := filepath.Join(destDir, "123456_medium.jpg")
		resizedPath := GetResizedPath(imgFile, destDir, size)

		assert.Equal(t, expectedPath, resizedPath)
	})

	t.Run("Without Size", func(t *testing.T) {
		size := ""

		expectedPath := filepath.Join(destDir, "123456.jpg")
		resizedPath := GetResizedPath(imgFile, destDir, size)

		assert.Equal(t, expectedPath, resizedPath)

	})
}

func TestFromGetRequset(t *testing.T) {
	baseDir := "test-from-get-request"

	t.Run("Valid Input, without size", func(t *testing.T) {
		receiptID := "123456"
		size := ""
		username := "user1"
		srcDir := filepath.Join(baseDir, "images")

		expectedFileName := receiptID + ".jpg"
		expectedPath := filepath.Join(srcDir, username, expectedFileName)

		imgMeta := FromGetRequset(receiptID, size, username, srcDir)

		assert.NotNil(t, imgMeta)
		assert.Equal(t, expectedFileName, imgMeta.FileName)
		assert.Equal(t, expectedPath, imgMeta.Path)
		assert.Equal(t, username, imgMeta.Username)
		assert.Equal(t, receiptID, imgMeta.ReceiptID)
		assert.Equal(t, ".jpg", imgMeta.Extension)
		assert.Equal(t, filepath.Join(srcDir, username), imgMeta.Dir)
	})

	t.Run("Valid Input, with size", func(t *testing.T) {
		receiptID := "123456"
		size := "large"
		username := "user1"
		srcDir := filepath.Join(baseDir, "images")

		expectedFileName := receiptID + "_" + size + ".jpg"
		expectedPath := filepath.Join(srcDir, username, expectedFileName)

		imgMeta := FromGetRequset(receiptID, size, username, srcDir)

		assert.NotNil(t, imgMeta)
		assert.Equal(t, expectedFileName, imgMeta.FileName)
		assert.Equal(t, expectedPath, imgMeta.Path)
		assert.Equal(t, username, imgMeta.Username)
		assert.Equal(t, receiptID, imgMeta.ReceiptID)
		assert.Equal(t, ".jpg", imgMeta.Extension)
		assert.Equal(t, filepath.Join(srcDir, username), imgMeta.Dir)
	})
}
