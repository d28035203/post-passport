FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Final stage
FROM alpine:latest

RUN addgroup -S goapp && adduser -S goapp -G goapp
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Fix permissions for .env file
RUN chmod 644 /app/.env

# Install Go for integration testing
RUN apk add --no-cache go@latest

USER goapp

EXPOSE 8080

CMD ["./main"]
