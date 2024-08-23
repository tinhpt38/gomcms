// 自动生成模板Group
package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Nhóm 结构体  Group
type Group struct {
    global.GVA_MODEL
    Name  string `json:"name" form:"name" gorm:"column:name;comment:;" binding:"required"`  //Tên nhóm 
    AttendanceClassId  *int `json:"attendanceClassId" form:"attendanceClassId" gorm:"column:attendance_class_id;comment:;"`  //Attendance Class 
}


// TableName Nhóm Group自定义表名 groups
func (Group) TableName() string {
    return "groups"
}

