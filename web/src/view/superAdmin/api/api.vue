<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="Đường dẫn">
          <el-input v-model="searchInfo.path" placeholder="Đường dẫn" />
        </el-form-item>
        <el-form-item label="Mô tả">
          <el-input v-model="searchInfo.description" placeholder="Mô tả" />
        </el-form-item>
        <el-form-item label="Nhóm API">
          <el-select v-model="searchInfo.apiGroup" clearable placeholder="Vui lòng chọn">
            <el-option v-for="item in apiGroupOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="Yêu cầu">
          <el-select v-model="searchInfo.method" clearable placeholder="Vui lòng chọn">
            <el-option v-for="item in methodOptions" :key="item.value" :label="`${item.label}(${item.value})`"
              :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">
            Tìm kiếm
          </el-button>
          <el-button icon="refresh" @click="onReset">
            Đặt lại
          </el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog('addApi')">
          Thêm mới
        </el-button>
        <el-button icon="delete" :disabled="!apis.length" @click="onDelete">
          Xóa
        </el-button>
        <el-button icon="Refresh" @click="onFresh">
          Làm mới bộ nhớ cache
        </el-button>
        <el-button icon="Compass" @click="onSync">
          Đồng bộ API
        </el-button>
        <ExportTemplate template-id="api" />
        <ExportExcel template-id="api" :limit="9999" />
        <ImportExcel template-id="api" @on-success="getTableData" />
      </div>
      <el-table :data="tableData" @sort-change="sortChange" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="id" min-width="60" prop="ID" sortable="custom" />
        <el-table-column align="left" label="Đường dẫn API" min-width="150" prop="path" sortable="custom" />
        <el-table-column align="left" label="Nhóm API" min-width="150" prop="apiGroup" sortable="custom" />
        <el-table-column align="left" label="Giới thiệu API" min-width="150" prop="description" sortable="custom" />
        <el-table-column align="left" label="Yêu cầu" min-width="150" prop="method" sortable="custom">
          <template #default="scope">
            <div>
              {{ scope.row.method }} / {{ methodFilter(scope.row.method) }}
            </div>
          </template>
        </el-table-column>

        <el-table-column align="left" fixed="right" label="Hành động" width="200">
          <template #default="scope">
            <el-button icon="edit" type="primary" link @click="editApiFunc(scope.row)">
              Chỉnh sửa
            </el-button>
            <el-button icon="delete" type="primary" link @click="deleteApiFunc(scope.row)">
              Xóa
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination :current-page="page" :page-size="pageSize" :page-sizes="[10, 30, 50, 100]" :total="total"
          layout="total, sizes, prev, pager, next, jumper" @current-change="handleCurrentChange"
          @size-change="handleSizeChange" />
      </div>
    </div>

    <el-drawer v-model="syncApiFlag" size="80%" :before-close="closeSyncDialog" :show-close="false">
      <warning-bar title="Đồng bộ API, không nhập nhóm định tuyến sẽ không được đồng bộ tự động" />
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">Đồng bộ định tuyến</span>
          <div>
            <el-button @click="closeSyncDialog">
              Hủy bỏ
            </el-button>
            <el-button type="primary" :loading="syncing" @click="enterSyncDialog">
              Xác nhận
            </el-button>
          </div>
        </div>
      </template>

      <h4>Định tuyến mới <span class="text-xs text-gray-500 ml-2 font-normal">Tồn tại trong định tuyến hiện tại nhưng
          không
          tồn tại trong bảng api</span></h4>
      <el-table :data="syncApiData.newApis">
        <el-table-column align="left" label="Đường dẫn API" min-width="150" prop="path" />
        <el-table-column align="left" label="Nhóm API" min-width="150" prop="apiGroup">
          <template #default="{ row }">
            <el-select v-model="row.apiGroup" placeholder="Vui lòng chọn hoặc thêm mới" allow-create filterable
              default-first-option>
              <el-option v-for="item in apiGroupOptions" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </template>
        </el-table-column>
        <el-table-column align="left" label="Giới thiệu API" min-width="150" prop="description">
          <template #default="{ row }">
            <el-input v-model="row.description" autocomplete="off" />
          </template>
        </el-table-column>
        <el-table-column align="left" label="Yêu cầu" min-width="150" prop="method">
          <template #default="scope">
            <div>
              {{ scope.row.method }} / {{ methodFilter(scope.row.method) }}
            </div>
          </template>
        </el-table-column>
        <el-table-column label="Hành động" min-width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" text @click="ignoreApiFunc(row, true)">
              Bỏ qua
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <h4>Định tuyến đã xóa <span class="text-xs text-gray-500 ml-2 font-normal">Không còn tồn tại trong dự án hiện tại,
          sau
          khi đồng bộ sẽ tự động xóa khỏi bảng apis</span></h4>
      <el-table :data="syncApiData.deleteApis">
        <el-table-column align="left" label="Đường dẫn API" min-width="150" prop="path" />
        <el-table-column align="left" label="Nhóm API" min-width="150" prop="apiGroup" />
        <el-table-column align="left" label="Giới thiệu API" min-width="150" prop="description" />
        <el-table-column align="left" label="Yêu cầu" min-width="150" prop="method">
          <template #default="scope">
            <div>
              {{ scope.row.method }} / {{ methodFilter(scope.row.method) }}
            </div>
          </template>
        </el-table-column>
      </el-table>

      <h4>Bỏ qua định tuyến <span class="text-xs text-gray-500 ml-2 font-normal">Bỏ qua định tuyến không tham gia đồng
          bộ
          api, thường là những định tuyến không cần xác thực</span></h4>
      <el-table :data="syncApiData.ignoreApis">
        <el-table-column align="left" label="Đường dẫn API" min-width="150" prop="path" />
        <el-table-column align="left" label="Nhóm API" min-width="150" prop="apiGroup" />
        <el-table-column align="left" label="Giới thiệu API" min-width="150" prop="description" />
        <el-table-column align="left" label="Yêu cầu" min-width="150" prop="method">
          <template #default="scope">
            <div>
              {{ scope.row.method }} / {{ methodFilter(scope.row.method) }}
            </div>
          </template>
        </el-table-column>
        <el-table-column label="Hành động" min-width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" text @click="ignoreApiFunc(row, false)">
              Hủy bỏ bỏ qua
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-drawer>

    <el-drawer v-model="dialogFormVisible" size="60%" :before-close="closeDialog" :show-close="false">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ dialogTitle }}</span>
          <div>
            <el-button @click="closeDialog">
              Hủy bỏ
            </el-button>
            <el-button type="primary" @click="enterDialog">
              Xác nhận
            </el-button>
          </div>
        </div>
      </template>

      <warning-bar title="Thêm mới API, cần cấu hình quyền trong quản lý vai trò để sử dụng" />
      <el-form ref="apiForm" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="Đường dẫn" prop="path">
          <el-input v-model="form.path" autocomplete="off" />
        </el-form-item>
        <el-form-item label="Yêu cầu" prop="method">
          <el-select v-model="form.method" placeholder="Vui lòng chọn" style="width:100%">
            <el-option v-for="item in methodOptions" :key="item.value" :label="`${item.label}(${item.value})`"
              :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="Nhóm API" prop="apiGroup">
          <el-select v-model="form.apiGroup" placeholder="Vui lòng chọn hoặc thêm mới" allow-create filterable
            default-first-option>
            <el-option v-for="item in apiGroupOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="Giới thiệu API" prop="description">
          <el-input v-model="form.description" autocomplete="off" />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  getApiById,
  getApiList,
  createApi,
  updateApi,
  deleteApi,
  deleteApisByIds,
  freshCasbin,
  syncApi,
  getApiGroups,
  ignoreApi,
  enterSyncApi
} from '@/api/api'
import { toSQLLine } from '@/utils/stringFun'
import WarningBar from '@/components/warningBar/warningBar.vue'
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import ExportExcel from '@/components/exportExcel/exportExcel.vue'
import ExportTemplate from '@/components/exportExcel/exportTemplate.vue'
import ImportExcel from '@/components/exportExcel/importExcel.vue'

