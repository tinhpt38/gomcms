package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
)

func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate(checkins.AttendanceClass{}, checkins.Group{})
	if err != nil {
		return err
	}
	return nil
}
