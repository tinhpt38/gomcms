package checkins

import (
	"encoding/base32"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
	checkinsReq "github.com/flipped-aurora/gin-vue-admin/server/model/checkins/request"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AttendanceCheckInService struct{}

type CheckinContext struct {
	ParticipantID uint
	AttendanceID  uint
	Timestamp     time.Time
	IP            string
	Lat           *float64
	Lng           *float64
	ConditionLogs []string
}

func logConditionCheck(ctx *CheckinContext, conditionID uint, result bool, message string) {
	log := fmt.Sprintf("Điều kiện #%d: %s - %v", conditionID, message, result)
	ctx.ConditionLogs = append(ctx.ConditionLogs, log)

	// Có thể lưu log này vào database trong tương lai
}

type conditionResult struct {
	ConditionID uint
	Pass        bool
	Error       error
}

func (attendanceCheckInService *AttendanceCheckInService) CreateAttendanceCheckIn(attendanceCheckIn *checkins.AttendanceCheckIn) (err error) {
	// Kiểm tra xem người tham gia đã điểm danh với điều kiện này trước đó chưa
	var existingCheckIn checkins.AttendanceCheckIn

	err = global.GVA_DB.Where(&checkins.AttendanceCheckIn{
		AttendanceId:     attendanceCheckIn.AttendanceId,
		PartpaticipantId: attendanceCheckIn.PartpaticipantId,
		ConditionId:      attendanceCheckIn.ConditionId,
	}).First(&existingCheckIn).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Nếu chưa điểm danh với điều kiện này trước đó, tạo record mới với counter = 1
			attendanceCheckIn.Counter = 1
			err = global.GVA_DB.Create(attendanceCheckIn).Error
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}

	// Nếu đã điểm danh trước đó, tăng counter và cập nhật
	// Đảm bảo counter không bao giờ là 0
	if existingCheckIn.Counter <= 0 {
		existingCheckIn.Counter = 1
	}
	existingCheckIn.Counter += 1

	// Cập nhật các thông tin mới nếu cần
	if attendanceCheckIn.IP != "" {
		existingCheckIn.IP = attendanceCheckIn.IP
	}
	if attendanceCheckIn.Lattidue != nil {
		existingCheckIn.Lattidue = attendanceCheckIn.Lattidue
	}
	if attendanceCheckIn.Longtidue != nil {
		existingCheckIn.Longtidue = attendanceCheckIn.Longtidue
	}
	if attendanceCheckIn.Accuracy != nil {
		existingCheckIn.Accuracy = attendanceCheckIn.Accuracy
	}
	if attendanceCheckIn.Agent != "" {
		existingCheckIn.Agent = attendanceCheckIn.Agent
	}
	if attendanceCheckIn.VisitorId != "" {
		existingCheckIn.VisitorId = attendanceCheckIn.VisitorId
	}

	// Cập nhật số may mắn nếu có
	if attendanceCheckIn.LuckyNumber > 0 {
		existingCheckIn.LuckyNumber = attendanceCheckIn.LuckyNumber
	}

	// Cập nhật thời gian điểm danh mới nhất
	existingCheckIn.CheckinDate = attendanceCheckIn.CheckinDate

	// Lưu cập nhật vào database
	err = global.GVA_DB.Save(&existingCheckIn).Error
	if err != nil {
		return err
	}

	// Cập nhật lại thông tin counter cho đối tượng ban đầu để hiển thị
	attendanceCheckIn.Counter = existingCheckIn.Counter
	attendanceCheckIn.ID = existingCheckIn.ID

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
	// First, get the existing record to ensure we don't lose data
	var existingCheckIn checkins.AttendanceCheckIn
	err = global.GVA_DB.Where("id = ?", attendanceCheckIn.ID).First(&existingCheckIn).Error
	if err != nil {
		return err
	}

	// Preserve lucky number if it exists
	if existingCheckIn.LuckyNumber > 0 && attendanceCheckIn.LuckyNumber == 0 {
		attendanceCheckIn.LuckyNumber = existingCheckIn.LuckyNumber
	}

	// Update the record
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

	if info.GroupId != nil {
		db = db.Where("group_id = ?", *info.GroupId)
	}

	if info.Agent != nil {
		db = db.Where("agent LIKE ?", "%"+*info.Agent+"%")
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
	// Sử dụng bộ tạo số ngẫu nhiên cục bộ
	localRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Khởi tạo biến cho số may mắn
	var luckyNumber int

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

	var listAgpsTemp []checkins.AttendanceGroupParticipant
	listAgpsTemp, _ = participantService.GetParticipantInAttendance(participant.ID, attendance.ID)
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
			msg := fmt.Sprintf("không thể thêm thông tin điểm danh của bạn: %s", err.Error())
			return nil, errors.New(msg)
		}

		listAgpsTemp = append(listAgpsTemp, *newagp)
	}

	var listAgps []checkins.AttendanceGroupParticipant
	listAgpIDSet := make(map[int]bool) // Sử dụng map để kiểm tra trùng lặp
	for _, agp := range listAgpsTemp {
		agpID := int(agp.ID)
		if !listAgpIDSet[agpID] { // Nếu chưa tồn tại trong map
			listAgps = append(listAgps, agp)
			listAgpIDSet[agpID] = true // Đánh dấu là đã tồn tại
		}
	}

	if len(listAgps) == 0 {
		return nil, errors.New("không tìm thấy thông tin điểm danh của bạn")
	}

	// Kiểm tra số lần điểm danh của thành viên đó
	// Tạo checkinLog mới
	checkinLog := checkins.CheckinLog{
		Email:        email,
		Code:         req.Code,
		AttendanceId: attendance.ID,
		VisitorId:    req.VisitorId,
		Lat:          req.Lat,
		Lng:          req.Lng,
		Accuracy:     req.Accuracy,
		Ip:           ip,
		Agent:        userAgent,
		FullName:     req.FullName,
	}

	// Kiểm tra xem đã có log với cùng thông tin hay không
	var existingLog checkins.CheckinLog
	err = global.GVA_DB.Where(&checkins.CheckinLog{
		Email:        email,
		Code:         req.Code,
		AttendanceId: attendance.ID,
		VisitorId:    req.VisitorId,
		Ip:           ip,
	}).Order("created_at DESC").First(&existingLog).Error

	// Nếu tìm thấy log trước đó, chỉ tạo log mới nếu có sự thay đổi về dữ liệu
	if err == nil {
		// Kiểm tra sự khác biệt giữa log hiện tại và log mới
		hasChanges := false

		// So sánh các thông tin về vị trí
		if (existingLog.Lat == nil && req.Lat != nil) ||
			(existingLog.Lat != nil && req.Lat == nil) ||
			(existingLog.Lat != nil && req.Lat != nil && *existingLog.Lat != *req.Lat) {
			hasChanges = true
		}

		if (existingLog.Lng == nil && req.Lng != nil) ||
			(existingLog.Lng != nil && req.Lng == nil) ||
			(existingLog.Lng != nil && req.Lng != nil && *existingLog.Lng != *req.Lng) {
			hasChanges = true
		}

		if (existingLog.Accuracy == nil && req.Accuracy != nil) ||
			(existingLog.Accuracy != nil && req.Accuracy == nil) ||
			(existingLog.Accuracy != nil && req.Accuracy != nil && *existingLog.Accuracy != *req.Accuracy) {
			hasChanges = true
		}

		// So sánh các thông tin khác
		// if existingLog.Agent != userAgent {
		// 	hasChanges = true
		// }

		// Nếu không có sự thay đổi, sử dụng lại log cũ
		if !hasChanges {
			checkinLog = existingLog
			global.GVA_LOG.Info("Sử dụng lại log hiện có thay vì tạo mới", zap.String("email", email), zap.String("visitorId", req.VisitorId))
		} else {
			// Có sự thay đổi, tạo log mới
			if err := global.GVA_DB.Create(&checkinLog).Error; err != nil {
				return nil, err
			}
		}
	} else {
		// Không tìm thấy log trước đó, tạo mới
		if err := global.GVA_DB.Create(&checkinLog).Error; err != nil {
			return nil, err
		}
	}

	if attendance.RestrictIp != nil && *attendance.RestrictIp != "" {
		ipString := attendance.RestrictIp
		ipRanges := strings.Split(*ipString, ",")
		if len(ipRanges) > 0 {
			if !isIPAllowed(ip, ipRanges) {
				// return nil, errors.New("địa chỉ IP không được phép")
				msg := fmt.Sprintf("địa chỉ ip của bạn không được phép điểm danh %s ", ip)
				checkinLog.MessageList += msg + "$$"
				global.GVA_DB.Where(checkins.CheckinLog{}).Where("id = ?", checkinLog.ID).Save(&checkinLog)
				return nil, errors.New(msg + ". Hệ thống đã ghi nhận lịch sử điểm danh của bạn")
			}
		}
	}

	if (req.Lat == nil || req.Lng == nil) || (*req.Lat == 0 || *req.Lng == 0) {
		msg := "bạn chưa cho phép truy cập vị trí trên thiết bị "
		checkinLog.MessageList += msg + "$$"
		global.GVA_DB.Where(checkins.CheckinLog{}).Where("id = ?", checkinLog.ID).Save(&checkinLog)
		return nil, errors.New(msg + ". Hệ thống đã ghi nhận lịch sử điểm danh của bạn")
	}

	if attendance.LimitCount > 0 {
		var checkinCount int64
		global.GVA_DB.Model(&checkins.AttendanceCheckIn{}).Where("partpaticipant_id = ? AND attendance_id = ?", participant.ID, attendance.ID).Count(&checkinCount)
		// var list []checkins.AttendanceCheckIn
		if int(checkinCount) >= attendance.LimitCount {
			msg := fmt.Sprintf("bạn đã điểm danh đủ số lần cho phép %d ", checkinCount)
			checkinLog.MessageList += msg + "$$"
			global.GVA_DB.Where(checkins.CheckinLog{}).Where("id = ?", checkinLog.ID).Save(&checkinLog)
			return nil, errors.New(msg + ". Hệ thống đã ghi nhận lịch sử điểm danh của bạn")
		}
	}

	// Kiểm tra số lần điểm danh của thiết bị đó
	if attendance.LimitClientCount > 0 {
		var limitClientCount int64
		global.GVA_DB.Where("visitor_id = ? and attendance_id = ?", req.VisitorId, attendance.ID).Model(&checkins.AttendanceCheckIn{}).Count(&limitClientCount)
		if limitClientCount >= int64(attendance.LimitClientCount) {
			msg := fmt.Sprintf("thiết bị của bạn đã điểm danh đủ số lần cho phép %d ", limitClientCount)
			checkinLog.MessageList += msg + "$$"
			global.GVA_DB.Where(checkins.CheckinLog{}).Where("id = ?", checkinLog.ID).Save(&checkinLog)
			return nil, errors.New(msg + ". Hệ thống đã ghi nhận lịch sử điểm danh của bạn")
			// return nil, errors.New("thiết bị đã điểm danh đủ số lần cho phép")
		}
	}

	// Kiểm tra điều kiện điểm danh
	var listAgpIDs []int
	agpIDSet := make(map[int]bool) // Sử dụng map để kiểm tra trùng lặp
	for _, agp := range listAgps {
		agpID := int(agp.ID)
		if !agpIDSet[agpID] { // Nếu chưa tồn tại trong map
			listAgpIDs = append(listAgpIDs, agpID)
			agpIDSet[agpID] = true // Đánh dấu là đã tồn tại
		}
	}

	var conditionCheckedIn []uint
	_ = global.GVA_DB.Model(&checkins.AttendanceCheckIn{}).
		Where("partpaticipant_id = ? and attendance_id = ?", participant.ID, attendance.ID).
		Pluck("condition_id", &conditionCheckedIn).Error

	rConditions, cerr := conditionService.GetConditionOfPartparticipant(attendance.ID, listAgpIDs)

	var validConditions []checkins.AGPCondition
	for _, condition := range rConditions {
		if condition.Condition != nil {
			validConditions = append(validConditions, condition)
		}
	}
	rConditions = validConditions

	if cerr != nil {
		// return nil, errors.New("không tìm thấy điều kiện điểm danh")
		msg := "không tìm thấy điều kiện điểm danh"
		checkinLog.MessageList += msg + "$$"
		global.GVA_DB.Where(checkins.CheckinLog{}).Where("id = ?", checkinLog.ID).Save(&checkinLog)
		return nil, errors.New(msg + ". Hệ thống đã ghi nhận lịch sử điểm danh của bạn")
	}

	var agpCheckins []checkins.AttendanceCheckIn
	globalMsg := "Điểm danh không thành công"
	var coreConditions []checkins.Condition

	if len(rConditions) == 0 {
		// Trường hợp không có điều kiện cụ thể, điểm danh tự do
		globalMsg = "Điểm danh thành công"
		for _, agp := range listAgps {
			attendanceCheckIn := checkins.AttendanceCheckIn{
				CheckinDate:      time.Now().UTC(),
				AttendanceId:     &attendance.ID,
				PartpaticipantId: &participant.ID,
				AreaId:           nil,
				GroupId:          agp.GroupId,
				ConditionId:      nil,
				IP:               ip,
				Lattidue:         req.Lat,
				Longtidue:        req.Lng,
				Agent:            userAgent,
				Accuracy:         req.Accuracy,
				VisitorId:        req.VisitorId,
			}
			agpCheckins = append(agpCheckins, attendanceCheckIn)
		}
	} else {
		// Chuẩn bị danh sách các điều kiện cần kiểm tra cho mỗi agp
		type groupCondition struct {
			agp        checkins.AttendanceGroupParticipant
			conditions []checkins.Condition
		}

		// Nhóm các điều kiện theo từng agp để xử lý song song hiệu quả hơn
		agpConditions := make(map[uint]groupCondition)

		for _, condition := range rConditions {
			for _, agp := range listAgps {
				if condition.AttendanceGroupParticipantId == int(agp.ID) {
					// Nếu điều kiện đã được điểm danh trước đó, đánh dấu là đã qua
					// nhưng KHÔNG tạo bản ghi mới - chỉ cập nhật counter ở cuối hàm
					if arrayContains(conditionCheckedIn, condition.Condition.ID) {
						tempCon := *condition.Condition
						tempCon.IsPass = true
						coreConditions = append(coreConditions, tempCon)
						continue
					}

					// Thêm vào danh sách kiểm tra
					gc, exists := agpConditions[agp.ID]
					if !exists {
						gc = groupCondition{
							agp:        agp,
							conditions: make([]checkins.Condition, 0),
						}
					}
					gc.conditions = append(gc.conditions, *condition.Condition)
					agpConditions[agp.ID] = gc
				}
			}
		}

		// Kiểm tra các điều kiện song song cho từng nhóm
		var anyConditionPassed bool
		for _, gc := range agpConditions {
			// Sử dụng hàm kiểm tra song song
			resultList := checkConditionsParallel(gc.agp, gc.conditions, req, ip) // Xử lý kết quả
			for _, result := range resultList {
				var tempCon checkins.Condition
				for _, c := range gc.conditions {
					if c.ID == result.ConditionID {
						tempCon = c
						break
					}
				}
				if result.Pass {
					tempCon.IsPass = true
					anyConditionPassed = true
					// Tạo bản ghi điểm danh
					attendanceCheckIn := checkins.AttendanceCheckIn{
						CheckinDate:      time.Now().UTC(),
						AttendanceId:     &attendance.ID,
						PartpaticipantId: &participant.ID,
						AreaId:           tempCon.AreaId,
						GroupId:          gc.agp.GroupId,
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
					tempCon.Message = formatErrorMessage(tempCon, result.Error)
				}

				coreConditions = append(coreConditions, tempCon)
			}
		}

		// Cập nhật thông báo nếu có ít nhất một điều kiện đạt
		if anyConditionPassed {
			globalMsg = "Điểm danh thành công"
		}
	}

	result = make(map[string]interface{})
	result["conditions"] = coreConditions
	result["attendance"] = attendance
	result["message"] = globalMsg

	// Kiểm tra và tạo số may mắn nếu điều kiện thỏa mãn
	if attendance.UseLuckyNumber && attendance.LuckyShowAfterMinCount > 0 {
		// Trong trường hợp này, các điểm danh đã được lưu trước đó, nên có thể truy vấn trực tiếp
		// từ cơ sở dữ liệu để lấy tổng số điểm danh
		type CounterSum struct {
			TotalCounter int64
		}
		var counterSum CounterSum
		global.GVA_DB.Model(&checkins.AttendanceCheckIn{}).
			Select("COALESCE(SUM(counter), 0) as total_counter").
			Where("partpaticipant_id = ? AND attendance_id = ?", participant.ID, attendance.ID).
			Scan(&counterSum)

		// Số điểm danh hiện tại (đã bao gồm các điểm danh mới)
		currentCheckinCount := int(counterSum.TotalCounter)

		// Kiểm tra xem có điều kiện nào đạt và cho phép hiển thị số may mắn không
		var hasShowLuckyNumberCondition bool

		// Kiểm tra xem có điểm danh thành công ít nhất một điều kiện không
		var hasAnyPassedCondition bool

		// Nếu có điều kiện đã điểm danh thành công hoặc điểm danh mới thành công
		if len(coreConditions) == 0 {
			// Nếu không có điều kiện, chỉ kiểm tra có điểm danh nào trong agpCheckins không
			hasAnyPassedCondition = len(agpCheckins) > 0
			hasShowLuckyNumberCondition = true
		} else {
			// Kiểm tra điều kiện đã pass
			for _, condition := range coreConditions {
				if condition.IsPass {
					hasAnyPassedCondition = true
					if condition.ShowLuckyNumber {
						hasShowLuckyNumberCondition = true
					}
				}
			}
		}

		// Nếu đủ điều kiện để hiển thị số may mắn (đủ số lần điểm danh và có ít nhất 1 điều kiện thành công)
		if hasAnyPassedCondition && hasShowLuckyNumberCondition && currentCheckinCount >= attendance.LuckyShowAfterMinCount {
			// Kiểm tra xem người tham gia đã có số may mắn chưa
			var existingLuckyCheckIn checkins.AttendanceCheckIn
			existingLuckyResult := global.GVA_DB.Where("partpaticipant_id = ? AND attendance_id = ? AND lucky_number > 0",
				participant.ID, attendance.ID).First(&existingLuckyCheckIn)

			if existingLuckyResult.Error == nil && existingLuckyCheckIn.LuckyNumber > 0 {
				// Nếu đã có số may mắn, sử dụng lại số đó
				luckyNumber = existingLuckyCheckIn.LuckyNumber
			} else {
				// Sử dụng hàm riêng để sinh số may mắn không trùng lặp
				luckyNumber, err = attendanceCheckInService.generateUniqueLuckyNumber(attendance.ID, participant.ID, localRand)
				if err != nil {
					global.GVA_LOG.Error("Lỗi khi sinh số may mắn", zap.Error(err))
				}
			}

			// Lưu số may mắn vào các bản ghi điểm danh mới
			for i := range agpCheckins {
				agpCheckins[i].LuckyNumber = luckyNumber
			}

			// Thêm số may mắn vào kết quả trả về
			result["luckyNumber"] = luckyNumber
		}
	}

	// Lưu tất cả các điểm danh vào database trước
	for _, agp := range agpCheckins {
		aciErr := attendanceCheckInService.CreateAttendanceCheckIn(&agp)
		if aciErr != nil {
			msg := fmt.Sprintf("thiết bị của bạn đã điểm danh đủ số lần cho phép %s ", aciErr.Error())
			checkinLog.MessageList += msg + "$$"
			global.GVA_DB.Where(checkins.CheckinLog{}).Where("id = ?", checkinLog.ID).Save(&checkinLog)
			return nil, errors.New(msg + ". Hệ thống đã ghi nhận lịch sử điểm danh của bạn")
			// return nil, errors.New("điểm danh thất bại" + aciErr.Error())
		}
	}

	// Sau khi lưu tất cả các điểm danh mới, lấy lại toàn bộ danh sách điểm danh từ database
	var allCheckins []checkins.AttendanceCheckIn
	err = global.GVA_DB.Where("partpaticipant_id = ? AND attendance_id = ?", participant.ID, attendance.ID).Find(&allCheckins).Error
	if err == nil {
		// Gán lại danh sách điểm danh đã được cập nhật counter vào kết quả
		result["checkins"] = allCheckins

		// Kiểm tra xem đã có số may mắn đã cấp trước đó hay chưa
		if luckyNumber == 0 {
			for _, checkin := range allCheckins {
				if checkin.LuckyNumber > 0 {
					luckyNumber = checkin.LuckyNumber
					result["luckyNumber"] = luckyNumber
					break
				}
			}
		}

		// Bổ sung thông tin counter vào conditions
		if len(allCheckins) > 0 && len(coreConditions) > 0 {
			// Tạo map để ánh xạ condition_id với counter
			counterMap := make(map[uint]int)
			for _, checkin := range allCheckins {
				if checkin.ConditionId != nil {
					counterMap[*checkin.ConditionId] = checkin.Counter
				}
			}

			// Bổ sung thông tin counter vào kết quả conditions
			var conditionsWithCounter []map[string]interface{}
			for _, condition := range coreConditions {
				conditionMap := make(map[string]interface{})
				for k, v := range map[string]interface{}{
					"ID":              condition.ID,
					"AttendanceId":    condition.AttendanceId,
					"GroupId":         condition.GroupId,
					"AreaId":          condition.AreaId,
					"StartAt":         condition.StartAt,
					"EndAt":           condition.EndAt,
					"IsPass":          condition.IsPass,
					"Message":         condition.Message,
					"ShowLuckyNumber": condition.ShowLuckyNumber,
				} {
					conditionMap[k] = v
				}

				// Thêm counter vào condition
				if counter, exists := counterMap[condition.ID]; exists {
					conditionMap["counter"] = counter
				} else {
					conditionMap["counter"] = 0
				}

				conditionsWithCounter = append(conditionsWithCounter, conditionMap)
			}

			// Thay thế conditions trong kết quả
			result["conditions"] = conditionsWithCounter
		}
	}

	// Cập nhật thông báo
	if luckyNumber > 0 {
		result["message"] = fmt.Sprintf("Điểm danh thành công. Số may mắn của bạn là: %d", luckyNumber)
	} else {
		result["message"] = "Điểm danh thành công"
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

// checkMatchArea là phiên bản cải tiến của hàm kiểm tra vị trí
func checkMatchArea(ip string, lat *float64, lng *float64, accuracy *float64, area checkins.AttendanceArea) (bool, error) {
	// Kiểm tra dữ liệu đầu vào
	if lat == nil || lng == nil {
		return false, errors.New("không tìm thấy vị trí")
	}

	// Lấy thông tin khu vực
	if area.Area == nil || area.Area.Latitude == nil || area.Area.Longitude == nil {
		return false, errors.New("thông tin khu vực không đầy đủ")
	}

	// Ưu tiên kiểm tra IP trước (nếu có) vì xác thực IP nhanh hơn tính toán khoảng cách
	if area.AllowRestrictIp && area.Area.RestrictIp != nil && *area.Area.RestrictIp != "" {
		ipRanges := strings.Split(*area.Area.RestrictIp, ",")
		if len(ipRanges) > 0 && !isIPAllowed(ip, ipRanges) {
			return false, errors.New("địa chỉ IP của bạn không được phép truy cập")
		}
	}

	// Kiểm tra khoảng cách
	latArea := area.Area.Latitude
	lngArea := area.Area.Longitude

	// Xác định bán kính: ưu tiên bán kính từ AttendanceArea, nếu không có thì lấy từ Area
	var radiusMeters float64
	if area.Radius != nil && *area.Radius > 0 {
		radiusMeters = *area.Radius
	} else if area.Area.Radius != nil && *area.Area.Radius > 0 {
		radiusMeters = *area.Area.Radius
	} else {
		radiusMeters = 100 // Mặc định bán kính 100m nếu không có thông tin
	}

	// Chuyển đổi radius từ mét sang km cho hàm tính khoảng cách
	radiusKm := radiusMeters / 1000.0

	// Tính khoảng cách và kiểm tra
	distance := haversine(*latArea, *lngArea, *lat, *lng)

	// Xem xét độ chính xác khi so sánh
	if distance > radiusKm {
		// Log chi tiết hơn về khoảng cách
		return false, fmt.Errorf("vị trí của bạn nằm ngoài phạm vi cho phép. Khoảng cách: %0.2f km, tối đa cho phép: %0.2f km",
			distance, radiusKm)
	}

	return true, nil
}

// haversine tính khoảng cách giữa hai điểm trên mặt đất (đơn vị km)
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	// Chuyển độ sang radian
	const pi = math.Pi
	lat1Rad := lat1 * pi / 180
	lon1Rad := lon1 * pi / 180
	lat2Rad := lat2 * pi / 180
	lon2Rad := lon2 * pi / 180

	// Công thức haversine
	dlat := lat2Rad - lat1Rad
	dlon := lon2Rad - lon1Rad
	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Bán kính trái đất (km)
	const earthRadius = 6371.0
	return earthRadius * c
}

// checkCondition là phiên bản cải tiến của hàm kiểm tra điều kiện
func checkCondition(participant checkins.AttendanceGroupParticipant, condition checkins.Condition, req checkinsReq.CheckinsReq, ip string) (bool, error) {
	// Thứ tự kiểm tra:
	// 1. Kiểm tra thời gian (nhanh nhất)
	// 2. Kiểm tra nhóm (nhanh)
	// 3. Kiểm tra vị trí (chậm nhất - có tính toán)

	// 1. Kiểm tra thời gian
	now := time.Now().UTC()
	if condition.StartAt != nil && now.Before(*condition.StartAt) {
		return false, errors.New("chưa đến thời gian điểm danh")
	}
	if condition.EndAt != nil && now.After(*condition.EndAt) {
		return false, errors.New("đã hết thời gian điểm danh")
	}

	// 2. Kiểm tra nhóm
	if condition.GroupId != nil {
		if participant.GroupId == nil {
			return false, errors.New("không tìm thấy thông tin nhóm")
		}
		if *condition.GroupId != *participant.GroupId {
			return false, errors.New("thông tin nhóm không thoả điều kiện")
		}
	}

	// 3. Kiểm tra vị trí (chậm nhất)
	if condition.AreaId != nil && condition.Area != nil {
		inArea, areaErr := checkMatchArea(ip, req.Lat, req.Lng, req.Accuracy, *condition.Area)
		if !inArea {
			return false, areaErr
		}
	}

	// Nếu qua được tất cả các điều kiện
	return true, nil
}

// checkConditionsParallel kiểm tra nhiều điều kiện đồng thời để tăng hiệu suất
func checkConditionsParallel(participant checkins.AttendanceGroupParticipant, conditions []checkins.Condition, req checkinsReq.CheckinsReq, ip string) []conditionResult {
	// Tạo channels cho kết quả và giới hạn số goroutine đồng thời
	results := make([]conditionResult, len(conditions))

	// Không chạy quá nhiều goroutine cùng lúc để tránh quá tải hệ thống
	maxConcurrent := 5
	if len(conditions) < maxConcurrent {
		maxConcurrent = len(conditions)
	}

	// Sử dụng worker pool pattern để kiểm soát số lượng goroutine
	var wg sync.WaitGroup
	jobs := make(chan struct {
		index     int
		condition checkins.Condition
	}, len(conditions))

	// Tạo một số lượng worker để xử lý song song
	for w := 0; w < maxConcurrent; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				// Thực hiện kiểm tra điều kiện
				pass, err := checkCondition(participant, job.condition, req, ip)

				// Lưu kết quả vào slice theo đúng vị trí ban đầu
				results[job.index] = conditionResult{
					ConditionID: job.condition.ID,
					Pass:        pass,
					Error:       err,
				}
			}
		}()
	}

	// Đưa công việc vào channel để các worker xử lý
	for i, condition := range conditions {
		jobs <- struct {
			index     int
			condition checkins.Condition
		}{i, condition}
	}
	close(jobs)

	// Đợi tất cả worker hoàn thành
	wg.Wait()

	return results
}

