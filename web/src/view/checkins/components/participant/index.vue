<template>
  <div>
    <div class="p-1 my-1 flex gap-2">
      <el-button type="primary" icon="plus" @click="openDialog">
        Thêm thành viên
      </el-button>
      <el-button type="success" icon="document" @click="openBulkAdd">
        Thêm hàng loạt
      </el-button>
    </div>

    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="Họ và tên">
          <el-input v-model="searchInfo.fullName" placeholder="Họ và tên" />
        </el-form-item>
        <el-form-item label="Email">
          <el-input v-model="searchInfo.email" placeholder="Email" />
        </el-form-item>
        <el-form-item label="Nhóm:">
          <el-select v-model="searchInfo.groupId" placeholder="Chọn nhóm" filterable clearable>
            <el-option v-for="item in groupOptions" :key="item.ID" :label="item.name" :value="item.ID" />
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
    <el-table :data="tableData" style="width: 100%" border>
      <el-table-column prop="fullName" label="Họ và tên" width="250">
        <template #default="scope">
          <span>{{ scope.row?.fullName?.replace("undefined", "") }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="email" label="Email" width="300" />
      <el-table-column label="Nhóm" min-width="300">
        <template #default="scope">
          <span>{{ scope.row.groups.map((e) => e.name).join(", ") }}</span>
        </template>
      </el-table-column>
      <el-table-column align="right" label="Hành động" fixed="right" min-width="240">
        <template #default="scope">
          <el-button type="primary" link icon="edit" class="table-button"
            @click="updateParticipantFunc(scope.row)">Chỉnh sửa</el-button>
          <el-button type="danger" link icon="delete" class="table-button"
            @click="deleteParticipantRow(scope.row)">Xoá</el-button>
        </template>
      </el-table-column>
    </el-table>
    <div class="flex justify-end">
      <el-pagination v-model:current-page="page" v-model:page-size="pageSize" :page-sizes="[20, 50, 100, 500]"
        :size="size" :background="true" layout="total, sizes, prev, pager, next, jumper" :total="total"
        @size-change="handleSizeChange" @current-change="handleCurrentChange" />
    </div>

    <!-- Drawer thêm/sửa đơn lẻ -->
    <el-drawer destroy-on-close size="800" v-model="dialogFormVisible" :show-close="false" :before-close="closeDialog">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? 'Thêm thành viên' : 'Chỉnh sửa thành viên' }}</span>
          <div>
            <el-button type="primary" @click="enterDialog">Đồng ý</el-button>
            <el-button @click="closeDialog">Hủy</el-button>
          </div>
        </div>
      </template>

      <el-form :model="formData" label-position="top" ref="elFormRef" :rules="rule" label-width="80px">
        <el-form-item label="Họ và tên:" prop="fullName">
          <el-input v-model="formData.fullName" :clearable="true" placeholder="Nhập họ và tên" />
        </el-form-item>
        <el-form-item label="Email:" prop="email">
          <el-input v-model="formData.email" :clearable="true" placeholder="Nhập email" />
        </el-form-item>
        <el-form-item label="Nhóm:" prop="groupId">
          <el-select v-model="formData.groupId" placeholder="Chọn nhóm" filterable clearable multiple style="width:100%">
            <el-option v-for="item in groupOptions" :key="item.ID" :label="item.name" :value="item.ID" />
          </el-select>
        </el-form-item>
      </el-form>
    </el-drawer>

    <!-- Drawer thêm hàng loạt -->
    <el-drawer destroy-on-close size="800" v-model="bulkAddVisible" :show-close="false">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">Thêm hàng loạt thành viên</span>
          <div>
            <el-button type="primary" @click="bulkAddParticipantsFunc">Đồng ý</el-button>
            <el-button @click="bulkAddVisible = false">Hủy</el-button>
          </div>
        </div>
      </template>

      <el-form :model="bulkFormData" ref="bulkFormRef" :rules="bulkRules" label-position="top" label-width="80px">
        <div class="font-bold py-2">Copy và dán danh sách email của thành viên vào đây. Mỗi email một dòng, tối đa
          1000 email / lần.</div>
        <el-form-item label="Danh sách email:" prop="list">
          <el-input type="textarea" :rows="20" v-model="bulkFormData.list" :clearable="true"
            placeholder="email1@example.com&#10;email2@example.com" />
        </el-form-item>
        <el-form-item label="Nhóm (có thể chọn nhiều):" prop="groupIds">
          <el-select v-model="bulkFormData.groupIds" placeholder="Chọn nhóm" filterable clearable multiple
            style="width:100%">
            <el-option v-for="item in groupOptions" :key="item.ID" :label="item.name" :value="item.ID" />
          </el-select>
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createParticipant,
  deleteParticipant,
  updateParticipant,
  bulkParticipants,
  getParticipantListByAttendance,
} from '@/api/checkins/participant'

