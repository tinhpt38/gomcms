<template>
    <div class="bg-white rounded">
        <div class="p-2 rounded">
            <h1>Bảng điều kiển</h1>
        </div>
        <div class="gva-search-box border">
            <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline"
                :rules="searchRule" @keyup.enter="onSubmit">
                <el-form-item label="Đơn vị" prop="agencyId" class="">
                    <el-select v-model="searchInfo.agencyId" placeholder="Chọn đơn vị" clearable filterable>
                        <el-option v-for="item in agencyOptions" :key="item.ID" :label="item.name" :value="item.ID" />
                    </el-select>
                </el-form-item>
                <el-form-item label="Danh mục" prop="categoryId" class="w-[300px]">
                    <el-tree-select class="w-full" v-model="searchInfo.categoryId" :data="categoryOptions"
                        check-on-click-node :render-after-expand="false" style="width: 240px" />
                </el-form-item>
                <el-form-item label="Ngày" prop="startAt">
                    <template #label>
                        <span>
                            Ngày tạo
                            <el-tooltip
                                content="Phạm vi tìm kiếm từ ngày bắt đầu (bao gồm) đến ngày kết thúc (không bao gồm)">
                                <el-icon>
                                    <QuestionFilled />
                                </el-icon>
                            </el-tooltip>
                        </span>
                    </template>
                    <el-date-picker v-model="searchInfo.startAt" type="date"
                        placeholder="Ngày bắt đầu"></el-date-picker>
                    —
                    <el-date-picker v-model="searchInfo.endAt" type="date" placeholder="Ngày kết thúc"></el-date-picker>
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" icon="search" @click="onSubmit">
                        Xem dữ liệu
                    </el-button>
                    <el-button icon="refresh" @click="onReset">
                        Đặt lại
                    </el-button>
                </el-form-item>
            </el-form>
        </div>
        <div class="rounded p-4">
            <div class="text-xl my-1 mx-1">Biểu đồ thể hiện tổng số điểm danh và tổng số thành viên điểm danh theo đơn
                vị và danh mục</div>
            <div class="mt-1 flex flex-row justify-between p-4 rounded bg-slate-100">
                <vue-echarts :option="statsByAgencyCategoryOptions" style="height: 440px; width: 100%;" />
            </div>
        </div>

        <div class="rounded p-4">
            <div class="text-xl my-1 mx-1">Thống kê tổng số điểm danh theo tháng hiện tại
                ngày.
            </div>
            <div class="mt-1 flex flex-row justify-between p-4 rounded bg-slate-100">
                <vue-echarts :option="scatterPlotOptions" style="height: 550px; width: 100%;" />
            </div>
        </div>
        <div class="rounded p-4">
            <div class="text-xl my-1 mx-1">Thống kê sự số lượng điểm danh của đơn vị, danh mục theo ngày
            </div>
            <div class="mt-1 flex flex-row justify-between p-4 rounded bg-slate-100">
                <vue-echarts :option="lineChartOptions" style="height: 550px; width: 100%;" />
            </div>
        </div>
    </div>
</template>

<script setup>

import {
    getAttendanceCategoryList
} from '@/api/checkins/attendanceCategory'

import {
    getAttendanceAgencyList
} from '@/api/checkins/attendanceAgency'

import {
    statsByAgencyCategory,
    statsScatterPlot,
    statsTrendLine
} from '@/api/checkins/attendance'

import { onMounted, ref } from 'vue';
import { VueEcharts } from 'vue3-echarts';
import { formatDate } from '@/utils/format'


defineOptions({
    name: 'CheckinDashboard',
})

// Search Form

const elSearchFormRef = ref()
const searchInfo = ref({
    agencyId: null,
    categoryId: null,
    startAt: new Date(new Date().getFullYear(), new Date().getMonth(), 1),
    endAt: new Date(new Date().getFullYear(), new Date().getMonth() + 1, 0),
});

const onSubmit = () => {
    statsFunc()
}

const onReset = () => {
    elSearchFormRef.value.resetFields()
    searchInfo.value = {
        agencyId: null,
        categoryId: null,
    }
}

const timeToDate = (dateString) => {
    const [day, month, year] = dateString.split('/').map(Number);
    return new Date(year, month - 1, day);
};

// Charts 
const statsByAgencyCategoryOptions = ref({})
const createStatsByAgencyCategory = (data) => {

    var xAxisData = data.map((e) => e.AgencyName)
    var seriesDataAgency = data.map((e) => e.CountByAgency)
    var seriesDataParts = data.map((e) => e.TotalParts)
    statsByAgencyCategoryOptions.value = {
        tooltip: {
            trigger: 'axis',
            axisPointer: {
                type: 'shadow'
            }
        },
        grid: {
            left: '3%',
            right: '4%',
            bottom: '3%',
            containLabel: true
        },
        xAxis: [
            {
                type: 'category',
                data: xAxisData,
                axisTick: {
                    alignWithLabel: true
                }
            }
        ],
        yAxis: [
            {
                type: 'value'
            }
        ],
        series: [
            {
                name: 'Tổng số điểm danh',
                type: 'bar',
                barWidth: '50%',
                data: seriesDataAgency
            },
            {
                name: 'Tổng số thành viên',
                type: 'bar',
                barWidth: '50%',
                data: seriesDataParts
            }
        ]
    }
}

