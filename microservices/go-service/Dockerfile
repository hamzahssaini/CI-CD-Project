# Stage 1: Build
FROM golang:1.21 AS builder

WORKDIR /app
COPY . .

RUN go mod init temp-module && go mod tidy

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Stage 2: Minimal final image
FROM alpine:latest

# Add CA certificates for MongoDB and HTTPS support
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 3002

CMD ["./main"]
