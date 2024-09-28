import service from '@/utils/request'

// @Tags Participant
// @Summary 创建Sinh viên (Người tham dự phiên điểm danh)
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Participant true "创建Sinh viên (Người tham dự phiên điểm danh)"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /participant/createParticipant [post]
export const createParticipant = (data) => {
  return service({
    url: '/participant/createParticipant',
    method: 'post',
    data
  })
}

export const bulkParticipants = (data) => {
  return service({
    url: '/participant/bulkParticipants',
    method: 'post',
    data
  })
}

// @Tags Participant
// @Summary 删除Sinh viên (Người tham dự phiên điểm danh)
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Participant true "删除Sinh viên (Người tham dự phiên điểm danh)"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /participant/deleteParticipant [delete]
export const deleteParticipant = (params) => {
  return service({
    url: '/participant/deleteParticipant',
    method: 'delete',
    params
  })
}

// @Tags Participant
// @Summary 批量删除Sinh viên (Người tham dự phiên điểm danh)
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Sinh viên (Người tham dự phiên điểm danh)"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /participant/deleteParticipant [delete]
export const deleteParticipantByIds = (params) => {
  return service({
    url: '/participant/deleteParticipantByIds',
    method: 'delete',
    params
  })
}

// @Tags Participant
// @Summary 更新Sinh viên (Người tham dự phiên điểm danh)
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Participant true "更新Sinh viên (Người tham dự phiên điểm danh)"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /participant/updateParticipant [put]
export const updateParticipant = (data) => {
  return service({
    url: '/participant/updateParticipant',
    method: 'put',
    data
  })
}

// @Tags Participant
// @Summary 用id查询Sinh viên (Người tham dự phiên điểm danh)
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.Participant true "用id查询Sinh viên (Người tham dự phiên điểm danh)"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /participant/findParticipant [get]
export const findParticipant = (params) => {
  return service({
    url: '/participant/findParticipant',
    method: 'get',
    params
  })
}

export const findLuckyParticipant = (params) => {
  return service({
    url: '/participant/findLuckyParticipant',
    method: 'get',
    params
  })
}

// @Tags Participant
// @Summary 分页获取Sinh viên (Người tham dự phiên điểm danh)列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取Sinh viên (Người tham dự phiên điểm danh)列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /participant/getParticipantList [get]
export const getParticipantList = (params) => {
  return service({
    url: '/participant/getParticipantList',
    method: 'get',
    params
  })
}


export const getParticipantListByAttendance = (params) => {
  return service({
    url: '/participant/getParticipantListByAttendance',
    method: 'get',
    params
  })
}
