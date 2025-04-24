<template>
  <el-button type="primary" icon="download" @click="exportExcelFunc">Tải xuống</el-button>
</template>

<!-- <script setup>
const props = defineProps({
  templateId: {
    type: String,
    required: true
  },
  condition: {
    type: Object,
    default: () => ({})
  },
  limit: {
    type: Number,
    default: 0
  },
  offset: {
    type: Number,
    default: 0
  },
  order: {
    type: String,
    default: ''
  }
})

import { ElMessage } from 'element-plus'

const exportExcelFunc = async () => {
  if (props.templateId === '') {
    ElMessage.error('Chưa thiết lập ID mẫu cho thành phần')
    return
  }
  const baseUrl = import.meta.env.VITE_BASE_API
  const paramsCopy = JSON.parse(JSON.stringify(props.condition))
  if (props.limit) {
    paramsCopy.limit = props.limit
  }
  if (props.offset) {
    paramsCopy.offset = props.offset
  }
  if (props.order) {
    paramsCopy.order = props.order
  }
  const params = Object.entries(paramsCopy)
    .map(([key, value]) => `${encodeURIComponent(key)}=${encodeURIComponent(value)}`)
    .join('&')
  const url = `${baseUrl}/sysExportTemplate/exportExcel?templateID=${props.templateId}${params ? '&' + params : ''}`

  window.open(url, '_blank')
}
</script> -->

<!-- <script setup>
import { exportAttendanceData } from './exportExcel' // Import hàm exportAttendanceData
import { ref } from 'vue'

const searchInfo = ref({
  attendanceId: 1, // Ví dụ, có thể lấy từ props hoặc từ context
})

const page = ref(1)
const pageSize = ref(10)

const exportExcelFunc = async () => {
  console.log('Tải xuống file excel')
  // Gọi hàm exportAttendanceData khi người dùng bấm tải xuống
  await exportAttendanceData(searchInfo.value, page.value, pageSize.value)
}
</script> -->

<script setup>
  const props = defineProps({
    formData: {
    type: Object,
    required: true
  }
})
  import { ElMessage } from 'element-plus'
  import { exportAttendanceData } from './exportExcel' 

  const exportExcelFunc = async () => {
   const attendanceId = props.formData.attendanceId;
   console.log("Form data received: ", attendanceId);
   await exportAttendanceData(attendanceId, 1, 500) // Gọi hàm exportAttendanceData với page = 1 và pageSize = 10
   ElMessage.success('Tải xuống file excel thành công')
}
</script>