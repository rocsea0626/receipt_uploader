package test_utils

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"receipt_uploader/internal/models/configs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func CreateTestImage(filePath string, width, height int) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
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

// GenerateUploadRequest prepares and uploads a test image to a specified URL
func GenerateUploadRequest(t *testing.T, url string, fileName string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("receipt", fileName)
	if err != nil {
		return nil, fmt.Errorf("error creating form file: %w", err)
	}

	tempFile, openErr := os.Open(fileName)
	assert.Nil(t, openErr)
	defer tempFile.Close()

	if _, err := io.Copy(part, tempFile); err != nil {
		return nil, fmt.Errorf("error writing to form file: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("error closing writer: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, nil
}

func InitTestServer(config *configs.Config) error {
	tmpErr := os.MkdirAll(config.DIR_TMP, 0755)
	if tmpErr != nil {
		err := fmt.Errorf("os.Mkdir() failed, err: %s", tmpErr.Error())
		return err
	}
	log.Printf("folder %s has been created", config.DIR_TMP)

	imagesErr := os.MkdirAll(config.DIR_IMAGES, 0755)
	if imagesErr != nil {
		err := fmt.Errorf("os.Mkdir() failed, err: %s", imagesErr.Error())
		return err
	}
	log.Printf("folder %s has been created", config.DIR_IMAGES)
	return nil
}
