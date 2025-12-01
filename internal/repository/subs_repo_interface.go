package repository

import (
	"context"
	"subs-service/internal/domain"
	"time"

	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, s *domain.Subscription) (int64, error)
	ReadByID(ctx context.Context, id int64) (*domain.Subscription, error)
	Update(ctx context.Context, s *domain.Subscription) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, filter SubscriptionFilter) ([]domain.Subscription, error)

	SumByPeriod(ctx context.Context, from, to time.Time, userID *uuid.UUID, serviceName *string) (int, error)
}

type SubscriptionFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	Limit       int
	Offset      int
}
