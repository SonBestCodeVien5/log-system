# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod .
COPY main.go .
RUN go mod download && go build -o demo-go .

# Run stage — image nhỏ, không có Go toolchain
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/demo-go .
CMD ["./demo-go"]
