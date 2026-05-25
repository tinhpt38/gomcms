import { ref } from 'vue'
import { findAttendanceArea } from '@/api/checkins/attendance'
import { getGroupList } from '@/api/checkins/group'

export function useAttendanceOptions(acId) {
  const groupOptions = ref([])
  const areaOptions = ref([])

  const refreshGroups = async () => {
    const res = await getGroupList({ page: 1, pageSize: -1, attendanceId: acId })
    if (res.code === 0) {
      groupOptions.value = res.data.list
    }
  }

  const refreshAreas = async () => {
    const res = await findAttendanceArea({ id: acId })
    if (res.code === 0) {
      areaOptions.value = res.data.map(item => ({
        ID: item.ID,
        name: item.area?.name,
      }))
    }
  }

  const refreshAll = async () => {
    await Promise.all([refreshGroups(), refreshAreas()])
  }

  return { groupOptions, areaOptions, refreshGroups, refreshAreas, refreshAll }
}
