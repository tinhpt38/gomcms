<template>
    <div>
        <el-button type="success" @click="onLuckyClick" class="mt-2 mb-4">
            Tìm người may mắn
        </el-button>
        <div class="flex flex-row justify-center align-middle w-full">
            <div class="rounded-md p-4 border">
                <p class="text-8xl p-2 text-center">{{ luckyMember.email }}</p>
                <p class="text-6xl p-2 text-center">{{ luckyMember.fullName?.replaceAll("undefined", "") ?? "--" }}</p>
                
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import {
    findLuckyParticipant,
} from '@/api/checkins/participant'

const props = defineProps({
    acId: {
        type: Number,
        required: true
    },
})


const searchInfo = ref({
    attendanceId: props.acId,
})

const luckyMember = ref({
    email: null,
    fullName: null,
})


const onLuckyClick = async () => {
    try {
        searchInfo.value.attendanceId = props.acId
        const res = await findLuckyParticipant(searchInfo.value)
        luckyMember.value = res.data
    } catch (error) {
        ElMessage.error('Có lỗi xảy ra')
    }
}


</script>