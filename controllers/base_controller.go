package controllers

import (
	"sort"
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/http"
	contractsValidation "github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/facades"
	validationPkg "github.com/goravel/framework/validation"

	"github.com/wcz0/goravel-admin/response"
)

// Controller 是所有后台控制器的基础实现。
type Controller struct {
	*response.ErrorHandler
}

// NewController 创建基础控制器。
func NewController() *Controller {
	return &Controller{
		ErrorHandler: response.NewErrorHandler(),
	}
}

// normalizeIntInput 将声明为 int 规则的字段值统一转换为整型，避免字符串导致的校验失败。
func normalizeIntInput(input map[string]any, rules map[string]string) {
	for field, rule := range rules {
		if !strings.Contains(rule, "int") {
			continue
		}
		v, ok := input[field]
		if !ok {
			continue
		}
		switch t := v.(type) {
		case float64:
			input[field] = int(t)
		case bool:
			if t {
				input[field] = 1
			} else {
				input[field] = 0
			}
		case string:
			switch strings.ToLower(t) {
			case "true":
				input[field] = 1
			case "false":
				input[field] = 0
			default:
				if iv, err := strconv.Atoi(t); err == nil {
					input[field] = iv
				}
			}
		}
	}
}

// HandleValidationErrors 统一处理验证错误，返回字符串格式的错误信息。
func (c *Controller) HandleValidationErrors(ctx http.Context, rules map[string]string, messages map[string]string) (bool, http.Response) {
	return handleValidationErrors(ctx, c.ErrorHandler, rules, messages)
}

// handleValidationErrors 供各类控制器复用的验证错误处理逻辑。
func handleValidationErrors(ctx http.Context, handler *response.ErrorHandler, rules map[string]string, messages map[string]string) (bool, http.Response) {
	input := ctx.Request().All()
	normalizeIntInput(input, rules)

	var (
		validator contractsValidation.Validator
		err       error
	)
	if messages != nil {
		validator, err = facades.Validation().Make(ctx, input, rules, validationPkg.Messages(messages))
	} else {
		validator, err = facades.Validation().Make(ctx, input, rules)
	}
	if err != nil {
		return true, handler.FailMsg(ctx, err.Error())
	}
	if validator.Fails() {
		all := validator.Errors().All()
		fields := make([]string, 0, len(all))
		for field := range all {
			fields = append(fields, field)
		}
		sort.Strings(fields)

		var errorMessages []string
		for _, field := range fields {
			for _, msg := range all[field] {
				errorMessages = append(errorMessages, msg)
			}
		}
		if len(errorMessages) > 0 {
			return true, handler.FailMsg(ctx, strings.Join(errorMessages, "; "))
		}
		return true, handler.FailMsg(ctx, "验证失败")
	}
	return false, nil
}
