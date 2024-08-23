package {{.Package}}

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/model/{{.Package}}"
    {{.Package}}Req "github.com/flipped-aurora/gin-vue-admin/server/model/{{.Package}}/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    {{- if .AutoCreateResource}}
    "github.com/flipped-aurora/gin-vue-admin/server/utils"
    {{- end }}
)

type {{.StructName}}Api struct {}

// Create{{.StructName}} 创建{{.Description}}
// @Tags {{.StructName}}
// @Summary 创建{{.Description}}
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body {{.Package}}.{{.StructName}} true "创建{{.Description}}"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /{{.Abbreviation}}/create{{.StructName}} [post]
func ({{.Abbreviation}}Api *{{.StructName}}Api) Create{{.StructName}}(c *gin.Context) {
	var {{.Abbreviation}} {{.Package}}.{{.StructName}}
	err := c.ShouldBindJSON(&{{.Abbreviation}})
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	{{- if .AutoCreateResource }}
    {{.Abbreviation}}.CreatedBy = utils.GetUserID(c)
	{{- end }}
	err = {{.Abbreviation}}Service.Create{{.StructName}}(&{{.Abbreviation}})
	if err != nil {
        global.GVA_LOG.Error("thất bại!", zap.Error(err))
		response.FailWithMessage("thất bại:" + err.Error(), c)
		return
	}
    response.OkWithMessage("thành công", c)
}

// Delete{{.StructName}} 删除{{.Description}}
// @Tags {{.StructName}}
// @Summary 删除{{.Description}}
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body {{.Package}}.{{.StructName}} true "删除{{.Description}}"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /{{.Abbreviation}}/delete{{.StructName}} [delete]
func ({{.Abbreviation}}Api *{{.StructName}}Api) Delete{{.StructName}}(c *gin.Context) {
	{{.PrimaryField.FieldJson}} := c.Query("{{.PrimaryField.FieldJson}}")
		{{- if .AutoCreateResource }}
    userID := utils.GetUserID(c)
        {{- end }}
	err := {{.Abbreviation}}Service.Delete{{.StructName}}({{.PrimaryField.FieldJson}} {{- if .AutoCreateResource -}},userID{{- end -}})
	if err != nil {
        global.GVA_LOG.Error("thất bại!", zap.Error(err))
		response.FailWithMessage("thất bại:" + err.Error(), c)
		return
	}
	response.OkWithMessage("thành công", c)
}

// Delete{{.StructName}}ByIds 批量删除{{.Description}}
// @Tags {{.StructName}}
// @Summary 批量删除{{.Description}}
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /{{.Abbreviation}}/delete{{.StructName}}ByIds [delete]
func ({{.Abbreviation}}Api *{{.StructName}}Api) Delete{{.StructName}}ByIds(c *gin.Context) {
	{{.PrimaryField.FieldJson}}s := c.QueryArray("{{.PrimaryField.FieldJson}}s[]")
    	{{- if .AutoCreateResource }}
    userID := utils.GetUserID(c)
        {{- end }}
	err := {{.Abbreviation}}Service.Delete{{.StructName}}ByIds({{.PrimaryField.FieldJson}}s{{- if .AutoCreateResource }},userID{{- end }})
	if err != nil {
        global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:" + err.Error(), c)
		return
	}
	response.OkWithMessage("Thành công", c)
}

// Update{{.StructName}} 更新{{.Description}}
// @Tags {{.StructName}}
// @Summary 更新{{.Description}}
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body {{.Package}}.{{.StructName}} true "更新{{.Description}}"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /{{.Abbreviation}}/update{{.StructName}} [put]
func ({{.Abbreviation}}Api *{{.StructName}}Api) Update{{.StructName}}(c *gin.Context) {
	var {{.Abbreviation}} {{.Package}}.{{.StructName}}
	err := c.ShouldBindJSON(&{{.Abbreviation}})
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	    {{- if .AutoCreateResource }}
    {{.Abbreviation}}.UpdatedBy = utils.GetUserID(c)
        {{- end }}
	err = {{.Abbreviation}}Service.Update{{.StructName}}({{.Abbreviation}})
	if err != nil {
        global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:" + err.Error(), c)
		return
	}
	response.OkWithMessage("Thành công", c)
}

// Find{{.StructName}} 用id查询{{.Description}}
// @Tags {{.StructName}}
// @Summary 用id查询{{.Description}}
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query {{.Package}}.{{.StructName}} true "用id查询{{.Description}}"
// @Success 200 {object} response.Response{data={{.Package}}.{{.StructName}},msg=string} "查询成功"
// @Router /{{.Abbreviation}}/find{{.StructName}} [get]
func ({{.Abbreviation}}Api *{{.StructName}}Api) Find{{.StructName}}(c *gin.Context) {
	{{.PrimaryField.FieldJson}} := c.Query("{{.PrimaryField.FieldJson}}")
	re{{.Abbreviation}}, err := {{.Abbreviation}}Service.Get{{.StructName}}({{.PrimaryField.FieldJson}})
	if err != nil {
        global.GVA_LOG.Error("Thất bại!", zap.Error(err))
		response.FailWithMessage("Thất bại:" + err.Error(), c)
		return
	}
	response.OkWithData(re{{.Abbreviation}}, c)
}

// Get{{.StructName}}List 分页获取{{.Description}}列表
// @Tags {{.StructName}}
// @Summary 分页获取{{.Description}}列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query {{.Package}}Req.{{.StructName}}Search true "分页获取{{.Description}}列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /{{.Abbreviation}}/get{{.StructName}}List [get]
func ({{.Abbreviation}}Api *{{.StructName}}Api) Get{{.StructName}}List(c *gin.Context) {
	var pageInfo {{.Package}}Req.{{.StructName}}Search
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := {{.Abbreviation}}Service.Get{{.StructName}}InfoList(pageInfo)
	if err != nil {
	    global.GVA_LOG.Error("Thất bại!", zap.Error(err))
        response.FailWithMessage("Thất bại:" + err.Error(), c)
        return
    }
    response.OkWithDetailed(response.PageResult{
        List:     list,
        Total:    total,
        Page:     pageInfo.Page,
        PageSize: pageInfo.PageSize,
    }, "Thành công", c)
}

{{- if .HasDataSource }}
// Get{{.StructName}}DataSource 获取{{.StructName}}的数据源
// @Tags {{.StructName}}
// @Summary 获取{{.StructName}}的数据源
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "查询成功"
// @Router /{{.Abbreviation}}/get{{.StructName}}DataSource [get]
func ({{.Abbreviation}}Api *{{.StructName}}Api) Get{{.StructName}}DataSource(c *gin.Context) {
    // 此接口为获取数据源定义的数据
    dataSource, err := {{.Abbreviation}}Service.Get{{.StructName}}DataSource()
    if err != nil {
        global.GVA_LOG.Error("Thất bại!", zap.Error(err))
   		response.FailWithMessage("Thất bại:" + err.Error(), c)
   		return
    }
   response.OkWithData(dataSource, c)
}
{{- end }}

// Get{{.StructName}}Public 不需要鉴权的{{.Description}}接口
// @Tags {{.StructName}}
// @Summary 不需要鉴权的{{.Description}}接口
// @accept application/json
// @Produce application/json
// @Param data query {{.Package}}Req.{{.StructName}}Search true "分页获取{{.Description}}列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /{{.Abbreviation}}/get{{.StructName}}Public [get]
func ({{.Abbreviation}}Api *{{.StructName}}Api) Get{{.StructName}}Public(c *gin.Context) {
    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的{{.Description}}接口信息",
    }, "Thành công", c)
}
