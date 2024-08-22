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

type AttendanceClassApi struct {}

// CreateAttendanceClass tạo AttendanceClass
// @Tags AttendanceClass
// @Summary tạo AttendanceClass
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body checkins.AttendanceClass true "tạo AttendanceClass"
// @Success 200 {object} response.Response{msg=string} "tạo thành công"
// @Router /attendanceClass/createAttendanceClass [post]
func (attendanceClassApi *AttendanceClassApi) CreateAttendanceClass(c *gin.Context) {
	var attendanceClass checkins.AttendanceClass
	err := c.ShouldBindJSON(&attendanceClass)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	attendanceClass.CreatedBy = utils.GetUserID(c)
	err = attendanceClassService.CreateAttendanceClass(&attendanceClass)
	if err != nil {
		global.GVA_LOG.Error("tạo thất bại!", zap.Error(err))
		response.FailWithMessage("tạo thất bại:" + err.Error(), c)
		return
	}
	response.OkWithMessage("tạo thành công", c)
}

// DeleteAttendanceClass xóa AttendanceClass
// @Tags AttendanceClass
// @Summary xóa AttendanceClass
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body checkins.AttendanceClass true "xóa AttendanceClass"
// @Success 200 {object} response.Response{msg=string} "xóa thành công"
// @Router /attendanceClass/deleteAttendanceClass [delete]
func (attendanceClassApi *AttendanceClassApi) DeleteAttendanceClass(c *gin.Context) {
	ID := c.Query("ID")
	userID := utils.GetUserID(c)
	err := attendanceClassService.DeleteAttendanceClass(ID,userID)
	if err != nil {
		global.GVA_LOG.Error("xóa thất bại!", zap.Error(err))
		response.FailWithMessage("xóa thất bại:" + err.Error(), c)
		return
	}
	response.OkWithMessage("xóa thành công", c)
}

// DeleteAttendanceClassByIds xóa nhiều AttendanceClass
// @Tags AttendanceClass
// @Summary xóa nhiều AttendanceClass
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "xóa thành công"
// @Router /attendanceClass/deleteAttendanceClassByIds [delete]
func (attendanceClassApi *AttendanceClassApi) DeleteAttendanceClassByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	userID := utils.GetUserID(c)
	err := attendanceClassService.DeleteAttendanceClassByIds(IDs,userID)
	if err != nil {
		global.GVA_LOG.Error("xóa thất bại!", zap.Error(err))
		response.FailWithMessage("xóa thất bại:" + err.Error(), c)
		return
	}
	response.OkWithMessage("xóa thành công", c)
}

// UpdateAttendanceClass cập nhật AttendanceClass
// @Tags AttendanceClass
// @Summary cập nhật AttendanceClass
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body checkins.AttendanceClass true "cập nhật AttendanceClass"
// @Success 200 {object} response.Response{msg=string} "cập nhật thành công"
// @Router /attendanceClass/updateAttendanceClass [put]
func (attendanceClassApi *AttendanceClassApi) UpdateAttendanceClass(c *gin.Context) {
	var attendanceClass checkins.AttendanceClass
	err := c.ShouldBindJSON(&attendanceClass)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	attendanceClass.UpdatedBy = utils.GetUserID(c)
	err = attendanceClassService.UpdateAttendanceClass(attendanceClass)
	if err != nil {
		global.GVA_LOG.Error("cập nhật thất bại!", zap.Error(err))
		response.FailWithMessage("cập nhật thất bại:" + err.Error(), c)
		return
	}
	response.OkWithMessage("cập nhật thành công", c)
}

// FindAttendanceClass tìm kiếm AttendanceClass theo ID
// @Tags AttendanceClass
// @Summary tìm kiếm AttendanceClass theo ID
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query checkins.AttendanceClass true "tìm kiếm AttendanceClass theo ID"
// @Success 200 {object} response.Response{data=checkins.AttendanceClass,msg=string} "tìm kiếm thành công"
// @Router /attendanceClass/findAttendanceClass [get]
func (attendanceClassApi *AttendanceClassApi) FindAttendanceClass(c *gin.Context) {
	ID := c.Query("id")
	reattendanceClass, err := attendanceClassService.GetAttendanceClass(ID)
	if err != nil {
		global.GVA_LOG.Error("tìm kiếm thất bại!", zap.Error(err))
		response.FailWithMessage("tìm kiếm thất bại:" + err.Error(), c)
		return
	}
	response.OkWithData(reattendanceClass, c)
}

// GetAttendanceClassList lấy danh sách AttendanceClass theo trang
// @Tags AttendanceClass
// @Summary lấy danh sách AttendanceClass theo trang
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query checkinsReq.AttendanceClassSearch true "lấy danh sách AttendanceClass theo trang"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "lấy thành công"
// @Router /attendanceClass/getAttendanceClassList [get]
func (attendanceClassApi *AttendanceClassApi) GetAttendanceClassList(c *gin.Context) {
	var pageInfo checkinsReq.AttendanceClassSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := attendanceClassService.GetAttendanceClassInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("lấy thất bại!", zap.Error(err))
		response.FailWithMessage("lấy thất bại:" + err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "lấy thành công", c)
}

// GetAttendanceClassPublic các API AttendanceClass không cần xác thực
// @Tags AttendanceClass
// @Summary các API AttendanceClass không cần xác thực
// @accept application/json
// @Produce application/json
// @Param data query checkinsReq.AttendanceClassSearch true "lấy danh sách AttendanceClass theo trang"
// @Success 200 {object} response.Response{data=object,msg=string} "lấy thành công"
// @Router /attendanceClass/getAttendanceClassPublic [get]
func (attendanceClassApi *AttendanceClassApi) GetAttendanceClassPublic(c *gin.Context) {
	// API này không cần xác thực
	// Ví dụ trả về một thông điệp cố định, thường được sử dụng cho dịch vụ phía C, cần triển khai logic kinh doanh của riêng mình
	response.OkWithDetailed(gin.H{
	   "info": "Thông tin các API AttendanceClass không cần xác thực",
	}, "lấy thành công", c)
}
