<!-- 
  This component displays a detailed view of the check-in counter
  with animated indicators for multiple check-ins.
-->
<template>
  <div class="checkin-counter">
    <div v-if="counter > 0" class="counter-display">
      <el-badge :value="counter" :type="counter > 5 ? 'warning' : 'success'" class="item">
        <span class="counter-label">
          {{ label || 'Số lần điểm danh' }}
        </span>
      </el-badge>
      
      <div v-if="counter > 1" class="counter-progress">
        <div class="progress-bar">
          <div 
            class="progress-fill" 
            :style="{ width: `${Math.min(100, (counter / maxDisplay) * 100)}%` }"
          ></div>
        </div>
        <div class="counter-details">
          <span v-if="lastCheckin" class="detail-text">
            Gần đây nhất: {{ formatDate(lastCheckin) }}
          </span>
        </div>
      </div>
    </div>
    <div v-else class="no-checkins">
      Chưa có điểm danh
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { formatDateTime } from '@/utils/format';

const props = defineProps({
  counter: {
    type: Number,
    default: 0
  },
  label: {
    type: String,
    default: ''
  },
  maxDisplay: {
    type: Number,
    default: 10
  },
  lastCheckin: {
    type: [Date, String],
    default: null
  }
});

const formatDate = (date) => {
  if (!date) return '';
  return formatDateTime(date);
};
</script>

<style scoped>
.checkin-counter {
  margin: 8px 0;
  padding: 6px;
  border-radius: 6px;
}

.counter-display {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.counter-label {
  font-size: 0.9rem;
  font-weight: 500;
}

.counter-progress {
  margin-top: 4px;
}

.progress-bar {
  height: 6px;
  background-color: #f0f0f0;
  border-radius: 3px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background-color: #67c23a;
  border-radius: 3px;
  transition: width 0.5s ease-out;
}

.counter-details {
  margin-top: 4px;
  font-size: 0.75rem;
  color: #909399;
}

.no-checkins {
  color: #909399;
  font-size: 0.85rem;
  font-style: italic;
}
</style>
