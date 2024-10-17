// 自动生成模板Participant
package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Sinh viên (Người tham dự phiên điểm danh) 结构体  Participant
type Participant struct {
	global.GVA_MODEL
	Email        string  `json:"email" form:"email" gorm:"column:email;comment:Email;" binding:"required"` //Email
	FullName     *string `json:"fullName" form:"fullName" gorm:"column:full_name;comment:Họ và tên;"`      //Họ và tên
	GroupId      *uint   `json:"groupId" form:"groupId" gorm:"-"`
	AttendanceId *uint   `json:"attendanceId" form:"attendanceId" gorm:"-"`
	Groups       []Group `json:"groups" gorm:"many2many:attendance_group_participants;"`
}

// TableName Sinh viên (Người tham dự phiên điểm danh) Participant自定义表名 participant
func (Participant) TableName() string {
	return "participants"
}

type AttendanceGroupParticipant struct {
	global.GVA_MODEL
	ParticipantId *uint          `json:"participantId" form:"participantId" gorm:"column:participant_id;comment:Sinh viên"`
	AttendanceId  *uint          `json:"attendanceId" form:"attendanceId" gorm:"column:attendance_id;comment:Phiên điểm danh"`
	GroupId       *uint          `json:"groupId" form:"groupId" gorm:"column:group_id;comment:Nhóm"`
	Participant   *Participant   `gorm:"foreignKey:ParticipantId"` // Thiết lập khóa ngoại đến struct Participant
	Attendance    *Attendance    `gorm:"foreignKey:AttendanceId"`  // Thiết lập khóa ngoại đến struct Attendance
	Group         *Group         `gorm:"foreignKey:GroupId"`       // Thiết lập khóa ngoại đến struct Group
	AGPCondition  []AGPCondition `json:"groups" gorm:"many2many:agp_conditions;"`
}

func (AttendanceGroupParticipant) TableName() string {
	return "attendance_group_participants"
}

type AGPCondition struct {
	global.GVA_MODEL
	AttendanceGroupParticipantId *uint `json:"attendanceGroupParticipantId" form:"attendanceGroupParticipantId" gorm:"column:agp_id;comment:Phiên điểm danh"`
	ConditionId                  *int  `json:"conditionId" form:"conditionId" gorm:"column:condition_id;comment:Điều kiện"`
	AttendanceId                 *int  `json:"attendanceId" form:"attendanceId" gorm:"column:attendance_id;comment:Phiên điểm danh"`
}

func (AGPCondition) TableName() string {
	return "agp_conditions"
}
