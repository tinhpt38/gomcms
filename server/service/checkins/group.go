package checkins

import (
	"math/rand"
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

	// if info.StartCreatedAt !=nil && info.EndCreatedAt !=nil {
	//  db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	// }
	if info.AttendanceId != nil {
		db = db.Where("attendance_id = ?", info.AttendanceId)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Debug().Find(&groups).Error
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

func (groupService *GroupService) AssignParticipantToGroupAuto(info checkinsReq.GroupAuto) (err error) {
	// Xem thử hiện tại có bao nhiêu nhóm trong DB rồi từ đó tính toán số nhóm cần tạo
	var count int64
	var groups []checkins.Group

	// Lấy danh sách nhóm hiện tại
	err = global.GVA_DB.Where("attendance_id = ?", info.AttendanceId).Find(&groups).Error
	if err != nil {
		return
	}

	getGroupNames := func(groups []checkins.Group) []string {
		names := make([]string, 0)
		for _, group := range groups {
			names = append(names, group.Name)
		}
		return names
	}

	// Đếm số lượng nhóm hiện tại
	count = int64(len(groups))

	// Tính toán số nhóm cần tạo
	// Số nhóm cần tạo = số nhóm cần tạo - số nhóm hiện tại
	validGroupQty := int64(info.GroupQty) - count

	if validGroupQty > 0 {
		// Tiến hành tạo nhóm theo số lượng cần tạo
		generator := NewGroupNameGenerator(getGroupNames(groups))

		for i := 0; i < int(validGroupQty); i++ {
			group := checkins.Group{
				Name:         generator.GenerateName(info.GroupNameType),
				AttendanceId: &info.AttendanceId,
			}

			err = global.GVA_DB.Create(&group).Error
			if err != nil {
				return
			}
		}
	}

	var participants []checkins.Participant
	err = global.GVA_DB.Find(&participants).Error

	if err != nil {
		return err
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(participants), func(i, j int) {
		participants[i], participants[j] = participants[j], participants[i]
	})

	// Phân bổ participants vào các nhóm và lưu vào bảng GroupParticipant
	for i, participant := range participants {
		groupID := uint(i%info.GroupQty) + 1

		attendanceID := uint(info.AttendanceId)

		participantGroup := checkins.AttendanceGroupParticipant{
			ParticipantId: &participant.ID,
			GroupId:       &groupID,
			AttendanceId:  &attendanceID,
		}

		err = global.GVA_DB.Create(&participantGroup).Error

		if err != nil {
			return err
		}
	}

	return
}
