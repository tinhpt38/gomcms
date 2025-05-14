package checkins

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
	checkinsReq "github.com/flipped-aurora/gin-vue-admin/server/model/checkins/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/config"
	cfgService "github.com/flipped-aurora/gin-vue-admin/server/service/config"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ParticipantService struct{}

// CreateParticipant 创建Sinh viên (Người tham dự phiên điểm danh)记录
// Author [piexlmax](https://github.com/piexlmax)
func (participantService *ParticipantService) CreateParticipant(participant *checkins.Participant) (err error) {
	// err = global.GVA_DB.FirstOrCreate(participant).Error

	err = global.GVA_DB.Model(&checkins.Participant{}).Where(&checkins.Participant{
		Email: participant.Email,
	}).FirstOrCreate(participant).Error

	return err
}

func (participantService *ParticipantService) BulkCreateParticipants(req checkinsReq.ListEmailParticipantRequest) (err error) {

	for _, email := range req.List {
		participant := checkins.Participant{
			Email: email,
		}
		err = global.GVA_DB.Model(&checkins.Participant{}).Where(&checkins.Participant{
			Email: email,
		}).FirstOrCreate(&participant).Error
		if err != nil {
			return err
		}
		agp := checkins.AttendanceGroupParticipant{
			ParticipantId: &participant.ID,
			AttendanceId:  req.AttendanceId,
			GroupId:       req.GroupId,
		}
		err = global.GVA_DB.Model(&checkins.AttendanceGroupParticipant{}).Where(&checkins.AttendanceGroupParticipant{
			ParticipantId: &participant.ID,
			GroupId:       req.GroupId,
			AttendanceId:  req.AttendanceId,
		}).FirstOrCreate(&agp).Error
	}

	return
	// err = global.GVA_DB.FirstOrCreate(participant).Error

	// err = global.GVA_DB.Model(&checkins.Participant{}).Where(&checkins.Participant{
	// 	Email: participant.Email,
	// }).FirstOrCreate(participant).Error

}

// DeleteParticipant 删除Sinh viên (Người tham dự phiên điểm danh)记录
// Author [piexlmax](https://github.com/piexlmax)
func (participantService *ParticipantService) DeleteParticipant(ID string) (err error) {
	err = global.GVA_DB.Delete(&checkins.Participant{}, "id = ?", ID).Error
	return err
}

func (participantService *ParticipantService) DeleteParticipantInAttendance(ID string, attId uint) (err error) {
	err = global.GVA_DB.Delete(&checkins.AttendanceGroupParticipant{}, "participant_id = ? AND attendance_id = ?", ID, attId).Debug().Error
	return err
}

// DeleteParticipantByIds 批量删除Sinh viên (Người tham dự phiên điểm danh)记录
// Author [piexlmax](https://github.com/piexlmax)
func (participantService *ParticipantService) DeleteParticipantByIds(IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]checkins.Participant{}, "id in ?", IDs).Error
	return err
}

// UpdateParticipant 更新Sinh viên (Người tham dự phiên điểm danh)记录
// Author [piexlmax](https://github.com/piexlmax)
func (participantService *ParticipantService) UpdateParticipant(participant checkins.Participant) (err error) {

	global.GVA_DB.Model(&checkins.Participant{}).Where("id = ?", participant.ID).Save(&participant)
	agpDb := global.GVA_DB.Table(checkins.AttendanceGroupParticipant{}.TableName())
	var existingGroupIds []uint
	if err := agpDb.Where("participant_id = ? AND attendance_id = ?", participant.ID, participant.AttendanceId).
		Pluck("group_id", &existingGroupIds).Error; err != nil {
		return err
	}

	newGroupIds := *participant.GroupId

	toAdd := difference(newGroupIds, existingGroupIds)

	// Tìm những GroupId cần xóa
	toDelete := difference(existingGroupIds, newGroupIds)

	if len(toDelete) > 0 {
		if err := agpDb.Where("attendance_id = ? AND group_id IN (?)", participant.AttendanceId, toDelete).
			Unscoped().
			Delete(&checkins.AttendanceGroupParticipant{}).Error; err != nil {
			return err
		}
	}
	// Tạo mới theo danh sách
	for _, grId := range toAdd {
		agbDb := global.GVA_DB.Table(checkins.AttendanceGroupParticipant{}.TableName())
		agp := checkins.AttendanceGroupParticipant{
			ParticipantId: &participant.ID,
			AttendanceId:  participant.AttendanceId,
			GroupId:       &grId,
		}
		rerr := agbDb.Create(&agp).Error
		if rerr != nil {
			return rerr
		}
	}

	return err
}

