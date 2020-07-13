package service

import (
	"context"
	"github.com/kniepok/corgi"
	"github.com/sirupsen/logrus"
)

type StarterService struct {
	subscriptionService corgi.SubscriptionService
	storage             corgi.SubscriptionStorage
	logger              *logrus.Logger
}

func NewStarterService(
	subscriptionService corgi.SubscriptionService,
	storage corgi.SubscriptionStorage,
) *StarterService {
	return &StarterService{
		subscriptionService: subscriptionService,
		storage:             storage,
		logger:              logrus.New(),
	}
}

func (r *StarterService) Start(ctx context.Context) error {
	allEntries, err := r.storage.GetAll(ctx)
	if err != nil {
		return err
	}

	for _, entry := range allEntries {
		err := r.subscriptionService.Subscribe(ctx, entry.Subscription)
		if err != nil {
			r.logger.Errorf("failed to subscribe :%s; err: %v", entry.User.Email, err.Error())
			continue
		}
	}

	return nil
}
