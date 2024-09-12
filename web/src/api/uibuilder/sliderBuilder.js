import service from '@/utils/request'

// @Tags SliderBuilder
// @Summary 创建SliderBuilder
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SliderBuilder true "创建SliderBuilder"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /sliderBuilder/createSliderBuilder [post]
export const createSliderBuilder = (data) => {
  return service({
    url: '/sliderBuilder/createSliderBuilder',
    method: 'post',
    data
  })
}

// @Tags SliderBuilder
// @Summary 删除SliderBuilder
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SliderBuilder true "删除SliderBuilder"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sliderBuilder/deleteSliderBuilder [delete]
export const deleteSliderBuilder = (params) => {
  return service({
    url: '/sliderBuilder/deleteSliderBuilder',
    method: 'delete',
    params
  })
}

// @Tags SliderBuilder
// @Summary 批量删除SliderBuilder
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除SliderBuilder"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sliderBuilder/deleteSliderBuilder [delete]
export const deleteSliderBuilderByIds = (params) => {
  return service({
    url: '/sliderBuilder/deleteSliderBuilderByIds',
    method: 'delete',
    params
  })
}

// @Tags SliderBuilder
// @Summary 更新SliderBuilder
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SliderBuilder true "更新SliderBuilder"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /sliderBuilder/updateSliderBuilder [put]
export const updateSliderBuilder = (data) => {
  return service({
    url: '/sliderBuilder/updateSliderBuilder',
    method: 'put',
    data
  })
}

// @Tags SliderBuilder
// @Summary 用id查询SliderBuilder
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.SliderBuilder true "用id查询SliderBuilder"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /sliderBuilder/findSliderBuilder [get]
export const findSliderBuilder = (params) => {
  return service({
    url: '/sliderBuilder/findSliderBuilder',
    method: 'get',
    params
  })
}

// @Tags SliderBuilder
// @Summary 分页获取SliderBuilder列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取SliderBuilder列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sliderBuilder/getSliderBuilderList [get]
export const getSliderBuilderList = (params) => {
  return service({
    url: '/sliderBuilder/getSliderBuilderList',
    method: 'get',
    params
  })
}

export const getSliderBuilderPublic = (params) => {
  return service({
    url: '/sliderBuilder/getSliderBuilderPublic',
    method: 'get',
    params
  })
}
