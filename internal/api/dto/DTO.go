// datatransfer для запроса данных
package datatransfer

import (
	"encoding/json"
	"net/http"
)

type DTOSubs struct {
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserId      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date,omitempty"`
}

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

type SumResponse struct {
	TotalPrice int `json:"total_price"`
}
