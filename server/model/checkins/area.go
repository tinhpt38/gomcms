// 自动生成模板Area
package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Khu vực điểm danh 结构体  Area
type Area struct {
	global.GVA_MODEL
	Name        string       `json:"name" form:"name" gorm:"column:name;comment:;" binding:"required"` //Tên khu vực
	Latitude    *float64     `json:"latitude" form:"latitude" gorm:"column:latitude;comment:;"`        //Latitude
	Longitude   *float64     `json:"longitude" form:"longitude" gorm:"column:longitude;comment:;"`     //Longitude
	Radius      *float64     `json:"radius" form:"radius" gorm:"column:radius;comment:;"`              //Bán kính
	CreatedBy   uint         `gorm:"column:created_by;comment:创建者"`
	UpdatedBy   uint         `gorm:"column:updated_by;comment:更新者"`
	DeletedBy   uint         `gorm:"column:deleted_by;comment:删除者"`
	Attendances []Attendance `gorm:"many2many:attendance_area;"` // Quan hệ nhiều-nhiều với Attendance thông qua bảng trung gian AttendanceArea

}

// TableName Khu vực điểm danh Area自定义表名 areas
func (Area) TableName() string {
	return "areas"
}

type AttendanceArea struct {
	global.GVA_MODEL
	AreaID       uint       `json:"area_id" form:"area_id" gorm:"column:area_id;comment:Khu vực"`
	AttendanceID uint       `json:"attendance_id" form:"attendance_id" gorm:"column:attendance_id;comment:Lớp điểm danh"`
	Area         Area       `gorm:"foreignKey:AreaID"`                                    // Thiết lập khóa ngoại đến struct Area
	Attendance   Attendance `gorm:"foreignKey:AttendanceID"`                              // Thiết lập khóa ngoại đến struct AttendanceClass
	Radius       *float64   `json:"radius" form:"radius" gorm:"column:radius;comment:半径"` // Bán kính
}

func (AttendanceArea) TableName() string {
	return "attendance_areas"
}
