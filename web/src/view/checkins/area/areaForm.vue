<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="Tên khu vực:" prop="name">
          <el-input v-model="formData.name" :clearable="true"  placeholder="请输入Tên khu vực" />
       </el-form-item>
        <el-form-item label="Latitude:" prop="latitude">
          <el-input-number v-model="formData.latitude" :precision="2" :clearable="true"></el-input-number>
       </el-form-item>
        <el-form-item label="Longitude:" prop="longitude">
          <el-input-number v-model="formData.longitude" :precision="2" :clearable="true"></el-input-number>
       </el-form-item>
        <el-form-item label="Bán kính:" prop="radius">
          <el-input-number v-model="formData.radius" :precision="2" :clearable="true"></el-input-number>
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
  createArea,
  updateArea,
  findArea
} from '@/api/checkins/area'

defineOptions({
    name: 'AreaForm'
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
            latitude: 0,
            longitude: 0,
            radius: 0,
        })
// 验证规则
const rule = reactive({
               name : [{
                   required: true,
                   message: 'Tên khu vực là bắt buộc ',
                   trigger: ['input','blur'],
               }],
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findArea({ ID: route.query.id })
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
               res = await createArea(formData.value)
               break
             case 'update':
               res = await updateArea(formData.value)
               break
             default:
               res = await createArea(formData.value)
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
