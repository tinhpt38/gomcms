<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule"
        @keyup.enter="onSubmit">
        <el-form-item label="Ngày tạo" prop="createdAt">
          <template #label>
            <span>
              Ngày tạo
              <el-tooltip content="Phạm vi tìm kiếm từ ngày bắt đầu (bao gồm) đến ngày kết thúc (không bao gồm)">
                <el-icon>
                  <QuestionFilled />
                </el-icon>
              </el-tooltip>
            </span>
          </template>
          <el-date-picker v-model="searchInfo.startCreatedAt" type="datetime" placeholder="Ngày bắt đầu"
            :disabled-date="time => searchInfo.endCreatedAt ? time.getTime() > searchInfo.endCreatedAt.getTime() : false" />
          —
          <el-date-picker v-model="searchInfo.endCreatedAt" type="datetime" placeholder="Ngày kết thúc"
            :disabled-date="time => searchInfo.startCreatedAt ? time.getTime() < searchInfo.startCreatedAt.getTime() : false" />
        </el-form-item>
        <el-form-item label="Tên mẫu" prop="name">
          <el-input v-model="searchInfo.name" placeholder="Điều kiện tìm kiếm" />

        </el-form-item>
        <el-form-item label="Tên bảng" prop="tableName">
          <el-input v-model="searchInfo.tableName" placeholder="Điều kiện tìm kiếm" />

        </el-form-item>
        <el-form-item label="ID mẫu" prop="templateID">
          <el-input v-model="searchInfo.templateID" placeholder="Điều kiện tìm kiếm" />

        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">Tìm kiếm</el-button>
          <el-button icon="refresh" @click="onReset">Đặt lại</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog">Thêm mới</el-button>

        <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length"
          @click="onDelete">Xóa</el-button>
      </div>
      <el-table ref="multipleTable" style="width: 100%" tooltip-effect="dark" :data="tableData" row-key="ID"
        @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="Ngày" width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="Cơ sở dữ liệu" width="120">
          <template #default="scope">
            <span>{{ scope.row.dbName || "Cơ sở dữ liệu GVA" }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="ID mẫu" prop="templateID" width="120" />
        <el-table-column align="left" label="Tên mẫu" prop="name" width="120" />
        <el-table-column align="left" label="Tên bảng" prop="tableName" width="120" />
        <el-table-column align="left" label="Thông tin mẫu" prop="templateInfo" min-width="120" />
        <el-table-column align="left" label="Thao tác" min-width="120">
          <template #default="scope">
            <el-button type="primary" link icon="edit" class="table-button"
              @click="updateSysExportTemplateFunc(scope.row)">Thay đổi</el-button>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">Xóa</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination layout="total, sizes, prev, pager, next, jumper" :current-page="page" :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]" :total="total" @current-change="handleCurrentChange"
          @size-change="handleSizeChange" />
      </div>
    </div>
    <el-drawer v-model="dialogFormVisible" size="60%" :before-close="closeDialog"
      :title="type === 'create' ? 'Thêm mới' : 'Chỉnh sửa'" :show-close="false" destroy-on-close>

      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? 'Thêm mới' : 'Chỉnh sửa' }}</span>
          <div>
            <el-button @click="closeDialog">Hủy</el-button>
            <el-button type="primary" @click="enterDialog">Đồng ý</el-button>
          </div>
        </div>
      </template>

      <el-form ref="elFormRef" :model="formData" label-position="right" :rules="rule" label-width="100px">

        <el-form-item label="Cơ sở dữ liệu kinh doanh" prop="dbName">
          <template #label>
            <el-tooltip
              content="Chú ý: Cần cấu hình nhiều cơ sở dữ liệu trước đó trong db-list. Nếu không cấu hình, cần cấu hình và khởi động lại dịch vụ để sử dụng. Nếu không thể chọn, hãy đặt disabled: false trong config.yaml để chọn cơ sở dữ liệu đích cho nhập và xuất."
              placement="bottom" effect="light">
              <div> Cơ sở dữ liệu kinh doanh <el-icon>
                  <QuestionFilled />
                </el-icon> </div>
            </el-tooltip>
          </template>
          <el-select v-model="formData.dbName" clearable @change="dbNameChange" placeholder="Chọn cơ sở dữ liệu kinh doanh">
            <el-option v-for="item in dbList" :key="item.aliasName" :value="item.aliasName" :label="item.aliasName"
              :disabled="item.disable">
              <div>
                <span>{{ item.aliasName }}</span>
                <span style="float:right;color:#8492a6;font-size:13px">{{ item.dbName }}</span>
              </div>
            </el-option>
          </el-select>
        </el-form-item>

        <el-form-item label="Tên mẫu:" prop="name">
          <el-input v-model="formData.name" :clearable="true" placeholder="Nhập tên mẫu" />
        </el-form-item>
        <el-form-item label="Tên bảng:" clearable prop="tableName">
          <!--          <el-input
            v-model="formData.tableName"
            :clearable="true"
            placeholder="Nhập tên bảng cần xuất"
          />-->
          <div class="w-full flex gap-4">
            <el-select v-model="formData.tableName" class="flex-1" filterable placeholder="Chọn bảng">
              <el-option v-for="item in tableOptions" :key="item.tableName" :label="item.tableName"
                :value="item.tableName" />
            </el-select>

            <el-button type="primary" @click="getColumnFunc">Tự động tạo mẫu</el-button>
          </div>

        </el-form-item>
        <el-form-item label="ID mẫu:" prop="templateID">
          <el-input v-model="formData.templateID" :clearable="true" placeholder="ID mẫu là thuộc tính gắn vào thành phần frontend" />
        </el-form-item>

        <el-form-item label="Điều kiện liên kết:">
          <div v-for="(join, key) in formData.joinTemplate" :key="key" class="flex gap-4 w-full mb-2">
            <el-select v-model="join.joins" placeholder="Chọn phương thức liên kết">
              <el-option label="LEFT JOIN" value="LEFT JOIN" />
              <el-option label="INNER JOIN" value="INNER JOIN" />
              <el-option label="RIGHT JOIN" value="RIGHT JOIN" />
            </el-select>
            <el-input v-model="join.table" placeholder="Nhập bảng liên kết" />
            <el-input v-model="join.on" placeholder="Điều kiện liên kết table1.a = table2.b" />
            <el-button type="danger" icon="delete" @click="() => formData.joinTemplate.splice(key, 1)">Xóa</el-button>
          </div>
          <div class="flex justify-end w-full">
            <el-button type="primary" icon="plus" @click="addJoin">Thêm điều kiện</el-button>
          </div>
        </el-form-item>

        <el-form-item label="Thông tin mẫu:" prop="templateInfo">
          <el-input v-model="formData.templateInfo" type="textarea" :rows="12" :clearable="true"
            :placeholder="templatePlaceholder" />
        </el-form-item>
        <el-form-item label="Số lượng xuất mặc định:">
          <el-input-number v-model="formData.limit" :step="1" :step-strictly="true" :precision="0" />
        </el-form-item>
        <el-form-item label="Điều kiện sắp xếp mặc định:">
          <el-input v-model="formData.order" placeholder="Ví dụ: id desc" />
        </el-form-item>
        <el-form-item label="Điều kiện xuất:">
          <div v-for="(condition, key) in formData.conditions" :key="key" class="flex gap-4 w-full mb-2">
            <el-input v-model="condition.from" placeholder="Lấy khóa json từ điều kiện truy vấn" />
            <el-input v-model="condition.column" placeholder="Cột tương ứng trong bảng" />
            <el-select v-model="condition.operator" placeholder="Chọn điều kiện truy vấn">
              <el-option v-for="item in typeSearchOptions" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
            <el-button type="danger" icon="delete" @click="() => formData.conditions.splice(key, 1)">Xóa</el-button>
          </div>
          <div class="flex justify-end w-full">
            <el-button type="primary" icon="plus" @click="addCondition">Thêm điều kiện</el-button>
          </div>
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createSysExportTemplate,
  deleteSysExportTemplate,
  deleteSysExportTemplateByIds,
  updateSysExportTemplate,
  findSysExportTemplate,
  getSysExportTemplateList
} from '@/api/exportTemplate.js'

