<template>
  <div class="lucky-number-feature">
    <!-- Simple notice when user is working toward a lucky number -->
    <div v-if="showProgress && !luckyNumber && remainingCheckins > 0" 
         class="mt-4 p-6 bg-gradient-to-r from-gray-50 to-gray-100 rounded-lg border border-gray-200 shadow-md">
      <h3 class="text-lg font-medium text-gray-700 mb-3 flex items-center">
        <i class="el-icon-star-on mr-2 text-yellow-500 lucky-star-pulse"></i>
        <span>Điểm danh để nhận số may mắn</span>
      </h3>
      
      <p class="mt-2 text-sm text-gray-600 italic">
        <i class="el-icon-info-filled mr-1 text-blue-500"></i>
        <strong>Lưu ý:</strong> Bạn phải điểm danh thành công ít nhất một điều kiện mới nhận được số may mắn
      </p>
    </div>
    
    <!-- Lucky Number Display - Show when user has a lucky number -->
    <div v-if="luckyNumber" 
         class="mt-4 lucky-number-container"
         :class="{'lucky-number-new': isNewLuckyNumber}">
         
      <!-- Confetti effect when receiving a new lucky number -->
      <div v-if="showConfetti" class="confetti-container">
        <div v-for="n in 50" :key="n" class="confetti" 
            :style="{
              left: Math.random() * 100 + '%',
              top: -10 * Math.random() + '%',
              backgroundColor: getRandomColor(),
              width: 5 + Math.random() * 10 + 'px',
              height: 5 + Math.random() * 10 + 'px',
              animationDelay: Math.random() * 3 + 's'
            }">
        </div>
      </div>
      
      <!-- Lucky number title with animation -->
      <h3 class="text-xl font-bold mb-4 flex items-center justify-center" 
          :class="{'lucky-number-shine': showAnimations}">
        <i class="el-icon-star-on mr-2 text-yellow-500 lucky-star"></i>
        Số may mắn của bạn
        <i class="el-icon-star-on ml-2 text-yellow-500 lucky-star"></i>
      </h3>
      
      <!-- Lucky number badge with decorative stars -->
      <div class="lucky-number-badge" :class="{'lucky-number-celebrate': showAnimations}">
        <span class="lucky-star-decoration star-top-left">★</span>
        <span class="lucky-star-decoration star-top-right">★</span>
        <span class="lucky-number-value">{{ luckyNumber }}</span>
        <span class="lucky-star-decoration star-bottom-left">★</span>
        <span class="lucky-star-decoration star-bottom-right">★</span>
      </div>
      
      <p class="mt-4 text-sm text-yellow-700 italic">
        Hãy lưu lại số may mắn này để tham gia các hoạt động khác!
      </p>
      
      <!-- Optional share buttons could be added here -->
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue';
import '../style/lucky-animation.css';

const props = defineProps({
  // The lucky number if available
  luckyNumber: {
    type: Number,
    default: null
  },
  // Current number of check-ins
  currentCheckins: {
    type: Number,
    default: 0
  },
  // Number of check-ins required to get a lucky number
  requiredCheckins: {
    type: Number,
    default: 0
  },
  // Whether the lucky number was just received (for animation)
  isNewLuckyNumber: {
    type: Boolean,
    default: false
  },
  // Whether to show the progress bar
  showProgress: {
    type: Boolean,
    default: true
  }
});

// Reactive references
const showAnimations = ref(false);
const showConfetti = ref(false);

// Computed properties
const remainingCheckins = computed(() => {
  return Math.max(0, props.requiredCheckins - props.currentCheckins);
});

const progressPercent = computed(() => {
  if (props.requiredCheckins <= 0) return 0;
  return Math.min(100, (props.currentCheckins / props.requiredCheckins) * 100);
});

// Function to generate random colors for confetti
const getRandomColor = () => {
  const colors = [
    '#ffd700', // Gold
    '#ffcc33', // Golden yellow
    '#ffc107', // Amber
    '#ffb300', // Amber darken-1
    '#ffa000', // Amber darken-2
    '#ff8f00', // Amber darken-3
    '#ff6f00'  // Amber darken-4
  ];
  return colors[Math.floor(Math.random() * colors.length)];
};

// Watch for new lucky number to trigger animations
watch(() => props.isNewLuckyNumber, (newValue) => {
  if (newValue && props.luckyNumber) {
    startAnimations();
  }
});

// Watch for lucky number to trigger animations when component mounts with existing lucky number
watch(() => props.luckyNumber, (newValue) => {
  if (newValue && props.isNewLuckyNumber) {
    startAnimations();
  }
});

// Start animations when component is mounted if a lucky number exists and is new
onMounted(() => {
  if (props.luckyNumber && props.isNewLuckyNumber) {
    startAnimations();
  }
});

// Function to start animations
const startAnimations = () => {
  showAnimations.value = true;
  showConfetti.value = true;
  
  // Stop animations after a while
  setTimeout(() => { 
    showAnimations.value = false;
  }, 5000);
  
  setTimeout(() => { 
    showConfetti.value = false;
  }, 6000);
};
</script>

<style scoped>
/* Base styling - detailed animations are in lucky-animation.css */
.lucky-number-feature {
  margin-top: 1.5rem;
}

.lucky-number-new {
  animation: fadeIn 1s ease-in-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
