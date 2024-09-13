package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
	checkinsReq "github.com/flipped-aurora/gin-vue-admin/server/model/checkins/request"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AttendanceService struct{}

func (attendanceService *AttendanceService) CreateAttendance(attendanceClass *checkins.Attendance) (err error) {
	err = global.GVA_DB.Create(attendanceClass).Error
	return err
}

func (attendanceService *AttendanceService) CreateAttendanceArea(attendanceArea *checkins.AttendanceArea) (err error) {
	err = global.GVA_DB.Create(attendanceArea).Error
	return err
}

func (attendanceService *AttendanceService) UpdateAttendanceArea(attendanceArea *checkins.AttendanceArea) (err error) {
	err = global.GVA_DB.Model(&checkins.AttendanceArea{}).Where("id = ?", attendanceArea.ID).Save(&attendanceArea).Error
	return err
}

func (attendanceService *AttendanceService) DeleteAttendanceArea(id string) (err error) {
	// err = global.GVA_DB.Delete(&checkins.AttendanceArea{}, "id = ? ", id).Error
	// Xoá luôn tham chiếu tới điều kiện
	attArea := checkins.AttendanceArea{}

	err = global.GVA_DB.Where("id = ?", id).First(&attArea).Error

	if err != nil {
		return err
	}

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("area_id = ?", attArea.AreaID).Delete(&checkins.Condition{}).Error; err != nil {
			return err
		}

		err = tx.Delete(&checkins.AttendanceArea{}, "id = ?", id).Error
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (attendanceService *AttendanceService) DeleteAttendance(ID string, userID uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&checkins.Attendance{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
			return err
		}
		if err = tx.Delete(&checkins.Attendance{}, "id = ?", ID).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (attendanceService *AttendanceService) DeleteAttendanceByIds(IDs []string, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&checkins.Attendance{}).Where("id in ?", IDs).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", IDs).Delete(&checkins.Attendance{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (attendanceService *AttendanceService) UpdateAttendance(attendanceClass checkins.Attendance) (err error) {
	err = global.GVA_DB.Model(&checkins.Attendance{}).Where("id = ?", attendanceClass.ID).Save(&attendanceClass).Error
	return err
}

func (attendanceService *AttendanceService) GetAttendance(ID string) (attendanceClass checkins.Attendance, err error) {
	err = global.GVA_DB.Where("id = ?", ID).Preload(clause.Associations).First(&attendanceClass).Error
	if err != nil {
		return
	}
	var totalParticipants, totalCheckin int64

	// Tính tổng số thành viên từ bảng attendance_group_participant
	err = global.GVA_DB.Table("attendance_group_participants").
		Where("attendance_id = ?", attendanceClass.ID).
		Count(&totalParticipants).Error
	if err != nil {
		return
	}

	// Tính tổng số thành viên đã check-in từ bảng attendance_check_in
	err = global.GVA_DB.Table("attendance_checkins").
		Select("count(DISTINCT partpaticipant_id) as total").
		Where("deleted_at is NULL").
		Where("attendance_id = ?", attendanceClass.ID).
		Scan(&totalCheckin).Error
	if err != nil {
		return
	}

	// Cập nhật vào model Attendance
	attendanceClass.Total = int(totalParticipants)
	attendanceClass.TotalCheckin = int(totalCheckin)

	return
}

func (attendanceService *AttendanceService) GetAttendanceArea(id string) (list []checkins.AttendanceArea, err error) {

	db := global.GVA_DB.Table(checkins.AttendanceArea{}.TableName())
	err = db.Where("attendance_id = ?", id).Preload("Area").Find(&list).Error
	return
}

func (attendanceService *AttendanceService) GetAttendanceInfoList(info checkinsReq.AttendanceSearch, createdBy uint) (list []checkins.Attendance, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.GVA_DB.Model(&checkins.Attendance{})
	var attendanceClasss []checkins.Attendance

	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}

	if info.StartDate != nil && info.EndDate != nil {
		db = db.Where("DATE(start_date) >= ? OR DATE(end_date) <= ?", info.StartDate.Format("2006-01-02"), info.EndDate.Format("2006-01-02"))
	}

	if createdBy != 0 && createdBy != 1 {
		db = db.Where("created_by = ?", createdBy)
	}
	if info.AgencyId != 0 {
		db = db.Where("agency_id = ?", info.AgencyId)
	}
	if info.CategoryId != 0 {
		db = db.Where("category_id = ?", info.CategoryId)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 && limit != -1 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Preload("Areas").Preload(clause.Associations).Order("created_at DESC").Find(&attendanceClasss).Error
	for i, attendance := range attendanceClasss {
		var totalParticipants, totalCheckin int64

		// Tính tổng số thành viên từ bảng attendance_group_participant
		err = global.GVA_DB.Table("attendance_group_participants").
			Where("attendance_id = ?", attendance.ID).
			Count(&totalParticipants).Error
		if err != nil {
			return
		}

		// Tính tổng số thành viên đã check-in từ bảng attendance_check_in
		err = global.GVA_DB.Table("attendance_checkins").
			Select("count(DISTINCT partpaticipant_id) as total").
			Where("deleted_at is NULL").
			Where("attendance_id = ?", attendance.ID).
			Scan(&totalCheckin).Error
		if err != nil {
			return
		}

		// Cập nhật vào model Attendance
		attendanceClasss[i].Total = int(totalParticipants)
		attendanceClasss[i].TotalCheckin = int(totalCheckin)
	}
	return attendanceClasss, total, err
}

func (attendanceService *AttendanceService) GetPublicAttendanceInfoList(info checkinsReq.AttendanceSearch, createdBy uint) (list []checkins.Attendance, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.GVA_DB.Model(&checkins.Attendance{})
	var attendanceClasss []checkins.Attendance

	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}

	if info.StartDate != nil && info.EndDate != nil {
		db = db.Where("DATE(start_date) >= ? OR DATE(end_date) <= ?", info.StartDate.Format("2006-01-02"), info.EndDate.Format("2006-01-02"))
	}

	if createdBy != 0 && createdBy != 1 {
		db = db.Where("created_by = ?", createdBy)
	}
	if info.AgencyId != 0 {
		db = db.Where("agency_id = ?", info.AgencyId)
	}
	if info.CategoryId != 0 {
		db = db.Where("category_id = ?", info.CategoryId)
	}

	db = db.Where("is_locked = 0")

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 && limit != -1 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Order("created_at DESC").Preload(clause.Associations).Find(&attendanceClasss).Error
	for i, attendance := range attendanceClasss {
		var totalParticipants, totalCheckin int64

		// Tính tổng số thành viên từ bảng attendance_group_participant
		err = global.GVA_DB.Table("attendance_group_participants").
			Where("attendance_id = ?", attendance.ID).
			Count(&totalParticipants).Error
		if err != nil {
			return
		}

		// Tính tổng số thành viên đã check-in từ bảng attendance_check_in
		err = global.GVA_DB.Table("attendance_checkins").
			Select("count(DISTINCT partpaticipant_id) as total").
			Where("deleted_at is NULL").
			Where("attendance_id = ?", attendance.ID).
			Scan(&totalCheckin).Error
		if err != nil {
			return
		}

		// Cập nhật vào model Attendance
		attendanceClasss[i].Total = int(totalParticipants)
		attendanceClasss[i].TotalCheckin = int(totalCheckin)
	}

	return attendanceClasss, total, err
}

