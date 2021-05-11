package server

import "time"

type PingResponse struct {
	App string    `json:"app"`
	Tm  time.Time `json:"time"`
}

type TimePeriod struct {
	Network string `json:"network"`
	Start   *int64 `json:"start,omitempty"`
	End     *int64 `json:"end,omitempty"`
}

func (p TimePeriod) IsValid() bool {
	// one of the borders must exist
	return p.End != nil || p.Start != nil
}

type BlocksPeriod struct {
	Network    string `json:"network"`
	BlockStart *int64 `json:"block_start,omitempty"`
	BlockEnd   *int64 `json:"block_end,omitempty"`
}

func (p BlocksPeriod) IsValid() bool {
	// one of them must exist
	return p.BlockEnd != nil || p.BlockStart != nil
}

type BlockStatsResponse struct {
	Network string           `json:"network"`
	Stats   map[string]int64 `json:"stats"`
}
