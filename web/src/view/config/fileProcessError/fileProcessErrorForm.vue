<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="Mã định danh quá trình xử lý tệp tin:" prop="fileProcessId">
          <el-input v-model.number="formData.fileProcessId" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="Mã định danh của tệp tin:" prop="fileProcessUuid">
          <el-input v-model="formData.fileProcessUuid" :clearable="true"  placeholder="请输入Mã định danh của tệp tin" />
       </el-form-item>
        <el-form-item label="Cột dữ liệu :" prop="fieldTitle">
          <el-input v-model="formData.fieldTitle" :clearable="true"  placeholder="请输入Cột dữ liệu " />
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
  createFileProcessError,
  updateFileProcessError,
  findFileProcessError
} from '@/api/config/fileProcessError'

defineOptions({
    name: 'FileProcessErrorForm'
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
            fileProcessId: undefined,
            fileProcessUuid: '',
            fieldTitle: '',
        })
// 验证规则
const rule = reactive({
               fileProcessId : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               fileProcessUuid : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               fieldTitle : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findFileProcessError({ ID: route.query.id })
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
               res = await createFileProcessError(formData.value)
               break
             case 'update':
               res = await updateFileProcessError(formData.value)
               break
             default:
               res = await createFileProcessError(formData.value)
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
