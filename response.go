package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

// ErrorField represents error of field
type ErrorField struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ErrorResponse represents the default error response
type ErrorResponse struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Fields  []*ErrorField `json:"fields,omitempty"`
}

type errorTrapper interface {
	Code() string
}

// ResponseJSON writes json http response
func ResponseJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// ResponseError writes error http response
func ResponseError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	var errorCode string

	if et, ok := err.(errorTrapper); ok {
		errorCode = et.Code()
	} else {
		switch status {
		case http.StatusUnauthorized:
			errorCode = "Unauthorized Request"
		case http.StatusNotFound:
			errorCode = "NotFound"
		case http.StatusBadRequest:
			errorCode = "BadRequest"
		case http.StatusUnprocessableEntity:
			errorCode = "UnprocessableEntity"
		case http.StatusTooManyRequests:
			errorCode = "TooManyRequests"
		default:
			errorCode = "InternalServerError"
		}
	}

	if status == http.StatusInternalServerError {
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:    errorCode,
			Message: "Server error",
		})

		errMessage := fmt.Sprintf("%+v\n\n%s", err, debug.Stack())
		log.Println(errMessage)
	} else {
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:    errorCode,
			Message: err.Error(),
		})
	}
}
