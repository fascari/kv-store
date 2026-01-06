package http

import "net/http"

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewErrorResponse(message string) ErrorResponse {
	return ErrorResponse{Error: message}
}

// Complete error response helpers

func BadRequest(w http.ResponseWriter, message string) {
	JSON(w, http.StatusBadRequest, NewErrorResponse(message))
}

func NotFound(w http.ResponseWriter, message string) {
	JSON(w, http.StatusNotFound, NewErrorResponse(message))
}

func InternalServerError(w http.ResponseWriter, message string) {
	JSON(w, http.StatusInternalServerError, NewErrorResponse(message))
}
