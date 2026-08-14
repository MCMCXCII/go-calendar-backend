package httpserver

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"project/internal/events/domain"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("write json response: %v", err)
	}
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrEventNotFound):
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "event not found"})

	case errors.Is(err, domain.ErrInvalidTimeRange),
		errors.Is(err, domain.ErrInvalidType),
		errors.Is(err, domain.ErrCustomTypeRequired),
		errors.Is(err, domain.ErrCustomTypeNotAllowed),
		errors.Is(err, domain.ErrTitleRequired),
		errors.Is(err, domain.ErrInvalidPeriod),
		errors.Is(err, domain.ErrPeriodRequired):
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})

	default:
		log.Printf("internal error: %v", err)
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
	}
}
