package checkins

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
	checkinsReq "github.com/flipped-aurora/gin-vue-admin/server/model/checkins/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GroupApi struct{}

func (groupApi *GroupApi) CreateGroup(c *gin.Context) {
	var group checkins.Group
	err := c.ShouldBindJSON(&group)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = groupService.CreateGroup(&group)
	if err != nil {
		global.GVA_LOG.Error("tạo thất bại!", zap.Error(err))
		response.FailWithMessage("tạo thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("tạo thành công", c)
}

func (groupApi *GroupApi) DeleteGroup(c *gin.Context) {
	ID := c.Query("ID")
	err := groupService.DeleteGroup(ID)
	if err != nil {
		global.GVA_LOG.Error("xóa thất bại!", zap.Error(err))
		response.FailWithMessage("xóa thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("xóa thành công", c)
}

func (groupApi *GroupApi) DeleteGroupByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	err := groupService.DeleteGroupByIds(IDs)
	if err != nil {
		global.GVA_LOG.Error("xóa nhiều thất bại!", zap.Error(err))
		response.FailWithMessage("xóa nhiều thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("xóa nhiều thành công", c)
}

func (groupApi *GroupApi) UpdateGroup(c *gin.Context) {
	var group checkins.Group
	err := c.ShouldBindJSON(&group)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = groupService.UpdateGroup(group)
	if err != nil {
		global.GVA_LOG.Error("cập nhật thất bại!", zap.Error(err))
		response.FailWithMessage("cập nhật thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("cập nhật thành công", c)
}

func (groupApi *GroupApi) FindGroup(c *gin.Context) {
	ID := c.Query("ID")
	regroup, err := groupService.GetGroup(ID)
	if err != nil {
		global.GVA_LOG.Error("tìm thất bại!", zap.Error(err))
		response.FailWithMessage("tìm thất bại:"+err.Error(), c)
		return
	}
	response.OkWithData(regroup, c)
}

func (groupApi *GroupApi) GetGroupList(c *gin.Context) {
	var pageInfo checkinsReq.GroupSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := groupService.GetGroupInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("lấy thất bại!", zap.Error(err))
		response.FailWithMessage("lấy thất bại:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "lấy thành công", c)
}

func (groupApi *GroupApi) GetGroupDataSource(c *gin.Context) {
	// Đây là một API để lấy dữ liệu nguồn định nghĩa của Nhóm
	dataSource, err := groupService.GetGroupDataSource()
	if err != nil {
		global.GVA_LOG.Error("lấy thất bại!", zap.Error(err))
		response.FailWithMessage("lấy thất bại:"+err.Error(), c)
		return
	}
	response.OkWithData(dataSource, c)
}

func (groupApi *GroupApi) AssignParticipantToGroupAuto(c *gin.Context) {
	var info checkinsReq.GroupAuto
	err := c.ShouldBindJSON(&info)

	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = groupService.AssignParticipantToGroupAuto(info)

	if err != nil {
		global.GVA_LOG.Error("lấy thất bại!", zap.Error(err))
		response.FailWithMessage("lấy thất bại:"+err.Error(), c)
		return
	}

	response.OkWithMessage("Gắn nhóm tự động cho sinh viên thành công", c)
}

func (groupApi *GroupApi) GetGroupPublic(c *gin.Context) {
	// API này không cần xác thực
	// Ví dụ trả về một thông điệp cố định, thường được sử dụng cho dịch vụ phía C, cần triển khai logic kinh doanh của riêng mình
	response.OkWithDetailed(gin.H{
		"info": "Thông tin API của Nhóm không cần xác thực",
	}, "lấy thành công", c)
}

func (groupApi *GroupApi) AddMemberToGroup(c *gin.Context) {
	var member checkinsReq.GroupMember
	err := c.ShouldBindJSON(&member)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = groupService.AddMemberToGroup(member)
	if err != nil {
		global.GVA_LOG.Error("Thêm thành viên thất bại!", zap.Error(err))
		response.FailWithMessage("Thêm thành viên thất bại:"+err.Error(), c)
		return
	}
	response.OkWithMessage("Thêm thành viên thành công", c)
}

// GetParticipantsForGroup lấy danh sách sinh viên chưa có nhóm theo attendanceId và query
func (groupApi *GroupApi) GetParticipantsForGroup(c *gin.Context) {
	attendanceIdStr := c.Query("attendanceId")
	query := c.Query("query")

	if attendanceIdStr == "" {
		response.FailWithMessage("attendanceId không được để trống", c)
		return
	}

	attendanceId, err := strconv.Atoi(attendanceIdStr)
	if err != nil {
		response.FailWithMessage("attendanceId không hợp lệ", c)
		return
	}

	// Khai báo slice để chứa kết quả
	var participants []checkins.AttendanceGroupParticipant

	// Truy vấn danh sách sinh viên chưa được phân nhóm (group_id IS NULL)
	db := global.GVA_DB.Table(checkins.AttendanceGroupParticipant{}.TableName()).Where("attendance_id = ? AND group_id IS NULL", attendanceId)

	// Nếu có query, lọc theo full_name hoặc email (giả sử các cột này tồn tại trong bảng)
	if query != "" {
		db = db.Where("full_name LIKE ? OR email LIKE ?", "%"+query+"%", "%"+query+"%")
	}

	if err := db.Find(&participants).Error; err != nil {
		response.FailWithMessage("Lấy danh sách thất bại: "+err.Error(), c)
		return
	}

	response.OkWithData(participants, c)
}
