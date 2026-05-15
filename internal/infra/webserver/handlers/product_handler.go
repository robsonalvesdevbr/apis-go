package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

// Create Product godoc
// @Summary Create a new product
// @Description Create a new product with the provided information
// @Tags products
// @Accept json
// @Produce json
// @Param request body dto.CreateProductInputDTO true "Product information"
// @Success 201 {object} entity.Product
// @Failure 400 {object} handlers.Error
// @Failure 500 {object} handlers.Error
// @Router /products [post]
// @Security ApiKeyAuth
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInputDTO
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newProduct, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.ProductDB.Create(newProduct)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newProduct)
}

// Get Product godoc
// @Summary Get a product by ID
// @Description Retrieve a product by its unique ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} entity.Product
// @Failure 400 {object} handlers.Error
// @Failure 404 {object} handlers.Error
// @Router /products/{id} [get]
// @Security ApiKeyAuth
func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// Update Product godoc
// @Summary Update a product by ID
// @Description Update a product's information by its unique ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body dto.CreateProductInputDTO true "Updated product information"
// @Success 200 {object} entity.Product
// @Failure 400 {object} handlers.Error
// @Failure 404 {object} handlers.Error
// @Failure 500 {object} handlers.Error
// @Router /products/{id} [put]
// @Security ApiKeyAuth
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var product entity.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product.ID, err = uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	existingProduct, err := h.ProductDB.FindByID(product.ID.String())
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	product.CreatedAt = existingProduct.CreatedAt
	product.UpdatedAt = time.Now().UTC()

	err = h.ProductDB.Update(&product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// Delete Product godoc
// @Summary Delete a product by ID
// @Description Delete a product by its unique ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 204 "No Content"
// @Failure 400 {object} handlers.Error
// @Failure 404 {object} handlers.Error
// @Failure 500 {object} handlers.Error
// @Router /products/{id} [delete]
// @Security ApiKeyAuth
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.ProductDB.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// List Products godoc
// @Summary List products with pagination and sorting
// @Description Retrieve a list of products with optional pagination and sorting parameters
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param sort query string false "Sort by field (e.g., name, price)"
// @Success 200 {array} entity.Product
// @Failure 500 {object} handlers.Error
// @Router /products [get]
// @Security ApiKeyAuth
func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	pageParam := r.URL.Query().Get("page")
	limitParam := r.URL.Query().Get("limit")
	sortParam := r.URL.Query().Get("sort")

	var page, limit int
	if pageParam != "" {
		page, _ = strconv.Atoi(pageParam)
	}
	if limitParam != "" {
		limit, _ = strconv.Atoi(limitParam)
	}

	products, err := h.ProductDB.FindAll(page, limit, sortParam)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}
