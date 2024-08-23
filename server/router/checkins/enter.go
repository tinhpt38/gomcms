package checkins

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	AttendanceRouter
	GroupRouter
	AreaRouter
}

var (
	attendanceApi = api.ApiGroupApp.CheckinsApiGroup.AttendanceApi
	groupApi           = api.ApiGroupApp.CheckinsApiGroup.GroupApi
	areaApi            = api.ApiGroupApp.CheckinsApiGroup.AreaApi
)
