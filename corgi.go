package corgi

import (
	"context"
	"fmt"
)

type ScheduledSubscription struct {
	Subscription
	EntryID int
}

// Subscription describes to which user a notification will be sent at a given interval
type Subscription struct {
	User     User
	Interval string
}

type User struct {
	ID    string
	Email string
}

type SubscriptionService interface {
	//	Subscribe adds user to notify list
	Subscribe(ctx context.Context, request Subscription) error
	// Unsubscribe removes user subscription
	Unsubscribe(ctx context.Context, user User) error
}

type StarterService interface {
	Start(ctx context.Context) error
}

// Scheduler can
type Scheduler interface {
	// Add registers new notification schedule
	Add(ctx context.Context, interval string, task func()) (entryID int, err error)
	// Remove deletes already existing notification schedule
	Remove(ctx context.Context, entryID int) error
}

// SubscriptionStorage can return a list of already created subscriptions
type SubscriptionStorage interface {
	// Create or update subscription entry
	Add(ctx context.Context, sub ScheduledSubscription) error
	// Get scheduled subscription by user ID
	Get(ctx context.Context, userID string) (ScheduledSubscription, error)
	// Remove scheduled subscription based on user ID
	Remove(ctx context.Context, userID string) error
	// Fetch all scheduled subscriptions
	GetAll(ctx context.Context) ([]ScheduledSubscription, error)
}

type NotificationService interface {
	Notify(ctx context.Context, user User) error
}

// NotificationSender can send notifications to given users
type NotificationSender interface {
	Send(ctx context.Context, user User, message string) error
}

// NotificationGenerator can create notification messages
type NotificationGenerator interface {
	CreateNotification(ctx context.Context, user User) (msg string, err error)
}

// IntervalValidator can verify if given interval is semantically correct
type IntervalValidator interface {
	Validate(ctx context.Context, interval string) error
}

type UserResolver interface {
	ResolveUserEmail(ctx context.Context, userID string) (email string, err error)
}

var ErrInvalidInterval = fmt.Errorf("invalid interval")
var ErrSubscriptionNotFound = fmt.Errorf("cannot find subscription")