defineOptions({
  name: 'Api',
})

const methodFilter = (value) => {
  const target = methodOptions.value.filter(item => item.value === value)[0]
  return target && `${target.label}`
}

const apis = ref([])
const form = ref({
  path: '',
  apiGroup: '',
  method: '',
  description: ''
})
const methodOptions = ref([
  {
    value: 'POST',
    label: 'Tạo mới',
    type: 'success'
  },
  {
    value: 'GET',
    label: 'Xem',
    type: ''
  },
  {
    value: 'PUT',
    label: 'Cập nhật',
    type: 'warning'
  },
  {
    value: 'DELETE',
    label: 'Xóa',
    type: 'danger'
  }
])

const type = ref('')
const rules = ref({
  path: [{ required: true, message: 'Vui lòng nhập đường dẫn API', trigger: 'blur' }],
  apiGroup: [
    { required: true, message: 'Vui lòng nhập tên nhóm', trigger: 'blur' }
  ],
  method: [
    { required: true, message: 'Vui lòng chọn phương thức', trigger: 'blur' }
  ],
  description: [
    { required: true, message: 'Vui lòng nhập mô tả API', trigger: 'blur' }
  ]
})

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const apiGroupOptions = ref([])
const apiGroupMap = ref({})

const getGroup = async () => {
  const res = await getApiGroups()
  if (res.code === 0) {
    const groups = res.data.groups
    apiGroupOptions.value = groups.map(item => ({ label: item, value: item }))
    apiGroupMap.value = res.data.apiGroupMap
  }
}

