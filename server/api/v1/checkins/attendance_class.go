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

// CreateAttendanceClass 创建AttendanceClass
// @Tags AttendanceClass
// @Summary 创建AttendanceClass
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body checkins.AttendanceClass true "创建AttendanceClass"
// @Success 200 {object} response.Response{msg=string} "创建成功"
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
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteAttendanceClass 删除AttendanceClass
// @Tags AttendanceClass
// @Summary 删除AttendanceClass
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body checkins.AttendanceClass true "删除AttendanceClass"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /attendanceClass/deleteAttendanceClass [delete]
func (attendanceClassApi *AttendanceClassApi) DeleteAttendanceClass(c *gin.Context) {
	ID := c.Query("ID")
    userID := utils.GetUserID(c)
	err := attendanceClassService.DeleteAttendanceClass(ID,userID)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteAttendanceClassByIds 批量删除AttendanceClass
// @Tags AttendanceClass
// @Summary 批量删除AttendanceClass
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /attendanceClass/deleteAttendanceClassByIds [delete]
func (attendanceClassApi *AttendanceClassApi) DeleteAttendanceClassByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
    userID := utils.GetUserID(c)
	err := attendanceClassService.DeleteAttendanceClassByIds(IDs,userID)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateAttendanceClass 更新AttendanceClass
// @Tags AttendanceClass
// @Summary 更新AttendanceClass
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body checkins.AttendanceClass true "更新AttendanceClass"
// @Success 200 {object} response.Response{msg=string} "更新成功"
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
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindAttendanceClass 用id查询AttendanceClass
// @Tags AttendanceClass
// @Summary 用id查询AttendanceClass
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query checkins.AttendanceClass true "用id查询AttendanceClass"
// @Success 200 {object} response.Response{data=checkins.AttendanceClass,msg=string} "查询成功"
// @Router /attendanceClass/findAttendanceClass [get]
func (attendanceClassApi *AttendanceClassApi) FindAttendanceClass(c *gin.Context) {
	ID := c.Query("ID")
	reattendanceClass, err := attendanceClassService.GetAttendanceClass(ID)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(reattendanceClass, c)
}

// GetAttendanceClassList 分页获取AttendanceClass列表
// @Tags AttendanceClass
// @Summary 分页获取AttendanceClass列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query checkinsReq.AttendanceClassSearch true "分页获取AttendanceClass列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
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
	    global.GVA_LOG.Error("获取失败!", zap.Error(err))
        response.FailWithMessage("获取失败:" + err.Error(), c)
        return
    }
    response.OkWithDetailed(response.PageResult{
        List:     list,
        Total:    total,
        Page:     pageInfo.Page,
        PageSize: pageInfo.PageSize,
    }, "获取成功", c)
}

// GetAttendanceClassPublic 不需要鉴权的AttendanceClass接口
// @Tags AttendanceClass
// @Summary 不需要鉴权的AttendanceClass接口
// @accept application/json
// @Produce application/json
// @Param data query checkinsReq.AttendanceClassSearch true "分页获取AttendanceClass列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /attendanceClass/getAttendanceClassPublic [get]
func (attendanceClassApi *AttendanceClassApi) GetAttendanceClassPublic(c *gin.Context) {
    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的AttendanceClass接口信息",
    }, "获取成功", c)
}
