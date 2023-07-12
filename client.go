package dydx

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-numb/go-dydx/helpers"
	"github.com/go-numb/go-dydx/onboard"
	"github.com/go-numb/go-dydx/private"
	"github.com/go-numb/go-dydx/public"
	"github.com/go-numb/go-dydx/types"
	"github.com/umbracle/ethgo/jsonrpc"
)

type Client struct {
	options types.Options

	Host                      string
	StarkPublicKey            string
	StarkPrivateKey           string
	StarkPublicKeyYCoordinate string
	ApiKeyCredentials         *types.ApiKeyCredentials
	ApiTimeout                time.Duration

	DefaultAddress string
	NetworkId      int
	Web3           *jsonrpc.Client
	EthSigner      helpers.EthSigner

	Private    *private.Private
	Public     *public.Public
	OnBoarding *onboard.OnBoarding

	Logger *log.Logger

	httpClient *http.Client
}

func New(options types.Options) *Client {
	client := &Client{
		Host:              strings.TrimPrefix(options.Host, "/"),
		ApiTimeout:        3 * time.Second,
		DefaultAddress:    options.DefaultEthereumAddress,
		StarkPublicKey:    options.StarkPublicKey,
		StarkPrivateKey:   options.StarkPrivateKey,
		ApiKeyCredentials: options.ApiKeyCredentials,
		NetworkId:         options.NetworkId,

		Logger: log.New(os.Stderr, "go-dydx ", log.LstdFlags),

		httpClient: &http.Client{
			Timeout: time.Second * 5,
		},
	}

	if options.Web3 != nil {
		networkId := options.NetworkId
		if networkId == 0 {
			net, _ := options.Web3.Net().Version()
			networkId = int(net)
		}

		client.Web3 = options.Web3
		client.EthSigner = &helpers.EthWeb3Signer{Web3: options.Web3}
		client.NetworkId = networkId
	}

	if client.NetworkId == 0 {
		client.NetworkId = types.NetworkIdMainnet
	}

	if options.StarkPrivateKey != "" {
		client.StarkPrivateKey = options.StarkPrivateKey
		client.EthSigner = &helpers.EthKeySinger{PrivateKey: options.StarkPrivateKey}
	}

	if options.HttpClient != nil {
		client.httpClient = options.HttpClient
	}

	if options.Logger != nil {
		client.Logger = options.Logger
	}

	client.OnBoarding = &onboard.OnBoarding{
		Host:       client.Host,
		EthSigner:  client.EthSigner,
		NetworkId:  client.NetworkId,
		EthAddress: client.DefaultAddress,
		Singer:     helpers.NewSigner(client.EthSigner, client.NetworkId),
		Logger:     client.Logger,
	}
	if options.ApiKeyCredentials == nil {
		client.ApiKeyCredentials = client.OnBoarding.RecoverDefaultApiCredentials(client.DefaultAddress)
	}

	client.Private = &private.Private{
		Host:              client.Host,
		NetworkId:         client.NetworkId,
		StarkPrivateKey:   client.StarkPrivateKey,
		DefaultAddress:    client.DefaultAddress,
		ApiKeyCredentials: client.ApiKeyCredentials,

		RateLimit:  new(types.RateLimit),
		Logger:     client.Logger,
		HttpClient: client.httpClient,
	}
	client.Public = &public.Public{
		Host:      client.Host,
		NetworkId: client.NetworkId,

		RateLimit:  new(types.RateLimit),
		Logger:     client.Logger,
		HttpClient: client.httpClient,
	}

	return client
}
