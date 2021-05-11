package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type IService interface {
	BlocksFromPeriods(context.Context, TimePeriod) (*BlocksPeriod, error)
	PeriodFromBlocks(context.Context, BlocksPeriod) (*TimePeriod, error)
	StatsDaily(ctx context.Context, network string) (*BlockStatsResponse, error)
	StatsYearly(ctx context.Context, network string) (*BlockStatsResponse, error)
}

type Controller struct {
	svc IService
}

func NewController(svc IService) *Controller {
	return &Controller{svc: svc}
}

// Ping godoc
// @Summary ping
// @Description Healthcheck endpoint
// @ID ping
// @Accept json
// @Produce json
// @Success 200 {object} PingResponse
// @Failure 400,404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Failure default {object} httputil.DefaultError
// @Router /ping [get]
func (c *Controller) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, &PingResponse{
		App: "blockstime",
		Tm:  time.Now(),
	})
}

// @Summary blocks
// @Description Getting blocks range from time period
// @ID blocks
// @Success 200 {object} BlocksPeriod
// @Accept json
// @Produce json
// @Router /blocks [get]
func (c *Controller) GetBlocksFromPeriods(ctx *gin.Context) {
	var input TimePeriod
	if err := ctx.BindQuery(&input); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
	} else if len(input.Network) == 0 {
		ctx.JSON(400, map[string]string{"error": "missing network parameter"})
	} else if !input.IsValid() {
		ctx.JSON(400, map[string]string{
			"error": "either start or end is expected",
		})
	}
	res, err := c.svc.BlocksFromPeriods(ctx, input)
	if err != nil {
		ctx.JSON(500, map[string]string{"error": err.Error()})
	}
	ctx.JSON(200, res)
}

// @Summary periods
// @Description Getting time period from blocks range
// @ID blocks
// @Accept json
// @Produce json
// @Success 200 {object} TimePeriod
// @Failure 400,404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Failure default {object} httputil.DefaultError
// @Router /periods [get]
func (c *Controller) GetPeriodFromBlocks(ctx *gin.Context) {
	var input BlocksPeriod
	if err := ctx.BindQuery(&input); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
	} else if len(input.Network) == 0 {
		ctx.JSON(400, map[string]string{"error": "missing network parameter"})
	} else if !input.IsValid() {
		ctx.JSON(400, map[string]string{
			"error": "either block_start or block_end is expected",
		})
	}
	res, err := c.svc.PeriodFromBlocks(ctx, input)
	if err != nil {
		ctx.JSON(500, map[string]string{"error": err.Error()})
	}
	ctx.JSON(200, res)
}

// @Summary stats daily
// @Description Getting totals of blocks processed daily
// @ID statsdaily
// @Accept json
// @Produce json
// @Success 200 {object} BlockStatsResponse
// @Failure 400,404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Failure default {object} httputil.DefaultError
// @Router /stats/daily [get]
func (c *Controller) GetStatsDaily(ctx *gin.Context) {
	network := ctx.Request.URL.Query().Get("network")
	if len(network) == 0 {
		ctx.JSON(400, map[string]string{"error": "missing network parameter"})
	}

	res, err := c.svc.StatsDaily(ctx, network)
	if err != nil {
		ctx.JSON(500, map[string]string{"error": err.Error()})
	}
	ctx.JSON(200, res)
}

// @Summary stats yearly
// @Description Getting totals of blocks processed yearly
// @ID statsyearly
// @Accept json
// @Produce json
// @Success 200 {object} BlockStatsResponse
// @Failure 400,404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Failure default {object} httputil.DefaultError
// @Router /stats/yearly [get]
func (c *Controller) GetStatsYearly(ctx *gin.Context) {
	network := ctx.Request.URL.Query().Get("network")
	if len(network) == 0 {
		ctx.JSON(400, map[string]string{"error": "missing network parameter"})
	}

	res, err := c.svc.StatsYearly(ctx, network)
	if err != nil {
		ctx.JSON(500, map[string]string{"error": err.Error()})
	}
	ctx.JSON(200, res)
}
