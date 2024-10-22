<template>
    <div>
        <div class="my-4">
            <el-form label-position="top" :inline="true" :model="searchInfo" class="demo-form-inline"
                @keyup.enter="onSubmit">
                <el-form-item label="Email" prop="email">
                    <el-input v-model="searchInfo.email" type="text" placeholder="Email" />
                </el-form-item>
                <el-form-item label="Hành động">
                    <el-button type="primary" icon="search" @click="onSubmit">
                        Tìm kiếm
                    </el-button>
                    <el-button icon="refresh" @click="onReset">
                        Đặt lại
                    </el-button>
                </el-form-item>
            </el-form>
        </div>
        <div class="mt-4">
            <el-table style="width: 100%" tooltip-effect="dark" :data="tableData" row-key="ID" border>
                <el-table-column align="left" label="Ngày giờ" prop="CreatedAt" width="180">
                    <template #default="scope">
                        {{ formatDateTime(scope.row.CreatedAt) }}
                    </template>
                </el-table-column>
                <el-table-column align="left" label="Họ tên" prop="fullName" width="200" />
                <el-table-column align="left" label="Email" prop="email" width="200"/>
                <el-table-column align="left" label="IP" prop="ip" width="140" />
                <el-table-column align="left" label="Vị trí" width="170">
                    <template #default="scope">
                        <a v-if="scope.row.lattidue" target="_blank"
                            :href="'https://www.google.com/maps?q=' + scope.row.lattidue + ',' + scope.row.longtidue">{{
                                scope.row.lattidue }}, {{ scope.row.longtidue }}</a>
                    </template>
                </el-table-column>
                <el-table-column align="left" label="Độ lệch" prop="accuracy" width="120">
                    <template #default="scope">
                        {{ scope.row.accuracy.toFixed(4) }}m
                    </template>
                </el-table-column>
                <el-table-column align="left" label="Client ID" prop="visitorId" width="160">
                    <template #default="scope">
                        <span>{{ scope.row.visitorId?.substring(0, 15) ?? '/' }}...</span>
                    </template>
                </el-table-column>
                <el-table-column label="Agent" prop="agent">
                    <template #default="scope">
                        {{ formatUserAgent(scope.row.agent) }}
                    </template>
                </el-table-column>
            </el-table>
            <div class="gva-pagination">
                <el-pagination layout="total, sizes, prev, pager, next, jumper" :current-page="page"
                    :page-size="pageSize" :page-sizes="[10, 30, 50, 100]" :total="total"
                    @current-change="handleCurrentChange" @size-change="handleSizeChange" />
            </div>
        </div>
    </div>
</template>

<script setup>
import {
    getAttendanceCheckInLogList
} from '@/api/checkins/attendanceCheckIn'
import { formatDateTime } from '@/utils/format'
import { formatUserAgent } from '@/utils/userAgent'
import { ref } from 'vue'

const props = defineProps({
    acId: {
        type: Number,
        required: true
    },
})

const page = ref(0)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])

const formData = ref({
})

const searchInfo = ref({
    email: null,
    attendanceId: props.acId,
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
    const table = await getAttendanceCheckInLogList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
    console.log(table)
    if (table.code === 0) {
        tableData.value = table.data.list
        total.value = table.data.total
        page.value = table.data.page
        pageSize.value = table.data.pageSize
    }
}

getTableData()

const onSubmit = () => {
    getTableData()
}

const onReset = () => {
    searchInfo.value = {
        email: null,
        attendanceId: props.acId
    }
    getTableData()
}
</script>