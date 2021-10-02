// Package bitflyer provides bitFlyer HTTP Client for Go.
package bitflyer

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/nasa9084/autotrade/jsonrpc"
)

const (
	DefaultHTTPEndpoint     = "https://api.bitflyer.com"
	DefaultRealtimeEndpoint = "wss://ws.lightstream.bitflyer.com/json-rpc"
)

// Client is a bitFlyer lightning HTTP API client.
type Client struct {
	httpEndpoint     string
	realtimeEndpoint string

	apiKey    string
	apiSecret string

	httpclient *http.Client
}

// Balance represents account asset balance.
type Balance struct {
	CurrencyCode string  `json:"currency_code"`
	Amount       float64 `json:"amount"`
	Available    float64 `json:"available"`
}

type subscribeParams struct {
	Channel string `json:"channel"`
}

type channelMessage struct {
	Channel string          `json:"channel"`
	Message json.RawMessage `json:"message"`
}

type Ticker struct {
	ProductCode     string    `json:"product_code"`
	Timestamp       time.Time `json:"timestamp"`
	State           string    `json:"state"`
	TickID          int64     `json:"tick_id"`
	BestBid         float64   `json:"best_bid"`
	BestAsk         float64   `json:"best_ask"`
	BestBidSize     float64   `json:"best_bid_size"`
	BestAskSize     float64   `json:"best_ask_size"`
	TotalBidDepth   float64   `json:"total_bid_depth"`
	TotalAskDepth   float64   `json:"total_ask_depth"`
	MarketBidSize   float64   `json:"market_bid_size"`
	MarketAskSize   float64   `json:"market_ask_size"`
	LastTradePrice  float64   `json:"ltp"`
	Volume          float64   `json:"volume"`
	VolumeByProduct float64   `json:"volume_by_product"`
}

// New returns a new Client.
func New(apiKey, apiSecret string) *Client {
	return &Client{
		httpEndpoint:     DefaultHTTPEndpoint,
		realtimeEndpoint: DefaultRealtimeEndpoint,

		apiKey:    apiKey,
		apiSecret: apiSecret,

		httpclient: http.DefaultClient,
	}
}

// GetBalance retrieve account asset balances.
func (c *Client) GetBalance(ctx context.Context) ([]Balance, error) {
	resp, err := c.get(ctx, "/v1/me/getbalance")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var balances []Balance
	if err := json.NewDecoder(resp.Body).Decode(&balances); err != nil {
		return nil, err
	}

	return balances, nil
}

// Ticker subscribe ticker channel using bitFlyer realtime API.
// The productCode can be obtained from https://lightning.bitflyer.com/docs?lang=ja#%E3%83%9E%E3%83%BC%E3%82%B1%E3%83%83%E3%83%88%E3%81%AE%E4%B8%80%E8%A6%A7
func (c *Client) Ticker(ctx context.Context, productCode string) (chan Ticker, error) {
	rpcc, err := jsonrpc.New(c.realtimeEndpoint)
	if err != nil {
		return nil, err
	}

	if err := rpcc.Do("subscribe", subscribeParams{Channel: "lightning_ticker_" + productCode}); err != nil {
		return nil, err
	}

	ch := make(chan Ticker, 100)
	go func() {
		for {
			select {
			case <-ctx.Done():
				rpcc.Close()
			default:
			}

			var chmsg channelMessage
			if err := rpcc.Read("channelMessage", &chmsg); err != nil {
				log.Println(err)
				continue
			}

			var t Ticker
			if err := json.Unmarshal(chmsg.Message, &t); err != nil {
				log.Println(err)
				continue
			}

			ch <- t
		}
	}()

	return ch, nil
}

func (c *Client) request(ctx context.Context, method, path string, body *bytes.Buffer) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.httpEndpoint+path, body)
	if err != nil {
		return nil, fmt.Errorf("creating a new HTTP request: %w", err)
	}

	now := strconv.FormatInt(time.Now().Unix(), 10)

	txt := now + method + path
	if body != nil {
		txt += body.String()
	}
	hasher := hmac.New(sha256.New, []byte(c.apiSecret))
	hasher.Write([]byte(txt))
	sign := hex.EncodeToString(hasher.Sum(nil))

	req.Header.Add("ACCESS-KEY", c.apiKey)
	req.Header.Add("ACCESS-TIMESTAMP", now)
	req.Header.Add("ACCESS-SIGN", sign)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.httpclient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http requesting: %w", err)
	}

	return resp, nil
}

func (c *Client) get(ctx context.Context, path string) (*http.Response, error) {
	return c.request(ctx, http.MethodGet, path, &bytes.Buffer{})
}

func (c *Client) post(ctx context.Context, path string, body interface{}) (*http.Response, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return nil, fmt.Errorf("endoding request body to json: %w", err)
	}

	return c.request(ctx, http.MethodPost, path, &buf)
}
