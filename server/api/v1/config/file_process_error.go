package config

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/model/config"
    configReq "github.com/flipped-aurora/gin-vue-admin/server/model/config/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type FileProcessErrorApi struct {}

// CreateFileProcessError 创建Lỗi ghi nhận từ quá trình xử lý tệp tin
// @Tags FileProcessError
// @Summary 创建Lỗi ghi nhận từ quá trình xử lý tệp tin
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body config.FileProcessError true "创建Lỗi ghi nhận từ quá trình xử lý tệp tin"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /fileProcessError/createFileProcessError [post]
func (fileProcessErrorApi *FileProcessErrorApi) CreateFileProcessError(c *gin.Context) {
	var fileProcessError config.FileProcessError
	err := c.ShouldBindJSON(&fileProcessError)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = fileProcessErrorService.CreateFileProcessError(&fileProcessError)
	if err != nil {
        global.GVA_LOG.Error("thất bại!", zap.Error(err))
		response.FailWithMessage("thất bại:" + err.Error(), c)
		return
	}
    response.OkWithMessage("thành công", c)
}

// DeleteFileProcessError 删除Lỗi ghi nhận từ quá trình xử lý tệp tin
// @Tags FileProcessError
// @Summary 删除Lỗi ghi nhận từ quá trình xử lý tệp tin
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body config.FileProcessError true "删除Lỗi ghi nhận từ quá trình xử lý tệp tin"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /fileProcessError/deleteFileProcessError [delete]
func (fileProcessErrorApi *FileProcessErrorApi) DeleteFileProcessError(c *gin.Context) {
	ID := c.Query("ID")
	err := fileProcessErrorService.DeleteFileProcessError(ID)
	if err != nil {
        global.GVA_LOG.Error("thất bại!", zap.Error(err))
		response.FailWithMessage("thất bại:" + err.Error(), c)
		return
	}
	response.OkWithMessage("thành công", c)
}

// DeleteFileProcessErrorByIds 批量删除Lỗi ghi nhận từ quá trình xử lý tệp tin
// @Tags FileProcessError
// @Summary 批量删除Lỗi ghi nhận từ quá trình xử lý tệp tin
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /fileProcessError/deleteFileProcessErrorByIds [delete]
func (fileProcessErrorApi *FileProcessErrorApi) DeleteFileProcessErrorByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	err := fileProcessErrorService.DeleteFileProcessErrorByIds(IDs)
	if err != nil {
        global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:" + err.Error(), c)
		return
	}
	response.OkWithMessage("Thành công", c)
}

// UpdateFileProcessError 更新Lỗi ghi nhận từ quá trình xử lý tệp tin
// @Tags FileProcessError
// @Summary 更新Lỗi ghi nhận từ quá trình xử lý tệp tin
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body config.FileProcessError true "更新Lỗi ghi nhận từ quá trình xử lý tệp tin"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /fileProcessError/updateFileProcessError [put]
func (fileProcessErrorApi *FileProcessErrorApi) UpdateFileProcessError(c *gin.Context) {
	var fileProcessError config.FileProcessError
	err := c.ShouldBindJSON(&fileProcessError)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = fileProcessErrorService.UpdateFileProcessError(fileProcessError)
	if err != nil {
        global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:" + err.Error(), c)
		return
	}
	response.OkWithMessage("Thành công", c)
}

// FindFileProcessError 用id查询Lỗi ghi nhận từ quá trình xử lý tệp tin
// @Tags FileProcessError
// @Summary 用id查询Lỗi ghi nhận từ quá trình xử lý tệp tin
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query config.FileProcessError true "用id查询Lỗi ghi nhận từ quá trình xử lý tệp tin"
// @Success 200 {object} response.Response{data=config.FileProcessError,msg=string} "查询成功"
// @Router /fileProcessError/findFileProcessError [get]
func (fileProcessErrorApi *FileProcessErrorApi) FindFileProcessError(c *gin.Context) {
	ID := c.Query("ID")
	refileProcessError, err := fileProcessErrorService.GetFileProcessError(ID)
	if err != nil {
        global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:" + err.Error(), c)
		return
	}
	response.OkWithData(refileProcessError, c)
}

// GetFileProcessErrorList 分页获取Lỗi ghi nhận từ quá trình xử lý tệp tin列表
// @Tags FileProcessError
// @Summary 分页获取Lỗi ghi nhận từ quá trình xử lý tệp tin列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query configReq.FileProcessErrorSearch true "分页获取Lỗi ghi nhận từ quá trình xử lý tệp tin列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /fileProcessError/getFileProcessErrorList [get]
func (fileProcessErrorApi *FileProcessErrorApi) GetFileProcessErrorList(c *gin.Context) {
	var pageInfo configReq.FileProcessErrorSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := fileProcessErrorService.GetFileProcessErrorInfoList(pageInfo)
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

// GetFileProcessErrorPublic 不需要鉴权的Lỗi ghi nhận từ quá trình xử lý tệp tin接口
// @Tags FileProcessError
// @Summary 不需要鉴权的Lỗi ghi nhận từ quá trình xử lý tệp tin接口
// @accept application/json
// @Produce application/json
// @Param data query configReq.FileProcessErrorSearch true "分页获取Lỗi ghi nhận từ quá trình xử lý tệp tin列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /fileProcessError/getFileProcessErrorPublic [get]
func (fileProcessErrorApi *FileProcessErrorApi) GetFileProcessErrorPublic(c *gin.Context) {
    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的Lỗi ghi nhận từ quá trình xử lý tệp tin接口信息",
    }, "Thành công", c)
}
