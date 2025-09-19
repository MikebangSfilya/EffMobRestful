// Package internal хранит внутри себя внутренний пакет для работы всего приложения, неэкспортируемый
// model.go содержит типы данных которые используются в приложении
package internal

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
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
	dataBase *pgxpool.Pool // база данных
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
func NewSubStore(dataBase *pgxpool.Pool) *SubscriptionStore {
	return &SubscriptionStore{
		dataBase: dataBase,
	}
}

// AddSub используется для добавления нашей подписки в хранилище(Store)
func (sub *SubscriptionStore) AddSub(ctx context.Context, subscription Subscription) error {
	query := `
		INSERT 
		INTO subscription 
		(user_id, service_name, price, start_date)
		VALUES ($1, $2, $3, $4)
	`

	_, err := sub.dataBase.Exec(ctx, query, subscription.UserId, subscription.ServiceName, subscription.Price, subscription.StartDate)
	return err
}

func (sub *SubscriptionStore) GetSubInfo(ctx context.Context, userId string) (Subscription, error) {

	query := `
	SELECT user_id, service_name, price, start_date 
	FROM subscription
	WHERE user_id=$1
	`
	var s Subscription

	if err := sub.dataBase.QueryRow(ctx, query, userId).Scan(&s.UserId, &s.ServiceName, &s.Price, &s.StartDate); err != nil {
		return Subscription{}, err
	}

	return s, nil
}

func (sub *SubscriptionStore) GetSubAllInfo(ctx context.Context) ([]Subscription, error) {

	query := `
	SELECT user_id, service_name, price, start_date 
	FROM subscription 
	`
	rows, err := sub.dataBase.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []Subscription
	for rows.Next() {
		var sub Subscription
		if err := rows.Scan(&sub.UserId, &sub.ServiceName, &sub.Price, &sub.StartDate); err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}
	return subs, nil
}

// Удаление записи из нашей базы данных
func (sub *SubscriptionStore) DeleteInfo(ctx context.Context, userId string) error {
	query := `
	DELETE F
	ROM subscription 
	WHERE user_id=$1`
	_, err := sub.dataBase.Exec(ctx, query, userId)
	return err
}

// Обновление информации из базы данных
func (sub *SubscriptionStore) UpdateSub(ctx context.Context, userId string, newSub Subscription) error {
	query := `
		UPDATE subscription 
		SET service_name=$1, price=$2, start_date=$3 
		WHERE user_id=$4
	`
	cmd, err := sub.dataBase.Exec(ctx, query, newSub.ServiceName, newSub.Price, newSub.StartDate, userId)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("subscription not found")
	}
	return nil
}
