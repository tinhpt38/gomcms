// 自动生成模板AttendanceCheckIn
package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"time"
	"gorm.io/datatypes"
)

// Thành viên checkins 结构体  AttendanceCheckIn
type AttendanceCheckIn struct {
    global.GVA_MODEL
    CheckinDate  *time.Time `json:"checkinDate" form:"checkinDate" gorm:"column:checkin_date;comment:;" binding:"required"`  //Ngày giờ 
    AttendanceId  *int `json:"attendanceId" form:"attendanceId" gorm:"column:attendance_id;comment:;" binding:"required"`  //Lớp điểm danh 
    PartpaticipantId  *int `json:"partpaticipantId" form:"partpaticipantId" gorm:"column:partpaticipant_id;comment:;" binding:"required"`  //Thành viên 
    AreaId  *int `json:"areaId" form:"areaId" gorm:"column:area_id;comment:;"`  //Khu vực 
    GroupId  *int `json:"groupId" form:"groupId" gorm:"column:group_id;comment:;"`  //Nhóm 
    ConditionId  *int `json:"conditionId" form:"conditionId" gorm:"column:condition_id;comment:;"`  //Điều kiện 
    Lattidue  *float64 `json:"lattidue" form:"lattidue" gorm:"column:lattidue;comment:;"`  //Kinh độ 
    IP  *int `json:"iP" form:"iP" gorm:"column:ip;comment:;"`  //IP 
    Longtidue  *float64 `json:"longtidue" form:"longtidue" gorm:"column:longtidue;comment:;"`  //Vĩ độ 
    Agent  datatypes.JSON `json:"agent" form:"agent" gorm:"column:agent;comment:;type:text;"swaggertype:"object"`  //Agent 
    Browser  string `json:"browser" form:"browser" gorm:"column:browser;comment:;"`  //Trình duyệt 
    CreatedBy  uint   `gorm:"column:created_by;comment:创建者"`
    UpdatedBy  uint   `gorm:"column:updated_by;comment:更新者"`
    DeletedBy  uint   `gorm:"column:deleted_by;comment:删除者"`
}


// TableName Thành viên checkins AttendanceCheckIn自定义表名 attendanceCheckIns
func (AttendanceCheckIn) TableName() string {
    return "attendanceCheckIns"
}

