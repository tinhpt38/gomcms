package response

import "time"

// AttendanceOverviewResponse is the single aggregate response for the Tổng quan tab.
type AttendanceOverviewResponse struct {
	Session            SessionSummary    `json:"session"`
	Geofences          []GeofenceItem    `json:"geofences"`
	CheckinPoints      []MapPoint        `json:"checkinPoints"`
	FailPoints         []MapPoint        `json:"failPoints"`
	CheckinsByHour     []HourBucket      `json:"checkinsByHour"`
	ConditionsTimeline []ConditionSlot   `json:"conditionsTimeline"`
	GroupStats         []GroupProgress   `json:"groupStats"`
	Coverage           CoverageSummary   `json:"coverage"`
	QualityPoints      []QualityPoint    `json:"qualityPoints"`
	TopIPs             []IPCount         `json:"topIPs"`
	SharedVisitors     []VisitorAlert    `json:"sharedVisitors"`
	FailReasons        []ReasonCount     `json:"failReasons"`
	MapTruncated       bool              `json:"mapTruncated"`
}

type SessionSummary struct {
	Title          string     `json:"title"`
	StartDate      *time.Time `json:"startDate"`
	EndDate        *time.Time `json:"endDate"`
	Total          int        `json:"total"`
	TotalCheckin   int        `json:"totalCheckin"`
	ElapsedPercent float64    `json:"elapsedPercent"` // 0-100
}

type GeofenceItem struct {
	Name    string   `json:"name"`
	Lat     float64  `json:"lat"`
	Lng     float64  `json:"lng"`
	RadiusM float64  `json:"radiusM"`
}

type MapPoint struct {
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	Accuracy float64 `json:"accuracy"`
	Label    string  `json:"label"`  // email / name
	At       string  `json:"at"`     // formatted time
}

type HourBucket struct {
	Label string `json:"label"` // "14:00" or "25/05 14:00"
	Count int    `json:"count"`
}

type ConditionSlot struct {
	ID        uint       `json:"id"`
	GroupName string     `json:"groupName"`
	AreaName  string     `json:"areaName"`
	StartAt   *time.Time `json:"startAt"`
	EndAt     *time.Time `json:"endAt"`
}

type GroupProgress struct {
	GroupID            uint   `json:"groupId"`
	Name               string `json:"name"`
	MemberCount        int    `json:"memberCount"`
	CheckinDistinct    int    `json:"checkinDistinct"`
}

type CoverageSummary struct {
	FullyMapped    int `json:"fullyMapped"`
	Partial        int `json:"partial"`
	NoMapping      int `json:"noMapping"`
}

type QualityPoint struct {
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	Accuracy    float64 `json:"accuracy"`
	DistanceM   float64 `json:"distanceM"`
	OutOfRadius bool    `json:"outOfRadius"`
	AreaName    string  `json:"areaName"`
}

type IPCount struct {
	IP    string `json:"ip"`
	Count int    `json:"count"`
}

type VisitorAlert struct {
	VisitorID    string   `json:"visitorId"`
	AttemptCount int      `json:"attemptCount"`
	Emails       []string `json:"emails"`
}

type ReasonCount struct {
	Reason string `json:"reason"`
	Count  int    `json:"count"`
}
