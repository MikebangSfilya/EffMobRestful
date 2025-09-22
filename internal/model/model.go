// Package internal хранит внутри себя внутренний пакет для работы всего приложения, неэкспортируемый
// model.go содержит типы данных которые используются в приложении
package model

import (
	"context"
	"database/sql"
	"errors"
	datatransfer "subscription/internal/dto"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
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

func nullableDate(cd *CustomDate) interface{} {
	if cd == nil {
		return nil
	}
	return cd.Time
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

// SubscriptionStore хранит подключение к базе данных
type SubscriptionStore struct {
	dataBase *pgxpool.Pool // база данных
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
		(id, user_id, service_name, price, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := sub.dataBase.Exec(
		ctx,
		query,
		subscription.ID,
		subscription.UserId,
		subscription.ServiceName,
		subscription.Price,
		subscription.StartDate.Time,
		nullableDate(subscription.EndDate))
	return err
}

func (sub *SubscriptionStore) GetSubInfo(ctx context.Context, Id string) (Subscription, error) {

	query := `
	SELECT id, user_id, service_name, price, start_date, end_date
	FROM subscription
	WHERE id=$1
	`
	var s Subscription
	var startDate time.Time
	var endDate sql.NullTime

	if err := sub.dataBase.QueryRow(ctx, query, Id).Scan(&s.ID, &s.UserId, &s.ServiceName, &s.Price, &startDate, &endDate); err != nil {
		return Subscription{}, err
	}
	s.StartDate = CustomDate{Time: startDate}
	if endDate.Valid {
		s.EndDate = &CustomDate{Time: endDate.Time}
	}

	return s, nil
}

func (sub *SubscriptionStore) GetSubAllInfo(ctx context.Context) ([]Subscription, error) {

	query := `
	SELECT id, user_id, service_name, price, start_date, end_date
	FROM subscription 
	`
	rows, err := sub.dataBase.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []Subscription
	var startDate time.Time
	var endDate sql.NullTime

	for rows.Next() {
		var sub Subscription
		if err := rows.Scan(&sub.ID, &sub.UserId, &sub.ServiceName, &sub.Price, &startDate, &endDate); err != nil {
			return nil, err
		}
		sub.StartDate = CustomDate{Time: startDate}
		if endDate.Valid {
			sub.EndDate = &CustomDate{Time: endDate.Time}
		}

		subs = append(subs, sub)
	}
	return subs, nil
}

// Удаление записи из нашей базы данных
func (sub *SubscriptionStore) DeleteInfo(ctx context.Context, id string) error {
	query := `
	DELETE FROM subscription 
	WHERE id=$1`
	_, err := sub.dataBase.Exec(ctx, query, id)
	return err
}

// Обновление информации из базы данных
func (sub *SubscriptionStore) UpdateSub(ctx context.Context, id string, newSub Subscription) error {
	query := `
		UPDATE subscription 
		SET service_name=$1, price=$2, start_date=$3
		WHERE id=$5
	`
	cmd, err := sub.dataBase.Exec(
		ctx,
		query,
		newSub.ServiceName,
		newSub.Price,
		newSub.StartDate.Time,
		id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("subscription not found")
	}
	return nil
}

// Метод для подсчета суммарной стоимости всех подписок за выбранный период
func (sub *SubscriptionStore) SumSubscriptions(ctx context.Context, userId, serviceName string, from, to CustomDate) (int, error) {
	query := `
		SELECT SUM(price) 
		FROM subscription
		WHERE ($1='' OR user_id=$1)
		  AND ($2='' OR service_name=$2)
		  AND start_date >= $3 AND start_date <= $4
	`
	var sum int
	err := sub.dataBase.QueryRow(
		ctx,
		query,
		userId,
		serviceName,
		from.Time, to.Time).Scan(&sum)
	return sum, err
}