func difference(slice1, slice2 []uint) []uint {
	diff := []uint{}
	// Tạo một map để tra cứu nhanh
	lookup := make(map[uint]struct{}, len(slice2))
	for _, v := range slice2 {
		lookup[v] = struct{}{}
	}

	// Tìm những phần tử không có trong map
	for _, v := range slice1 {
		if _, found := lookup[v]; !found {
			diff = append(diff, v)
		}
	}

	return diff
}

func (participantService *ParticipantService) GetParticipant(ID string) (participant checkins.Participant, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&participant).Error
	return
}

func (participantService *ParticipantService) GetLuckyParticipant(acId string) (participant checkins.Participant, luckyNumber int, err error) {
	// Lấy lịch sử số may mắn đã quay để không lặp lại
	var usedLuckyNumbers []checkins.UsedLuckyParticipant
	err = global.GVA_DB.Where("attendance_id = ?", acId).
		Order("created_at DESC").
		Find(&usedLuckyNumbers).Error

	if err != nil {
		global.GVA_LOG.Error("Lỗi khi lấy lịch sử số may mắn", zap.Error(err))
	}

	// Tạo map để lưu các số đã hiển thị
	usedNumbers := make(map[int]bool)
	for _, used := range usedLuckyNumbers {
		if used.LuckyNumber != nil {
			usedNumbers[*used.LuckyNumber] = true
		}
	}

	// Lấy số lần đã quay cho phiên điểm danh này
	var spinCount int64
	err = global.GVA_DB.Model(&checkins.UsedLuckyParticipant{}).
		Where("attendance_id = ?", acId).
		Count(&spinCount).Error

	if err != nil {
		global.GVA_LOG.Error("Lỗi khi đếm số lần quay", zap.Error(err))
		return participant, 0, err
	}
	// Lấy danh sách tất cả các số may mắn trong DB (có thể trúng)
	var existingLuckyNumbers []int
	err = global.GVA_DB.Model(&checkins.AttendanceCheckIn{}).
		Where("attendance_id = ? AND lucky_number IS NOT NULL", acId).
		Distinct().
		Pluck("lucky_number", &existingLuckyNumbers).Error

	if err != nil {
		global.GVA_LOG.Error("Có lỗi khi truy vấn số may mắn", zap.Error(err))
		return participant, 0, err
	}

	// Tạo random number generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Nếu chưa có người thắng, tiếp tục xử lý
	var chosenNumber int
	foundLuckyParticipant := false

	// Chiến lược quay số dựa trên số lần đã quay:
	// Lần 1: Chọn số không có trong DB (không có người trúng)
	// Lần 2: 50/50 giữa số trúng và không trúng
	// Lần 3: Đảm bảo có người trúng

	if spinCount >= 2 { // Lần quay thứ 3 trở đi: đảm bảo có người trúng
		// Tìm một số có trong database và chưa được hiển thị
		var potentialWinningNumbers []int
		for _, num := range existingLuckyNumbers {
			if !usedNumbers[num] {
				potentialWinningNumbers = append(potentialWinningNumbers, num)
			}
		}

		if len(potentialWinningNumbers) > 0 {
			// Chọn ngẫu nhiên từ các số có thể thắng chưa hiển thị
			idx := r.Intn(len(potentialWinningNumbers))
			chosenNumber = potentialWinningNumbers[idx]
		} else if len(existingLuckyNumbers) > 0 {
			// Nếu đã hiển thị hết các số có thể thắng, chọn ngẫu nhiên một người tham dự bất kỳ
			var randomParticipant checkins.AttendanceGroupParticipant
			randomParticipantErr := global.GVA_DB.Where("attendance_id = ?", acId).
				Order("RAND()").
				First(&randomParticipant).Error

			if randomParticipantErr == nil && randomParticipant.ParticipantId != nil {
				// Lấy thông tin người tham dự
				participantErr := global.GVA_DB.Where("id = ?", *randomParticipant.ParticipantId).
					First(&participant).Error

				if participantErr == nil {
					// Tìm một số chưa được hiển thị để gán cho người này
					for i := 1; i <= 99; i++ {
						if !usedNumbers[i] {
							chosenNumber = i
							// Lưu kết quả vào lịch sử
							usedLuckyParticipant := checkins.UsedLuckyParticipant{
								AttendanceId:  acId,
								ParticipantId: &participant.ID,
								LuckyNumber:   &chosenNumber,
							}
							global.GVA_DB.Create(&usedLuckyParticipant)
							return participant, chosenNumber, nil
						}
					}
					// Nếu đã hiển thị hết tất cả các số, reset lại
					chosenNumber = r.Intn(99) + 1
					usedLuckyParticipant := checkins.UsedLuckyParticipant{
						AttendanceId:  acId,
						ParticipantId: &participant.ID,
						LuckyNumber:   &chosenNumber,
					}
					global.GVA_DB.Create(&usedLuckyParticipant)
					return participant, chosenNumber, nil
				}
			}
			// Nếu không tìm được người tham dự, chọn một số chưa hiển thị
			for i := 1; i <= 99; i++ {
				if !usedNumbers[i] {
					chosenNumber = i
					break
				}
			}
			// Nếu đã hiển thị hết tất cả các số, chọn ngẫu nhiên
			if chosenNumber == 0 {
				chosenNumber = r.Intn(99) + 1
			}
		}
	} else if spinCount == 1 { // Lần quay thứ 2: 50/50
		// Tạo hai mảng: số thắng và số không thắng (chưa hiển thị)
		var winningNumbers, nonWinningNumbers []int

		// Số thắng là các số trong existingLuckyNumbers và chưa hiển thị
		for _, num := range existingLuckyNumbers {
			if !usedNumbers[num] {
				winningNumbers = append(winningNumbers, num)
			}
		}

		// Số không thắng là các số từ 1-99 trừ đi các số thắng và chưa hiển thị
		existingMap := make(map[int]bool)
		for _, num := range existingLuckyNumbers {
			existingMap[num] = true
		}

		for i := 1; i <= 99; i++ {
			if !existingMap[i] && !usedNumbers[i] {
				nonWinningNumbers = append(nonWinningNumbers, i)
			}
		}

		// 50/50 giữa số thắng và không thắng
		if len(winningNumbers) > 0 && len(nonWinningNumbers) > 0 {
			if r.Intn(2) == 0 { // 50% cơ hội chọn số thắng
				idx := r.Intn(len(winningNumbers))
				chosenNumber = winningNumbers[idx]
			} else { // 50% cơ hội chọn số không thắng
				idx := r.Intn(len(nonWinningNumbers))
				chosenNumber = nonWinningNumbers[idx]
			}
		} else if len(winningNumbers) > 0 {
			idx := r.Intn(len(winningNumbers))
			chosenNumber = winningNumbers[idx]
		} else if len(nonWinningNumbers) > 0 {
			idx := r.Intn(len(nonWinningNumbers))
			chosenNumber = nonWinningNumbers[idx]
		} else {
			// Nếu không còn số chưa hiển thị, reset lại
			// Tìm một số bất kỳ chưa được sử dụng
			var availableNumbers []int
			for i := 1; i <= 99; i++ {
				if !usedNumbers[i] {
					availableNumbers = append(availableNumbers, i)
				}
			}

			if len(availableNumbers) > 0 {
				idx := r.Intn(len(availableNumbers))
				chosenNumber = availableNumbers[idx]
			} else {
				// Nếu đã hiển thị hết tất cả các số, chọn ngẫu nhiên
				chosenNumber = r.Intn(99) + 1
			}
		}
	} else { // Lần quay đầu tiên: ưu tiên số không thắng
		// Tìm các số không có trong DB và chưa hiển thị
		var nonWinningNumbers []int

		// Số không thắng là các số từ 1-99 trừ đi các số trong DB và chưa hiển thị
		existingMap := make(map[int]bool)
		for _, num := range existingLuckyNumbers {
			existingMap[num] = true
		}

		for i := 1; i <= 99; i++ {
			if !existingMap[i] && !usedNumbers[i] {
				nonWinningNumbers = append(nonWinningNumbers, i)
			}
		}

		if len(nonWinningNumbers) > 0 {
			// Chọn ngẫu nhiên từ các số không thắng
			idx := r.Intn(len(nonWinningNumbers))
			chosenNumber = nonWinningNumbers[idx]
		} else {
			// Nếu không còn số không thắng chưa hiển thị, tìm số bất kỳ chưa hiển thị
			var availableNumbers []int
			for i := 1; i <= 99; i++ {
				if !usedNumbers[i] {
					availableNumbers = append(availableNumbers, i)
				}
			}

			if len(availableNumbers) > 0 {
				idx := r.Intn(len(availableNumbers))
				chosenNumber = availableNumbers[idx]
			} else {
				// Nếu đã hiển thị hết tất cả các số, reset và chọn ngẫu nhiên
				chosenNumber = r.Intn(99) + 1
			}
		}
	}

	// Kiểm tra xem số đã chọn có người thắng không
	var checkIn checkins.AttendanceCheckIn
	checkInErr := global.GVA_DB.Where("attendance_id = ? AND lucky_number = ?", acId, chosenNumber).
		Order("created_at DESC").
		First(&checkIn).Error

	luckyNumber = chosenNumber

	if checkInErr == nil && checkIn.ID > 0 {
		// Tìm thấy người tham gia có số may mắn này
		dbErr := global.GVA_DB.Where("id = ?", checkIn.PartpaticipantId).
			First(&participant).Error

		if dbErr == nil {
			foundLuckyParticipant = true

			// Lưu vào lịch sử
			usedLuckyParticipant := checkins.UsedLuckyParticipant{
				AttendanceId:  acId,
				ParticipantId: &participant.ID,
				LuckyNumber:   &chosenNumber,
			}
			global.GVA_DB.Create(&usedLuckyParticipant)
		}
	}

	// Nếu không tìm thấy người thắng, lưu vào lịch sử không có người thắng
	if !foundLuckyParticipant {
		usedLuckyParticipant := checkins.UsedLuckyParticipant{
			AttendanceId: acId,
			LuckyNumber:  &chosenNumber,
		}
		global.GVA_DB.Create(&usedLuckyParticipant)
	}

	return participant, luckyNumber, nil
}