// formatErrorMessage tạo thông báo lỗi chi tiết dựa trên loại điều kiện
func formatErrorMessage(condition checkins.Condition, err error) string {
	if err == nil {
		return ""
	}

	var msg string
	if condition.Area != nil && condition.Area.Area != nil {
		// Tạo thông báo về vị trí
		if condition.Area.Area.Latitude != nil && condition.Area.Area.Longitude != nil {
			if condition.Area.Radius != nil {
				// Có bán kính từ điều kiện khu vực
				msg = fmt.Sprintf("yêu cầu vị trí trong phạm vi %0.0f mét từ %s (%0.6f, %0.6f), ",
					*condition.Area.Radius, condition.Area.Area.Name,
					*condition.Area.Area.Latitude, *condition.Area.Area.Longitude)
			} else if condition.Area.Area.Radius != nil {
				// Có bán kính từ khu vực
				msg = fmt.Sprintf("yêu cầu vị trí trong phạm vi %0.0f mét từ %s (%0.6f, %0.6f), ",
					*condition.Area.Area.Radius, condition.Area.Area.Name,
					*condition.Area.Area.Latitude, *condition.Area.Area.Longitude)
			} else {
				// Không có thông tin bán kính
				msg = fmt.Sprintf("yêu cầu vị trí ở khu vực %s (%0.6f, %0.6f), ",
					condition.Area.Area.Name,
					*condition.Area.Area.Latitude, *condition.Area.Area.Longitude)
			}
		} else {
			// Chỉ có tên khu vực
			msg = fmt.Sprintf("yêu cầu vị trí nằm trong khu vực %s, ",
				condition.Area.Area.Name)
		}
	}

	// Thêm thông tin về thời gian
	if condition.StartAt != nil && condition.EndAt != nil {
		msg = msg + fmt.Sprintf("trong khoảng thời gian từ %s đến %s, ",
			condition.StartAt.Format("15:04 02/01/2006"),
			condition.EndAt.Format("15:04 02/01/2006"))
	} else if condition.StartAt != nil {
		msg = msg + fmt.Sprintf("sau thời điểm %s, ",
			condition.StartAt.Format("15:04 02/01/2006"))
	} else if condition.EndAt != nil {
		msg = msg + fmt.Sprintf("trước thời điểm %s, ",
			condition.EndAt.Format("15:04 02/01/2006"))
	}

	// Ghi log chi tiết về lỗi điều kiện
	if condition.ID > 0 {
		fmt.Printf("Điều kiện #%d thất bại: %s. %s\n", condition.ID, err.Error(), msg)
	}

	if msg != "" {
		return formatStandardError("", err.Error()+". "+msg)
	}
	return formatStandardError("", err.Error())
}

