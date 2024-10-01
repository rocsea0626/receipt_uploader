package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"receipt_uploader/internal/constants"
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

var numClients = constants.QUEUE_CAPACITY

func TestMainStess(t *testing.T) {

	baseDir := "stress-test-images"
	os.RemoveAll(baseDir)

	config := &configs.Config{
		Port:       ":8080",
		ResizedDir: filepath.Join(baseDir, "resized"),
		UploadsDir: filepath.Join(baseDir, "uploads"),
		Dimensions: configs.AllowedDimensions,
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

	t.Run("stress testing, multiple POST&GET /receipts request", func(t *testing.T) {

		waitTime := numClients * 1
		var wg sync.WaitGroup
		var mu sync.Mutex
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

				mu.Lock()
				receiptIDs[userToken] = uploadResp.ReceiptID
				mu.Unlock()
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

	t.Run("stress testing, multiple GET /receipts request for same receipt", func(t *testing.T) {
		waitTime := 2
		var wg sync.WaitGroup
		url := baseUrl + "/receipts"

		// upload 1 receipt
		userToken := "user_token_gets"
		req, reqErr := test_utils.GenerateUploadRequest(t, url, "test_image.jpg", userToken)
		assert.Nil(t, reqErr)

		resp, err := client.Do(req)
		assert.Nil(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var uploadResp http_responses.UploadResponse
		test_utils.ParseResponseBody(t, resp, &uploadResp)
		assert.NotEmpty(t, uploadResp.ReceiptID)

		logging.Infof("waiting %d s for server completes resizing of all uploaded images...", waitTime)
		time.Sleep(time.Duration(waitTime) * time.Second)

		size := "medium"
		for i := 0; i < numClients; i++ {
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
				assert.Equal(t, configs.AllowedDimensions[1].Height, height)
			}(userToken, uploadResp.ReceiptID)
		}
		wg.Wait()
	})

	t.Run("stress testing, multiple POST and GET inter-changeably", func(t *testing.T) {

		// to prepare test, upload 10 images sequentialy
		receiptIDs := make(map[string]string)
		url := baseUrl + "/receipts"
		size := "medium"
		numRequests := constants.QUEUE_CAPACITY
		var wg sync.WaitGroup

		for i := 0; i < numRequests; i++ {
			userToken := "inter_user_token_" + strconv.Itoa(i)
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
		}

		// generate GET and POST requests
		reqs := []*http.Request{}
		methods := []string{http.MethodGet, http.MethodPost}

		for i := 0; i < numRequests; i++ {
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

		time.Sleep(time.Duration(numRequests) * time.Second)

		// send request interchagnabley at same time
		for _, req := range reqs {
			wg.Add(1)
			go func() {
				defer wg.Done()

				resp, getErr := client.Do(req)
				assert.Nil(t, getErr)
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
