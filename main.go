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

	fmt.Printf("creating dir: %s to store uploaded receipts\n", config.UploadedDir)
	tmpErr := os.MkdirAll(config.UploadedDir, 0755)
	if tmpErr != nil {
		fmt.Printf("failed to start server, err: %s", tmpErr.Error())
		return
	}

	fmt.Printf("creating dir: %s to store generated images of receipts\n", config.GeneratedDir)
	imagesErr := os.MkdirAll(config.GeneratedDir, 0755)
	if imagesErr != nil {
		fmt.Printf("failed to start server, err: %s", imagesErr.Error())
		return
	}

	imagesService := images.NewService()

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
