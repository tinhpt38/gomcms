package checkins

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	AttendanceApi
	GroupApi
	AreaApi
}

var (
	attendanceService = service.ServiceGroupApp.CheckinsServiceGroup.AttendanceService
	groupService      = service.ServiceGroupApp.CheckinsServiceGroup.GroupService
	areaService       = service.ServiceGroupApp.CheckinsServiceGroup.AreaService
)
