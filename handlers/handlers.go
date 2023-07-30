package handlers

import (
	"anon-chat/models"
	_ "fmt"
	"github.com/google/uuid"
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseFiles("views/index.html"))

func generateUserID(w http.ResponseWriter) string {
	userID := uuid.New().String()

	cookie := http.Cookie{
		Name:     "user_id",
		Value:    userID,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)

	return userID
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
	if r.Method == http.MethodGet {
		// Get the user identifier (UUID) from the cookie.
		userID := getUserID(r)

		// Generate a new user identifier and store it in the cookie if not found.
		if userID == "" {
			userID = generateUserID(w)
		}

		// Get all messages associated with the current user's identifier.
		messages := models.GetMessagesByUserID(userID)

		// Render the template with the messages and the current user identifier.
		data := struct {
			Messages      []models.Message
			CurrentUserID string
		}{
			Messages:      messages,
			CurrentUserID: userID,
		}

		err := templates.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return // Return here to avoid the second call to response.WriteHeader()
	}

	// In case the method is not GET, return an error response.
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

func HandleSend(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		nickname := r.FormValue("nickname")
		messageText := r.FormValue("message")

		userID := getUserID(r) // Get the user's unique identifier from the cookie.

		message := models.Message{
			UserID:   userID, // Associate the message with the user's unique identifier.
			Nickname: nickname,
			Text:     messageText,
		}

		models.SaveMessage(message)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
