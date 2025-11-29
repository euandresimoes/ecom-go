package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/euandresimoes/ecom-go/internal/domain/auth"
	"github.com/euandresimoes/ecom-go/internal/domain/product"
	"github.com/euandresimoes/ecom-go/internal/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func (api *Api) routes() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middlewares.JSON)

	// health check endpoint
	r.Get("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"message": "all good!",
		})
	})

	// utils
	jwtManager := auth.NewJWTManager(api.jwtSecret, api.jwtExp)
	validator := validator.New()

	// handlers
	authRepo := auth.NewRepository(api.db, api.redis, jwtManager)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService, validator)
	r.Mount("/api/v1/auth", authHandler)

	productRepo := product.NewRepository(api.db, api.redis)
	productService := product.NewService(productRepo)
	productHandler := product.NewHandler(productService, validator, jwtManager)
	r.Mount("/api/v1/product", productHandler)

	return r
}

func (api *Api) Start() {
	r := api.routes()
	log.Printf("server running on %s", api.addr)
	if err := http.ListenAndServe(api.addr, r); err != nil {
		log.Fatal(err)
	}
}

type Api struct {
	addr      string
	db        *pgxpool.Pool
	redis     *redis.Client
	jwtSecret string
	jwtExp    time.Duration
}
