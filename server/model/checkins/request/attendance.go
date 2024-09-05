package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type AttendanceSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	StartDate      *time.Time `json:"startDate" form:"startDate"`
	EndDate        *time.Time `json:"endDate" form:"endDate"`
	request.PageInfo
}

type AttendanceAreaRequest struct {
	AttendanceID uint `json:"attendanceID" form:"attendanceID"`
	AreaID       uint `json:"areaID" form:"areaID"`
}

type AttendanceCloneRequest struct {
	AttendanceID uint `json:"attendanceID" form:"attendanceID"`
	WithData     bool `json:"withData" form:"withData"`
}