// Import formatting utilities as needed
import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import WarningBar from '@/components/warningBar/warningBar.vue'
import { getDB, getTable, getColumn } from '@/api/autoCode'

defineOptions({
  name: 'ExportTemplate'
})

const templatePlaceholder = `Thông tin mẫu: key đại diện cho tên cột trong cơ sở dữ liệu (trong chế độ liên kết, hãy viết dưới dạng table.column), value đại diện cho tên cột trong tệp excel xuất ra. Ví dụ: key là từ khóa hoặc hàm của cơ sở dữ liệu, hãy xử lý theo cách xử lý từ khóa hiện tại (ví dụ: mysql), như sau:
{
  "table_column1":"Cột thứ nhất",
  "table_column3":"Cột thứ ba",
  "table_column4":"Cột thứ tư",
  "\`rows\`":"Tôi là từ khóa hoặc hàm của cơ sở dữ liệu",
}
Nếu thêm JOINS, key xuất phải là {table_name1.table_column1:"Cột thứ nhất",table_name2.table_column2:"Cột thứ hai"}
Nếu có tên cột trùng lặp, định dạng xuất phải là {table_name1.table_column1 as key:"Cột thứ nhất",table_name2.table_column2 as key2:"Cột thứ hai"}
Không hỗ trợ nhập trong chế độ JOINS
`

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
  name: '',
  tableName: '',
  dbName: '',
  templateID: '',
  templateInfo: '',
  limit: 0,
  order: '',
  conditions: [],
  joinTemplate: []
})

