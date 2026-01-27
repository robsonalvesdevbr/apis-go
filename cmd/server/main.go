package main

import (
	"fmt"
	"net/http"

	"github.com/robsonalvesdevbr/apis-go/configs"
	"github.com/robsonalvesdevbr/apis-go/internal/entity"
	"github.com/robsonalvesdevbr/apis-go/internal/infra/database"
	"github.com/robsonalvesdevbr/apis-go/internal/infra/webserver/handlers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w", err))
	}

	db, err := gorm.Open(sqlite.Open("apis.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&entity.User{}, &entity.Product{})

	productDB := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDB)

	http.HandleFunc("/products", productHandler.CreateProduct)

	fmt.Printf("Server running on: %s:%s\n", config.DBHost, config.WebServerPort)
	http.ListenAndServe(fmt.Sprintf("%s:%s", config.DBHost, config.WebServerPort), nil)
}
