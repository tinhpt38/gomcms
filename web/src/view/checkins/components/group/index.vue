<template>
  <div>
    <!-- Phần quản lý nhóm -->
    <div class="p-1 my-1">
      <el-button type="primary" icon="plus" @click="openDialog()">
        Thêm nhóm
      </el-button>
      <el-tooltip class="box-item" effect="dark" content="Chức năng này sẽ xoá toàn bộ nhóm, tạo lại nhóm và tự động phân bổ thành viên vào nhóm" placement="top-start">
        <el-button type="success" icon="loading" @click="openAutoGroup()">
          Xoá và tạo lại nhóm Phân nhóm
        </el-button>
      </el-tooltip>
    </div>

    <!-- Danh sách nhóm -->
    <el-table :data="tableData" style="width: 100%" border>
      <el-table-column prop="name" label="Nhóm" />
      <el-table-column prop="total" label="Số thành viên" />
      <el-table-column label="Hành động">
        <template #default="scope">
          <el-button size="small" type="primary" plain round @click="updateGroupFunc(scope.row)">
            Sửa
          </el-button>
          <el-button size="small" type="danger" plain round @click="deleteRow(scope.row)">
            Xoá
          </el-button>
          <!-- Nút Thêm thành viên -->
          <el-button size="small" type="info" plain round @click="openAddMemberDialog(scope.row)">
            Thêm thành viên
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- Phân trang -->
    <div class="flex justify-end">
      <el-pagination v-model:current-page="page" v-model:page-size="pageSize" :page-sizes="[20, 50, 100, 500]"
        :size="size" :background="true" layout="total, sizes, prev, pager, next, jumper" :total="total"
        @size-change="handleSizeChange" @current-change="handleCurrentChange" />
    </div>

    <!-- Drawer Thêm/Cập nhật nhóm -->
    <el-drawer v-model="dialogFormVisible" destroy-on-close size="800" :show-close="false" :before-close="closeDialog">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? 'Thêm mới' : 'Chỉnh sửa' }}</span>
          <div>
            <el-button type="primary" @click="enterDialog">Đồng ý</el-button>
            <el-button @click="closeDialog">Huỷ</el-button>
          </div>
        </div>
      </template>
      <el-form ref="elFormRef" :model="formData" label-position="top" :rules="rule" label-width="80px">
        <el-form-item label="Tên nhóm:" prop="name">
          <el-input v-model="formData.name" clearable placeholder="Vui lòng nhập Tên nhóm" />
        </el-form-item>
        <el-form-item label="Attendance Class:" prop="attendanceId" class="hidden">
          <el-select v-model="formData.attendanceId" placeholder="Vui lòng chọn Attendance Class" style="width:100%" clearable>
            <el-option v-for="(item, key) in dataSource.attendanceId" :key="key" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
      </el-form>
    </el-drawer>

    <!-- Drawer Phân nhóm tự động -->
    <el-drawer v-model="autoGroupVisible" destroy-on-close size="800" :show-close="false" :before-close="closeAutoGroup">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">Phân nhóm</span>
          <div>
            <el-button type="primary" @click="enterGroupDialog">Đồng ý</el-button>
            <el-button @click="closeAutoGroup">Huỷ</el-button>
          </div>
        </div>
      </template>
      <div class="text-m py-4">Số thành viên sẽ được chia đều vào các nhóm</div>
      <el-form ref="autoGroupFormRef" :rules="groupRules" :model="autoGroupFormData" label-position="top" label-width="80px">
        <el-form-item label="Số nhóm:" prop="groupQty">
          <el-input v-model="autoGroupFormData.groupQty" clearable type="number" placeholder="Vui lòng số nhóm" />
        </el-form-item>
        <el-form-item label="Cách tạo tên:" prop="groupNameType">
          <el-select v-model="autoGroupFormData.groupNameType" placeholder="Vui lòng chọn" style="width:100%" clearable>
            <el-option label="Số thứ tự" value="baseOnNumberic" />
            <el-option label="Bảng chữ cái" value="baseOnAlphabet" />
          </el-select>
        </el-form-item>
      </el-form>
    </el-drawer>

    <!-- Dialog Thêm thành viên vào nhóm (Bảng danh sách sinh viên chưa tham gia nhóm) -->
    <el-dialog title="Thêm thành viên vào nhóm" v-model="addMemberDialogVisible" width="600px">
      <!-- Ô tìm kiếm -->
      <div style="margin-bottom: 10px;">
        <el-input 
          v-model="participantSearch" 
          placeholder="Tìm kiếm theo tên hoặc email" 
          clearable
          @input="fetchParticipants">
        </el-input>
      </div>
      <!-- Bảng danh sách sinh viên -->
      <el-table 
        ref="participantTable"
        :data="participantList" 
        style="width: 100%" 
        border 
        @selection-change="handleSelectionChange"
        row-key="ID">
        <el-table-column type="selection" width="55"></el-table-column>
        <el-table-column prop="fullName" label="Họ và tên" />
        <el-table-column prop="email" label="Email" />
      </el-table>
      <template #footer>
        <el-button @click="addMemberDialogVisible = false">Huỷ</el-button>
        <el-button type="primary" @click="handleBulkAdd">Thêm hàng loạt</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import {
  getGroupDataSource,
  createGroup,
  deleteGroup,
  deleteGroupByIds,
  updateGroup,
  findGroup,
  getGroupList,
  assignParticipantToGroupAuto,
  addMemberToGroupApi
} from '@/api/checkins/group'
import service from '@/utils/request'