// FormatStandardError định dạng thông báo lỗi theo chuẩn chung
func formatStandardError(code string, message string) string {
	return fmt.Sprintf("[%s] %s", code, message)
}

// isIPAllowed kiểm tra xem một địa chỉ IP có nằm trong các dải IP cho phép không
func isIPAllowed(clientIP string, ipRanges []string) bool {
	// Phân tích địa chỉ IP của người dùng
	clientIPObj := net.ParseIP(clientIP)
	if clientIPObj == nil {
		return false // IP không hợp lệ
	}

	// Tạo cache để tránh phân tích lại các dải IP
	// Cache này chỉ tồn tại trong phạm vi hàm
	cidrCache := make(map[string]*net.IPNet)

	for _, ipRange := range ipRanges {
		ipRange = strings.TrimSpace(ipRange)
		if ipRange == "" {
			continue
		}

		// Trường hợp 1: Dải CIDR (ví dụ: 192.168.1.0/24)
		if strings.Contains(ipRange, "/") {
			// Kiểm tra xem đã phân tích dải này chưa
			ipNet, exists := cidrCache[ipRange]
			if !exists {
				_, ipNetTmp, err := net.ParseCIDR(ipRange)
				if err != nil {
					continue // Bỏ qua dải không hợp lệ
				}
				ipNet = ipNetTmp
				cidrCache[ipRange] = ipNet
			}

			if ipNet.Contains(clientIPObj) {
				return true
			}
		} else {
			// Trường hợp 2: Địa chỉ IP đơn (ví dụ: 192.168.1.5)
			rangeIP := net.ParseIP(ipRange)
			if rangeIP != nil && rangeIP.Equal(clientIPObj) {
				return true
			}

			// Trường hợp 3: Dải IP (ví dụ: 192.168.1.1-192.168.1.10)
			if strings.Contains(ipRange, "-") {
				parts := strings.Split(ipRange, "-")
				if len(parts) == 2 {
					startIP := net.ParseIP(strings.TrimSpace(parts[0]))
					endIP := net.ParseIP(strings.TrimSpace(parts[1]))

					if startIP != nil && endIP != nil && isIPInRange(clientIPObj, startIP, endIP) {
						return true
					}
				}
			}
		}
	}

	return false
}

