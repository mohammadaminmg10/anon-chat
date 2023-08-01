package models

import (
	"time"
)

// Message represents an anonymous chat message.
type Message struct {
	UserID    string    // The unique identifier of the user who sent the message.
	Nickname  string    // The nickname of the user (optional).
	Text      string    // The content of the message.
	Timestamp time.Time // The timestamp when the message was created.
	Username  string
}

var messages []Message

// SaveMessage saves a new message to the list of messages.
func SaveMessage(message Message) {
	message.Timestamp = time.Now()
	messages = append(messages, message)
}

// GetMessages returns the list of all messages.
func GetMessages() []Message {
	return messages
}

func GetMessagesByUserID(userID string) []Message {
	var userMessages []Message
	for _, message := range messages {
		if message.UserID == userID {
			userMessages = append(userMessages, message)
		}
	}
	return userMessages
}

func GetMessagesByUsername(username string) []Message {
	var userMessages []Message
	for _, message := range messages {
		if message.Username == username {
			userMessages = append(userMessages, message)
		}
	}
	return userMessages
}

//func HandleIndex(w http.ResponseWriter, r *http.Request) {
//	// Check if the user is logged in
//	if !registration.IsLoggedIn(r) {
//		http.Redirect(w, r, "/register", http.StatusSeeOther)
//		return
//	}
//
//	userID := registration.GetUserID(r)
//
//	if strings.HasPrefix(r.URL.Path, "/user/") {
//		// Extract the username from the URL
//		username := strings.TrimPrefix(r.URL.Path, "/user/")
//
//		// Check if the logged-in user is the owner of the chat room
//		isOwner := (userID == username)
//
//		// Get all messages associated with the specified username
//		chatRoomMessages := models.GetMessagesByUsername(username)
//
//		data := struct {
//			Messages []models.Message
//			Username string
//			IsUser   bool
//			IsOwner  bool
//		}{
//			Messages: chatRoomMessages,
//			Username: username,
//			IsUser:   true,
//			IsOwner:  isOwner,
//		}
//
//		err := templates.Execute(w, data)
//		if err != nil {
//			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
//			return
//		}
//		return
//	}
//
//	// If the path is not "/user/username", display the main chat page
//	// with the message submission form.
//
//	messages := models.GetMessages()
//
//	data := struct {
//		Messages []models.Message
//		Username string
//		IsUser   bool
//	}{
//		Messages: messages,
//		Username: userID,
//		IsUser:   false,
//	}
//
//	err := templates.Execute(w, data)
//	if err != nil {
//		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
//		return
//	}
//}
