package watcher

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/bonedaddy/unibot/bclient"
	"github.com/bonedaddy/unibot/db"
	"github.com/bonedaddy/unibot/discord"
	"github.com/bonedaddy/unibot/utils"
)

// Service provides a price watcher service that updates a database
type Service struct {
	wg     *sync.WaitGroup
	db     *db.Database
	bc     *bclient.Client
	ctx    context.Context
	cancel context.CancelFunc
	period time.Duration
	items  []WatchItem
}

type WatchItem struct {
	Token0   string
	Token1   string
	Decimals int
}

func ConfigToWatchItmes(cfg *discord.Config) []WatchItem {
	items := make([]WatchItem, 0, len(cfg.Watchers))
	for _, watch := range cfg.Watchers {
		items = append(items, WatchItem{watch.Token0Address, watch.Token1Address, watch.Decimals})
	}
	return items
}

// New returns a new watcher service
func New(ctx context.Context, db *db.Database, bc *bclient.Client, tick time.Duration, watchItems []WatchItem) *Service {
	ctx, cancel := context.WithCancel(context.Background())
	return &Service{&sync.WaitGroup{}, db, bc, ctx, cancel, tick, watchItems}
}

func (s *Service) Start() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(s.period)
		defer ticker.Stop()
		for {
			select {
			case <-s.ctx.Done():
				return
			case <-ticker.C:
				for _, item := range s.items {
					price, err := s.bc.GetPrice(item.Token0, item.Token1)
					if err != nil {
						log.Printf("failed to get price for token0: %s token1:%s - %s\n", item.Token0, item.Token1, err)
						continue
					}
					log.Printf("token0: %s token1:%s - price: %v", item.Token0, item.Token1, price.Int64())
					priceF, _ := utils.ToDecimal(price, item.Decimals).Float64()
					if err := s.db.RecordPrice(item.Token0, item.Token1, priceF); err != nil {
						log.Printf("failed to record price for token0: %s token1: %s - %s\n", item.Token0, item.Token1, err)
						continue
					}
				}
			}
		}
	}()
}

func (s *Service) Stop() {
	s.cancel()
	s.wg.Wait()
	return
}
