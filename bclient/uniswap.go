package bclient

import (
	"math/big"

	"github.com/bonedaddy/unibot/uniswap"
	"github.com/ethereum/go-ethereum/common"
)

// EthDaiPrice returns the price of ETH in terms of DAI
func (c *Client) EthDaiPrice() (*big.Int, error) {
	reserves, err := c.Reserves(WETHTokenAddress.String(), DAITokenAddress.String())
	if err != nil {
		return nil, err
	}
	return new(big.Int).Div(reserves.Reserve1, reserves.Reserve0), nil
}

// Reserves returns available reserves in the pair
func (c *Client) Reserves(token0, token1 string) (*uniswap.Reserve, error) {
	return c.uc.GetReserves(common.HexToAddress(token0), common.HexToAddress(token1))
}

// ExchangeAmount returns the exchange amount for a variety of pairs
func (c *Client) ExchangeAmount(amount *big.Int, token0, token1 string) (*big.Int, error) {
	return c.uc.GetExchangeAmount(amount, common.HexToAddress(token0), common.HexToAddress(token1))
}
