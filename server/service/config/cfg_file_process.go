package config

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/config"
	configReq "github.com/flipped-aurora/gin-vue-admin/server/model/config/request"
)

type CfgFileProcessService struct{}

// CreateCfgFileProcess 创建Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống记录
// Author [piexlmax](https://github.com/piexlmax)
func (FileProcessService *CfgFileProcessService) CreateCfgFileProcess(FileProcess *config.CfgFileProcess) (err error) {
	err = global.GVA_DB.Create(FileProcess).Error
	return err
}

// DeleteCfgFileProcess 删除Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống记录
// Author [piexlmax](https://github.com/piexlmax)
func (FileProcessService *CfgFileProcessService) DeleteCfgFileProcess(ID string) (err error) {
	err = global.GVA_DB.Delete(&config.CfgFileProcess{}, "id = ?", ID).Error
	return err
}

// DeleteCfgFileProcessByIds 批量删除Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống记录
// Author [piexlmax](https://github.com/piexlmax)
func (FileProcessService *CfgFileProcessService) DeleteCfgFileProcessByIds(IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]config.CfgFileProcess{}, "id in ?", IDs).Error
	return err
}

// UpdateCfgFileProcess 更新Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống记录
// Author [piexlmax](https://github.com/piexlmax)
func (FileProcessService *CfgFileProcessService) UpdateCfgFileProcess(FileProcess config.CfgFileProcess) (err error) {
	err = global.GVA_DB.Model(&config.CfgFileProcess{}).Where("id = ?", FileProcess.ID).Updates(&FileProcess).Error
	return err
}

// GetCfgFileProcess 根据ID获取Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống记录
// Author [piexlmax](https://github.com/piexlmax)
func (FileProcessService *CfgFileProcessService) GetCfgFileProcess(ID string) (FileProcess config.CfgFileProcess, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&FileProcess).Error
	return
}

// GetCfgFileProcessInfoList 分页获取Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống记录
// Author [piexlmax](https://github.com/piexlmax)
func (FileProcessService *CfgFileProcessService) GetCfgFileProcessInfoList(info configReq.CfgFileProcessSearch) (list []config.CfgFileProcess, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&config.CfgFileProcess{})
	var FileProcesss []config.CfgFileProcess
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.Uuid != "" {
		db = db.Where("uuid = ?", info.Uuid)
	}
	if info.Percent != nil {
		db = db.Where("percent = ?", info.Percent)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&FileProcesss).Error
	return FileProcesss, total, err
}

func (fileProcessService *CfgFileProcessService) GetPercentFileProcess(f config.CfgFileProcess) (fileProcess config.CfgFileProcess, err error) {
	// get record with uuidFile
	uuid := f.Uuid
	userId := f.UserID

	err = global.GVA_DB.Where("uuid = ? AND user_id = ?", uuid, userId).First(&fileProcess).Error
	return
}
