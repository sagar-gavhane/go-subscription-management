package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/sagar-gavhane/go-subscription-management/internal/helpers"
	"github.com/sagar-gavhane/go-subscription-management/internal/services"
	"github.com/sagar-gavhane/go-subscription-management/internal/types"
	//"github.com/pelletier/go-toml/query"
)

type Repo struct {
	DB *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{
		DB: db,
	}
}

type Feature struct {
	Id      int    `json:"id"`
	Feature string `json:"feature"`
	PlanId  int    `json:"plan_id"`
}

type Plan struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	IconUrl     string    `json:"icon_url"`
	Features    []Feature `json:"features"`
	ReleaseDate string    `json:"release_date"`
	EndDate     string    `json:"end_date"`
}

func (repo *Repo) Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("index route"))
}

func (repo *Repo) PlansHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	_services := services.NewRepo(repo.DB)

	if r.Method == "GET" {
		var planId int
		planId, err = strconv.Atoi(r.URL.Query().Get("plan_id"))

		helpers.LogError(err)

		if planId > 0 {
			plan, err := _services.GetPlanById(planId)

			if err != nil {
				response := types.Response{
					Message: err.Error(),
					Data:    nil,
				}

				jsonData, err := json.Marshal(response)

				helpers.LogError(err)

				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(string(jsonData)))
				return
			}

			if plan.Id == 0 {
				response := types.Response{
					Message: errors.New(fmt.Sprintf("No plan found with id %v.", planId)).Error(),
					Data:    nil,
				}

				jsonData, err := json.Marshal(response)

				helpers.LogError(err)

				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(string(jsonData)))
				return
			}

			response := types.Response{
				Message: "All plans retrieved successfully.",
				Data:    plan,
			}

			jsonData, err := json.Marshal(response)

			helpers.LogError(err)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(string(jsonData)))
			return
		}

		plans, err := _services.GetPlans()

		if err != nil {
			fmt.Println(err)
			return
		}

		jsonData, err := json.Marshal(plans)

		if err != nil {
			log.Fatal(err)
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(string(jsonData)))
		return
	}

	var plan Plan

	if r.Method == "POST" {
		err = json.NewDecoder(r.Body).Decode(&plan)

		if err != nil {
			fmt.Println(err)
			return
		}

		var insertedId int
		query := `insert into "plans" ("name", "description", "type", "icon_url", "release_date", "end_date") values($1, $2, $3, $4, $5, $6) RETURNING id`
		err = repo.DB.QueryRow(query, plan.Name, plan.Description, plan.Type, plan.IconUrl, plan.ReleaseDate, plan.EndDate).Scan(&insertedId)

		for _, feature := range plan.Features {
			insertStmt := `insert into "plan_features" ("feature", "plan_id") values($1, $2)`
			_, err = repo.DB.Exec(insertStmt, feature, insertedId)
		}

		if err != nil {
			fmt.Println(err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("plan inserted successfully."))
		return
	}

	if r.Method == "DELETE" {
		fmt.Println(r.URL.Path)
		// 1. check plan_id is passed or not
		// 2. plan is exist in db?
		// 3. delete plan_features
		w.Write([]byte("delete plan"))
	}

	fmt.Println(r.Method)

	w.Write([]byte("working!!"))
}
