<template>
    <div class="p-2">
        <el-tabs v-model="tabsActiveTab" type="border-card">
            <el-tab-pane name="attendanceInfoTab" label="Chi tiết">
                <div class="card-container">
                    <el-row>
                        <el-col :span="12" class="grid-cell"><el-form-item label="Tiêu đề" prop="formData.title"
                                class="required">
                                <el-input v-model="formData.title" type="text" clearable></el-input>
                            </el-form-item></el-col>
                        <el-col :span="12" class="grid-cell flex justify-end">

                            <el-button type="primary" @click="saveAttendance">Lưu</el-button>
                        </el-col>
                    </el-row>

                    <el-row>
                        <el-col :span="12" class="grid-cell">
                            <el-form-item label="Ngày bắt đầu" label-width="150px" prop="startDate" class="required">
                                <el-date-picker v-model="formData.startDate" type="datetime" class="full-width-input"
                                    clearable></el-date-picker>
                            </el-form-item>
                        </el-col>
                        <el-col :span="12" class="grid-cell">
                            <el-form-item label="Ngày kết thúc" label-width="150px" prop="endDate" class="required">
                                <el-date-picker v-model="formData.endDate" type="datetime" class="full-width-input"
                                    clearable></el-date-picker>
                            </el-form-item>
                        </el-col>
                    </el-row>
                    <el-row>
                        <el-col :span="12" class="grid-cell">
                            <el-form-item label="Cho thử nghiệm" label-width="150px" prop="isTrailer">
                                <el-switch v-model="formData.isTrailer"></el-switch>
                            </el-form-item>
                        </el-col>
                        <el-col :span="12" class="grid-cell">
                            <el-form-item label="Khoá" label-width="150px" prop="isLocked">
                                <el-switch v-model="formData.isLocked"></el-switch>
                            </el-form-item>
                        </el-col>
                    </el-row>
                </div>
                <div class="table-container">
                    <table class="table-layout">
                        <tbody>
                            <tr>
                                <td class="table-cell">
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </el-tab-pane>
            <el-tab-pane name="partticipantsTab" label="Thành viên">
                <div class="table-container">
                    <el-table :data="partticipantsData" style="width: 100%">
                        <el-table-column prop="email" label="Email"></el-table-column>
                        <el-table-column prop="group" label="Nhóm"></el-table-column>
                        <el-table-column prop="totalCheckin" label="Checkins">
                            <template #default="scope">
                                <span>{{ scope.row.totalPass }} / {{ scope.row.totalCheckin }}</span>
                            </template>
                        </el-table-column>
                    </el-table>
                    <div class="flex justify-end">
                        <el-pagination v-model:current-page="partsPage" v-model:page-size="partsSize"
                            :page-sizes="[20, 50, 100, 500]" :size="partSize" :background="true"
                            layout="total, sizes, prev, pager, next, jumper" :total="partticipantsData.length"
                            @size-change="partsHandleSizeChange" @current-change="partsHandleCurrentChange" />
                    </div>
                </div>
            </el-tab-pane>
            <el-tab-pane name="groupTan" label="Nhóm">
                <div class="table-container">
                    <table class="table-layout">
                        <tbody>
                            <tr>
                                <td class="table-cell">
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </el-tab-pane>
            <el-tab-pane name="positionTab" label="Vị trí">
                <div class="table-container">
                    <table class="table-layout">
                        <tbody>
                            <tr>
                                <td class="table-cell">
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </el-tab-pane>
            <el-tab-pane name="conditionTab" label="Điều kiện">
                <div class="table-container">
                    <table class="table-layout">
                        <tbody>
                            <tr>
                                <td class="table-cell">
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </el-tab-pane>
        </el-tabs>
    </div>


</template>


<style lang="scss"></style>

<script setup>
import {
    updateAttendanceClass,
    findAttendanceClass,
} from '@/api/checkins/attendanceClass'
import { useRoute } from 'vue-router';

