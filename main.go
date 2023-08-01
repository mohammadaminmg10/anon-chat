package main

import (
	"anon-chat/handlers"
	"anon-chat/registration"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlers.HandleIndex)
	http.HandleFunc("/send", handlers.HandleSend)
	http.HandleFunc("/register", registration.HandleRegister)

	fmt.Println("Server listening on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
