package models

import "time"

// Message represents an anonymous chat message.
type Message struct {
	UserID    string    // The unique identifier of the user who sent the message.
	Nickname  string    // The nickname of the user (optional).
	Text      string    // The content of the message.
	Timestamp time.Time // The timestamp when the message was created.
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
