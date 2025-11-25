package product

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) http.Handler {
	h := &Handler{service: service}

	r := chi.NewRouter()

	r.Post("/", h.Create)
	r.Delete("/", h.Delete)
	r.Patch("/", h.Update)
	r.Get("/", h.GetAll)
	r.Get("/id", h.GetByID)
	r.Get("/public", h.GetByPublicID)

	return r
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var data ProductCreateDto

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status": http.StatusBadRequest,
			"error":  "invalid json",
		})
		return
	}

	product, err := h.service.Create(data)
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

	var data ProductUpdateDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	product, err := h.service.Update(id, data)
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
