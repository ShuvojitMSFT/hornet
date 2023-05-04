package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)




type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Content struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func main() {
	// Establish a database connection
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create a new router
	router := mux.NewRouter()

	// Handle sign up requests
	router.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Insert the user into the database
		_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, string(hashedPassword))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Return a success response
		w.WriteHeader(http.StatusCreated)
	}).Methods("POST")

	// Handle login requests
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Retrieve the user from the database
		var hashedPassword string
		err = db.QueryRow("SELECT password FROM users WHERE username = $1", user.Username).Scan(&hashedPassword)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Compare the password with the hash
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Generate a JWT token
		expirationTime := time.Now().Add(24 * time.Hour)
		claims := &Claims{
			Username: user.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			
			
			
			
			
			
			
