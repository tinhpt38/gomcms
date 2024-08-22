// 自动生成模板AttendanceClass
package checkins

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// AttendanceClass 结构体  AttendanceClass
type AttendanceClass struct {
	global.GVA_MODEL
	Title     string     `json:"title" form:"title" gorm:"column:title;comment:;" binding:"required"`              //Tiêu đề
	StartDate *time.Time `json:"startDate" form:"startDate" gorm:"column:start_date;comment:;" binding:"required"` //Ngày bắt đầu
	EndDate   *time.Time `json:"endDate" form:"endDate" gorm:"column:end_date;comment:;" binding:"required"`       //Ngày kết thúc
	IsTrial   *bool      `json:"isTrial" form:"isTrial" gorm:"column:is_trial;comment:;"`                          //Cho thử nghiệm
	IsLocked  *bool      `json:"isLocked" form:"isLocked" gorm:"column:is_locked;comment:;"`                       //Khoá
	CreatedBy uint       `gorm:"column:created_by;comment:创建者"`
	UpdatedBy uint       `gorm:"column:updated_by;comment:更新者"`
	DeletedBy uint       `gorm:"column:deleted_by;comment:删除者"`
}

// TableName AttendanceClass AttendanceClass自定义表名 attendance_class
func (AttendanceClass) TableName() string {
	return "attendance_class"
}
