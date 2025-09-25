package service

import (
	"context"
	datatransfer "subscription/internal/dto"
	"subscription/internal/model"
	"time"
)

type ServiceRepository interface {
	Create(ctx context.Context, dto datatransfer.DTOSubs) (model.Subscription, error)
	GetInfo(ctx context.Context, idSub string) (model.Subscription, error)
	GetAll(ctx context.Context) ([]model.Subscription, error)
	Delete(ctx context.Context, idSub string) error
	Update(ctx context.Context, id string, dto datatransfer.DTOSubs) (model.Subscription, error)
	Sum(ctx context.Context, userId, serviceName string, from, to time.Time) (int, error)
}

type ServiceStore struct {
	subscriptionStore model.SubscriptionRepository
}

func NewService(subStore model.SubscriptionRepository) *ServiceStore {
	return &ServiceStore{
		subscriptionStore: subStore,
	}
}
func (s *ServiceStore) Create(ctx context.Context, dto datatransfer.DTOSubs) (model.Subscription, error) {

	sub, err := model.NewSubscription(dto)
	if err != nil {
		return model.Subscription{}, err
	}
	if err := s.subscriptionStore.AddSub(ctx, sub); err != nil {
		return model.Subscription{}, err
	}

	return sub, nil

}

func (s *ServiceStore) GetInfo(ctx context.Context, idSub string) (model.Subscription, error) {

	sub, err := s.subscriptionStore.GetSubInfo(ctx, idSub)
	if err != nil {
		return model.Subscription{}, err
	}
	return sub, nil
}

func (s *ServiceStore) GetAll(ctx context.Context) ([]model.Subscription, error) {

	allSub, err := s.subscriptionStore.GetSubAllInfo(ctx)
	if err != nil {
		return []model.Subscription{}, err
	}
	return allSub, nil
}

func (s *ServiceStore) Delete(ctx context.Context, idSub string) error {
	return s.subscriptionStore.DeleteInfo(ctx, idSub)
}

func (s *ServiceStore) Update(ctx context.Context, id string, dto datatransfer.DTOSubs) (model.Subscription, error) {

	oldSub, err := s.subscriptionStore.GetSubInfo(ctx, id)
	if err != nil {
		return model.Subscription{}, err
	}
	updatedSub := model.Subscription{
		ServiceName: dto.ServiceName,
		Price:       dto.Price,
		UserId:      oldSub.UserId,
		StartDate:   oldSub.StartDate,
	}

	if err := s.subscriptionStore.UpdateSub(ctx, id, updatedSub); err != nil {
		return model.Subscription{}, err
	}

	return updatedSub, nil
}

func (s *ServiceStore) Sum(ctx context.Context, userId, serviceName string, from, to time.Time) (int, error) {
	fromCD := model.CustomDate{Time: from}
	toCD := model.CustomDate{Time: to}
	return s.subscriptionStore.SumSubscriptions(ctx, userId, serviceName, fromCD, toCD)
}
