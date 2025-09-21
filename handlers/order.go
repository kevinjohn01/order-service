package handlers

import (
    "context"
    "encoding/json"
    "net/http"
    "order-service/cache"
    "order-service/models"
    "order-service/queue"
    "order-service/repository"
    "order-service/services"
    "strconv"
    "time"
	"fmt"
    "github.com/gin-gonic/gin"
)

type OrderHandler struct {
    Repo *repository.OrderRepository
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
    var req struct {
        ProductID int `json:"productId"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // product, err := services.GetProduct(req.ProductID)
    // if err != nil {
    //     c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product"})
    //     return
    // }
	fmt.Printf("Request ProductID: %d\n", req.ProductID)

	product, err := services.GetProduct(req.ProductID)
	if err != nil {
		fmt.Printf("Error GetProduct: %v\n", err)
		return
	}
	fmt.Printf("Product from service: %+v\n", product)

    order := models.Order{
        ProductID:  req.ProductID,
        TotalPrice: product.Price,
        Status:     "pending",
        CreatedAt:  time.Now(),
    }

    err = h.Repo.Create(context.Background(), order)
	fmt.Printf("Error creating order: %v\n", err)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save order"})
        return
    }

    go queue.PublishOrder(struct {
        ProductID int `json:"productId"`
        Qty       int `json:"qty"`
    }{
        ProductID: req.ProductID,
        Qty:       1,
    })

    c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
    pidStr := c.Param("productId")
    pid, _ := strconv.Atoi(pidStr)

    cacheKey := "orders:product:" + pidStr
    cached, _ := cache.GetCache(cacheKey)
    if cached != "" {
        var orders []models.Order
        _ = json.Unmarshal([]byte(cached), &orders)
        c.JSON(http.StatusOK, orders)
        return
    }

    orders, err := h.Repo.FindByProductID(context.Background(), pid)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch orders"})
        return
    }

    bytes, _ := json.Marshal(orders)
    _ = cache.SetCache(cacheKey, string(bytes), time.Minute*5)

    c.JSON(http.StatusOK, orders)
}
