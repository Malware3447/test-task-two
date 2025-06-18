FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/runner/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/config ./config
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./main", "-config=./config/local.yml"]