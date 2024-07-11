package services

import (
	"goravel/app/response"
)

type Service struct {
	*response.ErrorHandler
}

func NewService() *Service {
	return &Service{
		ErrorHandler: response.NewErrorHandler(),
	}
}