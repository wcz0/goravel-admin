package services

import (
	"reflect"

	"github.com/goravel/framework/contracts/http"

	"github.com/goravel/framework/facades"
)

type AdminService[T any] struct {
	Model T
}

func NewAdminService[T any](model T) *AdminService[T] {
	return &AdminService[T]{
		Model: model,
	}
}

func (s *AdminService[T]) GetDetail(id any) (T, error) {
	var model T
	if err := facades.Orm().Query().Find(&model, id); err != nil {
		return model, err
	}
	return model, nil
}

// 模型保存
func (s *AdminService[T]) Store(ctx http.Context) error {

	allInput := ctx.Request().All()
	t := reflect.TypeOf(s.Model)
	v := reflect.ValueOf(s.Model)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()  // 获取指针指向的类型
		v = v.Elem()  // 获取指针指向的值
	}

	// 遍历结构体字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// 获取字段的 json 标签名作为键名
		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}
		// 如果没有 json 标签，使用字段名
		keyName := jsonTag
		if keyName == "" {
			keyName = field.Name
		}

		// 检查输入中是否有对应的值
		if value, ok := allInput[keyName]; ok {
			fieldValue := v.Field(i)
			if fieldValue.CanSet() {
				// 根据字段类型设置值
				if value != nil {
					setValue := reflect.ValueOf(value)
					if fieldValue.Type() != setValue.Type() && setValue.Type().ConvertibleTo(fieldValue.Type()) {
						// 需要类型转换
						setValue = setValue.Convert(fieldValue.Type())
					}
					fieldValue.Set(setValue)
				}
			}
		}
	}

	// 保存到数据库
	if err := facades.Orm().Query().Create(&s.Model); err != nil {
		return err
	}

	return nil
}

// Success 成功响应
func (s *AdminService[T]) Success(ctx http.Context, msg string) http.Response {
	return ctx.Response().Success().Json(http.Json{
		"code": 200,
		"msg":  msg,
		"data": map[string]any{},
	})
}

// SuccessData 成功响应（带数据）
func (s *AdminService[T]) SuccessData(ctx http.Context, data map[string]any) http.Response {
	return ctx.Response().Success().Json(http.Json{
		"code": 200,
		"msg":  "Success",
		"data": data,
	})
}

// FailMsg 失败响应
func (s *AdminService[T]) FailMsg(ctx http.Context, msg string) http.Response {
	return ctx.Response().Status(400).Json(http.Json{
		"code": 400,
		"msg":  msg,
		"data": map[string]any{},
	})
}
