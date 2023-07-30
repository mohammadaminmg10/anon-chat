package handlers

import (
	"fmt"
	"net/http"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "views/index.html")
	}
}

func HandleSend(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		message := r.FormValue("message")
		// process the 'message'
		fmt.Println("Received anonymous message:", message)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
