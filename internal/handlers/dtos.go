package handlers

type RefreshResponse struct {
	Message string `json:"message"`
	Total   int    `json:"total"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
