package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type CfgFileProcessSearch struct{
    StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
    EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
    Uuid  string `json:"uuid" form:"uuid" `
    Percent  *float64 `json:"percent" form:"percent" `
    request.PageInfo
}
