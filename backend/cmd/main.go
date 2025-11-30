package main

import (
	"log"
	"os"
	"time"

	"github.com/euandresimoes/ecom-go/backend/internal/infra/cache"
	"github.com/euandresimoes/ecom-go/backend/internal/infra/database"
)

func main() {
	// Load env variables - Dev
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	envs := map[int]string{
		0: os.Getenv("ADMIN_EMAIL"),
		1: os.Getenv("ADMIN_PASSWORD"),
		2: os.Getenv("API_ADDR"),
		3: os.Getenv("DATABASE_URL"),
		4: os.Getenv("REDIS_URL"),
		5: os.Getenv("JWT_SECRET"),
	}

	db, err := database.NewPostgres(envs[3])
	if err != nil {
		log.Fatalf("An error occurred while trying to connect to postgres: %s", err)
	}

	err = database.CreateAdmin(envs[0], envs[1], db)
	if err != nil {
		log.Fatalf("An error occurred while trying to create admin user: %s", err)
	}

	redis, err := cache.NewRedis(envs[4], "")
	if err != nil {
		log.Fatalf("An error occurred while trying to connect to redis: %s", err)
	}

	api := Api{
		addr:      envs[2],
		db:        db,
		redis:     redis,
		jwtSecret: envs[5],
		jwtExp:    time.Minute * 5,
	}

	api.Start()
}
