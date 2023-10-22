package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sagar-gavhane/go-subscription-management/internal/helpers"
	"github.com/sagar-gavhane/go-subscription-management/internal/services"
	"github.com/sagar-gavhane/go-subscription-management/internal/types"
)

type Repo struct {
	DB *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{
		DB: db,
	}
}

func (repo *Repo) Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("index route"))
}

func (repo *Repo) CreatePlan(w http.ResponseWriter, r *http.Request) {
	log.Println("handlers.CreatePlan executed")
	var err error
	services := services.NewRepo(repo.DB)
	var plan types.Plan

	err = json.NewDecoder(r.Body).Decode(&plan)

	log.Printf("handlers.CreatePlan req.body %+v\n", plan)

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

	log.Printf("handlers.CreatePlan plan inserted successfully and plan = %d\n", insertedId)

	newPlan, err := services.GetPlanById(insertedId)

	helpers.LogError(err)

	jsonData, err := json.Marshal(newPlan)

	helpers.LogError(err)

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(string(jsonData)))
	w.WriteHeader(http.StatusCreated)
	log.Println("handler.CreatePlan end")
}

func (repo *Repo) GetPlans(w http.ResponseWriter, r *http.Request) {
	log.Println("handlers.GetPlans executed")
	services := services.NewRepo(repo.DB)

	plans, err := services.GetPlans()

	if err != nil {
		fmt.Println(err)
		return
	}

	jsonData, err := json.Marshal(plans)

	if err != nil {
		helpers.LogError(err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(string(jsonData)))
	log.Println("handlers.GetPlans end")
}

func (repo *Repo) GetPlanById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	services := services.NewRepo(repo.DB)

	planId, err := strconv.Atoi(vars["planId"])

	if err != nil {
		w.Write([]byte("invalid plan id passed"))
		return
	}

	plan, err := services.GetPlanById(planId)

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
		Message: "Plan retrieved successfully.",
		Data:    plan,
	}

	jsonData, err := json.Marshal(response)

	helpers.LogError(err)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(string(jsonData)))
}

func (repo *Repo) UpdatePlanById(w http.ResponseWriter, r *http.Request) {
	_services := services.NewRepo(repo.DB)
	var plan types.Plan

	// 1. check plan exists or not in DB
	// 2. if yes, then patch new values
	planId, err := strconv.Atoi(r.URL.Query().Get("plan_id"))

	helpers.LogError(err)

	if planId == 0 {
		w.Write([]byte("invalid plan id received."))
		return
	}

	err = json.NewDecoder(r.Body).Decode(&plan)

	_services.UpdatePlanById(planId, plan)

	w.Write([]byte("patch plan"))

}

func (repo *Repo) DeletePlanById(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		fmt.Println(r.URL.Path)
		// 1. check plan_id is passed or not
		// 2. plan is exist in db?
		// 3. delete plan_features
	}

	w.Write([]byte("delete plan"))
}

func (repo *Repo) DeactivatePlanById(w http.ResponseWriter, r *http.Request) {
	var plan types.Plan

	json.NewDecoder(r.Body).Decode(&plan)

	w.Write([]byte("deactivate plan by id"))
}
