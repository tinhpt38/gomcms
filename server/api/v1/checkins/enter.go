package checkins

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	AttendanceClassApi
	GroupApi
}

var (
	attendanceClassService = service.ServiceGroupApp.CheckinsServiceGroup.AttendanceClassService
	groupService           = service.ServiceGroupApp.CheckinsServiceGroup.GroupService
)
