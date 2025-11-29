package main

import (
	"log"
	"os"
	"time"

	"github.com/euandresimoes/ecom-go/internal/infra/cache"
	"github.com/euandresimoes/ecom-go/internal/infra/database"
	"github.com/joho/godotenv"
)

func main() {
	// Load env variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	envs := map[int]string{
		0: os.Getenv("API_ADDR"),
		1: os.Getenv("DATABASE_URL"),
		2: os.Getenv("ADMIN_PASSWORD"),
		3: os.Getenv("REDIS_URL"),
		4: os.Getenv("JWT_SECRET"),
	}

	db, err := database.NewPostgres(envs[1])
	if err != nil {
		log.Fatalf("An error occurred while trying to connect to postgres: %s", err)
	}

	err = database.CreateAdmin(envs[2], db)
	if err != nil {
		log.Fatalf("An error occurred while trying to create admin user: %s", err)
	}

	redis, err := cache.NewRedis(envs[3], "")
	if err != nil {
		log.Fatalf("An error occurred while trying to connect to redis: %s", err)
	}

	api := Api{
		addr:      envs[0],
		db:        db,
		redis:     redis,
		jwtSecret: envs[4],
		jwtExp:    time.Minute * 5,
	}

	api.Start()
}
