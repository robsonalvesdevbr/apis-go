package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/robsonalvesdevbr/apis-go/internal/dto"
	"github.com/robsonalvesdevbr/apis-go/internal/entity"
	"github.com/robsonalvesdevbr/apis-go/internal/infra/database"
)

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(productDB database.ProductInterface) *ProductHandler {
	return &ProductHandler{ProductDB: productDB}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInputDTO
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newProduct, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.ProductDB.Create(newProduct)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newProduct)
}
