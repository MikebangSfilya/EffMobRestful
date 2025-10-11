package datatransfer

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	errServiceName    = errors.New("service name is required")
	errPriceNegative  = errors.New("price cannot be negative")
	errUserIDRequired = errors.New("user ID is required")
	errStartDate      = errors.New("start date is required")
	errInvalidDate    = errors.New("invalid date format")
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func WriteError(w http.ResponseWriter, err string, code int) {
	w.WriteHeader(code)
	resp := ErrorResponse{
		Error: err,
		Code:  code,
	}
	json.NewEncoder(w).Encode(resp)
}
