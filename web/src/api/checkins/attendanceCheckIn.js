import service from '@/utils/request'

// @Tags AttendanceCheckIn
// @Summary 创建Thành viên checkins
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.AttendanceCheckIn true "创建Thành viên checkins"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /attendanceCheckIn/createAttendanceCheckIn [post]
export const createAttendanceCheckIn = (data) => {
  return service({
    url: '/attendanceCheckIn/createAttendanceCheckIn',
    method: 'post',
    data
  })
}

export const publicAttendanceCheckIn = (data) => {
  return service({
    url: '/attendanceCheckIn/publicAttendanceCheckIn',
    method: 'post',
    data
  })
}

// @Tags AttendanceCheckIn
// @Summary 删除Thành viên checkins
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.AttendanceCheckIn true "删除Thành viên checkins"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /attendanceCheckIn/deleteAttendanceCheckIn [delete]
export const deleteAttendanceCheckIn = (params) => {
  return service({
    url: '/attendanceCheckIn/deleteAttendanceCheckIn',
    method: 'delete',
    params
  })
}

// @Tags AttendanceCheckIn
// @Summary 批量删除Thành viên checkins
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Thành viên checkins"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /attendanceCheckIn/deleteAttendanceCheckIn [delete]
export const deleteAttendanceCheckInByIds = (params) => {
  return service({
    url: '/attendanceCheckIn/deleteAttendanceCheckInByIds',
    method: 'delete',
    params
  })
}

// @Tags AttendanceCheckIn
// @Summary 更新Thành viên checkins
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.AttendanceCheckIn true "更新Thành viên checkins"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /attendanceCheckIn/updateAttendanceCheckIn [put]
export const updateAttendanceCheckIn = (data) => {
  return service({
    url: '/attendanceCheckIn/updateAttendanceCheckIn',
    method: 'put',
    data
  })
}

// @Tags AttendanceCheckIn
// @Summary 用id查询Thành viên checkins
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.AttendanceCheckIn true "用id查询Thành viên checkins"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /attendanceCheckIn/findAttendanceCheckIn [get]
export const findAttendanceCheckIn = (params) => {
  return service({
    url: '/attendanceCheckIn/findAttendanceCheckIn',
    method: 'get',
    params
  })
}

// @Tags AttendanceCheckIn
// @Summary 分页获取Thành viên checkins列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取Thành viên checkins列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /attendanceCheckIn/getAttendanceCheckInList [get]
export const getAttendanceCheckInList = (params) => {
  return service({
    url: '/attendanceCheckIn/getAttendanceCheckInList',
    method: 'get',
    params
  })
}


export const getAttendanceCheckInLogList = (params) => {
  return service({
    url: '/attendanceCheckIn/getAttendanceCheckInLogList',
    method: 'get',
    params
  })
}
