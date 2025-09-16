// Package internal хранит внутри себя внутренний пакет для работы всего приложения, неэкспортируемый
// model.go содержит типы данных которые используются в приложении
package internal

import (
	"time"

	"github.com/google/uuid"
)

// Subcriptoin godoc
// @Summary subscription type
// @Descriptoin Здесь хранится все наши записи о подписке, имя, стоимость, юзер дата и стартовое время
// Все с большой буквы для экспорта в json
type Subscription struct {
	ServiceName string
	Price       int
	UserId      string
	StartDate   string
}

func NewSubscription(serviceName string, price int) Subscription {
	userId := uuid.New()
	startedDate := time.Now()

	return Subscription{
		ServiceName: serviceName,
		Price:       price,
		UserId:      userId.String(),
		StartDate:   startedDate.Format("01-2006"),
	}

}
