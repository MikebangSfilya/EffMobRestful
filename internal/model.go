// Package internal хранит внутри себя внутренний пакет для работы всего приложения, неэкспортируемый
// model.go содержит типы данных которые используются в приложении
package internal

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Subscription хранит запись о подписке.
// Все поля с большой буквы для экспорта в JSON.
type Subscription struct {
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserId      string `json:"user_id"`
	StartDate   string `json:"start_date"`
}

// SubscriptionStore хранит слайс и мапу объектов Subscription.
// Можно использовать для хранения и работы с множеством подписок.
type SubscriptionStore struct {
	MapSub map[string]Subscription // key -- Subscription.UserId
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

// NewSubStore создает пустой контейнер SubscriptionSlice для добавления подписок.
func NewSubStore() SubscriptionStore {
	return SubscriptionStore{
		MapSub: make(map[string]Subscription),
	}
}

// AddSub используется для добавления нашей подписки в хранилище(Store)
func (sub *SubscriptionStore) AddSub(subscription Subscription) {
	sub.MapSub[subscription.UserId] = subscription
}

func (sub *SubscriptionStore) GetSubInfo(userId string) Subscription {
	copyMap := make(map[string]Subscription, len(sub.MapSub))

	for k, v := range sub.MapSub {
		copyMap[k] = v
	}
	v, ok := copyMap[userId]
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

func (sub *SubscriptionStore) DeleteInfo(userId string) {
	delete(sub.MapSub, userId)
}

func (sub *SubscriptionStore) UpdateSub(userId string, newSub Subscription) error {
	_, ok := sub.MapSub[userId]
	if !ok {
		return fmt.Errorf("subscription not found")
	}
	sub.MapSub[userId] = newSub
	return nil
}
