package service

import (
	"context"
	datatransfer "subscription/internal/dto"
	"subscription/internal/model"
)

type ServiceRepository interface {
	Create(ctx context.Context, dto datatransfer.DTOSubs) (model.Subscription, error)
	GetInfo(ctx context.Context, idSub string) (model.Subscription, error)
}

type ServiceStore struct {
	subscriptionsStore model.SubscriptionRepository
}

func NewService(subStore model.SubscriptionRepository) *ServiceStore {
	return &ServiceStore{
		subscriptionsStore: subStore,
	}
}
func (s *ServiceStore) Create(ctx context.Context, dto datatransfer.DTOSubs) (model.Subscription, error) {

	sub, err := model.NewSubscription(dto)
	if err != nil {
		return model.Subscription{}, err
	}
	if err := s.subscriptionsStore.AddSub(ctx, sub); err != nil {
		return model.Subscription{}, err
	}

	return sub, nil

}

func (s *ServiceStore) GetInfo(ctx context.Context, idSub string) (model.Subscription, error) {

	sub, err := s.subscriptionsStore.GetSubInfo(ctx, idSub)
	if err != nil {
		return model.Subscription{}, err
	}
	return sub, nil
}
