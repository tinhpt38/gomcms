<template>
  <div>
    <div class="flex items-center gap-2 mb-2">
      <span class="font-semibold text-gray-700">Bản đồ vị trí điểm danh</span>
      <el-tag v-if="truncated" type="warning" size="small">Hiển thị tối đa 2000 điểm</el-tag>
      <el-tag v-if="!geofences.length && !checkinPoints.length && !failPoints.length" type="info" size="small">
        Chưa có dữ liệu vị trí
      </el-tag>
    </div>

    <vue-echarts :option="chartOption" style="height: 420px; width: 100%;" autoresize />

    <div class="flex gap-4 mt-2 text-sm text-gray-500">
      <span class="flex items-center gap-1"><span class="inline-block w-3 h-3 rounded-full bg-green-500"></span>Điểm danh OK</span>
      <span class="flex items-center gap-1"><span class="inline-block w-3 h-3 rounded-full bg-red-500"></span>Thất bại</span>
      <span class="flex items-center gap-1"><span class="inline-block w-3 h-3 rounded-full bg-blue-300 opacity-50 border border-blue-400"></span>Khu vực</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { VueEcharts } from 'vue3-echarts'

const props = defineProps({
  geofences:     { type: Array, default: () => [] },
  checkinPoints: { type: Array, default: () => [] },
  failPoints:    { type: Array, default: () => [] },
  truncated:     { type: Boolean, default: false },
})

// Convert radius in metres to approximate degrees for the ECharts coordinate system.
// 1 degree latitude ≈ 111 320 m; longitude varies with lat but we approximate equally.
const mToDeg = (m) => m / 111320

const chartOption = computed(() => {
  // Collect all points to auto-scale axes
  const allLngs = [
    ...props.geofences.map(g => g.lng),
    ...props.checkinPoints.map(p => p.lng),
    ...props.failPoints.map(p => p.lng),
  ].filter(Boolean)
  const allLats = [
    ...props.geofences.map(g => g.lat),
    ...props.checkinPoints.map(p => p.lat),
    ...props.failPoints.map(p => p.lat),
  ].filter(Boolean)

  if (!allLngs.length) {
    return {
      graphic: [{ type: 'text', left: 'center', top: 'center', style: { text: 'Chưa có dữ liệu tọa độ', fontSize: 14, fill: '#aaa' } }]
    }
  }

  const pad = 0.001
  const xMin = Math.min(...allLngs) - pad
  const xMax = Math.max(...allLngs) + pad
  const yMin = Math.min(...allLats) - pad
  const yMax = Math.max(...allLats) + pad

  // Build geofence circles as ECharts graphic elements
  const circles = props.geofences.map(g => {
    const rDeg = mToDeg(g.radiusM || 50)
    // Map lat/lng to pixel via percentage — approximate linear mapping
    const cx = `${((g.lng - xMin) / (xMax - xMin)) * 100}%`
    const cy = `${(1 - (g.lat - yMin) / (yMax - yMin)) * 100}%`
    // Estimate pixel radius from chart width (~500px)
    const pxR = (rDeg / (xMax - xMin)) * 500
    return {
      type: 'circle',
      shape: { cx: 0, cy: 0, r: Math.max(6, Math.min(pxR, 200)) },
      position: [cx, cy],
      style: { fill: 'rgba(59,130,246,0.12)', stroke: '#3b82f6', lineWidth: 2 },
      silent: true,
      z: 0,
    }
  })

  return {
    tooltip: {
      trigger: 'item',
      formatter: (params) => {
        const d = params.data
        if (!d) return ''
        return `${d[2] ?? ''}<br/>${d[3] ?? ''}<br/>精度 ${d[4] ?? '-'}m`
      },
    },
    xAxis: {
      min: xMin, max: xMax,
      type: 'value',
      axisLabel: { formatter: (v) => v.toFixed(4) },
      name: 'Kinh độ',
    },
    yAxis: {
      min: yMin, max: yMax,
      type: 'value',
      axisLabel: { formatter: (v) => v.toFixed(4) },
      name: 'Vĩ độ',
    },
    graphic: circles,
    series: [
      {
        name: 'Điểm danh OK',
        type: 'scatter',
        symbolSize: 8,
        itemStyle: { color: '#22c55e', opacity: 0.8 },
        data: props.checkinPoints.map(p => [p.lng, p.lat, p.label, p.at, p.accuracy]),
      },
      {
        name: 'Thất bại',
        type: 'scatter',
        symbolSize: 8,
        itemStyle: { color: '#ef4444', opacity: 0.8 },
        data: props.failPoints.map(p => [p.lng, p.lat, p.label, p.at]),
      },
    ],
  }
})
</script>
