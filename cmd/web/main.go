package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

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

	mux := http.NewServeMux()

	_handlers := handlers.NewRepo(db)

	mux.HandleFunc("/api/plans", _handlers.PlansHandler)
	mux.HandleFunc("/", _handlers.Index)

	fmt.Println("application started on :3000")

	err = http.ListenAndServe(":3000", mux)

	if err != nil {
		log.Fatal("failed to start server", err)
	}
}
