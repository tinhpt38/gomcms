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
            :disabled-date="time => searchInfo.endCreatedAt ? time.getTime() > searchInfo.endCreatedAt.getTime() : false"></el-date-picker>
          —
          <el-date-picker v-model="searchInfo.endCreatedAt" type="datetime" placeholder="Ngày kết thúc"
            :disabled-date="time => searchInfo.startCreatedAt ? time.getTime() < searchInfo.startCreatedAt.getTime() : false"></el-date-picker>
        </el-form-item>

        <el-form-item label="Attendance Class" prop="attendanceClassId">

          <el-input v-model.number="searchInfo.attendanceClassId" placeholder="Điều kiện tìm kiếm" />

        </el-form-item>

        <template v-if="showAllQuery">
          <!-- Thêm các điều kiện tìm kiếm cần điều khiển hiển thị vào đây -->
        </template>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">Tìm kiếm</el-button>
          <el-button icon="refresh" @click="onReset">Đặt lại</el-button>
          <el-button link type="primary" icon="arrow-down" @click="showAllQuery = true" v-if="!showAllQuery">Mở
            rộng</el-button>
          <el-button link type="primary" icon="arrow-up" @click="showAllQuery = false" v-else>Thu gọn</el-button>
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

        <el-table-column align="left" label="Ngày" prop="createdAt" width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>

        <el-table-column align="left" label="Tên nhóm" prop="name" width="120" />
        <el-table-column align="left" label="Attendance Class" prop="attendanceClassId" width="120">
          <template #default="scope">

            <span>{{ filterDataSource(dataSource.attendanceClassId, scope.row.attendanceClassId) }}</span>

          </template>
        </el-table-column>
        <el-table-column align="left" label="Hành động" fixed="right" min-width="240">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)"><el-icon
                style="margin-right: 5px">
                <InfoFilled />
              </el-icon>Xem chi tiết</el-button>
            <el-button type="primary" link icon="edit" class="table-button" @click="updateGroupFunc(scope.row)">Chỉnh
              sửa</el-button>
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
    <el-drawer destroy-on-close size="800" v-model="dialogFormVisible" :show-close="false" :before-close="closeDialog">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? 'Thêm mới' : 'Chỉnh sửa' }}</span>
          <div>
            <el-button type="primary" @click="enterDialog">Đồng ý</el-button>
            <el-button @click="closeDialog">Huỷ</el-button>
          </div>
        </div>
      </template>

      <el-form :model="formData" label-position="top" ref="elFormRef" :rules="rule" label-width="80px">
        <el-form-item label="Tên nhóm:" prop="name">
          <el-input v-model="formData.name" :clearable="true" placeholder="Vui lòng nhập Tên nhóm" />
        </el-form-item>
        <el-form-item label="Attendance Class:" prop="attendanceClassId">
          <el-select v-model="formData.attendanceClassId" placeholder="Vui lòng chọn Attendance Class"
            style="width:100%" :clearable="true">
            <el-option v-for="(item, key) in dataSource.attendanceClassId" :key="key" :label="item.label"
              :value="item.value" />
          </el-select>
        </el-form-item>
      </el-form>
    </el-drawer>

    <el-drawer destroy-on-close size="800" v-model="detailShow" :show-close="true" :before-close="closeDetailShow">
      <el-descriptions column="1" border>
        <el-descriptions-item label="Tên nhóm">
          {{ detailFrom.name }}
        </el-descriptions-item>
        <el-descriptions-item label="Attendance Class">
          {{ detailFrom.attendanceClassId }}
        </el-descriptions-item>
      </el-descriptions>
    </el-drawer>

  </div>
</template>

<script setup>
import {
  getGroupDataSource,
  createGroup,
  deleteGroup,
  deleteGroupByIds,
  updateGroup,
  findGroup,
  getGroupList
} from '@/api/checkins/group'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict, filterDataSource, returnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

defineOptions({
  name: 'Group'
})

// 控制更多查询条件显示/隐藏状态
const showAllQuery = ref(false)

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
  name: '',
  attendanceClassId: undefined,
})
const dataSource = ref([])
const getDataSourceFunc = async () => {
  const res = await getGroupDataSource()
  if (res.code === 0) {
    dataSource.value = res.data
  }
}
getDataSourceFunc()



// 验证规则
const rule = reactive({
  name: [{
    required: true,
    message: 'Tên nhóm là bắt buộc',
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
  const table = await getGroupList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
    deleteGroupFunc(row)
  })
}

// 多选删除
const onDelete = async () => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const IDs = []
    if (multipleSelection.value.length === 0) {
      ElMessage({
        type: 'warning',
        message: '请选择要删除的数据'
      })
      return
    }
    multipleSelection.value &&
      multipleSelection.value.map(item => {
        IDs.push(item.ID)
      })
    const res = await deleteGroupByIds({ IDs })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '删除成功'
      })
      if (tableData.value.length === IDs.length && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  })
}

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')

// 更新行
const updateGroupFunc = async (row) => {
  const res = await findGroup({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data
    dialogFormVisible.value = true
  }
}


// 删除行
const deleteGroupFunc = async (row) => {
  const res = await deleteGroup({ ID: row.ID })
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
    attendanceClassId: undefined,
  }
}
// 弹窗确定
const enterDialog = async () => {
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return
    let res
    switch (type.value) {
      case 'create':
        res = await createGroup(formData.value)
        break
      case 'update':
        res = await updateGroup(formData.value)
        break
      default:
        res = await createGroup(formData.value)
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


const detailFrom = ref({})

// 查看详情控制标记
const detailShow = ref(false)


// 打开详情弹窗
const openDetailShow = () => {
  detailShow.value = true
}


// 打开详情
const getDetails = async (row) => {
  // 打开弹窗
  const res = await findGroup({ ID: row.ID })
  if (res.code === 0) {
    detailFrom.value = res.data
    openDetailShow()
  }
}


// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  detailFrom.value = {}
}


</script>

<style></style>
