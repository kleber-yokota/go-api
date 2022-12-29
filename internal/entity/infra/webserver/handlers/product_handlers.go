package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kleber-yokota/go-api/internal/entity"
	"github.com/kleber-yokota/go-api/internal/entity/dto"
	"github.com/kleber-yokota/go-api/internal/entity/infra/database"
	entityPkg "github.com/kleber-yokota/go-api/pkg/entity"
)


type ProductHandler struct {
	ProductDB database.ProdudctInterface
}

func NewProductHandler(db database.ProdudctInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

// CreateProduct godoc
// @Summary     create products
// @Description  create products
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        request   body      dto.CreateProductInput  true  "product request"
// @Success      201 
// @Failure      500  {object}  error
// @Router       /products [post]
// @Security ApiKeyAuth
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	p, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	err = h.ProductDB.Create(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GetProduct godoc
// @Summary     get specific product
// @Description  get specific product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        page   query      string	true "product id" Format(uuid)
// @Success      200	{object} entity.Product
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /products/{id} [get]
// @Security ApiKeyAuth
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: "Must have ID"}
		json.NewEncoder(w).Encode(error)
		return
	}
	product, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// UpdateProduct godoc
// @Summary     update product
// @Description  update product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        page   query      string	true "product id" Format(uuid)
// @Param        request   body      dto.CreateProductInput  true  "product request"
// @Success      200
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /products/{id} [put]
// @Security ApiKeyAuth
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: "Must have ID"}
		json.NewEncoder(w).Encode(error)
		return
	}

	var product entity.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	product.ID, err = entityPkg.ParseID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	_, err = h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	err = h.ProductDB.Update(&product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	w.WriteHeader(http.StatusOK)

}

// DeleteProduct godoc
// @Summary     delete product
// @Description  delete product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        page   query      string	true "product id" Format(uuid)
// @Success      200
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /products/{id} [delete]
// @Security ApiKeyAuth
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	err = h.ProductDB.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	w.WriteHeader(http.StatusOK)

}

// GetProducts godoc
// @Summary     list products
// @Description  get all product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        page   query      string	false "page number"
// @Param        page   query      string	false "limit"
// @Success      200	{array} entity.Product
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /products [get]
// @Security ApiKeyAuth
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request){
	page:=r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	pageInt, err:= strconv.Atoi(page)
	if err != nil {
		pageInt = 0
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 0
	}

	sort := r.URL.Query().Get("sort")

	products, err:= h.ProductDB.FindAll(pageInt, limitInt, sort)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)


}