const ignoreApiFunc = async (row, flag) => {
  const res = await ignoreApi({ path: row.path, method: row.method, flag })
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: res.msg
    })
    if (flag) {
      syncApiData.value.newApis = syncApiData.value.newApis.filter(item => !(item.path === row.path && item.method === row.method))
      syncApiData.value.ignoreApis.push(row)
      return
    }
    syncApiData.value.ignoreApis = syncApiData.value.ignoreApis.filter(item => !(item.path === row.path && item.method === row.method))
    syncApiData.value.newApis.push(row)
  }
}

const closeSyncDialog = () => {
  syncApiFlag.value = false
}

const syncing = ref(false)


const enterSyncDialog = async () => {
  syncing.value = true
  const res = await enterSyncApi(syncApiData.value)
  syncing.value = false
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: res.msg
    })
    syncApiFlag.value = false
    getTableData()
  }
}

const onReset = () => {
  searchInfo.value = {}
}
// Tìm kiếm

const onSubmit = () => {
  page.value = 1
  pageSize.value = 10
  getTableData()
}

// Phân trang
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// Sắp xếp
const sortChange = ({ prop, order }) => {
  if (prop) {
    if (prop === 'ID') {
      prop = 'id'
    }
    searchInfo.value.orderKey = toSQLLine(prop)
    searchInfo.value.desc = order === 'descending'
  }
  getTableData()
}

// Tìm kiếm
const getTableData = async () => {
  const table = await getApiList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()
getGroup()
// Thao tác hàng loạt
const handleSelectionChange = (val) => {
  apis.value = val
}

const onDelete = async () => {
  ElMessageBox.confirm('Bạn có chắc muốn xóa không?', 'Cảnh báo', {
    confirmButtonText: 'Đồng ý',
    cancelButtonText: 'Hủy',
    type: 'warning'
  }).then(async () => {
    const ids = apis.value.map(item => item.ID)
    const res = await deleteApisByIds({ ids })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: res.msg
      })
      if (tableData.value.length === ids.length && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  })
}
const onFresh = async () => {
  ElMessageBox.confirm('Bạn có chắc muốn làm mới cache không?', 'Cảnh báo', {
    confirmButtonText: 'Đồng ý',
    cancelButtonText: 'Hủy',
    type: 'warning'
  }).then(async () => {
    const res = await freshCasbin()
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: res.msg
      })
    }
  })
}

const syncApiData = ref({
  newApis: [],
  deleteApis: [],
  ignoreApis: []
})

const syncApiFlag = ref(false)

const onSync = async () => {
  const res = await syncApi()
  if (res.code === 0) {
    res.data.newApis.forEach(item => {
      item.apiGroup = apiGroupMap.value[item.path.split('/')[1]]
      console.log(apiGroupMap.value)
    })

    syncApiData.value = res.data
    syncApiFlag.value = true
  }
}

// Cửa sổ pop-up liên quan
const apiForm = ref(null)
const initForm = () => {
  apiForm.value.resetFields()
  form.value = {
    path: '',
    apiGroup: '',
    method: '',
    description: ''
  }
}

const dialogTitle = ref('Thêm mới API')
const dialogFormVisible = ref(false)
const openDialog = (key) => {
  switch (key) {
    case 'addApi':
      dialogTitle.value = 'Thêm mới API'
      break
    case 'edit':
      dialogTitle.value = 'Chỉnh sửa API'
      break
    default:
      break
  }
  type.value = key
  dialogFormVisible.value = true
}
const closeDialog = () => {
  initForm()
  dialogFormVisible.value = false
}

const editApiFunc = async (row) => {
  const res = await getApiById({ id: row.ID })
  form.value = res.data.api
  openDialog('edit')
}

const enterDialog = async () => {
  apiForm.value.validate(async valid => {
    if (valid) {
      switch (type.value) {
        case 'addApi':
          {
            const res = await createApi(form.value)
            if (res.code === 0) {
              ElMessage({
                type: 'success',
                message: 'Thêm thành công',
                showClose: true
              })
            }
            getTableData()
            getGroup()
            closeDialog()
          }

          break
        case 'edit':
          {
            const res = await updateApi(form.value)
            if (res.code === 0) {
              ElMessage({
                type: 'success',
                message: 'Chỉnh sửa thành công',
                showClose: true
              })
            }
            getTableData()
            closeDialog()
          }
          break
        default:
          // eslint-disable-next-line no-lone-blocks
          {
            ElMessage({
              type: 'error',
              message: 'Thao tác không xác định',
              showClose: true
            })
          }
          break
      }
    }
  })
}

const deleteApiFunc = async (row) => {
  ElMessageBox.confirm('Thao tác này sẽ xóa vĩnh viễn tất cả các vai trò chứa API này, bạn có muốn tiếp tục?', 'Cảnh báo', {
    confirmButtonText: 'Đồng ý',
    cancelButtonText: 'Hủy',
    type: 'warning'
  })
    .then(async () => {
      const res = await deleteApi(row)
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: 'Xóa thành công!'
        })
        if (tableData.value.length === 1 && page.value > 1) {
          page.value--
        }
        getTableData()
        getGroup()
      }
    })
}

</script>

<style scoped lang="scss">
.warning {
  color: #dc143c;
}
</style>
