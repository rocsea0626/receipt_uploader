package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"receipt_uploader/internal/logging"
	"receipt_uploader/internal/models/configs"
	"receipt_uploader/internal/models/http_responses"
	"receipt_uploader/internal/test_utils"
	"receipt_uploader/internal/utils"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var numClients = 15

func TestMainStess(t *testing.T) {

	baseDir := "stress-test-images"
	config := &configs.Config{
		Port:       ":8080",
		ResizedDir: filepath.Join(baseDir, "resized"),
		UploadsDir: filepath.Join(baseDir, "uploads"),
		Dimensions: configs.AllowedDimensions,
		Interval:   time.Duration(1) * time.Second,
		Mode:       "release",
	}
	baseUrl := "http://localhost" + config.Port
	defer os.RemoveAll(baseDir)

	client := &http.Client{}

	stopChan := make(chan struct{})
	t.Cleanup(func() {
		log.Println("Cleanup stress test")
		close(stopChan)
	})

	go utils.StartServer(config, stopChan)

	t.Run("stress testing, multiple POST /receipts request", func(t *testing.T) {
		waitTime := numClients * 2
		var wg sync.WaitGroup
		url := baseUrl + "/receipts"

		receiptIDs := make(map[string]string)

		for i := 0; i < numClients; i++ {
			wg.Add(1)
			go func(clientID int) {
				defer wg.Done()

				userToken := "user_token_" + strconv.Itoa(i)
				req, reqErr := test_utils.GenerateUploadRequest(t, url, "test_image.jpg", userToken)
				assert.Nil(t, reqErr)

				resp, err := client.Do(req)
				assert.Nil(t, err)
				defer resp.Body.Close()

				assert.Equal(t, http.StatusCreated, resp.StatusCode)

				var uploadResp http_responses.UploadResponse
				test_utils.ParseResponseBody(t, resp, &uploadResp)
				assert.NotEmpty(t, uploadResp.ReceiptID)
				receiptIDs[userToken] = uploadResp.ReceiptID

			}(i)
		}
		wg.Wait()

		logging.Infof("waiting %d s for server completes resizing of all uploaded images...", waitTime)
		time.Sleep(time.Duration(waitTime) * time.Second)

		size := "small"
		for token, receiptID := range receiptIDs {
			wg.Add(1)

			go func(token, receiptID string) {
				defer wg.Done()

				getUrl := fmt.Sprintf("%s/%s?size=%s", url, receiptID, size)

				getReq, getReqErr := http.NewRequest(http.MethodGet, getUrl, nil)
				getReq.Header.Set("username_token", token)
				assert.Nil(t, getReqErr)

				getResp, getErr := client.Do(getReq)
				assert.Nil(t, getErr)
				defer getResp.Body.Close()

				assert.Equal(t, http.StatusOK, getResp.StatusCode)
				header, headerErr := test_utils.ParseDownloadResponseHeader(getResp)
				assert.Nil(t, headerErr)
				assert.Equal(t, "image/jpeg", header.ContentType)

				fileName := receiptID + "_" + size + ".jpg"
				assert.Equal(t, fileName, header.Filename)

				_, height := test_utils.GetImageDimension(t, getResp)
				assert.Equal(t, configs.AllowedDimensions[0].Height, height)
			}(token, receiptID)
		}
		wg.Wait()
	})
}
