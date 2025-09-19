// DTO для запроса данных
package internal

import (
	"encoding/json"
	"net/http"
)

type DTOSubs struct {
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func writeError(w http.ResponseWriter, err string, code int) {
	w.WriteHeader(code)
	resp := ErrorResponse{
		Error: err,
		Code:  code,
	}
	json.NewEncoder(w).Encode(resp)
}
