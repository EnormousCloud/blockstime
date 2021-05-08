package indexer

import (
	"blockstime/internal/config"
	"blockstime/internal/engines"
	evmrpcclient "blockstime/internal/engines/evm"
	"blockstime/internal/timeslice"
	"errors"
	"time"

	"log"
)

type indexer struct {
	Network   *config.Network
	RpcClient engines.INetworkEngine
	Blocks    []int64
}

func New(network *config.Network) (*indexer, error) {
	if len(network.Nodes) == 0 {
		return nil, errors.New("no nodes in the network")
	}
	res := &indexer{
		Network: network,
	}
	// so far it uses only the first node in the network
	rpcClient, err := evmrpcclient.NewClient(&network.Nodes[0]) // define with (network.Engine)
	res.RpcClient = rpcClient

	blocks, err := timeslice.Load(network.LocalPath)
	if err != nil {
		log.Println("[WARN]", err.Error())
		maxBlockHeight, err := res.RpcClient.GetHeadBlockNumber()
		if err != nil {
			return nil, err
		}
		res.Blocks = make([]int64, maxBlockHeight)
	} else {
		res.Blocks = blocks
	}
	return res, nil
}

func (in *indexer) Save() error {
	return timeslice.Save(in.Blocks, in.Network.LocalPath)
}

func (in *indexer) Run() error {
	total := len(in.Blocks)
	missing := 0
	for _, v := range in.Blocks {
		if v == 0 {
			missing++
		}
	}
	log.Printf("[index] Missing %v out of %v blocks\n", missing, total)
	newcount := 0
	for index, v := range in.Blocks {
		if v == 0 {
			n, err := in.RpcClient.GetBlockTime(int64(index))
			if err != nil {
				log.Println(err)
			}
			newcount++
			if n != 0 {
				in.Blocks[index] = n
				log.Println("[index]", index, time.Unix(n, 0).UTC())
			}
			if newcount%10 == 0 {
				if err := in.Save(); err != nil {
					log.Println("[save]", err)
				} else {
					log.Println("[flushed]", in.Network.LocalPath)
				}
			}
		}
	}
	return nil
}
