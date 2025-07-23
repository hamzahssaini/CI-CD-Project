# Stage 1: Static build
FROM golang:1.21 AS builder

WORKDIR /app
COPY . .

RUN go mod init temp-module && go mod tidy

# 👇 Add flags to make static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Stage 2: Use scratch (zero size)
FROM scratch

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 3002
CMD ["./main"]
