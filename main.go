package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)


type hornetUser struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {

	dbConnStr := "host=34.72.216.95 port=5432 user=hornetapi password=KonamiCode001 dbname=hornet sslmode=disable"

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test the database connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new HTTP server
	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body into a User struct
		var user hornetUser
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Insert the user into the database
		stmt, err := db.Prepare("INSERT INTO hornet(name, email, password) VALUES($1, $2, $3) RETURNING id")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		var userID int
		err = stmt.QueryRow(user.Name, user.Email, user.Password).Scan(&userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with the new user's ID
		response := struct {
			ID int `json:"id"`
		}{userID}
		json.NewEncoder(w).Encode(response)
	})
	
	http.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body into a hornetUser struct
		var user hornetUser
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Query the database for the user's email and password
		row := db.QueryRow("SELECT id, name FROM users WHERE email = $1 AND password = $2", user.Email, user.Password)

		// Scan the row into a User struct
		var matchedUser User
		err = row.Scan(&matchedUser.ID, &matchedUser.Name)
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// Respond with the matched user's ID and name
		response := struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}{matchedUser.ID, matchedUser.Name}
		json.NewEncoder(w).Encode(response)
	})

	// Start the server
	log.Println("Server listening on 5432")
	err = http.ListenAndServe("34.72.216.95:5432", nil)
	if err != nil {
		log.Fatal(err)
	}
}
