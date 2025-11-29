package auth

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	service   *Service
	validator *validator.Validate
}

func NewHandler(service *Service, validator *validator.Validate) http.Handler {
	h := &Handler{service: service, validator: validator}

	r := chi.NewRouter()

	r.Post("/register", h.Register)
	r.Post("/login", h.Login)

	return r
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var data UserRegisterModel

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status": http.StatusBadRequest,
			"error":  "invalid json",
		})
		return
	}

	err = h.validator.Struct(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	err = h.service.register(data)
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
		"message": "account created successfully",
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var data UserLoginModel

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status": http.StatusBadRequest,
			"error":  "invalid json",
		})
		return
	}

	err = h.validator.Struct(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	token, err := h.service.login(data)
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
		"message": "login success",
		"data": map[string]any{
			"token": token,
		},
	})
}
