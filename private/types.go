package private

import (
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/go-numb/go-dydx/types"
)

type Private struct {
	NetworkId         int
	Host              string
	StarkPrivateKey   string
	DefaultAddress    string
	ApiKeyCredentials *types.ApiKeyCredentials

	RateLimit *types.RateLimit
	Logger    *log.Logger
}

type ApiBaseOrder struct {
	Signature  string `json:"signature"`
	Expiration string `json:"expiration"`
}

type ApiOrder struct {
	ApiBaseOrder
	Market          string `json:"market"`
	Side            string `json:"side"`
	Type            string `json:"type"`
	Size            string `json:"size"`
	Price           string `json:"price"`
	ClientId        string `json:"clientId"`
	TimeInForce     string `json:"timeInForce"`
	LimitFee        string `json:"limitFee"`
	CancelId        string `json:"cancelId,omitempty"`
	TriggerPrice    string `json:"triggerPrice,omitempty"`
	TrailingPercent string `json:"trailingPercent,omitempty"`
	PostOnly        bool   `json:"postOnly"`
}

type AccountResponse struct {
	Account Account `json:"account"`
}

type Account struct {
	PositionId         int64               `json:"positionId,string"`
	ID                 string              `json:"id"`
	StarkKey           string              `json:"starkKey"`
	Equity             string              `json:"equity"`
	FreeCollateral     string              `json:"freeCollateral"`
	QuoteBalance       string              `json:"quoteBalance"`
	PendingDeposits    string              `json:"pendingDeposits"`
	PendingWithdrawals string              `json:"pendingWithdrawals"`
	AccountNumber      string              `json:"accountNumber"`
	OpenPositions      map[string]Position `json:"openPositions"`
	CreatedAt          time.Time           `json:"createdAt"`
}

type PositionResponse struct {
	Positions []Position `json:"positions"`
}

type Position struct {
	Market        string      `json:"market"`
	Status        string      `json:"status"`
	Side          string      `json:"side"`
	Size          string      `json:"size"`
	MaxSize       string      `json:"maxSize"`
	EntryPrice    string      `json:"entryPrice"`
	ExitPrice     interface{} `json:"exitPrice"`
	UnrealizedPnl string      `json:"unrealizedPnl"`
	RealizedPnl   string      `json:"realizedPnl"`
	CreatedAt     time.Time   `json:"createdAt"`
	ClosedAt      interface{} `json:"closedAt"`
	NetFunding    string      `json:"netFunding"`
	SumOpen       string      `json:"sumOpen"`
	SumClose      string      `json:"sumClose"`
}

type OrderResponse struct {
	Order Order `json:"order"`
}

type CancelOrderResponse struct {
	CancelOrder Order `json:"cancelOrder"`
}

type CancelOrdersResponse struct {
	CancelOrders Order `json:"cancelOrder"`
}

type Order struct {
	ID              string    `json:"id"`
	ClientID        string    `json:"clientId"`
	AccountID       string    `json:"accountId"`
	Market          string    `json:"market"`
	Side            string    `json:"side"`
	Price           string    `json:"price"`
	TriggerPrice    string    `json:"triggerPrice"`
	TrailingPercent string    `json:"trailingPercent"`
	Size            string    `json:"size"`
	RemainingSize   string    `json:"remainingSize"`
	Type            string    `json:"type"`
	UnfillableAt    string    `json:"unfillableAt"`
	Status          string    `json:"status"`
	TimeInForce     string    `json:"timeInForce"`
	CancelReason    string    `json:"cancelReason"`
	PostOnly        bool      `json:"postOnly"`
	CreatedAt       time.Time `json:"createdAt"`
	ExpiresAt       time.Time `json:"expiresAt"`
}

type OrderListResponse struct {
	Orders []Order `json:"orders"`
}

type OrderQueryParam struct {
	Limit              int    `json:"limit"`
	Market             string `json:"market"`
	Status             string `json:"status"`
	Type               string `json:"type"`
	Side               string `json:"side"`
	CreatedBeforeOrAt  string `json:"createdAt"`
	ReturnLatestOrders string `json:"returnLatestOrders"`
}

type WithdrawResponse struct {
	Withdrawal []Withdrawal `json:"withdrawal"`
}

