package utils

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"receipt_uploader/internal/constants"
	"receipt_uploader/internal/handlers"
	"receipt_uploader/internal/images"
	"receipt_uploader/internal/logging"
	"receipt_uploader/internal/middlewares"
	"receipt_uploader/internal/models/configs"
	"receipt_uploader/internal/resize_queue"
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

	config := &configs.Config{
		Port:          os.Getenv("PORT"),
		ResizedDir:    filepath.Join(constants.ROOT_DIR_IMAGES, os.Getenv("DIR_RESIZED")),
		UploadsDir:    filepath.Join(constants.ROOT_DIR_IMAGES, os.Getenv("DIR_UPLOADS")),
		Dimensions:    configs.AllowedDimensions,
		Mode:          os.Getenv("MODE"),
		QueueCapacity: constants.QUEUE_CAPACITY,
	}

	return config, nil
}

func StartServer(config *configs.Config, stopChan chan struct{}) {
	fmt.Println("starting server...")
	if config.Mode == "release" {
		logging.SetGlobalLevel(logging.INFO_LEVEL)
		fmt.Println("running in release mode, set log level to INFO")
	}

	initErr := initDirs(config)
	if initErr != nil {
		fmt.Printf("failed to start server, err: %s", initErr.Error())
		return
	}

	imagesService := images.NewService(&config.Dimensions)
	resizeQueue := resize_queue.NewService(config.QueueCapacity, imagesService)
	go resizeQueue.Start(stopChan)

	srv := &http.Server{
		Addr:    config.Port,
		Handler: setupRouter(config, imagesService, resizeQueue),
	}

	go func() {
		fmt.Println("Starting server on ", config.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Error starting server:", err)
		}
	}()

	<-stopChan
	fmt.Println("Received shutdown signal, shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("Server shutdown failed:", err)
	} else {
		fmt.Println("Server exited gracefully")
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

func setupRouter(config *configs.Config, imagesService images.ServiceType, resizeQueue resize_queue.ServiceType) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers.HealthHandler())
	mux.Handle("/receipts", middlewares.Auth(http.HandlerFunc(handlers.UploadReceipt(config, imagesService, resizeQueue))))
	mux.Handle("/receipts/{receiptId}", middlewares.Auth(http.HandlerFunc(handlers.DownloadReceipt(config, imagesService))))
	return mux
}
