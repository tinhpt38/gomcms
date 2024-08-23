package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
	checkinsReq "github.com/flipped-aurora/gin-vue-admin/server/model/checkins/request"
	"gorm.io/gorm"
)

type AreaService struct{}

func (areaService *AreaService) CreateArea(area *checkins.Area) (err error) {
	err = global.GVA_DB.Create(area).Error
	return err
}

func (areaService *AreaService) DeleteArea(ID string, userID uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&checkins.Area{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
			return err
		}
		if err = tx.Delete(&checkins.Area{}, "id = ?", ID).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (areaService *AreaService) DeleteAreaByIds(IDs []string, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&checkins.Area{}).Where("id in ?", IDs).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", IDs).Delete(&checkins.Area{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (areaService *AreaService) UpdateArea(area checkins.Area) (err error) {
	err = global.GVA_DB.Model(&checkins.Area{}).Where("id = ?", area.ID).Updates(&area).Error
	return err
}

func (areaService *AreaService) GetArea(ID string) (area checkins.Area, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&area).Error
	return
}

func (areaService *AreaService) GetAreaInfoList(info checkinsReq.AreaSearch) (list []checkins.Area, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&checkins.Area{})
	var areas []checkins.Area

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

	err = db.Find(&areas).Error
	return areas, total, err
}
