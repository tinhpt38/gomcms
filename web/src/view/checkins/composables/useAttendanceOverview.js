import { ref } from 'vue'
import { getAttendanceOverview } from '@/api/checkins/attendance'

export function useAttendanceOverview(acId) {
  const overview = ref(null)
  const loading = ref(false)
  const error = ref(null)

  const fetchOverview = async () => {
    loading.value = true
    error.value = null
    try {
      const res = await getAttendanceOverview({ attendanceId: acId.value ?? acId })
      if (res.code === 0) {
        overview.value = res.data
      } else {
        error.value = res.msg ?? 'Không lấy được dữ liệu'
      }
    } catch (e) {
      error.value = e.message
    } finally {
      loading.value = false
    }
  }

  return { overview, loading, error, fetchOverview }
}
