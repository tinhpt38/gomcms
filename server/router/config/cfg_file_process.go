package config

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CfgFileProcessRouter struct{}

// InitCfgFileProcessRouter 初始化 Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống 路由信息
func (s *CfgFileProcessRouter) InitCfgFileProcessRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	FileProcessRouter := Router.Group("fileProcess").Use(middleware.OperationRecord())
	FileProcessRouterWithoutRecord := Router.Group("fileProcess")
	FileProcessRouterWithoutAuth := PublicRouter.Group("fileProcess")
	{
		FileProcessRouter.POST("createCfgFileProcess", fileProcessApi.CreateCfgFileProcess)             // 新建Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống
		FileProcessRouter.DELETE("deleteCfgFileProcess", fileProcessApi.DeleteCfgFileProcess)           // 删除Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống
		FileProcessRouter.DELETE("deleteCfgFileProcessByIds", fileProcessApi.DeleteCfgFileProcessByIds) // 批量删除Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống
		FileProcessRouter.PUT("updateCfgFileProcess", fileProcessApi.UpdateCfgFileProcess)              // 更新Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống
	}
	{
		FileProcessRouterWithoutRecord.GET("findCfgFileProcess", fileProcessApi.FindCfgFileProcess)       // 根据ID获取Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống
		FileProcessRouterWithoutRecord.GET("getCfgFileProcessList", fileProcessApi.GetCfgFileProcessList) // 获取Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống列表
		FileProcessRouterWithoutRecord.GET("getPercentFileProcess", fileProcessApi.GetPercentFileProcess) // Lấy tỷ lệ tiến trình nhập xuất dữ liệu Excel của hệ thống
	}
	{
		FileProcessRouterWithoutAuth.GET("getCfgFileProcessPublic", fileProcessApi.GetCfgFileProcessPublic) // 获取Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống列表
	}
}
