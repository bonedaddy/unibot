package db

import (
	"errors"
	"math"
	"time"

	"gorm.io/gorm"
)

// Price is a given price entry for an asset
type Price struct {
	gorm.Model
	Token0   string
	Token1   string
	USDPrice float64
}

// RecordPrice records the given asset price in the database
func (d *Database) RecordPrice(token0 string, token1 string, price float64) error {
	return d.db.Create(&Price{Token0: token0, Token1: token1, USDPrice: price}).Error
}

// LastPrice returns the last recorded price
func (d *Database) LastPrice(token0, token1 string) (float64, error) {
	var price Price
	if err := d.db.Model(&Price{}).Where("token0 = ? AND token1 = ?", token0, token1).Last(&price).Error; err != nil {
		return 0, err
	}
	return price.USDPrice, nil
}

// GetAllPrices returns all price entries for a given asset
func (d *Database) GetAllPrices(token0, token1 string) ([]*Price, error) {
	var prices []*Price
	return prices, d.db.Model(&Price{}).Where("token0 = ? AND token1 = ?", token0, token1).Find(&prices).Error
}

// PriceAvgInRange returns the average price of the given asset during the last N days
func (d *Database) PriceAvgInRange(token0, token1 string, windowInDays int) (float64, error) {
	prices, err := d.windowRangeQuery(token0, token1, windowInDays)
	if err != nil {
		return 0, err
	}
	var totalPrice float64
	for _, price := range prices {
		totalPrice += price.USDPrice
	}
	return totalPrice / float64(len(prices)), nil
}

// PriceChangeInRange returns the price change percentage in the last N days
func (d *Database) PriceChangeInRange(token0, token1 string, windowInDays int) (float64, error) {
	prices, err := d.windowRangeQuery(token0, token1, windowInDays)
	if err != nil {
		return 0, err
	}
	switch len(prices) {
	case 0:
		return 0, errors.New("no prices found")
	case 1:
		return 0, nil // no price change
	default:
	}
	startPrice := prices[0].USDPrice
	finalPrice := prices[len(prices)-1].USDPrice
	// get the percentage change
	percentChange := ((finalPrice - startPrice) / math.Abs(startPrice))
	return percentChange, nil
}

func (d *Database) windowRangeQuery(token0, token1 string, windowInDays int) ([]*Price, error) {
	windowEnd := time.Now()
	windowStart := windowEnd.AddDate(0, 0, -windowInDays)
	var prices []*Price
	return prices, d.db.Model(&Price{}).Where(
		"token0 = ? AND token1 = ? AND created_at BETWEEN ? AND ?",
		token0, token1, windowStart, windowEnd,
	).Find(&prices).Error
}
