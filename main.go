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

func setupDB(cfg config.DBConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	configFilename := "config/config.json"
	cfg, err := config.LoadConfig(configFilename)
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}

	DB, err := setupDB(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	// Assuming your static files are in the "static" directory
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", handlers.HandleIndex)

	http.HandleFunc("/send", func(writer http.ResponseWriter, request *http.Request) {
		handlers.HandleSend(writer, request, DB)
	})

	http.HandleFunc("/register", func(writer http.ResponseWriter, request *http.Request) {
		registration.HandleRegister(writer, request, DB, *cfg)
	})

	http.HandleFunc("/user/", func(writer http.ResponseWriter, request *http.Request) {
		handlers.HandleForm(writer, request, DB, *cfg)
	})

	fmt.Println("Server listening on port 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
