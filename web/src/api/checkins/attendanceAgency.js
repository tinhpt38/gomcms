import service from '@/utils/request'

// @Tags AttendanceAgency
// @Summary 创建Đơn vị
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.AttendanceAgency true "创建Đơn vị"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /attendanceAgency/createAttendanceAgency [post]
export const createAttendanceAgency = (data) => {
  return service({
    url: '/attendanceAgency/createAttendanceAgency',
    method: 'post',
    data
  })
}

// @Tags AttendanceAgency
// @Summary 删除Đơn vị
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.AttendanceAgency true "删除Đơn vị"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /attendanceAgency/deleteAttendanceAgency [delete]
export const deleteAttendanceAgency = (params) => {
  return service({
    url: '/attendanceAgency/deleteAttendanceAgency',
    method: 'delete',
    params
  })
}

// @Tags AttendanceAgency
// @Summary 批量删除Đơn vị
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Đơn vị"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /attendanceAgency/deleteAttendanceAgency [delete]
export const deleteAttendanceAgencyByIds = (params) => {
  return service({
    url: '/attendanceAgency/deleteAttendanceAgencyByIds',
    method: 'delete',
    params
  })
}

// @Tags AttendanceAgency
// @Summary 更新Đơn vị
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.AttendanceAgency true "更新Đơn vị"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /attendanceAgency/updateAttendanceAgency [put]
export const updateAttendanceAgency = (data) => {
  return service({
    url: '/attendanceAgency/updateAttendanceAgency',
    method: 'put',
    data
  })
}

// @Tags AttendanceAgency
// @Summary 用id查询Đơn vị
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.AttendanceAgency true "用id查询Đơn vị"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /attendanceAgency/findAttendanceAgency [get]
export const findAttendanceAgency = (params) => {
  return service({
    url: '/attendanceAgency/findAttendanceAgency',
    method: 'get',
    params
  })
}

// @Tags AttendanceAgency
// @Summary 分页获取Đơn vị列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取Đơn vị列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /attendanceAgency/getAttendanceAgencyList [get]
export const getAttendanceAgencyList = (params) => {
  return service({
    url: '/attendanceAgency/getAttendanceAgencyList',
    method: 'get',
    params
  })
}
