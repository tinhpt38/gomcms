package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type AttendanceAgencyRouter struct {}

// InitAttendanceAgencyRouter 初始化 Đơn vị 路由信息
func (s *AttendanceAgencyRouter) InitAttendanceAgencyRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	attendanceAgencyRouter := Router.Group("attendanceAgency").Use(middleware.OperationRecord())
	attendanceAgencyRouterWithoutRecord := Router.Group("attendanceAgency")
	attendanceAgencyRouterWithoutAuth := PublicRouter.Group("attendanceAgency")
	{
		attendanceAgencyRouter.POST("createAttendanceAgency", attendanceAgencyApi.CreateAttendanceAgency)   // 新建Đơn vị
		attendanceAgencyRouter.DELETE("deleteAttendanceAgency", attendanceAgencyApi.DeleteAttendanceAgency) // 删除Đơn vị
		attendanceAgencyRouter.DELETE("deleteAttendanceAgencyByIds", attendanceAgencyApi.DeleteAttendanceAgencyByIds) // 批量删除Đơn vị
		attendanceAgencyRouter.PUT("updateAttendanceAgency", attendanceAgencyApi.UpdateAttendanceAgency)    // 更新Đơn vị
	}
	{
		attendanceAgencyRouterWithoutRecord.GET("findAttendanceAgency", attendanceAgencyApi.FindAttendanceAgency)        // 根据ID获取Đơn vị
		attendanceAgencyRouterWithoutRecord.GET("getAttendanceAgencyList", attendanceAgencyApi.GetAttendanceAgencyList)  // 获取Đơn vị列表
	}
	{
	    attendanceAgencyRouterWithoutAuth.GET("getAttendanceAgencyPublic", attendanceAgencyApi.GetAttendanceAgencyPublic)  // 获取Đơn vị列表
	}
}
