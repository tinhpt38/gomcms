<template>
  <div class="p-2">
    <!-- Loading state -->
    <div v-if="loading" class="flex justify-center items-center py-16 text-gray-400">
      <el-icon class="animate-spin mr-2"><Loading /></el-icon>
      Đang tải dữ liệu tổng quan…
    </div>

    <!-- Error state -->
    <el-alert v-else-if="error" type="error" :title="error" :closable="false" show-icon class="mb-4" />

    <template v-else-if="overview">
      <!-- Two-column grid: Map (left) + Time (right) -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
        <div class="bg-white rounded shadow-sm p-4">
          <OverviewMapChart
            :geofences="overview.geofences ?? []"
            :checkin-points="overview.checkinPoints ?? []"
            :fail-points="overview.failPoints ?? []"
            :truncated="overview.mapTruncated"
          />
        </div>
        <div class="bg-white rounded shadow-sm p-4">
          <OverviewTimeSection
            :session="overview.session"
            :checkins-by-hour="overview.checkinsByHour ?? []"
            :conditions-timeline="overview.conditionsTimeline ?? []"
          />
        </div>
      </div>

      <!-- Two-column grid: Progress (left) + Quality (right) -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div class="bg-white rounded shadow-sm p-4">
          <OverviewProgressSection
            :session="overview.session"
            :group-stats="overview.groupStats ?? []"
            :coverage="overview.coverage ?? {}"
          />
        </div>
        <div class="bg-white rounded shadow-sm p-4">
          <OverviewQualitySection
            :quality-points="overview.qualityPoints ?? []"
            :fail-reasons="overview.failReasons ?? []"
            :top-i-ps="overview.topIPs ?? []"
            :shared-visitors="overview.sharedVisitors ?? []"
          />
        </div>
      </div>
    </template>

    <!-- Empty state (mounted but no data yet — should not normally show) -->
    <div v-else class="text-center text-gray-400 py-16">Chưa có dữ liệu</div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useAttendanceOverview } from '@/view/checkins/composables/useAttendanceOverview'
import OverviewMapChart       from './OverviewMapChart.vue'
import OverviewTimeSection    from './OverviewTimeSection.vue'
import OverviewProgressSection from './OverviewProgressSection.vue'
import OverviewQualitySection from './OverviewQualitySection.vue'

const props = defineProps({
  acId: { type: Number, required: true },
})

const { overview, loading, error, fetchOverview } = useAttendanceOverview(props.acId)

onMounted(() => {
  fetchOverview()
})
</script>
