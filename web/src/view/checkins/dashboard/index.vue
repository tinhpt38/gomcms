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
    </div>
</template>

<script setup>

import {
    getAttendanceCategoryList
} from '@/api/checkins/attendanceCategory'

import {
    getAttendanceAgencyList
} from '@/api/checkins/attendanceAgency'
import { ref } from 'vue';


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
}

const onReset = () => {
    elSearchFormRef.value.resetFields()
    searchInfo.value = {
        agencyId: null,
        categoryId: null,
    }
}

// Load Options
const categoryOptions = ref([])
const agencyOptions = ref([])


const getCategoryOptions = async () => {
    const table = await getAttendanceCategoryList({ page: 0, pageSize: -1 })
    if (table.code === 0) {
        categoryOptions.value = convertToTree(table.data.list)
    }
    // console.log("parent Options", categoryOptions.value)
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