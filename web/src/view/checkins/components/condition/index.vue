<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog">
          Thêm mới
        </el-button>
        <!-- <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length"
                    @click="onDelete">Xóa</el-button> -->
      </div>
      <el-table ref="multipleTable" style="width: 100%" tooltip-effect="dark" :data="tableData" row-key="ID" border
        @selection-change="handleSelectionChange">
        <!-- <el-table-column type="selection" width="55" /> -->

        <!-- <el-table-column align="left" label="Ngày" prop="createdAt" width="180">
                    <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
</el-table-column> -->
        <el-table-column align="left" label="Nhóm" prop="group.name" min-width="120" />
        <el-table-column align="left" label="Khu vực" prop="area.area.name" min-width="300">
          <template #default="scope">
            {{ scope.row.area?.area?.name }} <span v-if="scope.row.area?.radius">(Bán kính {{ scope.row.area?.radius
              }})</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="Bắt đầu" prop="startAt" min-width="180">
          <template #default="scope">
            {{ formatDateTime(scope.row.startAt) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="Kết thúc" prop="endAt" min-width="180">
          <template #default="scope">
            {{ formatDateTime(scope.row.endAt) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="Hành động" fixed="right" min-width="240">
          <template #default="scope">
            <!-- <el-button type="primary" link class="table-button" @click="getDetails(scope.row)"><el-icon
                                style="margin-right: 5px">
                                <InfoFilled />
                            </el-icon>Xem chi tiết</el-button> -->
            <el-button type="primary" link icon="edit" class="table-button" @click="updateConditionFunc(scope.row)">
              Chỉnh sửa
            </el-button>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">
              Xóa
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination layout="total, sizes, prev, pager, next, jumper" :current-page="page" :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]" :total="total" @current-change="handleCurrentChange"
          @size-change="handleSizeChange" />
      </div>
    </div>
    <el-drawer v-model="dialogFormVisible" destroy-on-close size="800" :show-close="false" :before-close="closeDialog">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? 'Thêm mới' : 'Chỉnh sửa' }}</span>
          <div>
            <el-button type="primary" @click="enterDialog">
              Đồng ý
            </el-button>
            <el-button @click="closeDialog">
              Hủy
            </el-button>
          </div>
        </div>
      </template>

      <el-form ref="elFormRef" :model="formData" label-position="top" :rules="searchRule" label-width="80px">
        <el-form-item label="Nhóm:" prop="groupId">
          <el-select v-model="formData.groupId" placeholder="Chọn nhóm" clearable filterable>
            <el-option v-for="item in groupOptions" :key="item.ID" :label="item.name" :value="item.ID" />
          </el-select>
        </el-form-item>
        <el-form-item label="Khu vực:" prop="areaId">
          <!-- <el-input v-model.number="formData.areaId" :clearable="true" placeholder="Nhập Khu vực" /> -->
          <el-select v-model="formData.areaId" placeholder="Chọn khu vực" clearable filterable>
            <el-option v-for="item in areaOptions" :key="item.ID" :label="item.name" :value="item.ID" />
          </el-select>
        </el-form-item>
        <el-form-item label="Bắt đầu:" prop="startAt">
          <el-date-picker v-model="formData.startAt" type="datetime" style="width:100%" placeholder="Chọn ngày giờ"
            :clearable="true" />
        </el-form-item>
        <el-form-item label="Kết thúc:" prop="endAt">
          <el-date-picker v-model="formData.endAt" type="datetime" style="width:100%" placeholder="Chọn ngày giờ"
            :clearable="true" />
        </el-form-item>
      </el-form>
    </el-drawer>

    <el-drawer v-model="detailShow" destroy-on-close size="800" :show-close="true" :before-close="closeDetailShow">
      <el-descriptions column="1" border>
        <el-descriptions-item label="Nhóm">
          {{ detailFrom.groupId }}
        </el-descriptions-item>
        <el-descriptions-item label="Khu vực">
          {{ detailFrom.areaId }}
        </el-descriptions-item>
        <el-descriptions-item label="Bắt đầu">
          {{ detailFrom.startAt }}
        </el-descriptions-item>
        <el-descriptions-item label="Kết thúc">
          {{ detailFrom.endAt }}
        </el-descriptions-item>
      </el-descriptions>
    </el-drawer>
  </div>
</template>


<script setup>
import {
  createCondition,
  deleteCondition,
  deleteConditionByIds,
  updateCondition,
  findCondition,
  getConditionList
} from '@/api/checkins/condition'

import {
  findAttendanceArea
} from '@/api/checkins/attendance'

import {
  getGroupList
} from '@/api/checkins/group'


import { formatDate, formatDateTime, } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import moment from 'moment'


const props = defineProps({
  acId: {
    type: Number,
    required: true
  }
})
const formData = ref({
  attendanceId: props.acId,
  groupId: null,
  areaId: null,
  startAt: null,
  endAt: null,
})

const elFormRef = ref()
const elSearchFormRef = ref()

// =========== Bảng điều khiển ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})



const areaOptions = ref([])
const getAreaListData = async () => {
  const table = await findAttendanceArea({ id: props.acId })
  if (table.code === 0) {
    areaOptions.value = table.data.map(item => {
      return {
        ID: item.ID,
        name: item.area?.name
      }
    })
  }

}
getAreaListData();

const groupOptions = ref([])
const getGroupOptions = async () => {
  searchInfo.value.attendanceId = props.acId
  const table = await getGroupList({ page: 1, pageSize: 99999, ...searchInfo.value })
  if (table.code === 0) {
    groupOptions.value = table.data.list
  }
  console.log('groupOptions', groupOptions.value)
}
getGroupOptions();


const searchRule = reactive({
  startAt: [
    {
      validator: (rule, value, callback) => {
        if (formData.value.startAt && !formData.value.endAt) {
          callback(new Error('Vui lòng nhập ngày kết thúc'))
        } else if (!formData.value.startAt && formData.value.endAt) {
          callback(new Error('Vui lòng nhập ngày bắt đầu'))
        } else if (moment(formData.value.startAt).isAfter(moment(formData.value.endAt))) {
          callback(new Error('Ngày bắt đầu phải trước ngày kết thúc'))
        } else {
          callback()
        }
      }, trigger: 'change'
    }
  ],
  endAt: [
    {
      validator: (rule, value, callback) => {
        if (formData.value.startAt && !formData.value.endAt) {
          callback(new Error('Vui lòng nhập ngày kết thúc'))
        } else if (!formData.value.startAt && formData.value.endAt) {
          callback(new Error('Vui lòng nhập ngày bắt đầu'))
        } else if (moment(formData.value.startAt).isAfter(moment(formData.value.endAt))) {
          callback(new Error('Ngày bắt đầu phải trước ngày kết thúc'))
        } else {
          callback()
        }
      }, trigger: 'change'
    }
  ],
})


// Đặt lại
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

// Tìm kiếm
const onSubmit = () => {
  elSearchFormRef.value?.validate(async (valid) => {
    if (!valid) return
    page.value = 1
    pageSize.value = 10
    getTableData()
  })
}

// Phân trang
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

// Thay đổi kích thước trang
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// Tìm kiếm
const getTableData = async () => {
  const table = await getConditionList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
  console.log("condition data", tableData.value)
}

getTableData()

// ============== Kết thúc bảng điều khiển ===============

// Lấy các từ điển cần thiết có thể trống, tùy theo nhu cầu
const setOptions = async () => {
}

// Lấy các từ điển cần thiết có thể trống, tùy theo nhu cầu
setOptions()


// Dữ liệu đa chọn
const multipleSelection = ref([])

const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

const deleteRow = (row) => {
  ElMessageBox.confirm('Bạn có chắc muốn xóa không?', 'Cảnh báo', {
    confirmButtonText: 'Đồng ý',
    cancelButtonText: 'Hủy',
    type: 'warning'
  }).then(() => {
    deleteConditionFunc(row)
  })
}

const onDelete = async () => {
  ElMessageBox.confirm('Bạn có chắc muốn xóa không?', 'Cảnh báo', {
    confirmButtonText: 'Đồng ý',
    cancelButtonText: 'Hủy',
    type: 'warning'
  }).then(async () => {
    const IDs = []
    if (multipleSelection.value.length === 0) {
      ElMessage({
        type: 'warning',
        message: 'Vui lòng chọn dữ liệu để xóa'
      })
      return
    }
    multipleSelection.value &&
      multipleSelection.value.map(item => {
        IDs.push(item.ID)
      })
    const res = await deleteConditionByIds({ IDs })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: 'Xóa thành công'
      })
      if (tableData.value.length === IDs.length && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  })
}

