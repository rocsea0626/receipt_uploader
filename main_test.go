package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"receipt_uploader/constants"
	"receipt_uploader/internal/models/configs"
	"receipt_uploader/internal/test_utils"
	"receipt_uploader/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	config := &configs.Config{
		Port:         ":8080",
		UploadedDir:  filepath.Join(constants.ROOT_DIR_IMAGES, "integ-test-tmp"),
		GeneratedDir: filepath.Join(constants.ROOT_DIR_IMAGES, "integ-test-images"),
	}
	baseUrl := "http://localhost" + config.Port

	defer os.RemoveAll(config.UploadedDir)
	defer os.RemoveAll(config.GeneratedDir)

	go utils.StartServer(config)

	t.Run("return 200, /health", func(t *testing.T) {
		resp, err := http.Get(baseUrl + "/health")
		assert.Nil(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("return 201, /receipts", func(t *testing.T) {
		filePath := "./integ-test.jpg"
		url := baseUrl + "/receipts"

		test_utils.CreateTestImage(filePath, 1000, 1200)
		defer os.Remove(filePath)
		req, reqErr := test_utils.GenerateUploadRequest(t, url, filePath)
		assert.Nil(t, reqErr)

		client := &http.Client{}
		resp, err := client.Do(req)
		assert.Nil(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		respBody, readErr := io.ReadAll(resp.Body)
		assert.Nil(t, readErr)
		assert.Contains(t, string(respBody), "receiptId")
	})
}
