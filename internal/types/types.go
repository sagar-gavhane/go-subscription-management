package types

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

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
