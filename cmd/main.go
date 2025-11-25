package main

import (
	"log"
	"os"

	"github.com/euandresimoes/ecom-go/internal/cache"
	"github.com/euandresimoes/ecom-go/internal/database"
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
		2: os.Getenv("REDIS_URL"),
	}

	db, err := database.NewPostgres(envs[1])
	if err != nil {
		log.Fatalf("An error occurred while trying to connect to postgres: %s", err)
	}

	redis, err := cache.NewRedis(envs[2], "")
	if err != nil {
		log.Fatalf("An error occurred while trying to connect to redis: %s", err)
	}

	api := Api{
		addr:  envs[0],
		db:    db,
		redis: redis,
	}

	api.Start()
}
