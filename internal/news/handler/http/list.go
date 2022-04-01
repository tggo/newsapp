package http

import (
	"net/http"

	"boostersNews/api/openapi"
	"boostersNews/internal/app/ginx"
	"boostersNews/internal/news/model"

	"github.com/gin-gonic/gin"
)

func (h *handler) list(c *gin.Context) {
	f := model.Filter{}
	err := ginx.Bind(c, f)
	if err != nil {
		ginx.ResponseError(c, http.StatusBadRequest, err,
			"could not bind filters for list",
			"wrong request",
			h.logger,
		)
		return
	}

	items, errList := h.service.List(c, &f)
	if errList != nil {
		ginx.ResponseError(c, http.StatusInternalServerError, errList,
			"could not get slice of posts",
			"something wrong",
			h.logger,
		)
	}

	response := openapi.ListPostsResponse{Posts: make([]openapi.Post, len(items))}
	for i := range items {
		response.Posts[i] = preparePostResponse(items[i])
	}

	c.JSON(http.StatusOK, response)
}
