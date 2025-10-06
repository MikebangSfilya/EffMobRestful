package datatransfer

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrValidation     = errors.New("validation error")
	ErrServiceName    = errors.New("service name is required")
	ErrPriceNegative  = errors.New("price cannot be negative")
	ErrUserIDRequired = errors.New("user ID is required")
	ErrStartDate      = errors.New("start date is required")
	ErrInvalidDate    = errors.New("invalid date format")
)

func WriteError(w http.ResponseWriter, err string, code int) {
	w.WriteHeader(code)
	resp := ErrorResponse{
		Error: err,
		Code:  code,
	}
	json.NewEncoder(w).Encode(resp)
}