// isIPInRange kiểm tra xem một IP có nằm trong khoảng [startIP, endIP] không
func isIPInRange(ip, startIP, endIP net.IP) bool {
	// Chuyển đổi sang dạng 4 byte để so sánh
	return bytesCompare(ip.To4(), startIP.To4()) >= 0 && bytesCompare(ip.To4(), endIP.To4()) <= 0
}

// bytesCompare so sánh hai mảng byte
func bytesCompare(a, b []byte) int {
	if a == nil || b == nil || len(a) != len(b) {
		return 0 // Không so sánh được
	}

	for i := 0; i < len(a); i++ {
		if a[i] < b[i] {
			return -1
		} else if a[i] > b[i] {
			return 1
		}
	}

	return 0 // Bằng nhau
}

func (attendanceCheckInService *AttendanceCheckInService) GetAttendanceCheckInLogInfoList(info checkinsReq.AttendanceCheckInSearch) (list []checkins.CheckinLog, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.GVA_DB.Model(&checkins.CheckinLog{})
	var checkinsLog []checkins.CheckinLog

	// if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
	// 	db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	// }

	if info.AttendanceId != nil {
		db = db.Where("attendance_id = ?", *info.AttendanceId)
	}

	if info.Email != nil {
		db = db.Where("email = ?", info.Email)
	}

	// if info.Agent != nil {
	// 	db = db.Where("agent LIKE ?", "%"+*info.Agent+"%")
	// }

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)

	}

	err = db.Order("created_at desc").Debug().Find(&checkinsLog).Error
	return checkinsLog, total, err
}

