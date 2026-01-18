package response

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Errors []string `json:"errors"`
}

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func ValidationError(err validator.ValidationErrors) Response {
	var errors []string

	for _, err := range err {
		switch err.Tag() {
		case "required":
			errors = append(errors, err.Field()+" is required")
		case "min":
			errors = append(errors, err.Field()+" must be at least "+err.Param())
		case "max":
			errors = append(errors, err.Field()+" must be at most "+err.Param())
		case "email":
			errors = append(errors, err.Field()+" must be a valid email address")
		}
	}
	return Response{Errors: errors}
}
