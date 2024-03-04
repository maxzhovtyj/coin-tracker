package sqlc

import (
	"context"
	"fmt"
	db "github.com/maxzhovtyj/coin-tracker/pkg/db/sqlc"
	"time"
)

type SubscriptionStorage struct {
	q *db.Queries
}

func (s *SubscriptionStorage) UpdateLastNotifiedAt(id int64) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelFunc()

	updated, err := s.q.UpdateLastNotifiedAt(ctx, id)
	if err != nil {
		return err
	}

	if updated.ID == 0 {
		return fmt.Errorf("no rows updated")
	}

	return nil
}

func (s *SubscriptionStorage) All() ([]db.Subscription, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	return s.q.GetSubscriptions(ctx)
}

func NewSubscriptionStorage(conn db.DBTX) *SubscriptionStorage {
	return &SubscriptionStorage{
		q: db.New(conn),
	}
}

func (s *SubscriptionStorage) Create(uid int64, subscriptionType, data, interval string) (db.Subscription, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelFunc()

	subscription, err := s.q.CreateSubscription(ctx, db.CreateSubscriptionParams{
		Type:           subscriptionType,
		UserID:         uid,
		Data:           data,
		NotifyInterval: interval,
	})
	if err != nil {
		return db.Subscription{}, err
	}

	return subscription, nil
}

func (s *SubscriptionStorage) UserSubscriptions(uid int64) ([]db.Subscription, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelFunc()

	return s.q.GetUserSubscriptions(ctx, uid)
}
