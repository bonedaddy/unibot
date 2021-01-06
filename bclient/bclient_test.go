package bclient

import (
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

var (
	// my personal address
	myAddress = common.HexToAddress("0x5a361A1dfd52538A158e352d21B5b622360a7C13")
)

func doSetup(t *testing.T) *Client {
	infuraAPIKey := os.Getenv("INFURA_API_KEY")
	if infuraAPIKey == "" {
		t.Fatal("INFURA_API_KEY env var is empty")
	}
	client, err := NewInfuraClient(infuraAPIKey, false)
	require.NoError(t, err)
	return client

}

func TestBClient(t *testing.T) {
	client := doSetup(t)
	t.Cleanup(func() {
		client.Close()
	})
	t.Run("Misc", func(t *testing.T) {
		_, err := client.CurrentBlock()
		require.NoError(t, err)
		require.NotNil(t, client.Uniswap())
		require.Equal(t, "Maker", guessTokenName("0x9f8F72aA9304c8B593d555F12eF6589cC3A579A2"))
		require.Equal(t, "MKR", guessTokenSymbol("0x9f8F72aA9304c8B593d555F12eF6589cC3A579A2"))
	})

}
