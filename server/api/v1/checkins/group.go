package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
	checkinsReq "github.com/flipped-aurora/gin-vue-admin/server/model/checkins/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GroupApi struct{}

// CreateGroup tạo Nhóm
// @Tags Group
// @Summary tạo Nhóm
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body checkins.Group true "tạo Nhóm"
// @Success 200 {object} response.Response{msg=string} "tạo thành công"
// @Router /group/createGroup [post]
func (groupApi *GroupApi) CreateGroup(c *gin.Context) {
	var group checkins.Group
	err := c.ShouldBindJSON(&group)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = groupService.CreateGroup(&group)
	if err != nil {
		global.GVA_LOG.Error("tạo thất bại!", zap.Error(err))
		response.FailWithMessage("tạo thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("tạo thành công", c)
}

// DeleteGroup xóa Nhóm
// @Tags Group
// @Summary xóa Nhóm
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body checkins.Group true "xóa Nhóm"
// @Success 200 {object} response.Response{msg=string} "xóa thành công"
// @Router /group/deleteGroup [delete]
func (groupApi *GroupApi) DeleteGroup(c *gin.Context) {
	ID := c.Query("ID")
	err := groupService.DeleteGroup(ID)
	if err != nil {
		global.GVA_LOG.Error("xóa thất bại!", zap.Error(err))
		response.FailWithMessage("xóa thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("xóa thành công", c)
}

// DeleteGroupByIds xóa nhiều Nhóm
// @Tags Group
// @Summary xóa nhiều Nhóm
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "xóa nhiều thành công"
// @Router /group/deleteGroupByIds [delete]
func (groupApi *GroupApi) DeleteGroupByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	err := groupService.DeleteGroupByIds(IDs)
	if err != nil {
		global.GVA_LOG.Error("xóa nhiều thất bại!", zap.Error(err))
		response.FailWithMessage("xóa nhiều thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("xóa nhiều thành công", c)
}

// UpdateGroup cập nhật Nhóm
// @Tags Group
// @Summary cập nhật Nhóm
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body checkins.Group true "cập nhật Nhóm"
// @Success 200 {object} response.Response{msg=string} "cập nhật thành công"
// @Router /group/updateGroup [put]
func (groupApi *GroupApi) UpdateGroup(c *gin.Context) {
	var group checkins.Group
	err := c.ShouldBindJSON(&group)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = groupService.UpdateGroup(group)
	if err != nil {
		global.GVA_LOG.Error("cập nhật thất bại!", zap.Error(err))
		response.FailWithMessage("cập nhật thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("cập nhật thành công", c)
}

// FindGroup tìm Nhóm bằng id
// @Tags Group
// @Summary tìm Nhóm bằng id
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query checkins.Group true "tìm Nhóm bằng id"
// @Success 200 {object} response.Response{data=checkins.Group,msg=string} "tìm thành công"
// @Router /group/findGroup [get]
func (groupApi *GroupApi) FindGroup(c *gin.Context) {
	ID := c.Query("ID")
	regroup, err := groupService.GetGroup(ID)
	if err != nil {
		global.GVA_LOG.Error("tìm thất bại!", zap.Error(err))
		response.FailWithMessage("tìm thất bại:"+err.Error(), c)
		return
	}
	response.OkWithData(regroup, c)
}

// GetGroupList lấy danh sách Nhóm theo trang
// @Tags Group
// @Summary lấy danh sách Nhóm theo trang
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query checkinsReq.GroupSearch true "lấy danh sách Nhóm theo trang"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "lấy thành công"
// @Router /group/getGroupList [get]
func (groupApi *GroupApi) GetGroupList(c *gin.Context) {
	var pageInfo checkinsReq.GroupSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := groupService.GetGroupInfoList(pageInfo)
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

// GetGroupDataSource lấy nguồn dữ liệu của Nhóm
// @Tags Group
// @Summary lấy nguồn dữ liệu của Nhóm
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "lấy thành công"
// @Router /group/getGroupDataSource [get]
func (groupApi *GroupApi) GetGroupDataSource(c *gin.Context) {
	// Đây là một API để lấy dữ liệu nguồn định nghĩa của Nhóm
	dataSource, err := groupService.GetGroupDataSource()
	if err != nil {
		global.GVA_LOG.Error("lấy thất bại!", zap.Error(err))
		response.FailWithMessage("lấy thất bại:"+err.Error(), c)
		return
	}
	response.OkWithData(dataSource, c)
}

// GetGroupPublic API của Nhóm không cần xác thực
// @Tags Group
// @Summary API của Nhóm không cần xác thực
// @accept application/json
// @Produce application/json
// @Param data query checkinsReq.GroupSearch true "lấy danh sách Nhóm theo trang"
// @Success 200 {object} response.Response{data=object,msg=string} "lấy thành công"
// @Router /group/getGroupPublic [get]
func (groupApi *GroupApi) GetGroupPublic(c *gin.Context) {
	// API này không cần xác thực
	// Ví dụ trả về một thông điệp cố định, thường được sử dụng cho dịch vụ phía C, cần triển khai logic kinh doanh của riêng mình
	response.OkWithDetailed(gin.H{
		"info": "Thông tin API của Nhóm không cần xác thực",
	}, "lấy thành công", c)
}
