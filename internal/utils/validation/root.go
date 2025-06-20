package validation

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	httpcommon "github.com/pna/order-app-backend/internal/domain/http_common"
	"github.com/pna/order-app-backend/internal/utils/error_utils"
	stringutils "github.com/pna/order-app-backend/internal/utils/string_utils"
)

func BindJsonAndValidate(c *gin.Context, dest interface{}) error {
	err := c.ShouldBindJSON(&dest)

	if err != nil {
		checkErr(c, err)
	}
	return err
}

func checkErr(c *gin.Context, err error) {
	switch t := err.(type) {
	case *json.UnmarshalTypeError:
		httpErr := httpcommon.Error{
			Message: "Invalid data type", Code: error_utils.ErrorCode.BAD_REQUEST, Field: t.Field,
		}
		c.JSON(http.StatusBadRequest, httpcommon.NewErrorResponse(httpErr))
		return
	case *json.SyntaxError:
		httpErr := httpcommon.Error{Message: err.Error(), Code: error_utils.ErrorCode.BAD_REQUEST}
		c.JSON(http.StatusBadRequest, httpcommon.NewErrorResponse(httpErr))
		return
	case validator.ValidationErrors:
		httpErrs := handleValidationErrors(err)
		c.JSON(http.StatusBadRequest, httpcommon.NewErrorResponse(httpErrs...))
		return
	default:
		httpErr := httpcommon.Error{Message: err.Error(), Code: error_utils.ErrorCode.BAD_REQUEST, Field: ""}
		c.JSON(http.StatusBadRequest, httpErr)
		return
	}
}

func handleValidationErrors(errs error) (httpErrs []httpcommon.Error) {
	for _, fieldErr := range errs.(validator.ValidationErrors) {
		field := stringutils.FirstLetterToLower(fieldErr.Field())
		httpErrs = append(httpErrs, httpcommon.Error{
			Message: "Invalid request", Code: error_utils.ErrorCode.BAD_REQUEST, Field: field,
		})
	}
	return httpErrs
}
