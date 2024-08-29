// 自动生成模板AttendanceAgency
package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Đơn vị 结构体  AttendanceAgency
type AttendanceAgency struct {
    global.GVA_MODEL
    Name  string `json:"name" form:"name" gorm:"column:name;comment:;" binding:"required"`  //Đơn vị 
    ParentId  *int `json:"parentId" form:"parentId" gorm:"column:parent_id;comment:;"`  //Đơn vị cha 
    CreatedBy  uint   `gorm:"column:created_by;comment:创建者"`
    UpdatedBy  uint   `gorm:"column:updated_by;comment:更新者"`
    DeletedBy  uint   `gorm:"column:deleted_by;comment:删除者"`
}


// TableName Đơn vị AttendanceAgency自定义表名 attendance_agencies
func (AttendanceAgency) TableName() string {
    return "attendance_agencies"
}