import { getDictFunc, formatDate, formatBoolean, filterDict, filterDataSource, returnArrImg, onDownloadFile } from '@/utils/format'
import { ElForm, ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

defineOptions({
    name: 'AttendanceClassDetail'
})

const $route = useRoute()
const tabsActiveTab = ref('attendanceInfoTab')
const currentId = ref($route.params.id)
const formData = ref({})
const getDetailData = async () => {
    var id = $route.params.id
    console.log('id', id)
    const res = await findAttendanceClass({ id: $route.params.id })
    if (res.code == 0) {
        formData.value = res.data
    }
    console.log('res', res)
}


getDetailData();

const saveAttendance = () => {
    console.log("currentId: ", currentId.value)
}


//region Participants
const partsPage = ref(4)
const partsSize = ref(20)
const partSize = ref(20)


const partticipantsData = ref([])

const initParticipantsData = () => {
    partticipantsData.value = [
        {
            id: 1,
            email: "ex1@example.com",
            group: "Group 1",
            totalCheckin: 10,
            totalPass: 5, // past condition
        },
        {
            id: 2,
            email: "ex2@example.com",
            group: "Group 2",
            totalCheckin: 8,
            totalPass: 3,
        },
        {
            id: 3,
            email: "ex3@example.com",
            group: "Group 3",
            totalCheckin: 12,
            totalPass: 7,
        },
        {
            id: 4,
            email: "ex4@example.com",
            group: "Group 4",
            totalCheckin: 6,
            totalPass: 2,
        },
        {
            id: 5,
            email: "ex5@example.com",
            group: "Group 5",
            totalCheckin: 15,
            totalPass: 10,
        },
        {
            id: 6,
            email: "ex6@example.com",
            group: "Group 6",
            totalCheckin: 9,
            totalPass: 4,
        },
        {
            id: 7,
            email: "ex7@example.com",
            group: "Group 7",
            totalCheckin: 11,
            totalPass: 6,
        },
        {
            id: 8,
            email: "ex8@example.com",
            group: "Group 8",
            totalCheckin: 7,
            totalPass: 3,
        },
        {
            id: 9,
            email: "ex9@example.com",
            group: "Group 9",
            totalCheckin: 13,
            totalPass: 8,
        },
        {
            id: 10,
            email: "ex10@example.com",
            group: "Group 10",
            totalCheckin: 5,
            totalPass: 2,
        },
        {
            id: 11,
            email: "ex11@example.com",
            group: "Group 11",
            totalCheckin: 14,
            totalPass: 9,
        },
        {
            id: 12,
            email: "ex12@example.com",
            group: "Group 12",
            totalCheckin: 7,
            totalPass: 4,
        },
        {
            id: 13,
            email: "ex13@example.com",
            group: "Group 13",
            totalCheckin: 10,
            totalPass: 6,
        },
        {
            id: 14,
            email: "ex14@example.com",
            group: "Group 14",
            totalCheckin: 9,
            totalPass: 5,
        },
        {
            id: 15,
            email: "ex15@example.com",
            group: "Group 15",
            totalCheckin: 12,
            totalPass: 8,
        },
        {
            id: 16,
            email: "ex16@example.com",
            group: "Group 16",
            totalCheckin: 6,
            totalPass: 3,
        },
        {
            id: 17,
            email: "ex17@example.com",
            group: "Group 17",
            totalCheckin: 11,
            totalPass: 7,
        },
        {
            id: 18,
            email: "ex18@example.com",
            group: "Group 18",
            totalCheckin: 8,
            totalPass: 4,
        },
        {
            id: 19,
            email: "ex19@example.com",
            group: "Group 19",
            totalCheckin: 13,
            totalPass: 9,
        },
        {
            id: 20,
            email: "ex20@example.com",
            group: "Group 20",
            totalCheckin: 7,
            totalPass: 5,
        }
    ]
}
initParticipantsData();
const partsHandleCurrentChange = (val) => {
    console.log('partsHandleCurrentChange', val)
}
const partsHandleSizeChange = (val) => {
    console.log('partsHandleSizeChange', val)
}
//endregion
</script>