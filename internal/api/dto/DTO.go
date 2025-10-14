// datatransfer для запроса данных
package datatransfer

import "github.com/google/uuid"

type DTOSubs struct {
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserId      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date,omitempty"`
}

type SumResponse struct {
	TotalPrice int `json:"total_price"`
}

func (d DTOSubs) Validate() error {

	_, err := uuid.Parse(d.UserId)
	if err != nil {
		return errNoUUID
	}
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
