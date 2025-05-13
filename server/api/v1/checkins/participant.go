package checkins

import (
	"strconv"
	"time"

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

func (participantApi *ParticipantApi) BulkCreateParticipants(c *gin.Context) {
	var req checkinsReq.ListEmailParticipantRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = participantService.BulkCreateParticipants(req)
	if err != nil {
		global.GVA_LOG.Error("Thêm hàng loạt thất bại!", zap.Error(err))
		response.FailWithMessage("Thêm hàng loạt thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("Thêm hàng loạt thành công", c)
}

func (participantApi *ParticipantApi) DeleteParticipant(c *gin.Context) {
	ID := c.Query("ID")
	attIdSting := c.Query("attendanceId")
	attId, _ := strconv.Atoi(attIdSting)
	if attId > 0 {
		err := participantService.DeleteParticipantInAttendance(ID, uint(attId))
		if err != nil {
			global.GVA_LOG.Error("thất bại!", zap.Error(err))
			response.FailWithMessage("thất bại:"+err.Error(), c)
			return
		}

	} else {
		err := participantService.DeleteParticipant(ID)
		if err != nil {
			global.GVA_LOG.Error("thất bại!", zap.Error(err))
			response.FailWithMessage("thất bại:"+err.Error(), c)
			return
		}
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

func (participantApi *ParticipantApi) FindLuckyParticipant(c *gin.Context) {
	acId := c.Query("attendanceId")

	// Lấy người may mắn hiện tại
	participant, luckyNumber, err := participantService.GetLuckyParticipant(acId)

	// Lấy lịch sử các số may mắn đã quay
	var luckyHistory []checkins.UsedLuckyParticipant
	global.GVA_DB.Where("attendance_id = ?", acId).
		Order("created_at DESC").
		Limit(20).
		Find(&luckyHistory)

	// Lấy thông tin chi tiết của mỗi người tham gia trong lịch sử
	type LuckyHistoryItem struct {
		ID            uint      `json:"id"`
		LuckyNumber   *int      `json:"luckyNumber"`
		ParticipantId *uint     `json:"participantId"`
		Email         string    `json:"email"`
		FullName      string    `json:"fullName"`
		CreatedAt     time.Time `json:"createdAt"`
	}

	var historyItems []LuckyHistoryItem
	for _, hist := range luckyHistory {
		item := LuckyHistoryItem{
			ID:            hist.ID,
			LuckyNumber:   hist.LuckyNumber,
			ParticipantId: hist.ParticipantId,
			CreatedAt:     hist.CreatedAt,
		}

		// Nếu có participant_id, lấy thông tin người tham gia
		if hist.ParticipantId != nil && *hist.ParticipantId > 0 {
			var p checkins.Participant
			if global.GVA_DB.Where("id = ?", *hist.ParticipantId).First(&p).Error == nil {
				item.Email = p.Email
				if p.FullName != nil {
					item.FullName = *p.FullName
				}
			}
		}

		historyItems = append(historyItems, item)
	}

	if err != nil {
		global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:"+err.Error(), c)
		return
	}

	response.OkWithData(gin.H{
		"participant":  participant,
		"luckyNumber":  luckyNumber,
		"luckyHistory": historyItems,
	}, c)
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
		"info": "Check thành viên (Người tham dự phiên điểm danh)",
	}, "Thành công", c)
}
