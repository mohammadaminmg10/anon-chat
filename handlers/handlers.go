package handlers

import (
	"anon-chat/models"
	"anon-chat/registration"
	"database/sql"
	_ "fmt"
	"html/template"
	"net/http"
	"strings"
)

var templates = template.Must(template.ParseFiles("views/index.html"))
var username string

func HandleIndex(w http.ResponseWriter, r *http.Request) {

	userID := registration.GetUserID(r)

	if registration.IsLoggedIn(r) {
		// If the user is logged in, redirect to their own chat room
		http.Redirect(w, r, "/user/"+userID, http.StatusSeeOther)
		return
	}

	// If the user is not logged in, redirect to the registration page
	http.Redirect(w, r, "/register", http.StatusSeeOther)
}

func HandleSend(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodPost {
		nickname := r.FormValue("nickname")
		messageText := r.FormValue("message")

		message := models.Message{
			UserID:   username,
			Nickname: nickname,
			Text:     messageText,
		}

		err := models.SaveMessage(db, message)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/user/"+username, http.StatusSeeOther)
		return
	}
}

func HandleForm(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userID := registration.GetUserID(r)

	// Extract the username from the URL
	username = strings.TrimPrefix(r.URL.Path, "/user/")

	if username == userID {
		// The user is the owner of the chat room, show chat history
		messages, err := models.GetMessagesByUsername(db, username)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		data := struct {
			Username string
			Messages []models.Message
			IsOwner  bool
			IsUser   bool
		}{
			Username: username,
			Messages: messages,
			IsOwner:  true,
			IsUser:   false,
		}

		err = templates.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	// The user is not the owner, show the message form
	data := struct {
		Username string
		IsOwner  bool
		IsUser   bool
	}{
		Username: username,
		IsOwner:  false,
		IsUser:   true,
	}

	err := templates.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	return
}