const scatterPlotOptions = ref({});

const updateScatterPlotOptions = (data) => {
    const sData = data.map(item => [timeToDate(item.CheckinTime), item.CheckinCount]);
    scatterPlotOptions.value = {
        xAxis: {
            type: 'time',  // Sử dụng trục thời gian
            name: 'Date',
            axisLabel: {
                formatter: (value) => {
                    const date = new Date(value);
                    return `${date.getDate()}/${date.getMonth() + 1}/${date.getFullYear()}`;  // Định dạng hiển thị ngày
                }
            }
        },
        yAxis: {
            type: 'value',  // Trục Y là số lượng điểm danh
            name: 'Tổng số lượt điểm danh'
        },
        series: [
            {
                symbolSize: 20,
                data: sData,  // Dữ liệu sẽ được cập nhật
                type: 'scatter'
            }
        ]
    }
}




const lineChartOptions = ref({});
const updateStatsTrendLineOptions = (data) => {
    // Gom nhóm dữ liệu theo AgencyName và CategoryName
    const groupedData = {};
    data.forEach(item => {
        const key = `${item.AgencyName} - ${item.CategoryName}`;
        if (!groupedData[key]) {
            groupedData[key] = { name: key, data: [], dates: [] };
        }
        groupedData[key].data.push(item.CheckinCount);
        groupedData[key].dates.push(formatDate(item.CheckinDate));
    });
    // Tạo danh sách các ngày cho trục X (xAxis)
    const xAxisData = [...new Set(data.map(item => formatDate(item.CheckinDate)))];

    // Cấu hình biểu đồ
    lineChartOptions.value = {
        
        tooltip: {
            trigger: 'axis'
        },
        legend: {
            data: Object.keys(groupedData)  // Hiển thị các dòng theo AgencyName - CategoryName
        },
        grid: {
            left: '3%',
            right: '4%',
            bottom: '3%',
            containLabel: true
        },
        toolbox: {
            feature: {
                saveAsImage: {}
            }
        },
        xAxis: {
            type: 'category',
            boundaryGap: false,
            data: xAxisData  // Dữ liệu trục X từ CheckinDate
        },
        yAxis: {
            type: 'value'
        },
        series: Object.keys(groupedData).map(key => ({
            name: key,
            type: 'line',
            stack: 'Total',
            data: groupedData[key].data  // Dữ liệu số liệu từ CheckinCount
        }))
    };

}



// Data

const statsByAgencyCategoryFun = async () => {
    console.log('statsByAgencyCategory', searchInfo.value)
    const table = await statsByAgencyCategory(searchInfo.value)
    console.log('statsByAgencyCategory', table)
    if (table.code === 0) {
        createStatsByAgencyCategory(table.data.data)
    }

}

const statsScatterPlotFun = async () => {
    const table = await statsScatterPlot(searchInfo.value)
    console.log('statsScatterPlot', table)
    if (table.code === 0) {
        updateScatterPlotOptions(table.data.data)
    }

}

const statsTrendLineFun = async () => {
    const table = await statsTrendLine(searchInfo.value)
    console.log('statsTrendLine', table)
    if (table.code === 0) {
        updateStatsTrendLineOptions(table.data.data)
    }

}


const statsFunc = async () => {
    statsByAgencyCategoryFun()
    statsScatterPlotFun()
    statsTrendLineFun()
}

// Load Options
const categoryOptions = ref([])
const agencyOptions = ref([])


const getCategoryOptions = async () => {
    const table = await getAttendanceCategoryList({ page: 0, pageSize: -1 })
    if (table.code === 0) {
        categoryOptions.value = convertToTree(table.data.list)
        var selectCurrent = table.data.list.filter((e) => e.isCurrent)
        searchInfo.value.categoryId = selectCurrent[0].ID
    }
    statsFunc()
}

getCategoryOptions()

const getAgencyOptions = async () => {
    const table = await getAttendanceAgencyList({ page: 0, pageSize: -1 })
    if (table.code === 0) {
        agencyOptions.value = table.data.list
    }
}
getAgencyOptions()


const convertToTree = (data) => {
    const map = {}
    const roots = []

    // Create a map of nodes using their ID as the key
    data.forEach((node) => {
        map[node.ID] = { ...node, value: node.ID, label: node.name, children: [] }
    })

    // Iterate over the nodes and assign children to their parent
    data.forEach((node) => {
        const parent = map[node.parentId]

        if (parent) {
            parent.children.push(map[node.ID])
        } else {
            roots.push(map[node.ID])
        }
    })
    console.log('roots', roots)
    return roots
}

</script>

<style></style>