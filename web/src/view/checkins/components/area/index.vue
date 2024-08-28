<template>
  <div>
    <div class="p-1 my-1">
      <el-button type="primary" icon="plus" @click="addNewArea()">
        Thêm khu vực
      </el-button>
    </div>

    <el-table :data="tableData" style="width: 100%" border>
      <el-table-column prop="area.name" label="Nhóm" />
      <el-table-column prop="area.latitude" label="Kinh độ" />
      <el-table-column prop="area.longitude" label="Vĩ độ" />
      <el-table-column prop="" label="Bán kính">
        <template #default="scope">
          {{ scope.row.radius ?? scope.row.area.radius }}
        </template>
      </el-table-column>
      <el-table-column label="Hành động">
        <template #default="scope">
          <el-button type="danger" plain round @click="onDeleteArea(scope.row.ID)">
            Xoá
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    <div class="flex justify-end">
      <el-pagination
        v-model:current-page="page" v-model:page-size="pageSize" :page-sizes="[20, 50, 100, 500]"
        :size="size" :background="true" layout="total, sizes, prev, pager, next, jumper" :total="total"
        @size-change="handleSizeChange" @current-change="handleCurrentChange"
      />
    </div>
    <el-drawer
      v-model="dialogFormVisible" destroy-on-close size="800" :show-close="false"
      :before-close="closeDialog"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? 'Thêm mới' : 'Chỉnh sửa' }}</span>
          <div>
            <el-button type="primary" @click="enterDialog">
              Đồng ý
            </el-button>
            <el-button @click="dialogFormVisible = false">
              Hủy
            </el-button>
          </div>
        </div>
      </template>

      <el-form ref="elFormRef" :model="formData" label-position="top" :rules="rules" label-width="80px">
        <el-form-item label="Tên khu vực:" prop="areaId">
          <el-select
            v-model="formData.areaId" placeholder="Chọn khu vực" clearable filterable
            @change="onSelectChange"
          >
            <el-option v-for="item in areaOptions" :key="item.ID" :label="item.name" :value="item.ID" />
          </el-select>
        </el-form-item>
        <el-form-item label="Bán kính:" prop="radius">
          <el-input-number v-model="formData.radius" style="width:100%" :precision="2" :clearable="true" />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import {
    deleteArea,
    getAreaList
} from '@/api/checkins/area'
import {
    createAttendanceArea,
    deleteAttendanceArea,
    findAttendanceArea
} from '@/api/checkins/attendance'
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus';

const props = defineProps(
    {
        acId: {
            type: Array,
            required: true
        }
    }
)
const elFormRef = ref()
const dialogFormVisible = ref(false)
const type = ref('')
const formData = ref({
    attendanceId: props.acId,
    areaId: null,
    radius: null,
})

const rules = ref({
    areaId: [
        { required: true, message: 'Vui lòng chọn khu vực', trigger: 'change' }
    ]
})

const tableData = ref([])
const total = ref(0)

const page = ref(4)
const pageSize = ref(20)
const size = ref(20)


const searchInfo = ref({
    attendanceId: null
})

const getAreaListData = async () => {
    searchInfo.value.attendanceId = props.acId
    const table = await findAttendanceArea({ id: props.acId })
    console.log('getAreaListData', table)
    if (table.code === 0) {
        tableData.value = table.data
        total.value = table.data.total
    }

}

getAreaListData();

const onDeleteArea = async (id) => {
    if (!id) return

    ElMessageBox.confirm('Khi xoá khu vực đồng thời cũng sẽ xoá điều kiện tham chiếu', 'Cảnh báo', {
        confirmButtonText: 'Đồng ý',
        cancelButtonText: 'Hủy',
        type: 'warning'
    }).then(async() => {
        const res = await deleteAttendanceArea({ id: id })
        if (res.code === 0) {
            ElMessage({
                type: 'success',
                message: 'Xoá thành công'
            })
            getAreaListData();
        }
    }).catch(() => {
        ElMessage({
            type: 'info',
            message: 'Đã hủy'
        })
    })

}

const addNewArea = () => {
    type.value = 'create'
    dialogFormVisible.value = true
}

const enterDialog = () => {
    elFormRef.value?.validate(async (valid) => {
        if (!valid) return
        let res
        formData.value.attendanceId = props.acId * 1
        switch (type.value) {
            case 'create':
                res = await createAttendanceArea(formData.value)
                break
            default:
                break
        }
        if (res.code === 0) {
            ElMessage({
                type: 'success',
                message: 'Tạo/cập nhật thành công'
            })
            dialogFormVisible.value = false;
            getAreaListData();
        }
    })
}

const onSelectChange = (val) => {
    var item = areaOptions.value.find(item => item.ID === val)
    if (item === undefined) return
    formData.value.radius = item.radius
}

const areaOptions = ref([])

const getAreaOptions = async () => {
    const table = await getAreaList({ page: 1, pageSize: 9999 })
    if (table.code === 0) {
        areaOptions.value = table.data.list
    }
    console.log('areaOptions', areaOptions)
}

getAreaOptions();



const handleCurrentChange = (val) => {
    console.log('handleCurrentChange', val)
}

const handleSizeChange = (val) => {
    console.log('partsHahandleSizeChangendleSizeChange', val)
}


</script>

<style scoped>
/* Add your component styles here */
</style>