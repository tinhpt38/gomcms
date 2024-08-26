<template>
  <div>
    <div class="p-1 my-1">
      <!-- <el-button type="primary" icon="plus" @click="() => {}">
        Thêm thành viên
      </el-button> -->
    </div>
    <el-table :data="tableData" style="width: 100%">
      <el-table-column prop="fullName" label="Họ và tên" />
      <el-table-column prop="email" label="Email" />
      <el-table-column label="Nhóm">
        <template #default="scope">
          <span>{{ scope.row.groups.map((e) => e.name).join(", ") }}</span>
        </template>
      </el-table-column>
      <!-- <el-table-column prop="totalCheckin" label="Checkins">
        <template #default="scope">
          <span>{{ scope.row.totalPass }} / {{ scope.row.totalCheckin }}</span>
        </template>
      </el-table-column> -->
    </el-table>
    <div class="flex justify-end">
      <el-pagination
        v-model:current-page="page" v-model:page-size="pageSize" :page-sizes="[20, 50, 100, 500]"
        :size="size" :background="true" layout="total, sizes, prev, pager, next, jumper"
        :total="total" @size-change="handleSizeChange" @current-change="handleCurrentChange"
      />
    </div>
  </div>
</template>

<script setup>
import {
  createParticipant,
  deleteParticipant,
  deleteParticipantByIds,
  getParticipantListByAttendance
} from '@/api/checkins/participant'
import { getDictFunc, formatDate, formatBoolean, filterDict ,filterDataSource, returnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

const props = defineProps({
    acId: {
        type: Number,
        required: true
    }
})
const page = ref(0)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})


const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

const getTableData = async() => {
  searchInfo.value.attendanceId = props.acId
  const table = await getParticipantListByAttendance({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
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

const deleteRow = (row) => {
    ElMessageBox.confirm('Bạn có chắc muốn xóa không?', 'Cảnh báo', {
      confirmButtonText: 'Đồng ý',
      cancelButtonText: 'Hủy',
      type: 'warning'
    }).then(() => {
      deleteParticipantFunc(row)
    })
  }

const onDelete = async() => {
  ElMessageBox.confirm('Bạn có chắc muốn xóa không?', 'Cảnh báo', {
    confirmButtonText: 'Đồng ý',
    cancelButtonText: 'Hủy',
    type: 'warning'
  }).then(async() => {
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


const updateParticipantFunc = async(row) => {
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

const enterDialog = async () => {
     elFormRef.value?.validate( async (valid) => {
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



</script>

<style scoped>
/* Add your component styles here */
</style>