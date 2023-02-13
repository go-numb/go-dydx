package realtime

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/buger/jsonparser"
	"github.com/huangjosh/go-dydx/private"
	"github.com/huangjosh/go-dydx/public"

	"github.com/gorilla/websocket"
)

const (
	PRODUCTION = "wss://api.dydx.exchange/v3/ws"
	STAGING    = "wss://api.stage.dydx.exchange/v3/ws"
)

const (
	UNDEFINED = "undefined"
	ERROR     = "error"
	ACCOUNT   = "v3_accounts"
	MARKETS   = "v3_markets"
	ORDERBOOK = "v3_orderbook"
	TRADES    = "v3_trades"
)

type request struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`

	// Public
	ID             string `json:"id,omitempty"`
	IncludeOffsets bool   `json:"includeOffsets,omitempty"`

	// Private
	AccountNumber string `json:"accountNumber,omitempty"`
	ApiKey        string `json:"apiKey,omitempty"`
	Signature     string `json:"signature,omitempty"`
	Timestamp     string `json:"timestamp,omitempty"`
	Passphrase    string `json:"passphrase,omitempty"`
}

type Response struct {
	Type         string `json:"type"`
	Channel      string `json:"channel"`
	ConnectionID string `json:"connection_id,omitempty"`
	ID           string `json:"id,omitempty"`
	MessageID    int    `json:"message_id,omitempty"`
	Contents     any    `json:"-"`

	Account   Account                  `json:"-"`
	Markets   public.MarketsResponse   `json:"-"`
	Trades    public.TradesResponse    `json:"-"`
	Orderbook public.OrderbookResponse `json:"-"`

	Results error
}

type Account struct {
	Orders []private.Order `json:"orders"`
	private.AccountResponse
	private.TransfersResponse
	private.FundingPaymentsResponse
	private.FillsResponse
	private.PositionResponse
}

func subscribe(conn *websocket.Conn, private *private.Private, channels, symbols []string) error {
	for i := range channels {
		if private != nil && channels[i] == ACCOUNT {
			r := &request{
				Type:    "subscribe",
				Channel: channels[i],
				// single account
				AccountNumber: "0",
				ApiKey:        private.ApiKeyCredentials.Key,
				Passphrase:    private.ApiKeyCredentials.Passphrase,
			}

			isoTimestamp := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
			r.Signature = private.Sign("/ws/accounts", http.MethodGet, isoTimestamp, "")
			r.Timestamp = isoTimestamp
			if err := conn.WriteJSON(r); err != nil {
				return err
			}
			time.Sleep(time.Second)
			fmt.Println("account request")
			continue
		}

		for j := range symbols {
			r := &request{ // "{"type":"subscribe","channel":"v3_trades","id":"BTC-USD"}
				Type:    "subscribe",
				Channel: channels[i],
				ID:      symbols[j],
			}

			if r.Channel == ORDERBOOK {
				r.IncludeOffsets = true
			}

			if err := conn.WriteJSON(r); err != nil {
				return err
			}
			time.Sleep(time.Second)
		}
	}

	return nil
}

func unsubscribe(conn *websocket.Conn, channels, symbols []string) error {
	for i := range channels {
		if err := conn.WriteJSON(&request{
			Type:    "unsubscribe",
			Channel: channels[i],
		}); err != nil {
			return err
		}
	}
	return nil
}

func ping(conn *websocket.Conn) (err error) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := conn.WriteMessage(websocket.PingMessage, []byte(`pong`)); err != nil {
				goto EXIT
			}
		}
	}
EXIT:
	return err
}

func Connect(ctx context.Context, ch chan Response, channels, symbols []string, private *private.Private, l *log.Logger) error {
RESTART:

	if l == nil {
		l = log.New(os.Stdout, "dYdX websocket ", log.Llongfile)
	}

	conn, _, err := websocket.DefaultDialer.Dial(PRODUCTION, nil)
	if err != nil {
		return err
	}

	// Initial connect
	_, msg, err := conn.ReadMessage()
	if err != nil {
		l.Printf("[ERROR]: msg error: %+v", err)
		return err
	}
	if err := subscribe(conn, private, channels, symbols); err != nil {
		l.Printf("[ERROR]: connected error: %+v", string(msg))
		return err
	}

	// Dorument: The server will send pings every 30s and expects a pong within 10s. The server does not expect pings, but will respond with a pong if sent one.
	// Pong every 10 seconds without Ping.
	go ping(conn)

	func() {
		defer conn.Close()
		defer unsubscribe(conn, channels, symbols)

		for {
			var res Response
			_, msg, err := conn.ReadMessage()
			if err != nil {
				l.Printf("[ERROR]: msg error: %+v", err)
				res.Channel = ERROR
				res.Results = fmt.Errorf("%v", err)
				ch <- res
			}

			if err := json.Unmarshal(msg, &res); err != nil {
				l.Printf("[ERROR]: unmarshal error: %+v", string(msg))
				res.Channel = ERROR
				res.Results = fmt.Errorf("%v", string(msg))
				ch <- res
			}

			data, _, _, err := jsonparser.Get(msg, "contents")
			if err != nil {
				err = fmt.Errorf("[ERROR]: data err: %v %s", err, string(msg))
				l.Println(err)
				res.Channel = ERROR
				res.Results = err
				ch <- res
			}

			switch res.Channel {
			case ACCOUNT:
				if err := json.Unmarshal(data, &res.Account); err != nil {
					l.Printf("[WARN]: cant unmarshal accounts %+v", err)
					continue
				}

			case MARKETS:
				if err := json.Unmarshal([]byte(data), &res.Markets); err != nil {
					l.Printf("[WARN]: cant unmarshal markets %+v", err)
					continue
				}
				// handle case where market data response if keys differently in json payload after inital connection
				if len(res.Markets.Markets) == 0 {
					var marketData map[string]public.Market
					if err := json.Unmarshal([]byte(data), &marketData); err != nil {
						l.Printf("[WARN]: cant unmarshal markets %+v", err)
						continue
					}
					res.Markets.Markets = marketData
				}

			case TRADES:
				if err := json.Unmarshal(data, &res.Trades); err != nil {
					l.Printf("[WARN]: cant unmarshal trades %+v", err)
					continue
				}

			case ORDERBOOK:
				if err := json.Unmarshal(data, &res.Orderbook); err != nil {
					l.Printf("[WARN]: cant unmarshal orderbook %+v", err)
					continue
				}

			case ERROR:
				message, _, _, err := jsonparser.Get(msg, "message")
				err = fmt.Errorf("[ERROR]: type err: %v %s", err, string(message))
				l.Println(err)
				res.Channel = ERROR
				res.Results = err

			default:
				res.Channel = UNDEFINED
				res.Results = fmt.Errorf("%v", string(msg))
			}

			ch <- res

		}
	}()

	goto RESTART

}
