import service from '@/utils/request'

// @Tags CfgFileProcess
// @Summary 创建Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.CfgFileProcess true "创建Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /fileProcess/createCfgFileProcess [post]
export const createCfgFileProcess = (data) => {
  return service({
    url: '/fileProcess/createCfgFileProcess',
    method: 'post',
    data
  })
}

// @Tags CfgFileProcess
// @Summary 删除Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.CfgFileProcess true "删除Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /fileProcess/deleteCfgFileProcess [delete]
export const deleteCfgFileProcess = (params) => {
  return service({
    url: '/fileProcess/deleteCfgFileProcess',
    method: 'delete',
    params
  })
}

// @Tags CfgFileProcess
// @Summary 批量删除Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /fileProcess/deleteCfgFileProcess [delete]
export const deleteCfgFileProcessByIds = (params) => {
  return service({
    url: '/fileProcess/deleteCfgFileProcessByIds',
    method: 'delete',
    params
  })
}

// @Tags CfgFileProcess
// @Summary 更新Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.CfgFileProcess true "更新Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /fileProcess/updateCfgFileProcess [put]
export const updateCfgFileProcess = (data) => {
  return service({
    url: '/fileProcess/updateCfgFileProcess',
    method: 'put',
    data
  })
}

// @Tags CfgFileProcess
// @Summary 用id查询Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.CfgFileProcess true "用id查询Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /fileProcess/findCfgFileProcess [get]
export const findCfgFileProcess = (params) => {
  return service({
    url: '/fileProcess/findCfgFileProcess',
    method: 'get',
    params
  })
}

// @Tags CfgFileProcess
// @Summary 分页获取Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取Lưu trữ thông tin nhập xuất dữ liệu Excel của hệ thống列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /fileProcess/getCfgFileProcessList [get]
export const getCfgFileProcessList = (params) => {
  return service({
    url: '/fileProcess/getCfgFileProcessList',
    method: 'get',
    params
  })
}


export const getPercentFileProcess = (params) => {
  return service({
    url: '/fileProcess/getPercentFileProcess',
    method: 'get',
    params
  })
}