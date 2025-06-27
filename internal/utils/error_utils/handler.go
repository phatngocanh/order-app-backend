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
	case ErrorCode.INVENTORY_VERSION_MISMATCH:
		statusCode = http.StatusBadRequest
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "Version mismatch, please try again",
			Field:   field,
			Code:    ErrorCode.INVENTORY_VERSION_MISMATCH,
		})
	case ErrorCode.INVENTORY_QUANTITY_NEGATIVE:
		statusCode = http.StatusBadRequest
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "Quantity cannot be negative",
			Field:   field,
			Code:    ErrorCode.INVENTORY_QUANTITY_NEGATIVE,
		})
	case ErrorCode.INVENTORY_QUANTITY_EXCEEDED:
		statusCode = http.StatusBadRequest
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "Order quantity exceeds available inventory",
			Field:   field,
			Code:    ErrorCode.INVENTORY_QUANTITY_EXCEEDED,
		})
	case ErrorCode.DUPLICATE_ORDER_ITEMS:
		statusCode = http.StatusBadRequest
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "Each product can have at most one order item from inventory and one from external",
			Field:   field,
			Code:    ErrorCode.DUPLICATE_ORDER_ITEMS,
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
	case ErrorCode.USERNAME_NOT_FOUND:
		statusCode = http.StatusNotFound
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "Username not found",
			Field:   field,
			Code:    ErrorCode.USERNAME_NOT_FOUND,
		})
	case ErrorCode.NOT_FOUND:
		statusCode = http.StatusNotFound
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "Resource not found",
			Field:   field,
			Code:    ErrorCode.NOT_FOUND,
		})
	case ErrorCode.UNAUTHORIZED:
		statusCode = http.StatusUnauthorized
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "Unauthorized",
			Field:   field,
			Code:    ErrorCode.UNAUTHORIZED,
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
