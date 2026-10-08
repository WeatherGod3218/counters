package models

type ErrorResponse struct {
	Error string `json:"error"`
}

type ResetResponse struct {
	Exists bool `json:"exists"`
}

func NewErrorResponse() *ErrorResponse {
	return &ErrorResponse{
		Error: "There was an error with this request",
	}
}
