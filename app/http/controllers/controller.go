package controllers

import (
	"goravel/app/response"

)

type Controller struct {
	*response.ErrorHandler
	// Controller AdminController
}

func NewController() *Controller {
	return &Controller{
		ErrorHandler: response.NewErrorHandler(),
		// Controller:   nil,
	}
}
