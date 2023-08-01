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
	http.HandleFunc("/usr/", handlers.HandleIndex)

	fmt.Println("Server listening on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
