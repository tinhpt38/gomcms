package checkins

import (
	"math"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
	overviewResp "github.com/flipped-aurora/gin-vue-admin/server/model/checkins/response"
)

const (
	maxCheckinPoints = 2000
	maxFailPoints    = 500
)

// haversineM returns the distance in metres between two lat/lng points.
func haversineM(lat1, lng1, lat2, lng2 float64) float64 {
	const r = 6371000 // Earth radius in metres
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	return r * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

func (attendanceService *AttendanceService) GetAttendanceOverview(attendanceId uint) (result overviewResp.AttendanceOverviewResponse, err error) {
	db := global.GVA_DB

	// ── Session summary ──────────────────────────────────────────────────────
	var att checkins.Attendance
	if err = db.First(&att, attendanceId).Error; err != nil {
		return
	}
	elapsed := 0.0
	if att.StartDate != nil && att.EndDate != nil {
		total := att.EndDate.Sub(*att.StartDate).Seconds()
		if total > 0 {
			passed := time.Since(*att.StartDate).Seconds()
			elapsed = math.Max(0, math.Min(100, passed/total*100))
		}
	}
	result.Session = overviewResp.SessionSummary{
		Title:          att.Title,
		StartDate:      att.StartDate,
		EndDate:        att.EndDate,
		Total:          att.Total,
		TotalCheckin:   att.TotalCheckin,
		ElapsedPercent: elapsed,
	}

	// ── Geofences ────────────────────────────────────────────────────────────
	type gfRow struct {
		Name    string
		Lat     *float64
		Lng     *float64
		RadiusM *float64
	}
	var gfRows []gfRow
	db.Raw(`
		SELECT a.name, a.latitude AS lat, a.longitude AS lng,
		       COALESCE(aa.radius, a.radius) AS radius_m
		FROM attendance_areas aa
		JOIN areas a ON a.id = aa.area_id
		WHERE aa.attendance_id = ? AND aa.deleted_at IS NULL AND a.deleted_at IS NULL`,
		attendanceId).Scan(&gfRows)
	for _, g := range gfRows {
		if g.Lat == nil || g.Lng == nil {
			continue
		}
		r := 0.0
		if g.RadiusM != nil {
			r = *g.RadiusM
		}
		result.Geofences = append(result.Geofences, overviewResp.GeofenceItem{
			Name: g.Name, Lat: *g.Lat, Lng: *g.Lng, RadiusM: r,
		})
	}

	// ── Check-in points (OK) ─────────────────────────────────────────────────
	type checkinRow struct {
		Lat      *float64
		Lng      *float64
		Accuracy *float64
		Email    *string
		FullName *string
		At       time.Time
	}
	var checkinRows []checkinRow
	db.Raw(`
		SELECT ac.lattidue AS lat, ac.longtidue AS lng, ac.accuracy,
		       p.email, p.full_name, ac.checkin_date AS at
		FROM attendance_checkins ac
		LEFT JOIN participants p ON p.id = ac.partpaticipant_id
		WHERE ac.attendance_id = ? AND ac.deleted_at IS NULL
		  AND ac.lattidue IS NOT NULL
		ORDER BY ac.id DESC
		LIMIT ?`, attendanceId, maxCheckinPoints).Scan(&checkinRows)

	result.MapTruncated = len(checkinRows) == maxCheckinPoints
	for _, r := range checkinRows {
		if r.Lat == nil || r.Lng == nil {
			continue
		}
		label := ""
		if r.FullName != nil && *r.FullName != "" {
			label = *r.FullName
		} else if r.Email != nil {
			label = *r.Email
		}
		acc := 0.0
		if r.Accuracy != nil {
			acc = *r.Accuracy
		}
		result.CheckinPoints = append(result.CheckinPoints, overviewResp.MapPoint{
			Lat: *r.Lat, Lng: *r.Lng, Accuracy: acc, Label: label,
			At: r.At.Format("02/01 15:04"),
		})
	}

	// ── Fail points ───────────────────────────────────────────────────────────
	type failRow struct {
		Lat   *float64
		Lng   *float64
		Email string
		At    time.Time
	}
	var failRows []failRow
	db.Raw(`
		SELECT lat, lng, email, created_at AS at
		FROM checkin_logs
		WHERE attendance_id = ? AND lat IS NOT NULL
		ORDER BY id DESC
		LIMIT ?`, attendanceId, maxFailPoints).Scan(&failRows)
	for _, r := range failRows {
		if r.Lat == nil || r.Lng == nil {
			continue
		}
		result.FailPoints = append(result.FailPoints, overviewResp.MapPoint{
			Lat: *r.Lat, Lng: *r.Lng, Label: r.Email,
			At: r.At.Format("02/01 15:04"),
		})
	}

	// ── Check-ins by hour ─────────────────────────────────────────────────────
	// If session spans more than 1 day use "dd/mm HH:00", else "HH:00"
	multiDay := att.StartDate != nil && att.EndDate != nil &&
		att.EndDate.Sub(*att.StartDate) > 24*time.Hour
	format := "%H:00"
	if multiDay {
		format = "%d/%m %H:00"
	}
	type hourRow struct {
		Label string
		Count int
	}
	var hourRows []hourRow
	db.Raw(`
		SELECT DATE_FORMAT(checkin_date, ?) AS label, COUNT(*) AS count
		FROM attendance_checkins
		WHERE attendance_id = ? AND deleted_at IS NULL
		GROUP BY DATE_FORMAT(checkin_date, ?)
		ORDER BY MIN(checkin_date)`, format, attendanceId, format).Scan(&hourRows)
	for _, r := range hourRows {
		result.CheckinsByHour = append(result.CheckinsByHour, overviewResp.HourBucket{Label: r.Label, Count: r.Count})
	}

	// ── Conditions timeline ───────────────────────────────────────────────────
	type condRow struct {
		ID        uint
		GroupName *string
		AreaName  *string
		StartAt   *time.Time
		EndAt     *time.Time
	}
	var condRows []condRow
	db.Raw(`
		SELECT c.id, g.name AS group_name, a.name AS area_name, c.start_at, c.end_at
		FROM conditions c
		LEFT JOIN groups g ON g.id = c.group_id
		LEFT JOIN attendance_areas aa ON aa.id = c.area_id
		LEFT JOIN areas a ON a.id = aa.area_id
		WHERE c.attendance_id = ? AND c.deleted_at IS NULL`, attendanceId).Scan(&condRows)
	for _, r := range condRows {
		slot := overviewResp.ConditionSlot{ID: r.ID, StartAt: r.StartAt, EndAt: r.EndAt}
		if r.GroupName != nil {
			slot.GroupName = *r.GroupName
		}
		if r.AreaName != nil {
			slot.AreaName = *r.AreaName
		}
		result.ConditionsTimeline = append(result.ConditionsTimeline, slot)
	}

	// ── Group stats ───────────────────────────────────────────────────────────
	type groupRow struct {
		GroupID         uint
		Name            string
		MemberCount     int
		CheckinDistinct int
	}
	var groupRows []groupRow
	db.Raw(`
		SELECT g.id AS group_id, g.name,
		       COUNT(DISTINCT agp.participant_id) AS member_count,
		       COUNT(DISTINCT ac.partpaticipant_id) AS checkin_distinct
		FROM groups g
		LEFT JOIN attendance_group_participants agp
		       ON agp.group_id = g.id AND agp.attendance_id = ? AND agp.deleted_at IS NULL
		LEFT JOIN attendance_checkins ac
		       ON ac.group_id = g.id AND ac.attendance_id = ? AND ac.deleted_at IS NULL
		WHERE g.attendance_id = ? AND g.deleted_at IS NULL
		GROUP BY g.id, g.name`, attendanceId, attendanceId, attendanceId).Scan(&groupRows)
	for _, r := range groupRows {
		result.GroupStats = append(result.GroupStats, overviewResp.GroupProgress{
			GroupID: r.GroupID, Name: r.Name,
			MemberCount: r.MemberCount, CheckinDistinct: r.CheckinDistinct,
		})
	}

	// ── Coverage summary ──────────────────────────────────────────────────────
	type coverageRow struct {
		FullyMapped int
		HasAny      int
		Total       int
	}
	var cov coverageRow
	db.Raw(`
		SELECT
		  SUM(CASE WHEN ac_count > 0 AND ac_count >= agp_count THEN 1 ELSE 0 END) AS fully_mapped,
		  SUM(CASE WHEN ac_count > 0 AND ac_count < agp_count  THEN 1 ELSE 0 END) AS has_any,
		  COUNT(*) AS total
		FROM (
		  SELECT agp.participant_id,
		         COUNT(DISTINCT agp.id)              AS agp_count,
		         COUNT(DISTINCT agpc.id)             AS ac_count
		  FROM attendance_group_participants agp
		  LEFT JOIN agp_conditions agpc
		         ON agpc.agp_id = agp.id AND agpc.deleted_at IS NULL
		  WHERE agp.attendance_id = ? AND agp.deleted_at IS NULL
		  GROUP BY agp.participant_id
		) t`, attendanceId).Scan(&cov)
	result.Coverage = overviewResp.CoverageSummary{
		FullyMapped: cov.FullyMapped,
		Partial:     cov.HasAny,
		NoMapping:   cov.Total - cov.FullyMapped - cov.HasAny,
	}

	// ── Quality scatter ───────────────────────────────────────────────────────
	// For each check-in with GPS: compute distance to the check-in's area (or closest geofence)
	type qualRaw struct {
		Lat      *float64
		Lng      *float64
		Accuracy *float64
		AreaLat  *float64
		AreaLng  *float64
		RadiusM  *float64
		AreaName *string
	}
	var qualRows []qualRaw
	db.Raw(`
		SELECT ac.lattidue AS lat, ac.longtidue AS lng, ac.accuracy,
		       ar.latitude AS area_lat, ar.longitude AS area_lng,
		       COALESCE(aa.radius, ar.radius) AS radius_m, ar.name AS area_name
		FROM attendance_checkins ac
		LEFT JOIN attendance_areas aa ON aa.id = ac.area_id AND aa.deleted_at IS NULL
		LEFT JOIN areas ar ON ar.id = aa.area_id AND ar.deleted_at IS NULL
		WHERE ac.attendance_id = ? AND ac.deleted_at IS NULL AND ac.lattidue IS NOT NULL
		LIMIT 500`, attendanceId).Scan(&qualRows)
	for _, r := range qualRows {
		if r.Lat == nil || r.Lng == nil {
			continue
		}
		acc := 0.0
		if r.Accuracy != nil {
			acc = *r.Accuracy
		}
		dist := 0.0
		outOfRadius := false
		areaName := ""
		if r.AreaLat != nil && r.AreaLng != nil {
			dist = haversineM(*r.Lat, *r.Lng, *r.AreaLat, *r.AreaLng)
			if r.RadiusM != nil && *r.RadiusM > 0 {
				outOfRadius = dist > *r.RadiusM
			}
		}
		if r.AreaName != nil {
			areaName = *r.AreaName
		}
		result.QualityPoints = append(result.QualityPoints, overviewResp.QualityPoint{
			Lat: *r.Lat, Lng: *r.Lng, Accuracy: acc,
			DistanceM: math.Round(dist), OutOfRadius: outOfRadius, AreaName: areaName,
		})
	}

	// ── Top IPs ───────────────────────────────────────────────────────────────
	type ipRow struct {
		IP    string
		Count int
	}
	var ipRows []ipRow
	db.Raw(`
		SELECT ip, COUNT(*) AS count
		FROM attendance_checkins
		WHERE attendance_id = ? AND deleted_at IS NULL AND ip != ''
		GROUP BY ip
		ORDER BY count DESC
		LIMIT 10`, attendanceId).Scan(&ipRows)
	for _, r := range ipRows {
		result.TopIPs = append(result.TopIPs, overviewResp.IPCount{IP: r.IP, Count: r.Count})
	}

	// ── Shared visitors ───────────────────────────────────────────────────────
	type visitorRow struct {
		VisitorID    string
		AttemptCount int
		Emails       string // comma-joined
	}
	var visitorRows []visitorRow
	db.Raw(`
		SELECT ac.visitor_id, COUNT(DISTINCT ac.partpaticipant_id) AS attempt_count,
		       GROUP_CONCAT(DISTINCT p.email) AS emails
		FROM attendance_checkins ac
		LEFT JOIN participants p ON p.id = ac.partpaticipant_id
		WHERE ac.attendance_id = ? AND ac.deleted_at IS NULL AND ac.visitor_id != ''
		GROUP BY ac.visitor_id
		HAVING COUNT(DISTINCT ac.partpaticipant_id) > 1
		LIMIT 20`, attendanceId).Scan(&visitorRows)
	for _, r := range visitorRows {
		emails := strings.Split(r.Emails, ",")
		result.SharedVisitors = append(result.SharedVisitors, overviewResp.VisitorAlert{
			VisitorID: r.VisitorID, AttemptCount: r.AttemptCount, Emails: emails,
		})
	}

	// ── Fail reasons ──────────────────────────────────────────────────────────
	type msgRow struct {
		MessageList string
	}
	var msgRows []msgRow
	db.Raw(`
		SELECT message_list FROM checkin_logs
		WHERE attendance_id = ? AND message_list != ''`, attendanceId).Scan(&msgRows)
	reasonMap := map[string]int{}
	for _, r := range msgRows {
		parts := strings.Split(r.MessageList, "$$")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				reasonMap[p]++
			}
		}
	}
	// Return top 10 sorted by count
	type kv struct {
		k string
		v int
	}
	var sorted []kv
	for k, v := range reasonMap {
		sorted = append(sorted, kv{k, v})
	}
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].v > sorted[i].v {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	if len(sorted) > 10 {
		sorted = sorted[:10]
	}
	for _, kv := range sorted {
		result.FailReasons = append(result.FailReasons, overviewResp.ReasonCount{Reason: kv.k, Count: kv.v})
	}

	return
}
