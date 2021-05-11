package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
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
	} else if !input.IsValid() {
		ctx.JSON(400, map[string]string{
			"error": "either start or end is expected",
		})
	}
	ctx.JSON(501, "not implemented")
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
	} else if !input.IsValid() {
		ctx.JSON(400, map[string]string{
			"error": "either block_start or block_end is expected",
		})
	}
	ctx.JSON(501, "not implemented")
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
	ctx.JSON(501, "not implemented")
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
	ctx.JSON(501, "not implemented")
}
