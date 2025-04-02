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
      <!-- Mô tả nếu cần -->
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
      
      <!-- Cột hiển thị điểm danh -->
      <el-table-column label="Điểm danh" width="180">
        <template #default="scope">
          <span>{{ scope.row.passCount }} / {{ scope.row.conditionCount }}</span>
        </template>
      </el-table-column>

      <!-- Cột hiển thị danh sách điều kiện -->
      <el-table-column label="Tên điều kiện" min-width="300">
        <template #default="scope">
          <div v-if="scope.row.conditions && scope.row.conditions.length">
            <div v-for="(cond, index) in scope.row.conditions" :key="index">
              {{ cond.name }} ({{ cond.status }} - {{ cond.count }} lần)
            </div>
          </div>
          <div v-else>
            Không có dữ liệu
          </div>
        </template>
      </el-table-column>
      
      <el-table-column align="right" label="Hành động" fixed="right" min-width="240">
        <template #default="scope">
          <el-button type="primary" link icon="edit" class="table-button"
            @click="updateParticipantFunc(scope.row)">Chỉnh sửa</el-button>
          <el-button type="danger" link icon="delete" class="table-button"
            @click="deleteParticipantRow(scope.row)">Xoá</el-button>
          <el-button type="info" link icon="eye" class="table-button"
            @click="viewConditions(scope.row)">Xem điều kiện</el-button>
        </template>
      </el-table-column>
    </el-table>
    
    <div class="flex justify-end">
      <el-pagination v-model:current-page="page" v-model:page-size="pageSize" :page-sizes="[20, 50, 100, 500]"
        :size="size" :background="true" layout="total, sizes, prev, pager, next, jumper" :total="total"
        @size-change="handleSizeChange" @current-change="handleCurrentChange" />
    </div>

    <!-- Drawer chỉnh sửa/thêm thành viên -->
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
          <el-select v-model="formData.groupId" placeholder="Chọn nhóm" filterable clearable multiple>
            <el-option v-for="item in groupOptions" :key="item.ID" :label="item.name" :value="item.ID" />
          </el-select>
        </el-form-item>
      </el-form>
    </el-drawer>

    <!-- Drawer thêm hàng loạt thành viên -->
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
        <div class="font-bold py-2">
          Copy và dán danh sách email của thành viên vào đây. Mỗi email một dòng, tối đa 1000 email/ lần
        </div>
        <el-form-item label="Nhập danh sách thành viên:" prop="list">
          <el-input type="textarea" :rows="20" v-model="bulkFormData.list" :clearable="true"
            placeholder="Nhập danh sách email thành viên" />
        </el-form-item>
        <el-form-item label="Nhóm:" prop="groupId">
          <el-select v-model="bulkFormData.groupId" placeholder="Chọn nhóm" filterable clearable>
            <el-option v-for="item in groupOptions" :key="item.ID" :label="item.name" :value="item.ID" />
          </el-select>
        </el-form-item>
      </el-form>
    </el-drawer>

    <!-- Drawer hiển thị thông tin điều kiện -->
    <el-drawer destroy-on-close size="600" v-model="conditionDrawerVisible" :show-close="false">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">Điều kiện điểm danh</span>
          <el-button @click="conditionDrawerVisible = false">Đóng</el-button>
        </div>
      </template>
      <div>
        <p>Trạng thái: {{ currentCondition.status }}</p>
        <p>Số lần thành công: {{ currentCondition.passCount }}</p>
        <p>Tổng điều kiện: {{ currentCondition.conditionCount }}</p>
        <el-table :data="currentCondition.conditions" style="width: 100%">
          <el-table-column prop="name" label="Tên điều kiện" />
          <el-table-column prop="status" label="Trạng thái" />
          <el-table-column prop="count" label="Số lần" />
        </el-table>
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  createParticipant,
  deleteParticipant,
  deleteParticipantByIds,
  findParticipant,
  updateParticipant,
  bulkParticipants,
  getParticipantListByAttendance,
  getParticipantConditions
} from '@/api/checkins/participant'
import { syncCondition } from '@/api/checkins/condition'
import { findAttendanceArea } from '@/api/checkins/attendance'

const props = defineProps({
  acId: {
    type: Number,
    required: true
  },
  groupOptions: {
    type: Array,
    default: () => []
  }
})

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])

const formData = ref({})

const searchInfo = ref({
  fullName: '',
  email: '',
  attendanceId: props.acId,
  groupId: []
})

