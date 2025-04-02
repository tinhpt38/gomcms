package checkins

import (
	"errors"
	"fmt"
	"os"

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

// ParticipantService chứa các phương thức xử lý nghiệp vụ cho Sinh viên (Participant)
type ParticipantService struct{}

// ------------------- Các hàm CRUD hiện có --------------------

// CreateParticipant 创建Sinh viên (Người tham dự phiên điểm danh)记录
func (participantService *ParticipantService) CreateParticipant(participant *checkins.Participant) (err error) {
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
}

func (participantService *ParticipantService) DeleteParticipant(ID string) (err error) {
	err = global.GVA_DB.Delete(&checkins.Participant{}, "id = ?", ID).Error
	return err
}

func (participantService *ParticipantService) DeleteParticipantInAttendance(ID string, attId uint) (err error) {
	err = global.GVA_DB.Delete(&checkins.AttendanceGroupParticipant{}, "participant_id = ? AND attendance_id = ?", ID, attId).Debug().Error
	return err
}

func (participantService *ParticipantService) DeleteParticipantByIds(IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]checkins.Participant{}, "id in ?", IDs).Error
	return err
}

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
	toDelete := difference(existingGroupIds, newGroupIds)

	if len(toDelete) > 0 {
		if err := agpDb.Where("attendance_id = ? AND group_id IN (?)", participant.AttendanceId, toDelete).
			Unscoped().
			Delete(&checkins.AttendanceGroupParticipant{}).Error; err != nil {
			return err
		}
	}
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
	lookup := make(map[uint]struct{}, len(slice2))
	for _, v := range slice2 {
		lookup[v] = struct{}{}
	}
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

func (participantService *ParticipantService) GetLuckyParticipant(acId string) (participant checkins.Participant, err error) {
	db := global.GVA_DB.Model(&checkins.Participant{})
	err = db.Joins("JOIN attendance_group_participants ON participants.id = attendance_group_participants.participant_id").
		Where("attendance_group_participants.attendance_id = ?", acId).
		Order("RAND()").
		First(&participant).Error

	return
}

func (participantService *ParticipantService) GetParticipantByEmail(email string) (participant checkins.Participant, err error) {
	err = global.GVA_DB.Where("email = ?", email).First(&participant).Error
	return
}

// Sử dụng scope để lấy thông tin các nhóm của một thành viên trong phiên điểm danh
func AttendanceGroupFilters(participantId uint, attendanceId uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Joins("LEFT JOIN `groups` ON groups.id = attendance_group_participants.group_id").
			Where("attendance_group_participants.participant_id = ? AND attendance_group_participants.attendance_id = ? AND groups.attendance_id = ?", participantId, attendanceId, attendanceId).
			Where("attendance_group_participants.deleted_at IS NULL").
			Order("attendance_group_participants.id")
	}
}

func (participantService *ParticipantService) GetParticipantInAttendance(participantId uint, attendanceId uint) (memberOfAttendance []checkins.AttendanceGroupParticipant, err error) {
	db := global.GVA_DB.Table(checkins.AttendanceGroupParticipant{}.TableName()).
		Scopes(AttendanceGroupFilters(participantId, attendanceId))
	err = db.Preload(clause.Associations).Find(&memberOfAttendance).Error
	return
}

// ------------------- Các hàm mới bổ sung --------------------

// ParticipantMetadata mở rộng Participant với các trường tính toán
type ParticipantMetadata struct {
	checkins.Participant
	SuccessCount  int `json:"successCount" gorm:"-"`
	RequiredCount int `json:"requiredCount" gorm:"-"`
}

