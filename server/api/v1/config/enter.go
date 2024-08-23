package config

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	FileProcessErrorApi
	CfgFileProcessApi
}

var (
	fileProcessErrorService = service.ServiceGroupApp.ConfigServiceGroup.FileProcessErrorService
	fileProcessService      = service.ServiceGroupApp.ConfigServiceGroup.CfgFileProcessService
	participateService      = service.ServiceGroupApp.CheckinsServiceGroup.ParticipantService
)
