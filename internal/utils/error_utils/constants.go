package error_utils

import (
	"database/sql"

	"github.com/go-redis/redis/v8"
)

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
	FORBIDDEN             string
	INTERNAL_SERVER_ERROR string
	BAD_REQUEST           string
	ACCESS_TOKEN_INVALID  string
}

var ErrorCode = errorCode{
	DB_DOWN:               "DB_DOWN",
	FORBIDDEN:             "FORBIDDEN",
	BAD_REQUEST:           "BAD_REQUEST",
	INTERNAL_SERVER_ERROR: "INTERNAL_SERVER_ERROR",
	ACCESS_TOKEN_INVALID:  "ACCESS_TOKEN_INVALID",
}
