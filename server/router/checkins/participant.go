package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ParticipantRouter struct{}

// InitParticipantRouter 初始化 Sinh viên (Người tham dự phiên điểm danh) 路由信息
func (s *ParticipantRouter) InitParticipantRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	participantRouter := Router.Group("participant").Use(middleware.OperationRecord())
	participantRouterWithoutRecord := Router.Group("participant")
	participantRouterWithoutAuth := PublicRouter.Group("participant")
	{
		participantRouter.POST("createParticipant", participantApi.CreateParticipant)             // 新建Sinh viên (Người tham dự phiên điểm danh)
		participantRouter.DELETE("deleteParticipant", participantApi.DeleteParticipant)           // 删除Sinh viên (Người tham dự phiên điểm danh)
		participantRouter.DELETE("deleteParticipantByIds", participantApi.DeleteParticipantByIds) // 批量删除Sinh viên (Người tham dự phiên điểm danh)
		participantRouter.PUT("updateParticipant", participantApi.UpdateParticipant)              // 更新Sinh viên (Người tham dự phiên điểm danh)
	}
	{
		participantRouterWithoutRecord.GET("findParticipant", participantApi.FindParticipant)       // 根据ID获取Sinh viên (Người tham dự phiên điểm danh)
		participantRouterWithoutRecord.GET("getParticipantList", participantApi.GetParticipantList) // 获取Sinh viên (Người tham dự phiên điểm danh)列表
		participantRouterWithoutRecord.GET("getParticipantListByAttendance", participantApi.GetParticipantListByAttendance)
	}
	{
		participantRouterWithoutAuth.GET("getParticipantPublic", participantApi.GetParticipantPublic) // 获取Sinh viên (Người tham dự phiên điểm danh)列表
	}
}
