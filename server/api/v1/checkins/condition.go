package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
    checkinsReq "github.com/flipped-aurora/gin-vue-admin/server/model/checkins/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "github.com/flipped-aurora/gin-vue-admin/server/utils"
)

type ConditionApi struct {}

// CreateCondition 创建Điều kiện để checkins
// @Tags Condition
// @Summary 创建Điều kiện để checkins
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body checkins.Condition true "创建Điều kiện để checkins"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /condition/createCondition [post]
func (conditionApi *ConditionApi) CreateCondition(c *gin.Context) {
	var condition checkins.Condition
	err := c.ShouldBindJSON(&condition)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    condition.CreatedBy = utils.GetUserID(c)
	err = conditionService.CreateCondition(&condition)
	if err != nil {
        global.GVA_LOG.Error("thất bại!", zap.Error(err))
		response.FailWithMessage("thất bại:" + err.Error(), c)
		return
	}
    response.OkWithMessage("thành công", c)
}

// DeleteCondition 删除Điều kiện để checkins
// @Tags Condition
// @Summary 删除Điều kiện để checkins
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body checkins.Condition true "删除Điều kiện để checkins"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /condition/deleteCondition [delete]
func (conditionApi *ConditionApi) DeleteCondition(c *gin.Context) {
	ID := c.Query("ID")
    userID := utils.GetUserID(c)
	err := conditionService.DeleteCondition(ID,userID)
	if err != nil {
        global.GVA_LOG.Error("thất bại!", zap.Error(err))
		response.FailWithMessage("thất bại:" + err.Error(), c)
		return
	}
	response.OkWithMessage("thành công", c)
}

// DeleteConditionByIds 批量删除Điều kiện để checkins
// @Tags Condition
// @Summary 批量删除Điều kiện để checkins
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /condition/deleteConditionByIds [delete]
func (conditionApi *ConditionApi) DeleteConditionByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
    userID := utils.GetUserID(c)
	err := conditionService.DeleteConditionByIds(IDs,userID)
	if err != nil {
        global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:" + err.Error(), c)
		return
	}
	response.OkWithMessage("Thành công", c)
}

// UpdateCondition 更新Điều kiện để checkins
// @Tags Condition
// @Summary 更新Điều kiện để checkins
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body checkins.Condition true "更新Điều kiện để checkins"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /condition/updateCondition [put]
func (conditionApi *ConditionApi) UpdateCondition(c *gin.Context) {
	var condition checkins.Condition
	err := c.ShouldBindJSON(&condition)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    condition.UpdatedBy = utils.GetUserID(c)
	err = conditionService.UpdateCondition(condition)
	if err != nil {
        global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:" + err.Error(), c)
		return
	}
	response.OkWithMessage("Thành công", c)
}

// FindCondition 用id查询Điều kiện để checkins
// @Tags Condition
// @Summary 用id查询Điều kiện để checkins
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query checkins.Condition true "用id查询Điều kiện để checkins"
// @Success 200 {object} response.Response{data=checkins.Condition,msg=string} "查询成功"
// @Router /condition/findCondition [get]
func (conditionApi *ConditionApi) FindCondition(c *gin.Context) {
	ID := c.Query("ID")
	recondition, err := conditionService.GetCondition(ID)
	if err != nil {
        global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:" + err.Error(), c)
		return
	}
	response.OkWithData(recondition, c)
}

// GetConditionList 分页获取Điều kiện để checkins列表
// @Tags Condition
// @Summary 分页获取Điều kiện để checkins列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query checkinsReq.ConditionSearch true "分页获取Điều kiện để checkins列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /condition/getConditionList [get]
func (conditionApi *ConditionApi) GetConditionList(c *gin.Context) {
	var pageInfo checkinsReq.ConditionSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := conditionService.GetConditionInfoList(pageInfo)
	if err != nil {
	    global.GVA_LOG.Error("Thất bại!", zap.Error(err))
        response.FailWithMessage("Thất bại:" + err.Error(), c)
        return
    }
    response.OkWithDetailed(response.PageResult{
        List:     list,
        Total:    total,
        Page:     pageInfo.Page,
        PageSize: pageInfo.PageSize,
    }, "Thành công", c)
}

// GetConditionPublic 不需要鉴权的Điều kiện để checkins接口
// @Tags Condition
// @Summary 不需要鉴权的Điều kiện để checkins接口
// @accept application/json
// @Produce application/json
// @Param data query checkinsReq.ConditionSearch true "分页获取Điều kiện để checkins列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /condition/getConditionPublic [get]
func (conditionApi *ConditionApi) GetConditionPublic(c *gin.Context) {
    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的Điều kiện để checkins接口信息",
    }, "Thành công", c)
}
