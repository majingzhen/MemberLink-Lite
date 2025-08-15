package common

import "fmt"

// CustomError 自定义错误类型
type CustomError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Error 实现error接口
func (e *CustomError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// NewCustomError 创建自定义错误
func NewCustomError(code int, message string, details ...string) *CustomError {
	err := &CustomError{
		Code:    code,
		Message: message,
	}
	if len(details) > 0 {
		err.Details = details[0]
	}
	return err
}

// 预定义错误
var (
	// 通用错误
	ErrInvalidParams = NewCustomError(CodeBadRequest, "参数错误")
	ErrUnauthorized  = NewCustomError(CodeUnauthorized, "未授权访问")
	ErrForbidden     = NewCustomError(CodeForbidden, "禁止访问")
	ErrNotFound      = NewCustomError(CodeNotFound, "资源不存在")
	ErrConflict      = NewCustomError(CodeConflict, "资源冲突")
	ErrServerError   = NewCustomError(CodeServerError, "服务器内部错误")

	// 用户相关错误
	ErrUserNotFound    = NewCustomError(CodeNotFound, "用户不存在")
	ErrUserExists      = NewCustomError(CodeConflict, "用户已存在")
	ErrPhoneExists     = NewCustomError(CodeConflict, "手机号已存在")
	ErrEmailExists     = NewCustomError(CodeConflict, "邮箱已存在")
	ErrInvalidPassword = NewCustomError(CodeBadRequest, "密码错误")
	ErrPasswordTooWeak = NewCustomError(CodeBadRequest, "密码强度不足")
	ErrInvalidEmail    = NewCustomError(CodeBadRequest, "邮箱格式错误")
	ErrInvalidPhone    = NewCustomError(CodeBadRequest, "手机号格式错误")
	ErrUserDisabled    = NewCustomError(CodeForbidden, "用户已被禁用")
	ErrEditUserAvatar  = NewCustomError(CodeServerError, "编辑用户头像失败")

	// 认证相关错误
	ErrInvalidToken   = NewCustomError(CodeUnauthorized, "令牌无效")
	ErrTokenExpired   = NewCustomError(CodeUnauthorized, "令牌已过期")
	ErrTokenMalformed = NewCustomError(CodeUnauthorized, "令牌格式错误")

	// 文件相关错误
	ErrFileNotFound          = NewCustomError(CodeNotFound, "文件不存在")
	ErrFileTooBig            = NewCustomError(CodeBadRequest, "文件过大")
	ErrInvalidFileType       = NewCustomError(CodeBadRequest, "文件类型不支持")
	ErrUploadFailed          = NewCustomError(CodeServerError, "文件上传失败")
	ErrOpenFileFailed        = NewCustomError(CodeServerError, "文件打开失败")
	ErrGetFileUrlFailed      = NewCustomError(CodeServerError, "获取文件URL失败")
	ErrSearchFileFailed      = NewCustomError(CodeServerError, "搜索文件失败")
	ErrDeleteFileFailed      = NewCustomError(CodeServerError, "删除文件失败")
	ErrSearchFileListFailed  = NewCustomError(CodeServerError, "搜索文件列表失败")
	ErrSearchFileCountFailed = NewCustomError(CodeServerError, "搜索文件数量失败")
	ErrSaveFileFailed        = NewCustomError(CodeServerError, "保存文件失败")
	// 计算文件哈希失败
	ErrCalculateFileHashFailed = NewCustomError(CodeServerError, "计算文件哈希失败")
	// 生成签名URL失败
	ErrGetSignedURLFailed = NewCustomError(CodeServerError, "生成签名URL失败")

	// 数据库相关错误
	ErrDatabaseError  = NewCustomError(CodeServerError, "数据库操作失败")
	ErrRecordNotFound = NewCustomError(CodeNotFound, "记录不存在")
	ErrDuplicateKey   = NewCustomError(CodeConflict, "数据重复")

	// 业务相关错误
	ErrInsufficientBalance = NewCustomError(CodeBadRequest, "余额不足")
	ErrInsufficientPoints  = NewCustomError(CodeBadRequest, "积分不足")
	ErrInvalidOperation    = NewCustomError(CodeBadRequest, "无效操作")
)

// ValidationError 参数验证错误
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors 多个验证错误
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

// Error 实现error接口
func (ve *ValidationErrors) Error() string {
	if len(ve.Errors) == 0 {
		return "validation failed"
	}
	return fmt.Sprintf("validation failed: %s", ve.Errors[0].Message)
}

// Add 添加验证错误
func (ve *ValidationErrors) Add(field, message string) {
	ve.Errors = append(ve.Errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

// HasErrors 检查是否有错误
func (ve *ValidationErrors) HasErrors() bool {
	return len(ve.Errors) > 0
}

// NewValidationErrors 创建验证错误集合
func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		Errors: make([]ValidationError, 0),
	}
}
