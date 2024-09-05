package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
	checkinsReq "github.com/flipped-aurora/gin-vue-admin/server/model/checkins/request"
	"gorm.io/gorm"
)

type AttendanceCategoryService struct{}

// CreateAttendanceCategory 创建Danh mục điểm danh记录
// Author [piexlmax](https://github.com/piexlmax)
func (attendanceCategoryService *AttendanceCategoryService) CreateAttendanceCategory(attendanceCategory *checkins.AttendanceCategory) (err error) {

	if attendanceCategory.IsCurrent {
		global.GVA_DB.Model(&checkins.AttendanceCategory{}).Where("is_current = ?", true).Update("is_current", false)
	}
	err = global.GVA_DB.Create(attendanceCategory).Error
	return err
}

// DeleteAttendanceCategory 删除Danh mục điểm danh记录
// Author [piexlmax](https://github.com/piexlmax)
func (attendanceCategoryService *AttendanceCategoryService) DeleteAttendanceCategory(ID string, userID uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&checkins.AttendanceCategory{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
			return err
		}
		if err = tx.Delete(&checkins.AttendanceCategory{}, "id = ?", ID).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteAttendanceCategoryByIds 批量删除Danh mục điểm danh记录
// Author [piexlmax](https://github.com/piexlmax)
func (attendanceCategoryService *AttendanceCategoryService) DeleteAttendanceCategoryByIds(IDs []string, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&checkins.AttendanceCategory{}).Where("id in ?", IDs).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", IDs).Delete(&checkins.AttendanceCategory{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateAttendanceCategory 更新Danh mục điểm danh记录
// Author [piexlmax](https://github.com/piexlmax)
func (attendanceCategoryService *AttendanceCategoryService) UpdateAttendanceCategory(attendanceCategory checkins.AttendanceCategory) (err error) {
	err = global.GVA_DB.Model(&checkins.AttendanceCategory{}).Where("id = ?", attendanceCategory.ID).Updates(&attendanceCategory).Error
	return err
}

// GetAttendanceCategory 根据ID获取Danh mục điểm danh记录
// Author [piexlmax](https://github.com/piexlmax)
func (attendanceCategoryService *AttendanceCategoryService) GetAttendanceCategory(ID string) (attendanceCategory checkins.AttendanceCategory, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&attendanceCategory).Error
	return
}

// GetAttendanceCategoryInfoList 分页获取Danh mục điểm danh记录
// Author [piexlmax](https://github.com/piexlmax)
func (attendanceCategoryService *AttendanceCategoryService) GetAttendanceCategoryInfoList(info checkinsReq.AttendanceCategorySearch) (list []checkins.AttendanceCategory, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&checkins.AttendanceCategory{})
	var attendanceCategorys []checkins.AttendanceCategory
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.Name != nil {
		db = db.Where("name = ?", info.Name)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 && limit != -1 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&attendanceCategorys).Error
	return attendanceCategorys, total, err
}
