// 自动生成模板AttendanceCategory
package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Danh mục điểm danh 结构体  AttendanceCategory
type AttendanceCategory struct {
	global.GVA_MODEL
	Name      string `json:"name" form:"name" gorm:"column:name;comment:;" binding:"required"` //Tên phân loại
	ParentId  *uint  `json:"parentId" form:"parentId" gorm:"column:parent_id;comment:;"`       //Danh mục cha
	IsCurrent bool   `json:"isCurrent" form:"isCurrent" gorm:"column:is_current;comment:;"`    //Đang sử dụng
	CreatedBy uint   `gorm:"column:created_by;comment:创建者"`
	UpdatedBy uint   `gorm:"column:updated_by;comment:更新者"`
	DeletedBy uint   `gorm:"column:deleted_by;comment:删除者"`
}

// TableName Danh mục điểm danh AttendanceCategory自定义表名 attendance_categories
func (AttendanceCategory) TableName() string {
	return "attendance_categories"
}