const props = defineProps({
  acId: { type: Number, required: true }
})

const tableData = ref([])
const searchInfo = ref({ attendanceId: null })
const type = ref('')

const page = ref(1)
const pageSize = ref(20)
const size = ref(20)
const total = ref(0)
const emits = defineEmits(['onSuccess'])

const getTableData = async () => {
  searchInfo.value.attendanceId = props.acId
  const table = await getGroupList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}
getTableData()

// Quy tắc validate cho form nhóm
const rule = reactive({
  name: [
    { required: true, message: 'Tên nhóm là bắt buộc', trigger: ['input', 'blur'] },
    { whitespace: true, message: 'Bạn không thể nhập mỗi khoảng trắng', trigger: ['input', 'blur'] }
  ]
})

// Quy tắc validate cho form phân nhóm tự động
const groupRules = reactive({
  groupQty: [{ required: true, message: 'Số lượng nhóm là bắt buộc', trigger: ['input', 'blur'] }],
  groupNameType: [{ required: true, message: 'Cách đặt tên là bắt buộc', trigger: ['input', 'blur'] }]
})

const formData = ref({ name: '', attendanceId: props.acId })
const autoGroupFormData = ref({ attendanceId: props.acId, groupNameType: "baseOnAlphabet", groupQty: '' })

const dataSource = ref([])
const getDataSourceFunc = async () => {
  const res = await getGroupDataSource()
  if (res.code === 0) {
    dataSource.value = res.data
  }
}
getDataSourceFunc()

const deleteRow = (row) => {
  ElMessageBox.confirm('Bạn có chắc muốn xoá không?', 'Cảnh báo', {
    confirmButtonText: 'Đồng ý',
    cancelButtonText: 'Hủy',
    type: 'warning'
  }).then(() => {
    deleteGroupFunc(row)
  })
  emits('onSuccess')
}

const dialogFormVisible = ref(false)
const openDialog = () => {
  type.value = 'create'
  dialogFormVisible.value = true
}
const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = { name: '', attendanceId: props.acId }
}
const elFormRef = ref()
const autoGroupFormRef = ref()

const autoGroupVisible = ref(false)
const openAutoGroup = () => { autoGroupVisible.value = true }
const closeAutoGroup = () => { autoGroupVisible.value = false }

const enterGroupDialog = async () => {
  autoGroupFormRef.value?.validate(async (valid) => {
    if (!valid) return
    const res = await assignParticipantToGroupAuto({
      groupQty: +autoGroupFormData.value.groupQty,
      groupNameType: autoGroupFormData.value.groupNameType,
      attendanceId: +props.acId
    })
    if (res.code === 0) {
      ElMessage({ type: 'success', message: 'Phân nhóm thành công' })
      closeAutoGroup()
      getTableData()
    }
  })
}

