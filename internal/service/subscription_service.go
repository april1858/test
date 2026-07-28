package service

import (
	"context"
	"errors"
	"time"

	"github.com/april1858/test/internal/domain"
	"github.com/google/uuid"
	// Пакет логирования будет настроен далее, пока используем стандартный контекст или заглушку логера
)

// SubscriptionService реализует интерфейс domain.SubscriptionService.
type SubscriptionService struct {
	repo domain.SubscriptionRepository
}

// NewSubscriptionService — конструктор для внедрения зависимости репозитория (Dependency Injection).
func NewSubscriptionService(repo domain.SubscriptionRepository) domain.SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, req *domain.CreateSubscriptionReq) (*domain.Subscription, error) {
	// Конвертируем строки дат в типы time.Time для БД
	startDate, err := parseMMYYYY(req.StartDate)
	if err != nil {
		return nil, errors.New("invalid start_date format, use MM-YYYY")
	}

	var endDatePtr *time.Time
	if req.EndDate != "" {
		endDate, err := parseMMYYYY(req.EndDate)
		if err != nil {
			return nil, errors.New("invalid end_date format, use MM-YYYY")
		}
		// Бизнес-валидация: дата окончания не может быть раньше даты начала
		if endDate.Before(startDate) {
			return nil, errors.New("end_date cannot be before start_date")
		}
		endDatePtr = &endDate
	}
	// Формируем чистую доменную модель
	sub := &domain.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   startDate,
		EndDate:     endDatePtr,
	}
	// Передаем управление слою работы с данными
	if err := s.repo.Create(ctx, sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *SubscriptionService) GetSubscription(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *SubscriptionService) UpdateSubscription(ctx context.Context, id uuid.UUID, req *domain.UpdateSubscriptionReq) (*domain.Subscription, error) {
	startDate, err := parseMMYYYY(req.StartDate)
	if err != nil {
		return nil, errors.New("invalid start_date format, use MM-YYYY")
	}
	var endDatePtr *time.Time
	if req.EndDate != "" {
		endDate, err := parseMMYYYY(req.EndDate)
		if err != nil {
			return nil, errors.New("invalid end_date format, use MM-YYYY")
		}
		if endDate.Before(startDate) {
			return nil, errors.New("end_date cannot be before start_date")
		}
		endDatePtr = &endDate
	}
	// Создаем модель для обновления
	sub := &domain.Subscription{
		ID:          id,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   startDate,
		EndDate:     endDatePtr,
	}
	if err := s.repo.Update(ctx, sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *SubscriptionService) DeleteSubscription(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *SubscriptionService) ListSubscriptions(ctx context.Context, limit, offset int) ([]domain.Subscription, error) {
	// Защита от некорректных значений пагинации
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.List(ctx, limit, offset)
}

// CalculateTotalCost подготавливает даты и агрегирует стоимость из репозитория.
func (s *SubscriptionService) CalculateTotalCost(ctx context.Context, userID uuid.UUID, serviceName string, fromStr, toStr string) (int, error) {
	from, err := parseMMYYYY(fromStr)
	if err != nil {
		return 0, errors.New("invalid 'from' date format, use MM-YYYY")
	}
	to, err := parseMMYYYY(toStr)
	if err != nil {
		return 0, errors.New("invalid 'to' date format, use MM-YYYY")
	}
	if to.Before(from) {
		return 0, errors.New("'to' date cannot be before 'from' date")
	}
	return s.repo.GetTotalCost(ctx, userID, serviceName, from, to)
}
