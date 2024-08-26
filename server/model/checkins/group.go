// 自动生成模板Group
package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Nhóm 结构体  Group
type Group struct {
	global.GVA_MODEL
	Name         string `json:"name" form:"name" gorm:"column:name;comment:;" binding:"required"`       //Tên nhóm
	AttendanceId *int   `json:"attendanceId" form:"attendanceId" gorm:"column:attendance_id;comment:;"` //Attendance Class
	Total        int    `json:"total" gorm:"-"`                                                         //Tổng số người tham gia
}

// TableName Nhóm Group自定义表名 groups
func (Group) TableName() string {
	return "groups"
}
