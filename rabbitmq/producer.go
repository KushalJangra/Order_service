package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

const (
	QueueName   = "orderQueue"
	RabbitMQURL = "amqp://guest:guest@localhost:5672/"
)

type Order struct {
	ID         int    `json:"id"`
	CustomerID int    `json:"customer_id"`
	ProductID  int    `json:"product_id"`
	Quantity   int    `json:"quantity"`
	OrderDate  string `json:"order_date"`
}

func main() {
	
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


	orders := []Order{
		{CustomerID: 23, ProductID: 2, Quantity: 1, OrderDate: time.Now().Format("2006-01-02")},
		{CustomerID: 24, ProductID: 1, Quantity: 9, OrderDate: time.Now().Format("2006-01-02")},
	}


	for i, order := range orders {
		order.ID = i + 1 // Assign an order ID
		orderJSON, err := json.Marshal(order) // Convert order to JSON
		if err != nil {
			log.Fatalf("Failed to marshal order to JSON: %v", err)
		}

		
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        orderJSON,
			},
		)
		if err != nil {
			log.Fatalf("Failed to publish order: %v", err)
		}

		
		fmt.Printf("Order %d sent: %s\n", order.ID, string(orderJSON))
	}
}
