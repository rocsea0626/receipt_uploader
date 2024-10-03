package main

import (
	"fmt"
	"os"
	"os/signal"
	"receipt_uploader/internal/logging"
	"receipt_uploader/internal/utils"
	"sync"
	"syscall"
)

func main() {
	config, configErr := utils.LoadConfig()
	logging.SetGlobalLevel(logging.DEBUG_LEVEL)

	var wg sync.WaitGroup

	if configErr != nil {
		fmt.Printf("utils.LoadConfig() failed, err: %s", configErr.Error())
		return
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	stopChan := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		utils.StartServer(config, stopChan)
	}()

	<-signalChan
	close(stopChan)
	fmt.Println("Shutting down server...")
	wg.Wait()
}