// Helper function to check if a slice contains a value
func contains(slice []uint, value uint) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func (participantService *ParticipantService) GetParticipantByEmail(email string) (participant checkins.Participant, err error) {
	err = global.GVA_DB.Where("email = ?", email).First(&participant).Error
	return
}

// func (participantService *ParticipantService) GetParticipantInAttendance(participantId uint, attendanceId uint) (memberOfAttendance checkins.AttendanceGroupParticipant, err error) {
// 	db := global.GVA_DB.Table(checkins.AttendanceGroupParticipant{}.TableName())
// 	err = db.Joins("left join `groups` on groups.id = attendance_group_participants.group_id").
// 		Where("attendance_group_participants.participant_id = ? AND attendance_group_participants.attendance_id = ? AND groups.attendance_id = ?", participantId, attendanceId, attendanceId).
// 		Where("attendance_group_participants.deleted_at IS NULL").
// 		Order("attendance_group_participants.id").
// 		Preload(clause.Associations).
// 		First(&memberOfAttendance).
// 		Error
// 	return
// }

func (participantService *ParticipantService) GetParticipantInAttendance(participantId uint, attendanceId uint) (memberOfAttendance []checkins.AttendanceGroupParticipant, err error) {
	db := global.GVA_DB.Table(checkins.AttendanceGroupParticipant{}.TableName())
	err = db.Joins("left join `groups` on groups.id = attendance_group_participants.group_id").
		Where("attendance_group_participants.participant_id = ? AND attendance_group_participants.attendance_id = ? AND groups.attendance_id = ?", participantId, attendanceId, attendanceId).
		Where("attendance_group_participants.deleted_at IS NULL").
		Order("attendance_group_participants.id").
		Preload(clause.Associations).
		Find(&memberOfAttendance).
		Error
	return
}

