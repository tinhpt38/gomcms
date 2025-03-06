<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule"
        @keyup.enter="onSubmit">
        <el-form-item label="Đơn vị" prop="agencyId" class="">
          <el-select v-model="searchInfo.agencyId" placeholder="Chọn đơn vị" clearable filterable>
            <el-option v-for="item in agencyOptions" :key="item.ID" :label="item.name" :value="item.ID" />
          </el-select>
        </el-form-item>
        <el-form-item label="Danh mục" prop="categoryId" class="w-[300px]">
          <el-tree-select class="w-full" v-model="searchInfo.categoryId" :data="categoryOptions" check-on-click-node
            :render-after-expand="false" style="width: 240px" />
        </el-form-item>


        <template v-if="showAllQuery">
          <el-form-item label="Ngày tạo" prop="createdAt">
            <template #label>
              <span>
                Ngày tạo˝
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
        </template>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">
            Tìm kiếm
          </el-button>
          <el-button icon="refresh" @click="onReset">
            Đặt lại
          </el-button>
          <el-button v-if="!showAllQuery" link type="primary" icon="arrow-down" @click="showAllQuery = true">
            Mở rộng
          </el-button>
          <el-button v-else link type="primary" icon="arrow-up" @click="showAllQuery = false">
            Thu gọn
          </el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog">
          Thêm mới
        </el-button>
        <!-- <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="onDelete">
          Xóa
        </el-button> -->
      </div>
      <el-table ref="multipleTable" style="width: 100%" tooltip-effect="dark" :data="tableData" row-key="ID"
        @selection-change="handleSelectionChange">
        <!-- <el-table-column type="selection" width="55" /> -->

        <el-table-column align="left" label="Tiêu đề" prop="title" width="250" />
        <el-table-column align="left" label="Bắt đầu" prop="startDate" width="120">
          <template #default="scope">
            {{ formatDateTime(scope.row.startDate) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="Kết thúc" prop="endDate" width="120">
          <template #default="scope">
            {{ formatDateTime(scope.row.endDate) }}
          </template>
        </el-table-column>
        <el-table-column label="Đơn vị" prop="agency.name" width="150" />
        <el-table-column align="left" label="Danh mục" prop="category.name" width="100" />
        <el-table-column align="left" label="Hệ số" prop="weight" width="100" />
        <el-table-column align="left" label="Điểm danh" width="100">
          <template #default="scope">
            {{ scope.row.totalCheckin }} / {{ scope.row.total }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="Trạng thái" prop="isLocked" width="100">
          <template #default="scope">
            {{ scope.row.isLocked ? "Đã đóng" : "Đang mở" }}
          </template>
        </el-table-column>
        <!-- <el-table-column align="left" label="Ngày tạo" prop="createdAt" width="120">
          <template #default="scope">
            {{ formatDate(scope.row.CreatedAt) }}
          </template>
        </el-table-column> -->
        <el-table-column align="left" label="Thao tác" fixed="right" min-width="240">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)"><el-icon
                style="margin-right: 5px">
                <InfoFilled />
              </el-icon>Xem chi tiết</el-button>
            <!-- <el-button type="primary" link icon="edit" class="table-button"
              @click="updateAttendanceFunc(scope.row)">Chỉnh sửa</el-button> -->
            <el-button type="primary" link icon="documentCopy" class="table-button" @click="cloneRow(scope.row)">Nhân
              bản</el-button>
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

      <el-form ref="elFormRef" :model="formData" label-position="top" :rules="rule" label-width="80px">
        <el-form-item label="Tiêu đề:" prop="title">
          <el-input v-model="formData.title" :clearable="true" placeholder="Nhập tiêu đề" />
        </el-form-item>
        <el-form-item label="Hệ số:" prop="weight">
          <el-input v-model="formData.weight" :clearable="true" type="number" placeholder="Nhập hệ số" />
        </el-form-item>
        <el-form-item label="Ngày bắt đầu:" prop="startDate">
          <el-date-picker v-model="formData.startDate" type="datetime" style="width:100%" placeholder="Chọn ngày"
            :clearable="true" />
        </el-form-item>
        <el-form-item label="Ngày kết thúc:" prop="endDate">
          <el-date-picker v-model="formData.endDate" type="datetime" style="width:100%" placeholder="Chọn ngày"
            :clearable="true" />
        </el-form-item>
        <el-form-item label="Đơn vị" prop="agencyId" class="w-full required">
          <el-select v-model="formData.agencyId" placeholder="Chọn đơn vị" clearable filterable>
            <el-option v-for="item in agencyOptions" :key="item.ID" :label="item.name" :value="item.ID" />
          </el-select>
        </el-form-item>
        <el-form-item label="Danh mục" prop="categoryId" class="w-full required">
          <el-tree-select class="w-full" v-model="formData.categoryId" :data="categoryOptions" check-on-click-node
            :render-after-expand="false" style="width: 240px" />
        </el-form-item>
        <el-form-item prop="description" label="Mô tả" class="w-full">
          <el-input type="text" v-model="formData.description"></el-input>
        </el-form-item>
        <el-form-item label="Cho phép khách:" prop="allowGuest">
          <el-switch v-model="formData.allowGuest" active-color="#13ce66" inactive-color="#ff4949" active-text="Có"
            inactive-text="Không" clearable />
        </el-form-item>
        <el-form-item label="Khoá:" prop="isLocked">
          <el-switch v-model="formData.isLocked" active-color="#13ce66" inactive-color="#ff4949" active-text="Có"
            inactive-text="Không" clearable />
        </el-form-item>
      </el-form>
    </el-drawer>

    <el-drawer v-model="detailShow" destroy-on-close size="800" :show-close="true" :before-close="closeDetailShow">
      <el-descriptions column="1" border>
        <el-descriptions-item label="Tiêu đề">
          {{ detailFrom.title }}
        </el-descriptions-item>
        <el-descriptions-item label="Ngày bắt đầu">
          {{ detailFrom.startDate }}
        </el-descriptions-item>
        <el-descriptions-item label="Ngày kết thúc">
          {{ detailFrom.endDate }}
        </el-descriptions-item>
        <el-descriptions-item label="Cho thử nghiệm">
          {{ detailFrom.isTrial }}
        </el-descriptions-item>
        <el-descriptions-item label="Khoá">
          {{ detailFrom.isLocked }}
        </el-descriptions-item>
      </el-descriptions>
    </el-drawer>
    <el-dialog v-model="dialogClone">
      <template #title>
        <span>Xác nhận</span>
      </template>
      <p>Bạn có muốn nhân bản mục này với dữ liệu thành viên không?</p>
      <br>
      <span slot="footer" class="dialog-footer mt-4">
        <el-button type="danger" @click="cloneAttendanceFun(false)">Không</el-button>
        <el-button type="primary" @click="cloneAttendanceFun(true)">Có</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script setup>
import {
  createAttendance,
  deleteAttendance,
  deleteAttendanceByIds,
  updateAttendance,
  findAttendance,
  getAttendanceList,
  cloneAttendance
} from '@/api/checkins/attendance'
import router from '@/router';
import {
  getAttendanceCategoryList
} from '@/api/checkins/attendanceCategory'

import {
  getAttendanceAgencyList
} from '@/api/checkins/attendanceAgency'


import { getDictFunc, formatDate, formatBoolean, filterDict, filterDataSource, returnArrImg, onDownloadFile, formatDateTime } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

defineOptions({
  name: 'Attendance'
})


const showAllQuery = ref(false)


const formData = ref({
  title: '',
  startDate: new Date(),
  endDate: new Date(),
  isTrial: false,
  isLocked: false,
})

const dialogClone = ref(false)


// 验证规则
const rule = reactive({
  title: [{
    required: true,
    message: 'Tiêu đề không được để trống',
    trigger: ['input', 'blur'],
  },
  {
    whitespace: true,
    message: 'Bạn không thể chỉ nhập dấu cách',
    trigger: ['input', 'blur'],
  }
  ],
  startDate: [{
    required: true,
    message: 'Ngày bắt đầu là bắt buộc',
    trigger: ['input', 'blur'],
  },
  ],
  endDate: [{
    required: true,
    message: 'Ngày kết thúc là bắt buộc',
    trigger: ['input', 'blur'],
  },
  ],
  categoryId: [{
    required: true,
    message: 'Danh mục không được để trống',
    trigger: ['input', 'blur'],
  }],
  agencyId: [{
    required: true,
    message: 'Đơn vị không được để trống',
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

// =========== Bảng điều khiển ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const categoryOptions = ref([])
const agencyOptions = ref([])
const clientURL = ref(import.meta.env.VITE_CLIENT_URL)

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

// Trang
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
  const table = await getAttendanceList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
  //console.log('table, ', table.data.list)
}

getTableData()


const getCategoryOptions = async () => {
  const table = await getAttendanceCategoryList({ page: 0, pageSize: -1 })
  if (table.code === 0) {
    categoryOptions.value = convertToTree(table.data.list)
    var selectCurrent = table.data.list.filter((e) => e.isCurrent)
    formData.value.categoryId = selectCurrent[0].ID
  }

}

getCategoryOptions()

const getAgencyOptions = async () => {
  const table = await getAttendanceAgencyList({ page: 0, pageSize: -1 })
  if (table.code === 0) {
    agencyOptions.value = table.data.list
  }

}
getAgencyOptions()

const convertToTree = (data) => {
  const map = {}
  const roots = []

  // Create a map of nodes using their ID as the key
  data.forEach((node) => {
    map[node.ID] = { ...node, value: node.ID, label: node.name, children: [] }
  })

  // Iterate over the nodes and assign children to their parent
  data.forEach((node) => {
    const parent = map[node.parentId]

    if (parent) {
      parent.children.push(map[node.ID])
    } else {
      roots.push(map[node.ID])
    }
  })
  //console.log('roots', roots)
  return roots
}



// ============== Kết thúc bảng điều khiển ===============

// Lấy các từ điển cần thiết (có thể trống)
const setOptions = async () => {
}

// Lấy các từ điển cần thiết (có thể trống)
setOptions()


// Dữ liệu đa chọn
const multipleSelection = ref([])
// Đa chọn
const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

// Xóa hàng
const deleteRow = (row) => {
  ElMessageBox.confirm('Khi xoá, sẽ mất hết toàn bộ dữ liệu và không thể phục hồi. Bạn có chắc muốn xóa không?', 'Cảnh báo', {
    confirmButtonText: 'Đồng ý',
    cancelButtonText: 'Hủy',
    type: 'warning'
  }).then(() => {
    deleteAttendanceFunc(row)
  })
}

const formClone = ref({
  id: null,
  withData: false
})
const cloneRow = (row) => {
  dialogClone.value = true
  formClone.value.id = row.ID
}

const cloneAttendanceFun = async (withData) => {
  dialogClone.value = false
  var res = await cloneAttendance({ attendanceId: formClone.value.id, withData: withData })
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: 'Nhân bản thành công'
    })
    getTableData()
  } else {
    ElMessage({
      type: 'error',
      message: res.message
    })
  }
  //console.log('res', res)
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
    const res = await deleteAttendanceByIds({ IDs })
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

// Đánh dấu hành vi (thêm hoặc sửa) hàng
const type = ref('')

// Cập nhật hàng
const updateAttendanceFunc = async (row) => {
  const res = await findAttendance({ id: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data
    dialogFormVisible.value = true
  }
}


// Xóa hàng
const deleteAttendanceFunc = async (row) => {
  const res = await deleteAttendance({ ID: row.ID })
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

// Đánh dấu điều khiển hộp thoại
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
    title: '',
    startDate: new Date(),
    endDate: new Date(),
    allowGuest: false,
    isLocked: false,
  }
}
// Xác nhận hộp thoại
const enterDialog = async () => {
  if (formData.value.weight) {
    formData.value.weight = Number(formData.value.weight)
  }

  elFormRef.value?.validate(async (valid) => {
    if (!valid) return
    let res
    switch (type.value) {
      case 'create':
        res = await createAttendance(formData.value)
        break
      case 'update':
        res = await updateAttendance(formData.value)
        break
      default:
        res = await createAttendance(formData.value)
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

// Đánh dấu điều khiển hiển thị chi tiết
const detailShow = ref(false)


// Mở chi tiết
const openDetailShow = () => {
  detailShow.value = true
}


// Lấy chi tiết
const getDetails = async (row) => {
  // Mở hộp thoại
  // const res = await findAttendance({ ID: row.ID })
  // if (res.code === 0) {
  // detailFrom.value = res.data
  // // openDetailShow()
  router.push({
    name: 'attendanceDetail',
    params: {
      id: row.ID
    }
  })
  // }
}


// Đóng hộp thoại chi tiết
const closeDetailShow = () => {
  detailShow.value = false
  detailFrom.value = {}
}


</script>

<style></style>
