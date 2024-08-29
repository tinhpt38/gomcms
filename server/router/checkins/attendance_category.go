package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type AttendanceCategoryRouter struct {}

// InitAttendanceCategoryRouter 初始化 Danh mục điểm danh 路由信息
func (s *AttendanceCategoryRouter) InitAttendanceCategoryRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	attendanceCategoryRouter := Router.Group("attendanceCategory").Use(middleware.OperationRecord())
	attendanceCategoryRouterWithoutRecord := Router.Group("attendanceCategory")
	attendanceCategoryRouterWithoutAuth := PublicRouter.Group("attendanceCategory")
	{
		attendanceCategoryRouter.POST("createAttendanceCategory", attendanceCategoryApi.CreateAttendanceCategory)   // 新建Danh mục điểm danh
		attendanceCategoryRouter.DELETE("deleteAttendanceCategory", attendanceCategoryApi.DeleteAttendanceCategory) // 删除Danh mục điểm danh
		attendanceCategoryRouter.DELETE("deleteAttendanceCategoryByIds", attendanceCategoryApi.DeleteAttendanceCategoryByIds) // 批量删除Danh mục điểm danh
		attendanceCategoryRouter.PUT("updateAttendanceCategory", attendanceCategoryApi.UpdateAttendanceCategory)    // 更新Danh mục điểm danh
	}
	{
		attendanceCategoryRouterWithoutRecord.GET("findAttendanceCategory", attendanceCategoryApi.FindAttendanceCategory)        // 根据ID获取Danh mục điểm danh
		attendanceCategoryRouterWithoutRecord.GET("getAttendanceCategoryList", attendanceCategoryApi.GetAttendanceCategoryList)  // 获取Danh mục điểm danh列表
	}
	{
	    attendanceCategoryRouterWithoutAuth.GET("getAttendanceCategoryPublic", attendanceCategoryApi.GetAttendanceCategoryPublic)  // 获取Danh mục điểm danh列表
	}
}
