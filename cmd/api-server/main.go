package main

import (
	"context"
	pb "distributed-rate-limiter/proto"
	"fmt"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var rlsClient pb.RateLimiterServiceClient

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error: Couldn't connect to the RPC server at: localhost:50051, %v", err)
	}

	defer conn.Close()
	rlsClient = pb.NewRateLimiterServiceClient(conn)
	http.HandleFunc("/data", rateLimitedEndpoint)
	fmt.Println("Protected API Server Listening on port: 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func rateLimitedEndpoint(w http.ResponseWriter, r *http.Request) {
	userId := "user-123"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := rlsClient.ShouldAllow(ctx, &pb.RateLimitRequest{UserId: userId})
	if err != nil {
		log.Printf("WARN: Could not check Rate Limit: %v. Allowing Req as fallback.", err)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK (RLS Down)"))
	}

	if res.GetAllowed() {
		log.Printf("Request ALLOWED for user %s", userId)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello! You are within your rate limit."))
	} else {
		log.Printf("Request DENIED for user %s", userId)
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("Too many Requests."))
	}

}