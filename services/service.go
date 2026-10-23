package services

import (
	"errors"

	"github.com/goravel/framework/facades"

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

// ErrRecordNotFound 用于 firstByID 找不到记录时的返回值。
var ErrRecordNotFound = errors.New("记录不存在")

// firstByID 按主键取一条记录。
//
// 不能直接用 Find/First：id 为 0 或记录不存在时两者都不会报错——Find 返回零值，
// First 则把 gorm 的 ErrRecordNotFound 吞掉直接返回 nil。拿零值继续 Save 会把
// 它当成新记录 INSERT（曾经因此写出一条 name 为空的扩展记录），所以这里先
// 用 Count 确认存在性，再取数据。
func firstByID[T integer](id T, target any) error {
	if id <= 0 {
		return ErrRecordNotFound
	}

	total, err := facades.Orm().Query().Model(target).Where("id", id).Count()
	if err != nil {
		return err
	}
	if total == 0 {
		return ErrRecordNotFound
	}

	return facades.Orm().Query().Where("id", id).First(target)
}

// integer 覆盖各表主键可能用到的整数类型。
type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// 列表分页的边界。perPage 直接来自请求，不设上限会变成一次全表查询。
const (
	defaultPageSize = 15
	maxPageSize     = 500
)

// clampPage 把页码收敛到 >= 1。
func clampPage(page int) int {
	if page < 1 {
		return 1
	}

	return page
}

// clampPerPage 把每页条数收敛到 [1, maxPageSize]，非法值回落到默认值。
func clampPerPage(perPage int) int {
	if perPage <= 0 {
		return defaultPageSize
	}
	if perPage > maxPageSize {
		return maxPageSize
	}

	return perPage
}
