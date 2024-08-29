package checkins

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	AttendanceApi
	GroupApi
	AreaApi
	ParticipantApi
	ConditionApi
	AttendanceCheckInApi
	AttendanceCategoryApi
	AttendanceAgencyApi
}

var (
	attendanceService         = service.ServiceGroupApp.CheckinsServiceGroup.AttendanceService
	groupService              = service.ServiceGroupApp.CheckinsServiceGroup.GroupService
	areaService               = service.ServiceGroupApp.CheckinsServiceGroup.AreaService
	participantService        = service.ServiceGroupApp.CheckinsServiceGroup.ParticipantService
	conditionService          = service.ServiceGroupApp.CheckinsServiceGroup.ConditionService
	attendanceCheckInService  = service.ServiceGroupApp.CheckinsServiceGroup.AttendanceCheckInService
	attendanceCategoryService = service.ServiceGroupApp.CheckinsServiceGroup.AttendanceCategoryService
	attendanceAgencyService   = service.ServiceGroupApp.CheckinsServiceGroup.AttendanceAgencyService
)
