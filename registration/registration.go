package registration

import (
	"anon-chat/config"
	"database/sql"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
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

// setUserID stores the user's username in a cookie (session).
func setUserID(w http.ResponseWriter, username string, config config.Configuration) {
	cookie := http.Cookie{
		Name:     config.Cookie.Name,
		Value:    username,
		HttpOnly: true,
		Path:     "/",
	}

	http.SetCookie(w, &cookie)
}

// isLoggedIn checks if the user is already logged in by checking if the username cookie exists.
func IsLoggedIn(r *http.Request, config config.Configuration) bool {
	cookie, err := r.Cookie(config.Cookie.Name)
	if err == nil && cookie != nil {
		return true
	}
	return false
}

// getUserID retrieves the user's username from the session (cookie).
func GetUserID(r *http.Request, config config.Configuration) string {
	cookie, err := r.Cookie(config.Cookie.Name)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func HandleRegister(w http.ResponseWriter, r *http.Request, db *sql.DB, configuration config.Configuration) {
	if r.Method == http.MethodGet {
		// Check if the user is already logged in
		if IsLoggedIn(r, configuration) {
			// If the user is already logged in, redirect to their chat room
			http.Redirect(w, r, "/user/"+GetUserID(r, configuration), http.StatusSeeOther)
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

		// Check if the user exists and authenticate
		if AuthenticateUser(db, username, password) {
			// If the user exists and password is correct, log them in
			GenerateJWT(w, username, password, configuration)
			setUserID(w, username, configuration)
			http.Redirect(w, r, "/usr/"+username, http.StatusSeeOther)
			return
		}

		// Check if the username is unique and password is valid for registration
		if IsUniqueUsername(db, username) && len(password) >= 4 {
			// If username is unique and password is valid, register the user
			RegisterUser(db, username, password)
			// Log in the user after registration
			GenerateJWT(w, username, password, configuration)
			setUserID(w, username, configuration)
			http.Redirect(w, r, "/usr/"+username, http.StatusSeeOther)
			return
		}

		// Set the username in the user_cookie and redirect to their chat room
		// If neither authentication nor registration conditions are met
		errorMessage := "Invalid username or password"
		http.Redirect(w, r, "/register?error="+errorMessage, http.StatusSeeOther)
		return
	}

	// In case the method is not GET or POST, return an error response.
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// Check if the given username is unique in the database
func IsUniqueUsername(db *sql.DB, username string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM AnonChat.users WHERE username = $1", username).Scan(&count)
	if err != nil {
		fmt.Println("Error querying database:", err)
		return false
	}
	return count == 0
}

func RegisterUser(db *sql.DB, username, password string) error {
	// Hash the password before storing it in the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO AnonChat.users (username, password) VALUES ($1, $2)", username, hashedPassword)
	if err != nil {
		fmt.Println("Error inserting into database:", err)
		return err
	}

	return nil
}

func AuthenticateUser(db *sql.DB, username, password string) bool {
	var hashedPassword string

	// Query the database to retrieve the hashed password for the given username
	err := db.QueryRow("SELECT password FROM AnonChat.users WHERE username = $1", username).Scan(&hashedPassword)
	if err != nil {
		fmt.Println("Error querying database:", err)
		return false
	}

	// Compare the provided password with the hashed password from the database
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		// Passwords do not match
		return false
	}

	// Passwords match, authentication successful
	return true
}

func GenerateJWT(w http.ResponseWriter, username, password string, config config.Configuration) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: username,
		Password: password,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Jwt.JWTKey))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the JWT token as a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     config.Cookie.Name,
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
	})

}
