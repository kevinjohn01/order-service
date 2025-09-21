package repository

import (
    "context"
    "order-service/models"
    "github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
    DB *pgxpool.Pool
}

func (r *OrderRepository) Create(ctx context.Context, order models.Order) error {
    _, err := r.DB.Exec(ctx,
        "INSERT INTO orders (product_id, total_price, status, created_at) VALUES ($1, $2, $3, $4)",
        order.ProductID, order.TotalPrice, order.Status, order.CreatedAt)
    return err
}

func (r *OrderRepository) FindByProductID(ctx context.Context, productId int) ([]models.Order, error) {
    rows, err := r.DB.Query(ctx, "SELECT id, product_id, total_price, status, created_at FROM orders WHERE product_id=$1", productId)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var orders []models.Order
    for rows.Next() {
        var o models.Order
        err = rows.Scan(&o.ID, &o.ProductID, &o.TotalPrice, &o.Status, &o.CreatedAt)
        if err != nil {
            return nil, err
        }
        orders = append(orders, o)
    }
    return orders, nil
}
