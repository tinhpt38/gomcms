package checkins

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct{ AttendanceClassApi }

var attendanceClassService = service.ServiceGroupApp.CheckinsServiceGroup.AttendanceClassService
