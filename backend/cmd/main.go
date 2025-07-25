package main

import (
	"fmt"
	"google-auth-demo/backend/internal/httpserver"
	"google-auth-demo/backend/internal/service"
)

func main() {
	fmt.Println("Start the program")

	// db, err := postgresdb.New()
	// if err != nil {
	// 	fmt.Println("Error initializing PostgresDB:", err)
	// 	return
	// }

	svc := service.New()

	server := httpserver.New(svc)

	server.Start() // handle error
}
