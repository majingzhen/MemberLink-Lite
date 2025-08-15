package common

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// PaginationUsageExample 分页组件使用示例
// 这个文件展示了如何在服务层中使用优化后的分页组件

// ExampleService 示例服务
type ExampleService struct {
	db *gorm.DB
}

// ExampleModel 示例模型
type ExampleModel struct {
	ID   uint64 `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// ExampleRequest 示例请求
type ExampleRequest struct {
	PageRequest
	Name string `json:"name" form:"name"`
	Type string `json:"type" form:"type"`
}

// GetExampleList 获取示例列表 - 基础用法
func (s *ExampleService) GetExampleList(ctx context.Context, req *ExampleRequest) (*PaginateResult, error) {
	var records []ExampleModel

	// 验证并设置默认分页参数
	if err := req.PageRequest.ValidateAndSetDefaults(); err != nil {
		return nil, err
	}

	// 构建查询条件
	conditions := []func(*gorm.DB) *gorm.DB{}

	// 添加名称筛选
	if req.Name != "" {
		conditions = append(conditions, func(db *gorm.DB) *gorm.DB {
			return db.Where("name LIKE ?", "%"+req.Name+"%")
		})
	}

	// 添加类型筛选
	if req.Type != "" {
		conditions = append(conditions, func(db *gorm.DB) *gorm.DB {
			return db.Where("type = ?", req.Type)
		})
	}

	// 执行分页查询
	result, err := PaginateQueryWithModel(
		s.db.WithContext(ctx),
		&req.PageRequest,
		&ExampleModel{},
		&records,
		conditions...,
	)
	if err != nil {
		return nil, fmt.Errorf("查询示例记录失败: %w", err)
	}

	return result, nil
}

// GetExampleListSimple 获取示例列表 - 简化用法
func (s *ExampleService) GetExampleListSimple(ctx context.Context, req *ExampleRequest) (*PaginateResult, error) {
	var records []ExampleModel

	// 直接使用PaginateQuery进行简单查询
	query := s.db.WithContext(ctx).Model(&ExampleModel{})

	// 添加筛选条件
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	// 执行分页查询
	result, err := PaginateQuery(query, &req.PageRequest, &records)
	if err != nil {
		return nil, fmt.Errorf("查询示例记录失败: %w", err)
	}

	return result, nil
}

// GetExampleListWithService 使用BasePaginateService
func (s *ExampleService) GetExampleListWithService(ctx context.Context, req *ExampleRequest) (*PaginateResult, error) {
	// 创建分页服务
	paginateService := NewBasePaginateService(s.db.WithContext(ctx), &ExampleModel{})

	// 构建查询条件
	conditions := []func(*gorm.DB) *gorm.DB{}

	if req.Name != "" {
		conditions = append(conditions, func(db *gorm.DB) *gorm.DB {
			return db.Where("name LIKE ?", "%"+req.Name+"%")
		})
	}

	if req.Type != "" {
		conditions = append(conditions, func(db *gorm.DB) *gorm.DB {
			return db.Where("type = ?", req.Type)
		})
	}

	// 执行分页查询
	return paginateService.Paginate(&req.PageRequest, conditions...)
}

// 使用示例：
// 1. 在控制器中解析分页参数
// func (c *ExampleController) GetList(ctx *gin.Context) {
//     req := &ExampleRequest{}
//     if err := ctx.ShouldBindQuery(req); err != nil {
//         common.BadRequest(ctx, err.Error())
//         return
//     }
//
//     result, err := c.service.GetExampleList(ctx, req)
//     if err != nil {
//         common.ServerError(ctx, err.Error())
//         return
//     }
//
//     common.Success(ctx, result)
// }

// 2. 或者使用ParsePageRequest从上下文解析
// func (c *ExampleController) GetListAuto(ctx *gin.Context) {
//     pageReq, err := common.ParsePageRequest(ctx)
//     if err != nil {
//         common.BadRequest(ctx, err.Error())
//         return
//     }
//
//     req := &ExampleRequest{
//         PageRequest: *pageReq,
//         Name:        ctx.Query("name"),
//         Type:        ctx.Query("type"),
//     }
//
//     result, err := c.service.GetExampleList(ctx, req)
//     if err != nil {
//         common.ServerError(ctx, err.Error())
//         return
//     }
//
//     common.Success(ctx, result)
// }
