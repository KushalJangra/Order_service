package main

import (
	"fmt"
	"log"
	"net/http"

	"Service1f/controllers"
	"Service1f/model"
)

func main() {
	// Initialize database connection
	models.InitDB()
	defer models.CloseDB()

	// Set up routes
	http.HandleFunc("/orders", controllers.OrderHandler)

	fmt.Println("Order service running on port 8082...")
	log.Fatal(http.ListenAndServe(":8082", nil))
}