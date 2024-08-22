import service from '@/utils/request'

// @Tags AttendanceClass
// @Summary 创建AttendanceClass
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.AttendanceClass true "创建AttendanceClass"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /attendanceClass/createAttendanceClass [post]
export const createAttendanceClass = (data) => {
  return service({
    url: '/attendanceClass/createAttendanceClass',
    method: 'post',
    data
  })
}

// @Tags AttendanceClass
// @Summary 删除AttendanceClass
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.AttendanceClass true "删除AttendanceClass"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /attendanceClass/deleteAttendanceClass [delete]
export const deleteAttendanceClass = (params) => {
  return service({
    url: '/attendanceClass/deleteAttendanceClass',
    method: 'delete',
    params
  })
}

// @Tags AttendanceClass
// @Summary 批量删除AttendanceClass
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除AttendanceClass"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /attendanceClass/deleteAttendanceClass [delete]
export const deleteAttendanceClassByIds = (params) => {
  return service({
    url: '/attendanceClass/deleteAttendanceClassByIds',
    method: 'delete',
    params
  })
}

// @Tags AttendanceClass
// @Summary 更新AttendanceClass
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.AttendanceClass true "更新AttendanceClass"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /attendanceClass/updateAttendanceClass [put]
export const updateAttendanceClass = (data) => {
  return service({
    url: '/attendanceClass/updateAttendanceClass',
    method: 'put',
    data
  })
}

// @Tags AttendanceClass
// @Summary 用id查询AttendanceClass
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.AttendanceClass true "用id查询AttendanceClass"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /attendanceClass/findAttendanceClass [get]
export const findAttendanceClass = (params) => {
  return service({
    url: '/attendanceClass/findAttendanceClass',
    method: 'get',
    params
  })
}

// @Tags AttendanceClass
// @Summary 分页获取AttendanceClass列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取AttendanceClass列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /attendanceClass/getAttendanceClassList [get]
export const getAttendanceClassList = (params) => {
  return service({
    url: '/attendanceClass/getAttendanceClassList',
    method: 'get',
    params
  })
}