const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// Đồng bộ điều kiện cho toàn bộ thành viên trong phiên điểm danh
const syncConditionsForAttendance = async () => {
  const data = { AttendanceId: props.acId }
  try {
    const res = await syncCondition(data)
    if (res.code !== 0) {
      if (res.msg && res.msg.includes("Không đủ quyền")) {
        console.warn("Sync condition skipped: Không đủ quyền")
        return
      } else {
        ElMessage({
          type: 'error',
          message: 'Đồng bộ điều kiện thất bại: ' + res.msg
        })
      }
    }
  } catch (error) {
    console.error("Sync condition error:", error)
    ElMessage({
      type: 'error',
      message: 'Đồng bộ điều kiện thất bại'
    })
  }
}

const getTableData = async () => {
  await syncConditionsForAttendance()
  searchInfo.value.attendanceId = props.acId
  const table = await getParticipantListByAttendance({
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo.value
  })
  if (table.code === 0) {
    const list = table.data.list.map((e) => {
      e.groupId = e.groups?.map(k => k.ID)
      return e
    })

    await Promise.all(list.map(async (row) => {
      const res = await getParticipantConditions({ participantId: row.ID, attendanceId: props.acId })
      if (res.code === 0) {
        row.passCount = res.data.successCount
        row.conditionCount = res.data.requiredCount
        row.conditions = res.data.conditions || []
      } else {
        row.passCount = 0
        row.conditionCount = 0
        row.conditions = []
      }
    }))

    tableData.value = list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}
getTableData()

const multipleSelection = ref([])

const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

const deleteParticipantRow = (row) => {
  ElMessageBox.confirm('Bạn có chắc muốn xóa thành viên này không?', 'Cảnh báo', {
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
    multipleSelection.value.forEach(item => {
      IDs.push(item.ID)
    })
    const res = await deleteParticipantByIds({ IDs })
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

const type = ref('')

const updateParticipantFunc = async (row) => {
  type.value = 'update'
  formData.value = row
  dialogFormVisible.value = true
}

const deleteParticipantFunc = async (row) => {
  const res = await deleteParticipant({ ID: row.ID, attendanceId: Number(props.acId) })
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

const dialogFormVisible = ref(false)

const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = {
    fullName: '',
    email: '',
    groupId: []
  }
}

const elFormRef = ref()
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
    attendanceId: props.acId,
    groupId: []
  }
  getTableData()
}

const bulkAddVisible = ref(false)
const bulkFormData = ref({
  list: '',
  groupId: ''
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
      },
      trigger: 'blur'
    },
  ],
})

const openBulkAdd = () => {
  bulkAddVisible.value = true
}

const bulkAddParticipantsFunc = async () => {
  bulkFormRef.value?.validate(async (valid) => {
    if (!valid) return
    const list = bulkFormData.value.list.split("\n")
    const res = await bulkParticipants({ list, attendanceId: Number(props.acId), groupId: Number(bulkFormData.value.groupId) })
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

const conditionDrawerVisible = ref(false)
const currentCondition = reactive({
  status: '',
  passCount: 0,
  conditionCount: 0,
  conditions: []
})

const viewConditions = async (row) => {
  // console.log("Fetching conditions for:", row.ID);
  const res = await getParticipantConditions({ participantId: row.ID, attendanceId: props.acId });
  // console.log("API Response:", res);

  if (res.code !== 0 || !res.data) {
    ElMessage({ type: 'error', message: 'Lỗi API khi lấy điều kiện' });
    return;
  }

  row.passCount = res.data.successCount || 0;
  row.conditionCount = res.data.requiredCount || 0;
  currentCondition.status = (res.data.successCount && res.data.successCount > 0) ? "Đã điểm danh" : "Chưa điểm danh";
  currentCondition.passCount = res.data.successCount || 0;
  currentCondition.conditionCount = res.data.requiredCount || 0;

  let conditions = res.data.conditions || [];
  if (conditions.length === 0 && res.data.requiredCount > 0) {
    let areaName = "Khu vực";
    const areaRes = await findAttendanceArea({ id: props.acId });
    if (areaRes.code === 0 && areaRes.data && areaRes.data.length > 0) {
      areaName = areaRes.data[0].area?.name || areaName;
    }
    conditions.push({
      name: areaName,
      status: currentCondition.status,
      count: res.data.successCount || 0
    });
  }
  currentCondition.conditions = conditions;
  conditionDrawerVisible.value = true;
};
</script>

<style scoped>
/* Add your component styles here */
</style>
