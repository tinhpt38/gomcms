<template>
    <div>
        <div class="p-1 my-1">
            <el-button type="primary" icon="plus" @click="openDialog()">Thêm nhóm</el-button>
        </div>

        <el-table :data="tableData" style="width: 100%">
            <el-table-column prop="name" label="Nhóm"></el-table-column>
            <!-- TODO: Handle totalParts after -->
            <el-table-column prop="totalParts" label="Số thành viên"></el-table-column>
            <el-table-column label="Hành động">
                <template #default="scope">
                    <el-button type="primary" plain round @click="updateGroupFunc(scope.row)">Sửa</el-button>
                    <el-button type="danger" plain round @click="deleteRow(scope.row)">Xoá</el-button>
                </template>
            </el-table-column>
        </el-table>
        <div class="flex justify-end">
            <el-pagination v-model:current-page="page" v-model:page-size="pageSize" :page-sizes="[20, 50, 100, 500]"
                :size="size" :background="true" layout="total, sizes, prev, pager, next, jumper" :total="total"
                @size-change="handleSizeChange" @current-change="handleCurrentChange" />
        </div>

    </div>

    <el-drawer destroy-on-close size="800" v-model="dialogFormVisible" :show-close="false" :before-close="closeDialog">
        <template #header>
            <div class="flex justify-between items-center">
                <span class="text-lg">{{ type === 'create' ? 'Thêm mới' : 'Chỉnh sửa' }}</span>
                <div>
                    <el-button type="primary" @click="enterDialog">Đồng ý</el-button>
                    <el-button @click="closeDialog">Huỷ</el-button>
                </div>
            </div>
        </template>

        <el-form :model="formData" label-position="top" ref="elFormRef" :rules="rule" label-width="80px">
            <el-form-item label="Tên nhóm:" prop="name">
                <el-input v-model="formData.name" :clearable="true" placeholder="Vui lòng nhập Tên nhóm" />
            </el-form-item>
            <el-form-item label="Attendance Class:" prop="attendanceClassId" class="hidden">
                <el-select v-model="formData.attendanceClassId" placeholder="Vui lòng chọn Attendance Class"
                    style="width:100%" :clearable="true">
                    <el-option v-for="(item, key) in dataSource.attendanceClassId" :key="key" :label="item.label"
                        :value="item.value" />
                </el-select>
            </el-form-item>
        </el-form>
    </el-drawer>
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
    getGroupList
} from '@/api/checkins/group'


const props = defineProps({
    acId: {
        type: Number,
        required: true
    }
})

const tableData = ref([])
const searchInfo = ref({
    attendanceClassId: null
})
const type = ref('')

const page = ref(1)
const pageSize = ref(20)
const size = ref(20)
const total = ref(0)

const getTableData = async () => {
    searchInfo.value.attendanceClassId = props.acId
    const table = await getGroupList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
    if (table.code === 0) {
        tableData.value = table.data.list
        total.value = table.data.total
        page.value = table.data.page
        pageSize.value = table.data.pageSize
    }
}
const rule = reactive({
    name: [{
        required: true,
        message: 'Tên nhóm là bắt buộc',
        trigger: ['input', 'blur'],
    },
    {
        whitespace: true,
        message: 'Bạn không thể nhập mỗi khoảng trắng',
        trigger: ['input', 'blur'],
    }
    ],
})

const formData = ref({
    name: '',
    attendanceClassId: props.acId *1,
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

const deleteRow = (row) => {
    ElMessageBox.confirm('Bạn có chắc muốn xoá không?', 'Cảnh báo', {
        confirmButtonText: 'Đồng ý',
        cancelButtonText: 'Hủy',
        type: 'warning'
    }).then(() => {
        deleteGroupFunc(row)
    })
}
const dialogFormVisible = ref(false)
const openDialog = () => {
    type.value = 'create'
    dialogFormVisible.value = true
}

const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = {
        name: '',
        attendanceClassId: props.acId * 1,
    }
}
const elFormRef = ref()


const enterDialog = async () => {
    elFormRef.value?.validate(async (valid) => {
        if (!valid) return
        let res
        switch (type.value) {
            case 'create':
                res = await createGroup(formData.value)
                break
            case 'update':
                res = await updateGroup(formData.value)
                break
            default:
                res = await createGroup(formData.value)
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

const onDelete = async () => {
    ElMessageBox.confirm('Bạn có chắc muốn xoá không?', 'Cảnh báo', {
        confirmButtonText: 'Đồng ý',
        cancelButtonText: 'Hủy',
        type: 'warning'
    }).then(async () => {
        const IDs = []
        if (multipleSelection.value.length === 0) {
            ElMessage({
                type: 'warning',
                message: 'Vui lòng chọn dữ liệu để xoá'
            })
            return
        }
        multipleSelection.value &&
            multipleSelection.value.map(item => {
                IDs.push(item.ID)
            })
        const res = await deleteGroupByIds({ IDs })
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


const deleteGroupFunc = async (row) => {
    const res = await deleteGroup({ ID: row.ID })
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