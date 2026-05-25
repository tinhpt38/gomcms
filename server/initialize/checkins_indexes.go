package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"
)

// ensureCheckinsIndexes creates performance indexes for checkins tables.
// Uses IF NOT EXISTS so it is safe to run on an existing database.
func ensureCheckinsIndexes() {
	db := global.GVA_DB
	indexes := []struct {
		name  string
		query string
	}{
		// agp_conditions: fast lookup by attendance_id and deduplicate (agp_id, condition_id)
		{
			name:  "idx_agp_conditions_attendance_id",
			query: "CREATE INDEX IF NOT EXISTS idx_agp_conditions_attendance_id ON agp_conditions (attendance_id)",
		},
		{
			name:  "idx_agp_conditions_agp_id",
			query: "CREATE INDEX IF NOT EXISTS idx_agp_conditions_agp_id ON agp_conditions (agp_id)",
		},
		// attendance_group_participants: filter by attendance + participant or group
		{
			name:  "idx_agp_attendance_participant",
			query: "CREATE INDEX IF NOT EXISTS idx_agp_attendance_participant ON attendance_group_participants (attendance_id, participant_id)",
		},
		{
			name:  "idx_agp_attendance_group",
			query: "CREATE INDEX IF NOT EXISTS idx_agp_attendance_group ON attendance_group_participants (attendance_id, group_id)",
		},
		// conditions: filter by attendance_id
		{
			name:  "idx_conditions_attendance_id",
			query: "CREATE INDEX IF NOT EXISTS idx_conditions_attendance_id ON conditions (attendance_id)",
		},
	}

	for _, idx := range indexes {
		if err := db.Exec(idx.query).Error; err != nil {
			global.GVA_LOG.Warn("failed to create index "+idx.name, zap.Error(err))
		}
	}
}
