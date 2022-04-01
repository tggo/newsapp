// Package services simple all service object
package services

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"

	news "boostersNews/internal/news/service"
)

var ErrNotInitialize = errors.New("not initialize service")

// Services interface for services
type Services struct {
	newsService news.Service
}

func NewServices(logger *zap.Logger) *Services {
	return &Services{}
}
func (s *Services) SetNews(newsService news.Service) {
	s.newsService = newsService
}

func (s *Services) GetNews() news.Service {
	if s.newsService == nil {
		panic(ErrNotInitialize)
	}
	return s.newsService
}