// FillMetadata tính toán số lần checkin thành công (SuccessCount) và số điều kiện cần điểm danh (RequiredCount)
// Ví dụ: SuccessCount là số bản ghi AttendanceCheckIn có condition_id khác 0,
// RequiredCount có thể là số lần checkin được ghi nhận (hoặc điều kiện áp dụng) – tùy theo nghiệp vụ cụ thể.
func (participantService *ParticipantService) FillMetadata(participants []checkins.Participant, attendanceId uint) ([]ParticipantMetadata, error) {
	var result []ParticipantMetadata
	for _, p := range participants {
		// Đếm số bản ghi checkin (yêu cầu điểm danh) của thành viên tại phiên điểm danh
		var totalRequest int64
		if err := global.GVA_DB.Model(&checkins.AttendanceCheckIn{}).
			Where("partpaticipant_id = ? AND attendance_id = ? AND deleted_at IS NULL", p.ID, attendanceId).
			Count(&totalRequest).Error; err != nil {
			return nil, err
		}
		// Đếm số lần checkin thành công (ví dụ: condition_id khác 0)
		var totalPass int64
		if err := global.GVA_DB.Model(&checkins.AttendanceCheckIn{}).
			Where("partpaticipant_id = ? AND attendance_id = ? AND deleted_at IS NULL AND condition_id != 0", p.ID, attendanceId).
			Count(&totalPass).Error; err != nil {
			return nil, err
		}
		meta := ParticipantMetadata{
			Participant:   p,
			SuccessCount:  int(totalPass),
			RequiredCount: int(totalRequest),
		}
		result = append(result, meta)
	}
	return result, nil
}

// GetParticipantConditions trả về danh sách các điều kiện điểm danh của một thành viên
// Qua join bảng conditions, agp_conditions và attendance_group_participants
func (participantService *ParticipantService) GetParticipantConditions(participantId uint, attendanceId uint) (conditions []checkins.Condition, err error) {
	err = global.GVA_DB.
		Table("conditions").
		Joins("JOIN agp_conditions ON agp_conditions.condition_id = conditions.id").
		Joins("JOIN attendance_group_participants ON agp_conditions.agp_id = attendance_group_participants.id").
		Where("attendance_group_participants.participant_id = ? AND attendance_group_participants.attendance_id = ?", participantId, attendanceId).
		Find(&conditions).Error
	return
}

// ImportExcel và các hàm liên quan dưới đây giữ nguyên như cũ để xử lý nhập liệu từ file Excel
func (participantService *ParticipantService) ImportExcel(info config.CfgFileProcess) (err error) {
	path := global.GVA_CONFIG.Local.Path + "/" + info.UniqueFileName

	file, err := excelize.OpenFile(path)
	if err != nil {
		return
	}

	// Đóng file sau khi đọc
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
		err := os.Remove(path)
		if err != nil {
			global.GVA_LOG.Error("Có lỗi khi xoá tập tin excel: "+info.FileName, zap.Error(err))
		}
	}()

	sheetIndex := 0
	sheet := file.WorkBook.Sheets.Sheet[sheetIndex].Name
	rows, err := file.GetRows(sheet)
	if err != nil {
		return err
	}
	uuid := info.Uuid

	fileProcessService := cfgService.CfgFileProcessServiceApp
	var beginDataRowIndex int = 1
	var lenRows float64 = float64(len(rows) - beginDataRowIndex)

	if lenRows < 1 {
		fileProcessService.UpdateCfgFileProcess(config.CfgFileProcess{
			Uuid:     uuid,
			Percent:  100,
			Msg:      "Tập tin trống, vui lòng bổ sung dữ liệu trước khi nhập",
			Status:   config.FILE_PROCESS_STATUS_ERROR,
			TotalRow: int(lenRows),
		})
	}

	var groups []checkins.Group
	err = global.GVA_DB.Find(&groups).Error
	if err != nil {
		return
	}

	groupMap := make(map[string]checkins.Group)
	for _, group := range groups {
		groupMap[group.Name] = group
	}

	ratio := 100.0 / lenRows
	validationPass := true
	sheetHeaders := rows[0]
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
			return err
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
			fullName := utils.GetArrayValue(row, 0)
			email := utils.GetArrayValue(row, 1)
			group := utils.GetArrayValue(row, 2)
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
		fullName := utils.GetArrayValue(row, 0)
		email := utils.GetArrayValue(row, 1)
		groupName := utils.GetArrayValue(row, 2)
		itemToSave := checkins.Participant{
			GVA_MODEL: global.GVA_MODEL{},
			FullName:  &fullName,
			Email:     email,
		}
		count += 1
		percent := utils.RoundFloat(float64(count)*ratio, 2)
		result := global.GVA_DB.Model(&checkins.Participant{}).Where(&checkins.Participant{
			Email: email,
		}).FirstOrCreate(&itemToSave)
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
			fileProcessService.UpdateCfgFileProcess(config.CfgFileProcess{
				Uuid:    uuid,
				Percent: percent,
				Status:  config.FILE_PROCESS_STATUS_PROCESSING,
			})
		}
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

