package main

import (
	"distributed-rate-limiter/pkg/rate-limiter-service/simple"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
)



func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("WARN: Error loading .env file: %v", err)
	}
}

func main(){
	testUserId := "user 1"
	numOfReqAllowed := 5
	windowLen := 10*time.Second
	// simpleRL := simple.NewRateLimiter(numOfReqAllowed, windowLen)
	redisRL,_ := simple.NewRateLimiter(numOfReqAllowed, windowLen)

	fmt.Println("Testing the Sliding Window Log Rate Limiter")
	
	for i := 1; i <= 7; i++ {
		allowed := redisRL.ShouldAllow(testUserId)
		fmt.Printf("Request %d: Allowed: %v\n", i, allowed)
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("\nWaiting for window to slide...")
	time.Sleep(windowLen)

	fmt.Printf("\nTesting %s again after window has passed", testUserId)
	allowed := redisRL.ShouldAllow(testUserId)
	fmt.Printf("\nRequest 8: Allowed: %v\n", allowed)
}
