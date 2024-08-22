package checkins

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct{ AttendanceClassRouter }

var attendanceClassApi = api.ApiGroupApp.CheckinsApiGroup.AttendanceClassApi
