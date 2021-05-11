package server

import (
	"blockstime/internal/config"
	"context"
	"errors"
)

type service struct {
}

func NewService(cfg config.Config) *service {
	return &service{}
}

func (s *service) BlocksFromPeriods(ctx context.Context, rq TimePeriod) (*BlocksPeriod, error) {
	return nil, errors.New("Not implemented")
}

func (s *service) PeriodFromBlocks(ctx context.Context, rq BlocksPeriod) (*TimePeriod, error) {
	return nil, errors.New("Not implemented")
}

func (s *service) StatsDaily(ctx context.Context, network string) (*BlockStatsResponse, error) {
	return nil, errors.New("Not implemented")
}

func (s *service) StatsYearly(ctx context.Context, network string) (*BlockStatsResponse, error) {
	return nil, errors.New("Not implemented")
}