func (attendanceService *AttendanceService) CloneAttendance(req checkinsReq.AttendanceCloneRequest, createdBy uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// Retrieve the current attendance
		curAttendance := checkins.Attendance{}
		err = tx.First(&curAttendance, req.AttendanceID).Error
		if err != nil {
			return err
		}
		newAtendance := curAttendance
		newAtendance.CreatedBy = createdBy
		newAtendance.ID = 0
		newAtendance.Title = newAtendance.Title + " - Copy"
		if err := tx.Create(&newAtendance).Error; err != nil {
			return err
		}

		curAttendanceAreas := []checkins.AttendanceArea{}
		if err := tx.Where("attendance_id = ?", curAttendance.ID).Find(&curAttendanceAreas).Error; err != nil {
			return err
		}

		for _, area := range curAttendanceAreas {
			newArea := area
			newArea.ID = 0
			newArea.AttendanceID = &newAtendance.ID
			if err := tx.Create(&newArea).Error; err != nil {
				return err
			}
		}

		curConditions := []checkins.Condition{}
		if err := tx.Where("attendance_id = ?", curAttendance.ID).Find(&curConditions).Error; err != nil {
			return err
		}
		for _, condition := range curConditions {
			newCondition := condition
			newCondition.ID = 0
			newCondition.AttendanceId = &newAtendance.ID
			if err := tx.Create(&newCondition).Error; err != nil {
				return err
			}
		}

		if req.WithData {

			curGroups := []checkins.Group{}
			if err := tx.Where("attendance_id = ?", curAttendance.ID).Find(&curGroups).Error; err != nil {
				return err
			}
			newGroups := []checkins.Group{}
			for _, group := range curGroups {
				newGroup := group
				newGroup.ID = 0
				newGroup.AttendanceId = newAtendance.ID
				if err := tx.Create(&newGroup).Error; err != nil {
					return err
				}
				newGroups = append(newGroups, newGroup)
			}

			agps := []checkins.AttendanceGroupParticipant{}
			if err := tx.Where("attendance_id = ?", curAttendance.ID).Find(&agps).Error; err != nil {
				return err
			}
			for _, agp := range agps {
				newAgp := agp
				newAgp.ID = 0
				newAgp.AttendanceId = &newAtendance.ID
				for index, group := range curGroups {
					if group.ID == *agp.GroupId {
						newAgp.GroupId = &newGroups[index].ID
						break
					}
				}
				if err := tx.Create(&newAgp).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	return
}

func (attendanceService *AttendanceService) StatsByAgencyCategory(req checkinsReq.StatsByInfoRequest) (result map[string]interface{}, err error) {

	category := checkins.AttendanceCategory{}
	err = global.GVA_DB.Where("id = ?", req.CategoryId).First(&category).Error
	if err != nil {
		return
	}
	var listCatIs []uint
	if *category.ParentId == 0 {
		err = global.GVA_DB.Model(&checkins.AttendanceCategory{}).Select("id").Where("id = ?", category.ParentId).Scan(&listCatIs).Error
		if err != nil {
			return
		}
	} else {
		listCatIs = append(listCatIs, req.CategoryId)
	}

	var resultRows []struct {
		AgencyName    string
		CountByAgency int
		TotalParts    int
	}

	rawDB := global.GVA_DB.Table("attendances as att").
		Select("attagen.name as AgencyName, count(att.agency_id) as CountByAgency, acheckins.totalParts as TotalParts").
		Joins("LEFT JOIN attendance_agencies as attagen on attagen.id = att.agency_id").
		Joins("LEFT JOIN attendance_categories as attcat on attcat.id = att.category_id").
		Joins("LEFT JOIN (SELECT attendance_checkins.attendance_id as attid, COUNT(DISTINCT attendance_checkins.partpaticipant_id) as totalParts FROM attendance_checkins WHERE attendance_checkins.deleted_at is NULL GROUP by attid) as acheckins on acheckins.attid = att.id").
		Where("att.deleted_at is NULL")
	if req.AgencyId != 0 {
		rawDB.Where("att.agency_id = ?", req.AgencyId)
	}
	rawDB.Where("att.category_id IN (?)", listCatIs).Group("att.agency_id")
	err = rawDB.Debug().Find(&resultRows).Error
	if err != nil {
		return
	}
	result = make(map[string]interface{})
	result["data"] = resultRows
	return result, err
}

func (attendanceService *AttendanceService) StatsScatterPlot(req checkinsReq.StatsByInfoRequest) (result map[string]interface{}, err error) {
	// Xác định category
	category := checkins.AttendanceCategory{}
	err = global.GVA_DB.Where("id = ?", req.CategoryId).First(&category).Error
	if err != nil {
		return
	}

	var listCatIds []uint
	// Kiểm tra xem category có parent hay không
	if *category.ParentId == 0 {
		// Lấy tất cả các category con nếu có
		err = global.GVA_DB.Model(&checkins.AttendanceCategory{}).Select("id").Where("parent_id = ?", category.ID).Scan(&listCatIds).Error
		if err != nil {
			return
		}
	} else {
		listCatIds = append(listCatIds, req.CategoryId)
	}

	// Cấu trúc để nhận kết quả
	var resultRows []struct {
		CheckinTime  string
		CheckinCount int
	}

	// Xây dựng truy vấn thô
	rawDB := global.GVA_DB.Table("attendance_checkins as attci").
		Select("DATE_FORMAT(attci.created_at, '%d/%m/%Y') as CheckinTime, COUNT(attci.id) as CheckinCount").
		Joins("LEFT JOIN attendances as att on att.id = attci.attendance_id").
		Where("attci.deleted_at IS NULL").
		Where("att.category_id IN (?)", listCatIds).
		Group("DATE_FORMAT(attci.created_at, '%d/%m/%Y')").
		Order("CheckinTime")

	// Lọc theo agency nếu có
	// if req.AgencyId != 0 {
	// 	rawDB = rawDB.Where("att.agency_id = ?", req.AgencyId)
	// }

	// Lọc theo khoảng thời gian (nếu có trong request)
	if req.StartAt != nil && req.EndAt != nil {
		rawDB = rawDB.Where("attci.created_at BETWEEN ? AND ?", req.StartAt, req.EndAt)
	}

	// Thực thi truy vấn và lưu kết quả
	err = rawDB.Debug().Find(&resultRows).Error
	if err != nil {
		return
	}

	// Định dạng kết quả trả về
	result = make(map[string]interface{})
	result["data"] = resultRows
	return result, err
}

func (attendanceService *AttendanceService) StatsTrendLine(req checkinsReq.StatsByInfoRequest) (result map[string]interface{}, err error) {
	// Xác định category
	category := checkins.AttendanceCategory{}
	err = global.GVA_DB.Where("id = ?", req.CategoryId).First(&category).Error
	if err != nil {
		return
	}

	var listCatIds []uint
	// Kiểm tra xem category có parent hay không
	if *category.ParentId == 0 {
		// Lấy tất cả các category con nếu có
		err = global.GVA_DB.Model(&checkins.AttendanceCategory{}).Select("id").Where("parent_id = ?", category.ID).Scan(&listCatIds).Error
		if err != nil {
			return
		}
	} else {
		listCatIds = append(listCatIds, req.CategoryId)
	}

	// Cấu trúc để nhận kết quả
	var resultRows []struct {
		CheckinDate  string
		CheckinCount int
		AgencyName   string
		CategoryName string
	}
	rawDB := global.GVA_DB.Table("attendance_checkins as ac").
		Select("DATE(ac.created_at) as CheckinDate, COUNT(ac.partpaticipant_id) as CheckinCount, attagen.name as AgencyName, attcat.name as CategoryName").
		Joins("LEFT JOIN attendances att ON att.id = ac.attendance_id").
		Joins("LEFT JOIN attendance_agencies attagen ON attagen.id = att.agency_id").
		Joins("LEFT JOIN attendance_categories attcat ON attcat.id = att.category_id").
		Where("ac.deleted_at IS NULL").
		Where("att.deleted_at IS NULL").
		Where("attagen.deleted_at IS NULL").
		Where("attcat.deleted_at IS NULL").
		Where("att.category_id IN (?)", listCatIds).
		Group("DATE(ac.created_at), attagen.name, attcat.name").
		Order("CheckinDate")

	// Lọc theo agency nếu có
	if req.AgencyId != 0 {
		rawDB = rawDB.Where("att.agency_id = ?", req.AgencyId)
	}

	// Lọc theo khoảng thời gian (nếu có trong request)
	if req.StartAt != nil && req.EndAt != nil {
		rawDB = rawDB.Where("ac.created_at BETWEEN ? AND ?", req.StartAt, req.EndAt)
	}

	// Thực thi truy vấn và lưu kết quả
	err = rawDB.Debug().Find(&resultRows).Error
	if err != nil {
		return
	}

	// Định dạng kết quả trả về
	result = make(map[string]interface{})
	result["data"] = resultRows
	return result, err
}
