package main

import (
	"anon-chat/handlers"
	"anon-chat/registration"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func setupDB() (*sql.DB, error) {
	connStr := "postgres://moamin:1TLkY3LVPsffik92zVZf33eiXv2wacMM@dpg-cjad7f6e546c738c72o0-a.oregon-postgres.render.com/anonchat"
	db, err := sql.Open("postgres", connStr)
	//	db, err := sql.Open("postgres", "dbname=anonchat user=moamin password=1TLkY3LVPsffik92zVZf33eiXv2wacMM host=dpg-cjad7f6e546c738c72o0-a port=5432 sslmode=disable")
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
	_, err = DB.Exec(`
       CREATE SCHEMA anonchat;
    `)
	// Create the "messages" table
	_, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS AnonChat.messages (
            id SERIAL PRIMARY KEY,
            user_id VARCHAR(50) NOT NULL,
            nickname VARCHAR(50),
            text TEXT NOT NULL,
            timestamp TIMESTAMPTZ NOT NULL
        );
    `)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Messages table created successfully.")

	// Create the "users" table
	_, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS AnonChat.users (
            id SERIAL PRIMARY KEY,
            username VARCHAR(100) NOT NULL UNIQUE,
            password VARCHAR(255) NOT NULL
        );
    `)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Users table created successfully.")
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
