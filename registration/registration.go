package registration

import (
	"github.com/golang-jwt/jwt/v4"
	"html/template"
	"net/http"
	"time"
)

var registerTemplate = template.Must(template.ParseFiles("views/register.html"))

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.RegisteredClaims
}

const jwtCookieName = "jwt_token"

var jwtKey = []byte("n7hUjNdqoF9q5jSk3zDpW2vR1zHtY8dA")

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
func IsLoggedIn(r *http.Request) bool {
	cookie, err := r.Cookie(cookieName)
	if err == nil && cookie != nil {
		return true
	}
	return false
}

// getUserID retrieves the user's username from the session (cookie).
func GetUserID(r *http.Request) string {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Check if the user is already logged in
		if IsLoggedIn(r) {
			// If the user is already logged in, redirect to their chat room
			http.Redirect(w, r, "/user/"+GetUserID(r), http.StatusSeeOther)
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
		password := r.FormValue("password")

		// Perform any necessary validation on the username (e.g., check if it's unique, alphanumeric, etc.)
		// For simplicity, let's assume the username is always valid and unique for now.

		// Generate a JWT token for the user
		expirationTime := time.Now().Add(5 * time.Minute)
		claims := &Claims{
			Username: username,
			Password: password,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Set the JWT token as a cookie
		http.SetCookie(w, &http.Cookie{
			Name:     jwtCookieName,
			Value:    tokenString,
			Expires:  expirationTime,
			HttpOnly: true,
		})

		// Set the username in the user_cookie and redirect to their chat room
		setUserID(w, username)
		http.Redirect(w, r, "/user/"+username, http.StatusSeeOther)
		return
	}

	// In case the method is not GET or POST, return an error response.
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}
