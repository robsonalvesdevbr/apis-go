package main

import (
	"fmt"

	"github.com/robsonalvesdevbr/apis-go/configs"
)

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w", err))
	}

	fmt.Printf("Server running on: %s:%s\n", config.DBHost, config.WebServerPort)
}
