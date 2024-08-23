// 自动生成模板FileProcessError
package config

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Lỗi ghi nhận từ quá trình xử lý tệp tin 结构体  FileProcessError
type FileProcessError struct {
	global.GVA_MODEL
	FileProcessId   *uint  `json:"fileProcessId" form:"fileProcessId" gorm:"column:file_process_id;comment:;" binding:"required"`       //Mã định danh quá trình xử lý tệp tin
	FileProcessUuid string `json:"fileProcessUuid" form:"fileProcessUuid" gorm:"column:file_process_uuid;comment:;" binding:"required"` //Mã định danh của tệp tin
	FieldTitle      string `json:"fieldTitle" form:"fieldTitle" gorm:"column:field_title;comment:;" binding:"required"`                 //Cột dữ liệu
	ExpectedValue   string `json:"expectedValue" form:"expectedValue" gorm:"column:expected_value;comment:;size:191;"`
	ReceivedValue   string `json:"receivedValue" form:"receivedValue" gorm:"column:received_value;comment:;size:191;"`
	Note            string `json:"note" form:"note" gorm:"column:note;comment:;size:191;"`
	Row             int    `json:"row" form:"row" gorm:"column:row;comment:;size:19;"`
	CreatedBy       uint   `gorm:"column:created_by;comment:Creator"`
	UpdatedBy       uint   `gorm:"column:updated_by;comment:Updater"`
	DeletedBy       uint   `gorm:"column:deleted_by;comment:Deleter"`
}

// TableName Lỗi ghi nhận từ quá trình xử lý tệp tin FileProcessError自定义表名 cfg_file_process_error
func (FileProcessError) TableName() string {
	return "cfg_file_process_error"
}
