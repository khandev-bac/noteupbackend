package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJsonError(w http.ResponseWriter, message string, statusCode int, err error) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(&ErrorMessage{
		Message:    message,
		StatusCode: statusCode,
		Error:      err,
	})
}

func WriteJsonESuccess(w http.ResponseWriter, message string, statusCode int, data any) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(&SuccessMessage{
		Message:    message,
		StatusCode: statusCode,
		Data:       data,
	})
}
