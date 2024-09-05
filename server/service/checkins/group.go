package checkins

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
	checkinsReq "github.com/flipped-aurora/gin-vue-admin/server/model/checkins/request"
)

type GroupService struct{}

func (groupService *GroupService) CreateGroup(group *checkins.Group) (err error) {
	err = global.GVA_DB.Create(group).Error
	return err
}

func (groupService *GroupService) DeleteGroup(ID string) (err error) {
	err = global.GVA_DB.Delete(&checkins.Group{}, "id = ?", ID).Error
	return err
}

func (groupService *GroupService) DeleteGroupByIds(IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]checkins.Group{}, "id in ?", IDs).Error
	return err
}

func (groupService *GroupService) UpdateGroup(group checkins.Group) (err error) {
	err = global.GVA_DB.Model(&checkins.Group{}).Where("id = ?", group.ID).Updates(&group).Error
	return err
}

func (groupService *GroupService) GetGroup(ID string) (group checkins.Group, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&group).Error
	return
}

func (groupService *GroupService) GetGroupInfoList(info checkinsReq.GroupSearch) (list []checkins.Group, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.GVA_DB.Model(&checkins.Group{})
	var groups []checkins.Group

	if info.AttendanceId != nil {
		db = db.Where("attendance_id = ?", info.AttendanceId)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 && limit != -1 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&groups).Error

	for i, group := range groups {
		var count int64
		err = global.GVA_DB.Model(&checkins.AttendanceGroupParticipant{}).Where("group_id = ? and attendance_id = ?", group.ID, info.AttendanceId).Count(&count).Error
		if err != nil {
			return
		}
		groups[i].Total = int(count)
	}
	return groups, total, err
}

func (groupService *GroupService) GetGroupDataSource() (res map[string][]map[string]any, err error) {
	res = make(map[string][]map[string]any)

	attendanceClassId := make([]map[string]any, 0)
	global.GVA_DB.Table("attendances").Select("title as label,id as value").Scan(&attendanceClassId)
	res["attendanceClassId"] = attendanceClassId
	return
}

func NewGroupNameGenerator(existingNames []string) *checkins.GroupNameGenerator {
	generator := &checkins.GroupNameGenerator{
		NumericCount:    1,
		AlphabeticCount: 0,
		ExistingNames:   make(map[string]bool),
	}

	if len(existingNames) == 0 {
		return generator
	}

	// Thêm các tên đã tồn tại vào map existingNames
	for _, name := range existingNames {
		generator.ExistingNames[name] = true
	}

	return generator
}

func (groupService *GroupService) GenerateAlphabetic(len int) (result []string) {
	var alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := 0; i < len; i++ {
		var temp string
		num := i
		for {
			temp = string(alphabet[num%26]) + temp
			num = num/26 - 1
			if num < 0 {
				break
			}
		}
		result = append(result, temp)
	}

	return result
}

func (groupService *GroupService) AssignParticipantToGroupAuto(info checkinsReq.GroupAuto) (err error) {

	var groups []checkins.Group
	// Xoá hết các nhóm hiện tại
	err = global.GVA_DB.Delete(&checkins.Group{}, "attendance_id = ?", info.AttendanceId).Error
	if err != nil {
		return
	}
	groupNamePrefix := "Nhóm tự động - "
	alphabetSuggestions := groupService.GenerateAlphabetic(info.GroupQty)
	// Tạo các nhóm mới

	for i := 0; i < info.GroupQty; i++ {
		groupSufix := ""
		if info.GroupNameType == "baseOnNumberic" {
			groupSufix = strconv.Itoa(i + 1)
		} else {
			groupSufix = alphabetSuggestions[i]
		}
		group := checkins.Group{
			Name:         groupNamePrefix + groupSufix,
			AttendanceId: uint(info.AttendanceId),
		}
		err = global.GVA_DB.Create(&group).Error
		if err != nil {
			return
		}
		groups = append(groups, group)
	}

	rDb := global.GVA_DB.Table(checkins.AttendanceGroupParticipant{}.TableName())
	var agps []checkins.AttendanceGroupParticipant
	err = rDb.Where("attendance_id = ?", info.AttendanceId).Find(&agps).Error

	if err != nil {
		return err
	}

	totalParticipants := len(agps)
	groupCount := info.GroupQty

	// Calculate base number of participants per group and the remainder
	baseParticipantsPerGroup := totalParticipants / groupCount
	remainder := totalParticipants % groupCount

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(agps), func(i, j int) {
		agps[i], agps[j] = agps[j], agps[i]
	})

	currentIndex := 0

	for _, group := range groups {
		// Determine number of participants for this group
		numParticipants := baseParticipantsPerGroup
		if remainder > 0 {
			numParticipants++
			remainder--
		}

		// Assign participants to the group
		for i := 0; i < numParticipants; i++ {
			if currentIndex >= totalParticipants {
				break
			}
			agp := agps[currentIndex]
			agp.GroupId = &group.ID
			err = global.GVA_DB.Save(&agp).Error
			if err != nil {
				return err
			}
			currentIndex++
		}
	}

	return nil
}

// func (groupService *GroupService) AssignParticipantToGroupAuto(info checkinsReq.GroupAuto) (err error) {
// 	// Xem thử hiện tại có bao nhiêu nhóm trong DB rồi từ đó tính toán số nhóm cần tạo
// 	var count int64
// 	var groups []checkins.Group

// 	// Lấy danh sách nhóm hiện tại
// 	err = global.GVA_DB.Where("attendance_id = ?", info.AttendanceId).Find(&groups).Error
// 	if err != nil {
// 		return
// 	}

// 	getGroupNames := func(groups []checkins.Group) []string {
// 		names := make([]string, 0)
// 		for _, group := range groups {
// 			names = append(names, group.Name)
// 		}
// 		return names
// 	}

// 	// Đếm số lượng nhóm hiện tại
// 	count = int64(len(groups))

// 	// Tính toán số nhóm cần tạo
// 	// Số nhóm cần tạo = số nhóm cần tạo - số nhóm hiện tại
// 	validGroupQty := int64(info.GroupQty) - count

// 	if validGroupQty > 0 {
// 		// Tiến hành tạo nhóm theo số lượng cần tạo
// 		generator := NewGroupNameGenerator(getGroupNames(groups))

// 		for i := 0; i < int(validGroupQty); i++ {
// 			group := checkins.Group{
// 				Name:         generator.GenerateName(info.GroupNameType),
// 				AttendanceId: uint(info.AttendanceId),
// 			}

// 			err = global.GVA_DB.Create(&group).Error
// 			if err != nil {
// 				return
// 			}
// 		}
// 	}

// 	var participants []checkins.Participant
// 	err = global.GVA_DB.Find(&participants).Error

// 	if err != nil {
// 		return err
// 	}

// 	rand.Seed(time.Now().UnixNano())
// 	rand.Shuffle(len(participants), func(i, j int) {
// 		participants[i], participants[j] = participants[j], participants[i]
// 	})

// 	// Phân bổ participants vào các nhóm và lưu vào bảng GroupParticipant
// 	for i, participant := range participants {
// 		groupID := uint(i%info.GroupQty) + 1

// 		attendanceID := uint(info.AttendanceId)

// 		participantGroup := checkins.AttendanceGroupParticipant{
// 			ParticipantId: &participant.ID,
// 			GroupId:       &groupID,
// 			AttendanceId:  &attendanceID,
// 		}

// 		err = global.GVA_DB.Create(&participantGroup).Error

// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return
// }
