package bclient

import (
	"github.com/ethereum/go-ethereum/common"
)

var (

	// preset tokens

	// DEFI5TokenAddress is the address of the DEFI5 token/pool contract
	DEFI5TokenAddress = common.HexToAddress("0xfa6de2697d59e88ed7fc4dfe5a33dac43565ea41")
	// CC10TokenAddress is the address of the CC10 token/pool contract
	CC10TokenAddress = common.HexToAddress("0x17ac188e09a7890a1844e5e65471fe8b0ccfadf3")
	// WETHTokenAddress is the address of the WETH token contract
	WETHTokenAddress = common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
	// DAITokenAddress is the address of the MCD (Multi Collateral DAI) contract
	DAITokenAddress = common.HexToAddress("0x6b175474e89094c44da98b954eedeac495271d0f")
	// NDXTokenAddress is the address of the NDX contract
	NDXTokenAddress = common.HexToAddress("0x86772b1409b61c639eaac9ba0acfbb6e238e5f83")

	// misc variables

	// InfuraWSURL is the URL for INFURA websockets access
	InfuraWSURL = "wss://mainnet.infura.io/ws/v3/"
	// InfuraHTTPURL is the URL for INFURA HTTP access
	InfuraHTTPURL = "https://mainnet.infura.io/v3/"
)
