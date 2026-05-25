package checkins

import (
	"testing"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
	"github.com/stretchr/testify/assert"
)

func cleanupOverviewTestData(t *testing.T) {
	t.Helper()
	tables := []string{
		"attendance_checkins",
		"checkin_logs",
		"agp_conditions",
		"attendance_areas",
		"conditions",
		"attendance_group_participants",
		"participants",
		"groups",
		"areas",
		"attendances",
	}
	for _, table := range tables {
		global.GVA_DB.Exec("DELETE FROM " + table)
	}
}

// TestGetAttendanceOverview_Geofences verifies geofence rows are returned.
func TestGetAttendanceOverview_Geofences(t *testing.T) {
	cleanupOverviewTestData(t)
	svc := AttendanceService{}

	attId := uint(9001)
	att := checkins.Attendance{Title: "Overview Geofence Test"}
	att.ID = attId
	global.GVA_DB.Create(&att)

	lat, lng, radius := 15.97, 108.24, float64(100)
	area := checkins.Area{Name: "Khu A", Latitude: &lat, Longitude: &lng, Radius: &radius}
	global.GVA_DB.Create(&area)

	attArea := checkins.AttendanceArea{AttendanceID: &attId, AreaID: &area.ID, Radius: &radius}
	global.GVA_DB.Create(&attArea)

	res, err := svc.GetAttendanceOverview(attId)
	assert.NoError(t, err)
	assert.Len(t, res.Geofences, 1)
	assert.Equal(t, "Khu A", res.Geofences[0].Name)
	assert.InDelta(t, radius, res.Geofences[0].RadiusM, 0.01)
}

// TestGetAttendanceOverview_HourBuckets verifies checkin-by-hour aggregation.
func TestGetAttendanceOverview_HourBuckets(t *testing.T) {
	cleanupOverviewTestData(t)
	svc := AttendanceService{}

	attId := uint(9002)
	att := checkins.Attendance{Title: "Overview Hour Test"}
	att.ID = attId
	global.GVA_DB.Create(&att)

	p := checkins.Participant{Email: "overview@test.com"}
	global.GVA_DB.Create(&p)

	now := time.Now()
	for i := 0; i < 3; i++ {
		t0 := now
		checkin := checkins.AttendanceCheckIn{
			CheckinDate:      t0,
			AttendanceId:     &attId,
			PartpaticipantId: &p.ID,
		}
		global.GVA_DB.Create(&checkin)
	}
	// One checkin 2 hours earlier
	earlier := now.Add(-2 * time.Hour)
	checkin := checkins.AttendanceCheckIn{
		CheckinDate:      earlier,
		AttendanceId:     &attId,
		PartpaticipantId: &p.ID,
	}
	global.GVA_DB.Create(&checkin)

	res, err := svc.GetAttendanceOverview(attId)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(res.CheckinsByHour), 2, "Phải có ít nhất 2 bucket giờ")
	total := 0
	for _, b := range res.CheckinsByHour {
		total += b.Count
	}
	assert.Equal(t, 4, total, "Tổng count phải bằng số check-in")
}

// TestGetAttendanceOverview_FailReasons verifies messageList parsing.
func TestGetAttendanceOverview_FailReasons(t *testing.T) {
	cleanupOverviewTestData(t)
	svc := AttendanceService{}

	attId := uint(9003)
	att := checkins.Attendance{Title: "Overview Fail Test"}
	att.ID = attId
	global.GVA_DB.Create(&att)

	msgs := []string{
		"Ngoài khu vực$$Hết giờ",
		"Ngoài khu vực",
		"Ngoài khu vực$$Thiết bị không hợp lệ",
	}
	for _, msg := range msgs {
		log := checkins.CheckinLog{
			Email:        "fail@test.com",
			Code:         "test",
			AttendanceId: attId,
			VisitorId:    "v1",
			MessageList:  msg,
		}
		global.GVA_DB.Create(&log)
	}

	res, err := svc.GetAttendanceOverview(attId)
	assert.NoError(t, err)
	assert.NotEmpty(t, res.FailReasons)
	// "Ngoài khu vực" should appear 3 times — first in the sorted list
	assert.Equal(t, "Ngoài khu vực", res.FailReasons[0].Reason)
	assert.Equal(t, 3, res.FailReasons[0].Count)
}
