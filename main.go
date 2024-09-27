package main

import (
	"fmt"
	"log"
	"receipt_uploader/internal/utils"
)

func main() {
	config, configErr := utils.LoadConfig()
	if configErr != nil {
		log.Printf("utils.LoadConfig() failed, err: %s", configErr.Error())
		fmt.Println("failed to load config")
		return
	}

	utils.StartServer(config)
}
