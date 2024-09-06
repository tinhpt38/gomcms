import service from '@/utils/request'

export const createAttendance = (data) => {
  return service({
    url: '/attendance/createAttendance',
    method: 'post',
    data
  })
}


export const cloneAttendance = (data) => {
  return service({
    url: '/attendance/cloneAttendance',
    method: 'post',
    data
  })
}

export const createAttendanceArea = (data) => {
  return service({
    url: '/attendance/createAttendanceArea',
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

export const deleteAttendanceArea = (params) => {
  return service({
    url: '/attendance/deleteAttendanceArea',
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

export const findAttendanceArea = (params) => {
  return service({
    url: '/attendance/findAttendanceArea',
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

export const getAttendancePublic = (params) => {
  return service({
    url: '/attendance/getAttendancePublic',
    method: 'get',
    params
  })
}

export const statsByAgencyCategory = (data) => {
  return service({
    url: '/attendance/statsByAgencyCategory',
    method: 'post',
    data
  })
}


export const statsScatterPlot = (data) => {
  return service({
    url: '/attendance/statsScatterPlot',
    method: 'post',
    data
  })
}

export const statsTrendLine = (data) => {
  return service({
    url: '/attendance/statsTrendLine',
    method: 'post',
    data
  })
}


