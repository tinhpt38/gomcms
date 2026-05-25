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

// TestGetSyncStatusNoRules: phiên không có điều kiện nào.
func TestGetSyncStatusNoRules(t *testing.T) {
	cleanupSyncTestData(t)
	svc := ConditionService{}

	attId := uint(2003)
	att := checkins.Attendance{Title: "No Rules Test"}
	att.ID = attId
	global.GVA_DB.Create(&att)

	status, err := svc.GetSyncStatusFull(int(attId))
	assert.NoError(t, err)
	assert.Equal(t, SyncStateNoRules, status.SyncState)
	assert.False(t, status.NeedsSync)
}

// TestGetSyncStatusNeverSynced: có điều kiện nhưng chưa sync lần nào.
func TestGetSyncStatusNeverSynced(t *testing.T) {
	cleanupSyncTestData(t)
	svc := ConditionService{}

	attId := uint(2003)
	att := checkins.Attendance{Title: "Never Synced Test"}
	att.ID = attId
	global.GVA_DB.Create(&att)

	cond := checkins.Condition{AttendanceId: &attId}
	global.GVA_DB.Create(&cond)

	status, err := svc.GetSyncStatusFull(int(attId))
	assert.NoError(t, err)
	assert.Equal(t, int64(1), status.ConditionCount)
	assert.Equal(t, int64(0), status.AgpConditionCount)
	assert.Equal(t, SyncStateNeverSynced, status.SyncState)
	assert.True(t, status.NeedsSync)
}

// TestGetSyncStatusStale: đã sync, sau đó thêm AGP → stale.
func TestGetSyncStatusStale(t *testing.T) {
	cleanupSyncTestData(t)
	svc := ConditionService{}

	attId := uint(2006)
	att := checkins.Attendance{Title: "Stale Test"}
	att.ID = attId
	global.GVA_DB.Create(&att)

	group := checkins.Group{Name: "Nhóm X", AttendanceId: attId}
	global.GVA_DB.Create(&group)

	p := checkins.Participant{Email: "stale@test.com"}
	global.GVA_DB.Create(&p)

	agp := checkins.AttendanceGroupParticipant{
		ParticipantId: &p.ID,
		AttendanceId:  &attId,
		GroupId:       &group.ID,
	}
	global.GVA_DB.Create(&agp)

	cond := checkins.Condition{AttendanceId: &attId, GroupId: &group.ID}
	global.GVA_DB.Create(&cond)

	// Sync lần đầu → synced
	err := svc.SyncCondtionForAllMember(int(attId))
	assert.NoError(t, err)
	status, err := svc.GetSyncStatusFull(int(attId))
	assert.NoError(t, err)
	assert.Equal(t, SyncStateSynced, status.SyncState, "Sau sync phải là synced")

	// Thêm thành viên mới chưa sync → stale
	p2 := checkins.Participant{Email: "stale2@test.com"}
	global.GVA_DB.Create(&p2)
	agp2 := checkins.AttendanceGroupParticipant{
		ParticipantId: &p2.ID,
		AttendanceId:  &attId,
		GroupId:       &group.ID,
	}
	global.GVA_DB.Create(&agp2)

	status, err = svc.GetSyncStatusFull(int(attId))
	assert.NoError(t, err)
	assert.Equal(t, SyncStateStale, status.SyncState, "Sau thêm AGP mới phải là stale")
	assert.True(t, status.NeedsSync)
	assert.Greater(t, status.AgpWithoutMappingCount, int64(0))
}

// TestSyncNoNullConditions: sau sync không có dòng condition_id NULL.
func TestSyncNoNullConditions(t *testing.T) {
	cleanupSyncTestData(t)
	svc := ConditionService{}

	attId := uint(2007)
	att := checkins.Attendance{Title: "No Null Test"}
	att.ID = attId
	global.GVA_DB.Create(&att)

	group := checkins.Group{Name: "Nhóm Y", AttendanceId: attId}
	global.GVA_DB.Create(&group)

	p := checkins.Participant{Email: "nonull@test.com"}
	global.GVA_DB.Create(&p)
	agp := checkins.AttendanceGroupParticipant{
		ParticipantId: &p.ID,
		AttendanceId:  &attId,
		GroupId:       &group.ID,
	}
	global.GVA_DB.Create(&agp)

	// Điều kiện cho nhóm khác — AGP này sẽ không match → không được insert NULL
	otherGroup := checkins.Group{Name: "Nhóm Z", AttendanceId: attId}
	global.GVA_DB.Create(&otherGroup)
	cond := checkins.Condition{AttendanceId: &attId, GroupId: &otherGroup.ID}
	global.GVA_DB.Create(&cond)

	err := svc.SyncCondtionForAllMember(int(attId))
	assert.NoError(t, err)

	var nullCount int64
	global.GVA_DB.Raw(
		"SELECT COUNT(*) FROM agp_conditions WHERE attendance_id = ? AND condition_id IS NULL", attId,
	).Scan(&nullCount)
	assert.Equal(t, int64(0), nullCount, "Không được có dòng condition_id NULL sau sync")
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
