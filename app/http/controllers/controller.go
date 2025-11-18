package controllers

import (
    "goravel/app/response"
    "strings"
    "strconv"

    "github.com/goravel/framework/contracts/http"
    validationPkg "github.com/goravel/framework/validation"
    contractsValidation "github.com/goravel/framework/contracts/validation"
    "github.com/goravel/framework/facades"
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
func (c *Controller) HandleValidationErrors(ctx http.Context, rules map[string]string, messages map[string]string) (bool, http.Response) {
    input := ctx.Request().All()
    for field, rule := range rules {
        if strings.Contains(rule, "int") {
            if v, ok := input[field]; ok {
                switch t := v.(type) {
                case float64:
                    input[field] = int(t)
                case bool:
                    if t { input[field] = 1 } else { input[field] = 0 }
                case string:
                    lv := strings.ToLower(t)
                    switch lv {
                    case "true":
                        input[field] = 1
                    case "false":
                        input[field] = 0
                    default:
                        if iv, err := strconv.Atoi(t); err == nil { input[field] = iv }
                    }
                }
            }
        }
    }
    var validator contractsValidation.Validator
    var err error
    if messages != nil {
        validator, err = facades.Validation().Make(input, rules, validationPkg.Messages(messages))
    } else {
        validator, err = facades.Validation().Make(input, rules)
    }
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
