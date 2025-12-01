package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"subs-service/internal/domain"

	"github.com/google/uuid"
)

type PostgresSubscriptionRepository struct {
	db *sql.DB
}

func NewPostgresSubscriptionRepository(db *sql.DB) *PostgresSubscriptionRepository {
	return &PostgresSubscriptionRepository{db: db}
}

func (r *PostgresSubscriptionRepository) Create(ctx context.Context, s *domain.Subscription) (int64, error) {
	query := `
        INSERT INTO subscriptions
        (service_name, price, user_id, start_date, end_date)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id;
    `
	var id int64

	err := r.db.QueryRowContext(
		ctx,
		query,
		s.ServiceName,
		s.Price,
		s.UserID,
		s.StartDate,
		s.EndDate,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PostgresSubscriptionRepository) ReadByID(ctx context.Context, id int64) (*domain.Subscription, error) {
	query := `
        SELECT id, service_name, price, user_id, start_date, end_date
        FROM subscriptions
        WHERE id = $1;
    `
	var s domain.Subscription

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID,
		&s.ServiceName,
		&s.Price,
		&s.UserID,
		&s.StartDate,
		&s.EndDate,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *PostgresSubscriptionRepository) Update(ctx context.Context, s *domain.Subscription) error {
	query := `
        UPDATE subscriptions
        SET service_name = $1,
            price = $2,
            user_id = $3,
            start_date = $4,
            end_date = $5,
            updated_at = NOW()
        WHERE id = $6;
    `

	_, err := r.db.ExecContext(
		ctx,
		query,
		s.ServiceName,
		s.Price,
		s.UserID,
		s.StartDate,
		s.EndDate,
		s.ID,
	)

	return err
}

func (r *PostgresSubscriptionRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM subscriptions WHERE id = $1;`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresSubscriptionRepository) List(ctx context.Context, f SubscriptionFilter) ([]domain.Subscription, error) {
	var args []interface{}
	var conditions []string
	baseQuery := `
        SELECT id, service_name, price, user_id, start_date, end_date
        FROM subscriptions
    `

	i := 1

	if f.UserID != nil {
		conditions = append(conditions, "user_id = $"+fmt.Sprint(i))
		args = append(args, *f.UserID)
		i++
	}

	if f.ServiceName != nil {
		conditions = append(conditions, "service_name = $"+fmt.Sprint(i))
		args = append(args, *f.ServiceName)
		i++
	}

	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	baseQuery += " ORDER BY id ASC"

	if f.Limit == 0 {
		f.Limit = 100
	}
	if f.Limit > 0 {
		baseQuery += " LIMIT " + fmt.Sprint(f.Limit)
	}
	if f.Offset > 0 {
		baseQuery += " OFFSET " + fmt.Sprint(f.Offset)
	}

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []domain.Subscription{}

	for rows.Next() {
		var s domain.Subscription
		err := rows.Scan(
			&s.ID,
			&s.ServiceName,
			&s.Price,
			&s.UserID,
			&s.StartDate,
			&s.EndDate,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}

	return result, rows.Err()
}

func (r *PostgresSubscriptionRepository) SumByPeriod(
	ctx context.Context,
	from, to time.Time,
	userID *uuid.UUID,
	serviceName *string,
) (int, error) {

	query := `
        SELECT COALESCE(SUM(price), 0)
        FROM subscriptions
        WHERE start_date <= $1
          AND (end_date IS NULL OR end_date >= $2)
    `

	args := []any{to, from}
	i := 3

	if userID != nil {
		query += " AND user_id = $" + fmt.Sprint(i)
		args = append(args, *userID)
		i++
	}

	if serviceName != nil {
		query += " AND service_name = $" + fmt.Sprint(i)
		args = append(args, *serviceName)
		i++
	}

	var sum int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&sum)

	return sum, err
}
