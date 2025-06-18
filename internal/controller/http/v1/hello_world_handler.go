package v1

import (
	httpcommon "github.com/pna/order-app-backend/internal/domain/http_common"
	"github.com/pna/order-app-backend/internal/domain/model"
	"github.com/pna/order-app-backend/internal/service"
	"github.com/pna/order-app-backend/internal/utils/error_utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HelloWorldHandler struct {
	helloWorldService service.HelloWorldService
}

func NewHelloWorldHandler(helloWorldService service.HelloWorldService) *HelloWorldHandler {
	return &HelloWorldHandler{
		helloWorldService: helloWorldService,
	}
}

// @Summary Hello World
// @Description Hello World
// @Tags Healths
// @Accept json
// @Produce json
// @Router /hello-world [get]
func (h *HelloWorldHandler) HelloWorld(ctx *gin.Context) {
	helloWorldResponse, errCode := h.helloWorldService.HelloWorld(ctx)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}
	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse[model.HelloWorldResponse](helloWorldResponse))
}
