package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type AttendanceCheckInSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	AttendanceId   *uint      `json:"attendanceId" form:"attendanceId"`
	Email          *string    `json:"email" form:"email"`
	request.PageInfo
}
