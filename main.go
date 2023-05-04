package main

import (
	"database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"

    _ "github.com/lib/pq"
    "github.com/gorilla/mux"
    "github.com/gorilla/context"
)


type hornetUserSignin struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type hornetUserSignup struct {
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

	// Create the signup api for the client.
	
	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body into a User struct
		var user hornetUserSignup
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
	
	
	// This part will expose the signin api to the client.
	
	http.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body into a hornetUser struct
		var user hornetUserSignin
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Query the database for the user's email and password
		row := db.QueryRow("SELECT id, name FROM users WHERE email = $1 AND password = $2", user.Email, user.Password)

		
		var matchedUser hornetUserSignin
		err = row.Scan(&matchedUser.ID, &matchedUser.Name)
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		
		response := struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}{matchedUser.ID, matchedUser.Name}
		json.NewEncoder(w).Encode(response)
	})
	
	



// This part is if the server is on the local host. I am using a GCP instance of postgre

/*		
	// Start the server
	log.Println("Server listening on 5432")
	err = http.ListenAndServe("34.72.216.95:5432", nil)
	if err != nil {
		log.Fatal(err)
	}
	*/
}