type Withdrawal struct {
	ID              string      `json:"id"`
	Type            string      `json:"type"`
	DebitAsset      string      `json:"debitAsset"`
	CreditAsset     string      `json:"creditAsset"`
	DebitAmount     string      `json:"debitAmount"`
	CreditAmount    string      `json:"creditAmount"`
	TransactionHash string      `json:"transactionHash"`
	Status          string      `json:"status"`
	ClientID        string      `json:"clientId"`
	FromAddress     string      `json:"fromAddress"`
	ToAddress       interface{} `json:"toAddress"`
	ConfirmedAt     interface{} `json:"confirmedAt"`
	CreatedAt       time.Time   `json:"createdAt"`
}

type WithdrawalParam struct {
	ClientID     string `json:"clientId"`
	ToAddress    string `json:"toAddress"`
	CreditAsset  string `json:"creditAsset"`
	CreditAmount string `json:"creditAmount"`

	DebitAmount string `json:"debitAmount"`

	LpPositionId string `json:"lpPositionId"`
	Expiration   string `json:"expiration"`
	Signature    string `json:"signature"`
}

type FillsResponse struct {
	Fills []Fill `json:"fills"`
}

type Fill struct {
	ID        string    `json:"id"`
	Side      string    `json:"side"`
	Liquidity string    `json:"liquidity"`
	Type      string    `json:"type"`
	Market    string    `json:"market"`
	OrderID   string    `json:"orderId"`
	Price     string    `json:"price"`
	Size      string    `json:"size"`
	Fee       string    `json:"fee"`
	CreatedAt time.Time `json:"createdAt"`
}

type FillsParam struct {
	Market            string `json:"market,omitempty"`
	OrderId           string `json:"order_id,omitempty"`
	Limit             string `json:"limit,omitempty"`
	CreatedBeforeOrAt string `json:"createdBeforeOrAt,omitempty"`
}

type FundingPaymentsResponse struct {
	FundingPayments []FundingPayment `json:"fundingPayments`
}

type FundingPayment struct {
	Market       string    `json:"market"`
	Payment      string    `json:"payment"`
	Rate         string    `json:"rate"`
	PositionSize string    `json:"positionSize"`
	Price        string    `json:"price"`
	EffectiveAt  time.Time `json:"effectiveAt"`
}

type FundingPaymentsParam struct {
	Market              string `json:"market,omitempty"`
	Limit               string `json:"limit,omitempty"`
	EffectiveBeforeOrAt string `json:"effectiveBeforeOrAt,omitempty"`
}

type HistoricalPnLResponse struct {
	HistoricalPnLs []HistoricalPnL `json:"historicalPnl"`
}

type HistoricalPnL struct {
	AccountID    string    `json:"accountId"`
	Equity       string    `json:"equity"`
	TotalPnl     string    `json:"totalPnl"`
	NetTransfers string    `json:"netTransfers"`
	CreatedAt    time.Time `json:"createdAt"`
}

type HistoricalPnLParam struct {
	EffectiveBeforeOrAt string `json:"effectiveBeforeOrAt,omitempty"`
	EffectiveAtOrAfter  string `json:"effectiveAtOrAfter,omitempty"`
}

type TradingRewardsResponse struct {
	TradingRewardss []TradingRewards `json:"historicalPnl"`
}

type TradingRewards struct {
	AccountID    string    `json:"accountId"`
	Equity       string    `json:"equity"`
	TotalPnl     string    `json:"totalPnl"`
	NetTransfers string    `json:"netTransfers"`
	CreatedAt    time.Time `json:"createdAt"`
}

type TradingRewardsParam struct {
	EffectiveBeforeOrAt string `json:"effectiveBeforeOrAt,omitempty"`
	EffectiveAtOrAfter  string `json:"effectiveAtOrAfter,omitempty"`
}

func (o OrderQueryParam) ToParams() url.Values {
	params := url.Values{}
	if o.Market != "" {
		params.Add("market", o.Market)
	}
	if o.Status != "" {
		params.Add("status", o.Status)
	}
	if o.Side != "" {
		params.Add("side", o.Side)
	}
	if o.Type != "" {
		params.Add("type", o.Type)
	}
	if o.Limit != 0 {
		params.Add("limit", strconv.Itoa(o.Limit))
	}
	if o.CreatedBeforeOrAt != "" {
		params.Add("createdBeforeOrAt", o.CreatedBeforeOrAt)
	}
	if o.ReturnLatestOrders != "" {
		params.Add("returnLatestOrders", o.ReturnLatestOrders)
	}
	return params
}
