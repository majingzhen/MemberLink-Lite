package common

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Example: 如何在控制器中使用分页组件
func ExampleUserController(c *gin.Context, db *gorm.DB) {
	// 1. 解析分页参数
	pageReq, err := ParsePageRequest(c)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	// 2. 执行分页查询
	var users []interface{} // 实际使用时替换为具体的用户模型
	result, err := PaginateQuery(db.Model(&users), pageReq, &users)
	if err != nil {
		ServerError(c, "查询失败")
		return
	}

	// 3. 返回分页结果
	SuccessPage(c, result.List, result.Total, result.Page, result.PageSize)
}

// Example: 如何在服务层使用分页组件
func ExampleUserService(db *gorm.DB, pageReq *PageRequest, keyword string) (*PaginateResult, error) {
	// 定义查询条件
	conditions := []func(*gorm.DB) *gorm.DB{
		func(db *gorm.DB) *gorm.DB {
			if keyword != "" {
				return db.Where("username LIKE ? OR nickname LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
			}
			return db
		},
		func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", 1) // 只查询激活用户
		},
	}

	// 执行分页查询
	var users []interface{} // 实际使用时替换为具体的用户模型
	return PaginateQueryWithCondition(db.Model(&users), pageReq, &users, conditions...)
}

// Example: 如何使用分页服务
func ExamplePaginateService(db *gorm.DB) {
	// 创建分页服务
	service := NewBasePaginateService(db, "User")

	// 使用服务进行分页查询
	pageReq := NewPageRequest(1, 10)
	result, err := service.Paginate(pageReq)
	if err != nil {
		// 处理错误
		return
	}

	// 使用结果
	_ = result
}

// Example: 如何检查分页有效性
func ExamplePageValidation(pageReq *PageRequest, total int64) bool {
	// 验证分页参数
	if err := pageReq.Validate(); err != nil {
		return false
	}

	// 检查页码是否有效
	if !pageReq.IsValidPage(total) {
		return false
	}

	return true
}

// Example: 如何获取分页信息
func ExamplePageInfo(pageReq *PageRequest, total int64) *PageInfo {
	// 获取详细的分页信息
	pageInfo := pageReq.ToPageInfo(total)

	// 可以用于前端显示分页控件
	// pageInfo.HasNext - 是否有下一页
	// pageInfo.HasPrev - 是否有上一页
	// pageInfo.Pages - 总页数

	return pageInfo
}
