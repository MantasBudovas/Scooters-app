package main

import (
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const ApiKey = "APIkey"

func main() {
	// Gets Database connection string
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://scooter:scooter@db:5432/scooters?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Checks if the database is reachable
	if err := db.Ping(); err != nil {
		log.Fatal("Unable to connect to the database:", err)
	}

	storage := &DBStorage{DB: db}
	rand.Seed(time.Now().UnixNano())

	// Launches client simulation after a delay
	go func() {
		time.Sleep(2 * time.Second)
		StartSimulation()
	}()

	// Registers API endpoints
	http.HandleFunc("/event", apiKeyMiddleware(EventHandlerDB(storage)))
	http.HandleFunc("/scooters", apiKeyMiddleware(ScootersHandler(storage)))

	log.Println("Server launched on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Function to check static API key
func apiKeyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != ApiKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
