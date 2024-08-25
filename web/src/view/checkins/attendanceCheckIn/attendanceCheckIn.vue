<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule"
        @keyup.enter="onSubmit">
        <el-form-item label="创建日期" prop="createdAt">
          <template #label>
            <span>
              创建日期
              <el-tooltip content="搜索范围是开始日期（包含）至结束日期（不包含）">
                <el-icon>
                  <QuestionFilled />
                </el-icon>
              </el-tooltip>
            </span>
          </template>
          <el-date-picker v-model="searchInfo.startCreatedAt" type="datetime" placeholder="开始日期"
            :disabled-date="time => searchInfo.endCreatedAt ? time.getTime() > searchInfo.endCreatedAt.getTime() : false"></el-date-picker>
          —
          <el-date-picker v-model="searchInfo.endCreatedAt" type="datetime" placeholder="结束日期"
            :disabled-date="time => searchInfo.startCreatedAt ? time.getTime() < searchInfo.startCreatedAt.getTime() : false"></el-date-picker>
        </el-form-item>


        <template v-if="showAllQuery">
          <!-- 将需要控制显示状态的查询条件添加到此范围内 -->
        </template>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
          <el-button link type="primary" icon="arrow-down" @click="showAllQuery = true"
            v-if="!showAllQuery">展开</el-button>
          <el-button link type="primary" icon="arrow-up" @click="showAllQuery = false" v-else>收起</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog">新增</el-button>
        <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length"
          @click="onDelete">删除</el-button>
      </div>
      <el-table ref="multipleTable" style="width: 100%" tooltip-effect="dark" :data="tableData" row-key="ID"
        @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />

        <el-table-column align="left" label="日期" prop="createdAt" width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>

        <el-table-column align="left" label="Ngày giờ" prop="checkinDate" width="180">
          <template #default="scope">{{ formatDate(scope.row.checkinDate) }}</template>
        </el-table-column>
        <el-table-column align="left" label="Lớp điểm danh" prop="attendanceId" width="120" />
        <el-table-column align="left" label="Thành viên" prop="partpaticipantId" width="120" />
        <el-table-column align="left" label="Khu vực" prop="areaId" width="120" />
        <el-table-column align="left" label="Nhóm" prop="groupId" width="120" />
        <el-table-column align="left" label="Điều kiện" prop="conditionId" width="120" />
        <el-table-column align="left" label="Kinh độ" prop="lattidue" width="120" />
        <el-table-column align="left" label="IP" prop="iP" width="120" />
        <el-table-column align="left" label="Vĩ độ" prop="longtidue" width="120" />
        <el-table-column label="Agent" prop="agent" width="200">
          <template #default="scope">
            [JSON]
          </template>
        </el-table-column>
        <el-table-column align="left" label="Trình duyệt" prop="browser" width="120" />
        <el-table-column align="left" label="操作" fixed="right" min-width="240">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)"><el-icon
                style="margin-right: 5px">
                <InfoFilled />
              </el-icon>查看详情</el-button>
            <el-button type="primary" link icon="edit" class="table-button"
              @click="updateAttendanceCheckInFunc(scope.row)">变更</el-button>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
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
          <span class="text-lg">{{ type === 'create' ? '添加' : '修改' }}</span>
          <div>
            <el-button type="primary" @click="enterDialog">Đồng ý</el-button>
            <el-button @click="closeDialog">Huỷ</el-button>
          </div>
        </div>
      </template>

      <el-form :model="formData" label-position="top" ref="elFormRef" :rules="rule" label-width="80px">
        <el-form-item label="Ngày giờ:" prop="checkinDate">
          <el-date-picker v-model="formData.checkinDate" type="date" style="width:100%" placeholder="选择日期"
            :clearable="true" />
        </el-form-item>
        <el-form-item label="Lớp điểm danh:" prop="attendanceId">
          <el-input v-model.number="formData.attendanceId" :clearable="true" placeholder="请输入Lớp điểm danh" />
        </el-form-item>
        <el-form-item label="Thành viên:" prop="partpaticipantId">
          <el-input v-model.number="formData.partpaticipantId" :clearable="true" placeholder="请输入Thành viên" />
        </el-form-item>
        <el-form-item label="Khu vực:" prop="areaId">
          <el-input v-model.number="formData.areaId" :clearable="true" placeholder="请输入Khu vực" />
        </el-form-item>
        <el-form-item label="Nhóm:" prop="groupId">
          <el-input v-model.number="formData.groupId" :clearable="true" placeholder="请输入Nhóm" />
        </el-form-item>
        <el-form-item label="Điều kiện:" prop="conditionId">
          <el-input v-model.number="formData.conditionId" :clearable="true" placeholder="请输入Điều kiện" />
        </el-form-item>
        <el-form-item label="Kinh độ:" prop="lattidue">
          <el-input-number v-model="formData.lattidue" style="width:100%" :precision="2" :clearable="true" />
        </el-form-item>
        <el-form-item label="IP:" prop="iP">
          <el-input v-model.number="formData.iP" :clearable="true" placeholder="请输入IP" />
        </el-form-item>
        <el-form-item label="Vĩ độ:" prop="longtidue">
          <el-input-number v-model="formData.longtidue" style="width:100%" :precision="2" :clearable="true" />
        </el-form-item>
        <el-form-item label="Agent:" prop="agent">
          // 此字段为json结构，可以前端自行控制展示和数据绑定模式 需绑定json的key为 formData.agent 后端会按照json的类型进行存取
          {{ formData.agent }}
        </el-form-item>
        <el-form-item label="Trình duyệt:" prop="browser">
          <el-input v-model="formData.browser" :clearable="true" placeholder="请输入Trình duyệt" />
        </el-form-item>
      </el-form>
    </el-drawer>

    <el-drawer destroy-on-close size="800" v-model="detailShow" :show-close="true" :before-close="closeDetailShow">
      <el-descriptions column="1" border>
        <el-descriptions-item label="Ngày giờ">
          {{ detailFrom.checkinDate }}
        </el-descriptions-item>
        <el-descriptions-item label="Lớp điểm danh">
          {{ detailFrom.attendanceId }}
        </el-descriptions-item>
        <el-descriptions-item label="Thành viên">
          {{ detailFrom.partpaticipantId }}
        </el-descriptions-item>
        <el-descriptions-item label="Khu vực">
          {{ detailFrom.areaId }}
        </el-descriptions-item>
        <el-descriptions-item label="Nhóm">
          {{ detailFrom.groupId }}
        </el-descriptions-item>
        <el-descriptions-item label="Điều kiện">
          {{ detailFrom.conditionId }}
        </el-descriptions-item>
        <el-descriptions-item label="Kinh độ">
          {{ detailFrom.lattidue }}
        </el-descriptions-item>
        <el-descriptions-item label="IP">
          {{ detailFrom.iP }}
        </el-descriptions-item>
        <el-descriptions-item label="Vĩ độ">
          {{ detailFrom.longtidue }}
        </el-descriptions-item>
        <el-descriptions-item label="Agent">
          {{ detailFrom.agent }}
        </el-descriptions-item>
        <el-descriptions-item label="Trình duyệt">
          {{ detailFrom.browser }}
        </el-descriptions-item>
      </el-descriptions>
    </el-drawer>

  </div>
