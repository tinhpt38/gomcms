<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="Tên nhóm:" prop="name">
          <el-input v-model="formData.name" :clearable="true"  placeholder="请输入Tên nhóm" />
       </el-form-item>
        <el-form-item label="Attendance Class:" prop="attendanceClassId">
        <el-select  v-model="formData.attendanceClassId" placeholder="请选择Attendance Class" style="width:100%" :clearable="true" >
          <el-option v-for="(item,key) in dataSource.attendanceClassId" :key="key" :label="item.label" :value="item.value" />
        </el-select>
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
    getGroupDataSource,
  createGroup,
  updateGroup,
  findGroup
} from '@/api/checkins/group'

defineOptions({
    name: 'GroupForm'
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
            name: '',
            attendanceClassId: undefined,
        })
// 验证规则
const rule = reactive({
               name : [{
                   required: true,
                   message: 'Tên nhóm là bắt buộc',
                   trigger: ['input','blur'],
               }],
})

const elFormRef = ref()
  const dataSource = ref([])
  const getDataSourceFunc = async()=>{
    const res = await getGroupDataSource()
    if (res.code === 0) {
      dataSource.value = res.data
    }
  }
  getDataSourceFunc()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findGroup({ ID: route.query.id })
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
               res = await createGroup(formData.value)
               break
             case 'update':
               res = await updateGroup(formData.value)
               break
             default:
               res = await createGroup(formData.value)
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
