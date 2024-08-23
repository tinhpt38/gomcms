package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type ConditionSearch struct{
    StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
    EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
    GroupId  *int `json:"groupId" form:"groupId" `
    AreaId  *int `json:"areaId" form:"areaId" `
    request.PageInfo
}
