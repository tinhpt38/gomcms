import service from '@/utils/request'

// @Tags AttendanceCategory
// @Summary 创建Danh mục điểm danh
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.AttendanceCategory true "创建Danh mục điểm danh"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /attendanceCategory/createAttendanceCategory [post]
export const createAttendanceCategory = (data) => {
  return service({
    url: '/attendanceCategory/createAttendanceCategory',
    method: 'post',
    data
  })
}

// @Tags AttendanceCategory
// @Summary 删除Danh mục điểm danh
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.AttendanceCategory true "删除Danh mục điểm danh"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /attendanceCategory/deleteAttendanceCategory [delete]
export const deleteAttendanceCategory = (params) => {
  return service({
    url: '/attendanceCategory/deleteAttendanceCategory',
    method: 'delete',
    params
  })
}

// @Tags AttendanceCategory
// @Summary 批量删除Danh mục điểm danh
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Danh mục điểm danh"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /attendanceCategory/deleteAttendanceCategory [delete]
export const deleteAttendanceCategoryByIds = (params) => {
  return service({
    url: '/attendanceCategory/deleteAttendanceCategoryByIds',
    method: 'delete',
    params
  })
}

// @Tags AttendanceCategory
// @Summary 更新Danh mục điểm danh
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.AttendanceCategory true "更新Danh mục điểm danh"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /attendanceCategory/updateAttendanceCategory [put]
export const updateAttendanceCategory = (data) => {
  return service({
    url: '/attendanceCategory/updateAttendanceCategory',
    method: 'put',
    data
  })
}

// @Tags AttendanceCategory
// @Summary 用id查询Danh mục điểm danh
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.AttendanceCategory true "用id查询Danh mục điểm danh"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /attendanceCategory/findAttendanceCategory [get]
export const findAttendanceCategory = (params) => {
  return service({
    url: '/attendanceCategory/findAttendanceCategory',
    method: 'get',
    params
  })
}

// @Tags AttendanceCategory
// @Summary 分页获取Danh mục điểm danh列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取Danh mục điểm danh列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /attendanceCategory/getAttendanceCategoryList [get]
export const getAttendanceCategoryList = (params) => {
  return service({
    url: '/attendanceCategory/getAttendanceCategoryList',
    method: 'get',
    params
  })
}


export const getAttendanceCategoryPublic = (params) => {
  return service({
    url: '/attendanceCategory/getAttendanceCategoryPublic',
    method: 'get',
    params
  })
}
