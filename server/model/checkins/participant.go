// 自动生成模板Participant
package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Sinh viên (Người tham dự phiên điểm danh) 结构体  Participant
type Participant struct {
	global.GVA_MODEL
	FullName string `json:"fullName" form:"fullName" gorm:"column:full_name;comment:Họ và tên;" binding:"required"` //Họ và tên
	Email    string `json:"email" form:"email" gorm:"column:email;comment:Email;" binding:"required"`               //Email
	// GroupId  uint   `json:"groupId" form:"groupId" gorm:"column:group_id;comment:Mã nhóm;" binding:"required"`      //Mã nhóm
	// Group Group `json:"group" gorm:"foreignKey:GroupId;references:ID;comment:Nhóm;"` //Nhóm
	
}

// TableName Sinh viên (Người tham dự phiên điểm danh) Participant自定义表名 participant
func (Participant) TableName() string {
	return "participants"
}

type ParticipantGroup struct {
	global.GVA_MODEL
	ParticipantId *uint        `json:"participantId" form:"participantId" gorm:"column:participant_id;comment:Sinh viên"`
	GroupId       *uint        `json:"groupId" form:"groupId" gorm:"column:group_id;comment:Nhóm"`
	Participant   *Participant `gorm:"foreignKey:ParticipantId"` // Thiết lập khóa ngoại đến struct Participant
	Group         *Group       `gorm:"foreignKey:GroupId"`       // Thiết lập khóa ngoại đến struct Group
}

func (ParticipantGroup) TableName() string {
	return "participant_groups"
}

type ParticipantAttendance struct {
	global.GVA_MODEL
	ParticipantId *uint        `json:"participantId" form:"participantId" gorm:"column:participant_id;comment:Sinh viên"`
	AttendanceId  *uint        `json:"attendanceId" form:"attendanceId" gorm:"column:attendance_id;comment:Phiên điểm danh"`
	Participant   *Participant `gorm:"foreignKey:ParticipantId"` // Thiết lập khóa ngoại đến struct Participant
	Attendance    *Attendance  `gorm:"foreignKey:AttendanceId"`  // Thiết lập khóa ngoại đến struct Attendance
}

func (ParticipantAttendance) TableName() string {
	return "participant_attendances"
}
