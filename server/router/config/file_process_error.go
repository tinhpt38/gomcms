package config

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type FileProcessErrorRouter struct {}

// InitFileProcessErrorRouter 初始化 Lỗi ghi nhận từ quá trình xử lý tệp tin 路由信息
func (s *FileProcessErrorRouter) InitFileProcessErrorRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	fileProcessErrorRouter := Router.Group("fileProcessError").Use(middleware.OperationRecord())
	fileProcessErrorRouterWithoutRecord := Router.Group("fileProcessError")
	fileProcessErrorRouterWithoutAuth := PublicRouter.Group("fileProcessError")
	{
		fileProcessErrorRouter.POST("createFileProcessError", fileProcessErrorApi.CreateFileProcessError)   // 新建Lỗi ghi nhận từ quá trình xử lý tệp tin
		fileProcessErrorRouter.DELETE("deleteFileProcessError", fileProcessErrorApi.DeleteFileProcessError) // 删除Lỗi ghi nhận từ quá trình xử lý tệp tin
		fileProcessErrorRouter.DELETE("deleteFileProcessErrorByIds", fileProcessErrorApi.DeleteFileProcessErrorByIds) // 批量删除Lỗi ghi nhận từ quá trình xử lý tệp tin
		fileProcessErrorRouter.PUT("updateFileProcessError", fileProcessErrorApi.UpdateFileProcessError)    // 更新Lỗi ghi nhận từ quá trình xử lý tệp tin
	}
	{
		fileProcessErrorRouterWithoutRecord.GET("findFileProcessError", fileProcessErrorApi.FindFileProcessError)        // 根据ID获取Lỗi ghi nhận từ quá trình xử lý tệp tin
		fileProcessErrorRouterWithoutRecord.GET("getFileProcessErrorList", fileProcessErrorApi.GetFileProcessErrorList)  // 获取Lỗi ghi nhận từ quá trình xử lý tệp tin列表
	}
	{
	    fileProcessErrorRouterWithoutAuth.GET("getFileProcessErrorPublic", fileProcessErrorApi.GetFileProcessErrorPublic)  // 获取Lỗi ghi nhận từ quá trình xử lý tệp tin列表
	}
}
