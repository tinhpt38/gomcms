package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type AttendanceCheckInRouter struct{}

// InitAttendanceCheckInRouter 初始化 Thành viên checkins 路由信息
func (s *AttendanceCheckInRouter) InitAttendanceCheckInRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	attendanceCheckInRouter := Router.Group("attendanceCheckIn").Use(middleware.OperationRecord())
	attendanceCheckInRouterWithoutRecord := Router.Group("attendanceCheckIn")
	attendanceCheckInRouterWithoutAuth := PublicRouter.Group("attendanceCheckIn")
	{
		attendanceCheckInRouter.POST("createAttendanceCheckIn", attendanceCheckInApi.CreateAttendanceCheckIn)             // 新建Thành viên checkins
		attendanceCheckInRouter.DELETE("deleteAttendanceCheckIn", attendanceCheckInApi.DeleteAttendanceCheckIn)           // 删除Thành viên checkins
		attendanceCheckInRouter.DELETE("deleteAttendanceCheckInByIds", attendanceCheckInApi.DeleteAttendanceCheckInByIds) // 批量删除Thành viên checkins
		attendanceCheckInRouter.PUT("updateAttendanceCheckIn", attendanceCheckInApi.UpdateAttendanceCheckIn)              // 更新Thành viên checkins
	}
	{
		attendanceCheckInRouterWithoutRecord.GET("findAttendanceCheckIn", attendanceCheckInApi.FindAttendanceCheckIn)             // 根据ID获取Thành viên checkins
		attendanceCheckInRouterWithoutRecord.GET("getAttendanceCheckInList", attendanceCheckInApi.GetAttendanceCheckInList)       // 获取Thành viên checkins列表
		attendanceCheckInRouterWithoutRecord.GET("getAttendanceCheckInLogList", attendanceCheckInApi.GetAttendanceCheckInLogList) // 获取Thành viên checkins列表

	}
	{
		attendanceCheckInRouterWithoutAuth.GET("getAttendanceCheckInPublic", attendanceCheckInApi.GetAttendanceCheckInPublic)
		attendanceCheckInRouterWithoutAuth.POST("publicAttendanceCheckIn", attendanceCheckInApi.CheckinAttendance)
	}
}
