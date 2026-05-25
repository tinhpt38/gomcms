<template>
  <div>
    <div>
      <div class="p-1 my-1 flex gap-2">
        <el-button type="primary" icon="plus" @click="openDialog()">
          Thêm nhóm
        </el-button>
        <el-button type="success" icon="sort" @click="openAutoGroup()">
          Phân nhóm tự động
        </el-button>
      </div>

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
          </template>
        </el-table-column>
      </el-table>
      <div class="flex justify-end">
        <el-pagination v-model:current-page="page" v-model:page-size="pageSize" :page-sizes="[20, 50, 100, 500]"
          :size="size" :background="true" layout="total, sizes, prev, pager, next, jumper" :total="total"
          @size-change="handleSizeChange" @current-change="handleCurrentChange" />
      </div>
    </div>

    <!-- Drawer tạo / sửa nhóm -->
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
          <el-input v-model="formData.name" :clearable="true" placeholder="Vui lòng nhập tên nhóm" />
        </el-form-item>
      </el-form>
    </el-drawer>

    <!-- Drawer phân nhóm tự động -->
    <el-drawer v-model="autoGroupVisible" destroy-on-close size="800" :show-close="false" :before-close="closeAutoGroup">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">Phân nhóm tự động</span>
          <div>
            <el-button type="primary" @click="enterGroupDialog">Đồng ý</el-button>
            <el-button @click="closeAutoGroup">Huỷ</el-button>
          </div>
        </div>
      </template>
      <el-form ref="autoGroupFormRef" :rules="groupRules" :model="autoGroupFormData" label-position="top"
        label-width="80px">
        <el-form-item label="Chế độ phân nhóm:" prop="mode">
          <el-radio-group v-model="autoGroupFormData.mode">
            <el-radio value="create">
              <span class="font-medium">Tạo nhóm mới</span>
              <span class="text-sm text-red-500 ml-1">(xoá toàn bộ nhóm cũ)</span>
            </el-radio>
            <el-radio value="reassign">
              <span class="font-medium">Giữ nhóm hiện tại</span>
              <span class="text-sm text-gray-500 ml-1">(chỉ chia lại thành viên)</span>
            </el-radio>
          </el-radio-group>
        </el-form-item>
        <template v-if="autoGroupFormData.mode === 'create'">
          <el-alert type="warning" :closable="false" class="mb-4">
            Thao tác này sẽ <strong>xoá toàn bộ nhóm hiện tại</strong> và tạo lại. Gõ <strong>XÁC NHẬN</strong> bên
            dưới để tiếp tục.
          </el-alert>
          <el-form-item label="Số nhóm:" prop="groupQty">
            <el-input v-model="autoGroupFormData.groupQty" :clearable="true" type="number"
              placeholder="Nhập số nhóm" />
          </el-form-item>
          <el-form-item label="Cách tạo tên:" prop="groupNameType">
            <el-select v-model="autoGroupFormData.groupNameType" placeholder="Vui lòng chọn" style="width:100%"
              :clearable="true">
              <el-option label="Số thứ tự" value="baseOnNumberic" />
              <el-option label="Bảng chữ cái" value="baseOnAlphabet" />
            </el-select>
          </el-form-item>
          <el-form-item label="Xác nhận:" prop="confirmText">
            <el-input v-model="autoGroupFormData.confirmText" placeholder='Nhập "XÁC NHẬN"' />
          </el-form-item>
        </template>
      </el-form>
    </el-drawer>
  </div>
</template>


<script setup>
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

import {
  getGroupDataSource,
  createGroup,
  deleteGroup,
  updateGroup,
  findGroup,
  getGroupList,
  assignParticipantToGroupAuto,
  reassignParticipantsOnly,
  getGroupDependencyCount,
} from '@/api/checkins/group'

import { syncCondition } from '@/api/checkins/condition'

