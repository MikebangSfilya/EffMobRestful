package service

import (
	"context"
	datatransfer "subscription/internal/api/dto"
	"subscription/internal/model"
	"time"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub model.Subscription) error
	GetByID(ctx context.Context, id string) (model.Subscription, error)
	GetAll(ctx context.Context) ([]model.Subscription, error)
	Update(ctx context.Context, id string, sub model.Subscription) error
	Delete(ctx context.Context, id string) error
	SumForPeriod(ctx context.Context, userId, serviceName string, from, to time.Time) (int, error)
}

type ServiceStore struct {
	subscriptionStore SubscriptionRepository
}

func NewService(subStore SubscriptionRepository) *ServiceStore {
	return &ServiceStore{
		subscriptionStore: subStore,
	}
}

func (s *ServiceStore) Create(ctx context.Context, dto datatransfer.DTOSubs) (model.Subscription, error) {

	sub, err := model.NewSubscription(dto)
	if err != nil {
		return model.Subscription{}, err
	}
	if err := s.subscriptionStore.Create(ctx, sub); err != nil {
		return model.Subscription{}, err
	}

	return sub, nil

}

func (s *ServiceStore) GetInfo(ctx context.Context, idSub string) (model.Subscription, error) {

	sub, err := s.subscriptionStore.GetByID(ctx, idSub)
	if err != nil {
		return model.Subscription{}, err
	}
	return sub, nil
}

func (s *ServiceStore) GetAll(ctx context.Context) ([]model.Subscription, error) {

	allSub, err := s.subscriptionStore.GetAll(ctx)
	if err != nil {
		return []model.Subscription{}, err
	}
	return allSub, nil
}

func (s *ServiceStore) Delete(ctx context.Context, idSub string) error {
	return s.subscriptionStore.Delete(ctx, idSub)
}

func (s *ServiceStore) Update(ctx context.Context, id string, dto datatransfer.DTOSubs) (model.Subscription, error) {

	oldSub, err := s.subscriptionStore.GetByID(ctx, id)
	if err != nil {
		return model.Subscription{}, err
	}
	updatedSub := model.Subscription{
		ServiceName: dto.ServiceName,
		Price:       dto.Price,
		UserId:      oldSub.UserId,
		StartDate:   oldSub.StartDate,
	}

	if err := s.subscriptionStore.Update(ctx, id, updatedSub); err != nil {
		return model.Subscription{}, err
	}

	return updatedSub, nil
}

func (s *ServiceStore) Sum(ctx context.Context, userId, serviceName string, from, to time.Time) (int, error) {

	return s.subscriptionStore.SumForPeriod(ctx, userId, serviceName, from, to)
}
