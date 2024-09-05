package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type ParticipantSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	FullName       string     `json:"fullName" form:"fullName" `
	Email          string     `json:"email" form:"email" `
	AttendanceId   *uint      `json:"attendanceId" form:"attendanceId"`
	GroupId        *uint      `json:"groupId" form:"groupId"`
	request.PageInfo
}

type ListEmailParticipantRequest struct {
	List         []string `json:"list" form:"list"`
	AttendanceId *uint    `json:"attendanceId" form:"attendanceId"`
}
