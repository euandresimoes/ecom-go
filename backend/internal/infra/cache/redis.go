package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis(url string, pwd string) (*redis.Client, error) {
	r := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: pwd,
		DB:       0,
		Protocol: 2,
	})

	ctx := context.Background()

	err := r.Set(
		ctx,
		"redis_test",
		"true",
		0,
	).Err()
	if err != nil {
		return nil, err
	}

	val, err := r.Get(
		ctx,
		"redis_test",
	).Result()
	if err != nil {
		return nil, err
	}

	fmt.Println("redis working?", val)

	return r, nil
}

func Get[T any](r *redis.Client, key string) (*T, error) {
	val, err := r.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	out := new(T)

	err = json.Unmarshal([]byte(val), out)
	if err != nil {
		return nil, err
	}

	return out, err
}

func Set(r *redis.Client, key string, data any) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return r.Set(context.Background(), key, bytes, time.Minute*30).Err()
}

func DeleteMany(r *redis.Client, key string) error {
	keys, _, err := r.Scan(
		context.Background(),
		uint64(0),
		key,
		100,
	).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		r.Del(context.Background(), keys...)
	}

	return nil
}

func DeleteUnique(r *redis.Client, key string) {
	r.Del(context.Background(), key)
}
