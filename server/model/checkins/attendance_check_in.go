// 自动生成模板AttendanceCheckIn
package checkins

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Thành viên checkins 结构体  AttendanceCheckIn
type AttendanceCheckIn struct {
	global.GVA_MODEL
	CheckinDate      time.Time    `json:"checkinDate" form:"checkinDate" gorm:"column:checkin_date;comment:;" binding:"required"`                //Ngày giờ
	AttendanceId     *uint        `json:"attendanceId" form:"attendanceId" gorm:"column:attendance_id;comment:;" binding:"required"`             //Lớp điểm danh
	PartpaticipantId *uint        `json:"partpaticipantId" form:"partpaticipantId" gorm:"column:partpaticipant_id;comment:;" binding:"required"` //Thành viên
	AreaId           *uint        `json:"areaId" form:"areaId" gorm:"column:area_id;comment:;"`                                                  //Khu vực
	GroupId          *uint        `json:"groupId" form:"groupId" gorm:"column:group_id;comment:;"`                                               //Nhóm
	ConditionId      *uint        `json:"conditionId" form:"conditionId" gorm:"column:condition_id;comment:;"`                                   //Điều kiện
	IP               string       `json:"iP" form:"iP" gorm:"column:ip;comment:;"`                                                               //IP
	Lattidue         *float64     `json:"lattidue" form:"lattidue" gorm:"column:lattidue;comment:;"`                                             //Kinh độ
	Longtidue        *float64     `json:"longtidue" form:"longtidue" gorm:"column:longtidue;comment:;"`                                          //Vĩ độ
	Agent            string       `json:"agent" form:"agent" gorm:"column:agent;comment:;type:text;`                                             //Agent
	Browser          string       `json:"browser" form:"browser" gorm:"column:browser;comment:;"`                                                //Trình duyệt
	CreatedBy        uint         `gorm:"column:created_by;comment:创建者"`
	UpdatedBy        uint         `gorm:"column:updated_by;comment:更新者"`
	DeletedBy        uint         `gorm:"column:deleted_by;comment:删除者"`
	Group            *Group       `json:"group" gorm:"foreignKey:GroupId;references:ID;comment:Nhóm;"`                      //Nhóm
	Area             *Area        `json:"area" gorm:"foreignKey:AreaId;references:ID;comment:Khu vực;"`                     //Khu vực
	Condition        *Condition   `json:"condition" gorm:"foreignKey:ConditionId;references:ID;comment:Điều kiện;"`         //Điều kiện
	Attendance       *Attendance  `json:"attendance" gorm:"foreignKey:AttendanceId;references:ID;comment:Lớp điểm danh;"`   //Lớp điểm danh
	Participant      *Participant `json:"participant" gorm:"foreignKey:PartpaticipantId;references:ID;comment:Thành viên;"` //Thành viên
}

// TableName Thành viên checkins AttendanceCheckIn自定义表名 attendanceCheckIns
func (AttendanceCheckIn) TableName() string {
	return "attendance_checkins"
}
