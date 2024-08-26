package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
	"github.com/flipped-aurora/gin-vue-admin/server/model/config"
)

func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate(
		checkins.Attendance{},
		checkins.Group{},
		checkins.Area{},
		checkins.AttendanceArea{},
		config.CfgFileProcess{},
		checkins.Participant{},
		config.FileProcessError{},
		checkins.Condition{},
		checkins.AttendanceCheckIn{},
		checkins.AttendanceGroupParticipant{})
	if err != nil {
		return err
	}
	return nil
}
