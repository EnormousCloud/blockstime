package server

import (
	"blockstime/internal/config"
	"blockstime/internal/timeslice"
	"context"
	"time"
)

type networkService struct {
	Network *config.Network
	Blocks  []int64
}

type service struct {
	Networks map[string]networkService
}

func NewService(cfg config.Config) *service {
	networks := map[string]networkService{}
	for _, c := range cfg.Networks {
		if c.Disabled {
			continue
		}
		blocks, _ := timeslice.Load(c.LocalPath)
		networks[c.Name] = networkService{
			Network: &c,
			Blocks:  blocks,
		}
	}
	return &service{Networks: networks}
}

func (s *service) HasNetwork(network string) bool {
	_, ok := s.Networks[network]
	return ok
}

func (s *service) BlocksFromPeriods(ctx context.Context, rq TimePeriod) (*BlocksPeriod, error) {
	return &BlocksPeriod{
		Network:    rq.Network,
		BlockStart: timeslice.BlockBefore(s.Networks[rq.Network].Blocks, rq.Start),
		BlockEnd:   timeslice.BlockAfter(s.Networks[rq.Network].Blocks, rq.End),
	}, nil
}

func (s *service) PeriodFromBlocks(ctx context.Context, rq BlocksPeriod) (*TimePeriod, error) {
	return &TimePeriod{
		Network: rq.Network,
		Start:   timeslice.TimeBefore(s.Networks[rq.Network].Blocks, rq.BlockStart),
		End:     timeslice.TimeAfter(s.Networks[rq.Network].Blocks, rq.BlockEnd),
	}, nil

}

func (s *service) StatsDaily(ctx context.Context, network string) (*BlockStatsResponse, error) {
	res := &BlockStatsResponse{
		Network: network,
		Stats:   map[string]int64{},
	}
	for _, tm := range s.Networks[network].Blocks {
		if tm > 0 {
			dt := time.Unix(tm, 0).UTC().Format("2006-01-02")
			if _, ok := res.Stats[dt]; !ok {
				res.Stats[dt] = 0
			}
			res.Stats[dt]++
		}
	}
	return res, nil
}

func (s *service) StatsYearly(ctx context.Context, network string) (*BlockStatsResponse, error) {
	res := &BlockStatsResponse{
		Network: network,
		Stats:   map[string]int64{},
	}
	for _, tm := range s.Networks[network].Blocks {
		if tm > 0 {
			dt := time.Unix(tm, 0).UTC().Format("2006")
			if _, ok := res.Stats[dt]; !ok {
				res.Stats[dt] = 0
			}
			res.Stats[dt]++
		}
	}
	return res, nil
}
