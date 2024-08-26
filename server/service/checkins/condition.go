package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
	checkinsReq "github.com/flipped-aurora/gin-vue-admin/server/model/checkins/request"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ConditionService struct{}

// CreateCondition 创建Điều kiện để checkins记录
// Author [piexlmax](https://github.com/piexlmax)
func (conditionService *ConditionService) CreateCondition(condition *checkins.Condition) (err error) {
	err = global.GVA_DB.Create(condition).Error
	return err
}

// DeleteCondition 删除Điều kiện để checkins记录
// Author [piexlmax](https://github.com/piexlmax)
func (conditionService *ConditionService) DeleteCondition(ID string, userID uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&checkins.Condition{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
			return err
		}
		if err = tx.Delete(&checkins.Condition{}, "id = ?", ID).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteConditionByIds 批量删除Điều kiện để checkins记录
// Author [piexlmax](https://github.com/piexlmax)
func (conditionService *ConditionService) DeleteConditionByIds(IDs []string, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&checkins.Condition{}).Where("id in ?", IDs).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", IDs).Delete(&checkins.Condition{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateCondition 更新Điều kiện để checkins记录
// Author [piexlmax](https://github.com/piexlmax)
func (conditionService *ConditionService) UpdateCondition(condition checkins.Condition) (err error) {
	err = global.GVA_DB.Model(&checkins.Condition{}).Where("id = ?", condition.ID).Updates(&condition).Error
	return err
}

// GetCondition 根据ID获取Điều kiện để checkins记录
// Author [piexlmax](https://github.com/piexlmax)
func (conditionService *ConditionService) GetCondition(ID string) (condition checkins.Condition, err error) {
	err = global.GVA_DB.Where("id = ?", ID).Preload(clause.Associations).Preload("Area.Area").First(&condition).Error
	return
}

func (conditionService *ConditionService) GetConditionsByAttendanceId(ID uint) (list []checkins.Condition, err error) {
	err = global.GVA_DB.Where("attendance_id = ?", ID).Preload("Group").Preload("Area.Area").Find(&list).Error
	return
}

// GetConditionInfoList 分页获取Điều kiện để checkins记录
// Author [piexlmax](https://github.com/piexlmax)
func (conditionService *ConditionService) GetConditionInfoList(info checkinsReq.ConditionSearch) (list []checkins.Condition, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&checkins.Condition{})
	var conditions []checkins.Condition
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.GroupId != nil {
		db = db.Where("group_id = ?", info.GroupId)
	}
	if info.AreaId != nil {
		db = db.Where("area_id = ?", info.AreaId)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Preload(clause.Associations).Preload("Area.Area").Find(&conditions).Error
	return conditions, total, err
}