// Đánh dấu hành vi (thêm hoặc sửa) trong cửa sổ pop-up
const type = ref('')


const updateConditionFunc = async (row) => {
  const res = await findCondition({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data
    dialogFormVisible.value = true
  }
}



const deleteConditionFunc = async (row) => {
  const res = await deleteCondition({ ID: row.ID })
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: 'Xoá thành công'
    })
    if (tableData.value.length === 1 && page.value > 1) {
      page.value--
    }
    getTableData()
  }
}


const dialogFormVisible = ref(false)


const openDialog = () => {
  type.value = 'create'
  dialogFormVisible.value = true
}

const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = {
    attendanceId: props.acId,
    groupId: null,
    areaId: null,
    startAt: null,
    endAt: null,
  }
}

const enterDialog = async () => {
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return
    let res
    formData.value.attendanceId = props.acId * 1
    switch (type.value) {
      case 'create':
        res = await createCondition(formData.value)
        break
      case 'update':
        res = await updateCondition(formData.value)
        break
      default:
        res = await createCondition(formData.value)
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


const detailShow = ref(false)



const openDetailShow = () => {
  detailShow.value = true
}



const getDetails = async (row) => {

  const res = await findCondition({ ID: row.ID })
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

<style scoped>
/* Add your component styles here */
</style>