const enterDialog = async () => {
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return
    let res
    if (type.value === 'create') {
      res = await createGroup(formData.value)
    } else if (type.value === 'update') {
      res = await updateGroup(formData.value)
    }
    if (res.code === 0) {
      ElMessage({ type: 'success', message: 'Tạo/cập nhật thành công' })
      closeDialog()
      getTableData()
    }
    emits('onSuccess')
  })
}

const onDelete = async () => {
  ElMessageBox.confirm('Bạn có chắc muốn xoá không?', 'Cảnh báo', {
    confirmButtonText: 'Đồng ý',
    cancelButtonText: 'Hủy',
    type: 'warning'
  }).then(async () => {
    const IDs = []
    if (!multipleSelection.value || multipleSelection.value.length === 0) {
      ElMessage({ type: 'warning', message: 'Vui lòng chọn dữ liệu để xoá' })
      return
    }
    multipleSelection.value.forEach(item => { IDs.push(item.ID) })
    const res = await deleteGroupByIds({ IDs })
    if (res.code === 0) {
      ElMessage({ type: 'success', message: 'Xoá thành công' })
      if (tableData.value.length === IDs.length && page.value > 1) page.value--
      getTableData()
    }
  })
  emits('onSuccess')
}

const deleteGroupFunc = async (row) => {
  const res = await deleteGroup({ ID: row.ID })
  if (res.code === 0) {
    ElMessage({ type: 'success', message: 'Xoá thành công' })
    if (tableData.value.length === 1 && page.value > 1) page.value--
    getTableData()
  }
  emits('onSuccess')
}

const updateGroupFunc = async (row) => {
  const res = await findGroup({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data
    dialogFormVisible.value = true
  }
  emits('onSuccess')
}

// ---------------------------
// Phần thêm thành viên vào nhóm
// ---------------------------
const addMemberDialogVisible = ref(false)
const participantSearch = ref('')          // Giá trị tìm kiếm
const participantList = ref([])            // Danh sách sinh viên chưa được phân nhóm
const selectedParticipants = ref([])         // Danh sách sinh viên được chọn
const currentGroupId = ref(null)             // Lưu ID của nhóm được chọn từ tag nhóm

// --- Hàm mở dialog thêm thành viên ---
const openAddMemberDialog = (group) => {
  console.log("openAddMemberDialog:", group)
  currentGroupId.value = group.ID
  addMemberDialogVisible.value = true
  participantSearch.value = ''
  selectedParticipants.value = []
  fetchParticipants()  // Tải danh sách sinh viên ban đầu
}

// --- Hàm gọi API lấy danh sách sinh viên chưa tham gia nhóm ---  
const fetchParticipants = async () => {
  try {
    const res = await service({
      url: '/group/getParticipantsForGroup',
      method: 'get',
      params: { attendanceId: props.acId, query: participantSearch.value }
    })
    if (res.code === 0 && res.data) {
      participantList.value = res.data
    } else {
      participantList.value = []
    }
  } catch (error) {
    console.error('Lỗi khi lấy danh sách sinh viên:', error)
    participantList.value = []
  }
}

// --- Hàm xử lý khi người dùng chọn sinh viên trong bảng ---
const handleSelectionChange = (val) => {
  selectedParticipants.value = val
}

// --- Hàm xử lý thêm hàng loạt sinh viên vào nhóm ---
const handleBulkAdd = async () => {
  if (selectedParticipants.value.length === 0) {
    ElMessage.warning("Vui lòng chọn ít nhất 1 sinh viên")
    return
  }
  try {
    // Lặp qua danh sách sinh viên được chọn và gọi API thêm
    for (const participant of selectedParticipants.value) {
      const res = await service({
        url: '/group/AddMember',
        method: 'post',
        data: {
          groupId: currentGroupId.value,
          email: participant.email,
          name: participant.fullName
        }
      })
      if (res.code !== 0) {
        ElMessage.error(`Lỗi khi thêm ${participant.fullName}`)
      }
    }
    ElMessage.success("Thêm thành viên thành công")
    addMemberDialogVisible.value = false
    // Nếu cần, bạn có thể gọi lại hàm cập nhật dữ liệu nhóm
  } catch (error) {
    ElMessage.error("Có lỗi xảy ra khi thêm thành viên")
  }
}
</script>

<style scoped>
/* Add your component styles here */
</style>
