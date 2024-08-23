package config

import (
	"path/filepath"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/config"
	configReq "github.com/flipped-aurora/gin-vue-admin/server/model/config/request"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"
)

type CfgFileProcessApi struct{}

// CreateCfgFileProcess 创建Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống
// @Tags CfgFileProcess
// @Summary 创建Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body config.CfgFileProcess true "创建Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /FileProcess/createCfgFileProcess [post]
func (FileProcessApi *CfgFileProcessApi) CreateCfgFileProcess(c *gin.Context) {
	action := c.PostForm("action")
	// typeProcess := c.PostForm("type")

	if isValid := config.CheckFileProcessAction(action); !isValid {
		response.FailWithMessage("Hành động không hợp lệ", c)
		return
	}

	// Init data
	u, err := uuid.NewV4()
	if err != nil {
		response.FailWithMessage("Tạo mã UUID thất bại", c)
		return
	}
	uuid := u.String()

	if err != nil {
		response.FailWithMessage("Tạo mã UUID thất bại", c)
		return
	}

	// var process config.CfgFileProcess = config.CfgFileProcess{
	// 	GVA_MODEL:      global.GVA_MODEL{},
	// 	FileName:       "",
	// 	Uuid:           uuid,
	// 	Msg:            "Tập tin đang được xử lý",
	// 	UniqueFileName: uuid + ".xlsx",
	// 	Extension:      "xlsx",
	// 	FileSize:       0,
	// 	TotalRow:       0,
	// 	Action:         action,
	// 	Type:           typeProcess,
	// 	Url:            "",
	// 	CreatedBy:      utils.GetUserID(c),
	// 	UpdatedBy:      0,
	// 	DeletedBy:      0,
	// }

	switch action {
	case "import":
		{
			// Parse file from form-data
			file, err := c.FormFile("file")
			if err != nil {
				response.FailWithMessage("Không tìm thấy tệp tin", c)
				return
			}

			// Init data
			extention := filepath.Ext(file.Filename)
			uniqueFileName := uuid + extention

			// Save file to server
			if err := c.SaveUploadedFile(file, global.GVA_CONFIG.Local.Path+"/"+uniqueFileName); err != nil {
				response.FailWithMessage(err.Error(), c)
				return
			}

			// Create new record process with 10% progress

			// Handle case of type example: "IMPORT_SAMPLE"
			break
		}
	case "export":
		break
	default:

		break
	}

	// var FileProcess config.CfgFileProcess
	// err := c.ShouldBindJSON(&FileProcess)
	// if err != nil {
	// 	response.FailWithMessage(err.Error(), c)
	// 	return
	// }
	// err = FileProcessService.CreateCfgFileProcess(&FileProcess)
	// if err != nil {
	// 	global.GVA_LOG.Error("Tạo thất bại!", zap.Error(err))
	// 	response.FailWithMessage("Tạo thất bại:"+err.Error(), c)
	// 	return
	// }
	// response.OkWithMessage("Tạo thành công", c)
}

// DeleteCfgFileProcess 删除Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống
// @Tags CfgFileProcess
// @Summary 删除Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body config.CfgFileProcess true "删除Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /FileProcess/deleteCfgFileProcess [delete]
func (FileProcessApi *CfgFileProcessApi) DeleteCfgFileProcess(c *gin.Context) {
	ID := c.Query("ID")
	err := FileProcessService.DeleteCfgFileProcess(ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteCfgFileProcessByIds 批量删除Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống
// @Tags CfgFileProcess
// @Summary 批量删除Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /FileProcess/deleteCfgFileProcessByIds [delete]
func (FileProcessApi *CfgFileProcessApi) DeleteCfgFileProcessByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	err := FileProcessService.DeleteCfgFileProcessByIds(IDs)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateCfgFileProcess 更新Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống
// @Tags CfgFileProcess
// @Summary 更新Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body config.CfgFileProcess true "更新Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /FileProcess/updateCfgFileProcess [put]
func (FileProcessApi *CfgFileProcessApi) UpdateCfgFileProcess(c *gin.Context) {
	var FileProcess config.CfgFileProcess
	err := c.ShouldBindJSON(&FileProcess)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = FileProcessService.UpdateCfgFileProcess(FileProcess)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindCfgFileProcess 用id查询Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống
// @Tags CfgFileProcess
// @Summary 用id查询Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query config.CfgFileProcess true "用id查询Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống"
// @Success 200 {object} response.Response{data=config.CfgFileProcess,msg=string} "查询成功"
// @Router /FileProcess/findCfgFileProcess [get]
func (FileProcessApi *CfgFileProcessApi) FindCfgFileProcess(c *gin.Context) {
	ID := c.Query("ID")
	reFileProcess, err := FileProcessService.GetCfgFileProcess(ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(reFileProcess, c)
}

// GetCfgFileProcessList 分页获取Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống列表
// @Tags CfgFileProcess
// @Summary 分页获取Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query configReq.CfgFileProcessSearch true "分页获取Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /FileProcess/getCfgFileProcessList [get]
func (FileProcessApi *CfgFileProcessApi) GetCfgFileProcessList(c *gin.Context) {
	var pageInfo configReq.CfgFileProcessSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := FileProcessService.GetCfgFileProcessInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetCfgFileProcessPublic 不需要鉴权的Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống接口
// @Tags CfgFileProcess
// @Summary 不需要鉴权的Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống接口
// @accept application/json
// @Produce application/json
// @Param data query configReq.CfgFileProcessSearch true "分页获取Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /FileProcess/getCfgFileProcessPublic [get]
func (FileProcessApi *CfgFileProcessApi) GetCfgFileProcessPublic(c *gin.Context) {
	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống接口信息",
	}, "获取成功", c)
}