const typeSearchOptions = ref([
  {
    label: '=',
    value: '='
  },
  {
    label: '<>',
    value: '<>'
  },
  {
    label: '>',
    value: '>'
  },
  {
    label: '<',
    value: '<'
  },
  {
    label: 'LIKE',
    value: 'LIKE'
  },
  {
    label: 'BETWEEN',
    value: 'BETWEEN'
  },
  {
    label: 'NOT BETWEEN',
    value: 'NOT BETWEEN'
  }
])

const addCondition = () => {
  formData.value.conditions.push({
    from: '',
    column: '',
    operator: ''
  })
}

const addJoin = () => {
  formData.value.joinTemplate.push({
    joins: 'LEFT JOIN',
    table: '',
    on: ''
  })
}

// 验证规则
const rule = reactive({
  name: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  {
    whitespace: true,
    message: '不能只输入空格',
    trigger: ['input', 'blur'],
  }
  ],
  tableName: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  {
    whitespace: true,
    message: '不能只输入空格',
    trigger: ['input', 'blur'],
  }
  ],
  templateID: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  {
    whitespace: true,
    message: '不能只输入空格',
    trigger: ['input', 'blur'],
  }
  ],
  templateInfo: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  {
    whitespace: true,
    message: '不能只输入空格',
    trigger: ['input', 'blur'],
  }
  ],
})

const searchRule = reactive({
  createdAt: [
    {
      validator: (rule, value, callback) => {
        if (searchInfo.value.startCreatedAt && !searchInfo.value.endCreatedAt) {
          callback(new Error('请填写结束日期'))
        } else if (!searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt) {
          callback(new Error('请填写开始日期'))
        } else if (searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt && (searchInfo.value.startCreatedAt.getTime() === searchInfo.value.endCreatedAt.getTime() || searchInfo.value.startCreatedAt.getTime() > searchInfo.value.endCreatedAt.getTime())) {
          callback(new Error('开始日期应当早于结束日期'))
        } else {
          callback()
        }
      }, trigger: 'change'
    }
  ],
})

const elFormRef = ref()
const elSearchFormRef = ref()

// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

const dbList = ref([])
const tableOptions = ref([])

const getDbFunc = async () => {
  const res = await getDB()
  if (res.code === 0) {
    dbList.value = res.data.dbList
  }
}

getDbFunc()

const dbNameChange = () => {
  formData.value.tableName = ''
  formData.value.templateInfo = ''
  getTableFunc()
}

const getTableFunc = async () => {
  const res = await getTable({ businessDB: formData.value.dbName })
  if (res.code === 0) {
    tableOptions.value = res.data.tables
  }
  formData.value.tableName = ''
}
getTableFunc()

