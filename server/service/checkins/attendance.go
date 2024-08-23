package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
	checkinsReq "github.com/flipped-aurora/gin-vue-admin/server/model/checkins/request"
	"gorm.io/gorm"
)

type AttendanceService struct{}

func (attendanceService *AttendanceService) CreateAttendance(attendanceClass *checkins.Attendance) (err error) {
	err = global.GVA_DB.Create(attendanceClass).Error
	return err
}

func (attendanceService *AttendanceService) CreateAttendanceArea(attendanceArea *checkins.AttendanceArea) (err error) {
	err = global.GVA_DB.Create(attendanceArea).Error
	return err
}

func (attendanceService *AttendanceService) DeleteAttendanceArea(id string) (err error) {
	err = global.GVA_DB.Delete(&checkins.AttendanceArea{}, "id = ? ", id).Error
	return err
}

func (attendanceService *AttendanceService) DeleteAttendance(ID string, userID uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&checkins.Attendance{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
			return err
		}
		if err = tx.Delete(&checkins.Attendance{}, "id = ?", ID).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (attendanceService *AttendanceService) DeleteAttendanceByIds(IDs []string, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&checkins.Attendance{}).Where("id in ?", IDs).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", IDs).Delete(&checkins.Attendance{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (attendanceService *AttendanceService) UpdateAttendance(attendanceClass checkins.Attendance) (err error) {
	err = global.GVA_DB.Model(&checkins.Attendance{}).Where("id = ?", attendanceClass.ID).Updates(&attendanceClass).Error
	return err
}

func (attendanceService *AttendanceService) GetAttendance(ID string) (attendanceClass checkins.Attendance, err error) {
	err = global.GVA_DB.Where("id = ?", ID).Preload("Areas").First(&attendanceClass).Error
	return
}

func (attendanceService *AttendanceService) GetAttendanceArea(id uint) (list []checkins.AttendanceArea, err error) {

	db := global.GVA_DB.Table(checkins.AttendanceArea{}.TableName())
	err = db.Where("attendance_id = ?", id).Preload("Area").Find(&list).Error
	return
}

func (attendanceService *AttendanceService) GetAttendanceInfoList(info checkinsReq.AttendanceSearch) (list []checkins.Attendance, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.GVA_DB.Model(&checkins.Attendance{})
	var attendanceClasss []checkins.Attendance

	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Preload("Areas").Debug().Find(&attendanceClasss).Error
	return attendanceClasss, total, err
}
