package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Subscription представляет основную сущность подписки в системе.
type Subscription struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	ServiceName string     `json:"service_name" db:"service_name"`
	Price       int        `json:"price" db:"price"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	StartDate   time.Time  `json:"start_date" db:"start_date"` // Внутреннее представление даты
	EndDate     *time.Time `json:"end_date" db:"end_date"`     // Указатель для поддержки значения NULL
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}

// SubscriptionRepository описывает методы абстрактного хранилища данных (Интерфейс слоя базы данных).
type SubscriptionRepository interface {
	Create(ctx context.Context, sub *Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*Subscription, error)
	Update(ctx context.Context, sub *Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]Subscription, error)
	GetTotalCost(ctx context.Context, userID uuid.UUID, serviceName string, from, to time.Time) (int, error)
}

// SubscriptionService описывает методы слоя бизнес-логики (Интерфейс для Handler).
type SubscriptionService interface {
	CreateSubscription(ctx context.Context, req *CreateSubscriptionReq) (*Subscription, error)
	GetSubscription(ctx context.Context, id uuid.UUID) (*Subscription, error)
	UpdateSubscription(ctx context.Context, id uuid.UUID, req *UpdateSubscriptionReq) (*Subscription, error)
	DeleteSubscription(ctx context.Context, id uuid.UUID) error
	ListSubscriptions(ctx context.Context, limit, offset int) ([]Subscription, error)
	CalculateTotalCost(ctx context.Context, userID uuid.UUID, serviceName string, fromStr, toStr string) (int, error)
}

// CreateSubscriptionReq DTO структура для валидации входящего запроса на создание.
type CreateSubscriptionReq struct {
	ServiceName string    `json:"service_name" binding:"required,min=1,max=100"`
	Price       int       `json:"price" binding:"required,gt=0"` // Валидация: строго больше 0
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	StartDate   string    `json:"start_date" binding:"required,customdate"` // "MM-YYYY" кастомный тег валидатора
	EndDate     string    `json:"end_date" binding:"omitempty,customdate"`  // Опционально
}

// UpdateSubscriptionReq DTO структура для валидации обновления.
type UpdateSubscriptionReq struct {
	ServiceName string    `json:"service_name" binding:"required,min=1,max=100"`
	Price       int       `json:"price" binding:"required,gt=0"`
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	StartDate   string    `json:"start_date" binding:"required,customdate"`
	EndDate     string    `json:"end_date" binding:"omitempty,customdate"`
}
