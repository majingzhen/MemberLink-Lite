package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse 统一API响应格式
type APIResponse struct {
	Code    int         `json:"code"`               // 状态码
	Message string      `json:"message"`            // 消息
	Data    interface{} `json:"data"`               // 数据
	TraceID string      `json:"trace_id,omitempty"` // 追踪ID
}

// PageResponse 分页响应格式
type PageResponse struct {
	List     interface{} `json:"list"`      // 数据列表
	Total    int64       `json:"total"`     // 总数
	Page     int         `json:"page"`      // 当前页
	PageSize int         `json:"page_size"` // 页大小
	Pages    int         `json:"pages"`     // 总页数
}

// 响应码定义
const (
	CodeSuccess      = 200 // 成功
	CodeBadRequest   = 400 // 请求错误
	CodeUnauthorized = 401 // 未授权
	CodeForbidden    = 403 // 禁止访问
	CodeNotFound     = 404 // 未找到
	CodeConflict     = 409 // 冲突
	CodeServerError  = 500 // 服务器错误
)

// 响应消息定义
const (
	MsgSuccess      = "操作成功"
	MsgBadRequest   = "请求参数错误"
	MsgUnauthorized = "未授权访问"
	MsgForbidden    = "禁止访问"
	MsgNotFound     = "资源不存在"
	MsgConflict     = "资源冲突"
	MsgServerError  = "服务器内部错误"
)

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	response := APIResponse{
		Code:    CodeSuccess,
		Message: MsgSuccess,
		Data:    data,
		TraceID: getTraceID(c),
	}
	c.JSON(http.StatusOK, response)
}

// SuccessWithMessage 带自定义消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	response := APIResponse{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
		TraceID: getTraceID(c),
	}
	c.JSON(http.StatusOK, response)
}

// SuccessPage 分页成功响应
func SuccessPage(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	pages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		pages++
	}

	pageData := PageResponse{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Pages:    pages,
	}

	response := APIResponse{
		Code:    CodeSuccess,
		Message: MsgSuccess,
		Data:    pageData,
		TraceID: getTraceID(c),
	}
	c.JSON(http.StatusOK, response)
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	response := APIResponse{
		Code:    code,
		Message: message,
		Data:    nil,
		TraceID: getTraceID(c),
	}
	c.JSON(http.StatusOK, response)
}

// BadRequest 请求错误响应
func BadRequest(c *gin.Context, message string) {
	if message == "" {
		message = MsgBadRequest
	}
	Error(c, CodeBadRequest, message)
}

// Unauthorized 未授权响应
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = MsgUnauthorized
	}
	Error(c, CodeUnauthorized, message)
}

// Forbidden 禁止访问响应
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = MsgForbidden
	}
	Error(c, CodeForbidden, message)
}

// NotFound 未找到响应
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = MsgNotFound
	}
	Error(c, CodeNotFound, message)
}

// Conflict 冲突响应
func Conflict(c *gin.Context, message string) {
	if message == "" {
		message = MsgConflict
	}
	Error(c, CodeConflict, message)
}

// ServerError 服务器错误响应
func ServerError(c *gin.Context, message string) {
	if message == "" {
		message = MsgServerError
	}
	Error(c, CodeServerError, message)
}

// ErrorResponse 错误响应（带数据）
func ErrorResponse(c *gin.Context, httpStatus int, message string, data interface{}) {
	response := APIResponse{
		Code:    httpStatus,
		Message: message,
		Data:    data,
		TraceID: getTraceID(c),
	}
	c.JSON(http.StatusOK, response) // 统一返回200状态码，错误信息在响应体中
}

// SuccessResponse 成功响应（带消息和数据）
func SuccessResponse(c *gin.Context, message string, data interface{}) {
	response := APIResponse{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
		TraceID: getTraceID(c),
	}
	c.JSON(http.StatusOK, response)
}

// getTraceID 获取追踪ID
func getTraceID(c *gin.Context) string {
	if traceID, exists := c.Get("trace_id"); exists {
		if id, ok := traceID.(string); ok {
			return id
		}
	}
	return ""
}
