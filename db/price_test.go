package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrice(t *testing.T) {
	t.Cleanup(func() {
		os.Remove("indexed.db")
	})
	db := newTestDB(t)
	t.Run("RecordPrice", func(t *testing.T) {
		type args struct {
			token0 string
			token1 string
			price  float64
		}
		tests := []struct {
			name    string
			args    args
			wantErr bool
		}{
			{"AB", args{"a", "b", 2.132}, false},
			{"BC", args{"b", "c", 3.1434}, false},
			{"DE", args{"d", "e", 4.123}, false},
			{"EF", args{"e", "f", 4.123}, false},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := db.RecordPrice(tt.args.token0, tt.args.token1, tt.args.price)
				if (err != nil) != tt.wantErr {
					t.Fatalf("RecordPrice() err %v, wantErr %v", err, tt.wantErr)
				}
			})
		}
	})
	t.Run("LastPrice", func(t *testing.T) {
		type args struct {
			token0      string
			token1      string
			firstPrice  float64
			secondPrice float64
		}
		tests := []struct {
			name    string
			args    args
			wantErr bool
		}{
			{"AB", args{"a", "b", 10.101, 11.23}, false},
			{"BC", args{"b", "c", 12.121, 13.31}, false},
			{"DE", args{"d", "e", 14.141, 15.81}, false},
			{"EF", args{"e", "f", 14.141, 15.81}, false},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := db.RecordPrice(tt.args.token0, tt.args.token1, tt.args.firstPrice)
				if (err != nil) != tt.wantErr {
					t.Fatalf("RecordPrice() err %v, wantErr %v", err, tt.wantErr)
				}
				if tt.wantErr {
					return
				}
				price, err := db.LastPrice(tt.args.token0, tt.args.token1)
				require.NoError(t, err)
				require.Equal(t, price, tt.args.firstPrice)

				require.NoError(t, db.RecordPrice(tt.args.token0, tt.args.token1, tt.args.secondPrice))

				price, err = db.LastPrice(tt.args.token0, tt.args.token1)
				require.NoError(t, err)
				require.Equal(t, price, tt.args.secondPrice)
			})
		}
	})
	t.Run("GetAllPrices", func(t *testing.T) {
		type args struct {
			token0      string
			token1      string
			wantEntries int
		}
		tests := []struct {
			name    string
			args    args
			wantErr bool
		}{
			{"AB", args{"a", "b", 3}, false},
			{"BC", args{"b", "c", 3}, false},
			{"DE", args{"d", "e", 3}, false},
			{"EF", args{"e", "f", 3}, false},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				entries, err := db.GetAllPrices(tt.args.token0, tt.args.token1)
				if (err != nil) != tt.wantErr {
					t.Fatalf("GetAllPrices() err %v, wantErr %v", err, tt.wantErr)
				}
				if tt.wantErr {
					return
				}
				require.Len(t, entries, tt.args.wantEntries)
			})
		}
	})

	t.Run("PriceAvgInRange", func(t *testing.T) {
		type args struct {
			token0    string
			token1    string
			window    int
			wantPrice float64
		}
		tests := []struct {
			name    string
			args    args
			wantErr bool
		}{
			{"AB", args{"a", "b", 1, 7.821000000000001}, false},
			{"BC", args{"b", "c", 1, 9.5248}, false},
			{"DE", args{"d", "e", 1, 11.357999999999999}, false},
			{"EF", args{"e", "f", 1, 11.357999999999999}, false},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				priceAvg, err := db.PriceAvgInRange(tt.args.token0, tt.args.token1, tt.args.window)
				if (err != nil) != tt.wantErr {
					t.Fatalf("GetAllPrices() err %v, wantErr %v", err, tt.wantErr)
				}
				if tt.wantErr {
					return
				}
				entries, err := db.GetAllPrices(tt.args.token0, tt.args.token1)
				require.NoError(t, err)
				var totalPrice float64
				for _, entry := range entries {
					totalPrice += entry.USDPrice
				}
				avgPrice := totalPrice / float64(len(entries))
				require.Equal(t, avgPrice, priceAvg)
				require.Equal(t, tt.args.wantPrice, priceAvg)

				// ensure recording a new price changes the average
				require.NoError(t, db.RecordPrice(tt.args.token0, tt.args.token1, 19))
				newPriceAvg, err := db.PriceAvgInRange(tt.args.token0, tt.args.token1, tt.args.window)
				require.NotEqual(t, newPriceAvg, priceAvg)
			})
		}
	})
	t.Run("PriceChangeInRange", func(t *testing.T) {
		type args struct {
			token0     string
			token1     string
			window     int
			wantChange float64
		}
		tests := []struct {
			name    string
			args    args
			wantErr bool
		}{
			{"AB", args{"a", "b", 1, 7.911819887429642}, false},
			{"BC", args{"b", "c", 1, 5.044410510911751}, false},
			{"DE", args{"d", "e", 1, 3.6082949308755756}, false},
			{"EF", args{"e", "f", 1, 3.6082949308755756}, false},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				change, err := db.PriceChangeInRange(tt.args.token0, tt.args.token1, tt.args.window)
				if (err != nil) != tt.wantErr {
					t.Fatalf("PriceChangeInRange() err %v, wantErr %v", err, tt.wantErr)
				}
				if tt.wantErr {
					return
				}
				require.Equal(t, tt.args.wantChange, change)

				// get the first price in the window
				prices, err := db.windowRangeQuery(tt.args.token0, tt.args.token1, tt.args.window)
				require.NoError(t, err)
				require.GreaterOrEqual(t, len(prices), 1)
				firstPrice := prices[0].USDPrice
				currPrice := prices[len(prices)-1].USDPrice
				// now record a price lower than first price to enforce negative percent change
				// we reduce its value to 2 less than firs the first price
				toReduce := (currPrice - firstPrice) + (firstPrice) + 0.123
				require.NoError(t, db.RecordPrice(tt.args.token0, tt.args.token1, currPrice-toReduce))

				// recalculate the price change
				newChange, err := db.PriceChangeInRange(tt.args.token0, tt.args.token1, tt.args.window)
				require.NoError(t, err)
				require.NotEqual(t, newChange, change)
				require.Less(t, newChange, change)

			})
		}
	})
}
