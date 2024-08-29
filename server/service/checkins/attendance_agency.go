package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
    checkinsReq "github.com/flipped-aurora/gin-vue-admin/server/model/checkins/request"
    "gorm.io/gorm"
)

type AttendanceAgencyService struct {}

// CreateAttendanceAgency 创建Đơn vị记录
// Author [piexlmax](https://github.com/piexlmax)
func (attendanceAgencyService *AttendanceAgencyService) CreateAttendanceAgency(attendanceAgency *checkins.AttendanceAgency) (err error) {
	err = global.GVA_DB.Create(attendanceAgency).Error
	return err
}

// DeleteAttendanceAgency 删除Đơn vị记录
// Author [piexlmax](https://github.com/piexlmax)
func (attendanceAgencyService *AttendanceAgencyService)DeleteAttendanceAgency(ID string,userID uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
	    if err := tx.Model(&checkins.AttendanceAgency{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
              return err
        }
        if err = tx.Delete(&checkins.AttendanceAgency{},"id = ?",ID).Error; err != nil {
              return err
        }
        return nil
	})
	return err
}

// DeleteAttendanceAgencyByIds 批量删除Đơn vị记录
// Author [piexlmax](https://github.com/piexlmax)
func (attendanceAgencyService *AttendanceAgencyService)DeleteAttendanceAgencyByIds(IDs []string,deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
	    if err := tx.Model(&checkins.AttendanceAgency{}).Where("id in ?", IDs).Update("deleted_by", deleted_by).Error; err != nil {
            return err
        }
        if err := tx.Where("id in ?", IDs).Delete(&checkins.AttendanceAgency{}).Error; err != nil {
            return err
        }
        return nil
    })
	return err
}

// UpdateAttendanceAgency 更新Đơn vị记录
// Author [piexlmax](https://github.com/piexlmax)
func (attendanceAgencyService *AttendanceAgencyService)UpdateAttendanceAgency(attendanceAgency checkins.AttendanceAgency) (err error) {
	err = global.GVA_DB.Model(&checkins.AttendanceAgency{}).Where("id = ?",attendanceAgency.ID).Updates(&attendanceAgency).Error
	return err
}

// GetAttendanceAgency 根据ID获取Đơn vị记录
// Author [piexlmax](https://github.com/piexlmax)
func (attendanceAgencyService *AttendanceAgencyService)GetAttendanceAgency(ID string) (attendanceAgency checkins.AttendanceAgency, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&attendanceAgency).Error
	return
}

// GetAttendanceAgencyInfoList 分页获取Đơn vị记录
// Author [piexlmax](https://github.com/piexlmax)
func (attendanceAgencyService *AttendanceAgencyService)GetAttendanceAgencyInfoList(info checkinsReq.AttendanceAgencySearch) (list []checkins.AttendanceAgency, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&checkins.AttendanceAgency{})
    var attendanceAgencys []checkins.AttendanceAgency
    // 如果有条件搜索 下方会自动创建搜索语句
    if info.StartCreatedAt !=nil && info.EndCreatedAt !=nil {
     db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
    }
    if info.Name != "" {
        db = db.Where("name = ?",info.Name)
    }
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }

	if limit != 0 {
       db = db.Limit(limit).Offset(offset)
    }
	
	err = db.Find(&attendanceAgencys).Error
	return  attendanceAgencys, total, err
}