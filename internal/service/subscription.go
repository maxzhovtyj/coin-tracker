package service

import (
	"encoding/json"
	"fmt"
	"github.com/maxzhovtyj/coin-tracker/internal/models"
	"github.com/maxzhovtyj/coin-tracker/internal/storage"
	"github.com/maxzhovtyj/coin-tracker/pkg/binance"
	"time"
)

const (
	CoinSubscriptionType = "coin"
)

type SubscriptionService struct {
	db  storage.Subscription
	api binance.API
}

func NewSubscriptionService(db storage.Subscription, api binance.API) *SubscriptionService {
	return &SubscriptionService{
		db:  db,
		api: api,
	}
}

func (s *SubscriptionService) UserSubscriptions(uid int64) ([]models.Subscription, error) {
	subscriptions, err := s.db.UserSubscriptions(uid)
	if err != nil {
		return nil, err
	}

	subs := make([]models.Subscription, len(subscriptions))
	for i, sub := range subscriptions {
		d, err := time.ParseDuration(sub.NotifyInterval)
		if err != nil {
			return nil, err
		}

		subs[i] = models.Subscription{
			ID:             sub.ID,
			Type:           sub.Type,
			ChatID:         sub.UserID,
			Data:           sub.Data,
			NotifyInterval: d,
			LastNotifiedAt: sub.LastNotifiedAt.Time,
		}
	}

	return subs, nil
}

func (s *SubscriptionService) CoinTicker(coin, windowSize string) (binance.SymbolTicker, error) {
	return s.api.CoinTicker(coin, windowSize)
}

func (s *SubscriptionService) Notified(id int64) error {
	return s.db.UpdateLastNotifiedAt(id)
}

func (s *SubscriptionService) All() ([]models.Subscription, error) {
	all, err := s.db.All()
	if err != nil {
		return nil, err
	}

	subs := make([]models.Subscription, len(all))
	for i, sub := range all {
		d, err := time.ParseDuration(sub.NotifyInterval)
		if err != nil {
			return nil, err
		}

		subs[i] = models.Subscription{
			ID:             sub.ID,
			Type:           sub.Type,
			ChatID:         sub.UserID,
			Data:           sub.Data,
			NotifyInterval: d,
			LastNotifiedAt: sub.LastNotifiedAt.Time,
		}
	}

	return subs, nil
}

func (s *SubscriptionService) NewCoinSubscription(uid int64, coinName string, interval time.Duration) error {
	var windowSize string
	if interval.Minutes() < 60 {
		windowSize = fmt.Sprintf("%dh", int(interval.Minutes()))
	} else if interval.Hours() < 24 {
		windowSize = fmt.Sprintf("%dh", int(interval.Hours()))
	} else {
		windowSize = fmt.Sprintf("%dd", int(interval.Hours()/24))
	}

	dataRaw, err := json.Marshal(models.CoinSubscriptionData{
		Interval: windowSize,
		CoinName: coinName,
	})
	if err != nil {
		return fmt.Errorf("cant marshal data: %w", err)
	}

	_, err = s.db.Create(uid, CoinSubscriptionType, string(dataRaw), interval.String())
	if err != nil {
		return err
	}

	return nil
}
