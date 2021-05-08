package evmrpcclient

import (
	"blockstime/internal/engines"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRpcClient(t *testing.T) {
	testNode, err := engines.NewTestNode()
	assert.NoError(t, err)

	client, err := NewClient(testNode)
	assert.NoError(t, err)

	block, err := client.GetHeadBlockNumber()
	assert.NoError(t, err)
	syncing, err := client.IsSyncing()
	assert.NoError(t, err)
	blockTime, err := client.GetBlockTime(block)
	assert.NoError(t, err)
	assert.NotEqual(t, blockTime, int64(0))

	fmt.Println("HEAD", block, "SYNCING", syncing,
		"TIME", time.Unix(blockTime, 0).UTC(), "vs", time.Now().UTC())
}
