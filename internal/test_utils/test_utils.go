package test_utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"receipt_uploader/internal/logging"
	"receipt_uploader/internal/models/configs"
	"receipt_uploader/internal/models/http_responses"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func CreateTestImage(filePath string, width, height int) error {
	logging.Debugf("CreateTestImage(filePath: %s, w: %d, h: %d)", filePath, width, height)

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
	dirErr := os.MkdirAll(config.ResizedDir, 0755)
	if dirErr != nil {
		err := fmt.Errorf("os.Mkdir() failed, err: %s", dirErr.Error())
		return err
	}
	logging.Debugf("folder %s has been created", config.ResizedDir)
	return nil
}

func GetFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}

	return fileInfo.Size(), nil
}

func GetImageDimension(t *testing.T, resp *http.Response) (int, int) {
	fileBytes, readErr := io.ReadAll(resp.Body)
	assert.Nil(t, readErr)

	img, _, decodeErr := image.Decode(bytes.NewReader(fileBytes))
	assert.Nil(t, decodeErr)

	return img.Bounds().Dx(), img.Bounds().Dy()
}

func ParseResponseBody(t *testing.T, resp *http.Response, response interface{}) {
	respBody, readErr := io.ReadAll(resp.Body)
	assert.Nil(t, readErr)

	unmarshalErr := json.Unmarshal(respBody, response)
	assert.Nil(t, unmarshalErr)
}

func ParseDownloadResponseHeader(resp *http.Response) (*http_responses.DownloadResponseHeader, error) {
	contentLen, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return nil, err
	}

	cp := resp.Header.Get("Content-Disposition")
	fileName := strings.TrimPrefix(cp, "attachment; filename=")

	return &http_responses.DownloadResponseHeader{
		Filename:      fileName,
		ContentType:   resp.Header.Get("Content-Type"),
		ContentLength: int64(contentLen),
	}, nil
}
