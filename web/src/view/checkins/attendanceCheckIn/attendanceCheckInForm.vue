<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="Ngày giờ:" prop="checkinDate">
          <el-date-picker v-model="formData.checkinDate" type="date" placeholder="选择日期" :clearable="true"></el-date-picker>
       </el-form-item>
        <el-form-item label="Lớp điểm danh:" prop="attendanceId">
          <el-input v-model.number="formData.attendanceId" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="Thành viên:" prop="partpaticipantId">
          <el-input v-model.number="formData.partpaticipantId" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="Khu vực:" prop="areaId">
          <el-input v-model.number="formData.areaId" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="Nhóm:" prop="groupId">
          <el-input v-model.number="formData.groupId" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="Điều kiện:" prop="conditionId">
          <el-input v-model.number="formData.conditionId" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="Kinh độ:" prop="lattidue">
          <el-input-number v-model="formData.lattidue" :precision="2" :clearable="true"></el-input-number>
       </el-form-item>
        <el-form-item label="IP:" prop="iP">
          <el-input v-model.number="formData.iP" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="Vĩ độ:" prop="longtidue">
          <el-input-number v-model="formData.longtidue" :precision="2" :clearable="true"></el-input-number>
       </el-form-item>
        <el-form-item label="Agent:" prop="agent">
          // 此字段为json结构，可以前端自行控制展示和数据绑定模式 需绑定json的key为 formData.agent 后端会按照json的类型进行存取
          {{ formData.agent }}
       </el-form-item>
        <el-form-item label="Trình duyệt:" prop="browser">
          <el-input v-model="formData.browser" :clearable="true"  placeholder="请输入Trình duyệt" />
       </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="save">保存</el-button>
          <el-button type="primary" @click="back">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import {
  createAttendanceCheckIn,
  updateAttendanceCheckIn,
  findAttendanceCheckIn
} from '@/api/checkins/attendanceCheckIn'

defineOptions({
    name: 'AttendanceCheckInForm'
})

// 自动获取字典
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'

const route = useRoute()
const router = useRouter()

const type = ref('')
const formData = ref({
            checkinDate: new Date(),
            attendanceId: undefined,
            partpaticipantId: undefined,
            areaId: undefined,
            groupId: undefined,
            conditionId: undefined,
            lattidue: 0,
            iP: undefined,
            longtidue: 0,
            agent: {},
            browser: '',
        })
// 验证规则
const rule = reactive({
               checkinDate : [{
                   required: true,
                   message: 'Ngày giờ không được để trống',
                   trigger: ['input','blur'],
               }],
               attendanceId : [{
                   required: true,
                   message: 'Lớp điểm danh không được để trống',
                   trigger: ['input','blur'],
               }],
               partpaticipantId : [{
                   required: true,
                   message: 'Thành viên không được để trống',
                   trigger: ['input','blur'],
               }],
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findAttendanceCheckIn({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
}

init()
// 保存按钮
const save = async() => {
      elFormRef.value?.validate( async (valid) => {
         if (!valid) return
            let res
           switch (type.value) {
             case 'create':
               res = await createAttendanceCheckIn(formData.value)
               break
             case 'update':
               res = await updateAttendanceCheckIn(formData.value)
               break
             default:
               res = await createAttendanceCheckIn(formData.value)
               break
           }
           if (res.code === 0) {
             ElMessage({
               type: 'success',
               message: 'Tạo/cập nhật thành công'
             })
           }
       })
}

// 返回按钮
const back = () => {
    router.go(-1)
}

</script>

<style>
</style>
