package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type AttendanceRouter struct{}

func (s *AttendanceRouter) InitAttendanceRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	attendanceRouter := Router.Group("attendance").Use(middleware.OperationRecord())
	attendanceRouterWithoutRecord := Router.Group("attendance")
	attendanceRouterWithoutAuth := PublicRouter.Group("attendance")
	{
		attendanceRouter.POST("createAttendance", attendanceApi.CreateAttendance)
		attendanceRouter.POST("createAttendanceArea", attendanceApi.CreateAttendanceArea)
		attendanceRouter.POST("cloneAttendance", attendanceApi.CloneAttendance)
		attendanceRouter.POST("statsByAgencyCategory", attendanceApi.StatsByAgencyCategory)
		attendanceRouter.POST("statsScatterPlot", attendanceApi.StatsScatterPlot)
		attendanceRouter.POST("statsTrendLine", attendanceApi.StatsTrendLine)

		attendanceRouter.DELETE("deleteAttendance", attendanceApi.DeleteAttendance)
		attendanceRouter.DELETE("deleteAttendanceArea", attendanceApi.DeleteAttendanceArea)
		attendanceRouter.DELETE("deleteAttendanceByIds", attendanceApi.DeleteAttendanceByIds)
		attendanceRouter.PUT("updateAttendance", attendanceApi.UpdateAttendance)

	}
	{
		attendanceRouterWithoutRecord.GET("findAttendance", attendanceApi.FindAttendance)
		attendanceRouterWithoutRecord.GET("findAttendanceArea", attendanceApi.FindAttendanceArea)
		attendanceRouterWithoutRecord.GET("getAttendanceList", attendanceApi.GetAttendanceList)
	}
	{
		attendanceRouterWithoutAuth.GET("getAttendancePublic", attendanceApi.GetPublicAttendanceList)
	}
}
