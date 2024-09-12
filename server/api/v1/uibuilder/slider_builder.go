package uibuilder

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/uibuilder"
	uibuilderReq "github.com/flipped-aurora/gin-vue-admin/server/model/uibuilder/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SliderBuilderApi struct{}

// CreateSliderBuilder 创建SliderBuilder
// @Tags SliderBuilder
// @Summary 创建SliderBuilder
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body uibuilder.SliderBuilder true "创建SliderBuilder"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /sliderBuilder/createSliderBuilder [post]
func (sliderBuilderApi *SliderBuilderApi) CreateSliderBuilder(c *gin.Context) {
	var sliderBuilder uibuilder.SliderBuilder
	err := c.ShouldBindJSON(&sliderBuilder)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	sliderBuilder.CreatedBy = utils.GetUserID(c)
	err = sliderBuilderService.CreateSliderBuilder(&sliderBuilder)
	if err != nil {
		global.GVA_LOG.Error("thất bại!", zap.Error(err))
		response.FailWithMessage("thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("thành công", c)
}

// DeleteSliderBuilder 删除SliderBuilder
// @Tags SliderBuilder
// @Summary 删除SliderBuilder
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body uibuilder.SliderBuilder true "删除SliderBuilder"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /sliderBuilder/deleteSliderBuilder [delete]
func (sliderBuilderApi *SliderBuilderApi) DeleteSliderBuilder(c *gin.Context) {
	ID := c.Query("ID")
	userID := utils.GetUserID(c)
	err := sliderBuilderService.DeleteSliderBuilder(ID, userID)
	if err != nil {
		global.GVA_LOG.Error("thất bại!", zap.Error(err))
		response.FailWithMessage("thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("thành công", c)
}

// DeleteSliderBuilderByIds 批量删除SliderBuilder
// @Tags SliderBuilder
// @Summary 批量删除SliderBuilder
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /sliderBuilder/deleteSliderBuilderByIds [delete]
func (sliderBuilderApi *SliderBuilderApi) DeleteSliderBuilderByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	userID := utils.GetUserID(c)
	err := sliderBuilderService.DeleteSliderBuilderByIds(IDs, userID)
	if err != nil {
		global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("Thành công", c)
}

// UpdateSliderBuilder 更新SliderBuilder
// @Tags SliderBuilder
// @Summary 更新SliderBuilder
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body uibuilder.SliderBuilder true "更新SliderBuilder"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /sliderBuilder/updateSliderBuilder [put]
func (sliderBuilderApi *SliderBuilderApi) UpdateSliderBuilder(c *gin.Context) {
	var sliderBuilder uibuilder.SliderBuilder
	err := c.ShouldBindJSON(&sliderBuilder)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	sliderBuilder.UpdatedBy = utils.GetUserID(c)
	err = sliderBuilderService.UpdateSliderBuilder(sliderBuilder)
	if err != nil {
		global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("Thành công", c)
}

// FindSliderBuilder 用id查询SliderBuilder
// @Tags SliderBuilder
// @Summary 用id查询SliderBuilder
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query uibuilder.SliderBuilder true "用id查询SliderBuilder"
// @Success 200 {object} response.Response{data=uibuilder.SliderBuilder,msg=string} "查询成功"
// @Router /sliderBuilder/findSliderBuilder [get]
func (sliderBuilderApi *SliderBuilderApi) FindSliderBuilder(c *gin.Context) {
	ID := c.Query("ID")
	resliderBuilder, err := sliderBuilderService.GetSliderBuilder(ID)
	if err != nil {
		global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:"+err.Error(), c)
		return
	}
	response.OkWithData(resliderBuilder, c)
}

// GetSliderBuilderList 分页获取SliderBuilder列表
// @Tags SliderBuilder
// @Summary 分页获取SliderBuilder列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query uibuilderReq.SliderBuilderSearch true "分页获取SliderBuilder列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /sliderBuilder/getSliderBuilderList [get]
func (sliderBuilderApi *SliderBuilderApi) GetSliderBuilderList(c *gin.Context) {
	var pageInfo uibuilderReq.SliderBuilderSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := sliderBuilderService.GetSliderBuilderInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "Thành công", c)
}

// @Router /sliderBuilder/getSliderBuilderPublic [get]
func (sliderBuilderApi *SliderBuilderApi) GetSliderBuilderPublic(c *gin.Context) {
	slider, err := sliderBuilderService.GetPublicSlider()
	if err != nil {
		global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại: "+err.Error(), c)
		return
	}
	response.OkWithData(slider, c)
}
