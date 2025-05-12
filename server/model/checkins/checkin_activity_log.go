package checkins

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// CheckinActivityLog lưu lại các hoạt động điểm danh
type CheckinActivityLog struct {
	global.GVA_MODEL
	AttendanceID  uint      `json:"attendance_id" gorm:"column:attendance_id;comment:ID của buổi điểm danh"`
	ParticipantID uint      `json:"participant_id" gorm:"column:participant_id;comment:ID của người tham gia"`
	Action        string    `json:"action" gorm:"column:action;type:varchar(100);comment:Hành động thực hiện"`
	Level         string    `json:"level" gorm:"column:level;type:varchar(50);comment:Mức độ quan trọng"`
	Message       string    `json:"message" gorm:"column:message;type:text;comment:Nội dung log"`
	Timestamp     time.Time `json:"timestamp" gorm:"column:timestamp;comment:Thời gian xảy ra"`
}
