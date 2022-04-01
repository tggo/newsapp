package http

import (
	"net/http"
	"strconv"

	"boostersNews/api/openapi"
	"boostersNews/internal/app/ginx"
	"boostersNews/internal/news/model"

	"github.com/gin-gonic/gin"
)

func (h *handler) get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ginx.ResponseError(c, http.StatusBadRequest, err,
			"could not bind id to int",
			"wrong request id",
			h.logger,
		)
		return
	}

	item, errGet := h.service.Get(c, int64(id))
	if errGet != nil {
		ginx.ResponseError(c, http.StatusInternalServerError, errGet,
			"could not get one posts",
			"something wrong when get post",
			h.logger,
		)
	}

	c.JSON(http.StatusOK, openapi.OnePostResponse{Post: preparePostResponse(item)})
}

func preparePostResponse(post *model.Post) (p openapi.Post) {
	p.Id = post.ID
	p.Title = post.Title
	p.Content = post.Body
	p.CreatedAt = post.CreatedAt.Unix()
	p.UpdatedAt = post.UpdatedAt.Unix()
	if post.UpdatedAt.IsZero() {
		p.UpdatedAt = post.CreatedAt.Unix()
	}
	return p
}
