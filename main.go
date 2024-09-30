package main

import (
	"fmt"
	"os"
	"os/signal"
	"receipt_uploader/internal/logging"
	"receipt_uploader/internal/utils"
	"syscall"
)

func main() {
	config, configErr := utils.LoadConfig()
	logging.SetGlobalLevel(logging.DEBUG_LEVEL)

	if configErr != nil {
		fmt.Printf("utils.LoadConfig() failed, err: %s", configErr.Error())
		return
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	stopChan := make(chan struct{})

	go utils.StartServer(config, stopChan)

	<-signalChan
	close(stopChan)
	fmt.Println("Shutting down server...")
}
