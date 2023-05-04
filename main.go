package main

import (
    "database/sql"
    "fmt"
    "log"


)

func main() {
    // Replace the values in <angle brackets> with your own details
    connStr := "host=34.72.216.95 port=5432 user=hornetapi password=KonamiCode001 dbname=hornet sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Connected to PostgreSQL instance on GCP!")
}
