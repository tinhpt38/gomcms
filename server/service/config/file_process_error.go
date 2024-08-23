package config

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/config"
    configReq "github.com/flipped-aurora/gin-vue-admin/server/model/config/request"
)

type FileProcessErrorService struct {}

// CreateFileProcessError 创建Lỗi ghi nhận từ quá trình xử lý tệp tin记录
// Author [piexlmax](https://github.com/piexlmax)
func (fileProcessErrorService *FileProcessErrorService) CreateFileProcessError(fileProcessError *config.FileProcessError) (err error) {
	err = global.GVA_DB.Create(fileProcessError).Error
	return err
}

// DeleteFileProcessError 删除Lỗi ghi nhận từ quá trình xử lý tệp tin记录
// Author [piexlmax](https://github.com/piexlmax)
func (fileProcessErrorService *FileProcessErrorService)DeleteFileProcessError(ID string) (err error) {
	err = global.GVA_DB.Delete(&config.FileProcessError{},"id = ?",ID).Error
	return err
}

// DeleteFileProcessErrorByIds 批量删除Lỗi ghi nhận từ quá trình xử lý tệp tin记录
// Author [piexlmax](https://github.com/piexlmax)
func (fileProcessErrorService *FileProcessErrorService)DeleteFileProcessErrorByIds(IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]config.FileProcessError{},"id in ?",IDs).Error
	return err
}

// UpdateFileProcessError 更新Lỗi ghi nhận từ quá trình xử lý tệp tin记录
// Author [piexlmax](https://github.com/piexlmax)
func (fileProcessErrorService *FileProcessErrorService)UpdateFileProcessError(fileProcessError config.FileProcessError) (err error) {
	err = global.GVA_DB.Model(&config.FileProcessError{}).Where("id = ?",fileProcessError.ID).Updates(&fileProcessError).Error
	return err
}

// GetFileProcessError 根据ID获取Lỗi ghi nhận từ quá trình xử lý tệp tin记录
// Author [piexlmax](https://github.com/piexlmax)
func (fileProcessErrorService *FileProcessErrorService)GetFileProcessError(ID string) (fileProcessError config.FileProcessError, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&fileProcessError).Error
	return
}

// GetFileProcessErrorInfoList 分页获取Lỗi ghi nhận từ quá trình xử lý tệp tin记录
// Author [piexlmax](https://github.com/piexlmax)
func (fileProcessErrorService *FileProcessErrorService)GetFileProcessErrorInfoList(info configReq.FileProcessErrorSearch) (list []config.FileProcessError, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&config.FileProcessError{})
    var fileProcessErrors []config.FileProcessError
    // 如果有条件搜索 下方会自动创建搜索语句
    if info.StartCreatedAt !=nil && info.EndCreatedAt !=nil {
     db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
    }
    if info.FileProcessId != nil {
        db = db.Where("file_process_id = ?",info.FileProcessId)
    }
    if info.FieldTitle != "" {
        db = db.Where("field_title LIKE ?","%"+ info.FieldTitle+"%")
    }
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }

	if limit != 0 {
       db = db.Limit(limit).Offset(offset)
    }
	
	err = db.Find(&fileProcessErrors).Error
	return  fileProcessErrors, total, err
}