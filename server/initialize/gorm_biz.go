package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
)

func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate(checkins.Attendance{}, checkins.Group{}, checkins.Area{}, checkins.AttendanceArea{})
	if err != nil {
		return err
	}
	return nil
}
