package handler

// ErrorResponse is the standard JSON error envelope returned by all handlers.
type ErrorResponse struct {
	Error string `json:"error"`
}
