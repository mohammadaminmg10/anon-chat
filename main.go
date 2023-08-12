package main

import (
	"anon-chat/config"
	"anon-chat/handlers"
	"anon-chat/registration"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func setupDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "dbname=AnonChat user=root password=DiKgtuqxrJU9cdJOZgtnJull host=alfie.iran.liara.ir port=34009 sslmode=disable")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	DB, err := setupDB()
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()
	configFilename := "config/config.json"
	cfg, err := config.LoadConfig(configFilename)
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}

	fmt.Println("Database Host:", cfg.Database.Host)
	fmt.Println("Database Port:", cfg.Database.Port)
	fmt.Println("JWT Secret:", cfg.Jwt.JWTKey)
	// Assuming your static files are in the "static" directory
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", handlers.HandleIndex)

	http.HandleFunc("/send", func(writer http.ResponseWriter, request *http.Request) { handlers.HandleSend(writer, request, DB) })

	http.HandleFunc("/register", func(writer http.ResponseWriter, request *http.Request) {
		registration.HandleRegister(writer, request, DB)
	})

	http.HandleFunc("/user/", func(writer http.ResponseWriter, request *http.Request) { handlers.HandleForm(writer, request, DB) })

	fmt.Println("Server listening on port 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
