package sqlite

import (
	"context"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/kniepok/corgi"
	"github.com/pkg/errors"
)

type SubscriptionStorage struct {
	db *gorm.DB
}

func NewSubscriptionStorage(db *gorm.DB) *SubscriptionStorage {
	db.AutoMigrate(&subscription{})
	return &SubscriptionStorage{
		db: db,
	}
}

// todo make that work with new details
func (s *SubscriptionStorage) Add(ctx context.Context, sub corgi.ScheduledSubscription) error {
	var alreadyExisting subscription
	err := s.db.Model(&alreadyExisting).Where("user_id = ?", sub.User.ID).Find(&alreadyExisting).Error
	if err != nil {
		if errors.Cause(err) == gorm.ErrRecordNotFound {
			newSub := toModel(sub)
			return s.db.Create(&newSub).Error
		}
		return err
	}

	alreadyExisting.EntryID = sub.EntryID
	alreadyExisting.Interval = sub.Interval

	return s.db.Save(&alreadyExisting).Error
}

func (s *SubscriptionStorage) Get(ctx context.Context, userID string) (corgi.ScheduledSubscription, error) {
	sub := subscription{}
	err := s.db.Model(&subscription{}).Where("user_id = ?", userID).First(&sub).Error
	if err != nil {
		if errors.Cause(err) == gorm.ErrRecordNotFound {
			return corgi.ScheduledSubscription{}, corgi.ErrSubscriptionNotFound
		}
		return corgi.ScheduledSubscription{}, err
	}
	return sub.fromModel(), nil
}

func (s *SubscriptionStorage) Remove(ctx context.Context, userID string) error {
	return s.db.Where("user_id = ?", userID).Delete(&subscription{}).Error
}

func (s *SubscriptionStorage) GetAll(ctx context.Context) ([]corgi.ScheduledSubscription, error) {
	var records []subscription

	err := s.db.Find(&records).Error
	if err != nil {
		return nil, err
	}
	models := make([]corgi.ScheduledSubscription, 0, len(records))
	for _, r := range records {
		models = append(models, r.fromModel())
	}
	return models, nil
}
