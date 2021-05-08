package indexer

import (
	"blockstime/internal/config"
	"blockstime/internal/engines"
	evmrpcclient "blockstime/internal/engines/evm"
	"blockstime/internal/timeslice"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"log"
)

type indexer struct {
	Network   *config.Network
	RpcClient engines.INetworkEngine
	Blocks    []int64

	closed chan struct{}
	wg     sync.WaitGroup
	ticker *time.Ticker
}

// graceful stop of the indexer - must was for the save to be finished
func (in *indexer) Stop() {
	close(in.closed)
	in.wg.Wait()
}

func New(network *config.Network) (*indexer, error) {
	if len(network.Nodes) == 0 {
		return nil, errors.New("no nodes in the network")
	}
	res := &indexer{
		Network: network,
		closed:  make(chan struct{}),
		ticker:  time.NewTicker(time.Second * 2),
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
	in.wg.Add(1)
	err := timeslice.Save(in.Blocks, in.Network.LocalPath)
	in.wg.Done()
	return err
}

func (in *indexer) Run() error {
	total := len(in.Blocks)
	if total == 0 {
		return errors.New("network error - no blocks")
	}
	missing := 0
	for _, v := range in.Blocks {
		if v == 0 {
			missing++
		}
	}
	log.Printf("[index] Missing %v out of %v blocks (%5.2f%%)\n",
		missing, total, 100.0*float64(missing)/float64(total))

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	go func() {
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
						log.Println("[flush]", in.Network.LocalPath)
					}
				}
			}
		}
	}()

	select {
	case sig := <-c:
		fmt.Printf("Got %s signal. Aborting...\n", sig)
		in.Stop()
	}
	return nil
}
