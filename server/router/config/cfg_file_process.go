package config

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CfgFileProcessRouter struct {}

// InitCfgFileProcessRouter 初始化 Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống 路由信息
func (s *CfgFileProcessRouter) InitCfgFileProcessRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	FileProcessRouter := Router.Group("FileProcess").Use(middleware.OperationRecord())
	FileProcessRouterWithoutRecord := Router.Group("FileProcess")
	FileProcessRouterWithoutAuth := PublicRouter.Group("FileProcess")
	{
		FileProcessRouter.POST("createCfgFileProcess", FileProcessApi.CreateCfgFileProcess)   // 新建Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống
		FileProcessRouter.DELETE("deleteCfgFileProcess", FileProcessApi.DeleteCfgFileProcess) // 删除Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống
		FileProcessRouter.DELETE("deleteCfgFileProcessByIds", FileProcessApi.DeleteCfgFileProcessByIds) // 批量删除Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống
		FileProcessRouter.PUT("updateCfgFileProcess", FileProcessApi.UpdateCfgFileProcess)    // 更新Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống
	}
	{
		FileProcessRouterWithoutRecord.GET("findCfgFileProcess", FileProcessApi.FindCfgFileProcess)        // 根据ID获取Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống
		FileProcessRouterWithoutRecord.GET("getCfgFileProcessList", FileProcessApi.GetCfgFileProcessList)  // 获取Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống列表
	}
	{
	    FileProcessRouterWithoutAuth.GET("getCfgFileProcessPublic", FileProcessApi.GetCfgFileProcessPublic)  // 获取Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống列表
	}
}
