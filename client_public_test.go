package dydx_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/huangjosh/go-dydx"
	"github.com/huangjosh/go-dydx/helpers"
	"github.com/huangjosh/go-dydx/public"
	"github.com/stretchr/testify/assert"
)

func TestGetMarkets(t *testing.T) {
	client := dydx.New(options)
	res, err := client.Public.GetMarkets("")
	assert.NoError(t, err)

	for k, v := range res.Markets {
		fmt.Printf("max position - %v\n", v.MaxPositionSize)
		fmt.Printf("%v - %#v\n", k, helpers.ToFloat(v.MinOrderSize)*helpers.ToFloat(v.IndexPrice))
	}
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
	start := -100 * 24 * time.Hour
	res, err := client.Public.GetCandles(&public.CandlesParam{
		Market:     "BTC-USD",
		Resolution: "5MINS",
		FromISO:    time.Now().UTC().Add(start).Format(time.RFC3339),
		ToISO:      time.Now().UTC().Add(start + time.Duration(5*100)*time.Minute).Format(time.RFC3339),
		Limit:      100,
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
