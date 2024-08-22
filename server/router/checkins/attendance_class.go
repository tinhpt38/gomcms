package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type AttendanceClassRouter struct {}

// InitAttendanceClassRouter 初始化 AttendanceClass 路由信息
func (s *AttendanceClassRouter) InitAttendanceClassRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	attendanceClassRouter := Router.Group("attendanceClass").Use(middleware.OperationRecord())
	attendanceClassRouterWithoutRecord := Router.Group("attendanceClass")
	attendanceClassRouterWithoutAuth := PublicRouter.Group("attendanceClass")
	{
		attendanceClassRouter.POST("createAttendanceClass", attendanceClassApi.CreateAttendanceClass)   // 新建AttendanceClass
		attendanceClassRouter.DELETE("deleteAttendanceClass", attendanceClassApi.DeleteAttendanceClass) // 删除AttendanceClass
		attendanceClassRouter.DELETE("deleteAttendanceClassByIds", attendanceClassApi.DeleteAttendanceClassByIds) // 批量删除AttendanceClass
		attendanceClassRouter.PUT("updateAttendanceClass", attendanceClassApi.UpdateAttendanceClass)    // 更新AttendanceClass
	}
	{
		attendanceClassRouterWithoutRecord.GET("findAttendanceClass", attendanceClassApi.FindAttendanceClass)        // 根据ID获取AttendanceClass
		attendanceClassRouterWithoutRecord.GET("getAttendanceClassList", attendanceClassApi.GetAttendanceClassList)  // 获取AttendanceClass列表
	}
	{
	    attendanceClassRouterWithoutAuth.GET("getAttendanceClassPublic", attendanceClassApi.GetAttendanceClassPublic)  // 获取AttendanceClass列表
	}
}
