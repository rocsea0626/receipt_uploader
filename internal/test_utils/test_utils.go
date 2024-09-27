package test_utils

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"mime/multipart"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func CreateTestImage(filePath string) error {
	img := image.NewRGBA(image.Rect(0, 0, 800, 600))
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			img.Set(x, y, color.RGBA{uint8(x % 256), uint8(y % 256), 0, 255})
		}
	}
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()
	return jpeg.Encode(out, img, nil)
}

// CreateImageForUpload will create a temporary JPEG image file and return
// a multipart form file for use in tests.
func CreateImageForUpload(t *testing.T, fileName string, width, height int) (*bytes.Buffer, *multipart.Writer) {
	tempFile, createErr := os.Create(fileName)
	assert.Nil(t, createErr)
	defer tempFile.Close()

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.RGBA{uint8(x * y % 256), 0, 0, 255})
		}
	}

	encodeErr := jpeg.Encode(tempFile, img, nil)
	assert.Nil(t, encodeErr)

	tempFile, openErr := os.Open(fileName)
	assert.Nil(t, openErr)
	defer tempFile.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, writerErr := writer.CreateFormFile("receipt", fileName)
	assert.Nil(t, writerErr)

	_, err := part.Write([]byte("fake content"))
	assert.Nil(t, err)

	writer.Close()

	return body, writer
}
