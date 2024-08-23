package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/service/checkins"
	"github.com/flipped-aurora/gin-vue-admin/server/service/config"
	"github.com/flipped-aurora/gin-vue-admin/server/service/example"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
)

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	SystemServiceGroup   system.ServiceGroup
	ExampleServiceGroup  example.ServiceGroup
	CheckinsServiceGroup checkins.ServiceGroup
	ConfigServiceGroup   config.ServiceGroup
}
