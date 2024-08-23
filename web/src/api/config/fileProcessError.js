import service from '@/utils/request'

// @Tags FileProcessError
// @Summary 创建Lỗi ghi nhận từ quá trình xử lý tệp tin
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.FileProcessError true "创建Lỗi ghi nhận từ quá trình xử lý tệp tin"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /fileProcessError/createFileProcessError [post]
export const createFileProcessError = (data) => {
  return service({
    url: '/fileProcessError/createFileProcessError',
    method: 'post',
    data
  })
}

// @Tags FileProcessError
// @Summary 删除Lỗi ghi nhận từ quá trình xử lý tệp tin
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.FileProcessError true "删除Lỗi ghi nhận từ quá trình xử lý tệp tin"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /fileProcessError/deleteFileProcessError [delete]
export const deleteFileProcessError = (params) => {
  return service({
    url: '/fileProcessError/deleteFileProcessError',
    method: 'delete',
    params
  })
}

// @Tags FileProcessError
// @Summary 批量删除Lỗi ghi nhận từ quá trình xử lý tệp tin
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Lỗi ghi nhận từ quá trình xử lý tệp tin"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /fileProcessError/deleteFileProcessError [delete]
export const deleteFileProcessErrorByIds = (params) => {
  return service({
    url: '/fileProcessError/deleteFileProcessErrorByIds',
    method: 'delete',
    params
  })
}

// @Tags FileProcessError
// @Summary 更新Lỗi ghi nhận từ quá trình xử lý tệp tin
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.FileProcessError true "更新Lỗi ghi nhận từ quá trình xử lý tệp tin"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /fileProcessError/updateFileProcessError [put]
export const updateFileProcessError = (data) => {
  return service({
    url: '/fileProcessError/updateFileProcessError',
    method: 'put',
    data
  })
}

// @Tags FileProcessError
// @Summary 用id查询Lỗi ghi nhận từ quá trình xử lý tệp tin
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.FileProcessError true "用id查询Lỗi ghi nhận từ quá trình xử lý tệp tin"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /fileProcessError/findFileProcessError [get]
export const findFileProcessError = (params) => {
  return service({
    url: '/fileProcessError/findFileProcessError',
    method: 'get',
    params
  })
}

// @Tags FileProcessError
// @Summary 分页获取Lỗi ghi nhận từ quá trình xử lý tệp tin列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取Lỗi ghi nhận từ quá trình xử lý tệp tin列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /fileProcessError/getFileProcessErrorList [get]
export const getFileProcessErrorList = (params) => {
  return service({
    url: '/fileProcessError/getFileProcessErrorList',
    method: 'get',
    params
  })
}
