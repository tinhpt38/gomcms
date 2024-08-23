package config

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	FileProcessErrorRouter
	CfgFileProcessRouter
}

var (
	fileProcessErrorApi = api.ApiGroupApp.ConfigApiGroup.FileProcessErrorApi
	fileProcessApi      = api.ApiGroupApp.ConfigApiGroup.CfgFileProcessApi
)
