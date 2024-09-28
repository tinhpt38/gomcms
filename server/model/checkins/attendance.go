// 自动生成模板Attendance
package checkins

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type Attendance struct {
	global.GVA_MODEL
	Title            string     `json:"title" form:"title" gorm:"column:title;comment:;" binding:"required"`                                         //Tiêu đề
	StartDate        *time.Time `json:"startDate" form:"startDate" gorm:"column:start_date;comment:;" binding:"required"`                            //Ngày bắt đầu
	EndDate          *time.Time `json:"endDate" form:"endDate" gorm:"column:end_date;comment:;" binding:"required"`                                  //Ngày kết thúc
	IsLocked         bool       `json:"isLocked" form:"isLocked" gorm:"column:is_locked;comment:;"`                                                  //Khoá
	EveryoneCanEdit  bool       `json:"everyoneCanEdit" form:"everyoneCanEdit" gorm:"column:everyone_can_edit;comment:Cho phép mọi người chỉnh sửa"` //Cho phép mọi người chỉnh sửa
	AllowGuest       bool       `json:"allowGuest" form:"allowGuest" gorm:"column:allow_guest;comment:Cho phép khách"`                               //Cho phép khách
	Weight           int        `json:"weight" form:"weight" gorm:"column:weight;comment:Hệ số"`
	ClientUrl        string     `json:"clientUrl" form:"clientUrl" gorm:"column:client_url;default=\"https://checkins.dlu.edu.vn\";comment:URL truy cập"` //URL của khách hàng
	RestrictIp       *string    `json:"restrictIp" form:"restrictIp" gorm:"column:restrict_ip;comment:IP giới hạn"`                                       //IPs giới hạn
	LimitCount       int        `json:"limitCount" form:"limitCount" gorm:"column:limit_count;comment:Số lần giới hạn"`
	LimitClientCount int        `json:"limitClientCount" form:"limitClientCount" gorm:"column:limit_client_count;comment:Số lần giới hạn theo máy"`
	RedirectUrl      *string    `json:"redirectUrl" form:"redirectUrl" gorm:"column:redirect_url;comment:URL chuyển hướng"` //URL chuyển hướng
	AgencyId         *uint      `json:"agencyId" form:"agencyId" gorm:"column:agency_id;comment:Đơn vị"`                    //Đơn vị
	CategoryId       *uint      `json:"categoryId" form:"categoryId" gorm:"column:category_id;comment:Loại"`                //Loại
	Description      string     `json:"description" form:"description" gorm:"column:description;comment:Ghi chú"`           //Ghi chú
	CreatedBy        uint       `gorm:"column:created_by;comment:创建者"`
	UpdatedBy        uint       `gorm:"column:updated_by;comment:更新者"`
	DeletedBy        uint       `gorm:"column:deleted_by;comment:删除者"`

	Category     *AttendanceCategory `json:"category" form:"category" gorm:"foreignKey:CategoryId;references:ID"`
	Agency       *AttendanceAgency   `json:"agency" form:"agency" gorm:"foreignKey:AgencyId;references:ID"` //Đơn vị
	Areas        []Area              `gorm:"many2many:attendance_areas;"`                                   // Quan hệ nhiều-nhiều với Area thông qua bảng trung gian AttendanceArea
	Participants []Participant       `gorm:"many2many:participant_attendances;"`                            // Quan hệ nhiều-nhiều với Participant thông qua bảng trung gian AttendanceGroupParticipant
	Total        int                 `json:"total" form:"total"`
	TotalCheckin int                 `json:"totalCheckin" form:"totalCheckin" `
}

func (Attendance) TableName() string {
	return "attendances"
}
