package http

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	"boostersNews/api/openapi"
	newsMocks "boostersNews/internal/news/mocks"
	"boostersNews/internal/news/model"
	"boostersNews/internal/news/service"
	"boostersNews/pkg/helpers"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var logger = zap.NewNop()
var swaggerFileName = "../../../../api/swagger.yaml"

func TestList(t *testing.T) {
	mockRepo := &newsMocks.Repository{}
	ginRouter, _, err := helpers.GetRouters(swaggerFileName)
	require.NoError(t, err)

	postOne := model.Post{ID: 1, Title: "title", Body: "Content", CreatedAt: time.Now()}

	newsService := service.NewService(mockRepo, logger)
	pHandler := NewHandler(newsService, logger)
	ginRouter.GET("/posts", pHandler.list)

	t.Run("just right return", func(t *testing.T) {
		mockRepo.On("Find", mock.Anything, mock.Anything).
			Return([]*model.Post{&postOne}, nil).
			Once()

		req, err := http.NewRequest("GET", "/posts", bytes.NewBuffer([]byte(``)))
		require.NoError(t, err)
		data := helpers.TestData{
			Expected: &openapi.ListPostsResponse{
				Posts: []openapi.Post{
					{
						Id:        1,
						Title:     postOne.Title,
						Content:   postOne.Body,
						CreatedAt: postOne.CreatedAt.Unix(),
						UpdatedAt: postOne.CreatedAt.Unix(),
					},
				},
			},
			Status:   http.StatusOK,
			Response: &openapi.ListPostsResponse{},
		}
		helpers.TestHTTPResponse(t, ginRouter, req, data, helpers.Check)
	})

	t.Run("empty response return", func(t *testing.T) {
		mockRepo.On("Find", mock.Anything, mock.Anything).
			Return(nil, nil).
			Once()

		req, err := http.NewRequest("GET", "/posts", bytes.NewBuffer([]byte(``)))
		require.NoError(t, err)
		data := helpers.TestData{
			Expected: &openapi.ListPostsResponse{
				Posts: nil,
			},
			Status:   http.StatusOK,
			Response: &openapi.ListPostsResponse{},
		}
		helpers.TestHTTPResponse(t, ginRouter, req, data, helpers.Check)
	})

}
