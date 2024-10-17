import service from '@/utils/request'

// @Tags Condition
// @Summary 创建Điều kiện để checkins
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Condition true "创建Điều kiện để checkins"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /condition/createCondition [post]
export const createCondition = (data) => {
  return service({
    url: '/condition/createCondition',
    method: 'post',
    data
  })
}

// @Tags Condition
// @Summary 删除Điều kiện để checkins
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Condition true "删除Điều kiện để checkins"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /condition/deleteCondition [delete]
export const deleteCondition = (params) => {
  return service({
    url: '/condition/deleteCondition',
    method: 'delete',
    params
  })
}

// @Tags Condition
// @Summary 批量删除Điều kiện để checkins
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Điều kiện để checkins"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /condition/deleteCondition [delete]
export const deleteConditionByIds = (params) => {
  return service({
    url: '/condition/deleteConditionByIds',
    method: 'delete',
    params
  })
}

// @Tags Condition
// @Summary 更新Điều kiện để checkins
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Condition true "更新Điều kiện để checkins"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /condition/updateCondition [put]
export const updateCondition = (data) => {
  return service({
    url: '/condition/updateCondition',
    method: 'put',
    data
  })
}

// @Tags Condition
// @Summary 用id查询Điều kiện để checkins
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.Condition true "用id查询Điều kiện để checkins"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /condition/findCondition [get]
export const findCondition = (params) => {
  return service({
    url: '/condition/findCondition',
    method: 'get',
    params
  })
}

// @Tags Condition
// @Summary 分页获取Điều kiện để checkins列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取Điều kiện để checkins列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /condition/getConditionList [get]
export const getConditionList = (params) => {
  return service({
    url: '/condition/getConditionList',
    method: 'get',
    params
  })
}

export const syncCondition = (data) => {
  return service({
    url: '/condition/syncCondition',
    method: 'post',
    data
  })
}