// GetParticipantInfoList 分页获取Sinh viên (Người tham dự phiên điểm danh)记录
// Author [piexlmax](https://github.com/piexlmax)
func (participantService *ParticipantService) GetParticipantInfoList(info checkinsReq.ParticipantSearch) (list []checkins.Participant, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&checkins.Participant{})
	var participants []checkins.Participant
	// Nếu có điều kiện tìm kiếm, câu lệnh tìm kiếm sẽ được tạo tự động ở dưới đây
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.FullName != "" {
		db = db.Where("full_name LIKE ?", "%"+info.FullName+"%")
	}
	if info.Email != "" {
		db = db.Where("email LIKE ?", "%"+info.Email+"%")
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&participants).Error
	return participants, total, err
}

func (participantService *ParticipantService) GetParticipantInfoListByAttendance(info checkinsReq.ParticipantSearch) (list []checkins.Participant, total int64, err error) {

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	agpDb := global.GVA_DB.Model(&checkins.AttendanceGroupParticipant{})
	var parIds []uint
	if info.GroupId != nil && *info.GroupId != 0 {
		agpDb = agpDb.Where("group_id = ?", *info.GroupId)
	}
	if info.AttendanceId != nil && *info.AttendanceId != 0 {
		agpDb = agpDb.Where("attendance_id = ?", *info.AttendanceId)
	}

	err = agpDb.Select("DISTINCT participant_id").Find(&parIds).Error
	if err != nil {
		return
	}

	db := global.GVA_DB.Model(&checkins.Participant{})
	var participants []checkins.Participant
	db = db.Where("id IN (?)", parIds)
	if info.FullName != "" {
		db = db.Where("full_name LIKE ?", "%"+info.FullName+"%")
	}
	if info.Email != "" {
		db = db.Where("email LIKE ?", "%"+info.Email+"%")
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Preload("Groups", "attendance_id = ?", info.AttendanceId).Debug().Find(&participants).Error
	// Xử lý lấy thông tin điều kiện điểm danh
	// newList := participantService.GetMetadata(participants, *info.AttendanceId)
	return participants, total, err
}

func (participantService *ParticipantService) ImportExcel(info config.CfgFileProcess) (err error) {
	path := global.GVA_CONFIG.Local.Path + "/" + info.UniqueFileName

	file, err := excelize.OpenFile(path)

	if err != nil {
		return
	}

	// * Close file after read
	defer func() {
		// Close the spreadsheet

		if err := file.Close(); err != nil {
			fmt.Println(err)
		}

		err := os.Remove(path)

		if err != nil {
			global.GVA_LOG.Error("Có lỗi khi xoá tập tin excel: "+info.FileName, zap.Error(err))
		}
	}()

	sheetIndex := 0
	sheet := file.WorkBook.Sheets.Sheet[sheetIndex].Name //! Get by sheet index
	rows, err := file.GetRows(sheet)
	if err != nil {
		return err
	}
	uuid := info.Uuid

	fileProcessService := cfgService.CfgFileProcessServiceApp

	var beginDataRowIndex int = 1
	var lenRows float64 = float64(len(rows) - beginDataRowIndex)

	// Kiểm tra trường hợp file trống
	if lenRows < 1 {
		fileProcessService.UpdateCfgFileProcess(config.CfgFileProcess{
			Uuid:     uuid,
			Percent:  100,
			Msg:      "Tập tin trống, vui lòng bổ sung dữ liệu trước khi nhập",
			Status:   config.FILE_PROCESS_STATUS_ERROR,
			TotalRow: int(lenRows),
		})
	}

	// Dữ liệu group
	var groups []checkins.Group
	err = global.GVA_DB.Find(&groups).Error

	if err != nil {
		return
	}

	groupMap := make(map[string]checkins.Group)
	for _, group := range groups {
		groupMap[group.Name] = group
	}

	ratio := 100.0 / lenRows // Tỷ lệ phần trăm
	validationPass := true
	sheetHeaders := rows[0] // Bỏ qua dòng đầu tiên
	errorRowAccept := lenRows
	columnsLength := len(sheetHeaders)

	if len(sheetHeaders) != columnsLength {
		errEntry := config.FileProcessError{
			GVA_MODEL:       global.GVA_MODEL{},
			FileProcessId:   &info.ID,
			FileProcessUuid: info.Uuid,
			FieldTitle:      "Lỗi",
			ExpectedValue:   "",
			ReceivedValue:   "",
			Note:            "Định dạnh dữ liệu không hợp lệ, kiểm tra tập tin mẫu",
			Row:             1,
		}

		err = global.GVA_DB.Create(&errEntry).Error

		if err != nil {
			return
		}

		fileProcessService.UpdateCfgFileProcess(config.CfgFileProcess{
			Uuid:     uuid,
			Percent:  100,
			Msg:      "Định dạnh dữ liệu không hợp lệ, kiểm tra tập tin mẫu",
			Status:   config.FILE_PROCESS_STATUS_ERROR,
			TotalRow: int(lenRows),
		})

		return errors.New("Định dạnh dữ liệu không hợp lệ, kiểm tra tập tin mẫu")
	}

	if lenRows > 100 {
		errorRowAccept = 20
	}
	if lenRows > 1000 {
		errorRowAccept = 10
	}

	fileProcessService.UpdateCfgFileProcess(config.CfgFileProcess{
		Uuid:     uuid,
		Percent:  0,
		Status:   config.FILE_PROCESS_STATUS_VALIDATING,
		TotalRow: int(lenRows),
	})

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		rowCounter := beginDataRowIndex + 1
		errorCounter := 0

		for _, row := range rows[beginDataRowIndex:] {
			//! Get data from row
			fullName := utils.GetArrayValue(row, 0)
			email := utils.GetArrayValue(row, 1)
			group := utils.GetArrayValue(row, 2)

			// error entries
			errorEntries := make([]config.FileProcessError, 0)

			if fullName == "" {
				errorEntries = append(errorEntries, config.FileProcessError{
					GVA_MODEL:       global.GVA_MODEL{},
					FileProcessId:   &info.ID,
					FileProcessUuid: uuid,
					FieldTitle:      sheetHeaders[0],
					ExpectedValue:   "Họ và tên",
					ReceivedValue:   fullName,
					Note:            "Họ và tên không được để trống",
					Row:             rowCounter,
				})
			}

			if email == "" {
				errorEntries = append(errorEntries, config.FileProcessError{
					GVA_MODEL:       global.GVA_MODEL{},
					FileProcessId:   &info.ID,
					FileProcessUuid: uuid,
					FieldTitle:      sheetHeaders[1],
					ExpectedValue:   "Email",
					ReceivedValue:   email,
					Note:            "Email không được để trống",
					Row:             rowCounter,
				})
			}

			if group != "" {
				if groupMap[group].ID == 0 {
					errorEntries = append(errorEntries, config.FileProcessError{
						GVA_MODEL:       global.GVA_MODEL{},
						FileProcessId:   &info.ID,
						FileProcessUuid: uuid,
						FieldTitle:      sheetHeaders[2],
						ExpectedValue:   "Nhóm",
						ReceivedValue:   group,
						Note:            "Nhóm không tồn tại",
						Row:             rowCounter,
					})
				}
			}

			// Error entries
			if len(errorEntries) > 0 {
				validationPass = false
				errorCounter += 1

				err = tx.Model(&config.FileProcessError{}).Create(&errorEntries).Error
				if err != nil {
					return err
				}

				if errorCounter >= int(errorRowAccept) {
					break
				}
			}

			// Increment row counter
			rowCounter += 1
		}

		return nil
	})

	if err != nil {
		fileProcessService.UpdateCfgFileProcess(config.CfgFileProcess{
			Uuid:     uuid,
			Percent:  100,
			Status:   config.FILE_PROCESS_STATUS_ERROR,
			Msg:      "Lỗi xảy ra khi thực hiện lưu thông tin xác thực dữ liệu, liên hệ quản trị viên",
			TotalRow: int(lenRows),
		})
	}

	if !validationPass {
		fileProcessService.UpdateCfgFileProcess(config.CfgFileProcess{
			Uuid:     uuid,
			Percent:  100,
			Status:   config.FILE_PROCESS_STATUS_ERROR,
			Msg:      "Lỗi xác thực dữ liệu, vui lòng chỉnh sửa tập tin excel theo danh sách lỗi chi tiết",
			TotalRow: int(lenRows),
		})

		return errors.New("Lỗi xác thực dữ liệu, vui lòng chỉnh sửa tập tin excel theo danh sách lỗi chi tiết")
	}

	fileProcessService.UpdateCfgFileProcess(config.CfgFileProcess{
		Uuid:     uuid,
		Percent:  0,
		Status:   config.FILE_PROCESS_STATUS_PROCESSING,
		TotalRow: int(lenRows),
		Msg:      "Tập tin đang được xử lý",
	})

	count := 0

	for _, row := range rows[beginDataRowIndex:] {
		//! Get data from row
		fullName := utils.GetArrayValue(row, 0)
		email := utils.GetArrayValue(row, 1)
		groupName := utils.GetArrayValue(row, 2)

		itemToSave := checkins.Participant{
			GVA_MODEL: global.GVA_MODEL{},
			FullName:  &fullName,
			Email:     email,
		}

		// Increment count
		count += 1
		percent := utils.RoundFloat(float64(count)*ratio, 2)

		result := global.GVA_DB.Model(&checkins.Participant{}).Where(&checkins.Participant{
			Email: email,
		}).FirstOrCreate(&itemToSave)

		// No record created
		if result.RowsAffected == 0 {
			continue
		}

		err = result.Error
		if err != nil {
			errorEntry := config.FileProcessError{
				GVA_MODEL:       global.GVA_MODEL{},
				FileProcessId:   &info.ID,
				FileProcessUuid: uuid,
				FieldTitle:      "Lỗi",
				ExpectedValue:   "",
				ReceivedValue:   "",
				Note:            "Lỗi xảy ra khi thực hiện lưu dữ liệu " + err.Error(),
				Row:             count,
			}

			global.GVA_DB.Model(&config.FileProcessError{}).Create(&errorEntry)

			fileProcessService.UpdateCfgFileProcess(config.CfgFileProcess{
				Uuid:    uuid,
				Percent: 100,
				Status:  config.FILE_PROCESS_STATUS_ERROR,
				Msg:     "Đã có lỗi xảy ra, vui lòng thử lại sau",
			})

			return
		}
		var groupId *uint
		if groupName != "" {
			groupWillCreate := checkins.Group{
				GVA_MODEL:    global.GVA_MODEL{},
				Name:         groupName,
				AttendanceId: info.AttendanceId,
			}
			groupResult := global.GVA_DB.Model(&checkins.Group{}).Where(&checkins.Group{
				Name:         groupName,
				AttendanceId: info.AttendanceId,
			}).FirstOrCreate(&groupWillCreate)

			err = groupResult.Error
			if err != nil {
				errorEntry := config.FileProcessError{
					GVA_MODEL:       global.GVA_MODEL{},
					FileProcessId:   &info.ID,
					FileProcessUuid: uuid,
					FieldTitle:      "Lỗi",
					ExpectedValue:   "",
					ReceivedValue:   "",
					Note:            "Lỗi xảy ra khi thực hiện lưu dữ liệu tại Group " + err.Error(),
					Row:             count,
				}

				global.GVA_DB.Model(&config.FileProcessError{}).Create(&errorEntry)

				fileProcessService.UpdateCfgFileProcess(config.CfgFileProcess{
					Uuid:    uuid,
					Percent: 100,
					Status:  config.FILE_PROCESS_STATUS_ERROR,
					Msg:     "Đã có lỗi xảy ra, vui lòng thử lại sau",
				})

				return
			}
			groupId = &groupWillCreate.ID

		}

		resultMany := global.GVA_DB.Model(&checkins.AttendanceGroupParticipant{}).Create(&checkins.AttendanceGroupParticipant{
			AttendanceId:  &info.AttendanceId,
			ParticipantId: &itemToSave.ID,
			GroupId:       groupId,
		})

		err = resultMany.Error
		if err != nil {
			errorEntry := config.FileProcessError{
				GVA_MODEL:       global.GVA_MODEL{},
				FileProcessId:   &info.ID,
				FileProcessUuid: uuid,
				FieldTitle:      "Lỗi",
				ExpectedValue:   "",
				ReceivedValue:   "",
				Note:            "Lỗi xảy ra khi thực hiện lưu dữ liệu tại Attendance - Group - Participant " + err.Error(),
				Row:             count,
			}

			global.GVA_DB.Model(&config.FileProcessError{}).Create(&errorEntry)

			fileProcessService.UpdateCfgFileProcess(config.CfgFileProcess{
				Uuid:    uuid,
				Percent: 100,
				Status:  config.FILE_PROCESS_STATUS_ERROR,
				Msg:     "Đã có lỗi xảy ra, vui lòng thử lại sau",
			})

			return
		}

		if count%5 == 0 {
			// update file process
			fileProcessService.UpdateCfgFileProcess(config.CfgFileProcess{
				Uuid:    uuid,
				Percent: percent,
				Status:  config.FILE_PROCESS_STATUS_PROCESSING,
			})
		}

		// condition return success
		if count == int(lenRows) {
			fileProcessService.UpdateCfgFileProcess(config.CfgFileProcess{
				Uuid:    uuid,
				Percent: 100,
				Status:  config.FILE_PROCESS_STATUS_FINISH,
				Msg:     "Đã xử lý xong file excel",
			})
		}
	}

	fileProcessService.UpdateCfgFileProcess(config.CfgFileProcess{
		Uuid:    uuid,
		Percent: 100,
		Status:  config.FILE_PROCESS_STATUS_FINISH,
		Msg:     "Đã xử lý xong file excel",
	})

	return nil
}

