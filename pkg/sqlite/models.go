package sqlite

import (
	"github.com/kniepok/corgi"
	"time"
)

type subscription struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    string `gorm:"unique;not null"`
	UserEmail string `gorm:"unique;not null"`
	EntryID   int
	Mode      int
	Weekday   int
	Hour      uint8
	Minute    uint8
}

func toModel(model corgi.ScheduledSubscription) subscription {
	return subscription{
		UserID:    model.User.ID,
		UserEmail: model.User.Email,
		EntryID:   model.EntryID,
		Mode:      int(model.Details.Mode),
		Weekday:   int(model.Details.DayOfWeek),
		Hour:      model.Details.Hour,
		Minute:    model.Details.Minute,
	}
}

func (s *subscription) fromModel() corgi.ScheduledSubscription {
	return corgi.ScheduledSubscription{
		Subscription: corgi.Subscription{
			User: corgi.User{
				ID:    s.UserID,
				Email: s.UserEmail,
			},
			Details: corgi.SubscriptionDetails{
				Mode:      corgi.SubscriptionMode(s.Mode),
				DayOfWeek: time.Weekday(s.Weekday),
				Hour:      s.Hour,
				Minute:    s.Minute,
			},
		},
		EntryID: s.EntryID,
	}
}
