package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type AreaRouter struct {}

// InitAreaRouter 初始化 Khu vực điểm danh 路由信息
func (s *AreaRouter) InitAreaRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	areaRouter := Router.Group("area").Use(middleware.OperationRecord())
	areaRouterWithoutRecord := Router.Group("area")
	areaRouterWithoutAuth := PublicRouter.Group("area")
	{
		areaRouter.POST("createArea", areaApi.CreateArea)   // 新建Khu vực điểm danh
		areaRouter.DELETE("deleteArea", areaApi.DeleteArea) // 删除Khu vực điểm danh
		areaRouter.DELETE("deleteAreaByIds", areaApi.DeleteAreaByIds) // 批量删除Khu vực điểm danh
		areaRouter.PUT("updateArea", areaApi.UpdateArea)    // 更新Khu vực điểm danh
	}
	{
		areaRouterWithoutRecord.GET("findArea", areaApi.FindArea)        // 根据ID获取Khu vực điểm danh
		areaRouterWithoutRecord.GET("getAreaList", areaApi.GetAreaList)  // 获取Khu vực điểm danh列表
	}
	{
	    areaRouterWithoutAuth.GET("getAreaPublic", areaApi.GetAreaPublic)  // 获取Khu vực điểm danh列表
	}
}
