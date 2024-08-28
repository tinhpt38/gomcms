<template>
  <div>
    <el-button type="success" icon="document-add" @click="openDialog">
      Nhập dữ liệu từ excel
    </el-button>

    <el-dialog v-model="dialogVisible" title="Nhập dữ liệu từ excel" width="60%" :before-close="closeDialog">
      <el-form v-if="isShow" ref="formRef" class="doit-form" :model="data">
        <el-row justify="space-around">
          <el-col :span="11">
            <el-form-item>
              <el-input-number v-model="data.month" placeholder="Tháng" class="w-100" :min="1" :max="12" />
            </el-form-item>
          </el-col>
          <el-col :span="11">
            <el-form-item>
              <el-input-number v-model="data.year" placeholder="Năm" class="w-100" :min="1000" :max="9999" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>

      <el-upload
        ref="upload" class="upload-demo" drag action="#"
        accept=".csv, application/vnd.openxmlformats-officedocument.spreadsheetml.sheet, application/vnd.ms-excel"
        :limit="1" :auto-upload="false" :on-remove="handleRemove" :on-exceed="handleLogoUploaderExceed"
        @change="handleFileChange"
      >
        <el-icon class="el-icon--upload">
          <upload-filled />
        </el-icon>
        <div class="el-upload__text">
          Kéo thả tập tin Excel vào đây <em>hoặc chọn từ máy tính</em>
        </div>
        <template #tip>
          <div class="el-upload__tip">
            File .xls, .xlsx
          </div>
        </template>
      </el-upload>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">
            Quay về
          </el-button>
          <el-button type="primary" @click="submitUpload">
            Lưu vào hệ thống
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: 'ImportExcel'
}
</script>

<script setup>
// Import
import { useConfirmBox, useNotification } from '@/use'
import { ref, onMounted } from 'vue'

import { createCfgFileProcess } from '@/api/config/cfg_file_process'
import { ElNotification } from 'element-plus'

// Use pinia
// Import strore
import { useFileProcessStore } from '@/pinia/bl'

const fileProcessStore = useFileProcessStore()

// props
const props = defineProps({
  formData: {
    type: Object,
    default() {
      return {
        action: ''
      }
    }
  },

})

onMounted(() => {
  if (props.formData.action === 'import_environment_entry') {
    isShow.value = true
  }
})

// Data
const dialogVisible = ref(false)
const upload = ref()
const selectedFile = ref(null)
const confirmBox = useConfirmBox()
const notification = useNotification()
const data = ref({})
const isShow = ref(false)

// Method
const handleLogoUploaderExceed = (newFiles) => {
  upload.value.clearFiles()
  const file = newFiles[0]
  upload.value.handleStart(file)
}

const openDialog = () => {
  dialogVisible.value = true
}

const closeDialog = () => {
  dialogVisible.value = false
  // Reset lại giá trị của
  upload.value?.clearFiles()
  data.value = {}
  selectedFile.value = null
}

const handleFileChange = (file) => {
  if (file.status === 'ready') {
    // Lưu file đã chọn vào biến selectedFile
    selectedFile.value = file.raw
    return
  }
}

const handleRemove = (file) => {
  selectedFile.value = null
}

const submitUpload = () => {
  if (selectedFile.value === null) {
    ElNotification(
      {
        title: 'Nhập dữ liệu từ file excel',
        message: 'Bạn chưa chọn file excel nào, vui lòng chọn file excel trước khi nhập dữ liệu!',
        duration: 2000,
        position: 'top-right',
      },
    )
    return
  }

  confirmBox.dialogConfirm('Bạn có chắc chắn muốn nhập dữ liệu từ file excel này không?', () => {
    const file = selectedFile.value

    const formData = new FormData()
    formData.append('file', file)
    formData.append('type', 'import')
    formData.append('filter', JSON.stringify(data.value))


    for (const k in props.formData) {
      formData.append(k, props.formData[k])
    }

    createCfgFileProcess(formData)
      .then((res) => {
        if (res.code === 0) {
          fileProcessStore.addFileProcess(res.data)
          closeDialog()
        } else if (res.code === 2) {
          notification.showErrorNotification(res.msg)
        }
      })
      .catch((err) => {
        console.error(err)
        notification.showErrorNotification('Nhập dữ liệu từ file excel thất bại!')
      })
  })
}
// Watch
</script>

<style lang="sass" scoped>
</style>
