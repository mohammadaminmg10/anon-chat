package handlers

import (
	"anon-chat/models"
	"anon-chat/registration"
	_ "fmt"
	"html/template"
	"net/http"
	"strings"
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
	userID := registration.GetUserID(r)

	if strings.HasPrefix(r.URL.Path, "/usr/") {
		// Extract the username from the URL
		username := strings.TrimPrefix(r.URL.Path, "/usr/")

		// Check if the logged-in user is the owner of the chat room or if it's an anonymous visitor
		isOwner := (userID != "" && userID == username)

		var messages []models.Message
		if isOwner {
			// If the user is the owner, fetch all messages associated with their username
			messages = models.GetMessagesByUsername(username)
		} else if userID != "" {
			// If the user is logged in (not anonymous), fetch messages they sent themselves
			messages = models.GetMessagesByUserID(userID)
		}

		data := struct {
			Username string
			IsOwner  bool
			Messages []models.Message
		}{
			Username: username,
			IsOwner:  isOwner,
			Messages: messages,
		}

		err := templates.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	// If the path is not "/usr/username", it's an anonymous visitor trying to send a message.
	// We'll display the chat history for all users and the chat field for anonymous messages.

	data := struct {
		Messages []models.Message
		Username string
		IsUser   bool
	}{
		Messages: models.GetMessages(),
		Username: userID,
		IsUser:   (userID != ""), // Check if the user is logged in or anonymous
	}

	err := templates.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func HandleSend(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		nickname := r.FormValue("nickname")
		messageText := r.FormValue("message")

		// Check if the user is logged in (not anonymous)
		userID := registration.GetUserID(r)

		message := models.Message{
			UserID:   userID, // Associate the message with the user's unique identifier.
			Nickname: nickname,
			Text:     messageText,
		}

		models.SaveMessage(message)

		// If the user is anonymous (not logged in), redirect to the chat history page
		if userID == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// If the user is logged in, redirect back to the same chat room page with their messages
		username := "mo@amin"
		http.Redirect(w, r, "/usr/"+username, http.StatusSeeOther)
	}
}