const props = defineProps({
  acId: {
    type: Number,
    required: true
  }
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

const rule = reactive({
  name: [
    { required: true, message: 'Tên nhóm là bắt buộc', trigger: ['input', 'blur'] },
    { whitespace: true, message: 'Tên nhóm không thể chỉ là khoảng trắng', trigger: ['input', 'blur'] },
  ],
})

const groupRules = reactive({
  groupQty: [{ required: true, message: 'Số lượng nhóm là bắt buộc', trigger: ['input', 'blur'] }],
  groupNameType: [{ required: true, message: 'Cách đặt tên là bắt buộc', trigger: ['input', 'blur'] }],
  confirmText: [
    {
      validator: (rule, value, callback) => {
        if (autoGroupFormData.value.mode === 'create' && value !== 'XÁC NHẬN') {
          callback(new Error('Vui lòng nhập đúng "XÁC NHẬN"'))
        } else {
          callback()
        }
      },
      trigger: ['input', 'blur'],
    },
  ],
})

const formData = ref({ name: '', attendanceId: props.acId * 1 })

const autoGroupFormData = ref({
  attendanceId: props.acId * 1,
  groupNameType: 'baseOnAlphabet',
  mode: 'reassign',
  confirmText: '',
})

const dataSource = ref([])
const getDataSourceFunc = async () => {
  const res = await getGroupDataSource()
  if (res.code === 0) {
    dataSource.value = res.data
  }
}
getDataSourceFunc()
getTableData()

// Guard: cảnh báo nếu nhóm còn thành viên hoặc điều kiện
const deleteRow = async (row) => {
  const depRes = await getGroupDependencyCount({ ID: row.ID })
  if (depRes.code === 0) {
    const { memberCount, conditionCount } = depRes.data
    let warningMsg = 'Bạn có chắc muốn xoá không?'
    if (memberCount > 0 || conditionCount > 0) {
      warningMsg = `Nhóm này đang có ${memberCount} thành viên và ${conditionCount} điều kiện liên kết. Xoá nhóm sẽ gỡ liên kết các bản ghi này. Bạn có chắc chắn không?`
    }
    ElMessageBox.confirm(warningMsg, 'Cảnh báo', {
      confirmButtonText: 'Đồng ý',
      cancelButtonText: 'Hủy',
      type: 'warning',
    }).then(() => {
      deleteGroupFunc(row)
    })
  }
}

const dialogFormVisible = ref(false)

const openDialog = () => {
  type.value = 'create'
  formData.value = { name: '', attendanceId: props.acId * 1 }
  dialogFormVisible.value = true
}

const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = { name: '', attendanceId: props.acId * 1 }
}

const elFormRef = ref()
const autoGroupFormRef = ref()

const autoGroupVisible = ref(false)
const openAutoGroup = () => {
  autoGroupFormData.value = {
    attendanceId: props.acId * 1,
    groupNameType: 'baseOnAlphabet',
    mode: 'reassign',
    confirmText: '',
  }
  autoGroupVisible.value = true
}
const closeAutoGroup = () => {
  autoGroupVisible.value = false
}

const enterGroupDialog = async () => {
  autoGroupFormRef.value?.validate(async (valid) => {
    if (!valid) return
    let res
    if (autoGroupFormData.value.mode === 'reassign') {
      res = await reassignParticipantsOnly({ attendanceId: +props.acId })
    } else {
      res = await assignParticipantToGroupAuto({
        groupQty: +autoGroupFormData.value.groupQty,
        groupNameType: autoGroupFormData.value.groupNameType,
        attendanceId: +props.acId,
      })
    }
    if (res.code === 0) {
      await syncCondition({ attendanceId: +props.acId })
      ElMessage({ type: 'success', message: 'Phân nhóm thành công và đã đồng bộ điều kiện' })
      closeAutoGroup()
      getTableData()
      emits('onSuccess')
    }
  })
}

const enterDialog = async () => {
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return
    let res
    if (type.value === 'create') {
      res = await createGroup(formData.value)
    } else {
      res = await updateGroup(formData.value)
    }
    if (res.code === 0) {
      ElMessage({ type: 'success', message: 'Tạo/cập nhật thành công' })
      closeDialog()
      getTableData()
      emits('onSuccess')
    }
  })
}

const deleteGroupFunc = async (row) => {
  const res = await deleteGroup({ ID: row.ID })
  if (res.code === 0) {
    ElMessage({ type: 'success', message: 'Xoá thành công' })
    if (tableData.value.length === 1 && page.value > 1) {
      page.value--
    }
    getTableData()
    emits('onSuccess')
  }
}

const updateGroupFunc = async (row) => {
  const res = await findGroup({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data
    dialogFormVisible.value = true
  }
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}
</script>

<style scoped>
/* Add your component styles here */
</style>
