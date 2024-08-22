package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
    checkinsReq "github.com/flipped-aurora/gin-vue-admin/server/model/checkins/request"
)

type GroupService struct {}

// CreateGroup 创建Nhóm记录
// Author [piexlmax](https://github.com/piexlmax)
func (groupService *GroupService) CreateGroup(group *checkins.Group) (err error) {
	err = global.GVA_DB.Create(group).Error
	return err
}

// DeleteGroup 删除Nhóm记录
// Author [piexlmax](https://github.com/piexlmax)
func (groupService *GroupService)DeleteGroup(ID string) (err error) {
	err = global.GVA_DB.Delete(&checkins.Group{},"id = ?",ID).Error
	return err
}

// DeleteGroupByIds 批量删除Nhóm记录
// Author [piexlmax](https://github.com/piexlmax)
func (groupService *GroupService)DeleteGroupByIds(IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]checkins.Group{},"id in ?",IDs).Error
	return err
}

// UpdateGroup 更新Nhóm记录
// Author [piexlmax](https://github.com/piexlmax)
func (groupService *GroupService)UpdateGroup(group checkins.Group) (err error) {
	err = global.GVA_DB.Model(&checkins.Group{}).Where("id = ?",group.ID).Updates(&group).Error
	return err
}

// GetGroup 根据ID获取Nhóm记录
// Author [piexlmax](https://github.com/piexlmax)
func (groupService *GroupService)GetGroup(ID string) (group checkins.Group, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&group).Error
	return
}

// GetGroupInfoList 分页获取Nhóm记录
// Author [piexlmax](https://github.com/piexlmax)
func (groupService *GroupService)GetGroupInfoList(info checkinsReq.GroupSearch) (list []checkins.Group, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&checkins.Group{})
    var groups []checkins.Group
    // 如果有条件搜索 下方会自动创建搜索语句
    if info.StartCreatedAt !=nil && info.EndCreatedAt !=nil {
     db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
    }
    if info.AttendanceClassId != nil {
        db = db.Where("attendance_class_id = ?",info.AttendanceClassId)
    }
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }

	if limit != 0 {
       db = db.Limit(limit).Offset(offset)
    }
	
	err = db.Find(&groups).Error
	return  groups, total, err
}
func (groupService *GroupService)GetGroupDataSource() (res map[string][]map[string]any, err error) {
	res = make(map[string][]map[string]any)
	
	   attendanceClassId := make([]map[string]any, 0)
       global.GVA_DB.Table("attendance_class").Select("id as label,id as value").Scan(&attendanceClassId)
	   res["attendanceClassId"] = attendanceClassId
	return
}