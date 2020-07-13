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
	Interval  string
	EntryID   int
}

func toModel(model corgi.ScheduledSubscription) subscription {
	return subscription{
		UserID:    model.User.ID,
		UserEmail: model.User.Email,
		Interval:  model.Interval,
		EntryID:   model.EntryID,
	}
}

func (s *subscription) fromModel() corgi.ScheduledSubscription {
	return corgi.ScheduledSubscription{
		Subscription: corgi.Subscription{
			User: corgi.User{
				ID:    s.UserID,
				Email: s.UserEmail,
			},
			Interval: s.Interval,
		},
		EntryID: s.EntryID,
	}
}
