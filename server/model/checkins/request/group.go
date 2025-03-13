package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type GroupSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	AttendanceId   *int       `json:"attendanceId" form:"attendanceId" `
	request.PageInfo
}

type GroupAuto struct {
	GroupQty      int    `json:"groupQty" form:"groupQty"`
	GroupNameType string `json:"groupNameType" form:"groupNameType"`
	AttendanceId  int    `json:"attendanceId" form:"attendanceId"`
}

type GroupMember struct {
	GroupId       uint   `json:"groupId" form:"groupId" gorm:"column:group_id;comment:'ID của nhóm';"`
	ParticipantId uint   `json:"participantId" form:"participantId" gorm:"column:participant_id;comment:'ID của sinh viên';" binding:"required"`
	AttendanceId  uint   `json:"attendanceId,string" form:"attendanceId" gorm:"column:attendance_id;comment:'ID của phiên điểm danh';"`
	Name          string `json:"name" form:"name"`
}
