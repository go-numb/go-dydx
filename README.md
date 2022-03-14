# go-dydx

dYdX exchange API version3.

part OnBoarding referred to [verichenn/dydx-v3-go](https://github.com/verichenn/dydx-v3-go).

## Description

go-dydx is a go client library for dYdX, [dYdX API Document](https://docs.dydx.exchange).

## Support
- [ ] private/accounts
- [ ] private/positions
- [ ] private/orders (get, post, delete)
- [ ] public/markets
- [ ] public/orderbook
- [ ] public/trades

## Usege
```go
package main

import (
	"fmt"
	"time"

	"github.com/labstack/gommon/log"

	"github.com/go-numb/go-dydx"
	"github.com/go-numb/go-dydx/helpers"
	"github.com/go-numb/go-dydx/private"
	"github.com/go-numb/go-dydx/types"
)

const (
	EthereumAddress = "0xtest"
)

var userID int64 = 11111
var options = types.Options{
	Host:                      types.ApiHostMainnet,
	StarkPublicKey:            "<please check Google Chrome Developer tool -> application starkkey>",
	StarkPrivateKey:           "<please check Google Chrome Developer tool -> application starkkey>",
	StarkPublicKeyYCoordinate: "<please check Google Chrome Developer tool -> application starkkey>",
	DefaultEthereumAddress:    EthereumAddress,
	ApiKeyCredentials: &types.ApiKeyCredentials{
		Key:        "<please check Google Chrome Developer tool -> application apikey>",
		Secret:     "<please check Google Chrome Developer tool -> application secret>",
		Passphrase: "<please check Google Chrome Developer tool -> application passphrase>",
	},
}

func main() {
	start := time.Now()
	defer func() {
		fmt.Println("exec time: ", time.Since(start))
	}()

	client := dydx.New(options)
	account, err := client.Private.GetAccount(EthereumAddress)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(account)

    // print rate limit numbers
	fmt.Println(client.Public.RateLimit.ToNumber())


	params := &private.ApiOrder{
		ApiBaseOrder: private.ApiBaseOrder{Expiration: helpers.ExpireAfter(5 * time.Minute)},
		Market:       "ETH-USD",
		Side:         "BUY",
		Type:         "LIMIT",
		Size:         "1",
		Price:        "2000",
		ClientId:     helpers.RandomClientId(),
		TimeInForce:  "GTT",
		PostOnly:     true,
		LimitFee:     "0.01",
	}
	res, err := client.Private.CreateOrder(params, userID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)

	// rate limit updated above
    // print rate limit numbers
	fmt.Println(client.Private.RateLimit.ToNumber())
}

```

## Author

[@_numbP](https://twitter.com/_numbP)

## License

[MIT](https://github.com/go-numb/go-dydx/blob/master/LICENSE)