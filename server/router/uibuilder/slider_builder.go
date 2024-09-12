package uibuilder

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type SliderBuilderRouter struct {}

// InitSliderBuilderRouter 初始化 SliderBuilder 路由信息
func (s *SliderBuilderRouter) InitSliderBuilderRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	sliderBuilderRouter := Router.Group("sliderBuilder").Use(middleware.OperationRecord())
	sliderBuilderRouterWithoutRecord := Router.Group("sliderBuilder")
	sliderBuilderRouterWithoutAuth := PublicRouter.Group("sliderBuilder")
	{
		sliderBuilderRouter.POST("createSliderBuilder", sliderBuilderApi.CreateSliderBuilder)   // 新建SliderBuilder
		sliderBuilderRouter.DELETE("deleteSliderBuilder", sliderBuilderApi.DeleteSliderBuilder) // 删除SliderBuilder
		sliderBuilderRouter.DELETE("deleteSliderBuilderByIds", sliderBuilderApi.DeleteSliderBuilderByIds) // 批量删除SliderBuilder
		sliderBuilderRouter.PUT("updateSliderBuilder", sliderBuilderApi.UpdateSliderBuilder)    // 更新SliderBuilder
	}
	{
		sliderBuilderRouterWithoutRecord.GET("findSliderBuilder", sliderBuilderApi.FindSliderBuilder)        // 根据ID获取SliderBuilder
		sliderBuilderRouterWithoutRecord.GET("getSliderBuilderList", sliderBuilderApi.GetSliderBuilderList)  // 获取SliderBuilder列表
	}
	{
	    sliderBuilderRouterWithoutAuth.GET("getSliderBuilderPublic", sliderBuilderApi.GetSliderBuilderPublic)  // 获取SliderBuilder列表
	}
}
