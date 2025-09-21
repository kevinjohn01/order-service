package services

import (
    "encoding/json"
    "fmt"
    "net/http"
	"io"
	"log"
	"strconv"
)

type Product struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Price string `json:"price"`
}

func GetProduct(productID int) (*Product, error) {
    resp, err := http.Get(fmt.Sprintf("http://product-service:3000/products/%d", productID))
    fmt.Printf("HTTP GET Response: %+v\n", resp)
	if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("product not found")
    }

	// 👉 Debug: baca isi body mentah
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read body: %w", err)
    }
    log.Println("Raw body:", string(body)) // 🔍 cek JSON mentah

    var product Product
    if err := json.Unmarshal(body, &product); err != nil {
        return nil, fmt.Errorf("failed to unmarshal product: %w", err)
    }
	priceFloat, err := strconv.ParseFloat(product.Price, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid price format: %w", err)
	}
	log.Println("Parsed Price:", priceFloat)

	 log.Printf("Decoded Product: %+v\n", product) // 🔍 hasil decode
    return &product, nil

    // var product Product
    err = json.NewDecoder(resp.Body).Decode(&product)
	fmt.Printf("Decoded Product: %+v\n", product)
    if err != nil {
        return nil, err
    }
    return &product, nil
}
