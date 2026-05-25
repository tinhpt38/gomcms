import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/pinia/modules/user'

const downloadBlob = (blob, fileName) => {
  const url = window.URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = fileName
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  window.URL.revokeObjectURL(url)
}

export const buildExportFileName = (prefix, title, id) => {
  const safeTitle = String(title || id || 'export')
    .trim()
    .replace(/[/\\:*?"<>|]/g, '-')
    .replace(/\s+/g, '-')
    .slice(0, 80)

  return `${prefix}-${safeTitle || id}.xlsx`
}

export const downloadExcelFile = async (url, params, fileName) => {
  const userStore = useUserStore()

  try {
    const response = await axios({
      baseURL: import.meta.env.VITE_BASE_API,
      url,
      method: 'get',
      params,
      responseType: 'blob',
      headers: {
        'x-token': userStore.token,
        'x-user-id': userStore.userInfo.ID,
      },
    })

    const contentType = response.headers['content-type'] || ''
    if (contentType.includes('application/json')) {
      const text = await response.data.text()
      const result = JSON.parse(text)
      ElMessage.error(result.msg || 'Xuất Excel thất bại')
      return false
    }

    downloadBlob(response.data, fileName)
    ElMessage.success('Xuất Excel thành công')
    return true
  } catch (error) {
    if (error?.response?.data instanceof Blob) {
      try {
        const text = await error.response.data.text()
        const result = JSON.parse(text)
        ElMessage.error(result.msg || 'Xuất Excel thất bại')
        return false
      } catch {
        // fall through
      }
    }

    ElMessage.error(error?.response?.data?.msg || error?.message || 'Xuất Excel thất bại')
    return false
  }
}
