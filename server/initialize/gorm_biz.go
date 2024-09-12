package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
	"github.com/flipped-aurora/gin-vue-admin/server/model/config"
	"github.com/flipped-aurora/gin-vue-admin/server/model/uibuilder"
)

func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate(checkins.Attendance{}, checkins.Group{}, checkins.Area{}, checkins.AttendanceArea{}, config.CfgFileProcess{}, checkins.Participant{}, config.FileProcessError{}, checkins.Condition{}, checkins.AttendanceCheckIn{}, checkins.AttendanceGroupParticipant{}, checkins.AttendanceCategory{}, checkins.AttendanceAgency{}, uibuilder.SliderBuilder{})
	if err != nil {
		return err
	}
	return nil
}
