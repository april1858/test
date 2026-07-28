package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/april1858/test/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// PostgresRepository реализует интерфейс domain.SubscriptionRepository.
type PostgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository конструктор для создания репозитория.
func NewPostgresRepository(db *sqlx.DB) domain.SubscriptionRepository {
	return &PostgresRepository{db: db}
}
func (r *PostgresRepository) Create(ctx context.Context, sub *domain.Subscription) error {
	query := `INSERT INTO subscriptions (service_name, price, user_id,
	start_date, end_date)
	VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	// Выполняем запрос внутри контекста для поддержки таймаутов и логирования
	err := r.db.QueryRowContext(ctx, query, sub.ServiceName, sub.Price,
		sub.UserID, sub.StartDate, sub.EndDate).
		Scan(&sub.ID, &sub.CreatedAt)
	return err
}
func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date, created_at FROM subscriptions WHERE id = $1`
	var sub domain.Subscription
	err := r.db.GetContext(ctx, &sub, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("subscription not found")
		}
		return nil, err
	}
	return &sub, nil
}

func (r *PostgresRepository) Update(ctx context.Context, sub *domain.Subscription) error {
	query := `UPDATE subscriptions SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5 WHERE id = $6`
	res, err := r.db.ExecContext(ctx, query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate, sub.ID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("subscription to update not found")
	}
	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM subscriptions WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("subscription to delete not found")
	}
	return nil
}

func (r *PostgresRepository) List(ctx context.Context, limit, offset int) ([]domain.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date, created_at FROM subscriptions ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	var subs []domain.Subscription

	err := r.db.SelectContext(ctx, &subs, query, limit, offset)
	if err != nil {
		return nil, err
	}
	return subs, nil
}

// GetTotalCost считает суммарную стоимость подписок пользователя за период по названию сервиса.

func (r *PostgresRepository) GetTotalCost(ctx context.Context, userID uuid.UUID, serviceName string, from, to time.Time) (int, error) {
	// COALESCE гарантирует, что если записей нет, вернется 0, а не NULL.
	// Логика пересечения интервалов: подписка активна в периоде, если ее старт <= концу периода И (энд >= старту периода ИЛИ энд NULL)
	query := `SELECT COALESCE(SUM(price), 0) FROM subscriptions WHERE user_id = $1 AND service_name = $2 AND start_date <= $3 AND (end_date IS NULL OR end_date >= $4)`
	var total int
	err := r.db.GetContext(ctx, &total, query, userID, serviceName, to, from)
	if err != nil {
		return 0, err
	}
	return total, nil
}
