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
	"golang.org/x/exp/rand"
)

func TestMainStess(t *testing.T) {

	baseDir := "stress-test-images"
	os.RemoveAll(baseDir)

	config := &configs.Config{
		Port:          ":8080",
		ResizedDir:    filepath.Join(baseDir, "resized"),
		UploadsDir:    filepath.Join(baseDir, "uploads"),
		Dimensions:    configs.AllowedDimensions,
		Mode:          "release",
		QueueCapacity: 100,
		WorkerCount:   3,
	}
	numClients := config.QueueCapacity
	baseUrl := "http://localhost" + config.Port
	// defer os.RemoveAll(baseDir)

	var wg sync.WaitGroup
	stopChan := make(chan struct{})

	t.Cleanup(func() {
		log.Println("Cleanup stress test")
		close(stopChan)
		wg.Wait()
	})

	wg.Add(1)
	go func() {
		defer wg.Done()
		utils.StartServer(config, stopChan)
	}()

	t.Run("stress testing, multiple POST and GET inter-changeably", func(t *testing.T) {

		// to prepare test, upload 10 images sequentialy
		receiptIDs := make(map[string]string)
		url := baseUrl + "/receipts"
		size := "medium"
		var wg sync.WaitGroup

		for i := 0; i < numClients; i++ {
			userToken := "inter_user_token_" + strconv.Itoa(i)
			req, reqErr := test_utils.GenerateUploadRequest(t, url, "test_image.jpg", userToken)
			assert.Nil(t, reqErr)

			client := &http.Client{}
			resp, err := client.Do(req)
			assert.Nil(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusCreated, resp.StatusCode)

			var uploadResp http_responses.UploadResponse
			test_utils.ParseResponseBody(t, resp, &uploadResp)
			assert.NotEmpty(t, uploadResp.ReceiptID)
			receiptIDs[userToken] = uploadResp.ReceiptID
		}

		// generate GET and POST requests
		reqs := []*http.Request{}
		methods := []string{http.MethodGet, http.MethodPost}

		for i := 0; i < numClients; i++ {
			rand.Seed(uint64(time.Now().UnixNano()))
			randomInt := rand.Intn(2)

			method := methods[randomInt]
			if method == http.MethodGet {
				userToken := "inter_user_token_" + strconv.Itoa(i)
				receiptID := receiptIDs[userToken]
				getUrl := fmt.Sprintf("%s/%s?size=%s", url, receiptID, size)
				getReq, getReqErr := http.NewRequest(http.MethodGet, getUrl, nil)
				assert.Nil(t, getReqErr)
				getReq.Header.Set("username_token", userToken)
				reqs = append(reqs, getReq)
			}

			if method == http.MethodPost {
				postUserToken := "inter_post_user_token_" + strconv.Itoa(i)
				postReq, postReqErr := test_utils.GenerateUploadRequest(t, url, "test_image.jpg", postUserToken)
				assert.Nil(t, postReqErr)
				reqs = append(reqs, postReq)
			}
		}
		waitTime := time.Duration(numClients/2) * time.Second
		logging.Infof("waiting %d seconds to allowe server to resize all uploaded images", numClients/2)
		time.Sleep(waitTime)

		// send GET or POST request interchagnabley at same time
		for _, req := range reqs {
			wg.Add(1)
			go func() {
				defer wg.Done()

				client := &http.Client{}
				resp, err := client.Do(req)
				assert.Nil(t, err)
				if err != nil {
					logging.Errorf("err: %s", err.Error())
				}
				defer resp.Body.Close()

				if req.Method == http.MethodGet {
					assert.Equal(t, http.StatusOK, resp.StatusCode)
					header, headerErr := test_utils.ParseDownloadResponseHeader(resp)
					assert.Nil(t, headerErr)
					assert.Equal(t, "image/jpeg", header.ContentType)

					_, height := test_utils.GetImageDimension(t, resp)
					assert.Equal(t, configs.AllowedDimensions[1].Height, height)
				}

				if req.Method == http.MethodPost {
					assert.Equal(t, http.StatusCreated, resp.StatusCode)

					var uploadResp http_responses.UploadResponse
					test_utils.ParseResponseBody(t, resp, &uploadResp)
					assert.NotEmpty(t, uploadResp.ReceiptID)
				}
			}()
		}
		wg.Wait()
	})
}
