<template>
  <div class="flex flex-col gap-6">
    <!-- Session progress -->
    <div>
      <div class="font-semibold text-gray-700 mb-2">Tiến độ thời gian phiên</div>
      <div class="text-xs text-gray-400 mb-1 flex justify-between">
        <span>{{ session.startDate ? fmtDate(session.startDate) : '—' }}</span>
        <span>{{ session.endDate ? fmtDate(session.endDate) : '—' }}</span>
      </div>
      <el-progress
        :percentage="Math.round(session.elapsedPercent ?? 0)"
        :status="session.elapsedPercent >= 100 ? 'success' : undefined"
        :stroke-width="14"
      />
      <div class="text-xs text-gray-500 mt-1">
        {{ Math.round(session.elapsedPercent ?? 0) }}% thời gian đã trôi qua
      </div>
    </div>

    <!-- Checkins by hour histogram -->
    <div>
      <div class="font-semibold text-gray-700 mb-2">Lượt điểm danh theo giờ</div>
      <template v-if="checkinsByHour.length">
        <vue-echarts :option="hourOption" style="height: 200px; width: 100%;" autoresize />
      </template>
      <div v-else class="text-gray-400 text-sm py-4 text-center">Chưa có lượt điểm danh</div>
    </div>

    <!-- Conditions timeline -->
    <div>
      <div class="font-semibold text-gray-700 mb-2">Timeline điều kiện</div>
      <template v-if="timedConditions.length">
        <vue-echarts :option="timelineOption" :style="{ height: timelineHeight, width: '100%' }" autoresize />
      </template>
      <div v-else class="text-gray-400 text-sm py-4 text-center">Điều kiện không có khung giờ</div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { VueEcharts } from 'vue3-echarts'
import moment from 'moment'

const props = defineProps({
  session:             { type: Object, default: () => ({}) },
  checkinsByHour:      { type: Array,  default: () => [] },
  conditionsTimeline:  { type: Array,  default: () => [] },
})

const fmtDate = (d) => moment(d).format('DD/MM/YYYY HH:mm')

const hourOption = computed(() => ({
  tooltip: { trigger: 'axis' },
  grid: { left: '5%', right: '2%', top: 10, bottom: 30 },
  xAxis: {
    type: 'category',
    data: props.checkinsByHour.map(b => b.label),
    axisLabel: { fontSize: 11 },
  },
  yAxis: { type: 'value', minInterval: 1 },
  series: [{
    type: 'bar',
    data: props.checkinsByHour.map(b => b.count),
    itemStyle: { color: '#3b82f6', borderRadius: [3, 3, 0, 0] },
    label: { show: true, position: 'top', fontSize: 10 },
  }],
}))

// Only show conditions that have both startAt and endAt
const timedConditions = computed(() =>
  props.conditionsTimeline.filter(c => c.startAt && c.endAt)
)

const timelineHeight = computed(() => `${Math.max(100, timedConditions.value.length * 40 + 60)}px`)

const timelineOption = computed(() => {
  const sessionStart = props.session.startDate ? new Date(props.session.startDate).getTime() : null
  const sessionEnd   = props.session.endDate   ? new Date(props.session.endDate).getTime()   : null
  const now          = Date.now()

  const items = timedConditions.value.map((c, i) => ({
    name: [c.groupName || 'Tất cả', c.areaName].filter(Boolean).join(' · ') || `ĐK #${c.id}`,
    value: [
      i,
      new Date(c.startAt).getTime(),
      new Date(c.endAt).getTime(),
    ],
  }))

  const markLine = (sessionStart && sessionEnd && now >= sessionStart && now <= sessionEnd)
    ? { data: [{ xAxis: now, label: { formatter: 'Bây giờ', position: 'end' }, lineStyle: { color: '#f59e0b', type: 'dashed' } }] }
    : undefined

  return {
    tooltip: {
      formatter: (p) => {
        const [, start, end] = p.value
        return `${p.name}<br/>${moment(start).format('HH:mm DD/MM')} → ${moment(end).format('HH:mm DD/MM')}`
      },
    },
    grid: { left: '2%', right: '3%', top: 10, bottom: 20, containLabel: true },
    xAxis: {
      type: 'time',
      min: sessionStart ?? undefined,
      max: sessionEnd   ?? undefined,
      axisLabel: { formatter: '{MM}/{dd} {HH}:{mm}', fontSize: 10 },
    },
    yAxis: {
      type: 'category',
      data: items.map(i => i.name),
      axisLabel: { fontSize: 11 },
    },
    series: [{
      type: 'custom',
      renderItem: (_params, api) => {
        const start = api.value(1)
        const end   = api.value(2)
        const coord0 = api.coord([start, api.value(0)])
        const coord1 = api.coord([end, api.value(0)])
        const height = api.size([0, 1])[1] * 0.6
        return {
          type: 'rect',
          shape: { x: coord0[0], y: coord0[1] - height / 2, width: Math.max(2, coord1[0] - coord0[0]), height },
          style: api.style(),
        }
      },
      itemStyle: { color: '#6366f1', opacity: 0.8, borderRadius: 4 },
      encode: { x: [1, 2], y: 0 },
      data: items.map(i => i.value),
      markLine,
    }],
  }
})
</script>
