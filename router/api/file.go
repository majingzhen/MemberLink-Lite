package api

import (
	"MemberLink-Lite/controllers"
	"MemberLink-Lite/database"
	"MemberLink-Lite/middleware"
	"MemberLink-Lite/services"

	"github.com/gin-gonic/gin"
)

// RegisterFileRoutes 注册文件管理路由
func RegisterFileRoutes(r *gin.RouterGroup) {
	// 创建文件服务和控制器
	fileService := services.NewFileService(database.GetDB())
	fileController := controllers.NewFileController(fileService)

	// 文件管理路由组
	fileGroup := r.Group("/files")
	{
		// 公开路由 - 文件访问服务（用于本地存储）
		fileGroup.GET("/uploads/*filepath", fileController.ServeFile)

		// 需要认证的路由
		authGroup := fileGroup.Group("")
		authGroup.Use(middleware.JWTAuth())
		{
			// 文件上传
			authGroup.POST("/avatar", fileController.UploadAvatar) // 上传头像
			authGroup.POST("/image", fileController.UploadImage)   // 上传图片
			authGroup.POST("/upload", fileController.UploadFile)   // 上传通用文件

			// 文件管理
			authGroup.GET("", fileController.GetUserFiles)                // 获取用户文件列表
			authGroup.GET("/:id", fileController.GetFileInfo)             // 获取文件信息
			authGroup.GET("/:id/signed-url", fileController.GetSignedURL) // 获取签名URL
			authGroup.DELETE("/:id", fileController.DeleteFile)           // 删除文件
		}
	}
}
