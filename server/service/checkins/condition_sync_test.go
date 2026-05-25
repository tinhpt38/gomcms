package checkins

import (
	"fmt"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
	checkinsReq "github.com/flipped-aurora/gin-vue-admin/server/model/checkins/request"
	"github.com/stretchr/testify/assert"
)

func cleanupSyncTestData(t *testing.T) {
	t.Helper()
	tables := []string{
		"attendance_check_ins",
		"checkin_logs",
		"agp_conditions",
		"attendance_areas",
		"conditions",
		"attendance_group_participants",
		"participants",
		"groups",
		"areas",
		"attendances",
	}
	for _, table := range tables {
		err := global.GVA_DB.Exec(fmt.Sprintf("DELETE FROM %s", table)).Error
		assert.NoError(t, err)
	}
}

// TestSyncConditionForGroup kiểm tra sync điều kiện theo nhóm cụ thể.
func TestSyncConditionForGroup(t *testing.T) {
	cleanupSyncTestData(t)

	svc := ConditionService{}

	attId := uint(2001)
	att := checkins.Attendance{Title: "Sync Group Test"}
	att.ID = attId
	global.GVA_DB.Create(&att)

	group := checkins.Group{Name: "Nhóm A", AttendanceId: attId}
	global.GVA_DB.Create(&group)

	p := checkins.Participant{Email: "sync_group@test.com"}
	global.GVA_DB.Create(&p)

	agp := checkins.AttendanceGroupParticipant{
		ParticipantId: &p.ID,
		AttendanceId:  &attId,
		GroupId:       &group.ID,
	}
	global.GVA_DB.Create(&agp)

	attIdInt := int(attId)
	cond := checkins.Condition{
		AttendanceId: &attId,
		GroupId:      &group.ID,
	}
	global.GVA_DB.Create(&cond)

	err := svc.SyncCondtionForAllMember(attIdInt)
	assert.NoError(t, err)

	var count int64
	global.GVA_DB.Model(&checkins.AGPCondition{}).Where("attendance_id = ?", attId).Count(&count)
	assert.Equal(t, int64(1), count, "Phải có 1 AGPCondition sau sync theo nhóm")
}

// TestSyncConditionGlobal kiểm tra sync điều kiện toàn phiên (group_id IS NULL).
func TestSyncConditionGlobal(t *testing.T) {
	cleanupSyncTestData(t)

	svc := ConditionService{}

	attId := uint(2002)
	att := checkins.Attendance{Title: "Sync Global Test"}
	att.ID = attId
	global.GVA_DB.Create(&att)

	group := checkins.Group{Name: "Nhóm B", AttendanceId: attId}
	global.GVA_DB.Create(&group)

	p := checkins.Participant{Email: "sync_global@test.com"}
	global.GVA_DB.Create(&p)

	agp := checkins.AttendanceGroupParticipant{
		ParticipantId: &p.ID,
		AttendanceId:  &attId,
		GroupId:       &group.ID,
	}
	global.GVA_DB.Create(&agp)

	attIdInt := int(attId)
	// Điều kiện toàn phiên: group_id IS NULL
	cond := checkins.Condition{AttendanceId: &attId, GroupId: nil}
	global.GVA_DB.Create(&cond)

	err := svc.SyncCondtionForAllMember(attIdInt)
	assert.NoError(t, err)

	var count int64
	global.GVA_DB.Model(&checkins.AGPCondition{}).Where("attendance_id = ?", attId).Count(&count)
	assert.Equal(t, int64(1), count, "Điều kiện toàn phiên phải được gán cho mọi AGP")
}

// TestGetSyncStatus kiểm tra trạng thái needsSync.
func TestGetSyncStatus(t *testing.T) {
	cleanupSyncTestData(t)

	svc := ConditionService{}

	attId := uint(2003)
	att := checkins.Attendance{Title: "Sync Status Test"}
	att.ID = attId
	global.GVA_DB.Create(&att)

	cond := checkins.Condition{AttendanceId: &attId}
	global.GVA_DB.Create(&cond)

	condCount, agpCount, needsSync, err := svc.GetSyncStatus(int(attId))
	assert.NoError(t, err)
	assert.Equal(t, int64(1), condCount)
	assert.Equal(t, int64(0), agpCount)
	assert.True(t, needsSync, "Phải cần sync khi có condition nhưng chưa có agp_condition")
}

