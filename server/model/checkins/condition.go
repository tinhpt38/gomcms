// 自动生成模板Condition
package checkins

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Điều kiện để checkins 结构体  Condition
type Condition struct {
	global.GVA_MODEL
	AttendanceId *uint           `json:"attendanceId" form:"attendanceId" gorm:"column:attendance_id;comment:;"` //Khu vực
	GroupId      *uint           `json:"groupId" form:"groupId" gorm:"column:group_id;comment:;"`                //Nhóm
	AreaId       *uint           `json:"areaId" form:"areaId" gorm:"column:area_id;comment:;"`                   //Khu vực
	StartAt      *time.Time      `json:"startAt" form:"startAt" gorm:"column:start_at;comment:;"`                //Bắt đầu
	EndAt        *time.Time      `json:"endAt" form:"endAt" gorm:"column:end_at;comment:;"`                      //Kết thúc
	CreatedBy    uint            `gorm:"column:created_by;comment:创建者"`
	UpdatedBy    uint            `gorm:"column:updated_by;comment:更新者"`
	DeletedBy    uint            `gorm:"column:deleted_by;comment:删除者"`
	Group        *Group          `json:"group" gorm:"foreignKey:GroupId;references:ID;comment:Nhóm;"`  //Nhóm
	Area         *AttendanceArea `json:"area" gorm:"foreignKey:AreaId;references:ID;comment:Khu vực;"` //Khu vực
}

// TableName Điều kiện để checkins Condition自定义表名 conditions
func (Condition) TableName() string {
	return "conditions"
}
