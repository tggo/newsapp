package http

import (
	"net/http"
	"strconv"

	"boostersNews/api/openapi"
	"boostersNews/internal/app/ginx"
	"boostersNews/internal/news/model"

	"github.com/gin-gonic/gin"
)

func (h *handler) update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ginx.ResponseError(c, http.StatusBadRequest, err,
			"could not bind id to int",
			"wrong request id",
			h.logger,
		)
		return
	}

	requestData := openapi.PostUpdate{}
	err = ginx.Bind(c, &requestData)
	if err != nil {
		ginx.ResponseError(c, http.StatusBadRequest, err,
			"could not bind data for update post",
			"wrong request on update post",
			h.logger,
		)
		return
	}
	errCreate := h.service.Update(c, int64(id), &model.Post{
		ID:    int64(id),
		Title: requestData.Title,
		Body:  requestData.Content,
	})
	if errCreate != nil {
		ginx.ResponseError(c, http.StatusBadRequest, errCreate,
			"could not create post",
			"something wrong on creating process",
			h.logger,
		)
	}

	ginx.ResponseOK(c)
}