// TestBulkCreateParticipantsMultiGroup kiểm tra bulk add với nhiều nhóm.
func TestBulkCreateParticipantsMultiGroup(t *testing.T) {
	cleanupSyncTestData(t)

	svc := ParticipantService{}

	attId := uint(2004)
	att := checkins.Attendance{Title: "Bulk Multi Group"}
	att.ID = attId
	global.GVA_DB.Create(&att)

	groupA := checkins.Group{Name: "Group A", AttendanceId: attId}
	groupB := checkins.Group{Name: "Group B", AttendanceId: attId}
	global.GVA_DB.Create(&groupA)
	global.GVA_DB.Create(&groupB)

	req := checkinsReq.ListEmailParticipantRequest{
		List:         []string{"bulk1@test.com", "bulk2@test.com"},
		AttendanceId: &attId,
		GroupIds:     []uint{groupA.ID, groupB.ID},
	}

	err := svc.BulkCreateParticipants(req)
	assert.NoError(t, err)

	// 2 users × 2 groups = 4 AGP records
	var agpCount int64
	global.GVA_DB.Model(&checkins.AttendanceGroupParticipant{}).Where("attendance_id = ?", attId).Count(&agpCount)
	assert.Equal(t, int64(4), agpCount, "2 users × 2 groups = 4 AGP records")

	// Gọi lại không tạo duplicate
	err = svc.BulkCreateParticipants(req)
	assert.NoError(t, err)
	global.GVA_DB.Model(&checkins.AttendanceGroupParticipant{}).Where("attendance_id = ?", attId).Count(&agpCount)
	assert.Equal(t, int64(4), agpCount, "Không tạo duplicate khi bulk lại")
}

// TestReassignParticipantsOnly kiểm tra chia lại thành viên không xóa nhóm.
func TestReassignParticipantsOnly(t *testing.T) {
	cleanupSyncTestData(t)

	svc := GroupService{}

	attId := uint(2005)
	att := checkins.Attendance{Title: "Reassign Test"}
	att.ID = attId
	global.GVA_DB.Create(&att)

	groupA := checkins.Group{Name: "Alpha", AttendanceId: attId}
	groupB := checkins.Group{Name: "Beta", AttendanceId: attId}
	global.GVA_DB.Create(&groupA)
	global.GVA_DB.Create(&groupB)

	// 4 thành viên, tất cả trong groupA
	for i := 0; i < 4; i++ {
		p := checkins.Participant{Email: fmt.Sprintf("reassign%d@test.com", i)}
		global.GVA_DB.Create(&p)
		agp := checkins.AttendanceGroupParticipant{
			ParticipantId: &p.ID,
			AttendanceId:  &attId,
			GroupId:       &groupA.ID,
		}
		global.GVA_DB.Create(&agp)
	}

	err := svc.ReassignParticipantsOnly(int(attId))
	assert.NoError(t, err)

	var countA, countB int64
	global.GVA_DB.Model(&checkins.AttendanceGroupParticipant{}).
		Where("attendance_id = ? AND group_id = ?", attId, groupA.ID).Count(&countA)
	global.GVA_DB.Model(&checkins.AttendanceGroupParticipant{}).
		Where("attendance_id = ? AND group_id = ?", attId, groupB.ID).Count(&countB)

	assert.Equal(t, int64(4), countA+countB, "Tổng AGP phải giữ nguyên")
	assert.True(t, countA > 0 && countB > 0, "Cả 2 nhóm phải có thành viên sau reassign")

	// Tên nhóm không thay đổi
	var names []string
	global.GVA_DB.Model(&checkins.Group{}).Where("attendance_id = ?", attId).Pluck("name", &names)
	assert.Contains(t, names, "Alpha")
	assert.Contains(t, names, "Beta")
}
