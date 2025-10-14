package datatransfer

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

var (
	errServiceName    = errors.New("service name is required")
	errPriceNegative  = errors.New("price cannot be negative")
	errUserIDRequired = errors.New("user ID is required")
	errStartDate      = errors.New("start date is required")
	errInvalidDate    = errors.New("invalid date format")
	errNoUUID         = errors.New("user ID not UUID type")
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func WriteError(w http.ResponseWriter, err string, code int) {
	log.Printf("Sending error response: %s (code: %d)", err, code)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	resp := ErrorResponse{
		Error: err,
		Code:  code,
	}
	json.NewEncoder(w).Encode(resp)
}
