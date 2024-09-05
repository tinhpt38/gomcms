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
        <div class="text-base my-1 mx-1">Biểu đồ thể hiện tổng số điểm danh và tổng số thành viên điểm danh theo đơn vị và danh mục</div>
        <div class="mt-1 flex flex-row justify-between p-4 rounded bg-slate-100">
            <vue-echarts :option="statsByAgencyCategoryOptions" style="height: 440px; width: 100%;" />
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
    statsByAgencyCategory
} from '@/api/checkins/attendance'

import { onMounted, ref } from 'vue';
import { VueEcharts } from 'vue3-echarts';


defineOptions({
    name: 'CheckinDashboard',
})

// Search Form

const elSearchFormRef = ref()
const searchInfo = ref({
    agencyId: null,
    categoryId: null,
})

const onSubmit = () => {
    console.log('submit', searchInfo.value)
    statsByAgencyCategoryFun()
}

const onReset = () => {
    elSearchFormRef.value.resetFields()
    searchInfo.value = {
        agencyId: null,
        categoryId: null,
    }
}


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


// Data

const statsByAgencyCategoryFun = async () => {
    console.log('statsByAgencyCategory', searchInfo.value)
    const table = await statsByAgencyCategory(searchInfo.value)
    console.log('statsByAgencyCategory', table)
    if (table.code === 0) {
        createStatsByAgencyCategory(table.data.data)
    }

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
        statsByAgencyCategoryFun()
    }
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