package config

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct{ CfgFileProcessApi }

var FileProcessService = service.ServiceGroupApp.ConfigServiceGroup.CfgFileProcessService
