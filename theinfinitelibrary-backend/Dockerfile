# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git

# deps first for caching
COPY go.mod go.sum ./
RUN go mod download

# source
COPY . .

# Build binary (build the package, not a single file)
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Final stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/server /app/server

EXPOSE 8080
CMD ["/app/server"]
