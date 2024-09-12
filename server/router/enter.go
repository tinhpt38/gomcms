package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/router/checkins"
	"github.com/flipped-aurora/gin-vue-admin/server/router/config"
	"github.com/flipped-aurora/gin-vue-admin/server/router/example"
	"github.com/flipped-aurora/gin-vue-admin/server/router/system"
	"github.com/flipped-aurora/gin-vue-admin/server/router/uibuilder"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System    system.RouterGroup
	Example   example.RouterGroup
	Checkins  checkins.RouterGroup
	Config    config.RouterGroup
	Uibuilder uibuilder.RouterGroup
}
