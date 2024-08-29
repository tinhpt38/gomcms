package checkins

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	AttendanceRouter
	GroupRouter
	AreaRouter
	ParticipantRouter
	ConditionRouter
	AttendanceCheckInRouter
	AttendanceCategoryRouter
	AttendanceAgencyRouter
}

var (
	attendanceApi         = api.ApiGroupApp.CheckinsApiGroup.AttendanceApi
	groupApi              = api.ApiGroupApp.CheckinsApiGroup.GroupApi
	areaApi               = api.ApiGroupApp.CheckinsApiGroup.AreaApi
	participantApi        = api.ApiGroupApp.CheckinsApiGroup.ParticipantApi
	conditionApi          = api.ApiGroupApp.CheckinsApiGroup.ConditionApi
	attendanceCheckInApi  = api.ApiGroupApp.CheckinsApiGroup.AttendanceCheckInApi
	attendanceCategoryApi = api.ApiGroupApp.CheckinsApiGroup.AttendanceCategoryApi
	attendanceAgencyApi   = api.ApiGroupApp.CheckinsApiGroup.AttendanceAgencyApi
)
