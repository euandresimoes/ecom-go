package auth

import (
	"encoding/json"
	"net/http"

	"github.com/euandresimoes/ecom-go/backend/internal/infra/security"
	"github.com/euandresimoes/ecom-go/backend/internal/middlewares"
	"github.com/euandresimoes/ecom-go/backend/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	service    *Service
	validator  *validator.Validate
	jwtManager *security.JWTManager
}

func NewHandler(service *Service, validator *validator.Validate, jwtManager *security.JWTManager) http.Handler {
	h := &Handler{
		service:    service,
		validator:  validator,
		jwtManager: jwtManager,
	}

	r := chi.NewRouter()

	r.Post("/register", h.Register)
	r.Post("/login", h.Login)

	r.Group(func(protected chi.Router) {
		protected.Use(middlewares.Auth(h.jwtManager))
		protected.Get("/profile", h.Profile)
	})

	return r
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var data models.UserRegisterModel

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

	err = h.service.Register(data)
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
	var data models.UserLoginModel

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

	token, err := h.service.Login(data)
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

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(models.UserIDKey).(float64)

	profile, err := h.service.Profile(id)
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
		"message": "user profile",
		"data":    profile,
	})
}
