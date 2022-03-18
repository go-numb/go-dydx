package dydx_test

import (
	"fmt"
	"testing"

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

func TestGetHistoricalFunding(t *testing.T) {
	client := dydx.New(options)
	res, err := client.Public.GetHistoricalFunding(&public.HistoricalFundingsParam{
		MarketID: "BTC-USD",
	})
	assert.NoError(t, err)
	fmt.Printf("%v", res)
}
