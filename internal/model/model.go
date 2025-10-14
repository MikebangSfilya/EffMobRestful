// Package internal хранит внутри себя внутренний пакет для работы всего приложения, неэкспортируемый
// model.go содержит типы данных которые используются в приложении
package model

import (
	datatransfer "subscription/internal/api/dto"
	"time"

	"github.com/google/uuid"
)

type CustomDate struct {
	time.Time
}

// MarshalJSON реализует кастомный формат JSON "01-2006"
func (cd CustomDate) MarshalJSON() ([]byte, error) {
	return []byte(`"` + cd.Format("01-2006") + `"`), nil
}

// UnmarshalJSON парсит дату из формата "01-2006"
func (cd *CustomDate) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1]
	t, err := time.Parse("01-2006", s)
	if err != nil {
		return err
	}
	cd.Time = t
	return nil
}

// Subscription хранит запись о подписке.
// Все поля с большой буквы для экспорта в JSON.
type Subscription struct {
	ID          string      `json:"id"`
	ServiceName string      `json:"service_name"`
	Price       int         `json:"price"`
	UserId      string      `json:"user_id"`
	StartDate   CustomDate  `json:"start_date"`
	EndDate     *CustomDate `json:"end_date,omitempty"`
}

// NewSubscription создает новый объект Subscription с уникальным ID
func NewSubscription(dto datatransfer.DTOSubs) (Subscription, error) {
	startTime, err := time.Parse("01-2006", dto.StartDate)
	if err != nil {
		return Subscription{}, err
	}
	start := CustomDate{Time: startTime}

	var end *CustomDate
	if dto.EndDate != "" {
		endTime, err := time.Parse("01-2006", dto.EndDate)
		if err != nil {
			return Subscription{}, err
		}
		end = &CustomDate{Time: endTime}
	}

	return Subscription{
		ID:          uuid.New().String(),
		ServiceName: dto.ServiceName,
		Price:       dto.Price,
		UserId:      dto.UserId,
		StartDate:   start,
		EndDate:     end,
	}, nil

}
