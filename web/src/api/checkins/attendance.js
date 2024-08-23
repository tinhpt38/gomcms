import service from '@/utils/request'

export const createAttendance = (data) => {
  return service({
    url: '/attendance/createAttendance',
    method: 'post',
    data
  })
}

export const deleteAttendance = (params) => {
  return service({
    url: '/attendance/deleteAttendance',
    method: 'delete',
    params
  })
}

export const deleteAttendanceByIds = (params) => {
  return service({
    url: '/attendance/deleteAttendanceByIds',
    method: 'delete',
    params
  })
}

export const updateAttendance = (data) => {
  return service({
    url: '/attendance/updateAttendance',
    method: 'put',
    data
  })
}

export const findAttendance = (params) => {
  return service({
    url: '/attendance/findAttendance',
    method: 'get',
    params
  })
}

export const getAttendanceList = (params) => {
  return service({
    url: '/attendance/getAttendanceList',
    method: 'get',
    params
  })
}
