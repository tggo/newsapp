package http

import (
	"net/http"
	"strconv"

	"boostersNews/internal/app/ginx"

	"github.com/gin-gonic/gin"
)

func (h *handler) delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ginx.ResponseError(c, http.StatusBadRequest, err,
			"could not bind id to int",
			"wrong request id",
			h.logger,
		)
		return
	}

	errGet := h.service.Delete(c, int64(id))
	if errGet != nil {
		ginx.ResponseError(c, http.StatusInternalServerError, errGet,
			"could not delete post",
			"something wrong when delete post",
			h.logger,
		)
	}

	ginx.ResponseOK(c)
}
