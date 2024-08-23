package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
    checkinsReq "github.com/flipped-aurora/gin-vue-admin/server/model/checkins/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type ParticipantApi struct {}

// CreateParticipant 创建Sinh viên (Người tham dự phiên điểm danh)
// @Tags Participant
// @Summary 创建Sinh viên (Người tham dự phiên điểm danh)
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body checkins.Participant true "创建Sinh viên (Người tham dự phiên điểm danh)"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /participant/createParticipant [post]
func (participantApi *ParticipantApi) CreateParticipant(c *gin.Context) {
	var participant checkins.Participant
	err := c.ShouldBindJSON(&participant)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = participantService.CreateParticipant(&participant)
	if err != nil {
        global.GVA_LOG.Error("thất bại!", zap.Error(err))
		response.FailWithMessage("thất bại:" + err.Error(), c)
		return
	}
    response.OkWithMessage("thành công", c)
}

// DeleteParticipant 删除Sinh viên (Người tham dự phiên điểm danh)
// @Tags Participant
// @Summary 删除Sinh viên (Người tham dự phiên điểm danh)
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body checkins.Participant true "删除Sinh viên (Người tham dự phiên điểm danh)"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /participant/deleteParticipant [delete]
func (participantApi *ParticipantApi) DeleteParticipant(c *gin.Context) {
	ID := c.Query("ID")
	err := participantService.DeleteParticipant(ID)
	if err != nil {
        global.GVA_LOG.Error("thất bại!", zap.Error(err))
		response.FailWithMessage("thất bại:" + err.Error(), c)
		return
	}
	response.OkWithMessage("thành công", c)
}

// DeleteParticipantByIds 批量删除Sinh viên (Người tham dự phiên điểm danh)
// @Tags Participant
// @Summary 批量删除Sinh viên (Người tham dự phiên điểm danh)
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /participant/deleteParticipantByIds [delete]
func (participantApi *ParticipantApi) DeleteParticipantByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	err := participantService.DeleteParticipantByIds(IDs)
	if err != nil {
        global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:" + err.Error(), c)
		return
	}
	response.OkWithMessage("Thành công", c)
}

// UpdateParticipant 更新Sinh viên (Người tham dự phiên điểm danh)
// @Tags Participant
// @Summary 更新Sinh viên (Người tham dự phiên điểm danh)
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body checkins.Participant true "更新Sinh viên (Người tham dự phiên điểm danh)"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /participant/updateParticipant [put]
func (participantApi *ParticipantApi) UpdateParticipant(c *gin.Context) {
	var participant checkins.Participant
	err := c.ShouldBindJSON(&participant)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = participantService.UpdateParticipant(participant)
	if err != nil {
        global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:" + err.Error(), c)
		return
	}
	response.OkWithMessage("Thành công", c)
}

// FindParticipant 用id查询Sinh viên (Người tham dự phiên điểm danh)
// @Tags Participant
// @Summary 用id查询Sinh viên (Người tham dự phiên điểm danh)
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query checkins.Participant true "用id查询Sinh viên (Người tham dự phiên điểm danh)"
// @Success 200 {object} response.Response{data=checkins.Participant,msg=string} "查询成功"
// @Router /participant/findParticipant [get]
func (participantApi *ParticipantApi) FindParticipant(c *gin.Context) {
	ID := c.Query("ID")
	reparticipant, err := participantService.GetParticipant(ID)
	if err != nil {
        global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:" + err.Error(), c)
		return
	}
	response.OkWithData(reparticipant, c)
}

// GetParticipantList 分页获取Sinh viên (Người tham dự phiên điểm danh)列表
// @Tags Participant
// @Summary 分页获取Sinh viên (Người tham dự phiên điểm danh)列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query checkinsReq.ParticipantSearch true "分页获取Sinh viên (Người tham dự phiên điểm danh)列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /participant/getParticipantList [get]
func (participantApi *ParticipantApi) GetParticipantList(c *gin.Context) {
	var pageInfo checkinsReq.ParticipantSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := participantService.GetParticipantInfoList(pageInfo)
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

// GetParticipantPublic 不需要鉴权的Sinh viên (Người tham dự phiên điểm danh)接口
// @Tags Participant
// @Summary 不需要鉴权的Sinh viên (Người tham dự phiên điểm danh)接口
// @accept application/json
// @Produce application/json
// @Param data query checkinsReq.ParticipantSearch true "分页获取Sinh viên (Người tham dự phiên điểm danh)列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /participant/getParticipantPublic [get]
func (participantApi *ParticipantApi) GetParticipantPublic(c *gin.Context) {
    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的Sinh viên (Người tham dự phiên điểm danh)接口信息",
    }, "Thành công", c)
}
