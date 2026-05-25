package checkins

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
	checkinsReq "github.com/flipped-aurora/gin-vue-admin/server/model/checkins/request"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ConditionService struct{}

// CreateCondition 创建Điều kiện để checkins记录
// Author [piexlmax](https://github.com/piexlmax)
func (conditionService *ConditionService) CreateCondition(condition *checkins.Condition) (err error) {
	err = global.GVA_DB.Create(condition).Error
	return err
}

// DeleteCondition 删除Điều kiện để checkins记录
// Author [piexlmax](https://github.com/piexlmax)
func (conditionService *ConditionService) DeleteCondition(ID string, userID uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&checkins.Condition{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
			return err
		}
		if err = tx.Delete(&checkins.Condition{}, "id = ?", ID).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteConditionByIds 批量删除Điều kiện để checkins记录
// Author [piexlmax](https://github.com/piexlmax)
func (conditionService *ConditionService) DeleteConditionByIds(IDs []string, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&checkins.Condition{}).Where("id in ?", IDs).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", IDs).Delete(&checkins.Condition{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateCondition 更新Điều kiện để checkins记录
// Author [piexlmax](https://github.com/piexlmax)
func (conditionService *ConditionService) UpdateCondition(condition checkins.Condition) (err error) {
	err = global.GVA_DB.Model(&checkins.Condition{}).Where("id = ?", condition.ID).Updates(&condition).Error
	return err
}

// GetCondition 根据ID获取Điều kiện để checkins记录
// Author [piexlmax](https://github.com/piexlmax)
func (conditionService *ConditionService) GetCondition(ID string) (condition checkins.Condition, err error) {
	err = global.GVA_DB.Where("id = ?", ID).
		Preload(clause.Associations).
		Preload("Group").
		Preload("Area.Area").
		First(&condition).Error
	return
}

// func (conditionService *ConditionService) GetConditionOfPartparticipant(aId uint, pId uint) (list []checkins.Condition, err error) {
// 	// err = global.GVA_DB.Where("attendance_id = ?", ID).Preload(clause.Associations).Preload("Group").Preload("Area.Area").Find(&list).Error
// 	query := `
// 		SELECT
// 	c.*
// FROM
// 	conditions c
// LEFT JOIN (
// 	SELECT
// 		ac.condition_id as conid
// 	FROM
// 		attendance_checkins ac
// 	WHERE
// 		ac.partpaticipant_id = ?
// 		AND ac.attendance_id = ?
// 		AND ac.deleted_at IS NULL
// 		AND ac.condition_id != 0
// ) as ak
// ON
// 	c.id = ak.conid
// WHERE
// 	c.attendance_id = ?
// 	AND ak.conid IS NULL
// 	AND c.deleted_at IS NULL;
// `
// 	err = global.GVA_DB.Raw(query, pId, aId, aId).Preload(clause.Associations).Preload("Area.Area").Debug().Find(&list).Error
// 	return
// }

func (conditionService *ConditionService) GetConditionOfPartparticipant(aId uint, agpIds []int) (list []checkins.AGPCondition, err error) {
	newdb := global.GVA_DB.Model(&checkins.AGPCondition{})
	err = newdb.Joins("LEFT JOIN conditions con ON agp_conditions.condition_id = con.id").
		Where("agp_conditions.attendance_id = ? AND agp_conditions.agp_id IN (?) AND agp_conditions.deleted_at IS NULL AND con.deleted_at IS NULL", aId, agpIds).
		Preload(clause.Associations).
		Preload("Condition.Group").
		Preload("Condition.Area.Area").
		Find(&list).Error
	return
}

// GetConditionInfoList 分页获取Điều kiện để checkins记录
// Author [piexlmax](https://github.com/piexlmax)
func (conditionService *ConditionService) GetConditionInfoList(info checkinsReq.ConditionSearch) (list []checkins.Condition, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&checkins.Condition{})
	var conditions []checkins.Condition
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.GroupId != nil {
		db = db.Where("group_id = ?", info.GroupId)
	}
	if info.AreaId != nil {
		db = db.Where("area_id = ?", info.AreaId)
	}
	if info.AttendanceId != nil {
		db = db.Where("attendance_id = ?", info.AttendanceId)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Preload(clause.Associations).
		Preload("Group").
		Preload("Area.Area").
		Find(&conditions).Error
	return conditions, total, err
}

func (conditionService *ConditionService) GetConditionsByAttendanceId(attId uint) (list []checkins.Condition, err error) {
	db := global.GVA_DB.Model(&checkins.Condition{})
	var conditions []checkins.Condition

	db = db.Where("attendance_id = ?", attId)
	err = db.Preload(clause.Associations).
		Preload("Group").
		Preload("Area.Area").
		Find(&conditions).Error
	return conditions, err
}

// SyncState represents the sync health of agp_conditions for an attendance.
type SyncState string

const (
	SyncStateNoRules    SyncState = "no_rules"
	SyncStateNeverSynced SyncState = "never_synced"
	SyncStateStale      SyncState = "stale"
	SyncStateSynced     SyncState = "synced"
)

// SyncStatus holds the full sync health report for an attendance session.
type SyncStatus struct {
	// Legacy fields — kept for backward compatibility
	ConditionCount    int64 `json:"conditionCount"`
	AgpConditionCount int64 `json:"agpConditionCount"`
	NeedsSync         bool  `json:"needsSync"`

	// Extended fields
	SyncState                 SyncState `json:"syncState"`
	ExpectedAgpConditionCount int64     `json:"expectedAgpConditionCount"`
	AgpCount                  int64     `json:"agpCount"`
	ParticipantCount          int64     `json:"participantCount"`
	AgpWithoutMappingCount    int64     `json:"agpWithoutMappingCount"`
}

// countExpectedAgpConditions returns the number of (agp, condition) pairs that should
// exist after a correct sync.
func countExpectedAgpConditions(attId int) (int64, error) {
	var count int64
	err := global.GVA_DB.Raw(`
		SELECT COUNT(*)
		FROM attendance_group_participants agp
		INNER JOIN conditions c
		  ON c.attendance_id = ?
		  AND c.deleted_at IS NULL
		  AND (c.group_id = agp.group_id OR c.group_id IS NULL)
		WHERE agp.attendance_id = ?
		  AND agp.deleted_at IS NULL`,
		attId, attId).Scan(&count).Error
	return count, err
}

func (conditionService *ConditionService) GetSyncStatus(attId int) (conditionCount, agpConditionCount int64, needsSync bool, err error) {
	status, e := conditionService.GetSyncStatusFull(attId)
	if e != nil {
		err = e
		return
	}
	conditionCount = status.ConditionCount
	agpConditionCount = status.AgpConditionCount
	needsSync = status.NeedsSync
	return
}

func (conditionService *ConditionService) GetSyncStatusFull(attId int) (status SyncStatus, err error) {
	err = global.GVA_DB.Model(&checkins.Condition{}).
		Where("attendance_id = ? AND deleted_at IS NULL", attId).
		Count(&status.ConditionCount).Error
	if err != nil {
		return
	}

	err = global.GVA_DB.Model(&checkins.AGPCondition{}).
		Where("attendance_id = ? AND deleted_at IS NULL", attId).
		Count(&status.AgpConditionCount).Error
	if err != nil {
		return
	}

	err = global.GVA_DB.Model(&checkins.AttendanceGroupParticipant{}).
		Where("attendance_id = ? AND deleted_at IS NULL", attId).
		Count(&status.AgpCount).Error
	if err != nil {
		return
	}

	err = global.GVA_DB.Raw(
		`SELECT COUNT(DISTINCT participant_id)
		 FROM attendance_group_participants
		 WHERE attendance_id = ? AND deleted_at IS NULL`, attId).
		Scan(&status.ParticipantCount).Error
	if err != nil {
		return
	}

	status.ExpectedAgpConditionCount, err = countExpectedAgpConditions(attId)
	if err != nil {
		return
	}

	// AGP that have no agp_condition row despite conditions existing
	if status.ConditionCount > 0 && status.AgpCount > 0 {
		err = global.GVA_DB.Raw(`
			SELECT COUNT(*)
			FROM attendance_group_participants agp
			WHERE agp.attendance_id = ?
			  AND agp.deleted_at IS NULL
			  AND NOT EXISTS (
			    SELECT 1 FROM agp_conditions ac
			    WHERE ac.agp_id = agp.id AND ac.deleted_at IS NULL
			  )`, attId).Scan(&status.AgpWithoutMappingCount).Error
		if err != nil {
			return
		}
	}

	switch {
	case status.ConditionCount == 0:
		status.SyncState = SyncStateNoRules
	case status.AgpConditionCount == 0:
		status.SyncState = SyncStateNeverSynced
	case status.ExpectedAgpConditionCount != status.AgpConditionCount:
		status.SyncState = SyncStateStale
	default:
		status.SyncState = SyncStateSynced
	}

	status.NeedsSync = status.SyncState == SyncStateNeverSynced || status.SyncState == SyncStateStale
	return
}

// SyncAttendanceConditions is the single entry-point for syncing agp_conditions after any
// config change (conditions, AGP, groups). It runs a delta sync when data already exists,
// falling back to a full rebuild when the table is empty for this attendance.
// Sync errors are logged but do NOT fail the calling CRUD operation.
func (conditionService *ConditionService) SyncAttendanceConditions(attId int) {
	svc := conditionService
	var actual int64
	if err := global.GVA_DB.Model(&checkins.AGPCondition{}).
		Where("attendance_id = ? AND deleted_at IS NULL", attId).
		Count(&actual).Error; err != nil {
		global.GVA_LOG.Warn("auto-sync: count failed", zap.Int("attId", attId), zap.Error(err))
		return
	}

	var syncErr error
	if actual == 0 {
		syncErr = svc.SyncCondtionForAllMember(attId)
	} else {
		syncErr = svc.SyncDeltaForAttendance(attId)
	}
	if syncErr != nil {
		global.GVA_LOG.Warn("auto-sync: sync failed", zap.Int("attId", attId), zap.Error(syncErr))
	}
}

// SyncDeltaForAttendance performs an incremental sync:
//  1. Deletes agp_conditions that are no longer in the expected set.
//  2. Inserts missing (agp_id, condition_id) pairs (duplicate-safe via INSERT IGNORE).
func (conditionService *ConditionService) SyncDeltaForAttendance(attId int) (err error) {
	var conditionCount int64
	if err = global.GVA_DB.Model(&checkins.Condition{}).
		Where("attendance_id = ? AND deleted_at IS NULL", attId).
		Count(&conditionCount).Error; err != nil {
		return
	}
	if conditionCount == 0 {
		// No rules → wipe any stale mappings
		return global.GVA_DB.
			Where("attendance_id = ?", attId).
			Unscoped().
			Delete(&checkins.AGPCondition{}).Error
	}

	// Step 1: delete orphaned mappings (agp or condition no longer in expected set)
	deleteOrphans := `
		DELETE ac FROM agp_conditions ac
		WHERE ac.attendance_id = ?
		  AND NOT EXISTS (
		    SELECT 1
		    FROM attendance_group_participants agp
		    INNER JOIN conditions c
		      ON c.attendance_id = ?
		      AND c.deleted_at IS NULL
		      AND (c.group_id = agp.group_id OR c.group_id IS NULL)
		    WHERE agp.id = ac.agp_id
		      AND agp.deleted_at IS NULL
		      AND c.id = ac.condition_id
		  )`
	if err = global.GVA_DB.Exec(deleteOrphans, attId, attId).Error; err != nil {
		return
	}

	// Step 2: insert missing pairs
	insertMissing := `
		INSERT IGNORE INTO agp_conditions (agp_id, condition_id, attendance_id)
		SELECT agp.id, c.id, agp.attendance_id
		FROM attendance_group_participants agp
		INNER JOIN conditions c
		  ON c.attendance_id = ?
		  AND c.deleted_at IS NULL
		  AND (c.group_id = agp.group_id OR c.group_id IS NULL)
		WHERE agp.attendance_id = ?
		  AND agp.deleted_at IS NULL`
	err = global.GVA_DB.Exec(insertMissing, attId, attId).Error
	return
}

func (conditionService *ConditionService) SyncCondtionForAllMember(attId int) (err error) {
	var conditionCount int64
	err = global.GVA_DB.Model(&checkins.Condition{}).
		Where("attendance_id = ? AND deleted_at IS NULL", attId).
		Count(&conditionCount).Error
	if err != nil {
		return
	}
	if conditionCount == 0 {
		return nil
	}

	// Full rebuild: delete existing mappings then re-insert from expected pairs.
	err = global.GVA_DB.
		Where("attendance_id = ?", attId).
		Unscoped().
		Delete(&checkins.AGPCondition{}).Error
	if err != nil {
		return err
	}

	// INNER JOIN ensures only valid (agp, condition) pairs are inserted — no NULL condition_id rows.
	// Operator precedence: bind c.attendance_id explicitly so the OR covers only group_id matching.
	rawQuery := `
		INSERT INTO agp_conditions (agp_id, condition_id, attendance_id)
		SELECT agp.id, c.id, agp.attendance_id
		FROM attendance_group_participants agp
		INNER JOIN conditions c
		  ON c.attendance_id = ?
		  AND c.deleted_at IS NULL
		  AND (c.group_id = agp.group_id OR c.group_id IS NULL)
		WHERE agp.attendance_id = ?
		  AND agp.deleted_at IS NULL`
	err = global.GVA_DB.Exec(rawQuery, attId, attId).Error
	return
}
