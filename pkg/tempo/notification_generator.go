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

// todo make that real, for its a dummy
func (n NotificationGenerator) CreateNotification(ctx context.Context, user corgi.User) (string, error) {
	return fmt.Sprintf("%s! remember to fill timesheets", user.Email), nil
}
