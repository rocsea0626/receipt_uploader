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

func GenerateUploadRequest(t *testing.T, url string, fileName, userToken string) (*http.Request, error) {
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
	req.Header.Set("username_token", userToken)

	return req, nil
}

func InitTestServer(config *configs.Config) error {
	tmpErr := os.MkdirAll(config.UploadedDir, 0755)
	if tmpErr != nil {
		err := fmt.Errorf("os.Mkdir() failed, err: %s", tmpErr.Error())
		return err
	}
	log.Printf("folder %s has been created", config.UploadedDir)

	imagesErr := os.MkdirAll(config.GeneratedDir, 0755)
	if imagesErr != nil {
		err := fmt.Errorf("os.Mkdir() failed, err: %s", imagesErr.Error())
		return err
	}
	log.Printf("folder %s has been created", config.GeneratedDir)
	return nil
}

func GetFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}

	return fileInfo.Size(), nil
}
