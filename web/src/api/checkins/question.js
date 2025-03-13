import service from '@/utils/request'

// @Tags Question
// @Summary 创建 Câu hỏi
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Question true "创建 Câu hỏi"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /question/createQuestion [post]
export const createQuestion = (data) => {
  return service({
    url: '/question/createQuestion',
    method: 'post',
    data
  })
}

export const bulkQuestions = (data) => {
  return service({
    url: '/question/bulkQuestions',
    method: 'post',
    data
  })
}

// @Tags Question
// @Summary 删除 Câu hỏi
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Question true "删除 Câu hỏi"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /question/deleteQuestion [delete]
export const deleteQuestion = (params) => {
  return service({
    url: '/question/deleteQuestion',
    method: 'delete',
    params
  })
}

// @Tags Question
// @Summary 批量删除 Câu hỏi
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除 Câu hỏi"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /question/deleteQuestionByIds [delete]
export const deleteQuestionByIds = (params) => {
  return service({
    url: '/question/deleteQuestionByIds',
    method: 'delete',
    params
  })
}

// @Tags Question
// @Summary 更新 Câu hỏi
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Question true "更新 Câu hỏi"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /question/updateQuestion [put]
export const updateQuestion = (data) => {
  return service({
    url: '/question/updateQuestion',
    method: 'put',
    data
  })
}

// @Tags Question
// @Summary 用 id 查询 Câu hỏi
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.Question true "用 id 查询 Câu hỏi"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /question/findQuestion [get]
export const findQuestion = (params) => {
  return service({
    url: '/question/findQuestion',
    method: 'get',
    params
  })
}

// @Tags Question
// @Summary 分页获取 Câu hỏi 列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取 Câu hỏi 列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /question/getQuestionList [get]
export const getQuestionList = (params) => {
  return service({
    url: '/question/getQuestionList',
    method: 'get',
    params
  })
}
