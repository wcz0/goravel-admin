package services

import (
	"goravel/app/models/admin"
	"goravel/app/response"

	"github.com/goravel/framework/facades"
)

type AdminService[T admin.Model] struct {
	*response.ErrorHandler
	Model T
}

func NewAdminService[T admin.Model](model T) *AdminService[T] {
	return &AdminService[T]{
		ErrorHandler: response.NewErrorHandler(),
		Model:        model,
	}
}

func (s *AdminService[T]) GetDetail(id any) (T, error) {
	var model T
	if err := facades.Orm().Query().Find(&model, id); err != nil {
		return model, err
	}
	return model, nil
}
