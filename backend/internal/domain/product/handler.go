package product

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/euandresimoes/ecom-go/backend/internal/infra/security"
	"github.com/euandresimoes/ecom-go/backend/internal/middlewares"
	"github.com/euandresimoes/ecom-go/backend/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	service   *Service
	validator *validator.Validate
}

func NewHandler(service *Service, validator *validator.Validate, jwt *security.JWTManager) http.Handler {
	h := &Handler{service: service, validator: validator}

	r := chi.NewRouter()

	// public routes
	r.Get("/", h.GetAll)
	r.Get("/id", h.GetByID)
	r.Get("/public", h.GetByPublicID)
	r.Get("/category", h.GetAllCategories)

	// admin protected routes
	r.Group(func(protected chi.Router) {
		protected.Use(middlewares.Admin(jwt))

		protected.Post("/", h.Create)
		protected.Delete("/", h.Delete)
		protected.Patch("/", h.Update)
		protected.Post("/category", h.CreateCategory)
		protected.Delete("/category", h.DeleteCategory)
	})

	return r
}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var data models.CategoryCreateDto

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status": http.StatusBadRequest,
			"error":  "invalid json",
		})
		return
	}

	if err := h.validator.Struct(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	category, err := h.service.CreateCategory(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"status":  http.StatusCreated,
		"message": "category created",
		"data":    category,
	})
}

func (h *Handler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"status":  http.StatusOK,
		"message": "categories found",
		"data":    categories,
	})
}

func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	category, err := h.service.DeleteCategory(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"status":  http.StatusOK,
		"message": "category deleted",
		"data":    category,
	})
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var data models.ProductCreateDto

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status": http.StatusBadRequest,
			"error":  "invalid json",
		})
		return
	}

	if err := h.validator.Struct(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	product, err := h.service.Create(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"status":  http.StatusCreated,
		"message": "product created",
		"data":    product,
	})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	product, err := h.service.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"status":  http.StatusOK,
		"message": "product deleted",
		"data":    product,
	})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	var data models.ProductUpdateDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	if err := h.validator.Struct(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	product, err := h.service.Update(id, &data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"status":  http.StatusOK,
		"message": "product updated",
		"data":    product,
	})
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"status":  http.StatusOK,
		"message": "products found",
		"data":    products,
	})
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	product, err := h.service.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"status":  http.StatusOK,
		"message": "products found",
		"data":    product,
	})
}

func (h *Handler) GetByPublicID(w http.ResponseWriter, r *http.Request) {
	publicID := r.URL.Query().Get("public_id")

	product, err := h.service.GetByPublicID(publicID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"status":  http.StatusOK,
		"message": "product found",
		"data":    product,
	})
}
