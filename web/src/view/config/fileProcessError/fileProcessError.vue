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

        <el-form-item label="Mã quá trình xử lý tệp tin" prop="fileProcessId">

          <el-input v-model.number="searchInfo.fileProcessId" placeholder="Điều kiện tìm kiếm" />

        </el-form-item>
        <el-form-item label="Cột dữ liệu" prop="fieldTitle">
          <el-input v-model="searchInfo.fieldTitle" placeholder="Điều kiện tìm kiếm" />

        </el-form-item>

        <template v-if="showAllQuery">
          <!-- Thêm các điều kiện tìm kiếm cần điều khiển hiển thị ở đây -->
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

        <el-table-column align="left" label="Mã quá trình xử lý tệp tin" prop="fileProcessId" width="120" />
        <el-table-column align="left" label="Mã định danh của tệp tin" prop="fileProcessUuid" width="300" />
        <el-table-column align="left" label="Cột dữ liệu" prop="fieldTitle" width="120" />
        <el-table-column align="left" label="Hành động" fixed="right" min-width="240">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)"><el-icon
                style="margin-right: 5px">
                <InfoFilled />
              </el-icon>Xem chi tiết</el-button>
            <el-button type="primary" link icon="edit" class="table-button hidden"
              @click="updateFileProcessErrorFunc(scope.row)">Chỉnh sửa</el-button>
            <el-button type="primary" link icon="delete" class="hidden" @click="deleteRow(scope.row)">Xóa</el-button>
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
        <el-form-item label="Mã quá trình xử lý tệp tin:" prop="fileProcessId">
          <el-input v-model.number="formData.fileProcessId" :clearable="true"
            placeholder="Nhập mã quá trình xử lý tệp tin" />
        </el-form-item>
        <el-form-item label="Mã định danh của tệp tin:" prop="fileProcessUuid">
          <el-input v-model="formData.fileProcessUuid" :clearable="true" placeholder="Nhập mã định danh của tệp tin" />
        </el-form-item>
        <el-form-item label="Cột dữ liệu:" prop="fieldTitle">
          <el-input v-model="formData.fieldTitle" :clearable="true" placeholder="Nhập cột dữ liệu" />
        </el-form-item>
      </el-form>
    </el-drawer>

    <el-drawer destroy-on-close size="800" v-model="detailShow" :show-close="true" :before-close="closeDetailShow">
      <el-descriptions column="1" border>
        <el-descriptions-item label="Mã quá trình xử lý tệp tin">
          {{ detailFrom.fileProcessId }}
        </el-descriptions-item>
        <el-descriptions-item label="Mã định danh của tệp tin">
          {{ detailFrom.fileProcessUuid }}
        </el-descriptions-item>
        <el-descriptions-item label="Cột dữ liệu">
          {{ detailFrom.fieldTitle }}
        </el-descriptions-item>
      </el-descriptions>
    </el-drawer>

  </div>
</template>

<script setup>
import {
  createFileProcessError,
  deleteFileProcessError,
  deleteFileProcessErrorByIds,
  updateFileProcessError,
  findFileProcessError,
  getFileProcessErrorList
} from '@/api/config/fileProcessError'


import { getDictFunc, formatDate, formatBoolean, filterDict, filterDataSource, returnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

defineOptions({
  name: 'FileProcessError'
})


const showAllQuery = ref(false)

const formData = ref({
  fileProcessId: undefined,
  fileProcessUuid: '',
  fieldTitle: '',
})


const rule = reactive({
  fileProcessId: [{
    required: true,
    message: 'Vui lòng nhập mã quá trình xử lý tệp tin',
    trigger: ['input', 'blur'],
  }],
  fileProcessUuid: [{
    required: true,
    message: 'Vui lòng nhập mã định danh của tệp tin',
    trigger: ['input', 'blur'],
  },
  {
    whitespace: true,
    message: 'Không được chỉ nhập khoảng trắng',
    trigger: ['input', 'blur'],
  }],
  fieldTitle: [{
    required: true,
    message: 'Vui lòng nhập cột dữ liệu',
    trigger: ['input', 'blur'],
  },
  {
    whitespace: true,
    message: 'Không được chỉ nhập khoảng trắng',
    trigger: ['input', 'blur'],
  }],
})

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


const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})


const onReset = () => {
  searchInfo.value = {}
  getTableData()
}


const onSubmit = () => {
  elSearchFormRef.value?.validate(async (valid) => {
    if (!valid) return
    page.value = 1
    pageSize.value = 10
    getTableData()
  })
}


const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}


const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}


const getTableData = async () => {
  const table = await getFileProcessErrorList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

const setOptions = async () => {
}


setOptions()



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
    deleteFileProcessErrorFunc(row)
  })
}

// 多选删除
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
    const res = await deleteFileProcessErrorByIds({ IDs })
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

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')

// Cập nhật hàng
const updateFileProcessErrorFunc = async (row) => {
  const res = await findFileProcessError({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data
    dialogFormVisible.value = true
  }
}


// Xóa hàng
const deleteFileProcessErrorFunc = async (row) => {
  const res = await deleteFileProcessError({ ID: row.ID })
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

// Điều khiển hiển thị của hộp thoại
const dialogFormVisible = ref(false)

// Mở hộp thoại
const openDialog = () => {
  type.value = 'create'
  dialogFormVisible.value = true
}

// Đóng hộp thoại
const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = {
    fileProcessId: undefined,
    fileProcessUuid: '',
    fieldTitle: '',
  }
}
// Xác nhận hộp thoại
const enterDialog = async () => {
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return
    let res
    switch (type.value) {
      case 'create':
        res = await createFileProcessError(formData.value)
        break
      case 'update':
        res = await updateFileProcessError(formData.value)
        break
      default:
        res = await createFileProcessError(formData.value)
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

// Điều khiển hiển thị chi tiết
const detailShow = ref(false)


// Mở hộp thoại chi tiết
const openDetailShow = () => {
  detailShow.value = true
}


// Xem chi tiết
const getDetails = async (row) => {
  // Mở hộp thoại
  const res = await findFileProcessError({ ID: row.ID })
  if (res.code === 0) {
    detailFrom.value = res.data
    openDetailShow()
  }
}


// Đóng hộp thoại chi tiết
const closeDetailShow = () => {
  detailShow.value = false
  detailFrom.value = {}
}



</script>

<style></style>
