package types

import (
	"strconv"

	"github.com/umbracle/ethgo/jsonrpc"
)

const (
	ApiHostMainnet = "https://api.dydx.exchange"
	ApiHostRopsten = "https://api.stage.dydx.exchange"
	WsHostMainnet  = "wss://api.dydx.exchange/v3/ws"
	WsHostRopsten  = "wss://api.stage.dydx.exchange/v3/ws"
)

const (
	SignatureTypeNoPrepend   = 0
	SignatureTypeDecimal     = 1
	SignatureTypeHexadecimal = 2
)

const (
	Domain                       = "dYdX"
	Version                      = "1.0"
	Eip712DomainStringNoContract = "EIP712Domain(string name,string version,uint256 chainId)"
)

const (
	OffChainOnboardingAction    = "dYdX Onboarding"
	OffChainKeyDerivationAction = "dYdX STARK Key"
)

const (
	NetworkIdMainnet = 1
	NetworkIdRopsten = 3
)

const (
	OPEN        = "OPEN"
	CLOSED      = "CLOSED"
	LIQUIDATED  = "LIQUIDATED"
	LIQUIDATION = "LIQUIDATION"
	UNTRIGGERED = "UNTRIGGERED"
)

const (
	MARKET       = "MARKET"
	LIMIT        = "LIMIT"
	STOP         = "STOP"
	STOPLIMIT    = "STOP_LIMIT"
	TRAILINGSTOP = "TRAILING_STOP"
	TAKEPROFIT   = "TAKE_PROFIT"
)

const (
	BUY  = "BUY"
	SELL = "SELL"
)

const (
	LONG  = "LONG"
	SHORT = "SHORT"
)

const (
	TimeInForceGtt = "GTT"
	TimeInForceFok = "FOK"
	TimeInForceIoc = "IOC"
)

const (
	OrderStatusPending     = "PENDING"
	OrderStatusOpen        = "OPEN"
	OrderStatusFilled      = "FILLED"
	OrderStatusCanceled    = "CANCELED"
	OrderStatusUntriggered = "UNTRIGGERED"
)

const (
	Resolution1D     = "1DAY"
	Resolution4HOURS = "4HOURS"
	Resolution1HOUR  = "1HOUR"
	Resolution30MINS = "30MINS"
	Resolution15MINS = "15MINS"
	Resolution5MINS  = "5MINS"
	Resolution1MIN   = "1MIN"
)

type Options struct {
	NetworkId              int
	Host                   string
	DefaultEthereumAddress string

	StarkPublicKey            string
	StarkPrivateKey           string
	StarkPublicKeyYCoordinate string
	ApiKeyCredentials         *ApiKeyCredentials

	Web3 *jsonrpc.Client
}

type ApiKeyCredentials struct {
	Key        string
	Secret     string
	Passphrase string
}

type RateLimit struct {
	Remaining  string
	Reset      string
	RetryAfter string
	Limit      string
}

func (p *RateLimit) ToNumber() (remaining, reset, retryAfter, limit int64) {
	remaining, _ = strconv.ParseInt(p.Remaining, 10, 64)
	reset, _ = strconv.ParseInt(p.Reset, 10, 64)
	retryAfter, _ = strconv.ParseInt(p.RetryAfter, 10, 64)
	limit, _ = strconv.ParseInt(p.Limit, 10, 64)
	return
}
