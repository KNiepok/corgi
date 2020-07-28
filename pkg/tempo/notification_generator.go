package tempo

import (
	"context"
	"fmt"
	"github.com/kniepok/corgi"
)

type NotificationGenerator struct {
}

func NewNotificationGenerator() *NotificationGenerator {
	return &NotificationGenerator{}
}

func (n NotificationGenerator) CreateNotification(ctx context.Context, sub corgi.Subscription) (string, error) {
	// todo use user-schedule API to fetch how many hours should a user log
	// todo user worklogs API to fetch how much user logged
	// make that a message or something. Maybe something more detailed, or even8 introduce a structure to
	// render the data properly
	return fmt.Sprintf("%s! remember to fill timesheets", sub.User.Email), nil
}
