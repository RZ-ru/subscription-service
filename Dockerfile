# ---------------------------
# 1. Build stage
# ---------------------------
FROM golang:1.25 AS builder

WORKDIR /app

# Кэшируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем бинарь
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o subs-service ./cmd/sub-service/main.go


# ---------------------------
# 2. Final stage
# ---------------------------
FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/subs-service .
COPY ./docs ./docs
COPY ./internal/db/migrations ./internal/db/migrations
COPY .env .env

EXPOSE 8080

CMD ["./subs-service"]
