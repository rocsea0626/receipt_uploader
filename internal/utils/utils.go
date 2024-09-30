package utils

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"receipt_uploader/constants"
	"receipt_uploader/internal/handlers"
	"receipt_uploader/internal/image_worker"
	"receipt_uploader/internal/images"
	"receipt_uploader/internal/logging"
	"receipt_uploader/internal/middlewares"
	"receipt_uploader/internal/models/configs"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

func LoadConfig() (*configs.Config, error) {
	env := os.Getenv("env")
	envFile := ".env"
	if env == "dev" {
		envFile = ".env.development"
	}
	loadErr := godotenv.Load(envFile)
	if loadErr != nil {
		return nil, loadErr
	}

	config := &configs.Config{
		Port:       os.Getenv("PORT"),
		ResizedDir: filepath.Join(constants.ROOT_DIR_IMAGES, os.Getenv("DIR_RESIZED")),
		UploadsDir: filepath.Join(constants.ROOT_DIR_IMAGES, os.Getenv("DIR_UPLOADS")),
		Dimensions: configs.AllowedDimensions,
	}

	return config, nil
}

func startWorker(processFunc func(), stopChan <-chan struct{}) {
	logging.Debugf("startWorker()...")
	var wg sync.WaitGroup

	for {
		select {
		case <-stopChan:
			fmt.Println("Stopping startWorker()")
			wg.Wait()
			return
		default:
			wg.Add(1)
			defer wg.Done()
			processFunc()
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func StartServer(config *configs.Config, stopChan chan struct{}) {
	fmt.Println("starting server...")

	imagesErr := os.MkdirAll(config.ResizedDir, 0755)
	if imagesErr != nil {
		fmt.Printf("failed to start server, err: %s", imagesErr.Error())
		return
	}

	uploadsErr := os.MkdirAll(config.UploadsDir, 0755)
	if uploadsErr != nil {
		fmt.Printf("failed to start server, err: %s", uploadsErr.Error())
		return
	}

	imagesService := images.NewService(&config.Dimensions)
	imgWorkerService := image_worker.NewService(imagesService)

	go startWorker(
		func() {
			resizeErr := imgWorkerService.ResizeImages(config.UploadsDir, config.ResizedDir)
			if resizeErr != nil {
				logging.Errorf("workerService.ResizeImages() failed, err: %s", resizeErr.Error())
			}
		},
		stopChan,
	)

	http.HandleFunc("/health", handlers.HealthHandler())
	http.Handle("/receipts",
		middlewares.Auth(http.HandlerFunc(handlers.UploadReceipt(config, imagesService))),
	)
	http.Handle("/receipts/{receiptId}",
		middlewares.Auth(http.HandlerFunc(handlers.DownloadReceipt(config, imagesService))),
	)

	fmt.Println("Starting server on ", constants.PORT)
	if err := http.ListenAndServe(constants.PORT, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}

}
