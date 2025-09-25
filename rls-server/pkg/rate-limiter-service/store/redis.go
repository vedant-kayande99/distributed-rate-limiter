package store

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct{
	Client *redis.Client
}

func NewRedisStore() (*RedisStore, error) {
	host := getEnvOrDefault("REDIS_HOST", "localhost")
	port := getEnvOrDefault("REDIS_PORT", "6379")
	password := getEnvOrDefault("REDIS_PASSWORD", "")
	
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s",host,port),
		Password: password,
		DB: 0,
	})

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()

	if err != nil {
		return nil, err
	}
	
	log.Println("Successfully connected to Redis!")
	
	return &RedisStore{
		Client: client,
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}