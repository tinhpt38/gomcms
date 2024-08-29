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

type AttendanceCategoryApi struct {}

// CreateAttendanceCategory 创建Danh mục điểm danh
// @Tags AttendanceCategory
// @Summary 创建Danh mục điểm danh
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body checkins.AttendanceCategory true "创建Danh mục điểm danh"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /attendanceCategory/createAttendanceCategory [post]
func (attendanceCategoryApi *AttendanceCategoryApi) CreateAttendanceCategory(c *gin.Context) {
	var attendanceCategory checkins.AttendanceCategory
	err := c.ShouldBindJSON(&attendanceCategory)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    attendanceCategory.CreatedBy = utils.GetUserID(c)
	err = attendanceCategoryService.CreateAttendanceCategory(&attendanceCategory)
	if err != nil {
        global.GVA_LOG.Error("thất bại!", zap.Error(err))
		response.FailWithMessage("thất bại:" + err.Error(), c)
		return
	}
    response.OkWithMessage("thành công", c)
}

// DeleteAttendanceCategory 删除Danh mục điểm danh
// @Tags AttendanceCategory
// @Summary 删除Danh mục điểm danh
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body checkins.AttendanceCategory true "删除Danh mục điểm danh"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /attendanceCategory/deleteAttendanceCategory [delete]
func (attendanceCategoryApi *AttendanceCategoryApi) DeleteAttendanceCategory(c *gin.Context) {
	ID := c.Query("ID")
    userID := utils.GetUserID(c)
	err := attendanceCategoryService.DeleteAttendanceCategory(ID,userID)
	if err != nil {
        global.GVA_LOG.Error("thất bại!", zap.Error(err))
		response.FailWithMessage("thất bại:" + err.Error(), c)
		return
	}
	response.OkWithMessage("thành công", c)
}

// DeleteAttendanceCategoryByIds 批量删除Danh mục điểm danh
// @Tags AttendanceCategory
// @Summary 批量删除Danh mục điểm danh
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /attendanceCategory/deleteAttendanceCategoryByIds [delete]
func (attendanceCategoryApi *AttendanceCategoryApi) DeleteAttendanceCategoryByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
    userID := utils.GetUserID(c)
	err := attendanceCategoryService.DeleteAttendanceCategoryByIds(IDs,userID)
	if err != nil {
        global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:" + err.Error(), c)
		return
	}
	response.OkWithMessage("Thành công", c)
}

// UpdateAttendanceCategory 更新Danh mục điểm danh
// @Tags AttendanceCategory
// @Summary 更新Danh mục điểm danh
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body checkins.AttendanceCategory true "更新Danh mục điểm danh"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /attendanceCategory/updateAttendanceCategory [put]
func (attendanceCategoryApi *AttendanceCategoryApi) UpdateAttendanceCategory(c *gin.Context) {
	var attendanceCategory checkins.AttendanceCategory
	err := c.ShouldBindJSON(&attendanceCategory)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    attendanceCategory.UpdatedBy = utils.GetUserID(c)
	err = attendanceCategoryService.UpdateAttendanceCategory(attendanceCategory)
	if err != nil {
        global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:" + err.Error(), c)
		return
	}
	response.OkWithMessage("Thành công", c)
}

// FindAttendanceCategory 用id查询Danh mục điểm danh
// @Tags AttendanceCategory
// @Summary 用id查询Danh mục điểm danh
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query checkins.AttendanceCategory true "用id查询Danh mục điểm danh"
// @Success 200 {object} response.Response{data=checkins.AttendanceCategory,msg=string} "查询成功"
// @Router /attendanceCategory/findAttendanceCategory [get]
func (attendanceCategoryApi *AttendanceCategoryApi) FindAttendanceCategory(c *gin.Context) {
	ID := c.Query("ID")
	reattendanceCategory, err := attendanceCategoryService.GetAttendanceCategory(ID)
	if err != nil {
        global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:" + err.Error(), c)
		return
	}
	response.OkWithData(reattendanceCategory, c)
}

// GetAttendanceCategoryList 分页获取Danh mục điểm danh列表
// @Tags AttendanceCategory
// @Summary 分页获取Danh mục điểm danh列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query checkinsReq.AttendanceCategorySearch true "分页获取Danh mục điểm danh列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /attendanceCategory/getAttendanceCategoryList [get]
func (attendanceCategoryApi *AttendanceCategoryApi) GetAttendanceCategoryList(c *gin.Context) {
	var pageInfo checkinsReq.AttendanceCategorySearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := attendanceCategoryService.GetAttendanceCategoryInfoList(pageInfo)
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

// GetAttendanceCategoryPublic 不需要鉴权的Danh mục điểm danh接口
// @Tags AttendanceCategory
// @Summary 不需要鉴权的Danh mục điểm danh接口
// @accept application/json
// @Produce application/json
// @Param data query checkinsReq.AttendanceCategorySearch true "分页获取Danh mục điểm danh列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /attendanceCategory/getAttendanceCategoryPublic [get]
func (attendanceCategoryApi *AttendanceCategoryApi) GetAttendanceCategoryPublic(c *gin.Context) {
    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的Danh mục điểm danh接口信息",
    }, "Thành công", c)
}
