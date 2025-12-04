package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Response[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data"`
}

// It fetches a query parameter from the URL, converts it into an int64,
// and returns a default value if the parameter is missing.
func GetQueryInt64(ctx gin.Context, key string, defaultValue int64) (int64, error) {
	valueStr := ctx.Query(key)
	if valueStr == "" {
		return defaultValue, nil
	}
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}
