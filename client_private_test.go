package dydx_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	dydx "github.com/go-numb/go-dydx"
	"github.com/go-numb/go-dydx/helpers"
	"github.com/go-numb/go-dydx/private"
	"github.com/go-numb/go-dydx/realtime"
	"github.com/go-numb/go-dydx/types"
	"github.com/stretchr/testify/assert"
)

const (
	DefaultHost     = "http://localhost:8080"
	EthereumAddress = ""
	StarkKey        = ""
)

var userID int64 = 1

var options = types.Options{
	Host:                      types.ApiHostMainnet,
	StarkPublicKey:            "",
	StarkPrivateKey:           "",
	StarkPublicKeyYCoordinate: "",
	DefaultEthereumAddress:    EthereumAddress,
	ApiKeyCredentials: &types.ApiKeyCredentials{
		Key:        "",
		Secret:     "",
		Passphrase: "",
	},
}

func TestConnect(t *testing.T) {
	client := dydx.New(options)

	account, err := client.Private.GetAccount(client.Private.DefaultAddress)
	assert.NoError(t, err)

	fmt.Println(account)

	parent := context.Background()
	ctx, cancel := context.WithCancel(parent)

	r := make(chan realtime.Response)

	go realtime.Connect(ctx, r, []string{realtime.ACCOUNT, realtime.TRADES}, []string{"BTC-USD"}, client.Private, nil)

	for {
		select {
		case v := <-r:
			switch v.Channel {
			case realtime.ACCOUNT:
				fmt.Println(v.Account)
			case realtime.TRADES:
				fmt.Println(v.Trades)
			case realtime.ERROR:
				log.Println(v.Results)
				goto EXIT
			}

		}
	}

EXIT:
	cancel()
}

func TestUsers(t *testing.T) {
	client := dydx.New(options)
	res, err := client.Private.GetUsers()
	assert.NoError(t, err)

	fmt.Printf("%v", res)

	fmt.Printf("makerFee: %s, takeFee: %s\n", res.User.MakerFeeRate, res.User.TakerFeeRate)
}

func TestCreateOrder(t *testing.T) {
	client := dydx.New(options)
	o := &private.ApiOrder{
		ApiBaseOrder: private.ApiBaseOrder{Expiration: helpers.ExpireAfter(5 * time.Minute)},
		Market:       "ETH-USD",
		Side:         "BUY",
		Type:         "LIMIT",
		Size:         "1",
		Price:        "2500",
		ClientId:     helpers.RandomClientId(),
		TimeInForce:  "GTT",
		PostOnly:     true,
		LimitFee:     "0.01",
	}
	fmt.Printf("%+v\n", client.Private.NetworkId)
	res, err := client.Private.CreateOrder(o, userID)
	assert.NoError(t, err)

	fmt.Printf("%v", res)
}

// important!! WithDraw has not done any actual testing
func TestWithdrawFast(t *testing.T) {
	client := dydx.New(options)
	res, err := client.Private.WithdrawFast(&private.WithdrawalParam{})
	assert.NoError(t, err)

	fmt.Printf("%v", res)
}

func TestGetHistoricalPnL(t *testing.T) {
	client := dydx.New(options)
	res, err := client.Private.GetHistoricalPnL(&private.HistoricalPnLParam{})
	assert.NoError(t, err)

	fmt.Printf("%v", res)
}

func TestGetTradingRewards(t *testing.T) {
	client := dydx.New(options)
	res, err := client.Private.GetTradingRewards(&private.TradingRewardsParam{
		Epoch: 8,
	})
	assert.NoError(t, err)

	fmt.Printf("%v", res)
}
