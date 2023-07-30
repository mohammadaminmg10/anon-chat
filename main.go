package main

import (
	"anon-chat/handlers"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlers.HandleIndex)
	http.HandleFunc("/send", handlers.HandleSend)

	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
