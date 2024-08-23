package checkins

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	AttendanceClassRouter
	GroupRouter
	AreaRouter
}

var (
	attendanceClassApi = api.ApiGroupApp.CheckinsApiGroup.AttendanceClassApi
	groupApi           = api.ApiGroupApp.CheckinsApiGroup.GroupApi
	areaApi            = api.ApiGroupApp.CheckinsApiGroup.AreaApi
)
