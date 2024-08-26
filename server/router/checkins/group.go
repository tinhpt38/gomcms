package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type GroupRouter struct{}

// InitGroupRouter 初始化 Nhóm 路由信息
func (s *GroupRouter) InitGroupRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	groupRouter := Router.Group("group").Use(middleware.OperationRecord())
	groupRouterWithoutRecord := Router.Group("group")
	groupRouterWithoutAuth := PublicRouter.Group("group")
	{
		groupRouter.POST("createGroup", groupApi.CreateGroup)                                   // 新建Nhóm
		groupRouter.DELETE("deleteGroup", groupApi.DeleteGroup)                                 // 删除Nhóm
		groupRouter.DELETE("deleteGroupByIds", groupApi.DeleteGroupByIds)                       // 批量删除Nhóm
		groupRouter.PUT("updateGroup", groupApi.UpdateGroup)                                    // 更新Nhóm
		groupRouter.POST("assignParticipantToGroupAuto", groupApi.AssignParticipantToGroupAuto) // Tạo nhóm và gắn nhóm tự động cho người tham gia
	}
	{
		groupRouterWithoutRecord.GET("findGroup", groupApi.FindGroup)       // 根据ID获取Nhóm
		groupRouterWithoutRecord.GET("getGroupList", groupApi.GetGroupList) // 获取Nhóm列表
	}
	{
		groupRouterWithoutAuth.GET("getGroupDataSource", groupApi.GetGroupDataSource) // 获取Nhóm数据源
		groupRouterWithoutAuth.GET("getGroupPublic", groupApi.GetGroupPublic)         // 获取Nhóm列表
	}
}
