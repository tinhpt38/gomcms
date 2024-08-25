package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
	checkinsReq "github.com/flipped-aurora/gin-vue-admin/server/model/checkins/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ParticipantApi struct{}

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
		response.FailWithMessage("thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("thành công", c)
}

func (participantApi *ParticipantApi) DeleteParticipant(c *gin.Context) {
	ID := c.Query("ID")
	err := participantService.DeleteParticipant(ID)
	if err != nil {
		global.GVA_LOG.Error("thất bại!", zap.Error(err))
		response.FailWithMessage("thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("thành công", c)
}

func (participantApi *ParticipantApi) DeleteParticipantByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	err := participantService.DeleteParticipantByIds(IDs)
	if err != nil {
		global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("Thành công", c)
}

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
		response.FailWithMessage("Thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("Thành công", c)
}

func (participantApi *ParticipantApi) FindParticipant(c *gin.Context) {
	ID := c.Query("ID")
	reparticipant, err := participantService.GetParticipant(ID)
	if err != nil {
		global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:"+err.Error(), c)
		return
	}
	response.OkWithData(reparticipant, c)
}

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

func (participantApi *ParticipantApi) GetParticipantListByAttendance(c *gin.Context) {
	var pageInfo checkinsReq.ParticipantSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := participantService.GetParticipantInfoListByAttendance(pageInfo)
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

func (participantApi *ParticipantApi) GetParticipantPublic(c *gin.Context) {

	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的Sinh viên (Người tham dự phiên điểm danh)接口信息",
	}, "Thành công", c)
}
