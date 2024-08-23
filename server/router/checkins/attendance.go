package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type AttendanceRouter struct{}

// InitAttendanceRouter 初始化 AttendanceClass 路由信息
func (s *AttendanceRouter) InitAttendanceRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	attendanceRouter := Router.Group("attendance").Use(middleware.OperationRecord())
	attendanceRouterWithoutRecord := Router.Group("attendance")
	attendanceRouterWithoutAuth := PublicRouter.Group("attendance")
	{
		attendanceRouter.POST("createAttendance", attendanceApi.CreateAttendance)             // 新建AttendanceClass
		attendanceRouter.DELETE("deleteAttendance", attendanceApi.DeleteAttendance)           // 删除AttendanceClass
		attendanceRouter.DELETE("deleteAttendanceByIds", attendanceApi.DeleteAttendanceByIds) // 批量删除AttendanceClass
		attendanceRouter.PUT("updateAttendance", attendanceApi.UpdateAttendance)              // 更新AttendanceClass
	}
	{
		attendanceRouterWithoutRecord.GET("findAttendance", attendanceApi.FindAttendance)       // 根据ID获取AttendanceClass
		attendanceRouterWithoutRecord.GET("getAttendanceList", attendanceApi.GetAttendanceList) // 获取AttendanceClass列表
	}
	{
		attendanceRouterWithoutAuth.GET("getAttendancePublic", attendanceApi.GetAttendancePublic) // 获取AttendanceClass列表
	}
}
