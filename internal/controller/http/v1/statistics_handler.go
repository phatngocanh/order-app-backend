package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpcommon "github.com/pna/order-app-backend/internal/domain/http_common"
	"github.com/pna/order-app-backend/internal/service"
	"github.com/pna/order-app-backend/internal/utils/error_utils"
	log "github.com/sirupsen/logrus"
)

type StatisticsHandler struct {
	statisticsService service.StatisticsService
}

func NewStatisticsHandler(statisticsService service.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{
		statisticsService: statisticsService,
	}
}

// @Summary Get dashboard statistics
// @Description Get all statistics needed for the dashboard
// @Tags Statistics
// @Accept json
// @Produce json
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 200 {object} httpcommon.HttpResponse[model.DashboardStatsResponse]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /statistics/dashboard [get]
func (h *StatisticsHandler) GetDashboardStats(ctx *gin.Context) {
	context := ctx.Request.Context()

	stats, err := h.statisticsService.GetDashboardStats(context)
	if err != "" {
		log.Error("StatisticsHandler.GetDashboardStats Error: " + err)
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(err, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(&stats))
}
