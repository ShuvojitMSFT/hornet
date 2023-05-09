package main

import (
	// "database/sql"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type User struct {
	UserID   int    `json:"userid"`
	Username string `json:"username"`
	City     string `json:"city"`
	Email    string `json:"email"`
}

type JsonResponse struct {
	Type    string
	Data    []User
	Message string
}

// load .env file
func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

// connect to the database and return it as an object
// This is an exmple of connection leak
func dbConn() (db *sql.DB) {
	// pass the db credentials into variables
	host := goDotEnvVariable("DBHOST")
	port := goDotEnvVariable("DBPORT")
	dbUser := goDotEnvVariable("DBUSER")
	dbPass := goDotEnvVariable("DBPASS")
	dbname := goDotEnvVariable("DBNAME")
	// create a connection string
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, dbUser, dbPass, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return db
}

func main() {
	db := dbConn()

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	// Init the mux router
	router := mux.NewRouter()

	// Route handles & endpoints

	// Get all users
	router.HandleFunc("/users/", GetUsers).Methods("GET")

	// Create a user
	router.HandleFunc("/SignUp/", SignUp).Methods("POST")

	//login
	router.HandleFunc("/login/{user_id}/", Login).Methods("POST")

	// serve the app
	fmt.Println("Server at 10000")
	log.Fatal(http.ListenAndServe(":10000", router))

}

func Login(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["userid"]
	password := params["password"]

	var response = JsonResponse{}

	if userID == "" {
		response = JsonResponse{Type: "error", Message: "Enter a valid userID"}
	} else {
		var passcodematch string
		db := dbConn()
		getPassword := `SELECT "password" FROM public."users" WHERE id=$1`
		outPassword := db.QueryRow(getPassword, userID)
		err := outPassword.Scan(&passcodematch)
		if err != nil {
			panic(err)
		}

		if password == passcodematch {
			fmt.Println("The password is a match. You shall be supplied with a token shortly.")

		} else {
			fmt.Println("The password does not match")
		}

	}

	json.NewEncoder(w).Encode(response)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	fmt.Println("Getting users.")
	rows, err := db.Query(`SELECT "user_id","username","city","email" FROM public."users"`)

	if err != nil {
		panic(err)
	}

	var users []User

	// Foreach user
	for rows.Next() {
		var id int
		var username string
		var city string
		var email string

		err = rows.Scan(&id, &username, &city, &email)
		if err != nil {
			panic(err)
		}
		users = append(users, User{UserID: id, Username: username, City: city, Email: email})
		fmt.Println(id, username, city, email)
	}

	var response = JsonResponse{
		Type:    "success",
		Data:    users,
		Message: "Done!",
	}
	fmt.Printf("The response is: %+v", response)
	if json.NewEncoder(w).Encode(response) != nil {
		panic(err)
	}
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	user_id := r.FormValue("userID")
	username := r.FormValue("username")
	password := r.FormValue("password")
	city := r.FormValue("city")
	email := r.FormValue("email")

	var response = JsonResponse{}

	if user_id == "" || username == "" || email == "" || password == "" {
		response = JsonResponse{Type: "error", Message: "You are missing a mandatory field"}
	} else {
		db := dbConn()

		fmt.Println("Inserting new user with ID: " + user_id + " and name: " + username)

		var lastInsertID int
		err := db.QueryRow("INSERT INTO users(user_id, username, password, city, email) VALUES($1, $2, $3, $4, $5) returning user_id;", user_id, username, password, city, email).Scan(&lastInsertID)

		if err != nil {
			panic(err)
		}

		response = JsonResponse{Type: "success", Message: "User has been signed up."}
	}

	json.NewEncoder(w).Encode(response)
}
