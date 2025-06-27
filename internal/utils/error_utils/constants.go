package error_utils

import (
	"database/sql"

	"github.com/go-redis/redis/v8"
)

// ConstraintViolationError represents a database constraint violation
type ConstraintViolationError struct {
	Message string
}

func (e *ConstraintViolationError) Error() string {
	return e.Message
}

// VersionMismatchError represents an optimistic locking failure
type VersionMismatchError struct {
	Message string
}

func (e *VersionMismatchError) Error() string {
	return e.Message
}

type systemErrorMessage struct {
	SqlxNoRow string
	RedisNil  string
}

var SystemErrorMessage = systemErrorMessage{
	SqlxNoRow: sql.ErrNoRows.Error(),
	RedisNil:  redis.Nil.Error(),
}

type errorCode struct {
	// db related
	DB_DOWN string

	// auth related
	FORBIDDEN                   string
	INTERNAL_SERVER_ERROR       string
	BAD_REQUEST                 string
	ACCESS_TOKEN_INVALID        string
	USERNAME_NOT_FOUND          string
	UNAUTHORIZED                string
	INVENTORY_VERSION_MISMATCH  string
	INVENTORY_QUANTITY_NEGATIVE string
	INVENTORY_QUANTITY_EXCEEDED string
	DUPLICATE_ORDER_ITEMS       string

	// generic
	NOT_FOUND string
}

var ErrorCode = errorCode{
	DB_DOWN:                     "DB_DOWN",
	FORBIDDEN:                   "FORBIDDEN",
	BAD_REQUEST:                 "BAD_REQUEST",
	INTERNAL_SERVER_ERROR:       "INTERNAL_SERVER_ERROR",
	ACCESS_TOKEN_INVALID:        "ACCESS_TOKEN_INVALID",
	USERNAME_NOT_FOUND:          "USER_NOT_FOUND",
	UNAUTHORIZED:                "UNAUTHORIZED",
	NOT_FOUND:                   "NOT_FOUND",
	INVENTORY_VERSION_MISMATCH:  "INVENTORY_VERSION_MISMATCH",
	INVENTORY_QUANTITY_NEGATIVE: "INVENTORY_QUANTITY_NEGATIVE",
	INVENTORY_QUANTITY_EXCEEDED: "INVENTORY_QUANTITY_EXCEEDED",
	DUPLICATE_ORDER_ITEMS:       "DUPLICATE_ORDER_ITEMS",
}
