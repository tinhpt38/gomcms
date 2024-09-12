// 自动生成模板SliderBuilder
package uibuilder

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/datatypes"
)

// SliderBuilder 结构体  SliderBuilder
type SliderBuilder struct {
	global.GVA_MODEL
	Images    datatypes.JSON `json:"images" form:"images" gorm:"column:images;comment:;" binding:"required"swaggertype:"array,object"`   //Hỉnh ảnh
	Title     string         `json:"title" form:"title" gorm:"column:title;comment:Tiêu đề;type:varchar(255);" binding:"required"`       //Tiêu đề
	SubTitle  string         `json:"subTitle" form:"subTitle" gorm:"column:sub_title;comment:Tiêu đề phụ;type:varchar(255);"`            //Tiêu đề phụ
	Content   string         `json:"content" form:"content" gorm:"column:content;comment:Nội dung;type:longtext;"`                       //Nội dung
	Show      bool           `json:"show" form:"show" gorm:"column:show;comment:Hiển thị;type:tinyint(1);default:1;" binding:"required"` //Hiển thị
	Link      string         `json:"link" form:"link" gorm:"column:link;comment:Liên kết;type:varchar(255);"`                            //Liên kết
	CreatedBy uint           `gorm:"column:created_by;comment:创建者"`
	UpdatedBy uint           `gorm:"column:updated_by;comment:更新者"`
	DeletedBy uint           `gorm:"column:deleted_by;comment:删除者"`
}

// TableName SliderBuilder SliderBuilder自定义表名 uibuiler_sliders
func (SliderBuilder) TableName() string {
	return "uibuiler_sliders"
}
