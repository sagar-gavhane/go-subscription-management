package services

import (
	"database/sql"
	"log"

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

func (repo *Repo) GetPlans() ([]types.Plan, error) {
	var plans []types.Plan

	query := `select * from plans`
	rows, err := repo.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var plan types.Plan
		var features []types.Feature

		if err := rows.Scan(&plan.Id, &plan.Name, &plan.Description, &plan.Type, &plan.IconUrl, &plan.ReleaseDate, &plan.EndDate); err != nil {
			log.Println(err)
		}

		query = `select * from plan_features where plan_features.plan_id = $1`
		rows2, err := repo.DB.Query(query, plan.Id)

		if err != nil {
			log.Println(err)
		}

		defer rows2.Close()

		for rows2.Next() {
			var feature types.Feature

			if err = rows2.Scan(&feature.Id, &feature.Feature, &feature.PlanId); err != nil {
				log.Println(err)
			}

			features = append(features, feature)
		}

		plan.Features = features

		plans = append(plans, plan)
	}

	return plans, nil
}

// TODO: GetPlanById
func (repo *Repo) GetPlanById(planId int) (types.Plan, error) {
	var plan types.Plan

	query := `select * from plans where plans.id = $1`
	rows, err := repo.DB.Query(query, planId)

	if err != nil {
		return plan, err
	}

	defer rows.Close()

	for rows.Next() {
		var features []types.Feature

		if err := rows.Scan(&plan.Id, &plan.Name, &plan.Description, &plan.Type, &plan.IconUrl, &plan.ReleaseDate, &plan.EndDate); err != nil {
			log.Println(err)
		}

		query = `select * from plan_features where plan_features.plan_id = $1`
		rows2, err := repo.DB.Query(query, planId)

		if err != nil {
			log.Println(err)
		}

		defer rows2.Close()

		for rows2.Next() {
			var feature types.Feature

			if err = rows2.Scan(&feature.Id, &feature.Feature, &feature.PlanId); err != nil {
				log.Println(err)
			}

			features = append(features, feature)
		}

		plan.Features = features
	}

	// if plan == nil {
	// 	err := errors.New(fmt.Sprintf("No plan found with id %v.", planId))

	// 	return nil, err
	// }

	return plan, nil
}

// TODO: createPlans
// TODO: updatePlansById
// TODO: deletePlansById
