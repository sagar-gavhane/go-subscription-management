package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sagar-gavhane/go-subscription-management/internal/handlers"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("POSTGRESQL_HOST")
	user := os.Getenv("POSTGRESQL_USER")
	password := os.Getenv("POSTGRESQL_PASSWORD")
	dbname := os.Getenv("POSTGRESQL_DBNAME")

	port, err := strconv.Atoi(os.Getenv("POSTGRESQL_PORT"))

	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to parse port number %v", port))
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Println("Database connected successfully.")

	router := mux.NewRouter()

	_handlers := handlers.NewRepo(db)

	router.HandleFunc("/api/plans", _handlers.GetPlans).Methods("GET")
	router.HandleFunc("/api/plans", _handlers.CreatePlan).Methods("POST")
	router.HandleFunc("/api/plans/{planId}", _handlers.GetPlanById).Methods("GET")
	router.HandleFunc("/api/plans/{planId}", _handlers.UpdatePlanById).Methods("PATCH")
	router.HandleFunc("/api/plans/{planId}", _handlers.DeletePlanById).Methods("DELETE")
	router.HandleFunc("/api/plans/{planId}/deactivate", _handlers.DeactivatePlanById).Methods("POST")
	router.HandleFunc("/", _handlers.Index)

	port, err = strconv.Atoi(os.Getenv("PORT"))

	log.Printf("Application started on :%v", port)

	err = http.ListenAndServe(fmt.Sprintf(":%v", port), router)

	if err != nil {
		log.Fatal("failed to start server", err)
	}
}
