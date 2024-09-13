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
	AgencyId       uint       `json:"agencyId" form:"agencyId"`
	CategoryId     uint       `json:"categoryId" form:"categoryId"`
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

type StatsByInfoRequest struct {
	AgencyId   uint       `json:"agencyId" form:"agencyId"`
	CategoryId uint       `json:"categoryId" form:"categoryId"`
	StartAt    *time.Time `json:"startAt" form:"startAt"`
	EndAt      *time.Time `json:"endAt" form:"endAt"`
}

type AttendanceSearchHistory struct {
	request.PageInfo
	Email string `json:"email" form:"email"`
}