// Global mutex để đảm bảo thread-safety khi tạo số may mắn
var luckyNumberMutex sync.Mutex

// generateUniqueLuckyNumber tạo số may mắn không trùng lặp cho sự kiện điểm danh
// Hàm này sử dụng nhiều chiến lược khác nhau dựa trên quy mô của sự kiện
func (attendanceCheckInService *AttendanceCheckInService) generateUniqueLuckyNumber(attendanceID uint, participantID uint, localRand *rand.Rand) (int, error) {
	// Sử dụng mutex toàn cục để đảm bảo an toàn đồng thời giữa các request
	luckyNumberMutex.Lock()
	defer luckyNumberMutex.Unlock()

	// Khai báo các hằng số
	// Số lượng bucket trong bloom filter
	const bloomFilterSize = 1048576 // 2^20

	// Xác định kích thước tối đa của số may mắn dựa trên số lượng người tham gia
	var maxLuckyNumber int

	// Thống kê tổng số người tham gia đã có số may mắn
	var totalParticipantsWithLuckyNumber int64
	err := global.GVA_DB.Model(&checkins.AttendanceCheckIn{}).
		Where("attendance_id = ? AND lucky_number > 0", attendanceID).
		Group("partpaticipant_id").
		Count(&totalParticipantsWithLuckyNumber).Error

	if err != nil {
		global.GVA_LOG.Error("Lỗi khi đếm số người dùng đã có số may mắn", zap.Error(err))
	}

	// Đếm tổng số người đã check-in (thay vì đếm toàn bộ người tham gia)
	var totalCheckedInParticipants int64
	err = global.GVA_DB.Model(&checkins.AttendanceCheckIn{}).
		Where("attendance_id = ?", attendanceID).
		Group("partpaticipant_id").
		Count(&totalCheckedInParticipants).Error

	if err != nil {
		global.GVA_LOG.Error("Lỗi khi đếm số người đã check-in", zap.Error(err))
		// Fallback nếu không đếm được chính xác
		totalCheckedInParticipants = totalParticipantsWithLuckyNumber * 2
	}

	// Thêm hệ số dự phòng (3x) để đảm bảo luôn có đủ số may mắn
	const bufferFactor = 3

	// Tính toán số người đã check-in với hệ số dự phòng
	adjustedParticipantCount := totalCheckedInParticipants * bufferFactor

	// Điều chỉnh kích thước tối đa của số may mắn dựa trên số lượng người đã check-in (có dự phòng)
	if adjustedParticipantCount <= 10 {
		// Với số lượng người ít, chỉ cần số nhỏ 2 chữ số (1-99)
		maxLuckyNumber = 99
	} else if adjustedParticipantCount <= 100 {
		// Với số lượng vừa phải, dùng số 3 chữ số (1-999)
		maxLuckyNumber = 999
	} else if adjustedParticipantCount <= 1000 {
		// Với sự kiện lớn hơn, dùng số 4 chữ số (1-9999)
		maxLuckyNumber = 9999
	} else if adjustedParticipantCount <= 10000 {
		// Với sự kiện rất lớn, dùng số 5 chữ số (1-99999)
		maxLuckyNumber = 99999
	} else {
		// Với sự kiện cực lớn, giữ nguyên số 6 chữ số (1-1000000)
		maxLuckyNumber = 1000000
	}

	// Đảm bảo maxLuckyNumber luôn lớn hơn ít nhất 2 lần số người đã check-in
	if int64(maxLuckyNumber) <= totalCheckedInParticipants*2 {
		// Nâng lên mức số tiếp theo nếu không đủ
		if maxLuckyNumber <= 99 {
			maxLuckyNumber = 999
		} else if maxLuckyNumber <= 999 {
			maxLuckyNumber = 9999
		} else if maxLuckyNumber <= 9999 {
			maxLuckyNumber = 99999
		} else if maxLuckyNumber <= 99999 {
			maxLuckyNumber = 1000000
		}
	}

	// Lấy mẫu để xác định chiến lược tối ưu
	var strategy string

	// Nếu số người đã có số may mắn ít hơn 1000, sử dụng chiến lược ngẫu nhiên đơn giản
	if totalParticipantsWithLuckyNumber < 1000 {
		strategy = "random"
	} else if totalParticipantsWithLuckyNumber < 10000 {
		// Nếu có từ 1000-10000 người, sử dụng Bloom Filter + ánh xạ
		strategy = "bloom_filter"
	} else {
		// Cho số lượng lớn, sử dụng phương pháp phân đoạn số để tối ưu
		strategy = "segmented"
	}

	// Ghi log chiến lược được chọn
	global.GVA_LOG.Info("Chiến lược sinh số may mắn được chọn",
		zap.String("strategy", strategy),
		zap.Int64("participants_with_lucky_number", totalParticipantsWithLuckyNumber),
		zap.Int64("checked_in_participants", totalCheckedInParticipants),
		zap.Int64("adjusted_participant_count", adjustedParticipantCount),
		zap.Int("buffer_factor", bufferFactor),
		zap.Uint("attendanceID", attendanceID),
		zap.Int("maxLuckyNumber", maxLuckyNumber))

	// Biến lưu kết quả số may mắn
	var luckyNumber int

	// Dựa vào chiến lược để tạo số may mắn
	switch strategy {
	case "random":
		// Chiến lược 1: Sinh ngẫu nhiên và kiểm tra trùng lặp trực tiếp trong DB
		// Phù hợp cho số lượng người tham gia ít
		maxAttempts := 5 // Giảm số lần thử vì ta sẽ kiểm tra trực tiếp

		for attempts := 0; attempts < maxAttempts; attempts++ {
			// Tạo số may mắn ngẫu nhiên với phân phối đều 1-1000000
			candidateNumber := 1 + localRand.Intn(maxLuckyNumber)

			// Kiểm tra số này đã được dùng hay chưa bằng cách truy vấn trực tiếp
			var count int64
			err := global.GVA_DB.Model(&checkins.AttendanceCheckIn{}).
				Where("attendance_id = ? AND lucky_number = ?", attendanceID, candidateNumber).
				Count(&count).Error

			if err != nil {
				global.GVA_LOG.Error("Lỗi khi kiểm tra số may mắn", zap.Error(err))
				continue
			}

			if count == 0 {
				// Nếu số chưa được dùng, sử dụng nó
				luckyNumber = candidateNumber
				break
			}
		}

	case "bloom_filter":
		// Chiến lược 2: Sử dụng bloom filter để tối ưu kiểm tra
		// Phù hợp cho số lượng người tham gia vừa phải (1K-10K)

		// Lấy tất cả số may mắn đã dùng
		var usedLuckyNumbers []int
		err := global.GVA_DB.Model(&checkins.AttendanceCheckIn{}).
			Where("attendance_id = ? AND lucky_number > 0", attendanceID).
			Distinct("lucky_number").
			Pluck("lucky_number", &usedLuckyNumbers).Error

		if err != nil {
			global.GVA_LOG.Error("Lỗi khi lấy danh sách số may mắn", zap.Error(err))
		}

		// Tạo bloom filter đơn giản (mảng bit)
		bloomFilter := make([]bool, bloomFilterSize)

		// Hash function đơn giản cho bloom filter
		hashFunc := func(num int) int {
			return num % bloomFilterSize
		}

		// Tạo map để kiểm tra chính xác
		usedLuckyNumbersMap := make(map[int]bool)

		// Đánh dấu các số đã dùng vào bloom filter và map
		for _, num := range usedLuckyNumbers {
			bloomFilter[hashFunc(num)] = true
			usedLuckyNumbersMap[num] = true
		}

		// Số lần thử tối đa
		maxAttempts := 20

		// Tạo số may mắn không trùng
		for attempts := 0; attempts < maxAttempts; attempts++ {
			candidateNumber := 1 + localRand.Intn(maxLuckyNumber)

			// Kiểm tra nhanh với bloom filter trước
			if !bloomFilter[hashFunc(candidateNumber)] {
				// Nếu bloom filter nói không có, chắc chắn là chưa dùng
				luckyNumber = candidateNumber
				break
			} else if !usedLuckyNumbersMap[candidateNumber] {
				// Nếu bloom filter nói có nhưng map nói không (false positive), vẫn ok
				luckyNumber = candidateNumber
				break
			}
		}

	case "segmented":
		// Chiến lược 3: Phân đoạn số và tìm khoảng trống
		// Phù hợp cho số lượng rất lớn (>10K)

		// Tính số phân đoạn dựa trên số lượng người tham gia
		numSegments := 10 + int(totalParticipantsWithLuckyNumber/1000)
		if numSegments > 100 {
			numSegments = 100 // Giới hạn số phân đoạn
		}

		// Kích thước mỗi phân đoạn
		segmentSize := maxLuckyNumber / numSegments

		// Đếm số lượng số đã dùng trong mỗi phân đoạn
		segments := make([]int64, numSegments)

		for i := 0; i < numSegments; i++ {
			segmentStart := i*segmentSize + 1
			segmentEnd := (i + 1) * segmentSize

			// Đếm số lượng số may mắn đã dùng trong phân đoạn này
			err := global.GVA_DB.Model(&checkins.AttendanceCheckIn{}).
				Where("attendance_id = ? AND lucky_number >= ? AND lucky_number <= ?",
					attendanceID, segmentStart, segmentEnd).
				Distinct("lucky_number").
				Count(&segments[i]).Error

			if err != nil {
				global.GVA_LOG.Error("Lỗi khi đếm số may mắn theo phân đoạn", zap.Error(err))
			}
		}

		// Tìm phân đoạn có ít số đã dùng nhất
		leastUsedSegment := 0
		leastUsedCount := segments[0]

		for i := 1; i < numSegments; i++ {
			if segments[i] < leastUsedCount {
				leastUsedSegment = i
				leastUsedCount = segments[i]
			}
		}

		// Tỷ lệ sử dụng của phân đoạn này
		usageRatio := float64(leastUsedCount) / float64(segmentSize)

		// Nếu phân đoạn đã đầy trên 80%, sử dụng phương pháp khác
		if usageRatio >= 0.8 {
			// Sinh số ngẫu nhiên và kiểm tra trực tiếp
			for attempts := 0; attempts < 10; attempts++ {
				candidateNumber := 1 + localRand.Intn(maxLuckyNumber)

				var count int64
				err := global.GVA_DB.Model(&checkins.AttendanceCheckIn{}).
					Where("attendance_id = ? AND lucky_number = ?", attendanceID, candidateNumber).
					Count(&count).Error

				if err == nil && count == 0 {
					luckyNumber = candidateNumber
					break
				}
			}
		} else {
			// Sinh số ngẫu nhiên trong phân đoạn có ít số đã dùng nhất
			segmentStart := leastUsedSegment*segmentSize + 1
			segmentEnd := (leastUsedSegment + 1) * segmentSize

			// Lấy danh sách số đã dùng trong phân đoạn
			var usedInSegment []int
			err := global.GVA_DB.Model(&checkins.AttendanceCheckIn{}).
				Where("attendance_id = ? AND lucky_number >= ? AND lucky_number <= ?",
					attendanceID, segmentStart, segmentEnd).
				Distinct("lucky_number").
				Pluck("lucky_number", &usedInSegment).Error

			if err != nil {
				global.GVA_LOG.Error("Lỗi khi lấy số đã dùng trong phân đoạn", zap.Error(err))
			}

			// Tạo map cho phân đoạn
			usedInSegmentMap := make(map[int]bool)
			for _, num := range usedInSegment {
				usedInSegmentMap[num] = true
			}

			// Sinh số trong phân đoạn này
			maxAttempts := 20
			for attempts := 0; attempts < maxAttempts; attempts++ {
				candidateNumber := segmentStart + localRand.Intn(segmentEnd-segmentStart+1)

				if !usedInSegmentMap[candidateNumber] {
					luckyNumber = candidateNumber
					break
				}
			}
		}
	}

	// Nếu các phương pháp trên không tạo được số, sử dụng phương pháp cuối cùng
	if luckyNumber == 0 {
		// Tạo UUID-based number để đảm bảo không trùng lặp
		// Sử dụng timestamp nano + ID người dùng để tạo số đặc biệt
		timestamp := time.Now().UnixNano()
		// Điều chỉnh phương pháp tạo số đặc biệt dựa trên maxLuckyNumber
		if maxLuckyNumber <= 99 {
			// Với số 2 chữ số, đảm bảo số nằm trong khoảng 1-99
			specialNumber := (timestamp % 90) + 10 + int64(participantID%10)
			luckyNumber = int(specialNumber)
		} else if maxLuckyNumber <= 999 {
			// Với số 3 chữ số, đảm bảo số nằm trong khoảng 100-999
			specialNumber := (timestamp % 900) + 100 + int64(participantID%100)
			luckyNumber = int(specialNumber)
		} else if maxLuckyNumber <= 9999 {
			// Với số 4 chữ số, đảm bảo số nằm trong khoảng 1000-9999
			specialNumber := (timestamp % 9000) + 1000 + int64(participantID%100)*10
			luckyNumber = int(specialNumber)
		} else if maxLuckyNumber <= 99999 {
			// Với số 5 chữ số, đảm bảo số nằm trong khoảng 10000-99999
			specialNumber := (timestamp % 90000) + 10000 + int64(participantID%100)*100
			luckyNumber = int(specialNumber)
		} else {
			// Với số 6 chữ số, giữ nguyên phương pháp cũ
			specialNumber := (timestamp % 900000) + 100000 + int64(participantID%1000)*1000
			luckyNumber = int(specialNumber % 1000000)
			if luckyNumber < 100000 {
				luckyNumber += 100000 // Đảm bảo ít nhất 6 chữ số
			}
		}

		global.GVA_LOG.Warn("Sử dụng phương pháp đặc biệt để tạo số may mắn",
			zap.Int("luckyNumber", luckyNumber),
			zap.Uint("attendanceID", attendanceID),
			zap.Int("maxLuckyNumber", maxLuckyNumber))
	}

	return luckyNumber, nil
}
