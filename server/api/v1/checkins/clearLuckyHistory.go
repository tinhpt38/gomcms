package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (participantApi *ParticipantApi) ClearLuckyHistory(c *gin.Context) {
	acId := c.Query("attendanceId")

	// Xóa lịch sử số may mắn cho attendance này
	err := global.GVA_DB.Where("attendance_id = ?", acId).Delete(&checkins.UsedLuckyParticipant{}).Error

	if err != nil {
		global.GVA_LOG.Error("Thất bại khi xóa lịch sử!", zap.Error(err))
		response.FailWithMessage("Thất bại khi xóa lịch sử:"+err.Error(), c)
		return
	}

	response.OkWithMessage("Đã xóa lịch sử số may mắn thành công", c)
}
