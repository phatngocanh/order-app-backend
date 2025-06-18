package error_utils

import (
	"net/http"

	httpcommon "github.com/pna/order-app-backend/internal/domain/http_common"
)

func ErrorCodeToHttpResponse(errCode string, field string) (statusCode int, httpErrResponse httpcommon.HttpResponse[any]) {
	switch errCode {
	case ErrorCode.DB_DOWN:
		statusCode = http.StatusInternalServerError
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "Our database system is currently unavailable. Please try again in a few minutes",
			Field:   field,
			Code:    ErrorCode.DB_DOWN,
		})
	case ErrorCode.BAD_REQUEST:
		statusCode = http.StatusBadRequest
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "Invalid request parameters",
			Field:   field,
			Code:    ErrorCode.BAD_REQUEST,
		})
	case ErrorCode.FORBIDDEN:
		statusCode = http.StatusForbidden
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "You do not have permission to perform this action",
			Field:   field,
			Code:    ErrorCode.FORBIDDEN,
		})
	case ErrorCode.ACCESS_TOKEN_INVALID:
		statusCode = http.StatusUnauthorized
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "Invalid access token",
			Field:   field,
			Code:    ErrorCode.ACCESS_TOKEN_INVALID,
		})
	default:
		statusCode = http.StatusInternalServerError
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "An unexpected error occurred. Please try again later or contact support if the problem persists",
			Field:   field,
			Code:    ErrorCode.INTERNAL_SERVER_ERROR,
		})
	}

	return
}
