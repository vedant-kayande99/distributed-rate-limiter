FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /rls-server-build ./cmd/main.go

FROM alpine:latest

COPY --from=builder /rls-server-build /rls-server-build

EXPOSE 50051

CMD ["/rls-server-build"]