package service

import (
	"context"
	"github.com/kniepok/corgi"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type SubscriptionService struct {
	validator           corgi.IntervalValidator
	storage             corgi.SubscriptionStorage
	scheduler           corgi.Scheduler
	notificationService corgi.NotificationService
	logger              *logrus.Logger
}

func NewSubscriptionService(
	validator corgi.IntervalValidator,
	storage corgi.SubscriptionStorage,
	scheduler corgi.Scheduler,
	notificationService corgi.NotificationService) *SubscriptionService {
	return &SubscriptionService{
		validator:           validator,
		storage:             storage,
		scheduler:           scheduler,
		notificationService: notificationService,
		logger:              logrus.New(),
	}
}

func (s *SubscriptionService) Subscribe(ctx context.Context, subscription corgi.Subscription) error {
	if err := s.validator.Validate(ctx, subscription.Interval); err != nil {
		return corgi.ErrInvalidInterval
	}

	alreadyExisting, err := s.storage.Get(ctx, subscription.User.ID)
	if err != nil && errors.Cause(err) != corgi.ErrSubscriptionNotFound {
		return err
	}
	if alreadyExisting.EntryID != 0 {
		s.logger.Infof("overwriting subscription for %s", subscription.User.Email)
		err := s.Unsubscribe(ctx, subscription.User)
		if err != nil {
			return err
		}
	}

	entryID, err := s.scheduler.Add(ctx, subscription.Interval, func() {
		err := s.notificationService.Notify(ctx, subscription.User)
		if err != nil {
			s.logger.Errorf("failed to notify %s, err=%s", subscription.User.Email, err.Error())
		}
	})
	if err != nil {
		s.logger.Errorf("failed to schedule subscription for %s; err=%s", subscription.User.Email, err.Error())
		return err
	}

	err = s.storage.Add(ctx, corgi.ScheduledSubscription{
		Subscription: subscription,
		EntryID:      entryID,
	})
	if err != nil {
		_ = s.scheduler.Remove(ctx, entryID)
		return err
	}

	return nil
}

// Unsubscribe fetches entry from database, tries to stop scheduled job and removes subscription
func (s *SubscriptionService) Unsubscribe(ctx context.Context, user corgi.User) error {
	storedSub, err := s.storage.Get(ctx, user.ID)
	if err != nil {
		return err
	}

	if err = s.scheduler.Remove(ctx, storedSub.EntryID); err != nil {
		return err
	}

	return s.storage.Remove(ctx, user.ID)
}
