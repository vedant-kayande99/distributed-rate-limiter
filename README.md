# distributed-rate-limiter

Distributed Rate Limiter built with Go for controlling API traffic across multiple nodes.

## To run this application:

- docker-compose -f docker/redis.yml --env-file .env up -d
- go run cmd/main.go
