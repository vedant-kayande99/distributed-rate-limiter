package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"rls-server/pkg/rate-limiter-service/simple"
	"time"

	pb "github.com/vedant-kayande99/distributed-rate-limiter/proto"

	"google.golang.org/grpc"
)

type rlServer struct {
	pb.UnimplementedRateLimiterServiceServer
	limiter *simple.RateLimiter
}

func (server *rlServer) ShouldAllow(ctx context.Context, req *pb.RateLimitRequest) (*pb.RateLimitResponse, error) {
	log.Printf("Received request for user_id: %v", req.GetUserId())
	allowed := server.limiter.ShouldAllow(req.GetUserId())
	return &pb.RateLimitResponse{Allowed: allowed}, nil
}

func main() {
	rateLimiter, err := simple.NewRateLimiter(5, 10*time.Second)
	if err != nil {
		log.Fatalf("ERROR: Failed to create rate limiter: %v", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("ERROR: Failed to listen on port 50051: %v", err)
	}

	rpcServer := grpc.NewServer()
	pb.RegisterRateLimiterServiceServer(rpcServer, &rlServer{limiter: rateLimiter})

	fmt.Println("gRPC Rate Limiter Server Listening on port :50051")
	if err := rpcServer.Serve(lis); err != nil {
		log.Fatalf("ERROR: Failed to server: %v", err)
	}

}