package middleware

import (
	"MemberLink-Lite/common"
	"MemberLink-Lite/logger"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ErrorHandler 全局错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// 记录panic堆栈信息
		logger.Error("Panic recovered:", recovered)
		logger.Error("Stack trace:", string(debug.Stack()))

		// 处理不同类型的错误
		switch err := recovered.(type) {
		case *common.CustomError:
			handleCustomError(c, err)
		case *common.ValidationErrors:
			handleValidationErrors(c, err)
		case validator.ValidationErrors:
			handleValidatorErrors(c, err)
		case error:
			handleGenericError(c, err)
		default:
			handleUnknownError(c, recovered)
		}

		c.Abort()
	})
}

// handleCustomError 处理自定义错误
func handleCustomError(c *gin.Context, err *common.CustomError) {
	logger.Error("Custom error:", err.Error())

	response := common.APIResponse{
		Code:    err.Code,
		Message: err.Message,
		Data:    nil,
		TraceID: getTraceID(c),
	}

	// 如果有详细信息且是开发环境，可以包含详细信息
	if err.Details != "" && gin.Mode() == gin.DebugMode {
		response.Data = map[string]string{"details": err.Details}
	}

	c.JSON(http.StatusOK, response)
}

// handleValidationErrors 处理验证错误
func handleValidationErrors(c *gin.Context, err *common.ValidationErrors) {
	logger.Error("Validation errors:", err.Error())

	response := common.APIResponse{
		Code:    common.CodeBadRequest,
		Message: "参数验证失败",
		Data:    err.Errors,
		TraceID: getTraceID(c),
	}

	c.JSON(http.StatusOK, response)
}

// handleValidatorErrors 处理gin validator错误
func handleValidatorErrors(c *gin.Context, err validator.ValidationErrors) {
	logger.Error("Validator errors:", err.Error())

	validationErrors := common.NewValidationErrors()

	for _, fieldError := range err {
		message := getValidationErrorMessage(fieldError)
		validationErrors.Add(fieldError.Field(), message)
	}

	response := common.APIResponse{
		Code:    common.CodeBadRequest,
		Message: "参数验证失败",
		Data:    validationErrors.Errors,
		TraceID: getTraceID(c),
	}

	c.JSON(http.StatusOK, response)
}

// handleGenericError 处理通用错误
func handleGenericError(c *gin.Context, err error) {
	logger.Error("Generic error:", err.Error())

	response := common.APIResponse{
		Code:    common.CodeServerError,
		Message: common.MsgServerError,
		Data:    nil,
		TraceID: getTraceID(c),
	}

	// 开发环境下显示详细错误信息
	if gin.Mode() == gin.DebugMode {
		response.Data = map[string]string{"error": err.Error()}
	}

	c.JSON(http.StatusOK, response)
}

// handleUnknownError 处理未知错误
func handleUnknownError(c *gin.Context, recovered interface{}) {
	logger.Error("Unknown error:", recovered)

	response := common.APIResponse{
		Code:    common.CodeServerError,
		Message: common.MsgServerError,
		Data:    nil,
		TraceID: getTraceID(c),
	}

	// 开发环境下显示详细错误信息
	if gin.Mode() == gin.DebugMode {
		response.Data = map[string]interface{}{"error": fmt.Sprintf("%v", recovered)}
	}

	c.JSON(http.StatusOK, response)
}

// getValidationErrorMessage 获取验证错误消息
func getValidationErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "此字段为必填项"
	case "email":
		return "邮箱格式不正确"
	case "min":
		return fmt.Sprintf("长度不能少于%s个字符", fe.Param())
	case "max":
		return fmt.Sprintf("长度不能超过%s个字符", fe.Param())
	case "len":
		return fmt.Sprintf("长度必须为%s个字符", fe.Param())
	case "numeric":
		return "必须为数字"
	case "alpha":
		return "只能包含字母"
	case "alphanum":
		return "只能包含字母和数字"
	case "mobile":
		return "手机号格式不正确"
	default:
		return fmt.Sprintf("字段验证失败: %s", fe.Tag())
	}
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