const getColumnFunc = async () => {
  if (!formData.value.tableName) {
    ElMessage({
      type: 'error',
      message: '请先选择业务库及选择表后再进行操作'
    })
    return
  }
  formData.value.templateInfo = ""
  const res = await getColumn({
    businessDB: formData.value.dbName,
    tableName: formData.value.tableName
  })
  if (res.code === 0) {
    // 把返回值的data.columns做尊换，制作一组JSON数据，columnName做key，columnComment做value
    const templateInfo = {}
    res.data.columns.forEach(item => {
      templateInfo[item.columnName] = item.columnComment || item.columnName
    })
    formData.value.templateInfo = JSON.stringify(templateInfo, null, 2)
  }

}

// 重置
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

// 搜索
const onSubmit = () => {
  elSearchFormRef.value?.validate(async (valid) => {
    if (!valid) return
    page.value = 1
    pageSize.value = 10
    getTableData()
  })
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

// 修改页面容量
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 查询
const getTableData = async () => {
  const table = await getSysExportTemplateList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async () => {
}

// 获取需要的字典 可能为空 按需保留
setOptions()

// 多选数据
const multipleSelection = ref([])
// 多选
const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

// 删除行
const deleteRow = (row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    deleteSysExportTemplateFunc(row)
  })
}

// 多选删除
const onDelete = async () => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const ids = []
    if (multipleSelection.value.length === 0) {
      ElMessage({
        type: 'warning',
        message: '请选择要删除的数据'
      })
      return
    }
    multipleSelection.value &&
      multipleSelection.value.map(item => {
        ids.push(item.ID)
      })
    const res = await deleteSysExportTemplateByIds({ ids })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '删除成功'
      })
      if (tableData.value.length === ids.length && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  })
}

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')

// 更新行
const updateSysExportTemplateFunc = async (row) => {
  const res = await findSysExportTemplate({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data.resysExportTemplate
    if (!formData.value.conditions) {
      formData.value.conditions = []
    }
    if (!formData.value.joinTemplate) {
      formData.value.joinTemplate = []
    }
    dialogFormVisible.value = true
  }
}

// 删除行
const deleteSysExportTemplateFunc = async (row) => {
  const res = await deleteSysExportTemplate({ ID: row.ID })
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '删除成功'
    })
    if (tableData.value.length === 1 && page.value > 1) {
      page.value--
    }
    getTableData()
  }
}

// 弹窗控制标记
const dialogFormVisible = ref(false)

// 打开弹窗
const openDialog = () => {
  type.value = 'create'
  dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = {
    name: '',
    tableName: '',
    templateID: '',
    templateInfo: '',
    limit: 0,
    order: '',
    conditions: [],
    joinTemplate: [],
  }
}
// 弹窗确定
const enterDialog = async () => {
  // 判断 formData.templateInfo 是否为标准json格式 如果不是标准json 则辅助调整
  try {
    JSON.parse(formData.value.templateInfo)
  } catch (error) {
    ElMessage({
      type: 'error',
      message: '模板信息格式不正确，请检查'
    })
    return
  }

  const reqData = JSON.parse(JSON.stringify(formData.value))
  for (let i = 0; i < reqData.conditions.length; i++) {
    if (!reqData.conditions[i].from || !reqData.conditions[i].column || !reqData.conditions[i].operator) {
      ElMessage({
        type: 'error',
        message: '请填写完整的导出条件'
      })
      return
    }
    reqData.conditions[i].templateID = reqData.templateID
  }

  for (let i = 0; i < reqData.joinTemplate.length; i++) {
    if (!reqData.joinTemplate[i].joins || !reqData.joinTemplate[i].on) {
      ElMessage({
        type: 'error',
        message: '请填写完整的关联'
      })
      return
    }
    reqData.joinTemplate[i].templateID = reqData.templateID
  }

  elFormRef.value?.validate(async (valid) => {
    if (!valid) return
    let res
    switch (type.value) {
      case 'create':
        res = await createSysExportTemplate(reqData)
        break
      case 'update':
        res = await updateSysExportTemplate(reqData)
        break
      default:
        res = await createSysExportTemplate(reqData)
        break
    }
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: 'Tạo/cập nhật thành công'
      })
      closeDialog()
      getTableData()
    }
  })
}

</script>

<style></style>
