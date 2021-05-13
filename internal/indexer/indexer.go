package indexer

import (
	"blockstime/internal/config"
	"blockstime/internal/engines"
	evmrpcclient "blockstime/internal/engines/evm"
	"blockstime/internal/timeslice"
	"errors"
	"fmt"
	"math"
	"os"
	"os/signal"
	"sync"
	"time"

	"log"

	"go.uber.org/atomic"
)

type indexer struct {
	Network   *config.Network
	RpcClient engines.INetworkEngine
	Blocks    []int64
	Parallel  int

	wgSave     sync.WaitGroup
	wgRead     sync.WaitGroup
	newCounter atomic.Uint64
	errCounter atomic.Uint64
}

func New(network *config.Network) (*indexer, error) {
	if len(network.Nodes) == 0 {
		return nil, errors.New("no nodes in the network")
	}
	res := &indexer{
		Network:  network,
		Parallel: 600,
	}

	// so far it uses only the first node in the network
	rpcClient, err := evmrpcclient.NewClient(&network.Nodes[0]) // define with (network.Engine)
	if err != nil {
		return nil, err
	}
	res.RpcClient = rpcClient

	maxBlockHeight, err := res.RpcClient.GetHeadBlockNumber()
	if err != nil {
		return nil, err
	}
	res.Blocks = make([]int64, maxBlockHeight+1)
	blocks, err := timeslice.Load(network.LocalPath)
	if err != nil {
		log.Println("[warn]", err.Error())
	} else {
		for bi, bv := range blocks {
			res.Blocks[bi] = bv
		}
	}
	return res, nil
}

func (in *indexer) Save() error {
	in.wgSave.Add(1)
	err := timeslice.Save(in.Blocks, in.Network.LocalPath)
	in.wgSave.Done()
	return err
}

func (in *indexer) ReadBlock(index int64, wg *sync.WaitGroup) {
	// start := time.Now()
	n, err := in.RpcClient.GetBlockTime(int64(index))
	if err != nil {
		// log.Println(err)
		in.errCounter.Inc()
	} else if n != 0 {
		in.Blocks[index] = n
		// log.Debug("[index]", index, time.Unix(n, 0).UTC(), "took", time.Since(start))
		in.newCounter.Inc()
	}
	wg.Done()
}

func (in *indexer) getSyncPct() (int, float64) {
	missing := 0
	for index, v := range in.Blocks {
		if v == 0 && index > 0 {
			missing++
		}
	}
	return missing, 100.0 - 100.0*float64(missing)/float64(len(in.Blocks))
}

func (in *indexer) Stat() {
	total := len(in.Blocks)
	if total > 0 {
		missing, pct := in.getSyncPct()
		log.Printf("[index] Total %v blocks (synced %5.2f%%, missing %d)\n", total, pct, missing)
	}
}

func (in *indexer) Run() error {
	in.Stat()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	ch := make(chan int64, in.Parallel)
	// channel reader
	go func() {
		for {
			select {
			case index := <-ch:
				// fmt.Println("[index] block", index)
				in.wgRead.Add(1)
				go in.ReadBlock(int64(index), &in.wgRead)
			}
		}
	}()

	go func() {
		// flush every minute
		seconds := int64(15)
		period, _ := time.ParseDuration(fmt.Sprintf("%ds", seconds))

		for {
			time.Sleep(period)

			start := time.Now()
			resetCounter := uint64(0)
			current := in.newCounter.Swap(resetCounter)

			resetErrCounter := uint64(0)
			currentErr := in.errCounter.Swap(resetErrCounter)

			speed := math.Round(float64(current) / float64(seconds))
			if err := in.Save(); err != nil {
				log.Fatal("[save] ", err)
			} else {
				if currentErr > 0 {
					log.Println("[warn]", currentErr, "errors, blocks skipped")
				}
				missing, pct := in.getSyncPct()
				log.Println("[flush]", current, "blocks", speed, "/s",
					fmt.Sprintf("(synced %5.2f%%, missing %d)", pct, missing),
					" took ", time.Since(start))
				if missing == 0 {
					log.Println("[saver] done. stopping")
					os.Exit(0)
				}
			}
		}

	}()
	// channel writer - task generator
	go func() {
		for index, v := range in.Blocks {
			if v == 0 && index > 0 {
				in.wgRead.Wait()
				ch <- int64(index)
			}
		}
		log.Println("[index] tasks were distributed. cooling down.")
		time.Sleep(15 * time.Second)
		close(ch) // no more values to be sent to the channel
	}()

	select {
	case sig := <-c:
		fmt.Printf("Got %s signal. Aborting...\n", sig)
		in.Stop()
	}
	return nil
}

// graceful stop of the indexer -
// required to the Save function to be finished correctly
// (and do not loose database on interruption)
func (in *indexer) Stop() {
	in.wgSave.Wait()
}
