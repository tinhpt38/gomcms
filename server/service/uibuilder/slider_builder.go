package uibuilder

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/uibuilder"
    uibuilderReq "github.com/flipped-aurora/gin-vue-admin/server/model/uibuilder/request"
    "gorm.io/gorm"
)

type SliderBuilderService struct {}

// CreateSliderBuilder 创建SliderBuilder记录
// Author [piexlmax](https://github.com/piexlmax)
func (sliderBuilderService *SliderBuilderService) CreateSliderBuilder(sliderBuilder *uibuilder.SliderBuilder) (err error) {
	err = global.GVA_DB.Create(sliderBuilder).Error
	return err
}

// DeleteSliderBuilder 删除SliderBuilder记录
// Author [piexlmax](https://github.com/piexlmax)
func (sliderBuilderService *SliderBuilderService)DeleteSliderBuilder(ID string,userID uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
	    if err := tx.Model(&uibuilder.SliderBuilder{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
              return err
        }
        if err = tx.Delete(&uibuilder.SliderBuilder{},"id = ?",ID).Error; err != nil {
              return err
        }
        return nil
	})
	return err
}

// DeleteSliderBuilderByIds 批量删除SliderBuilder记录
// Author [piexlmax](https://github.com/piexlmax)
func (sliderBuilderService *SliderBuilderService)DeleteSliderBuilderByIds(IDs []string,deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
	    if err := tx.Model(&uibuilder.SliderBuilder{}).Where("id in ?", IDs).Update("deleted_by", deleted_by).Error; err != nil {
            return err
        }
        if err := tx.Where("id in ?", IDs).Delete(&uibuilder.SliderBuilder{}).Error; err != nil {
            return err
        }
        return nil
    })
	return err
}

// UpdateSliderBuilder 更新SliderBuilder记录
// Author [piexlmax](https://github.com/piexlmax)
func (sliderBuilderService *SliderBuilderService)UpdateSliderBuilder(sliderBuilder uibuilder.SliderBuilder) (err error) {
	err = global.GVA_DB.Model(&uibuilder.SliderBuilder{}).Where("id = ?",sliderBuilder.ID).Updates(&sliderBuilder).Error
	return err
}

// GetSliderBuilder 根据ID获取SliderBuilder记录
// Author [piexlmax](https://github.com/piexlmax)
func (sliderBuilderService *SliderBuilderService)GetSliderBuilder(ID string) (sliderBuilder uibuilder.SliderBuilder, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&sliderBuilder).Error
	return
}

// GetSliderBuilderInfoList 分页获取SliderBuilder记录
// Author [piexlmax](https://github.com/piexlmax)
func (sliderBuilderService *SliderBuilderService)GetSliderBuilderInfoList(info uibuilderReq.SliderBuilderSearch) (list []uibuilder.SliderBuilder, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&uibuilder.SliderBuilder{})
    var sliderBuilders []uibuilder.SliderBuilder
    // 如果有条件搜索 下方会自动创建搜索语句
    if info.StartCreatedAt !=nil && info.EndCreatedAt !=nil {
     db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
    }
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }

	if limit != 0 {
       db = db.Limit(limit).Offset(offset)
    }
	
	err = db.Find(&sliderBuilders).Error
	return  sliderBuilders, total, err
}