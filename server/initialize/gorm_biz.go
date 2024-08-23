package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
	"github.com/flipped-aurora/gin-vue-admin/server/model/config"
)

func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate(checkins.AttendanceClass{}, config.CfgFileProcess{})
	if err != nil {
		return err
	}
	return nil
}
