package common

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PageRequest 分页请求参数
// 用于处理分页查询的请求参数，支持页码和页大小的验证
type PageRequest struct {
	Page     int `json:"page" form:"page" binding:"min=1"`                   // 页码，从1开始
	PageSize int `json:"page_size" form:"page_size" binding:"min=1,max=100"` // 页大小，最大100
}

// DefaultPageRequest 默认分页参数
var DefaultPageRequest = PageRequest{
	Page:     1,
	PageSize: 10,
}

// NewPageRequest 创建分页请求
func NewPageRequest(page, pageSize int) *PageRequest {
	if page <= 0 {
		page = DefaultPageRequest.Page
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = DefaultPageRequest.PageSize
	}
	return &PageRequest{
		Page:     page,
		PageSize: pageSize,
	}
}

// GetOffset 获取偏移量
func (p *PageRequest) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// GetLimit 获取限制数量
func (p *PageRequest) GetLimit() int {
	return p.PageSize
}

// Validate 验证分页参数
func (p *PageRequest) Validate() error {
	if p.Page <= 0 {
		return NewCustomError(CodeBadRequest, "页码必须大于0")
	}
	if p.PageSize <= 0 {
		return NewCustomError(CodeBadRequest, "页大小必须大于0")
	}
	if p.PageSize > 100 {
		return NewCustomError(CodeBadRequest, "页大小不能超过100")
	}
	return nil
}

// ParsePageRequest 从Gin上下文解析分页参数
func ParsePageRequest(c *gin.Context) (*PageRequest, error) {
	page := 1
	pageSize := 10

	// 解析页码
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// 解析页大小
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	pageReq := &PageRequest{
		Page:     page,
		PageSize: pageSize,
	}

	return pageReq, pageReq.Validate()
}

// Paginate GORM分页插件
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 || pageSize > 100 {
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// PaginateWithRequest 使用PageRequest进行分页
func PaginateWithRequest(req *PageRequest) func(db *gorm.DB) *gorm.DB {
	return Paginate(req.Page, req.PageSize)
}

// PaginateResult 分页查询结果
// 包含分页数据和分页信息的完整结果结构
type PaginateResult struct {
	List     interface{} `json:"list"`      // 数据列表
	Total    int64       `json:"total"`     // 总记录数
	Page     int         `json:"page"`      // 当前页码
	PageSize int         `json:"page_size"` // 每页大小
	Pages    int         `json:"pages"`     // 总页数
}

// NewPaginateResult 创建分页结果
func NewPaginateResult(list interface{}, total int64, page, pageSize int) *PaginateResult {
	pages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		pages++
	}

	return &PaginateResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Pages:    pages,
	}
}

// PaginateQuery 执行分页查询
func PaginateQuery(db *gorm.DB, req *PageRequest, result interface{}) (*PaginateResult, error) {
	// 获取总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	// 执行分页查询
	if err := db.Scopes(PaginateWithRequest(req)).Find(result).Error; err != nil {
		return nil, err
	}

	return NewPaginateResult(result, total, req.Page, req.PageSize), nil
}

// PaginateQueryWithCondition 带条件的分页查询
func PaginateQueryWithCondition(db *gorm.DB, req *PageRequest, result interface{}, conditions ...func(*gorm.DB) *gorm.DB) (*PaginateResult, error) {
	// 应用查询条件
	query := db
	for _, condition := range conditions {
		query = condition(query)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 执行分页查询
	if err := query.Scopes(PaginateWithRequest(req)).Find(result).Error; err != nil {
		return nil, err
	}

	return NewPaginateResult(result, total, req.Page, req.PageSize), nil
}

// PaginateService 分页服务接口
type PaginateService interface {
	Paginate(req *PageRequest, conditions ...func(*gorm.DB) *gorm.DB) (*PaginateResult, error)
}

// BasePaginateService 基础分页服务
type BasePaginateService struct {
	db    *gorm.DB
	model interface{}
}

// NewBasePaginateService 创建基础分页服务
func NewBasePaginateService(db *gorm.DB, model interface{}) *BasePaginateService {
	return &BasePaginateService{
		db:    db,
		model: model,
	}
}

// Paginate 执行分页查询
func (s *BasePaginateService) Paginate(req *PageRequest, conditions ...func(*gorm.DB) *gorm.DB) (*PaginateResult, error) {
	// 创建结果切片
	resultType := s.model
	results := make([]interface{}, 0)

	// 执行分页查询
	return PaginateQueryWithCondition(s.db.Model(resultType), req, &results, conditions...)
}

// GetPageInfo 获取分页信息（不查询数据）
func GetPageInfo(total int64, page, pageSize int) *PageInfo {
	pages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		pages++
	}

	return &PageInfo{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Pages:    pages,
		HasNext:  page < pages,
		HasPrev:  page > 1,
	}
}

// PageInfo 分页信息
type PageInfo struct {
	Total    int64 `json:"total"`     // 总数
	Page     int   `json:"page"`      // 当前页
	PageSize int   `json:"page_size"` // 页大小
	Pages    int   `json:"pages"`     // 总页数
	HasNext  bool  `json:"has_next"`  // 是否有下一页
	HasPrev  bool  `json:"has_prev"`  // 是否有上一页
}

// IsValidPage 检查页码是否有效
func (p *PageRequest) IsValidPage(total int64) bool {
	if total == 0 {
		return p.Page == 1
	}
	maxPage := int(total) / p.PageSize
	if int(total)%p.PageSize > 0 {
		maxPage++
	}
	return p.Page <= maxPage
}

// ToPageInfo 转换为分页信息
func (p *PageRequest) ToPageInfo(total int64) *PageInfo {
	return GetPageInfo(total, p.Page, p.PageSize)
}

// CalculatePages 计算总页数
func CalculatePages(total int64, pageSize int) int {
	if total == 0 || pageSize <= 0 {
		return 0
	}
	pages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		pages++
	}
	return pages
}

// PaginateQueryWithModel 使用指定模型执行分页查询
func PaginateQueryWithModel(db *gorm.DB, req *PageRequest, model interface{}, result interface{}, conditions ...func(*gorm.DB) *gorm.DB) (*PaginateResult, error) {
	// 应用查询条件
	query := db.Model(model)
	for _, condition := range conditions {
		query = condition(query)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 执行分页查询
	if err := query.Scopes(PaginateWithRequest(req)).Find(result).Error; err != nil {
		return nil, err
	}

	return NewPaginateResult(result, total, req.Page, req.PageSize), nil
}

// ValidateAndSetDefaults 验证并设置默认分页参数
func (p *PageRequest) ValidateAndSetDefaults() error {
	if p.Page <= 0 {
		p.Page = DefaultPageRequest.Page
	}
	if p.PageSize <= 0 {
		p.PageSize = DefaultPageRequest.PageSize
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	return p.Validate()
}
