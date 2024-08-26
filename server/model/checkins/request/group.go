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
