package main

import (
	"fmt"
	"log"
	"net/http"
	"receipt_uploader/constants"
	"receipt_uploader/internal/handlers"
	"receipt_uploader/internal/utils"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	fmt.Fprintln(w, "Hello, World!")
}

func main() {
	log.Println("initializing server")
	config, configErr := utils.LoadConfig()
	if configErr != nil {
		log.Printf("utils.LoadConfig() failed, err: %s", configErr.Error())
		fmt.Println("failed to load config")
		return
	}
	utils.InitServer(config)

	http.HandleFunc("/health", helloHandler)
	http.HandleFunc("/receipts", handlers.ReceiptsHandler(config))

	fmt.Printf("Starting server on %s", constants.PORT)
	if err := http.ListenAndServe(constants.PORT, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
