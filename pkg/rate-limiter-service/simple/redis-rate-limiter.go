package simple

import (
	"context"
	"distributed-rate-limiter/pkg/rate-limiter-service/store"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/redis/go-redis/v9"
)

var rateLimitCheckScript *redis.Script

func init() {
	scriptPath := filepath.Join("pkg", "rate-limiter-service", "scripts","rate-limit-check.lua")
	scriptContent, err := os.ReadFile(scriptPath)
	if err != nil {
		panic(fmt.Sprintf("ERROR: Failed to load the rate-limit-check.lua script: %v", err))
	}
	rateLimitCheckScript = redis.NewScript(string(scriptContent))
}

type RateLimiter struct {
	rateLimit int
	windowSize time.Duration
	redisStore *store.RedisStore
}

func NewRateLimiter(limit int, window time.Duration) (*RateLimiter, error) {
	redisClient, err := store.NewRedisStore()
	if err != nil {
		return nil, fmt.Errorf("ERROR: Failed to create Redis Store: %v", err)		
	}
	return &RateLimiter{
		rateLimit: limit,
		windowSize: window,
		redisStore: redisClient,
	}, nil
}

func (rl *RateLimiter) ShouldAllow(userId string) bool{
	ctx := context.Background()
	key := fmt.Sprintf("rate-limit:%s", userId)
	curTime := time.Now().UnixNano()
	
	result, err := rateLimitCheckScript.Run(ctx, rl.redisStore.Client, 
		[]string{key},
		curTime,
		rl.windowSize.Nanoseconds(),
		rl.rateLimit,
	).Int64()

	if err != nil {
		fmt.Printf("ERROR: Error executing the Redis pipeline Tx: %v. Allowing Request!", err)
		return true
	}

	return result == 1
}
