import service from '@/utils/request'

// @Tags Group
// @Summary 创建Nhóm
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Group true "创建Nhóm"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /group/createGroup [post]
export const createGroup = (data) => {
  return service({
    url: '/group/createGroup',
    method: 'post',
    data
  })
}

// @Tags Group
// @Summary 删除Nhóm
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Group true "删除Nhóm"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /group/deleteGroup [delete]
export const deleteGroup = (params) => {
  return service({
    url: '/group/deleteGroup',
    method: 'delete',
    params
  })
}

// @Tags Group
// @Summary 批量删除Nhóm
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Nhóm"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /group/deleteGroup [delete]
export const deleteGroupByIds = (params) => {
  return service({
    url: '/group/deleteGroupByIds',
    method: 'delete',
    params
  })
}

// @Tags Group
// @Summary 更新Nhóm
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Group true "更新Nhóm"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /group/updateGroup [put]
export const updateGroup = (data) => {
  return service({
    url: '/group/updateGroup',
    method: 'put',
    data
  })
}

// @Tags Group
// @Summary 用id查询Nhóm
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.Group true "用id查询Nhóm"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /group/findGroup [get]
export const findGroup = (params) => {
  return service({
    url: '/group/findGroup',
    method: 'get',
    params
  })
}

// @Tags Group
// @Summary 分页获取Nhóm列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取Nhóm列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /group/getGroupList [get]
export const getGroupList = (params) => {
  return service({
    url: '/group/getGroupList',
    method: 'get',
    params
  })
}
// @Tags Group
// @Summary 获取数据源
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /group/findGroupDataSource [get]
export const getGroupDataSource = () => {
  return service({
    url: '/group/getGroupDataSource',
    method: 'get',
  })
}

export const assignParticipantToGroupAuto = (data) => {
  return service({
    url: '/group/assignParticipantToGroupAuto',
    method: 'post',
    data
  })
}
