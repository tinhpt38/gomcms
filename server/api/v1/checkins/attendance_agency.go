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

type AttendanceAgencyApi struct {}

// CreateAttendanceAgency 创建Đơn vị
// @Tags AttendanceAgency
// @Summary 创建Đơn vị
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body checkins.AttendanceAgency true "创建Đơn vị"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /attendanceAgency/createAttendanceAgency [post]
func (attendanceAgencyApi *AttendanceAgencyApi) CreateAttendanceAgency(c *gin.Context) {
	var attendanceAgency checkins.AttendanceAgency
	err := c.ShouldBindJSON(&attendanceAgency)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    attendanceAgency.CreatedBy = utils.GetUserID(c)
	err = attendanceAgencyService.CreateAttendanceAgency(&attendanceAgency)
	if err != nil {
        global.GVA_LOG.Error("thất bại!", zap.Error(err))
		response.FailWithMessage("thất bại:" + err.Error(), c)
		return
	}
    response.OkWithMessage("thành công", c)
}

// DeleteAttendanceAgency 删除Đơn vị
// @Tags AttendanceAgency
// @Summary 删除Đơn vị
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body checkins.AttendanceAgency true "删除Đơn vị"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /attendanceAgency/deleteAttendanceAgency [delete]
func (attendanceAgencyApi *AttendanceAgencyApi) DeleteAttendanceAgency(c *gin.Context) {
	ID := c.Query("ID")
    userID := utils.GetUserID(c)
	err := attendanceAgencyService.DeleteAttendanceAgency(ID,userID)
	if err != nil {
        global.GVA_LOG.Error("thất bại!", zap.Error(err))
		response.FailWithMessage("thất bại:" + err.Error(), c)
		return
	}
	response.OkWithMessage("thành công", c)
}

// DeleteAttendanceAgencyByIds 批量删除Đơn vị
// @Tags AttendanceAgency
// @Summary 批量删除Đơn vị
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /attendanceAgency/deleteAttendanceAgencyByIds [delete]
func (attendanceAgencyApi *AttendanceAgencyApi) DeleteAttendanceAgencyByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
    userID := utils.GetUserID(c)
	err := attendanceAgencyService.DeleteAttendanceAgencyByIds(IDs,userID)
	if err != nil {
        global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:" + err.Error(), c)
		return
	}
	response.OkWithMessage("Thành công", c)
}

// UpdateAttendanceAgency 更新Đơn vị
// @Tags AttendanceAgency
// @Summary 更新Đơn vị
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body checkins.AttendanceAgency true "更新Đơn vị"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /attendanceAgency/updateAttendanceAgency [put]
func (attendanceAgencyApi *AttendanceAgencyApi) UpdateAttendanceAgency(c *gin.Context) {
	var attendanceAgency checkins.AttendanceAgency
	err := c.ShouldBindJSON(&attendanceAgency)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    attendanceAgency.UpdatedBy = utils.GetUserID(c)
	err = attendanceAgencyService.UpdateAttendanceAgency(attendanceAgency)
	if err != nil {
        global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:" + err.Error(), c)
		return
	}
	response.OkWithMessage("Thành công", c)
}

// FindAttendanceAgency 用id查询Đơn vị
// @Tags AttendanceAgency
// @Summary 用id查询Đơn vị
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query checkins.AttendanceAgency true "用id查询Đơn vị"
// @Success 200 {object} response.Response{data=checkins.AttendanceAgency,msg=string} "查询成功"
// @Router /attendanceAgency/findAttendanceAgency [get]
func (attendanceAgencyApi *AttendanceAgencyApi) FindAttendanceAgency(c *gin.Context) {
	ID := c.Query("ID")
	reattendanceAgency, err := attendanceAgencyService.GetAttendanceAgency(ID)
	if err != nil {
        global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:" + err.Error(), c)
		return
	}
	response.OkWithData(reattendanceAgency, c)
}

// GetAttendanceAgencyList 分页获取Đơn vị列表
// @Tags AttendanceAgency
// @Summary 分页获取Đơn vị列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query checkinsReq.AttendanceAgencySearch true "分页获取Đơn vị列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /attendanceAgency/getAttendanceAgencyList [get]
func (attendanceAgencyApi *AttendanceAgencyApi) GetAttendanceAgencyList(c *gin.Context) {
	var pageInfo checkinsReq.AttendanceAgencySearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := attendanceAgencyService.GetAttendanceAgencyInfoList(pageInfo)
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

// GetAttendanceAgencyPublic 不需要鉴权的Đơn vị接口
// @Tags AttendanceAgency
// @Summary 不需要鉴权的Đơn vị接口
// @accept application/json
// @Produce application/json
// @Param data query checkinsReq.AttendanceAgencySearch true "分页获取Đơn vị列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /attendanceAgency/getAttendanceAgencyPublic [get]
func (attendanceAgencyApi *AttendanceAgencyApi) GetAttendanceAgencyPublic(c *gin.Context) {
    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的Đơn vị接口信息",
    }, "Thành công", c)
}
