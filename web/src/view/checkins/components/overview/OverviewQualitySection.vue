<template>
  <div class="flex flex-col gap-6">
    <!-- Shared visitor alert -->
    <el-alert
      v-if="sharedVisitors.length"
      type="warning"
      :closable="false"
      show-icon
      class="mb-2"
    >
      <template #title>
        Phát hiện {{ sharedVisitors.length }} thiết bị điểm danh nhiều tài khoản —
        có thể gian lận. <span class="font-semibold">Xem bảng bên dưới.</span>
      </template>
    </el-alert>

    <!-- GPS quality scatter -->
    <div>
      <div class="font-semibold text-gray-700 mb-2">Chất lượng GPS (khoảng cách vs độ chính xác)</div>
      <template v-if="qualityPoints.length">
        <vue-echarts :option="scatterOption" style="height: 260px; width: 100%;" autoresize />
        <div class="flex gap-4 mt-1 text-xs text-gray-500">
          <span><span class="inline-block w-2 h-2 rounded-full bg-green-500 mr-1"></span>Trong khu vực</span>
          <span><span class="inline-block w-2 h-2 rounded-full bg-red-400 mr-1"></span>Ngoài khu vực</span>
          <span><span class="inline-block w-2 h-2 rounded-full bg-gray-400 mr-1"></span>Không xác định</span>
        </div>
      </template>
      <div v-else class="text-gray-400 text-sm py-4 text-center">Không có điểm GPS trong phiên này</div>
    </div>

    <!-- Fail reasons pie -->
    <div v-if="failReasons.length">
      <div class="font-semibold text-gray-700 mb-2">Lý do điểm danh thất bại</div>
      <vue-echarts :option="pieOption" style="height: 240px; width: 100%;" autoresize />
    </div>

    <!-- Top IPs -->
    <div v-if="topIPs.length">
      <div class="font-semibold text-gray-700 mb-2">Top IP điểm danh</div>
      <el-table :data="topIPs" size="small" border style="width: 100%">
        <el-table-column prop="ip" label="IP" />
        <el-table-column prop="count" label="Lượt" width="80" align="right" />
      </el-table>
    </div>

    <!-- Shared visitors -->
    <div v-if="sharedVisitors.length">
      <div class="font-semibold text-gray-700 mb-2">Thiết bị dùng chung (cảnh báo gian lận)</div>
      <el-table :data="sharedVisitors" size="small" border style="width: 100%">
        <el-table-column prop="visitorId" label="Visitor ID">
          <template #default="scope">
            <span class="font-mono text-xs">{{ scope.row.visitorId?.substring(0, 20) }}...</span>
          </template>
        </el-table-column>
        <el-table-column prop="attemptCount" label="Tài khoản" width="90" align="right" />
        <el-table-column label="Email">
          <template #default="scope">
            {{ scope.row.emails?.join(', ') }}
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { VueEcharts } from 'vue3-echarts'

const props = defineProps({
  qualityPoints:  { type: Array, default: () => [] },
  failReasons:    { type: Array, default: () => [] },
  topIPs:         { type: Array, default: () => [] },
  sharedVisitors: { type: Array, default: () => [] },
})

const scatterOption = computed(() => ({
  tooltip: {
    trigger: 'item',
    formatter: (p) => {
      const [dist, acc, areaName] = p.data
      return `Khoảng cách: ${dist}m<br/>Độ lệch GPS: ${acc}m<br/>Khu vực: ${areaName || '—'}`
    },
  },
  grid: { left: '5%', right: '2%', top: 20, bottom: 40, containLabel: true },
  xAxis: { type: 'value', name: 'Khoảng cách (m)', nameLocation: 'end' },
  yAxis: { type: 'value', name: 'Độ lệch GPS (m)', nameLocation: 'end' },
  series: [
    {
      name: 'Trong khu vực',
      type: 'scatter',
      symbolSize: 7,
      itemStyle: { color: '#22c55e', opacity: 0.7 },
      data: props.qualityPoints
        .filter(p => !p.outOfRadius && p.areaName)
        .map(p => [p.distanceM, p.accuracy, p.areaName]),
    },
    {
      name: 'Ngoài khu vực',
      type: 'scatter',
      symbolSize: 7,
      itemStyle: { color: '#ef4444', opacity: 0.7 },
      data: props.qualityPoints
        .filter(p => p.outOfRadius)
        .map(p => [p.distanceM, p.accuracy, p.areaName]),
    },
    {
      name: 'Không xác định',
      type: 'scatter',
      symbolSize: 6,
      itemStyle: { color: '#9ca3af', opacity: 0.6 },
      data: props.qualityPoints
        .filter(p => !p.areaName)
        .map(p => [p.distanceM, p.accuracy, '']),
    },
  ],
}))

const pieOption = computed(() => ({
  tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
  legend: { bottom: 0, type: 'scroll' },
  series: [{
    type: 'pie',
    radius: ['35%', '65%'],
    label: { show: false },
    emphasis: { label: { show: true, fontSize: 12, fontWeight: 'bold' } },
    data: props.failReasons.map(r => ({ name: r.reason, value: r.count })),
  }],
}))
</script>
