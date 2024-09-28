package main

import (
	"fmt"
	"receipt_uploader/internal/logging"
	"receipt_uploader/internal/utils"
)

func main() {
	config, configErr := utils.LoadConfig()
	logging.SetGlobalLevel(logging.DEBUG)

	if configErr != nil {
		logging.Printf("utils.LoadConfig() failed, err: %s", configErr.Error())
		fmt.Println("failed to load config")
		return
	}

	utils.StartServer(config)
}
