// 自动生成模板CfgFileProcess
package config

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

// Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống 结构体  CfgFileProcess
type CfgFileProcess struct {
	global.GVA_MODEL
	Uuid           string  `json:"uuid" form:"uuid" gorm:"column:uuid;comment:;" binding:"required"` //Mã định danh tệp tin
	Percent        float64 `json:"percent" form:"percent" gorm:"column:percent;comment:;size:10;"`   //Tiến trình phần trăm
	Msg            string  `json:"msg" form:"msg" gorm:"column:msg;comment:;size:255;"`              //Thông báo
	Status         string  `json:"status" form:"status" gorm:"column:status;comment:;size:10;"`
	FileName       string  `json:"fileName" form:"fileName" gorm:"column:file_name;comment:;size:255;"`
	UniqueFileName string  `json:"uniqueFileName" form:"unique_file_name" gorm:"column:unique_file_name;comment:;size:255;"`
	Extension      string  `json:"extension" form:"extension" gorm:"column:extension;comment:;size:255;"`
	FileSize       float64 `json:"fileSize" form:"fileSize" gorm:"column:file_size;comment:;size:255;"`
	UserID         uint    `json:"userId" form:"user_id" gorm:"column:user_id;comment:;"`
	TotalRow       int     `json:"totalRow" form:"total_row" gorm:"column:total_row;comment:;"`
	Action         string  `json:"action" form:"action" gorm:"column:action;comment:;"`
	Type           string  `json:"type" form:"type" gorm:"column:type;comment:;size:255;"`
	Url            string  `json:"url" form:"url" gorm:"column:url;comment:;"`
	CreatedBy      uint    `gorm:"column:created_by;comment:Creator"`
	UpdatedBy      uint    `gorm:"column:updated_by;comment:Updater"`
	DeletedBy      uint    `gorm:"column:deleted_by;comment:Deleter"`

	User *system.SysUser `json:"user" form:"user" gorm:"foreignKey:CreatedBy"`
}

// TableName Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống CfgFileProcess自定义表名 cfg_file_process
func (CfgFileProcess) TableName() string {
	return "cfg_file_process"
}

// Tạo ra các hằng số cho hành động action
var (
	// ActionUpload Tải lên
	ACTION_IMPORT_SAMPLE = "IMPORT_SAMPLE"

	FILE_PROCESS_STATUS_PENDING    = "pending"
	FILE_PROCESS_STATUS_VALIDATING = "validating"
	FILE_PROCESS_STATUS_PROCESSING = "processing"
	FILE_PROCESS_STATUS_FINISH     = "finish"
	FILE_PROCESS_STATUS_ERROR      = "error"
)

func CheckFileProcessAction(action string) bool {
	valids := make([]string, 0)

	valids = append(valids, ACTION_IMPORT_SAMPLE)

	// Kiểm tra
	for _, valid := range valids {
		if valid == action {
			return true
		}
	}

	return false
}
