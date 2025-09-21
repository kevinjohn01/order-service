# Build stage
FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o order-service main.go

# Runtime stage
FROM gcr.io/distroless/base-debian12

WORKDIR /
COPY --from=builder /app/order-service /order-service

EXPOSE 8081

CMD ["/order-service"]
