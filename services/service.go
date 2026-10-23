package services

import (
	"github.com/wcz0/goravel-admin/response"
)

type Service struct {
	*response.ErrorHandler
}

func NewService() *Service {
	return &Service{
		ErrorHandler: response.NewErrorHandler(),
	}
}