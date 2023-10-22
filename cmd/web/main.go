package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sagar-gavhane/go-subscription-management/internal/handlers"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "sagar"
	password = "Sagar@123"
	dbname   = "subscription_mgmt_db"
)

func main() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	router := mux.NewRouter()

	_handlers := handlers.NewRepo(db)

	router.HandleFunc("/api/plans", _handlers.GetPlans).Methods("GET")
	router.HandleFunc("/api/plans", _handlers.CreatePlan).Methods("POST")
	router.HandleFunc("/api/plans/{planId}", _handlers.GetPlanById).Methods("GET")
	router.HandleFunc("/api/plans/{planId}", _handlers.UpdatePlanById).Methods("PATCH")
	router.HandleFunc("/api/plans/{planId}", _handlers.DeletePlanById).Methods("DELETE")
	router.HandleFunc("/api/plans/{planId}/deactivate", _handlers.DeactivatePlanById).Methods("POST")
	router.HandleFunc("/", _handlers.Index)

	fmt.Println("application started on :3000")

	err = http.ListenAndServe(":3000", router)

	if err != nil {
		log.Fatal("failed to start server", err)
	}
}
