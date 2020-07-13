package cron

import (
	"context"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

// Scheduler can user subscriptions to schedule notifications
type Scheduler struct {
	runner *cron.Cron
	logger *logrus.Logger
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		runner: cron.New(),
		logger: logrus.New(),
	}
}

func (s *Scheduler) Add(ctx context.Context, interval string, task func()) (int, error) {
	spec, err := parseIntervalToCron(interval)
	if err != nil {
		return 0, err
	}
	entryID, err := s.runner.AddFunc(spec, task)
	if err != nil {
		return 0, err
	}

	s.logger.Infof("scheduler: registering cron '%s' with entry id=%d", spec, entryID)
	return int(entryID), nil
}

func (s *Scheduler) Remove(ctx context.Context, entryID int) error {
	s.runner.Remove(cron.EntryID(entryID))
	s.logger.Infof("scheduler: removing entry id=%d", entryID)
	return nil
}
