<template>
    <div>
        <!-- Fail reasons pie chart -->
        <div v-if="failReasons.length" class="mb-6 bg-white rounded shadow-sm p-4">
            <div class="font-semibold text-gray-700 mb-3">Lý do điểm danh thất bại</div>
            <vue-echarts :option="pieOption" style="height: 220px; width: 100%;" autoresize />
        </div>

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
                    <el-button type="success" icon="download" :loading="exportLoading" @click="exportLogs">
                        Xuất Excel
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
                        <a v-if="scope.row.lat" target="_blank"
                            :href="'https://www.google.com/maps?q=' + scope.row.lat + ',' + scope.row.lng">
                            {{ scope.row.lat }}, {{ scope.row.lng }}
                        </a>
                    </template>
                </el-table-column>
                <el-table-column align="left" label="Độ lệch" prop="accuracy" width="120">
                    <template #default="scope">
                        {{ scope.row.accuracy?.toFixed(4) ?? '—' }}m
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
    getAttendanceCheckInLogList,
    exportAttendanceCheckInLogExcel,
} from '@/api/checkins/attendanceCheckIn'
import { getAttendanceOverview } from '@/api/checkins/attendance'
import { formatDateTime } from '@/utils/format'
import { formatUserAgent } from '@/utils/userAgent'
import { ref, computed, onMounted } from 'vue'
import { VueEcharts } from 'vue3-echarts'

const props = defineProps({
    acId: {
        type: Number,
        required: true
    },
    title: {
        type: String,
        default: '',
    },
})

const page = ref(0)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const exportLoading = ref(false)

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

const exportLogs = async () => {
    exportLoading.value = true
    try {
        await exportAttendanceCheckInLogExcel({
            attendanceId: props.acId,
            email: searchInfo.value.email || undefined,
        }, props.title)
    } finally {
        exportLoading.value = false
    }
}

// Fail reasons pie (lightweight — reuse overview endpoint)
const failReasons = ref([])
const pieOption = computed(() => ({
    tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
    legend: { bottom: 0, type: 'scroll' },
    series: [{
        type: 'pie',
        radius: ['35%', '65%'],
        label: { show: false },
        emphasis: { label: { show: true, fontSize: 12, fontWeight: 'bold' } },
        data: failReasons.value.map(r => ({ name: r.reason, value: r.count })),
    }],
}))

onMounted(async () => {
    try {
        const res = await getAttendanceOverview({ attendanceId: props.acId })
        if (res.code === 0 && res.data?.failReasons?.length) {
            failReasons.value = res.data.failReasons
        }
    } catch (_) {
        // silently ignore — chart is optional
    }
})
</script>