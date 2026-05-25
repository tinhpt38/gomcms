<template>
  <div class="flex flex-col gap-6">
    <!-- Donut: checked-in vs not -->
    <div>
      <div class="font-semibold text-gray-700 mb-2">Tiến độ điểm danh</div>
      <vue-echarts :option="donutOption" style="height: 220px; width: 100%;" autoresize />
    </div>

    <!-- Bar: per group -->
    <div v-if="groupStats.length">
      <div class="font-semibold text-gray-700 mb-2">Điểm danh theo nhóm</div>
      <vue-echarts :option="groupBarOption" :style="{ height: groupChartHeight, width: '100%' }" autoresize />
    </div>

    <!-- Coverage cards -->
    <div>
      <div class="font-semibold text-gray-700 mb-2">Phủ sóng điều kiện (AGP coverage)</div>
      <div class="flex gap-3 flex-wrap">
        <div class="bg-green-50 border border-green-200 rounded px-4 py-3 flex flex-col items-center min-w-[100px]">
          <span class="text-2xl font-bold text-green-600">{{ coverage.fullyMapped }}</span>
          <span class="text-xs text-gray-500 mt-1">Đủ điều kiện</span>
        </div>
        <div class="bg-orange-50 border border-orange-200 rounded px-4 py-3 flex flex-col items-center min-w-[100px]">
          <span class="text-2xl font-bold text-orange-500">{{ coverage.partial }}</span>
          <span class="text-xs text-gray-500 mt-1">Thiếu 1 phần</span>
        </div>
        <div class="bg-gray-50 border border-gray-200 rounded px-4 py-3 flex flex-col items-center min-w-[100px]">
          <span class="text-2xl font-bold text-gray-500">{{ coverage.noMapping }}</span>
          <span class="text-xs text-gray-500 mt-1">Không có mapping</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { VueEcharts } from 'vue3-echarts'

const props = defineProps({
  session:    { type: Object, default: () => ({}) },
  groupStats: { type: Array,  default: () => [] },
  coverage:   { type: Object, default: () => ({ fullyMapped: 0, partial: 0, noMapping: 0 }) },
})

const donutOption = computed(() => {
  const total   = props.session.total ?? 0
  const checked = props.session.totalCheckin ?? 0
  const notYet  = Math.max(0, total - checked)
  return {
    tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
    legend: { bottom: 0 },
    series: [{
      type: 'pie',
      radius: ['45%', '70%'],
      label: { show: false },
      emphasis: { label: { show: true, fontSize: 14, fontWeight: 'bold' } },
      data: [
        { name: 'Đã điểm danh', value: checked, itemStyle: { color: '#22c55e' } },
        { name: 'Chưa điểm danh', value: notYet, itemStyle: { color: '#e5e7eb' } },
      ],
    }],
  }
})

const groupChartHeight = computed(() => `${Math.max(120, props.groupStats.length * 32 + 60)}px`)

const groupBarOption = computed(() => ({
  tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
  legend: { bottom: 0 },
  grid: { left: '2%', right: '3%', top: 10, bottom: 40, containLabel: true },
  xAxis: { type: 'value', minInterval: 1 },
  yAxis: {
    type: 'category',
    data: props.groupStats.map(g => g.name),
    axisLabel: { fontSize: 11, width: 100, overflow: 'truncate' },
  },
  series: [
    {
      name: 'Thành viên',
      type: 'bar',
      data: props.groupStats.map(g => g.memberCount),
      itemStyle: { color: '#93c5fd', borderRadius: [0, 3, 3, 0] },
    },
    {
      name: 'Đã điểm danh',
      type: 'bar',
      data: props.groupStats.map(g => g.checkinDistinct),
      itemStyle: { color: '#22c55e', borderRadius: [0, 3, 3, 0] },
    },
  ],
}))
</script>
