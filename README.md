# distributed-rate-limiter

Distributed Rate Limiter built with Go for controlling API traffic across multiple nodes.

## To run this application:

### Run with Docker Compose

Environment files:

- Create a repo-level `.env`:

```bash
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=${REDIS_PASSWORD}
RATE_LIMITER_HOST=rls-server-lb
```

Start services (from repo root):

```bash
docker compose -f docker/docker-compose.yml --env-file .env up -d
```

Services:

- Redis (data source)
- rls-server (gRPC)
- api-server (protected resource)

Stop:

```bash
docker compose -f docker/docker-compose.yml down
```

Clean volumes (Redis data):

```bash
docker compose -f docker/docker-compose.yml down -v
```

### Rebuild after code changes

Rebuild only `rls-server`:

```bash
docker compose -f docker/docker-compose.yml build rls-server
docker compose -f docker/docker-compose.yml up -d
```

Rebuild everything:

```bash
docker compose -f docker/docker-compose.yml build
docker compose -f docker/docker-compose.yml up -d
```

### Run the test client

Create an external network test-net -> Run the test-client image

```bash
docker network create test-net
docker compose -f test-client/docker-compose.yml up -d
```
