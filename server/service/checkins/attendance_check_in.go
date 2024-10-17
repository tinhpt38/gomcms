package checkins

import (
	"encoding/base32"
	"errors"
	"fmt"
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
	// var count int64
	// err = global.GVA_DB.Where(&checkins.AttendanceCheckIn{
	// 	AttendanceId:     attendanceCheckIn.AttendanceId,
	// 	PartpaticipantId: attendanceCheckIn.PartpaticipantId,
	// 	ConditionId:      attendanceCheckIn.ConditionId,
	// }).Count(&count).Error
	// if count == 0 {
	// 	err = global.GVA_DB.Create(attendanceCheckIn).Error
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	// return nil
	err = global.GVA_DB.Create(attendanceCheckIn).Error
	if err != nil {
		return err
	}
	return nil
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

func (attendanceCheckInService *AttendanceCheckInService) GetConditionPassCheckIn(aId uint, pId uint) (list []checkins.AttendanceCheckIn, err error) {
	err = global.GVA_DB.Where("attendance_id = ? AND participant_id = ? AND condition_id != 0", aId, pId).Find(&list).Error
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
	if info.Email != nil {
		partSer := new(ParticipantService)
		participant, _ := partSer.GetParticipantByEmail(*info.Email)
		db = db.Where("partpaticipant_id = ?", participant.ID)
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
		return nil, errors.New("mã QR không hợp lệ")
	}

	//Attendance is _
	attendance, grr := attendanceService.GetAttendance(atId)
	if grr != nil {
		return nil, errors.New("không tìm thấy thông tin điểm danh")
	}

	if attendance.IsLocked {
		return nil, errors.New("phiên điểm danh đã bị khóa")
	}

	if attendance.RestrictIp != nil && *attendance.RestrictIp != "" {
		ipString := attendance.RestrictIp
		ipRanges := strings.Split(*ipString, ",")
		if len(ipRanges) > 0 {
			if !isIPAllowed(ip, ipRanges) {
				return nil, errors.New("địa chỉ IP không được phép")
			}
		}
	}

	var now = time.Now()
	if attendance.StartDate != nil {
		if now.Before(*attendance.StartDate) {
			return nil, errors.New("chưa đến thời gian điểm danh")
		}
	}
	if attendance.EndDate != nil {
		if now.After(*attendance.EndDate) {
			return nil, errors.New("đã hết thời gian điểm danh")
		}
	}

	email := req.Email

	//Participant is _
	participant, perr := participantService.GetParticipantByEmail(email)
	if perr != nil {
		if attendance.AllowGuest {
			participant = checkins.Participant{
				Email:    email,
				FullName: req.FullName,
			}
			perr = participantService.CreateParticipant(&participant)
			if perr != nil {
				return nil, errors.New("tạo mới thông tin thất bại")
			}
		} else {
			return nil, errors.New("email của bạn không phải là thành viên của hệ thống")
		}
	}
	var listAgps []checkins.AttendanceGroupParticipant
	listAgps, gerr := participantService.GetParticipantInAttendance(participant.ID, attendance.ID)
	if gerr != nil {
		if attendance.AllowGuest {
			agpDb := global.GVA_DB.Table(checkins.AttendanceGroupParticipant{}.TableName())
			newagp := &checkins.AttendanceGroupParticipant{
				ParticipantId: &participant.ID,
				GroupId:       nil,
				AttendanceId:  &attendance.ID,
			}

			err := agpDb.Where(checkins.AttendanceGroupParticipant{
				ParticipantId: &participant.ID,
				GroupId:       nil,
				AttendanceId:  &attendance.ID,
			}).FirstOrCreate(&newagp).Error
			if err != nil {
				msg := fmt.Sprintln(`không thể thêm thông tin điểm danh của bạn: %s`, err.Error())
				return nil, errors.New(msg)
			}
			listAgps = append(listAgps, *newagp)
		}
	} else {
		return nil, errors.New("không tìm thấy thông tin điểm danh của bạn")
	}

	// Kiểm tra số lần điểm danh của thành viên đó

	if attendance.LimitCount > 0 {
		var list []checkins.AttendanceCheckIn
		global.GVA_DB.Where("partpaticipant_id = ? AND attendance_id = ?", participant.ID, attendance.ID).Find(&list)
		if len(list) >= attendance.LimitCount {
			return nil, errors.New("bạn đã điểm danh đủ số lần cho phép")
		}
	}

	// Kiểm tra số lần điểm danh của thiết bị đó

	if attendance.LimitClientCount > 0 {
		var limitClientCount int64
		global.GVA_DB.Where("visitor_id = ? and attendance_id = ?", req.VisitorId, attendance.ID).Model(&checkins.AttendanceCheckIn{}).Count(&limitClientCount)
		if limitClientCount >= int64(attendance.LimitClientCount) {
			return nil, errors.New("thiết bị đã điểm danh đủ số lần cho phép")
		}
	}

	// Kiểm tra điều kiện điểm danh
	var listAgpIDs []int
	for _, agp := range listAgps {
		listAgpIDs = append(listAgpIDs, int(agp.ID))
	}
	var conditionCheckedIn []uint
	err = global.GVA_DB.Model(&checkins.AttendanceCheckIn{}).
		Where("partpaticipant_id = ? and attendance_id = ?", participant.ID, attendance.ID).Debug().
		Pluck("condition_id", &conditionCheckedIn).Error

	if err != nil {
		return nil, errors.New("không thể kiểm tra điều kiện điểm danh")
	}

	rConditions, cerr := conditionService.GetConditionOfPartparticipant(attendance.ID, listAgpIDs)
	if cerr != nil || len(rConditions) == 0 {
		return nil, errors.New("không tìm thấy điều kiện điểm danh")
	}

	var agpCheckins []checkins.AttendanceCheckIn
	var coreConditions []checkins.Condition
	globalMsg := "Điểm danh không thành công"
	for _, condition := range rConditions {
		var tempCon checkins.Condition

		for _, agp := range listAgps {
			// tempCon = *rConditions[i].Condition
			tempCon = *condition.Condition
			if condition.AttendanceGroupParticipantId == int(agp.ID) && !arrayContains(conditionCheckedIn, tempCon.ID) {
				result, cerr := checkCondition(agp, tempCon, req, ip)
				if result {
					tempCon.IsPass = true
					globalMsg = "Điểm danh thành công"
					attendanceCheckIn := checkins.AttendanceCheckIn{
						CheckinDate:      time.Now().UTC(),
						AttendanceId:     &attendance.ID,
						PartpaticipantId: &participant.ID,
						AreaId:           tempCon.AreaId,
						GroupId:          agp.GroupId,
						ConditionId:      &tempCon.ID,
						IP:               ip,
						Lattidue:         req.Lat,
						Longtidue:        req.Lng,
						Agent:            userAgent,
						Accuracy:         req.Accuracy,
						VisitorId:        req.VisitorId,
					}
					agpCheckins = append(agpCheckins, attendanceCheckIn)
				} else {
					tempCon.IsPass = false
					tempCon.Message = cerr.Error()
				}
				coreConditions = append(coreConditions, tempCon)
			}

			if arrayContains(conditionCheckedIn, tempCon.ID) {
				tempCon.IsPass = true
				coreConditions = append(coreConditions, tempCon)
			}

		}
	}

	result = make(map[string]interface{})
	result["conditions"] = coreConditions
	result["attendance"] = attendance
	result["message"] = globalMsg
	if len(coreConditions) == 0 {
		result["message"] = "Không có điều kiện điểm danh nào thoả mãn"
	}

	for _, agp := range agpCheckins {
		aciErr := attendanceCheckInService.CreateAttendanceCheckIn(&agp)
		if aciErr != nil {
			return nil, errors.New("điểm danh thất bại")
		}
	}

	return
}

func arrayContains(arr []uint, element uint) bool {
	for _, e := range arr {
		if e == element {
			return true
		}
	}
	return false
}

func (attendanceCheckInService *AttendanceCheckInService) DecodeBase32(encoded string) (string, error) {
	decoded, err := base32.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

func checkMatchArea(ip string, lat *float64, lng *float64, accuracy *float64, area checkins.AttendanceArea) (bool, error) {
	if lat == nil || lng == nil {
		return false, errors.New("không tìm thấy vị trí")
	}
	latArea := area.Area.Latitude
	lngArea := area.Area.Longitude
	defaultRadius := area.Area.Radius // default meter
	if accuracy == nil || *accuracy == 0 {
		*accuracy = float64(0.1)
	}

	radiusArea := area.Radius
	if radiusArea == nil {
		radiusArea = defaultRadius
	}
	radius := float32(float64(*radiusArea) / float64(1000))
	result := isWithinRadius(*latArea, *lngArea, float64(radius), *lat, *lng, *accuracy)
	if area.AllowRestrictIp {
		ipString := area.Area.RestrictIp
		ipRanges := strings.Split(*ipString, ",")
		if len(ipRanges) > 0 {
			if !isIPAllowed(ip, ipRanges) {
				return false, errors.New("IP vị trí của bạn không khớp")
			}
		}
	}
	if !result {
		return false, errors.New("Vị trí của bạn không khớp")
	}
	return result, nil
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
func isWithinRadius(lat, lon, radius, xLat, xLng, accuracy float64) bool {
	distance := haversine(lat, lon, xLat, xLng)
	return distance <= radius
	// return math.Abs(distance-radius) <= accuracy
}

func checkCondition(participant checkins.AttendanceGroupParticipant, condition checkins.Condition, req checkinsReq.CheckinsReq, ip string) (bool, error) {

	lat := req.Lat
	lng := req.Lng
	// Trường hợp 1 -- Điều kiện không phân nhóm, có khu vực, có thời gian
	if condition.GroupId == nil && condition.AreaId != nil && condition.StartAt != nil && condition.EndAt != nil {
		inArea, err := checkMatchArea(ip, lat, lng, req.Accuracy, *condition.Area)
		matchTime := time.Now().UTC().After(*condition.StartAt) && time.Now().UTC().Before(*condition.EndAt)
		return inArea && matchTime, err
	}

	// Trường hợp 2 -- Điều kiện không phân nhóm, không khu vực, có thời gian
	if condition.GroupId == nil && condition.AreaId == nil && condition.StartAt != nil && condition.EndAt != nil {
		result := time.Now().UTC().After(*condition.StartAt) && time.Now().UTC().Before(*condition.EndAt)
		if !result {
			return false, errors.New("Thời gian không khớp")
		}
		return result, nil
	}

	// // Trường hợp 3 -- Điều kiện không phân nhóm, không khu vực, không thời gian bắt đầu, có thời gian kết thúc -- NL
	// if condition.GroupId == nil && condition.AreaId == nil && condition.StartAt == nil && condition.EndAt != nil {
	// 	result := time.Now().UTC().Before(*condition.EndAt)
	// 	if !result {
	// 		return false, errors.New("Thời gian không khớp")
	// 	}
	// 	return result, nil
	// }

	// Trường hợp 4 -- Điều kiện có phân nhóm, khu vực, có thời gian
	if condition.GroupId != nil && condition.AreaId != nil && condition.StartAt != nil && condition.EndAt != nil {
		// Kiểm tra xem có participant.GroupId không
		if participant.GroupId == nil {
			return false, errors.New("Không tìm thấy thông tin nhóm")
		}

		inArea, aerr := checkMatchArea(ip, lat, lng, req.Accuracy, *condition.Area)
		if aerr != nil {
			return false, aerr
		}
		matchTime := time.Now().UTC().After(*condition.StartAt) && time.Now().UTC().Before(*condition.EndAt)
		result := inArea == matchTime && (*condition.GroupId == *participant.GroupId)
		if !result {
			return false, errors.New("Thông tin nhóm và thời gian không thoả điều kiện")
		}

	}

	// Trường hợp 5 -- Điều kiện có phân nhóm, không khu vực, có thời gian
	if condition.GroupId != nil && condition.AreaId == nil && condition.StartAt != nil && condition.EndAt != nil {
		// Kiểm tra xem có participant.GroupId không
		if participant.GroupId == nil {
			return false, errors.New("Không tìm thấy thông tin nhóm")
		}

		matchTime := time.Now().UTC().After(*condition.StartAt) && time.Now().UTC().Before(*condition.EndAt)
		result := (*condition.GroupId == *participant.GroupId) && matchTime
		if !result {
			return false, errors.New("Thông tin nhóm và thời gian không thoả điều kiện")
		}
		return result, nil
	}

	// // Trường hợp 6 -- Điều kiện có phân nhóm, không khu vực, không thời gian bắt đầu, có thời gian kết thúc  -- NL
	// if condition.GroupId != nil && condition.AreaId == nil && condition.StartAt == nil && condition.EndAt != nil {
	// 	// Kiểm tra xem có participant.GroupId không
	// 	if participant.GroupId == nil {
	// 		return false, errors.New("Không tìm thấy thông tin nhóm")
	// 	}

	// 	matchTime := time.Now().UTC().Before(*condition.EndAt)
	// 	result := (*condition.GroupId == *participant.GroupId) && matchTime
	// 	if !result {
	// 		return false, errors.New("Thông tin nhóm và thời gian không thoả điều kiện")
	// 	}
	// 	return result, nil
	// }
	// Có thể gộp giữ 7 và 6
	// Trường hợp 7 - Điều kiện có phân nhóm, không khu vực, không thời gian bắt đầu, không thời gian kết thúc
	if condition.GroupId != nil && condition.AreaId == nil && condition.StartAt == nil && condition.EndAt == nil {
		// Kiểm tra xem có participant.GroupId không
		if participant.GroupId == nil {
			return false, errors.New("Không tìm thấy thông tin nhóm")
		}

		result := *condition.GroupId == *participant.GroupId
		if !result {
			return false, errors.New("Thông tin nhóm không thoả điều kiện")
		}
		return result, nil
	}

	// Trường hợp 8 -- Điều kiện không phân nhóm, không khu vực, có thời gian bắt đầu, không thời gian kết thúc -- NL
	if condition.GroupId == nil && condition.AreaId == nil && condition.StartAt != nil && condition.EndAt == nil {
		result := time.Now().UTC().After(*condition.StartAt)
		if !result {
			return false, errors.New("Thời gian không khớp")
		}
		return result, nil

	}

	// Trường hợp 9 -- Điều kiện có phân nhóm, khu vực, không thời gian bắt đầu, không thời gian kết thúc
	if condition.GroupId != nil && condition.AreaId != nil && condition.StartAt == nil && condition.EndAt == nil {
		// Kiểm tra xem có participant.GroupId không
		if participant.GroupId == nil {
			return false, errors.New("Không tìm thấy thông tin nhóm")
		}

		inArea, aerr := checkMatchArea(ip, lat, lng, req.Accuracy, *condition.Area)
		if aerr != nil {
			return false, aerr
		}

		result := inArea && (*condition.GroupId == *participant.GroupId)
		if !result {
			return false, errors.New("Thông tin nhóm không thoả điều kiện")
		}
	}
	// Trường hợp 10 -- Điều kiện không phân nhóm, có khu vực, không thời gian bắt đầu, không thời gian kết thúc
	if condition.GroupId == nil && condition.AreaId != nil && condition.StartAt == nil && condition.EndAt == nil {
		inArea, err := checkMatchArea(ip, lat, lng, req.Accuracy, *condition.Area)
		return inArea, err
	}

	return true, nil
}

// TODO: Allow 1 trong 2 dải IP trong ipRangs
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
