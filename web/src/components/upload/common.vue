<template>
  <div>
    <el-upload
      :action="`${getBaseUrl()}/fileUploadAndDownload/upload`"
      :before-upload="checkFile"
      :on-error="uploadError"
      :on-success="uploadSuccess"
      :show-file-list="false"
      class="upload-btn"
    >
      <el-button type="primary">Tải lên thông thường</el-button>
    </el-upload>
  </div>
</template>

<script setup>

import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { isVideoMime, isImageMime } from '@/utils/image'
import { getBaseUrl } from '@/utils/format'

defineOptions({
  name: 'UploadCommon',
})

const emit = defineEmits(['on-success'])
const path = ref(import.meta.env.VITE_BASE_API)

const fullscreenLoading = ref(false)

const checkFile = (file) => {
  fullscreenLoading.value = true
  const isLt500K = file.size / 1024 / 1024 < 1.5 // 500K, @todo should be configurable in the project
  const isLt5M = file.size / 1024 / 1024 < 5 // 5MB, @todo should be configurable in the project
  const isVideo = isVideoMime(file.type)
  const isImage = isImageMime(file.type)
  let pass = true
  if (!isVideo && !isImage) {
    ElMessage.error('Chỉ có thể tải lên hình ảnh định dạng jpg, png, svg, webp hoặc video định dạng mp4, webm!')
    fullscreenLoading.value = false
    pass = false
  }
  if (!isLt5M && isVideo) {
    ElMessage.error('Kích thước video tải lên không được vượt quá 5MB')
    fullscreenLoading.value = false
    pass = false
  }
  if (!isLt500K && isImage) {
    ElMessage.error('Kích thước hình ảnh tải lên chưa được nén không được vượt quá 500KB, vui lòng sử dụng chức năng nén trước khi tải lên')
    fullscreenLoading.value = false
    pass = false
  }

  console.log('Kết quả kiểm tra tệp tải lên: ', pass)

  return pass
}

const uploadSuccess = (res) => {
  const { data } = res
  if (data.file) {
    emit('on-success', data.file.url)
  }
}

const uploadError = () => {
  ElMessage({
    type: 'error',
    message: 'Tải lên thất bại'
  })
  fullscreenLoading.value = false
}
</script>