// Scope cho các bộ lọc tìm kiếm chung đối với Participant
func ParticipantFilters(info checkinsReq.ParticipantSearch) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
			db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
		}
		if info.FullName != "" {
			db = db.Where("full_name LIKE ?", info.FullName+"%")
		}
		if info.Email != "" {
			db = db.Where("email LIKE ?", info.Email+"%")
		}
		if info.AttendanceId != nil {
			db = db.Joins("JOIN attendance_group_participants ON participants.id = attendance_group_participants.participant_id").
				Where("attendance_group_participants.attendance_id = ? AND attendance_group_participants.deleted_at IS NULL", info.AttendanceId)
		}
		if info.GroupId != nil && *info.GroupId != 0 {
			db = db.Joins("JOIN `groups` ON groups.id = attendance_group_participants.group_id").
				Where("groups.id = ?", info.GroupId)
		}
		return db
	}
}

func (participantService *ParticipantService) GetParticipantInfoList(info checkinsReq.ParticipantSearch) (list []checkins.Participant, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&checkins.Participant{}).Scopes(ParticipantFilters(info))
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Preload("Groups", "attendance_id = ?", info.AttendanceId).Find(&list).Error
	return list, total, err
}

func (participantService *ParticipantService) GetParticipantInfoListByAttendance(info checkinsReq.ParticipantSearch) (list []checkins.Participant, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&checkins.Participant{}).Scopes(ParticipantFilters(info))
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Preload("Groups", "attendance_id = ?", info.AttendanceId).Debug().Find(&list).Error
	return list, total, err
}

// GetParticipantConditionData trả về dữ liệu điều kiện của 1 thành viên trong phiên điểm danh
func (participantService *ParticipantService) GetParticipantConditionData(participantId uint, attendanceId uint) (data map[string]interface{}, err error) {
	// Lấy danh sách điều kiện theo participant
	conditions, err := participantService.GetParticipantConditions(participantId, attendanceId)
	if err != nil {
		return nil, err
	}

	// Tính tổng số yêu cầu điểm danh (requiredCount) và số lần checkin thành công (successCount)
	var totalRequest int64
	if err := global.GVA_DB.Model(&checkins.AttendanceCheckIn{}).
		Where("partpaticipant_id = ? AND attendance_id = ? AND deleted_at IS NULL", participantId, attendanceId).
		Count(&totalRequest).Error; err != nil {
		return nil, err
	}

	var totalPass int64
	if err := global.GVA_DB.Model(&checkins.AttendanceCheckIn{}).
		Select("DISTINCT condition_id").
		Where("partpaticipant_id = ? AND attendance_id = ? AND deleted_at IS NULL AND condition_id IS NOT NULL", participantId, attendanceId).
		Count(&totalPass).Error; err != nil {
		return nil, err
	}

	status := "Có dữ liệu"
	if len(conditions) == 0 {
		status = "Không có dữ liệu"
	}

	data = map[string]interface{}{
		"status":        status,
		"successCount":  totalPass,
		"requiredCount": totalRequest,
		"conditions":    conditions,
	}
	return data, nil
}
