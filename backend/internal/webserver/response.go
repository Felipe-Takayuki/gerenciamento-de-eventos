package webserver

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Felipe-Takayuki/Adamas/adamas-api/internal/utils"
)

type MessageResponse struct {
	Message string `json:"message"`
}

func jsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func errorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(utils.ErrorMessage{Message: message})
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, utils.ErrNotFound):
		errorResponse(w, http.StatusNotFound, err.Error())
	case errors.Is(err, utils.ErrInvalidCredentials):
		errorResponse(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, utils.ErrUnauthorized):
		errorResponse(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, utils.ErrForbidden):
		errorResponse(w, http.StatusForbidden, err.Error())
	case errors.Is(err, utils.ErrBadRequest):
		errorResponse(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, utils.ErrInstitutionNotOwner):
		errorResponse(w, http.StatusForbidden, err.Error())
	default:
		errorResponse(w, http.StatusBadRequest, err.Error())
	}
}