</template>

<script setup>
import {
  createAttendanceCheckIn,
  deleteAttendanceCheckIn,
  deleteAttendanceCheckInByIds,
  updateAttendanceCheckIn,
  findAttendanceCheckIn,
  getAttendanceCheckInList
} from '@/api/checkins/attendanceCheckIn'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict, filterDataSource, returnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

defineOptions({
  name: 'AttendanceCheckIn'
})

// 控制更多查询条件显示/隐藏状态
const showAllQuery = ref(false)

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
  checkinDate: new Date(),
  attendanceId: undefined,
  partpaticipantId: undefined,
  areaId: undefined,
  groupId: undefined,
  conditionId: undefined,
  lattidue: 0,
  iP: undefined,
  longtidue: 0,
  agent: {},
  browser: '',
})



// 验证规则
const rule = reactive({
  checkinDate: [{
    required: true,
    message: 'Ngày giờ không được để trống',
    trigger: ['input', 'blur'],
  },
  ],
  attendanceId: [{
    required: true,
    message: 'Lớp điểm danh không được để trống',
    trigger: ['input', 'blur'],
  },
  ],
  partpaticipantId: [{
    required: true,
    message: 'Thành viên không được để trống',
    trigger: ['input', 'blur'],
  },
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
  const table = await getAttendanceCheckInList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
    deleteAttendanceCheckInFunc(row)
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
    const res = await deleteAttendanceCheckInByIds({ IDs })
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
const updateAttendanceCheckInFunc = async (row) => {
  const res = await findAttendanceCheckIn({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data
    dialogFormVisible.value = true
  }
}


// 删除行
const deleteAttendanceCheckInFunc = async (row) => {
  const res = await deleteAttendanceCheckIn({ ID: row.ID })
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
    checkinDate: new Date(),
    attendanceId: undefined,
    partpaticipantId: undefined,
    areaId: undefined,
    groupId: undefined,
    conditionId: undefined,
    lattidue: 0,
    iP: undefined,
    longtidue: 0,
    agent: {},
    browser: '',
  }
}
// 弹窗确定
const enterDialog = async () => {
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return
    let res
    switch (type.value) {
      case 'create':
        res = await createAttendanceCheckIn(formData.value)
        break
      case 'update':
        res = await updateAttendanceCheckIn(formData.value)
        break
      default:
        res = await createAttendanceCheckIn(formData.value)
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
  const res = await findAttendanceCheckIn({ ID: row.ID })
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
