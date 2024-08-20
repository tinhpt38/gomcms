<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list justify-between">
        <span class="text font-bold">Chi tiết từ điển</span>
        <el-button
          type="primary"
          icon="plus"
          @click="openDrawer"
        >
          Thêm mục từ điển
        </el-button>
      </div>
      <el-table
        ref="multipleTable"
        :data="tableData"
        style="width: 100%"
        tooltip-effect="dark"
        row-key="ID"
      >
        <el-table-column
          type="selection"
          width="55"
        />
        <el-table-column
          align="left"
          label="Ngày tạo"
          width="180"
        >
          <template #default="scope">
            {{ formatDate(scope.row.CreatedAt) }}
          </template>
        </el-table-column>

        <el-table-column
          align="left"
          label="Giá trị hiển thị"
          prop="label"
        />

        <el-table-column
          align="left"
          label="Giá trị từ điển"
          prop="value"
        />

        <el-table-column
          align="left"
          label="Giá trị mở rộng"
          prop="extend"
        />

        <el-table-column
          align="left"
          label="Trạng thái"
          prop="status"
          width="120"
        >
          <template #default="scope">
            {{ formatBoolean(scope.row.status) }}
          </template>
        </el-table-column>

        <el-table-column
          align="left"
          label="Đánh dấu sắp xếp"
          prop="sort"
          width="120"
        />

        <el-table-column
          align="left"
          label="Hành động"
          width="180"
        >
          <template #default="scope">
            <el-button
              type="primary"
              link
              icon="edit"
              @click="updateSysDictionaryDetailFunc(scope.row)"
            >
              Chỉnh sửa
            </el-button>
            <el-button
              type="primary"
              link
              icon="delete"
              @click="deleteSysDictionaryDetailFunc(scope.row)"
            >
              Xóa
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <el-drawer
      v-model="drawerFormVisible"
      size="30%"
      :show-close="false"
      :before-close="closeDrawer"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type==='create' ? 'Thêm mục từ điển' : 'Chỉnh sửa mục từ điển' }}</span>
          <div>
            <el-button @click="closeDrawer">
              Hủy
            </el-button>
            <el-button type="primary" @click="enterDrawer">
              Xác nhận
            </el-button>
          </div>
        </div>
      </template>
      <el-form
        ref="drawerForm"
        :model="formData"
        :rules="rules"
        label-width="110px"
      >
        <el-form-item
          label="Giá trị hiển thị"
          prop="label"
        >
          <el-input
            v-model="formData.label"
            placeholder="Nhập giá trị hiển thị"
            clearable
            :style="{width: '100%'}"
          />
        </el-form-item>
        <el-form-item
          label="Giá trị từ điển"
          prop="value"
        >
          <el-input
            v-model="formData.value"
            placeholder="Nhập giá trị từ điển"
            clearable
            :style="{width: '100%'}"
          />
        </el-form-item>
        <el-form-item
          label="Giá trị mở rộng"
          prop="extend"
        >
          <el-input
            v-model="formData.extend"
            placeholder="Nhập giá trị mở rộng"
            clearable
            :style="{width: '100%'}"
          />
        </el-form-item>
        <el-form-item
          label="Trạng thái"
          prop="status"
          required
        >
          <el-switch
            v-model="formData.status"
            active-text="Bật"
            inactive-text="Tắt"
          />
        </el-form-item>
        <el-form-item
          label="Đánh dấu sắp xếp"
          prop="sort"
        >
          <el-input-number
            v-model.number="formData.sort"
            placeholder="Đánh dấu sắp xếp"
          />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createSysDictionaryDetail,
  deleteSysDictionaryDetail,
  updateSysDictionaryDetail,
  findSysDictionaryDetail,
  getSysDictionaryDetailList
} from '@/api/sysDictionaryDetail' // Replace this with your own address
import { ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatBoolean, formatDate } from '@/utils/format'

defineOptions({
  name: 'SysDictionaryDetail'
})

const props = defineProps({
  sysDictionaryID: {
    type: Number,
    default: 0
  }
})

const formData = ref({
  label: null,
  value: null,
  status: true,
  sort: null
})
const rules = ref({
  label: [
    {
      required: true,
      message: 'Nhập giá trị hiển thị',
      trigger: 'blur'
    }
  ],
  value: [
    {
      required: true,
      message: 'Nhập giá trị từ điển',
      trigger: 'blur'
    }
  ],
  sort: [
    {
      required: true,
      message: 'Nhập đánh dấu sắp xếp',
      trigger: 'blur'
    }
  ]
})

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])

// Pagination
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// Query
const getTableData = async() => {
  const table = await getSysDictionaryDetailList({
    page: page.value,
    pageSize: pageSize.value,
    sysDictionaryID: props.sysDictionaryID
  })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

const type = ref('')
const drawerFormVisible = ref(false)
const updateSysDictionaryDetailFunc = async(row) => {
  drawerForm.value && drawerForm.value.clearValidate()
  const res = await findSysDictionaryDetail({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data.reSysDictionaryDetail
    drawerFormVisible.value = true
  }
}

const closeDrawer = () => {
  drawerFormVisible.value = false
  formData.value = {
    label: null,
    value: null,
    status: true,
    sort: null,
    sysDictionaryID: props.sysDictionaryID
  }
}
const deleteSysDictionaryDetailFunc = async(row) => {
  ElMessageBox.confirm('Bạn có chắc muốn xóa không?', 'Thông báo', {
    confirmButtonText: 'Đồng ý',
    cancelButtonText: 'Hủy',
    type: 'warning'
  }).then(async() => {
    const res = await deleteSysDictionaryDetail({ ID: row.ID })
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
  })
}

const drawerForm = ref(null)
const enterDrawer = async() => {
  drawerForm.value.validate(async valid => {
    formData.value.sysDictionaryID = props.sysDictionaryID
    if (!valid) return
    let res
    switch (type.value) {
      case 'create':
        res = await createSysDictionaryDetail(formData.value)
        break
      case 'update':
        res = await updateSysDictionaryDetail(formData.value)
        break
      default:
        res = await createSysDictionaryDetail(formData.value)
        break
    }
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: 'Tạo/Chỉnh sửa thành công'
      })
      closeDrawer()
      getTableData()
    }
  })
}
const openDrawer = () => {
  type.value = 'create'
  drawerForm.value && drawerForm.value.clearValidate()
  drawerFormVisible.value = true
}

watch(() => props.sysDictionaryID, () => {
  getTableData()
})

</script>

<style>
</style>
