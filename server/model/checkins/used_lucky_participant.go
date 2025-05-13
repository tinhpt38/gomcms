package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// UsedLuckyParticipant là cấu trúc để lưu trữ thông tin người tham gia đã được chọn làm người may mắn
type UsedLuckyParticipant struct {
	global.GVA_MODEL
	AttendanceId  string `json:"attendanceId" form:"attendanceId" gorm:"column:attendance_id;comment:ID phiên điểm danh"`
	ParticipantId *uint  `json:"participantId" form:"participantId" gorm:"column:participant_id;comment:ID người tham gia"`
	LuckyNumber   *int   `json:"luckyNumber" form:"luckyNumber" gorm:"column:lucky_number;comment:Số may mắn"`
}

// TableName định nghĩa tên bảng trong cơ sở dữ liệu
func (UsedLuckyParticipant) TableName() string {
	return "used_lucky_participants"
}
