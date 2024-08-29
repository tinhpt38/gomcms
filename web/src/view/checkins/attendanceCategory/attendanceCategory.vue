<template>
  <div>
    <div class="gva-search-box hidden">
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

        <el-form-item label="Tên phân loại" prop="name">

          <el-input v-model.number="searchInfo.name" placeholder="Điều kiện tìm kiếm" />

        </el-form-item>

        <template v-if="showAllQuery">
          <!-- Thêm các điều kiện tìm kiếm cần điều khiển hiển thị vào đây -->
        </template>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">Tìm kiếm</el-button>
          <el-button icon="refresh" @click="onReset">Làm mới</el-button>
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

        <!-- <el-table-column align="left" label="Ngày" prop="createdAt" width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
         -->
        <el-table-column align="left" label="Tên phân loại" prop="name" width="200" />
        <!-- <el-table-column align="left" label="Danh mục cha" prop="parentId" width="120" /> -->
        <el-table-column align="left" label="Hành động" fixed="right" min-width="240">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)"><el-icon
                style="margin-right: 5px">
                <InfoFilled />
              </el-icon>Xem chi tiết</el-button>
            <el-button type="primary" link icon="edit" class="table-button"
              @click="updateAttendanceCategoryFunc(scope.row)">Chỉnh sửa</el-button>
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
            <el-button @click="closeDialog">Hủy</el-button>
          </div>
        </div>
      </template>

      <el-form :model="formData" label-position="top" ref="elFormRef" :rules="rule" label-width="80px">
        <el-form-item label="Tên phân loại:" prop="name">
          <el-input v-model="formData.name" type="text" :clearable="true" placeholder="Nhập tên phân loại" />
        </el-form-item>
        <el-form-item label="Danh mục cha:" prop="parentId">
          <!-- <el-input v-model.number="formData.parentId" :clearable="true" placeholder="Nhập danh mục cha" /> -->
          <el-select v-model="formData.parentId" clearable filterable placeholder="Chọn danh mục cha">
            <el-option v-for="item in parentOptions" :key="item.ID" :value="item.ID" :label="item.name"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
    </el-drawer>

    <el-drawer destroy-on-close size="800" v-model="detailShow" :show-close="true" :before-close="closeDetailShow">
      <el-descriptions column="1" border>
        <el-descriptions-item label="Tên phân loại">
          {{ detailFrom.name }}
        </el-descriptions-item>
        <el-descriptions-item label="Danh mục cha">
          {{ detailFrom.parentId }}
        </el-descriptions-item>
      </el-descriptions>
    </el-drawer>

  </div>
</template>

<script setup>
import { deleteApisByIds } from '@/api/api';
import {
  createAttendanceCategory,
  deleteAttendanceCategory,
  deleteAttendanceCategoryByIds,
  updateAttendanceCategory,
  findAttendanceCategory,
  getAttendanceCategoryList
} from '@/api/checkins/attendanceCategory'

import { getDictFunc, formatDate, formatBoolean, filterDict, filterDataSource, returnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

defineOptions({
  name: 'AttendanceCategory'
})


const showAllQuery = ref(false)


const formData = ref({
  name: null,
  parentId: null,
})

const rule = reactive({
  name: [{
    required: true,
    message: 'Tên không được để trống',
    trigger: ['input', 'blur'],
  },
  ],
})

const parentOptions = ref([])

const searchRule = reactive({
  createdAt: [
    {
      validator: (rule, value, callback) => {
        if (searchInfo.value.startCreatedAt && !searchInfo.value.endCreatedAt) {
          callback(new Error('Vui lòng nhập ngày kết thúc'))
        } else if (!searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt) {
          callback(new Error('Vui lòng nhập ngày bắt đầu'))
        } else if (searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt && (searchInfo.value.startCreatedAt.getTime() === searchInfo.value.endCreatedAt.getTime() || searchInfo.value.startCreatedAt.getTime() > searchInfo.value.endCreatedAt.getTime())) {
          callback(new Error('Ngày bắt đầu phải trước ngày kết thúc'))
        } else {
          callback()
        }
      }, trigger: 'change'
    }
  ],
})

const elFormRef = ref()
const elSearchFormRef = ref()

// =========== Phần điều khiển bảng ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

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

// Truy vấn
const getTableData = async () => {
  const table = await getAttendanceCategoryList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

const getParentOptions = async () => {
  const table = await getAttendanceCategoryList({ page: 0, pageSize: -1 })
  if (table.code === 0) {
    parentOptions.value = table.data.list
  }
  console.log("parent Options", parentOptions.value)
}

getParentOptions()

// ============== Phần điều khiển bảng kết thúc ===============

// Lấy các từ điển cần thiết có thể trống, giữ theo nhu cầu
const setOptions = async () => {
}

// Lấy các từ điển cần thiết có thể trống, giữ theo nhu cầu
setOptions()


// Dữ liệu đa chọn
const multipleSelection = ref([])
// Đa chọn
const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

// Xóa hàng
const deleteRow = (row) => {
  ElMessageBox.confirm('Bạn có chắc muốn xóa không?', 'Cảnh báo', {
    confirmButtonText: 'Đồng ý',
    cancelButtonText: 'Hủy',
    type: 'warning'
  }).then(() => {
    deleteAttendanceCategoryFunc(row)
  })
}

// Xóa nhiều hàng
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
    const res = await deleteAttendanceCategoryByIds({ IDs })
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

// Đánh dấu hành động (thêm hoặc sửa) trong cửa sổ pop-up
const type = ref('')

// Cập nhật hàng
const updateAttendanceCategoryFunc = async (row) => {
  const res = await findAttendanceCategory({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data
    dialogFormVisible.value = true
  }
}


// Xóa hàng
const deleteAttendanceCategoryFunc = async (row) => {
  const res = await deleteAttendanceCategory({ ID: row.ID })
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: 'Xóa thành công'
    })
    if (tableData.value.length === 1 && page.value > 1) {
      page.value--
    }
    getTableData()
  }
}

// Đánh dấu cửa sổ pop-up
const dialogFormVisible = ref(false)

// Mở cửa sổ pop-up
const openDialog = () => {
  type.value = 'create'
  dialogFormVisible.value = true
}

// Đóng cửa sổ pop-up
const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = {
    name: undefined,
    parentId: undefined,
  }
}
// Xác nhận cửa sổ pop-up
const enterDialog = async () => {
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return
    let res
    switch (type.value) {
      case 'create':
        res = await createAttendanceCategory(formData.value)
        break
      case 'update':
        res = await updateAttendanceCategory(formData.value)
        break
      default:
        res = await createAttendanceCategory(formData.value)
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

// Mở cửa sổ chi tiết
const detailShow = ref(false)


// Mở cửa sổ chi tiết
const openDetailShow = () => {
  detailShow.value = true
}


// Mở chi tiết
const getDetails = async (row) => {
  // Mở cửa sổ pop-up
  const res = await findAttendanceCategory({ ID: row.ID })
  if (res.code === 0) {
    detailFrom.value = res.data
    openDetailShow()
  }
}


// Đóng cửa sổ chi tiết
const closeDetailShow = () => {
  detailShow.value = false
  detailFrom.value = {}
}



</script>

<style></style>
