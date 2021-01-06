package bclient

import (
	"errors"
	"math/big"
	"strings"

	"github.com/bonedaddy/go-indexed/uniswap"
	"github.com/bonedaddy/go-indexed/utils"
	"github.com/ethereum/go-ethereum/common"
)

// NdxEthPairAddress returns the UniswapV2 pair of NDX-ETH
func (c *Client) NdxEthPairAddress() common.Address {
	return uniswap.GeneratePairAddress(NDXTokenAddress, WETHTokenAddress)
}

// Defi5EthPairAddress returns the UniswapV2 pair of DEFI5-ETH
func (c *Client) Defi5EthPairAddress() common.Address {
	return uniswap.GeneratePairAddress(DEFI5TokenAddress, WETHTokenAddress)
}

// Cc10EthPairAddress returns the UniswapV2 pair of CC10-ETH
func (c *Client) Cc10EthPairAddress() common.Address {
	return uniswap.GeneratePairAddress(CC10TokenAddress, WETHTokenAddress)
}

// NdxDaiPrice returns the price of NDX in terms of DAI
func (c *Client) NdxDaiPrice() (float64, error) {
	ndxEthPrice, err := c.ExchangeAmount(utils.ToWei("1.0", 18), "ndx-eth")
	if err != nil {
		return 0, err
	}
	ndxEthPriceDec := utils.ToDecimal(ndxEthPrice, 18)
	ethDaiPrice, err := c.EthDaiPrice()
	if err != nil {
		return 0, err
	}
	ethDaiPriceDec := utils.ToDecimal(utils.ToWei(ethDaiPrice.Int64(), 18), 18)
	// derive the price of DEFI5 by getting the amount of ETH you would get from
	// 1 DEFI5 token, and converting that into DAI
	edF, _ := ethDaiPriceDec.Float64()
	neF, _ := ndxEthPriceDec.Float64()
	return edF * neF, nil
}

// Cc10DaiPrice returns the price of CC10 in terms of DAI
func (c *Client) Cc10DaiPrice() (float64, error) {
	cc10EthPrice, err := c.ExchangeAmount(utils.ToWei("1.0", 18), "cc10-eth")
	if err != nil {
		return 0, err
	}
	cc10EthPriceDec := utils.ToDecimal(cc10EthPrice, 18)
	ethDaiPrice, err := c.EthDaiPrice()
	if err != nil {
		return 0, err
	}
	ethDaiPriceDec := utils.ToDecimal(utils.ToWei(ethDaiPrice.Int64(), 18), 18)
	// derive the price of DEFI5 by getting the amount of ETH you would get from
	// 1 DEFI5 token, and converting that into DAI
	edF, _ := ethDaiPriceDec.Float64()
	ceF, _ := cc10EthPriceDec.Float64()
	return edF * ceF, nil
}

// Defi5DaiPrice returns the price of DEFI5 in terms of DAI
func (c *Client) Defi5DaiPrice() (float64, error) {
	defi5EthPrice, err := c.ExchangeAmount(utils.ToWei("1.0", 18), "defi5-eth")
	if err != nil {
		return 0, err
	}
	defi5EthPriceDec := utils.ToDecimal(defi5EthPrice, 18)
	ethDaiPrice, err := c.EthDaiPrice()
	if err != nil {
		return 0, err
	}
	ethDaiPriceDec := utils.ToDecimal(utils.ToWei(ethDaiPrice.Int64(), 18), 18)
	// derive the price of DEFI5 by getting the amount of ETH you would get from
	// 1 DEFI5 token, and converting that into DAI
	edF, _ := ethDaiPriceDec.Float64()
	deF, _ := defi5EthPriceDec.Float64()
	return edF * deF, nil
}

// EthDaiPrice returns the price of ETH in terms of DAI
func (c *Client) EthDaiPrice() (*big.Int, error) {
	reserves, err := c.Reserves("eth-dai")
	if err != nil {
		return nil, err
	}
	return new(big.Int).Div(reserves.Reserve1, reserves.Reserve0), nil
}

// Reserves returns available reserves in the pair
func (c *Client) Reserves(pair string) (*uniswap.Reserve, error) {
	switch strings.ToLower(pair) {
	case "ndx-eth":
		return c.uc.GetReserves(NDXTokenAddress, WETHTokenAddress)
	case "eth-ndx":
		return c.uc.GetReserves(WETHTokenAddress, NDXTokenAddress)
	case "cc10-eth":
		return c.uc.GetReserves(CC10TokenAddress, WETHTokenAddress)
	case "eth-cc10":
		return c.uc.GetReserves(WETHTokenAddress, CC10TokenAddress)
	case "defi5-eth":
		return c.uc.GetReserves(DEFI5TokenAddress, WETHTokenAddress)
	case "eth-defi5":
		return c.uc.GetReserves(WETHTokenAddress, DEFI5TokenAddress)
	case "eth-dai":
		return c.uc.GetReserves(WETHTokenAddress, DAITokenAddress)
	default:
		return nil, errors.New("unsupported pair")
	}
}

// ExchangeAmount returns the exchange amount for a variety of pairs
func (c *Client) ExchangeAmount(amount *big.Int, pair string) (*big.Int, error) {
	switch strings.ToLower(pair) {
	case "ndx-eth":
		return c.uc.GetExchangeAmount(amount, NDXTokenAddress, WETHTokenAddress)
	case "eth-ndx":
		return c.uc.GetExchangeAmount(amount, WETHTokenAddress, NDXTokenAddress)
	case "cc10-eth":
		return c.uc.GetExchangeAmount(amount, CC10TokenAddress, WETHTokenAddress)
	case "eth-cc10":
		return c.uc.GetExchangeAmount(amount, WETHTokenAddress, CC10TokenAddress)
	case "defi5-eth":
		return c.uc.GetExchangeAmount(amount, DEFI5TokenAddress, WETHTokenAddress)
	case "eth-defi5":
		return c.uc.GetExchangeAmount(amount, WETHTokenAddress, DEFI5TokenAddress)
	default:
		return nil, errors.New("unsupported pair")
	}
}

// PairDecimals returns the decimals for the corresponding token pair
func (c *Client) PairDecimals(pair string) int {
	switch strings.ToLower(pair) {
	case "ndx-eth":
		return 18
	case "eth-ndx":
		return 18
	case "cc10-eth":
		return 18
	case "eth-cc10":
		return 18
	case "defi5-eth":
		return 18
	case "eth-defi5":
		return 18
	default:
		return 0
	}
}