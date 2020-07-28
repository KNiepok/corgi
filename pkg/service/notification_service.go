package service

import (
	"context"
	"github.com/kniepok/corgi"
)

type NotificationService struct {
	messageGenerator corgi.NotificationGenerator
	messageSender    corgi.NotificationSender
}

func NewNotificationService(
	messageGenerator corgi.NotificationGenerator,
	messageSender corgi.NotificationSender) *NotificationService {
	return &NotificationService{
		messageGenerator: messageGenerator,
		messageSender:    messageSender,
	}
}

// Notify asks generator to create a message
// Then passes it to sender
func (n *NotificationService) Notify(ctx context.Context, sub corgi.Subscription) error {
	msg, err := n.messageGenerator.CreateNotification(ctx, sub)
	if err != nil {
		return err
	}

	return n.messageSender.Send(ctx, sub.User, msg)
}
