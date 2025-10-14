package repository

import (
	"context"
	"database/sql"
	"errors"
	"subscription/internal/model"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxRepository struct {
	db *pgxpool.Pool
}

func nullableDate(cd *model.CustomDate) interface{} {
	if cd == nil {
		return nil
	}
	return cd.Time
}

func NewPgxRepository(db *pgxpool.Pool) *pgxRepository {
	return &pgxRepository{
		db: db,
	}
}

// AddSub используется для добавления нашей подписки в хранилище(Store)
func (sub *pgxRepository) Create(ctx context.Context, subscription model.Subscription) error {
	query := `
		INSERT 
		INTO subscription 
		(id, user_id, service_name, price, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := sub.db.Exec(
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

func (sub *pgxRepository) GetByID(ctx context.Context, Id string) (model.Subscription, error) {

	query := `
	SELECT id, user_id, service_name, price, start_date, end_date
	FROM subscription
	WHERE id=$1
	`
	var s model.Subscription
	var startDate time.Time
	var endDate sql.NullTime

	if err := sub.db.QueryRow(ctx, query, Id).Scan(&s.ID, &s.UserId, &s.ServiceName, &s.Price, &startDate, &endDate); err != nil {
		return model.Subscription{}, err
	}
	s.StartDate = model.CustomDate{Time: startDate}
	if endDate.Valid {
		s.EndDate = &model.CustomDate{Time: endDate.Time}
	}

	return s, nil
}

func (sub *pgxRepository) GetAll(ctx context.Context) ([]model.Subscription, error) {

	query := `
	SELECT id, user_id, service_name, price, start_date, end_date
	FROM subscription 
	`
	rows, err := sub.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []model.Subscription
	var startDate time.Time
	var endDate sql.NullTime

	for rows.Next() {
		var sub model.Subscription
		if err := rows.Scan(&sub.ID, &sub.UserId, &sub.ServiceName, &sub.Price, &startDate, &endDate); err != nil {
			return nil, err
		}
		sub.StartDate = model.CustomDate{Time: startDate}
		if endDate.Valid {
			sub.EndDate = &model.CustomDate{Time: endDate.Time}
		}

		subs = append(subs, sub)
	}
	return subs, nil
}

// Удаление записи из нашей базы данных
func (sub *pgxRepository) Delete(ctx context.Context, id string) error {

	query := `
	DELETE FROM subscription 
	WHERE id=$1`
	_, err := sub.db.Exec(ctx, query, id)
	return err
}

// Обновление информации из базы данных
func (sub *pgxRepository) Update(ctx context.Context, id string, newSub model.Subscription) error {

	query := `
		UPDATE subscription 
		SET service_name=$1, price=$2, start_date=$3
		WHERE id=$4
	`
	cmd, err := sub.db.Exec(
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
func (sub *pgxRepository) SumForPeriod(ctx context.Context, userId, serviceName string, from, to time.Time) (int, error) {

	query := `
		SELECT SUM(price) 
		FROM subscription
		WHERE ($1='' OR user_id=$1)
		  AND ($2='' OR service_name=$2)
		  AND start_date >= $3 AND start_date <= $4
	`
	var sum int
	err := sub.db.QueryRow(
		ctx,
		query,
		userId,
		serviceName,
		from, to).Scan(&sum)
	return sum, err
}
