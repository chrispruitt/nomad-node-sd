# Stage 1: Build the Go binary
FROM golang:1.20 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Stage 2: Create a lightweight image
FROM alpine:3.14

WORKDIR /app

COPY --from=builder /app/app .

CMD ["./app"]