import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

const props = defineProps({
  acId: {
    type: Number,
    required: true
  },
  groupOptions: {
    type: Array,
    default: () => [],
  },
})

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const size = ref(20)
const tableData = ref([])
const formData = ref({})

const searchInfo = ref({
  fullName: '',
  email: '',
  attendanceId: props.acId,
  groupId: null,
})

const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

const getTableData = async () => {
  searchInfo.value.attendanceId = props.acId
  const table = await getParticipantListByAttendance({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list.map((e) => {
      e.groupId = e.groups?.map(k => k.ID)
      return e
    })
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

const deleteParticipantRow = (row) => {
  ElMessageBox.confirm('Bạn có chắc muốn xóa thành viên này không?', 'Cảnh báo', {
    confirmButtonText: 'Đồng ý',
    cancelButtonText: 'Hủy',
    type: 'warning'
  }).then(() => {
    deleteParticipantFunc(row)
  })
}

const type = ref('')

const updateParticipantFunc = async (row) => {
  type.value = 'update'
  formData.value = { ...row }
  dialogFormVisible.value = true
}

const deleteParticipantFunc = async (row) => {
  const res = await deleteParticipant({ ID: row.ID, attendanceId: Number(props.acId) })
  if (res.code === 0) {
    ElMessage({ type: 'success', message: 'Xoá thành công' })
    if (tableData.value.length === 1 && page.value > 1) {
      page.value--
    }
    getTableData()
  }
}

const dialogFormVisible = ref(false)

const openDialog = () => {
  type.value = 'create'
  formData.value = { fullName: '', email: '', groupId: [] }
  dialogFormVisible.value = true
}

const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = {}
}

const elFormRef = ref()

const rule = reactive({
  email: [
    { required: true, message: 'Email không được để trống', trigger: ['input', 'blur'] },
    { type: 'email', message: 'Email không đúng định dạng', trigger: ['input', 'blur'] },
  ],
})

const enterDialog = async () => {
  formData.value.attendanceId = Number(props.acId)
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return
    let res
    if (type.value === 'create') {
      res = await createParticipant(formData.value)
    } else {
      res = await updateParticipant(formData.value)
    }
    if (res.code === 0) {
      ElMessage({ type: 'success', message: 'Tạo/cập nhật thành công' })
      closeDialog()
      getTableData()
    }
  })
}

const onSubmit = () => {
  getTableData()
}

const onReset = () => {
  searchInfo.value = { fullName: '', email: '', attendanceId: props.acId, groupId: null }
  getTableData()
}

// Bulk add
const bulkAddVisible = ref(false)
const bulkFormData = ref({ list: '', groupIds: [] })
const bulkFormRef = ref()

const bulkRules = reactive({
  list: [
    {
      validator: (rule, value, callback) => {
        if (!value || !value.trim()) {
          callback(new Error('Danh sách email không được để trống'))
        } else if (value.trim().split('\n').length > 1000) {
          callback(new Error('Giới hạn tối đa 1000 thành viên / lần'))
        } else {
          callback()
        }
      },
      trigger: 'blur',
    },
  ],
})

const openBulkAdd = () => {
  bulkFormData.value = { list: '', groupIds: [] }
  bulkAddVisible.value = true
}

const bulkAddParticipantsFunc = async () => {
  bulkFormRef.value?.validate(async (valid) => {
    if (!valid) return
    const list = bulkFormData.value.list.trim().split('\n').map(e => e.trim()).filter(Boolean)
    const groupIds = bulkFormData.value.groupIds?.map(Number) ?? []
    const res = await bulkParticipants({
      list,
      attendanceId: Number(props.acId),
      groupIds,
      groupId: groupIds[0] ?? undefined,
    })
    if (res.code === 0) {
      ElMessage({ type: 'success', message: 'Thêm thành công' })
      bulkAddVisible.value = false
      getTableData()
    }
  })
}
</script>

<style scoped>
/* Add your component styles here */
</style>
