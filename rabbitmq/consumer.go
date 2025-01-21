package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/streadway/amqp"
	_ "github.com/go-sql-driver/mysql" 
)

const (
	DBHost         = "127.0.0.1"
	DBUser         = "root"
	DBPass         = "Kush@123456"
	DBDbase        = "pro"
	BaseURLCustomer = "http://localhost:8080/customers"
	BaseURLProduct  = "http://localhost:8080/products"
	RabbitMQURL    = "amqp://guest:guest@localhost:5672/"
	QueueName      = "orderQueue"
)

var db *sql.DB


type Order struct {
	ID         int    `json:"id"`
	CustomerID int    `json:"customer_id"`
	ProductID  int    `json:"product_id"`
	Quantity   int    `json:"quantity"`
	OrderDate  string `json:"order_date,omitempty"`
}

func initDB() {
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

func CreateOrder(order Order) error {
	query := `INSERT INTO orders (customer_id, product_id, quantity, order_date) VALUES (?, ?, ?, ?)`
	_, err := db.Exec(query, order.CustomerID, order.ProductID, order.Quantity, order.OrderDate)
	if err != nil {
		return fmt.Errorf("failed to insert order: %v", err)
	}
	return nil
}

func FetchCustomer(customerID int) (bool, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%d", BaseURLCustomer, customerID))
	if err != nil {
		return false, fmt.Errorf("failed to fetch customer: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("customer not found: %v", resp.Status)
	}

	return true, nil
}

func FetchProduct(productID int) (bool, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%d", BaseURLProduct, productID))
	if err != nil {
		return false, fmt.Errorf("failed to fetch product: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("product not found: %v", resp.Status)
	}

	return true, nil
}

func ProcessOrder(order Order) {

	customerExists, err := FetchCustomer(order.CustomerID)
	if err != nil || !customerExists {
		log.Printf("Customer not found: %v", err)
		return
	}

	
	productExists, err := FetchProduct(order.ProductID)
	if err != nil || !productExists {
		log.Printf("Product not found: %v", err)
		return
	}

	
	err = CreateOrder(order)
	if err != nil {
		log.Printf("Failed to create order: %v", err)
		return
	}

	log.Printf("Order processed successfully: %+v", order)
}

func StartConsumer() {
	conn, err := amqp.Dial(RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		QueueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Println("Waiting for messages...")
	for msg := range msgs {
		var order Order
		err := json.Unmarshal(msg.Body, &order)
		if err != nil {
			log.Printf("Failed to decode order: %v", err)
			continue
		}

		log.Printf("Received order: %+v", order)
		ProcessOrder(order)
	}
}

func main() {
	initDB()
	defer db.Close()

	go StartConsumer()

	fmt.Println("Order consumer running...")
	select {} // Keep the main goroutine alive
}
