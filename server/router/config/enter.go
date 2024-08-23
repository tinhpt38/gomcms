package config

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct{ CfgFileProcessRouter }

var FileProcessApi = api.ApiGroupApp.ConfigApiGroup.CfgFileProcessApi
