package checkins

import (
	"encoding/base32"
	"errors"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
	checkinsReq "github.com/flipped-aurora/gin-vue-admin/server/model/checkins/request"
	"gorm.io/gorm"
)

type AttendanceCheckInService struct{}

func (attendanceCheckInService *AttendanceCheckInService) CreateAttendanceCheckIn(attendanceCheckIn *checkins.AttendanceCheckIn) (err error) {
	err = global.GVA_DB.Create(attendanceCheckIn).Error
	return err
}

func (attendanceCheckInService *AttendanceCheckInService) DeleteAttendanceCheckIn(ID string, userID uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&checkins.AttendanceCheckIn{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
			return err
		}
		if err = tx.Delete(&checkins.AttendanceCheckIn{}, "id = ?", ID).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (attendanceCheckInService *AttendanceCheckInService) DeleteAttendanceCheckInByIds(IDs []string, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&checkins.AttendanceCheckIn{}).Where("id in ?", IDs).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", IDs).Delete(&checkins.AttendanceCheckIn{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (attendanceCheckInService *AttendanceCheckInService) UpdateAttendanceCheckIn(attendanceCheckIn checkins.AttendanceCheckIn) (err error) {
	err = global.GVA_DB.Model(&checkins.AttendanceCheckIn{}).Where("id = ?", attendanceCheckIn.ID).Updates(&attendanceCheckIn).Error
	return err
}

func (attendanceCheckInService *AttendanceCheckInService) GetAttendanceCheckIn(ID string) (attendanceCheckIn checkins.AttendanceCheckIn, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&attendanceCheckIn).Error
	return
}

func (attendanceCheckInService *AttendanceCheckInService) GetAttendanceCheckInInfoList(info checkinsReq.AttendanceCheckInSearch) (list []checkins.AttendanceCheckIn, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.GVA_DB.Model(&checkins.AttendanceCheckIn{})
	var attendanceCheckIns []checkins.AttendanceCheckIn

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

	err = db.Find(&attendanceCheckIns).Error
	return attendanceCheckIns, total, err
}

func (attendanceCheckInService *AttendanceCheckInService) CheckinAttendance(req checkinsReq.CheckinsReq) (err error) {
	atId, derr := attendanceCheckInService.DecodeBase32(req.Code)
	if derr != nil {
		return errors.New("Mã QR không hợp lệ")
	}
	attendanceService := new(AttendanceService)
	attendance, grr := attendanceService.GetAttendance(atId)
	if grr != nil {
		return errors.New("Không tìm thấy thông tin điểm danh")
	}
	fmt.Println(attendance)

	return
}

func (attendanceCheckInService *AttendanceCheckInService) DecodeBase32(encoded string) (string, error) {
	decoded, err := base32.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
