package main

import (
	"fmt"
	"net/http"
	"os"
	"receipt_uploader/constants"
	"receipt_uploader/internal/handlers"
	"receipt_uploader/internal/images"
	"receipt_uploader/internal/logging"
	"receipt_uploader/internal/middlewares"
	"receipt_uploader/internal/models/configs"
	"receipt_uploader/internal/utils"
)

func StartServer(config *configs.Config) {
	fmt.Println("starting server...")

	fmt.Printf("creating dir: %s to store images\n", config.ImagesDir)
	imagesErr := os.MkdirAll(config.ImagesDir, 0755)
	if imagesErr != nil {
		fmt.Printf("failed to start server, err: %s", imagesErr.Error())
		return
	}

	imagesService := images.NewService(&config.Dimensions)

	http.HandleFunc("/health", handlers.HealthHandler())
	http.Handle("/receipts",
		middlewares.Auth(http.HandlerFunc(handlers.UploadReceipt(config, imagesService))),
	)
	http.Handle("/receipts/{receiptId}",
		middlewares.Auth(http.HandlerFunc(handlers.DownloadReceipt(config, imagesService))),
	)

	fmt.Printf("Starting server on %s", constants.PORT)
	if err := http.ListenAndServe(constants.PORT, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func main() {
	config, configErr := utils.LoadConfig()
	logging.SetGlobalLevel(logging.DEBUG_LEVEL)

	if configErr != nil {
		fmt.Printf("utils.LoadConfig() failed, err: %s", configErr.Error())
		fmt.Println("failed to load config")
		return
	}

	StartServer(config)
}
