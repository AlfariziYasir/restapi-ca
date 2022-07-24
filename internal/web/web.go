package web

import (
	"restapi/internal/constant"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUrlQueryString(c *gin.Context, key string) string {
	return c.Param(key)
}

func GetUrlQueryInt(key string) (int, error) {
	i, err := strconv.Atoi(key)
	if err != nil {
		return 0, constant.ErrUrlQueryParameter
	}
	return i, nil
}

func GetUrlQueryInt64(c *gin.Context, key string) (int, error) {
	i, err := strconv.Atoi(c.Param(key))
	if err != nil {
		return 0, constant.ErrUrlQueryParameter
	}
	return i, nil
}
