// Package jsonrpc provides JSON-RPC 2.0 over WebSocket client.
package jsonrpc

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

const jsonrpcVersion = "2.0"

type Client struct {
	conn *websocket.Conn
}

type outgoingRequest struct {
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type incomingRequest struct {
	Version string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

// New returns a new Client which connects given endpoint using WebSocket.
// Calling (*Client.Close) is caller's responsibility.
func New(endpoint string) (*Client, error) {
	conn, _, err := websocket.DefaultDialer.Dial(endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	return &Client{
		conn: conn,
	}, nil
}

// Close do close WebSocket connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

// Do request.
func (c *Client) Do(method string, params interface{}) error {
	err := c.conn.WriteJSON(outgoingRequest{
		Version: jsonrpcVersion,
		Method:  method,
		Params:  params,
	})

	return err
}

// Read incoming request, verify it, then parse its params.
func (c *Client) Read(method string, params interface{}) error {
	var req incomingRequest
	if err := c.conn.ReadJSON(&req); err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	if req.Version != jsonrpcVersion {
		return fmt.Errorf("jsonrpc version mismatch: %s != %s", req.Version, jsonrpcVersion)
	}

	if req.Method != method {
		return fmt.Errorf("method mismatch: %s != %s", req.Method, method)
	}

	if err := json.Unmarshal(req.Params, params); err != nil {
		return fmt.Errorf("unmarshaling result: %w", err)
	}

	return nil
}
