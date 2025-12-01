package service

import (
	"context"
	"fmt"
	"time"

	"subs-service/internal/domain"
	"subs-service/internal/repository"
	"subs-service/logger"

	"github.com/google/uuid"
)

type SubscriptionService struct {
	repo   repository.SubscriptionRepository
	logger *logger.Logger
}

func (s *SubscriptionService) Logger() *logger.Logger {
	return s.logger
}

func NewSubscriptionService(repo repository.SubscriptionRepository, log *logger.Logger) *SubscriptionService {
	return &SubscriptionService{
		repo:   repo,
		logger: log,
	}
}

func parseMonthYear(s string) (time.Time, error) {
	t, err := time.Parse("01-2006", s)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format: %s, expected MM-YYYY", s)
	}

	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC), nil
}

func (s *SubscriptionService) Create(
	ctx context.Context,
	serviceName string,
	price int,
	userID string,
	start string,
	end *string,
) (int64, error) {

	if serviceName == "" {
		s.logger.Error("service name empty")
		return 0, fmt.Errorf("service name cannot be empty")
	}

	if price <= 0 {
		s.logger.Error("price must be > 0")
		return 0, fmt.Errorf("price must be > 0")
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		s.logger.Error("invalid user_id", "error", err)
		return 0, fmt.Errorf("invalid user_id: %v", err)
	}

	startDate, err := parseMonthYear(start)
	if err != nil {
		s.logger.Error("invalid start_date", "error", err)
		return 0, err
	}

	var endDate *time.Time
	if end != nil && *end != "" {
		parsed, err := parseMonthYear(*end)
		if err != nil {
			s.logger.Error("invalid end_date", "error", err)
			return 0, err
		}
		endDate = &parsed
	}

	sub := domain.Subscription{
		ServiceName: serviceName,
		Price:       price,
		UserID:      uid,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	id, err := s.repo.Create(ctx, &sub)
	if err != nil {
		s.logger.Error("failed to create subscription", "error", err)
		return 0, err
	}

	s.logger.Info("subscription created", "id", id)
	return id, nil
}

func (s *SubscriptionService) ReadByID(ctx context.Context, id int64) (*domain.Subscription, error) {
	result, err := s.repo.ReadByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to read subscription", "id", id, "error", err)
		return nil, err
	}
	if result == nil {
		s.logger.Error("subscription not found", "id", id)
		return nil, fmt.Errorf("subscription not found")
	}
	return result, nil
}

func (s *SubscriptionService) Update(
	ctx context.Context,
	id int64,
	serviceName string,
	price int,
	userID string,
	start string,
	end *string,
) error {

	if serviceName == "" {
		s.logger.Error("service name empty")
		return fmt.Errorf("service name cannot be empty")
	}

	if price <= 0 {
		s.logger.Error("price must be > 0")
		return fmt.Errorf("price must be > 0")
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		s.logger.Error("invalid user_id", "error", err)
		return fmt.Errorf("invalid user_id: %v", err)
	}

	startDate, err := parseMonthYear(start)
	if err != nil {
		s.logger.Error("invalid start_date", "error", err)
		return err
	}

	var endDate *time.Time
	if end != nil && *end != "" {
		parsed, err := parseMonthYear(*end)
		if err != nil {
			s.logger.Error("invalid end_date", "error", err)
			return err
		}
		endDate = &parsed
	}

	sub := domain.Subscription{
		ID:          id,
		ServiceName: serviceName,
		Price:       price,
		UserID:      uid,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	err = s.repo.Update(ctx, &sub)
	if err != nil {
		s.logger.Error("failed to update subscription", "id", id, "error", err)
		return err
	}

	s.logger.Info("subscription updated", "id", id)
	return nil
}

func (s *SubscriptionService) Delete(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("failed to delete subscription", "id", id, "error", err)
		return err
	}
	s.logger.Info("subscription deleted", "id", id)
	return nil
}

func (s *SubscriptionService) List(
	ctx context.Context,
	userID *string,
	serviceName *string,
	limit int,
	offset int,
) ([]domain.Subscription, error) {

	var uid *uuid.UUID
	if userID != nil {
		parsed, err := uuid.Parse(*userID)
		if err != nil {
			s.logger.Error("invalid user_id", "error", err)
			return nil, fmt.Errorf("invalid user_id: %v", err)
		}
		uid = &parsed
	}

	f := repository.SubscriptionFilter{
		UserID:      uid,
		ServiceName: serviceName,
		Limit:       limit,
		Offset:      offset,
	}

	list, err := s.repo.List(ctx, f)
	if err != nil {
		s.logger.Error("failed to list subscriptions", "error", err)
		return nil, err
	}

	return list, nil
}

func (s *SubscriptionService) SumByPeriod(
	ctx context.Context,
	fromStr, toStr string,
	userID *string,
	serviceName *string,
) (int, error) {

	from, err := parseMonthYear(fromStr)
	if err != nil {
		s.logger.Error("invalid from date", "error", err)
		return 0, err
	}

	to, err := parseMonthYear(toStr)
	if err != nil {
		s.logger.Error("invalid to date", "error", err)
		return 0, err
	}

	var uid *uuid.UUID
	if userID != nil {
		parsed, err := uuid.Parse(*userID)
		if err != nil {
			s.logger.Error("invalid user_id", "error", err)
			return 0, fmt.Errorf("invalid user_id: %v", err)
		}
		uid = &parsed
	}

	sum, err := s.repo.SumByPeriod(ctx, from, to, uid, serviceName)
	if err != nil {
		s.logger.Error("failed to sum subscriptions", "error", err)
		return 0, err
	}

	return sum, nil
}
