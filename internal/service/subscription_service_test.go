package service

import (
	"context"
	"testing"
	"time"

	"github.com/april1858/test/internal/domain"
	"github.com/google/uuid"
)

// MockSubscriptionRepository реализует domain.SubscriptionRepository для тестов.
type MockSubscriptionRepository struct {
	CreateFunc       func(ctx context.Context, sub *domain.Subscription) error
	GetByIDFunc      func(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	UpdateFunc       func(ctx context.Context, sub *domain.Subscription) error
	DeleteFunc       func(ctx context.Context, id uuid.UUID) error
	ListFunc         func(ctx context.Context, limit, offset int) ([]domain.Subscription, error)
	GetTotalCostFunc func(ctx context.Context, userID uuid.UUID, serviceName string, from, to time.Time) (int, error)
}

func (m *MockSubscriptionRepository) Create(ctx context.Context, sub *domain.Subscription) error {
	return m.CreateFunc(ctx, sub)
}
func (m *MockSubscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	return m.GetByIDFunc(ctx, id)
}
func (m *MockSubscriptionRepository) Update(ctx context.Context, sub *domain.Subscription) error {
	return m.UpdateFunc(ctx, sub)
}
func (m *MockSubscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.DeleteFunc(ctx, id)
}
func (m *MockSubscriptionRepository) List(ctx context.Context, limit, offset int) ([]domain.Subscription, error) {
	return m.ListFunc(ctx, limit, offset)
}
func (m *MockSubscriptionRepository) GetTotalCost(ctx context.Context, userID uuid.UUID, serviceName string, from, to time.Time) (int, error) {
	return m.GetTotalCostFunc(ctx, userID, serviceName, from, to)
}

// TestCreateSubscription_Success проверяет успешный сценарий создания подписки.
func TestCreateSubscription_Success(t *testing.T) {
	mockRepo := &MockSubscriptionRepository{
		CreateFunc: func(ctx context.Context, sub *domain.Subscription) error {
			sub.ID = uuid.New()
			return nil
		},
	}
	svc := NewSubscriptionService(mockRepo)
	req := &domain.CreateSubscriptionReq{
		ServiceName: "Netflix",
		Price:       800,
		UserID:      uuid.New(),
		StartDate:   "01-2026",
		EndDate:     "12-2026",
	}

	res, err := svc.CreateSubscription(context.Background(), req)
	if err != nil {
		t.Fatalf("Ожидалось отсутствие ошибки, получено: %v", err)
	}
	if res.ServiceName != req.ServiceName {
		t.Errorf("Ожидалось имя сервиса %s, получено %s", req.ServiceName, res.ServiceName)
	}
	if res.Price != req.Price {
		t.Errorf("Ожидалась цена %d, получена %d", req.Price, res.Price)
	}
}

// TestCreateSubscription_ValidationError проверяет валидацию некорректных дат.
func TestCreateSubscription_ValidationError(t *testing.T) {
	mockRepo := &MockSubscriptionRepository{}
	svc := NewSubscriptionService(mockRepo)
	req := &domain.CreateSubscriptionReq{
		ServiceName: "Spotify",
		Price:       300,
		UserID:      uuid.New(),
		StartDate:   "12-2026",
		EndDate:     "01-2026", // Дата окончания раньше даты начала
	}
	_, err := svc.CreateSubscription(context.Background(), req)
	if err == nil {
		t.Fatal("Ожидалась ошибка валидации бизнес-логики дат, но её нет")
	}
	expectedErr := "end_date cannot be before start_date"
	if err.Error() != expectedErr {
		t.Errorf("Ожидалась ошибка '%s', получена '%s'", expectedErr,
			err.Error())
	}
}
