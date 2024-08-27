package checkins

import (
	"encoding/base32"
	"errors"
	"math"
	"net"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
	checkinsReq "github.com/flipped-aurora/gin-vue-admin/server/model/checkins/request"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AttendanceCheckInService struct{}

func (attendanceCheckInService *AttendanceCheckInService) CreateAttendanceCheckIn(attendanceCheckIn *checkins.AttendanceCheckIn) (err error) {
	err = global.GVA_DB.Create(attendanceCheckIn).Error
	return err
}

func (attendanceCheckInService *AttendanceCheckInService) DeleteAttendanceCheckIn(ID string, userID uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&checkins.AttendanceCheckIn{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
			return err
		}
		if err = tx.Delete(&checkins.AttendanceCheckIn{}, "id = ?", ID).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (attendanceCheckInService *AttendanceCheckInService) DeleteAttendanceCheckInByIds(IDs []string, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&checkins.AttendanceCheckIn{}).Where("id in ?", IDs).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", IDs).Delete(&checkins.AttendanceCheckIn{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (attendanceCheckInService *AttendanceCheckInService) UpdateAttendanceCheckIn(attendanceCheckIn checkins.AttendanceCheckIn) (err error) {
	err = global.GVA_DB.Model(&checkins.AttendanceCheckIn{}).Where("id = ?", attendanceCheckIn.ID).Updates(&attendanceCheckIn).Error
	return err
}

func (attendanceCheckInService *AttendanceCheckInService) GetAttendanceCheckIn(ID string) (attendanceCheckIn checkins.AttendanceCheckIn, err error) {
	err = global.GVA_DB.Where("id = ?", ID).Preload(clause.Associations).First(&attendanceCheckIn).Error
	return
}

func (attendanceCheckInService *AttendanceCheckInService) GetAttendanceCheckInInfoList(info checkinsReq.AttendanceCheckInSearch) (list []checkins.AttendanceCheckIn, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.GVA_DB.Model(&checkins.AttendanceCheckIn{})
	var attendanceCheckIns []checkins.AttendanceCheckIn

	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.AttendanceId != nil {
		db = db.Where("attendance_id = ?", *info.AttendanceId)
	}
	
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Order("created_at desc").Preload(clause.Associations).Find(&attendanceCheckIns).Error
	return attendanceCheckIns, total, err
}

func (attendanceCheckInService *AttendanceCheckInService) CheckinAttendance(req checkinsReq.CheckinsReq, ip string, userAgent string) (result map[string]interface{}, err error) {
	conditionService := new(ConditionService)
	attendanceService := new(AttendanceService)
	participantService := new(ParticipantService)

	atId, derr := attendanceCheckInService.DecodeBase32(req.Code)
	if derr != nil {
		return nil, errors.New("Mã QR không hợp lệ")
	}

	//Attendance is _
	attendance, grr := attendanceService.GetAttendance(atId)
	if grr != nil {
		return nil, errors.New("Không tìm thấy thông tin điểm danh")
	}

	if attendance.IsLocked {
		return nil, errors.New("Phiên điểm danh đã bị khóa")
	}

	if attendance.RestrictIp != nil && *attendance.RestrictIp != "" {
		ipString := attendance.RestrictIp
		ipRanges := strings.Split(*ipString, ",")
		if len(ipRanges) > 0 {
			if !isIPAllowed(ip, ipRanges) {
				return nil, errors.New("Địa chỉ IP không được phép")
			}
		}
	}

	var now = time.Now()
	if attendance.StartDate != nil {
		if now.Before(*attendance.StartDate) {
			return nil, errors.New("Chưa đến thời gian điểm danh")
		}
	}
	if attendance.EndDate != nil {
		if now.After(*attendance.EndDate) {
			return nil, errors.New("Đã hết thời gian điểm danh")
		}
	}

	email := req.Email

	//Participant is _
	participant, perr := participantService.GetParticipantByEmail(email)
	if perr != nil {
		return nil, errors.New("Email của bạn không phải là thành viên của hệ thống")
	}

	agp, gerr := participantService.GetParticipantInAttendance(participant.ID, attendance.ID)
	if gerr != nil {
		return nil, errors.New("Không tìm thấy thông tin điểm danh của bạn")
	}

	// Đã đạt số lần giới hạn hay chưa
	var list []checkins.AttendanceCheckIn
	global.GVA_DB.Where("partpaticipant_id = ? AND attendance_id = ?", participant.ID, attendance.ID).Find(&list)
	if attendance.LimitCount != nil && len(list) >= *attendance.LimitCount {
		return nil, errors.New("Bạn đã điểm danh đủ số lần cho phép")
	}

	conditions, cerr := conditionService.GetConditionsByAttendanceId(attendance.ID)
	if cerr != nil {
		return nil, errors.New("Không tìm thấy điều kiện điểm danh")
	}
	var satisfiedConditions []checkins.Condition
	var conditionsStatus []checkins.Condition
	if len(conditions) > 0 {
		satisfiedConditions, conditionsStatus = func() ([]checkins.Condition, []checkins.Condition) {

			var usedConditon []checkins.Condition
			// var unUsedConditon []checkins.Condition
			var conditionsStatus []checkins.Condition

			for _, condition := range conditions {
				if checkCondition(agp, condition, req.Lat, req.Lng) {
					condition.IsPass = true
					usedConditon = append(usedConditon, condition)
				} else {
					condition.IsPass = false
					// unUsedConditon = append(unUsedConditon, condition)
				}
				conditionsStatus = append(conditionsStatus, condition)
			}
			return usedConditon, conditionsStatus
		}()
	}
	result = make(map[string]interface{})
	// result["pass"] = satisfiedConditions
	// result["fail"] = unsatisfiedConditions
	result["conditions"] = conditionsStatus
	// xử lý những thằng nào đã pass
	firstCondition := &checkins.Condition{}
	if len(conditions) > 0 {
		if len(satisfiedConditions) == 0 {
			result["message"] = "Không thoả điều kiện điểm danh"
		} else {
			firstCondition = &satisfiedConditions[0]
			result["message"] = "Điểm danh thành công"
		}
	}

	attendanceCheckIn := checkins.AttendanceCheckIn{
		CheckinDate:      time.Now().UTC(),
		AttendanceId:     &attendance.ID,
		PartpaticipantId: &participant.ID,
		AreaId:           firstCondition.AreaId,
		GroupId:          agp.GroupId,
		ConditionId:      &firstCondition.ID,
		IP:               ip,
		Lattidue:         req.Lat,
		Longtidue:        req.Lng,
		Agent:            userAgent,
		Browser:          "",
	}
	aciErr := attendanceCheckInService.CreateAttendanceCheckIn(&attendanceCheckIn)
	if aciErr != nil {
		return nil, errors.New("Điểm danh thất bại")
	}
	return
}

func (attendanceCheckInService *AttendanceCheckInService) DecodeBase32(encoded string) (string, error) {
	decoded, err := base32.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

func checkMatchArea(lat *float64, lng *float64, area checkins.AttendanceArea) bool {
	if lat == nil || lng == nil {
		return false
	}
	latArea := area.Area.Latitude
	lngArea := area.Area.Longitude
	defaultRadius := area.Area.Radius // default meter
	radiusArea := area.Radius
	if radiusArea == nil {
		radiusArea = defaultRadius
	}
	radius := float32(float64(*radiusArea) / float64(1000))
	result := isWithinRadius(*latArea, *lngArea, float64(radius), *lat, *lng)
	return result
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert degrees to radians
	const pi = math.Pi
	lat1Rad := lat1 * pi / 180
	lon1Rad := lon1 * pi / 180
	lat2Rad := lat2 * pi / 180
	lon2Rad := lon2 * pi / 180

	// Haversine formula
	dlat := lat2Rad - lat1Rad
	dlon := lon2Rad - lon1Rad
	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Earth radius in kilometers
	const earthRadius = 6371
	return earthRadius * c
}

// isWithinRadius checks if the point (xLat, xLng) is within the given radius (in kilometers) of (lat, lon).
func isWithinRadius(lat, lon, radius, xLat, xLng float64) bool {
	distance := haversine(lat, lon, xLat, xLng)
	return distance <= radius
}

func checkCondition(participant checkins.AttendanceGroupParticipant, condition checkins.Condition, lat *float64, lng *float64) bool {
	// Trường hợp 1

	if condition.GroupId == nil && condition.AreaId != nil && condition.StartAt != nil && condition.EndAt != nil {
		inArea := checkMatchArea(lat, lng, *condition.Area)
		matchTime := time.Now().UTC().After(*condition.StartAt) && time.Now().UTC().Before(*condition.EndAt)
		return inArea && matchTime
	}

	// Trường hợp 2
	if condition.GroupId == nil && condition.AreaId == nil && condition.StartAt != nil && condition.EndAt != nil {
		return time.Now().UTC().After(*condition.StartAt) && time.Now().UTC().Before(*condition.EndAt)
	}

	// Trường hợp 3
	if condition.GroupId == nil && condition.AreaId == nil && condition.StartAt == nil && condition.EndAt != nil {
		return time.Now().UTC().Before(*condition.EndAt)
	}

	// Trường hợp 4
	if condition.GroupId != nil && condition.AreaId != nil && condition.StartAt != nil && condition.EndAt != nil {
		inArea := checkMatchArea(lat, lng, *condition.Area)
		matchTime := time.Now().UTC().After(*condition.StartAt) && time.Now().UTC().Before(*condition.EndAt)
		return inArea == matchTime && (*condition.GroupId == *participant.GroupId)
	}

	// Trường hợp 5
	if condition.GroupId != nil && condition.AreaId == nil && condition.StartAt != nil && condition.EndAt != nil {
		matchTime := time.Now().UTC().After(*condition.StartAt) && time.Now().UTC().Before(*condition.EndAt)
		return (*condition.GroupId == *participant.GroupId) && matchTime
	}

	// Trường hợp 6
	if condition.GroupId != nil && condition.AreaId == nil && condition.StartAt == nil && condition.EndAt != nil {
		matchTime := time.Now().UTC().Before(*condition.EndAt)
		return (*condition.GroupId == *participant.GroupId) && matchTime
	}

	// Trường hợp 7
	if condition.GroupId != nil && condition.AreaId == nil && condition.StartAt == nil && condition.EndAt == nil {
		return *condition.GroupId == *participant.GroupId
	}

	// Trường hợp 8
	if condition.GroupId == nil && condition.AreaId == nil && condition.StartAt != nil && condition.EndAt == nil {
		return time.Now().UTC().After(*condition.StartAt)
		// now := time.Now().UTC()
		// fmt.Printf("now: %v", now)
		// fmt.Printf("startAt: %v", *condition.StartAt)
		// return now.After(*condition.StartAt)
	}

	// Trường hợp 9
	if condition.GroupId != nil && condition.AreaId != nil && condition.StartAt == nil && condition.EndAt == nil {
		inArea := checkMatchArea(lat, lng, *condition.Area)
		return inArea && (*condition.GroupId == *participant.GroupId)
	}

	return false
}

func isIPAllowed(clientIP string, ipRanges []string) bool {
	for _, ipRange := range ipRanges {
		_, ipNet, err := net.ParseCIDR(ipRange)
		if err != nil {
			ip := net.ParseIP(ipRange)
			if ip != nil && ip.String() == clientIP {
				return true
			}
		} else if ipNet.Contains(net.ParseIP(clientIP)) {
			return true
		}
	}
	return false
}
