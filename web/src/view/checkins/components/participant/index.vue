<template>
    <el-table :data="participants" style="width: 100%">
        <el-table-column prop="email" label="Email"></el-table-column>
        <el-table-column prop="group" label="Nhóm"></el-table-column>
        <el-table-column prop="totalCheckin" label="Checkins">
            <template #default="scope">
                <span>{{ scope.row.totalPass }} / {{ scope.row.totalCheckin }}</span>
            </template>
        </el-table-column>
    </el-table>
    <div class="flex justify-end">
        <el-pagination v-model:current-page="page" v-model:page-size="pageSize" :page-sizes="[20, 50, 100, 500]"
            :size="size" :background="true" layout="total, sizes, prev, pager, next, jumper"
            :total="participants.length" @size-change="handleSizeChange" @current-change="handleCurrentChange" />
    </div>
</template>

<script>
export default {
    props: {
        participants: {
            type: Array,
            required: true
        }
    }
}

</script>

<script setup>
import { ref } from 'vue'

const page = ref(4)
const pageSize = ref(20)
const size = ref(20)


const handleCurrentChange = (val) => {
    console.log('partsHandleCurrentChange', val)
}

const handleSizeChange = (val) => {
    console.log('partsHandleSizeChange', val)
}
</script>

<style scoped>
/* Add your component styles here */
</style>