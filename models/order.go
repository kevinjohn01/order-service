package models

import "time"

type Order struct {
    ID         int       `json:"id"`
    ProductID  int       `json:"productId"`
    TotalPrice string   `json:"price"`
    Status     string    `json:"status"`
    CreatedAt  time.Time `json:"createdAt"`
}