// func (participantService *ParticipantService) GetMetadata(list []checkins.Participant, attId uint) (newList []checkins.Participant) {
// 	conditionService := new(ConditionService)

// 	conditions, _ := conditionService.GetConditionsByAttendanceId(attId)
// 	flagSub := 0
// 	for index, participant := range list {
// 		newList = append(newList, participant)
// 		if len(conditions) > 0 {
// 			for _, condition := range conditions {
// 				if condition.GroupId == nil && participant.GroupId != condition.GroupId {
// 					flagSub += 1
// 				}
// 			}
// 		}
// 		newList[index].ConditionCount = len(conditions) - flagSub
// 		// participant.PassCount = 0

// 		db := global.GVA_DB.Model(&checkins.AttendanceCheckIn{})
// 		var totalRequest int64
// 		db.Where("partpaticipant_id = ? AND attendance_id = ? AND deleted_at IS NULL", participant.ID, attId).Count(&totalRequest)
// 		newList[index].RequestCount = int(totalRequest)

// 		pdb := global.GVA_DB.Model(&checkins.AttendanceCheckIn{})
// 		var totalPass int64
// 		pdb.Where("partpaticipant_id = ? AND attendance_id = ? AND deleted_at IS NULL AND condition_id != 0", participant.ID, attId).Debug().Count(&totalPass)
// 		newList[index].PassCount = int(totalPass)

// 	}
// 	return newList

// }
