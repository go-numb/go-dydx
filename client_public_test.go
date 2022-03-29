package dydx_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-numb/go-dydx"
	"github.com/go-numb/go-dydx/public"
	"github.com/stretchr/testify/assert"
)

func TestGetMarkets(t *testing.T) {
	client := dydx.New(options)
	res, err := client.Public.GetMarkets("BTC-USD")
	assert.NoError(t, err)
	fmt.Printf("%v", res)
}

func TestGetTrades(t *testing.T) {
	client := dydx.New(options)
	res, err := client.Public.GetTrades(&public.TradesParam{
		MarketID: "BTC-USD",
	})
	assert.NoError(t, err)
	fmt.Printf("%v", res)
}

func TestGetCandles(t *testing.T) {
	client := dydx.New(options)
	res, err := client.Public.GetCandles(&public.CandlesParam{
		Market:     "BTC-USD",
		Resolution: "1MIN",
		FromISO:    time.Now().UTC().Add(-1*time.Minute - 24*time.Hour).Format(time.RFC3339),
		ToISO:      time.Now().UTC().Add(-24 * time.Hour).Format(time.RFC3339),
	})

	assert.NoError(t, err)
	fmt.Printf("%v, lenght: %d\n", res, len(res.Candles))
}
func TestGetHistoricalFunding(t *testing.T) {
	client := dydx.New(options)
	res, err := client.Public.GetHistoricalFunding(&public.HistoricalFundingsParam{
		Market: "BTC-USD",
	})
	assert.NoError(t, err)
	fmt.Printf("%v", res)
}
