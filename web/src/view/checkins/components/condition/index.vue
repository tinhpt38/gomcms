<template>
    <div>
        <div class="p-1 my-1">
            <el-button type="primary" icon="plus" @click="addNewCondition()">Thêm điều kiện</el-button>
        </div>

        <el-table :data="conditions" style="width: 100%">
            <el-table-column prop="areaId" label="Khu vực"></el-table-column>
            <el-table-column prop="groupId" label="Nhóm"></el-table-column>
            <el-table-column prop="startAt" label="Bắt đầu"></el-table-column>
            <el-table-column prop="endAt" label="Kết thúc"></el-table-column>
            <el-table-column label="Hành động">
                <template #default="scope">
                    <el-button type="primary" plain round>Sửa</el-button>
                    <el-button type="danger" plain round>Xoá</el-button>
                </template>
            </el-table-column>
        </el-table>
        <div class="flex justify-end">
            <el-pagination v-model:current-page="page" v-model:page-size="pageSize" :page-sizes="[20, 50, 100, 500]"
                :size="size" :background="true" layout="total, sizes, prev, pager, next, jumper" :total="conditions.length"
                @size-change="handleSizeChange" @current-change="handleCurrentChange" />
        </div>

    </div>
</template>

<script>
export default {
    props: {
        conditions: {
            type: Array,
            required: true
        }
    }
}

</script>

<script setup>
import { ref } from 'vue'


const addNewCondition = () => {
    console.log('addNewCondition')
}



const page = ref(4)
const pageSize = ref(20)
const size = ref(20)


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