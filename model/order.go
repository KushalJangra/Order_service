package models

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

const (
	DBHost  = "127.0.0.1"
	DBUser  = "root"
	DBPass  = "Kush@123456"
	DBDbase = "pro"
	BaseURLCustomer = "http://localhost:8080/customers"
	BaseURLProduct  = "http://localhost:8080/products"
)

var db *sql.DB

// Initialize database connection
func InitDB() {
	dbConn := fmt.Sprintf("%s:%s@tcp(%s)/%s", DBUser, DBPass, DBHost, DBDbase)
	var err error
	db, err = sql.Open("mysql", dbConn)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Database connected successfully")
}

// CloseDB closes the database connection
func CloseDB() {
	db.Close()
}

// Order represents an order
type Order struct {
	ID         int    `json:"id"`
	CustomerID int    `json:"customer_id"`
	ProductID  int    `json:"product_id"`
	Quantity   int    `json:"quantity"`
	OrderDate  string `json:"order_date,omitempty"`
}

// CreateOrder inserts a new order into the database
func CreateOrder(order Order) (int, error) {
	query := `INSERT INTO orders (customer_id, product_id, quantity, order_date) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(query, order.CustomerID, order.ProductID, order.Quantity, order.OrderDate)
	if err != nil {
		return 0, fmt.Errorf("failed to insert order: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get inserted order ID: %v", err)
	}
	return int(id), nil
}
