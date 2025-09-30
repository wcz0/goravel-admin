package controllers

import (
	"goravel/app/response"
	"strings"

	"github.com/goravel/framework/contracts/http"
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

// HandleValidationErrors 统一处理验证错误，返回字符串格式的错误信息
func (c *Controller) HandleValidationErrors(ctx http.Context, rules map[string]string) (bool, http.Response) {
	validator, err := ctx.Request().Validate(rules)
	if err != nil {
		return true, c.FailMsg(ctx, err.Error())
	}
	if validator.Fails() {
		errors := validator.Errors().All()
		var errorMessages []string
		for _, errMap := range errors {
			for _, msg := range errMap {
				errorMessages = append(errorMessages, msg)
			}
		}
		return true, c.FailMsg(ctx, strings.Join(errorMessages, "; "))
	}
	return false, nil
}
