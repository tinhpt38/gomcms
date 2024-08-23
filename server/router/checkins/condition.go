package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ConditionRouter struct {}

// InitConditionRouter 初始化 Điều kiện để checkins 路由信息
func (s *ConditionRouter) InitConditionRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	conditionRouter := Router.Group("condition").Use(middleware.OperationRecord())
	conditionRouterWithoutRecord := Router.Group("condition")
	conditionRouterWithoutAuth := PublicRouter.Group("condition")
	{
		conditionRouter.POST("createCondition", conditionApi.CreateCondition)   // 新建Điều kiện để checkins
		conditionRouter.DELETE("deleteCondition", conditionApi.DeleteCondition) // 删除Điều kiện để checkins
		conditionRouter.DELETE("deleteConditionByIds", conditionApi.DeleteConditionByIds) // 批量删除Điều kiện để checkins
		conditionRouter.PUT("updateCondition", conditionApi.UpdateCondition)    // 更新Điều kiện để checkins
	}
	{
		conditionRouterWithoutRecord.GET("findCondition", conditionApi.FindCondition)        // 根据ID获取Điều kiện để checkins
		conditionRouterWithoutRecord.GET("getConditionList", conditionApi.GetConditionList)  // 获取Điều kiện để checkins列表
	}
	{
	    conditionRouterWithoutAuth.GET("getConditionPublic", conditionApi.GetConditionPublic)  // 获取Điều kiện để checkins列表
	}
}
