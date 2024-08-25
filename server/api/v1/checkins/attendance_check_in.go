package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
	checkinsReq "github.com/flipped-aurora/gin-vue-admin/server/model/checkins/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AttendanceCheckInApi struct{}

func (attendanceCheckInApi *AttendanceCheckInApi) CreateAttendanceCheckIn(c *gin.Context) {
	var attendanceCheckIn checkins.AttendanceCheckIn
	err := c.ShouldBindJSON(&attendanceCheckIn)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	attendanceCheckIn.CreatedBy = utils.GetUserID(c)
	err = attendanceCheckInService.CreateAttendanceCheckIn(&attendanceCheckIn)
	if err != nil {
		global.GVA_LOG.Error("thất bại!", zap.Error(err))
		response.FailWithMessage("thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("thành công", c)
}

func (attendanceCheckInApi *AttendanceCheckInApi) DeleteAttendanceCheckIn(c *gin.Context) {
	ID := c.Query("ID")
	userID := utils.GetUserID(c)
	err := attendanceCheckInService.DeleteAttendanceCheckIn(ID, userID)
	if err != nil {
		global.GVA_LOG.Error("thất bại!", zap.Error(err))
		response.FailWithMessage("thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("thành công", c)
}

func (attendanceCheckInApi *AttendanceCheckInApi) DeleteAttendanceCheckInByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	userID := utils.GetUserID(c)
	err := attendanceCheckInService.DeleteAttendanceCheckInByIds(IDs, userID)
	if err != nil {
		global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("Thành công", c)
}

func (attendanceCheckInApi *AttendanceCheckInApi) UpdateAttendanceCheckIn(c *gin.Context) {
	var attendanceCheckIn checkins.AttendanceCheckIn
	err := c.ShouldBindJSON(&attendanceCheckIn)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	attendanceCheckIn.UpdatedBy = utils.GetUserID(c)
	err = attendanceCheckInService.UpdateAttendanceCheckIn(attendanceCheckIn)
	if err != nil {
		global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("Thành công", c)
}

func (attendanceCheckInApi *AttendanceCheckInApi) FindAttendanceCheckIn(c *gin.Context) {
	ID := c.Query("ID")
	reattendanceCheckIn, err := attendanceCheckInService.GetAttendanceCheckIn(ID)
	if err != nil {
		global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:"+err.Error(), c)
		return
	}
	response.OkWithData(reattendanceCheckIn, c)
}

func (attendanceCheckInApi *AttendanceCheckInApi) GetAttendanceCheckInList(c *gin.Context) {
	var pageInfo checkinsReq.AttendanceCheckInSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := attendanceCheckInService.GetAttendanceCheckInInfoList(pageInfo)
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

func (attendanceCheckInApi *AttendanceCheckInApi) GetAttendanceCheckInPublic(c *gin.Context) {
	response.OkWithDetailed(gin.H{
		"info": "Thông tin API checkins của thành viên không cần xác thực",
	}, "Thành công", c)
}

func (attendanceCheckInApi *AttendanceCheckInApi) CheckinAttendance(c *gin.Context) {
	var checkinReq checkinsReq.CheckinsReq
	err := c.ShouldBindJSON(&checkinReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	result, err := attendanceCheckInService.CheckinAttendance(checkinReq, ip, userAgent)
	if err != nil {
		global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:"+err.Error(), c)
		return
	}
	response.OkWithData(result, c)

}
