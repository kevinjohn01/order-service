package main

import (
    "context"
    "log"
    "order-service/cache"
    "order-service/handlers"
    "order-service/queue"
    "order-service/repository"

    "github.com/gin-gonic/gin"
    "github.com/jackc/pgx/v5/pgxpool"
)

func main() {
    dbpool, err := pgxpool.New(context.Background(), "postgres://postgres:postgre@orderdb:5432/orderdb?sslmode=disable")
    if err != nil {
        log.Fatal("Failed to connect database:", err)
    }

    repo := &repository.OrderRepository{DB: dbpool}
    orderHandler := &handlers.OrderHandler{Repo: repo}

    cache.InitRedis()
    queue.InitRabbitMQ()

    r := gin.Default()
    r.POST("/orders", orderHandler.CreateOrder)
    r.GET("/orders/product/:productId", orderHandler.GetOrders)

    r.Run(":3001")
}