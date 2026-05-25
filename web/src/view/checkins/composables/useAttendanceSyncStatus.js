import { ref, computed } from 'vue'
import { getSyncStatus } from '@/api/checkins/condition'

export function useAttendanceSyncStatus(acId) {
  const syncData = ref({
    syncState: 'no_rules',
    needsSync: false,
    conditionCount: 0,
    agpConditionCount: 0,
    expectedAgpConditionCount: 0,
    agpCount: 0,
    participantCount: 0,
    agpWithoutMappingCount: 0,
  })

  const tagType = computed(() => {
    switch (syncData.value.syncState) {
      case 'synced': return 'success'
      case 'stale': return 'warning'
      case 'never_synced': return 'warning'
      default: return 'info'
    }
  })

  const tagLabel = computed(() => {
    switch (syncData.value.syncState) {
      case 'synced': return 'Đã đồng bộ điều kiện'
      case 'stale': return 'Cấu hình đã thay đổi'
      case 'never_synced': return 'Chưa đồng bộ điều kiện'
      default: return 'Chưa có điều kiện'
    }
  })

  const tooltipContent = computed(() => {
    const d = syncData.value
    return `~${d.expectedAgpConditionCount} mapping · ${d.agpCount} AGP · ${d.participantCount} thành viên · ${d.conditionCount} điều kiện`
  })

  const fetchSyncStatus = async () => {
    const res = await getSyncStatus({ attendanceId: acId.value ?? acId })
    if (res.code === 0) {
      syncData.value = res.data
    }
  }

  return { syncData, tagType, tagLabel, tooltipContent, fetchSyncStatus }
}
