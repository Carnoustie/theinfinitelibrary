# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git

# Copy dependency files first
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build binary
RUN go build -o server server.go

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
