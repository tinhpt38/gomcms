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

type AreaApi struct {}

// CreateArea tạo Khu vực điểm danh
// @Tags Area
// @Summary tạo Khu vực điểm danh
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body checkins.Area true "tạo Khu vực điểm danh"
// @Success 200 {object} response.Response{msg=string} "tạo thành công"
// @Router /area/createArea [post]
func (areaApi *AreaApi) CreateArea(c *gin.Context) {
	var area checkins.Area
	err := c.ShouldBindJSON(&area)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	area.CreatedBy = utils.GetUserID(c)
	err = areaService.CreateArea(&area)
	if err != nil {
		global.GVA_LOG.Error("thất bại!", zap.Error(err))
		response.FailWithMessage("thất bại:" + err.Error(), c)
		return
	}
	response.OkWithMessage("thành công", c)
}

// DeleteArea xóa Khu vực điểm danh
// @Tags Area
// @Summary xóa Khu vực điểm danh
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body checkins.Area true "xóa Khu vực điểm danh"
// @Success 200 {object} response.Response{msg=string} "xóa thành công"
// @Router /area/deleteArea [delete]
func (areaApi *AreaApi) DeleteArea(c *gin.Context) {
	ID := c.Query("ID")
	userID := utils.GetUserID(c)
	err := areaService.DeleteArea(ID,userID)
	if err != nil {
		global.GVA_LOG.Error("thất bại!", zap.Error(err))
		response.FailWithMessage("thất bại:" + err.Error(), c)
		return
	}
	response.OkWithMessage("thành công", c)
}

// DeleteAreaByIds xóa nhiều Khu vực điểm danh
// @Tags Area
// @Summary xóa nhiều Khu vực điểm danh
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "xóa thành công"
// @Router /area/deleteAreaByIds [delete]
func (areaApi *AreaApi) DeleteAreaByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	userID := utils.GetUserID(c)
	err := areaService.DeleteAreaByIds(IDs,userID)
	if err != nil {
		global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:" + err.Error(), c)
		return
	}
	response.OkWithMessage("Thành công", c)
}

// UpdateArea cập nhật Khu vực điểm danh
// @Tags Area
// @Summary cập nhật Khu vực điểm danh
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body checkins.Area true "cập nhật Khu vực điểm danh"
// @Success 200 {object} response.Response{msg=string} "cập nhật thành công"
// @Router /area/updateArea [put]
func (areaApi *AreaApi) UpdateArea(c *gin.Context) {
	var area checkins.Area
	err := c.ShouldBindJSON(&area)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	area.UpdatedBy = utils.GetUserID(c)
	err = areaService.UpdateArea(area)
	if err != nil {
		global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:" + err.Error(), c)
		return
	}
	response.OkWithMessage("Thành công", c)
}

// FindArea tìm Khu vực điểm danh bằng id
// @Tags Area
// @Summary tìm Khu vực điểm danh bằng id
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query checkins.Area true "tìm Khu vực điểm danh bằng id"
// @Success 200 {object} response.Response{data=checkins.Area,msg=string} "tìm thành công"
// @Router /area/findArea [get]
func (areaApi *AreaApi) FindArea(c *gin.Context) {
	ID := c.Query("ID")
	rearea, err := areaService.GetArea(ID)
	if err != nil {
		global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:" + err.Error(), c)
		return
	}
	response.OkWithData(rearea, c)
}

// GetAreaList lấy danh sách Khu vực điểm danh theo trang
// @Tags Area
// @Summary lấy danh sách Khu vực điểm danh theo trang
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query checkinsReq.AreaSearch true "lấy danh sách Khu vực điểm danh theo trang"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "lấy thành công"
// @Router /area/getAreaList [get]
func (areaApi *AreaApi) GetAreaList(c *gin.Context) {
	var pageInfo checkinsReq.AreaSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := areaService.GetAreaInfoList(pageInfo)
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

// GetAreaPublic API Khu vực điểm danh không cần xác thực
// @Tags Area
// @Summary API Khu vực điểm danh không cần xác thực
// @accept application/json
// @Produce application/json
// @Param data query checkinsReq.AreaSearch true "lấy danh sách Khu vực điểm danh theo trang"
// @Success 200 {object} response.Response{data=object,msg=string} "lấy thành công"
// @Router /area/getAreaPublic [get]
func (areaApi *AreaApi) GetAreaPublic(c *gin.Context) {
	// API này không cần xác thực
	// Ví dụ trả về một thông điệp cố định, thường được sử dụng cho dịch vụ C端, cần tự triển khai logic kinh doanh của riêng mình
	response.OkWithDetailed(gin.H{
	   "info": "Thông tin API Khu vực điểm danh không cần xác thực",
	}, "Thành công", c)
}

