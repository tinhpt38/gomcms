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

type AttendanceApi struct{}

func (attendanceApi *AttendanceApi) CreateAttendance(c *gin.Context) {
	var attendance checkins.Attendance
	err := c.ShouldBindJSON(&attendance)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	attendance.CreatedBy = utils.GetUserID(c)
	err = attendanceService.CreateAttendance(&attendance)
	if err != nil {
		global.GVA_LOG.Error("tạo thất bại!", zap.Error(err))
		response.FailWithMessage("tạo thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("tạo thành công", c)
}

func (attendanceApi *AttendanceApi) CreateAttendanceArea(c *gin.Context) {
	var attendanceArea checkins.AttendanceArea
	err := c.ShouldBindJSON(&attendanceArea)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = attendanceService.CreateAttendanceArea(&attendanceArea)
	if err != nil {
		global.GVA_LOG.Error("tạo thất bại!", zap.Error(err))
		response.FailWithMessage("tạo thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("tạo thành công", c)
}

func (attendanceApi *AttendanceApi) DeleteAttendance(c *gin.Context) {
	ID := c.Query("ID")
	userID := utils.GetUserID(c)
	err := attendanceService.DeleteAttendance(ID, userID)
	if err != nil {
		global.GVA_LOG.Error("xóa thất bại!", zap.Error(err))
		response.FailWithMessage("xóa thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("xóa thành công", c)
}

func (attendanceApi *AttendanceApi) DeleteAttendanceArea(c *gin.Context) {
	id := c.Query("id")
	err := attendanceService.DeleteAttendanceArea(id)
	if err != nil {
		global.GVA_LOG.Error("xóa thất bại!", zap.Error(err))
		response.FailWithMessage("xóa thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("xóa thành công", c)
}

// DeleteAttendanceByIds xóa nhiều AttendanceClass
// @Tags AttendanceClass
// @Summary xóa nhiều AttendanceClass
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "xóa thành công"
// @Router /attendance/deleteAttendanceByIds [delete]
func (attendanceApi *AttendanceApi) DeleteAttendanceByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	userID := utils.GetUserID(c)
	err := attendanceService.DeleteAttendanceByIds(IDs, userID)
	if err != nil {
		global.GVA_LOG.Error("xóa thất bại!", zap.Error(err))
		response.FailWithMessage("xóa thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("xóa thành công", c)
}

// UpdateAttendance cập nhật AttendanceClass
// @Tags AttendanceClass
// @Summary cập nhật AttendanceClass
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body checkins.AttendanceClass true "cập nhật AttendanceClass"
// @Success 200 {object} response.Response{msg=string} "cập nhật thành công"
// @Router /attendance/updateAttendance [put]
func (attendanceApi *AttendanceApi) UpdateAttendance(c *gin.Context) {
	var attendance checkins.Attendance
	err := c.ShouldBindJSON(&attendance)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	attendance.UpdatedBy = utils.GetUserID(c)
	err = attendanceService.UpdateAttendance(attendance)
	if err != nil {
		global.GVA_LOG.Error("cập nhật thất bại!", zap.Error(err))
		response.FailWithMessage("cập nhật thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("cập nhật thành công", c)
}

func (attendanceApi *AttendanceApi) FindAttendance(c *gin.Context) {
	ID := c.Query("id")
	reattendance, err := attendanceService.GetAttendance(ID)
	if err != nil {
		global.GVA_LOG.Error("tìm kiếm thất bại!", zap.Error(err))
		response.FailWithMessage("tìm kiếm thất bại:"+err.Error(), c)
		return
	}
	response.OkWithData(reattendance, c)
}

func (attendanceApi *AttendanceApi) FindAttendanceArea(c *gin.Context) {
	id := c.Query("id")
	areas, err := attendanceService.GetAttendanceArea(id)
	if err != nil {
		global.GVA_LOG.Error("tìm kiếm thất bại!", zap.Error(err))
		response.FailWithMessage("tìm kiếm thất bại:"+err.Error(), c)
		return
	}
	response.OkWithData(areas, c)
}

func (attendanceApi *AttendanceApi) GetAttendanceList(c *gin.Context) {
	var pageInfo checkinsReq.AttendanceSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	craeteBy := utils.GetUserID(c)
	list, total, err := attendanceService.GetAttendanceInfoList(pageInfo, craeteBy)
	if err != nil {
		global.GVA_LOG.Error("lấy thất bại!", zap.Error(err))
		response.FailWithMessage("lấy thất bại:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "lấy thành công", c)
}

func (attendanceApi *AttendanceApi) GetAttendancePublic(c *gin.Context) {
	// API này không cần xác thực
	// Ví dụ trả về một thông điệp cố định, thường được sử dụng cho dịch vụ phía C, cần triển khai logic kinh doanh của riêng mình
	response.OkWithDetailed(gin.H{
		"info": "Thông tin các API AttendanceClass không cần xác thực",
	}, "lấy thành công", c)
}

func (attendanceApi *AttendanceApi) CloneAttendance(c *gin.Context) {
	var clone checkinsReq.AttendanceCloneRequest
	err := c.ShouldBindJSON(&clone)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	createdBy := utils.GetUserID(c)
	err = attendanceService.CloneAttendance(clone, createdBy)
	if err != nil {
		global.GVA_LOG.Error("Nhân bản thất bại!", zap.Error(err))
		response.FailWithMessage("Nhân bản thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("Nhân bản thành công", c)
}

func (attendanceApi *AttendanceApi) StatsByAgencyCategory(c *gin.Context) {
	var stats checkinsReq.StatsByAgencyCategory
	err := c.ShouldBindJSON(&stats)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	res, err := attendanceService.StatsByAgencyCategory(stats)
	if err != nil {
		global.GVA_LOG.Error("Không lấy được dữ liệu!", zap.Error(err))
		response.FailWithMessage("Không lấy được dữ liệu:"+err.Error(), c)
		return
	}
	response.OkWithData(res, c)
}
