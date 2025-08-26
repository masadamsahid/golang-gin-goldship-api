package helpers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParsePaginationFromQueryParams(ctx *gin.Context) (int, int, error) {
	strPage := ctx.Query("page")
	strPageSize := ctx.Query("pageSize")

	if strPage == "" {
		strPage = "1"
	}
	if strPageSize == "" {
		strPageSize = "10"
	}

	page, err := strconv.Atoi(strPage)
	if err != nil || page < 1 {
		return 0, 0, errors.New("invalid 'page' query")
	}
	pageSize, err := strconv.Atoi(strPageSize)
	if err != nil || pageSize < 1 {
		return 0, 0, errors.New("invalid 'pageSize' query")
	}

	return page, pageSize, nil
}
