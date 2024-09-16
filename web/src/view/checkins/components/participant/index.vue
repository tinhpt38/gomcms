<template>
  <div>
    <div class="p-1 my-1">
      <el-button type="primary" icon="plus" @click="openBulkAdd">
        Thêm thành viên
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
    <div class="p-2 my-2">
      <span>Nếu <strong> số Tổng thành công (pass) </strong> lớn hơn <strong> số Tổng điều kiện (total) </strong> có
        nghĩa là đã có <strong>(pass - total)</strong> điều kiện bị xoá bởi người tạo.</span>
    </div>
    <el-table :data="tableData" style="width: 100%" border>
      <el-table-column prop="fullName" label="Họ và tên" />
      <el-table-column prop="email" label="Email" />
      <el-table-column label="Nhóm">
        <template #default="scope">
          <span>{{ scope.row.groups.map((e) => e.name).join(", ") }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="passCount" label="Điểm danh thành công" />
      <el-table-column prop="conditionCount" label="Tổng điều kiện cần" />
      <el-table-column prop="requestCount" label="Tổng truy cập" />
      <el-table-column align="left" label="Hành động" fixed="right" min-width="240">
        <template #default="scope">
          <el-button type="primary" link icon="edit" class="table-button"
            @click="updateParticipantFunc(scope.row)">Chỉnh sửa</el-button>
        </template>
      </el-table-column>
    </el-table>
    <div class="flex justify-end">
      <el-pagination v-model:current-page="page" v-model:page-size="pageSize" :page-sizes="[20, 50, 100, 500]"
        :size="size" :background="true" layout="total, sizes, prev, pager, next, jumper" :total="total"
        @size-change="handleSizeChange" @current-change="handleCurrentChange" />
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
        <el-form-item label="Họ và tên:" prop="fullName">
          <el-input v-model="formData.fullName" :clearable="true" placeholder="Nhập Họ và tên" />
        </el-form-item>
        <el-form-item label="Email:" prop="email">
          <el-input v-model="formData.email" :clearable="true" placeholder="Nhập Email" />
        </el-form-item>
        <el-form-item label="Nhóm:" prop="groupId">
          <el-select v-model="formData.groupId" placeholder="Chọn nhóm" filterable clearable>
            <el-option v-for="item in groupOptions" :key="item.ID" :label="item.name" :value="item.ID" />
          </el-select>
        </el-form-item>
      </el-form>
    </el-drawer>

    <el-drawer destroy-on-close size="800" v-model="bulkAddVisible" :show-close="false">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">Thêm hàng loạt thành viên</span>
          <div>
            <el-button type="primary" @click="bulkAddParticipantsFunc">Đồng ý</el-button>
            <el-button @click="bulkAddVisible == false">Hủy</el-button>
          </div>
        </div>
      </template>

      <el-form :model="bulkFormData" ref="bulkFormRef" :rules="bulkRules" label-position="top" label-width="80px">
        <div class="font-bold py-2">Copy và dán danh sách email của thành viên vào đây. Mỗi email một dòng, tối đa 1000
          email/ lần</div>
        <el-form-item label="Nhập danh sách thành viên:" prop="list">
          <el-input type="textarea" :rows="20" v-model="bulkFormData.list" :clearable="true"
            placeholder="Nhập mã số email thành viên" />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createParticipant,
  deleteParticipant,
  deleteParticipantByIds,
  findParticipant,
  updateParticipant,
  bulkParticipants,
  getParticipantListByAttendance
} from '@/api/checkins/participant'
import { getDictFunc, formatDate, formatBoolean, filterDict, filterDataSource, returnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

const props = defineProps({
  acId: {
    type: Number,
    required: true
  },
  groupOptions: {
    type: Array,
    required: false
  }
})
const page = ref(0)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])

const formData = ref({
})

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
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
  console.log("participant components")
  console.log(tableData.value)
}

getTableData()


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
    deleteParticipantFunc(row)
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
    const res = await deleteParticipantByIds({ IDs })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: 'Xoá thành công'
      })
      if (tableData.value.length === IDs.length && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  })
}


const type = ref('')


const updateParticipantFunc = async (row) => {
  const res = await findParticipant({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data
    dialogFormVisible.value = true
  }
}



const deleteParticipantFunc = async (row) => {
  const res = await deleteParticipant({ ID: row.ID })
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
    fullName: '',
    email: '',
  }
}

const elFormRef = ref()
const enterDialog = async () => {
  if (formData.value.groupId) {
    formData.value.groupId = Number(formData.value.groupId)
  }
  formData.value.attendanceId = Number(props.acId)
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return
    let res
    switch (type.value) {
      case 'create':
        res = await createParticipant(formData.value)
        break
      case 'update':
        res = await updateParticipant(formData.value)
        break
      default:
        res = await createParticipant(formData.value)
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

const onSubmit = () => {
  getTableData()
}

const onReset = () => {
  searchInfo.value = {
    fullName: '',
    email: '',
    attendanceId: props.acId
  }
  getTableData()
}

const bulkAddVisible = ref(false)
const bulkFormData = ref({
  list: ''
})
const bulkFormRef = ref()

const bulkRules = reactive({
  list: [
    {
      validator: (rule, value, callback) => {
        if (!value) {
          callback(new Error('Không được để trống'))
        } else if (value.split("\n").length > 1000) {
          callback(new Error('Giới hạn tối đa 1000 thành viên / lần'))
        } else {
          callback()
        }
      }, trigger: 'blur'
    }
  ],
})

const openBulkAdd = () => {
  bulkAddVisible.value = true
}

const bulkAddParticipantsFunc = async () => {
  bulkFormRef.value?.validate(async (valid) => {
    if (!valid) return
    var list = bulkFormData.value.list.split("\n")
    var res = await bulkParticipants({ list, attendanceId: Number(props.acId) })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: 'Thêm thành công'
      })
      bulkAddVisible.value = false
      getTableData()
    }

  })
}

</script>

<style scoped>
/* Add your component styles here */
</style>