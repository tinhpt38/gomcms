package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
    checkinsReq "github.com/flipped-aurora/gin-vue-admin/server/model/checkins/request"
    "gorm.io/gorm"
)

type AttendanceClassService struct {}


func (attendanceClassService *AttendanceClassService) CreateAttendanceClass(attendanceClass *checkins.AttendanceClass) (err error) {
	err = global.GVA_DB.Create(attendanceClass).Error
	return err
}


func (attendanceClassService *AttendanceClassService)DeleteAttendanceClass(ID string,userID uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
	    if err := tx.Model(&checkins.AttendanceClass{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
              return err
        }
        if err = tx.Delete(&checkins.AttendanceClass{},"id = ?",ID).Error; err != nil {
              return err
        }
        return nil
	})
	return err
}


func (attendanceClassService *AttendanceClassService)DeleteAttendanceClassByIds(IDs []string,deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
	    if err := tx.Model(&checkins.AttendanceClass{}).Where("id in ?", IDs).Update("deleted_by", deleted_by).Error; err != nil {
            return err
        }
        if err := tx.Where("id in ?", IDs).Delete(&checkins.AttendanceClass{}).Error; err != nil {
            return err
        }
        return nil
    })
	return err
}


func (attendanceClassService *AttendanceClassService)UpdateAttendanceClass(attendanceClass checkins.AttendanceClass) (err error) {
	err = global.GVA_DB.Model(&checkins.AttendanceClass{}).Where("id = ?",attendanceClass.ID).Updates(&attendanceClass).Error
	return err
}


func (attendanceClassService *AttendanceClassService)GetAttendanceClass(ID string) (attendanceClass checkins.AttendanceClass, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&attendanceClass).Error
	return
}


func (attendanceClassService *AttendanceClassService)GetAttendanceClassInfoList(info checkinsReq.AttendanceClassSearch) (list []checkins.AttendanceClass, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    
	db := global.GVA_DB.Model(&checkins.AttendanceClass{})
    var attendanceClasss []checkins.AttendanceClass
    
    if info.StartCreatedAt !=nil && info.EndCreatedAt !=nil {
     db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
    }
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }

	if limit != 0 {
       db = db.Limit(limit).Offset(offset)
    }
	
	err = db.Find(&attendanceClasss).Error
	return  attendanceClasss, total, err
}