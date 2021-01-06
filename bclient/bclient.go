package bclient

import (
	"context"
	"fmt"
	"time"

	"github.com/bonedaddy/unibot/uniswap"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client wraps ethclient and provides helper functions for interacting with the indexed finance smart contracts
type Client struct {
	ec *ethclient.Client
	uc *uniswap.Client
}

// NewInfuraClient returns an eth client connected to infura
func NewInfuraClient(token string, websockets bool) (*Client, error) {
	var url string
	if websockets {
		url = InfuraWSURL + token
	} else {
		url = InfuraHTTPURL + token
	}
	return NewClient(url)
}

// NewClient returns an eth client connected to an RPC
func NewClient(url string) (*Client, error) {
	ec, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	return &Client{ec, uniswap.NewClient(ec)}, nil
}

// CurrentBlock returns the current block known by the ethereum client
func (c *Client) CurrentBlock() (uint64, error) {
	return c.ec.BlockNumber(context.Background())
}

// Uniswap returns a uniswap client helper
func (c *Client) Uniswap() *uniswap.Client { return c.uc }

// Close terminates the blockchain connection
func (c *Client) Close() {
	c.ec.Close()
}

func guessTokenName(address string) string {
	if address == "0x9f8F72aA9304c8B593d555F12eF6589cC3A579A2" {
		return "Maker"
	}
	return fmt.Sprintf("unknown-%v-%v", time.Now().UnixNano(), len(address))
}

func guessTokenSymbol(address string) string {
	if address == "0x9f8F72aA9304c8B593d555F12eF6589cC3A579A2" {
		return "MKR"
	}
	return fmt.Sprintf("unknown-%v-%v", time.Now().UnixNano(), len(address))
}
