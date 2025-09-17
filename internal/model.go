// Package internal хранит внутри себя внутренний пакет для работы всего приложения, неэкспортируемый
// model.go содержит типы данных которые используются в приложении
package internal

import (
	"time"

	"github.com/google/uuid"
)

// Subscription хранит запись о подписке.
// Все поля с большой буквы для экспорта в JSON.
type Subscription struct {
	ServiceName string
	Price       int
	UserId      string
	StartDate   string
}

// SubscriptionStore хранит слайс и мапу объектов Subscription.
// Можно использовать для хранения и работы с множеством подписок.
type SubscriptionStore struct {
	MapSub map[string]Subscription // key -- Subscription.ServiceName
}

// NewSubscription создает новый объект Subscription с уникальным UserId и текущей датой старта.
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

// NewSubMap создает пустой контейнер SubscriptionSlice для добавления подписок.
func NewSubMap() SubscriptionStore {
	return SubscriptionStore{
		MapSub: make(map[string]Subscription),
	}
}

// AddSub используется для добавления нашей подписки в хранилище(Store)
func (sub *SubscriptionStore) AddSub(subscription Subscription) {
	sub.MapSub[subscription.ServiceName] = subscription
}

func (sub *SubscriptionStore) GetSubInfo(seviceName string) Subscription {
	copyMap := make(map[string]Subscription, len(sub.MapSub))

	for k, v := range sub.MapSub {
		copyMap[k] = v
	}
	v, ok := copyMap[seviceName]
	if !ok {
		return Subscription{}
	}
	return v
}

func (sub *SubscriptionStore) GetSubAllInfo() map[string]Subscription {

	copyMap := make(map[string]Subscription, len(sub.MapSub))

	for k, v := range sub.MapSub {
		copyMap[k] = v
	}

	return copyMap

}

func (sub *SubscriptionStore) DeleteInfo(serviceName string) {
	delete(sub.MapSub, serviceName)
}
