package resource

import (
	"github.com/MISingularity/deepdash/pkg/statistics"
	"github.com/gin-gonic/gin"
	"net/http"
)

func TotalPvHandler(c *gin.Context) {
	pv := statistics.QueryTotalPvCount()

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, map[string]int{"total_pv": pv})
}
