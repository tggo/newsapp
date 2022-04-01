package ginx

import (
	"encoding/json"
	"net/http"
	"strconv"

	openapiClient "boostersNews/api/openapi"
	"boostersNews/internal/app/errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	prefix           = "gin"
	LoggerReqBodyKey = prefix + "/logger-req-body"
)

func Bind(ctx *gin.Context, filter interface{}) error {
	if ctx.ContentType() == "application/json" {
		if err := ctx.BindJSON(filter); err != nil {
			return err
		}
	} else if err := ctx.BindQuery(filter); err != nil {
		return err
	}
	return nil
}

func ResponseOK(c *gin.Context, data ...interface{}) {
	if len(data) > 0 {
		v := make(map[string]interface{})
		for i := range data {
			v[strconv.Itoa(i)] = data[i]
		}
		ResponseSuccess(c, openapiClient.SuccessResponse{Status: errors.StatusOK, Data: v})
	} else {
		ResponseSuccess(c, openapiClient.SuccessResponse{Status: errors.StatusOK})
	}
}

func ResponseError(c *gin.Context, status int, logErr error, logComment, msgToUserDisplay string, logger *zap.Logger) {
	logger.Error(logComment,
		zap.Error(logErr),
		zap.Any("header", c.Request.Header),
		zap.Any("request_params", c.Params))

	ResponseJSON(c, status, openapiClient.BadResponse{
		Status:  errors.StatusError,
		Message: msgToUserDisplay,
	})
}

func ResponseSuccess(c *gin.Context, v interface{}) {
	ResponseJSON(c, http.StatusOK, v)
}

func ResponseJSON(c *gin.Context, status int, v interface{}) {
	buf, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	// c.Set(ResBodyKey, buf)
	c.Data(status, "application/json; charset=utf-8", buf)
	c.Abort()
}
