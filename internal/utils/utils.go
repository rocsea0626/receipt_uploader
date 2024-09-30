package utils

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"receipt_uploader/internal/constants"
	"receipt_uploader/internal/handlers"
	"receipt_uploader/internal/image_worker"
	"receipt_uploader/internal/images"
	"receipt_uploader/internal/logging"
	"receipt_uploader/internal/middlewares"
	"receipt_uploader/internal/models/configs"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

func LoadConfig() (*configs.Config, error) {
	env := os.Getenv("env")
	envFile := ".env"
	if env == "dev" {
		envFile = ".env.dev"
	}
	loadErr := godotenv.Load(envFile)
	if loadErr != nil {
		return nil, loadErr
	}

	interval, intervalErr := strconv.Atoi(os.Getenv("INTERVAL"))
	if intervalErr != nil {
		return nil, intervalErr
	}

	config := &configs.Config{
		Port:       os.Getenv("PORT"),
		ResizedDir: filepath.Join(constants.ROOT_DIR_IMAGES, os.Getenv("DIR_RESIZED")),
		UploadsDir: filepath.Join(constants.ROOT_DIR_IMAGES, os.Getenv("DIR_UPLOADS")),
		Dimensions: configs.AllowedDimensions,
		Interval:   time.Duration(interval) * time.Second,
		Mode:       os.Getenv("MODE"),
	}

	return config, nil
}

func StartServer(config *configs.Config, stopChan chan struct{}) {
	fmt.Println("starting server...")
	if config.Mode == "release" {
		fmt.Println("running in release mode, set log level to INFO")
		logging.SetGlobalLevel(logging.INFO_LEVEL)
	}

	initErr := initDirs(config)
	if initErr != nil {
		fmt.Printf("failed to start server, err: %s", initErr.Error())
		return
	}

	imagesService := images.NewService(&config.Dimensions)
	imgWorkerService := image_worker.NewService(imagesService)

	go startWorker(
		config.Interval,
		func() {
			resizeErr := imgWorkerService.ResizeImages(config.UploadsDir, config.ResizedDir)
			if resizeErr != nil {
				logging.Errorf("workerService.ResizeImages() failed, err: %s", resizeErr.Error())
			}
		},
		stopChan,
	)

	setupRouter(config, imagesService)

	fmt.Println("Starting server on ", constants.PORT)
	if err := http.ListenAndServe(constants.PORT, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}

}

func startWorker(internal time.Duration, processFunc func(), stopChan <-chan struct{}) {
	logging.Debugf("startWorker()...")
	var wg sync.WaitGroup

	for {
		select {
		case <-stopChan:
			fmt.Println("Stopping startWorker()")
			wg.Wait()
			return
		case <-time.After(internal):
			logging.Debugf("Wake up after sleeping, starting processing...")
			wg.Add(1)
			defer wg.Done()
			processFunc()
		}
	}
}

func initDirs(config *configs.Config) error {
	imagesErr := os.MkdirAll(config.ResizedDir, 0755)
	if imagesErr != nil {
		return imagesErr
	}

	uploadsErr := os.MkdirAll(config.UploadsDir, 0755)
	if uploadsErr != nil {
		return uploadsErr
	}
	return nil
}

func setupRouter(config *configs.Config, imagesService images.ServiceType) {

	http.HandleFunc("/health", handlers.HealthHandler())
	http.Handle("/receipts",
		middlewares.Auth(http.HandlerFunc(handlers.UploadReceipt(config, imagesService))),
	)
	http.Handle("/receipts/{receiptId}",
		middlewares.Auth(http.HandlerFunc(handlers.DownloadReceipt(config, imagesService))),
	)
}
