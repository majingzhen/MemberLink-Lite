package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// BindAndValidate 绑定并验证请求参数的通用函数
func BindAndValidate(c *gin.Context, req interface{}) error {
	// 绑定JSON参数
	if err := c.ShouldBindJSON(req); err != nil {
		// 处理绑定错误
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			// 处理验证错误
			validationErrs := NewValidationErrors()
			for _, fieldError := range validationErrors {
				message := getValidationErrorMessage(fieldError)
				validationErrs.Add(fieldError.Field(), message)
			}

			BadRequestWithData(c, "参数验证失败", validationErrs.Errors)
			return err
		}

		// 其他绑定错误
		BadRequest(c, "请求参数格式错误: "+err.Error())
		return err
	}

	return nil
}

// BindQueryAndValidate 绑定并验证查询参数的通用函数
func BindQueryAndValidate(c *gin.Context, req interface{}) error {
	// 绑定查询参数
	if err := c.ShouldBindQuery(req); err != nil {
		// 处理绑定错误
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			// 处理验证错误
			validationErrs := NewValidationErrors()
			for _, fieldError := range validationErrors {
				message := getValidationErrorMessage(fieldError)
				validationErrs.Add(fieldError.Field(), message)
			}

			BadRequestWithData(c, "查询参数验证失败", validationErrs.Errors)
			return err
		}

		// 其他绑定错误
		BadRequest(c, "查询参数格式错误: "+err.Error())
		return err
	}

	return nil
}

// BindURIAndValidate 绑定并验证URI参数的通用函数
func BindURIAndValidate(c *gin.Context, req interface{}) error {
	// 绑定URI参数
	if err := c.ShouldBindUri(req); err != nil {
		// 处理绑定错误
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			// 处理验证错误
			validationErrs := NewValidationErrors()
			for _, fieldError := range validationErrors {
				message := getValidationErrorMessage(fieldError)
				validationErrs.Add(fieldError.Field(), message)
			}

			BadRequestWithData(c, "路径参数验证失败", validationErrs.Errors)
			return err
		}

		// 其他绑定错误
		BadRequest(c, "路径参数格式错误: "+err.Error())
		return err
	}

	return nil
}

// BadRequestWithData 带数据的错误请求响应
func BadRequestWithData(c *gin.Context, message string, data interface{}) {
	response := APIResponse{
		Code:    CodeBadRequest,
		Message: message,
		Data:    data,
		TraceID: getTraceID(c),
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
		return "长度不能少于" + fe.Param() + "个字符"
	case "max":
		return "长度不能超过" + fe.Param() + "个字符"
	case "len":
		return "长度必须为" + fe.Param() + "个字符"
	case "numeric":
		return "必须为数字"
	case "alpha":
		return "只能包含字母"
	case "alphanum":
		return "只能包含字母和数字"
	case "mobile":
		return "手机号格式不正确"
	case "gte":
		return "值不能小于" + fe.Param()
	case "lte":
		return "值不能大于" + fe.Param()
	case "gt":
		return "值必须大于" + fe.Param()
	case "lt":
		return "值必须小于" + fe.Param()
	case "oneof":
		return "值必须是以下之一: " + fe.Param()
	default:
		return "字段验证失败: " + fe.Tag()
	}
}
