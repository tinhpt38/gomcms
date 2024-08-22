<template>
    <div class="tab-container">
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

const saveAttendance = ()=>{
    console.log("currentId: ", currentId.value)
}

</script>