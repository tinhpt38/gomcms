// 自动生成模板Group
package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Nhóm 结构体  Group
type Group struct {
	global.GVA_MODEL
	Name           string        `json:"name" form:"name" gorm:"column:name;comment:;" binding:"required"`       //Tên nhóm
	AttendanceId   uint          `json:"attendanceId" form:"attendanceId" gorm:"column:attendance_id;comment:;"` //Attendance Class
	Total          int           `json:"total" gorm:"-"`                                                         //Tổng số người tham gia
	Participants   []Participant `json:"participants" gorm:"many2many:attendance_group_participants;"`           // Danh sách thành viên trong nhóm
	ParticipantIDs []uint        `json:"participantIds" form:"participantIds" gorm:"-"`                          // Danh sách ID thành viên được tag từ giao diện (không lưu DB)
}

// TableName Nhóm Group自定义表名 groups
func (Group) TableName() string {
	return "groups"
}

// GroupMember đại diện cho thành viên trong nhóm (hoặc sinh viên tham gia điểm danh)
type GroupMember struct {
	global.GVA_MODEL
	ParticipantId uint `json:"participantId" form:"participantId" gorm:"column:participant_id;comment:'ID của sinh viên';" binding:"required"`
	GroupId       uint `json:"groupId" form:"groupId" gorm:"column:group_id;comment:'ID của nhóm';"`
}

// TableName chỉ định tên bảng trong database
func (GroupMember) TableName() string {
	return "attendance_group_participants"
}
