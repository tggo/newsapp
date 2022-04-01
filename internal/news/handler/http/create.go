package http

import (
	"net/http"
	"time"

	"boostersNews/api/openapi"
	"boostersNews/internal/app/ginx"
	"boostersNews/internal/news/model"

	"github.com/gin-gonic/gin"
)

func (h *handler) create(c *gin.Context) {
	requestData := openapi.PostCreate{}
	err := ginx.Bind(c, &requestData)
	if err != nil {
		ginx.ResponseError(c, http.StatusBadRequest, err,
			"could not bind data for create post",
			"wrong request on create post",
			h.logger,
		)
		return
	}
	id, errCreate := h.service.Create(c, requestToPost(requestData))
	if errCreate != nil {
		ginx.ResponseError(c, http.StatusBadRequest, errCreate,
			"could not create post",
			"something wrong on creating process",
			h.logger,
		)
	}

	ginx.ResponseOK(c, id)
}

func requestToPost(post openapi.PostCreate) *model.Post {
	p := model.Post{
		Title:     post.Title,
		Body:      post.Content,
		CreatedAt: time.Now().UTC(),
	}
	return &p
}
