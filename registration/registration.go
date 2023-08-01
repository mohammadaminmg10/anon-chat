package registration

import (
	"github.com/google/uuid"
	"html/template"
	"net/http"
)

var registerTemplate = template.Must(template.ParseFiles("views/register.html"))

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

// getUserID retrieves the user's username from the session (cookie).
func getUserID(r *http.Request) string {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}

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

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Check if the user is already logged in
		if isLoggedIn(r) {
			// If the user is already logged in, redirect to their chat room
			http.Redirect(w, r, "/user/"+getUserID(r), http.StatusSeeOther)
			return
		}

		// If the user is not logged in, display the user registration page
		err := registerTemplate.Execute(w, nil)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	if r.Method == http.MethodPost {
		// Handle user registration form submission
		username := r.FormValue("username")

		// Perform any necessary validation on the username (e.g., check if it's unique, alphanumeric, etc.)
		// For simplicity, let's assume the username is always valid and unique for now.

		// Save the username and its associated user identifier in the session (cookie-based)
		setUserID(w, username)

		// Redirect the user to their chat room page
		http.Redirect(w, r, "/user/"+username, http.StatusSeeOther)
	}

	// In case the method is not GET or POST, return an error response.
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}
