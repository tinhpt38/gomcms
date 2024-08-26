package checkins

import (
	"encoding/base32"
	"errors"
	"fmt"
	"math"
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
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Preload(clause.Associations).Find(&attendanceCheckIns).Error
	return attendanceCheckIns, total, err
}

func (attendanceCheckInService *AttendanceCheckInService) CheckinAttendance(req checkinsReq.CheckinsReq, ip string, userAgent string) (result map[string]string, err error) {
	// conditionService := new(ConditionService)
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

	//TODO: Xử lý IsTrail sau

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

	pa, gerr := participantService.GetParticipantInAttendance(participant.ID, attendance.ID)
	if gerr != nil {
		return nil, errors.New("Không tìm thấy thông tin điểm danh của bạn")
	}
	fmt.Printf("partticipant in attendance: %v", pa)
	



	/*
		Thành viên có nhóm hay không?
		Nếu có ? Có điều kiện điểm danh hay không?
		Nếu ko có đk thì pass => không quan trọng
		Nếu có điều kiện thì kiểm tra điều kiện
	*/

	attendanceCheckIn := checkins.AttendanceCheckIn{
		CheckinDate:      time.Now().UTC(),
		AttendanceId:     &attendance.ID,
		PartpaticipantId: &participant.ID,
		AreaId:           nil,
		GroupId:          nil,
		ConditionId:      nil,
		IP:               ip,
		Lattidue:         nil,
		Longtidue:        nil,
		Agent:            userAgent,
		Browser:          "",
	}
	aciErr := attendanceCheckInService.CreateAttendanceCheckIn(&attendanceCheckIn)
	if aciErr != nil {
		return nil, errors.New("Điểm danh thất bại")
	}

	// conditions, cerr := conditionService.GetConditionsByAttendanceId(attendance.ID)
	// hasCondition := true
	// if conditions == nil || cerr != nil {
	// 	hasCondition = false
	// }
	// if !hasCondition {
	// 	attendanceCheckIn := checkins.AttendanceCheckIn{
	// 		CheckinDate:      time.Now().UTC(),
	// 		AttendanceId:     &attendance.ID,
	// 		PartpaticipantId: &participant.ID,
	// 		AreaId:           nil,
	// 		GroupId:          &participant.Group.ID,
	// 		ConditionId:      nil,
	// 		IP:               "req.IP",
	// 		Lattidue:         nil,
	// 		Longtidue:        nil,
	// 		Agent:            nil,
	// 		Browser:          "req.Browser",
	// 	}
	// 	aciErr := attendanceCheckInService.CreateAttendanceCheckIn(&attendanceCheckIn)
	// 	if aciErr != nil {
	// 		return errors.New("Điểm danh thất bại")
	// 	}

	// } else {
	// go func() {
	// 	errorMap := make(map[string]string)
	// 	attendanceCheckIn := checkins.AttendanceCheckIn{
	// 		CheckinDate:      time.Now().UTC(),
	// 		AttendanceId:     &attendance.ID,
	// 		PartpaticipantId: &participant.ID,
	// 		AreaId:           nil,
	// 		GroupId:          &participant.Group.ID,
	// 		ConditionId:      nil,
	// 		IP:               "req.IP",
	// 		Lattidue:         nil,
	// 		Longtidue:        nil,
	// 		Agent:            nil,
	// 		Browser:          "req.Browser",
	// 	}
	// 	for _, condition := range conditions {
	// 		if condition.AreaId != nil {
	// 			areaResult := attendanceCheckInService.CheckinCheckArea(*req.Lat, *req.Lng, *condition.Area)
	// 			if !areaResult {
	// 				errorMap["Area"] = "Bạn không ở trong khu vực điểm danh"
	// 			} else {
	// 				attendanceCheckIn.AreaId = condition.AreaId
	// 			}
	// 		}
	// 		if condition.GroupId != nil {
	// 			if *condition.GroupId != participant.Group.ID {
	// 				errorMap["Group"] = "Nhóm của bạn không phải là nhóm của phiên điểm danh"
	// 			}
	// 		}
	// 		if condition.StartAt != nil {
	// 			if time.Now().UTC().Before(*condition.StartAt) {
	// 				errorMap["StartAt"] = "Chưa đến thời gian điểm danh"
	// 			}
	// 		}
	// 		if condition.EndAt != nil {
	// 			if time.Now().UTC().After(*condition.EndAt) {
	// 				errorMap["EndAt"] = "Đã hết thời gian điểm danh"
	// 			}
	// 		}

	// 	}

	// 	attendanceCheckInService.CreateAttendanceCheckIn(&attendanceCheckIn)

	// }()

	// }
	result = make(map[string]string)
	result["message"] = "Điểm danh thành công"
	return
}

func (attendanceCheckInService *AttendanceCheckInService) DecodeBase32(encoded string) (string, error) {
	decoded, err := base32.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

func (attendanceCheckInService *AttendanceCheckInService) CheckinCheckArea(lat float64, lng float64, area checkins.AttendanceArea) bool {
	latArea := area.Area.Latitude
	lngArea := area.Area.Longitude
	defaultRadius := area.Area.Radius // default meter
	radiusArea := area.Radius
	if radiusArea == nil {
		radiusArea = defaultRadius
	}
	radius := float32(float64(*radiusArea) / float64(1000))
	result := isWithinRadius(*latArea, *lngArea, float64(radius), lat, lng)
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
