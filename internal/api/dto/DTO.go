// datatransfer для запроса данных
package datatransfer

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

type SumResponse struct {
	TotalPrice int `json:"total_price"`
}

func (d DTOSubs) Validate() error {
	if d.ServiceName == "" {
		return errServiceName
	}
	if d.Price < 0 {
		return errPriceNegative
	}
	if d.UserId == "" {
		return errUserIDRequired
	}
	if d.StartDate == "" {
		return errStartDate
	}
	return nil

}
