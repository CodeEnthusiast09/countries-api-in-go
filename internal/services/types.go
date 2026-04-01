package services

type StatusResponse struct {
	Total           int64   `json:"total_countries"`
	LastRefreshedAt *string `json:"last_refreshed_at"`
}
