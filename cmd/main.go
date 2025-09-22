package main

import (
	"distributed-rate-limiter/pkg/rate-limiter-service/simple"
	"fmt"
	"time"
)


func main(){
	testUserId := "user 1"
	numOfReqAllowed := 5
	windowLen := 10*time.Second
	simpleRL := simple.NewRateLimiter(numOfReqAllowed, windowLen)

	fmt.Println("Testing the Sliding Window Log Rate Limiter")
	
	for i := 1; i <= 7; i++ {
		allowed := simpleRL.ShouldAllow(testUserId)
		fmt.Printf("Request %d: Allowed: %v\n", i, allowed)
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("\nWaiting for window to slide...")
	time.Sleep(windowLen)

	fmt.Printf("\nTesting %s again after window has passed", testUserId)
	allowed := simpleRL.ShouldAllow(testUserId)
	fmt.Printf("\nRequest 8: Allowed: %v\n", allowed)
}
