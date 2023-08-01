package handlers

import (
	"anon-chat/models"
	"anon-chat/registration"
	_ "fmt"
	"html/template"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseFiles("views/index.html"))

const cookieName = "user_cookie"

// setUserID stores the user's username in a cookie (session).
func setUserID(w http.ResponseWriter, username string) {
	cookie := http.Cookie{
		Name:     cookieName,
		Value:    username,
		HttpOnly: true,
		Path:     "/",
		// Add additional secure options as needed, such as Secure and SameSite,
		// depending on your deployment environment.
	}

	http.SetCookie(w, &cookie)
}

// isLoggedIn checks if the user is already logged in by checking if the username cookie exists.
func isLoggedIn(r *http.Request) bool {
	cookie, err := r.Cookie(cookieName)
	if err == nil && cookie != nil {
		return true
	}
	return false
}

// getUserID retrieves the user identifier (UUID) from the cookie.
func getUserID(r *http.Request) string {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		// Check if the user is logged in
		if !isLoggedIn(r) {
			// If the user is not logged in, redirect to the registration page
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		// Get the username from the session (cookie)
		username := getUserID(r)

		// Get all messages associated with the current user's username
		chatRoomMessages := models.GetMessagesByUsername(username)

		// Render the template with the messages and the current user's username
		data := struct {
			Messages []models.Message
			Username string
		}{
			Messages: chatRoomMessages,
			Username: username,
		}

		err := templates.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	// In case the path is not "/", return an error response.
	http.Error(w, "Page Not Found", http.StatusNotFound)
}

func HandleSend(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		nickname := r.FormValue("nickname")
		messageText := r.FormValue("message")

		// Extract useful information from the request
		userAgent := r.UserAgent()
		ipAddress := r.RemoteAddr
		// You can access other headers as needed: r.Header.Get("Header-Name")

		userID := getUserID(r) // Get the user's unique identifier from the cookie.

		// Log the received message and associated information
		log.Printf("Received message from UserID: %s, Nickname: %s, IP: %s, UserAgent: %s",
			userID, nickname, ipAddress, userAgent)

		message := models.Message{
			UserID:   userID, // Associate the message with the user's unique identifier.
			Nickname: nickname,
			Text:     messageText,
		}

		models.SaveMessage(message)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	registration.HandleRegister(w, r)
}
