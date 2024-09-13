package example

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
	exampleRes "github.com/flipped-aurora/gin-vue-admin/server/model/example/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type FileUploadAndDownloadApi struct{}

// UploadFile
// @Tags      ExaFileUploadAndDownload
// @Summary   Ví dụ tải lên tệp
// @Security  ApiKeyAuth
// @accept    multipart/form-data
// @Produce   application/json
// @Param     file  formData  file                                                           true  "Ví dụ tải lên tệp"
// @Success   200   {object}  response.Response{data=exampleRes.ExaFileResponse,msg=string}  "Ví dụ tải lên tệp, trả về chi tiết tệp"
// @Router    /fileUploadAndDownload/upload [post]
func (b *FileUploadAndDownloadApi) UploadFile(c *gin.Context) {
	var file example.ExaFileUploadAndDownload
	noSave := c.DefaultQuery("noSave", "0")
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		global.GVA_LOG.Error("Không thể nhận tệp!", zap.Error(err))
		response.FailWithMessage("Không thể nhận tệp", c)
		return
	}
	file, err = fileUploadAndDownloadService.UploadFile(header, noSave) // Sau khi tải lên tệp, nhận được đường dẫn tệp
	if err != nil {
		global.GVA_LOG.Error("Không thể chỉnh sửa liên kết cơ sở dữ liệu!", zap.Error(err))
		response.FailWithMessage("Không thể chỉnh sửa liên kết cơ sở dữ liệu", c)
		return
	}
	response.OkWithDetailed(exampleRes.ExaFileResponse{File: file}, "Tải lên thành công", c)
}

// EditFileName chỉnh sửa tên tệp hoặc ghi chú
func (b *FileUploadAndDownloadApi) EditFileName(c *gin.Context) {
	var file example.ExaFileUploadAndDownload
	err := c.ShouldBindJSON(&file)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = fileUploadAndDownloadService.EditFileName(file)
	if err != nil {
		global.GVA_LOG.Error("Chỉnh sửa thất bại!", zap.Error(err))
		response.FailWithMessage("Chỉnh sửa thất bại", c)
		return
	}
	response.OkWithMessage("Chỉnh sửa thành công", c)
}

// DeleteFile
// @Tags      ExaFileUploadAndDownload
// @Summary   Xóa tệp
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body      example.ExaFileUploadAndDownload  true  "Chỉ cần truyền id của tệp vào"
// @Success   200   {object}  response.Response{msg=string}     "Xóa tệp"
// @Router    /fileUploadAndDownload/deleteFile [post]
func (b *FileUploadAndDownloadApi) DeleteFile(c *gin.Context) {
	var file example.ExaFileUploadAndDownload
	err := c.ShouldBindJSON(&file)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := fileUploadAndDownloadService.DeleteFile(file); err != nil {
		global.GVA_LOG.Error("Xóa thất bại!", zap.Error(err))
		response.FailWithMessage("Xóa thất bại", c)
		return
	}
	response.OkWithMessage("Xóa thành công", c)
}

// GetFileList
// @Tags      ExaFileUploadAndDownload
// @Summary   Danh sách tệp phân trang
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.PageInfo                                        true  "Trang, kích thước mỗi trang"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "Danh sách tệp phân trang, trả về danh sách, tổng số, trang, kích thước mỗi trang"
// @Router    /fileUploadAndDownload/getFileList [post]
func (b *FileUploadAndDownloadApi) GetFileList(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := fileUploadAndDownloadService.GetFileRecordInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("Lấy danh sách thất bại!", zap.Error(err))
		response.FailWithMessage("Lấy danh sách thất bại", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "Lấy danh sách thành công", c)
}
