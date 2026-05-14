package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
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

	userDB := database.NewUser(db)
	userHandler := handlers.NewUserHandler(userDB)

	productDB := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDB)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", config.TokenAuth))
	r.Use(middleware.WithValue("exp", config.JwtExpiresIn))
	r.Use(LogRequest)

	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.TokenAuth))
		r.Use(jwtauth.Authenticator(config.TokenAuth))
		r.Get("/{id}", productHandler.GetProduct)
		r.Get("/", productHandler.ListProducts)
		r.Post("/", productHandler.CreateProduct)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser)
		r.Post("/generate-token", userHandler.GetJWT)
	})

	fmt.Printf("Server running on: %s:%s\n", config.DBHost, config.WebServerPort)
	http.ListenAndServe(fmt.Sprintf("%s:%s", config.DBHost, config.WebServerPort), r)
}

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetReqID(r.Context())
		fmt.Printf("Received request: %s %s %s\n", r.Method, r.URL.Path, requestID)
		next.ServeHTTP(w, r)
	})
}
