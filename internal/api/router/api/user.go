package api

import (
	"github.com/gin-gonic/gin"
	"member-link-lite/internal/api/controllers"
	"member-link-lite/internal/api/middleware"
)

// RegisterUserRoutes 注册用户相关路由
func RegisterUserRoutes(rg *gin.RouterGroup) {
	userController := controllers.NewUserController()

	user := rg.Group("/user")
	user.Use(middleware.JWTAuth()) // 所有用户路由都需要认证
	{
		// 获取个人信息
		user.GET("/profile", userController.GetProfile)

		// 更新个人信息
		user.PUT("/profile", userController.UpdateProfile)

		// 修改密码
		user.PUT("/password", userController.ChangePassword)

		// 上传头像
		user.POST("/avatar", userController.UploadAvatar)
	}
}
