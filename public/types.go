package public

import (
	"encoding/json"
	"log"
	"time"

	"github.com/huangjosh/go-dydx/types"
)

type Public struct {
	Host      string
	NetworkId int

	RateLimit *types.RateLimit
	Logger    *log.Logger
}

type MarketsResponse struct {
	Markets map[string]Market `json:"markets"`
}

type Market struct {
	Market                           string    `json:"market"`
	BaseAsset                        string    `json:"baseAsset"`
	QuoteAsset                       string    `json:"quoteAsset"`
	StepSize                         string    `json:"stepSize"`
	TickSize                         string    `json:"tickSize"`
	IndexPrice                       string    `json:"indexPrice"`
	OraclePrice                      string    `json:"oraclePrice"`
	PriceChange24H                   string    `json:"priceChange24H"`
	NextFundingRate                  string    `json:"nextFundingRate"`
	MinOrderSize                     string    `json:"minOrderSize"`
	Type                             string    `json:"type"`
	InitialMarginFraction            string    `json:"initialMarginFraction"`
	MaintenanceMarginFraction        string    `json:"maintenanceMarginFraction"`
	BaselinePositionSize             string    `json:"baselinePositionSize"`
	IncrementalPositionSize          string    `json:"incrementalPositionSize"`
	IncrementalInitialMarginFraction string    `json:"incrementalInitialMarginFraction"`
	Volume24H                        string    `json:"volume24H"`
	Trades24H                        string    `json:"trades24H"`
	OpenInterest                     string    `json:"openInterest"`
	MaxPositionSize                  string    `json:"maxPositionSize"`
	AssetResolution                  string    `json:"assetResolution"`
	SyntheticAssetID                 string    `json:"syntheticAssetId"`
	Status                           string    `json:"status"`
	NextFundingAt                    time.Time `json:"nextFundingAt"`
}

type TradesResponse struct {
	Trades []Trade `json:"trades"`
}

type Trade struct {
	Side      string    `json:"side"`
	Size      string    `json:"size"`
	Price     string    `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
}

type TradesParam struct {
	MarketID           string `url:"-"`
	Limit              int    `url:"limit,omitempty"`
	StartingBeforeOrAt string `url:"startingBeforeOrAt,omitempty"`
}

type OrderbookResponse struct {
	Offset string `json:"offset"`
	Bids   []Book `json:"bids"`
	Asks   []Book `json:"asks"`
}

type Book struct {
	Price  string
	Size   string
	Offset string
}

func (p *Book) UnmarshalJSON(data []byte) error {
	var s []string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	l := len(s)
	switch l {
	case 2:
		p.Price = s[0]
		p.Size = s[1]
	case 3:
		p.Price = s[0]
		p.Size = s[1]
		p.Offset = s[2]
	}

	return nil
}

type CandlesResponse struct {
	Candles []Candle `json:"candles"`
}

type Candle struct {
	Market               string    `json:"market"`
	Resolution           string    `json:"resolution"`
	Low                  string    `json:"low"`
	High                 string    `json:"high"`
	Open                 string    `json:"open"`
	Close                string    `json:"close"`
	BaseTokenVolume      string    `json:"baseTokenVolume"`
	Trades               string    `json:"trades"`
	UsdVolume            string    `json:"usdVolume"`
	StartingOpenInterest string    `json:"startingOpenInterest"`
	StartedAt            time.Time `json:"startedAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
}

type CandlesParam struct {
	Market     string `url:"-"`
	Resolution string `url:"resolution,omitempty"`
	FromISO    string `url:"fromISO,omitempty"`
	ToISO      string `url:"toISO,omitempty"`
	// Max:100
	Limit int `url:"limit,omitempty"`
}

type HistoricalFundingsResponse struct {
	HistoricalFundings []HistoricalFunding `json:"historicalFunding"`
}

type HistoricalFunding struct {
	Market      string    `json:"-"`
	Rate        string    `json:"rate"`
	Price       string    `json:"price"`
	EffectiveAt time.Time `json:"effectiveAt"`
}

type HistoricalFundingsParam struct {
	Market              string `url:"-"`
	EffectiveBeforeOrAt string `url:"effectiveBeforeOrAt,omitempty"`
}
