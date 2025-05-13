<template>
    <div class="random-wheel-container my-8">
        <el-button type="success" @click="onLuckyClick" class="mt-2 mb-4" :disabled="isSpinning">
            <span v-if="!isSpinning">Tìm số may mắn</span>
            <span v-else>Đang quay...</span>
        </el-button>
        
        <div class="wheel-container mb-8" v-if="isSpinning || hasResult">
            <div class="wheel" ref="wheelRef" :class="{ 'spinning': isSpinning }">
                <div class="wheel-inner">
                    <div v-for="(item, index) in wheelItems" :key="index" class="wheel-item" 
                        :style="getWheelItemStyle(index)">
                        {{ item }}
                    </div>
                    <div class="wheel-center"></div>
                </div>
                <div class="wheel-pointer"></div>
            </div>
        </div>
        
        <div v-if="hasResult" class="result-container">
            <div class="winner-card" :class="{ 'animate__animated animate__tada': showAnimation }">
                <div class="confetti-container" v-if="showConfetti">
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
                
                <div class="winner-info">
                    <div class="lucky-number-display">
                        <span class="lucky-number-label">Số May Mắn</span>
                        <span class="lucky-number">{{ luckyNumber }}</span>
                    </div>
                    <h2 class="winner-title">Người may mắn</h2>
                    <template v-if="luckyMember && luckyMember.email">
                        <p class="winner-email">{{ luckyMember.email }}</p>
                        <p class="winner-name">{{ luckyMember.fullName?.replaceAll("undefined", "") ?? "--" }}</p>
                    </template>
                    <template v-else>
                        <p class="winner-email no-winner">Chưa tìm được chủ nhân</p>
                        <p class="winner-name">--</p>
                    </template>
                </div>
            </div>
            
            <!-- History of lucky numbers -->
            <div class="history-container" v-if="luckyNumbersHistory.length > 0">
                <div class="history-header">
                    <h3 class="history-title">Lịch sử số may mắn</h3>
                    <el-button type="danger" size="small" @click="clearHistory" icon="Delete">Xóa lịch sử</el-button>
                </div>
                <div class="history-list">
                    <div v-for="(item, idx) in luckyNumbersHistory" :key="idx" class="history-item">
                        <div class="history-number">{{ item.luckyNumber }}</div>
                        <div class="history-info" v-if="item.email">
                            <span class="history-email">{{ item.email }}</span>
                            <span class="history-name">{{ item.fullName || "--" }}</span>
                        </div>
                        <div class="history-info" v-else>
                            <span class="history-email no-winner">Chưa tìm được chủ nhân</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
    findLuckyParticipant,
    clearLuckyHistory
} from '@/api/checkins/participant'
import 'animate.css'

const props = defineProps({
    acId: {
        type: Number,
        required: true
    },
})

const searchInfo = ref({
    attendanceId: props.acId,
})

const luckyMember = ref({
    email: null,
    fullName: null,
})

// Lucky number storage
const luckyNumber = ref(null)
// History of selected lucky numbers
const luckyNumbersHistory = ref([])

// Wheel related states
const isSpinning = ref(false)
const hasResult = ref(false)
const showAnimation = ref(false)
const showConfetti = ref(false)
const wheelRef = ref(null)
const wheelItems = ref([
    '1', '2', '3', '4', '5', '6', '7', '8',
    '9', '10', '11', '12', '13', '14', '15', '16'
])

// Compute if we have a result to show
const hasLuckyMember = computed(() => {
    return luckyMember.value && luckyMember.value.email !== null;
})

// Function to get random colors for confetti
const getRandomColor = () => {
    const colors = [
        '#ffd700', // Gold
        '#ffcc33', // Golden yellow
        '#ffc107', // Amber
        '#ff8f00', // Amber darken-3
        '#ff6f00'  // Amber darken-4
    ];
    return colors[Math.floor(Math.random() * colors.length)];
};

// Function to get wheel item style based on index
const getWheelItemStyle = (index) => {
    const angle = (360 / wheelItems.value.length) * index;
    return {
        transform: `rotate(${angle}deg) translateY(-125px) rotate(-${angle}deg)`,
        backgroundColor: index % 2 === 0 ? '#ff9800' : '#ffeb3b',
    }
}

// Start animations when result is shown
const startAnimations = () => {
    showAnimation.value = true;
    showConfetti.value = true;
    
    // Stop animations after a while
    setTimeout(() => { 
        showAnimation.value = false;
    }, 5000);
    
    setTimeout(() => { 
        showConfetti.value = false;
    }, 6000);
};

// Clear history of lucky numbers
const clearHistory = () => {
    ElMessageBox.confirm(
        'Bạn có chắc chắn muốn xóa lịch sử số may mắn?',
        'Xác nhận xóa',
        {
            confirmButtonText: 'Xóa',
            cancelButtonText: 'Hủy',
            type: 'warning',
        }
    ).then(() => {
        clearLuckyHistory({ attendanceId: props.acId })
            .then(() => {
                ElMessage.success('Đã xóa lịch sử số may mắn');
                luckyNumbersHistory.value = [];
            })
            .catch(err => {
                ElMessage.error('Có lỗi xảy ra khi xóa lịch sử');
                console.error(err);
            });
    }).catch(() => {
        // User canceled
    });
};

const onLuckyClick = async () => {
    try {
        isSpinning.value = true;
        hasResult.value = false;
        
        // Randomize wheel rotation
        if (wheelRef.value) {
            const turns = 5 + Math.floor(Math.random() * 5); // 5-10 turns
            const extraDegrees = Math.floor(Math.random() * 360);
            const totalRotation = turns * 360 + extraDegrees;
            
            wheelRef.value.style.transition = 'transform 4s cubic-bezier(0.1, 0.7, 0.1, 1)';
            wheelRef.value.style.transform = `rotate(${totalRotation}deg)`;
        }
        
        // Fetch data after a delay to simulate wheel spinning
        searchInfo.value.attendanceId = props.acId;
        
        // Add a delay to match the wheel animation
        setTimeout(async () => {
            const res = await findLuckyParticipant(searchInfo.value);
            luckyMember.value = res.data.participant;
            luckyNumber.value = res.data.luckyNumber;
            
            // Load history directly from API response
            if (res.data.luckyHistory && res.data.luckyHistory.length > 0) {
                luckyNumbersHistory.value = res.data.luckyHistory;
            }
            
            // Show results after spinning finishes
            isSpinning.value = false;
            hasResult.value = true;
            startAnimations();
        }, 4500); // Slightly less than the wheel animation time
        
    } catch (error) {
        isSpinning.value = false;
        ElMessage.error('Có lỗi xảy ra');
    }
}

onMounted(() => {
    // Check if we already have results to display
    if (hasLuckyMember.value) {
        hasResult.value = true;
    }
    
    // Tải danh sách lịch sử từ server khi component được mount
    findLuckyParticipant({ attendanceId: props.acId }).then(res => {
        if (res.data && res.data.luckyHistory) {
            luckyNumbersHistory.value = res.data.luckyHistory;
            
            // Nếu có số may mắn mới nhất, hiển thị nó
            if (res.data.luckyNumber) {
                luckyNumber.value = res.data.luckyNumber;
                luckyMember.value = res.data.participant || {};
                hasResult.value = true;
            }
        }
    }).catch(err => {
        console.error("Không thể tải lịch sử số may mắn", err);
    });
})
</script>

<style scoped>
.random-wheel-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 100%;
}

/* Wheel Styles */
.wheel-container {
    position: relative;
    width: 300px;
    height: 300px;
    margin: 0 auto;
}

.wheel {
    position: relative;
    width: 100%;
    height: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
}

.wheel-inner {
    position: relative;
    width: 100%;
    height: 100%;
    border-radius: 50%;
    background: #f5f5f5;
    overflow: hidden;
    box-shadow: 0 0 20px rgba(0, 0, 0, 0.3);
    border: 5px solid #e0e0e0;
}

.wheel-item {
    position: absolute;
    top: 50%;
    left: 50%;
    width: 60px;
    height: 60px;
    text-align: center;
    display: flex;
    justify-content: center;
    align-items: center;
    font-weight: bold;
    font-size: 20px;
    border-radius: 50%;
    transform-origin: center;
    color: #333;
}

.wheel-center {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 50px;
    height: 50px;
    background: #fff;
    border-radius: 50%;
    border: 5px solid #ccc;
    z-index: 2;
}

.wheel-pointer {
    position: absolute;
    top: -30px;
    left: 50%;
    transform: translateX(-50%);
    width: 40px;
    height: 40px;
    background: #f44336;
    clip-path: polygon(50% 100%, 0% 0%, 100% 0%);
    z-index: 3;
}

.spinning {
    animation: spinning 4s cubic-bezier(0.1, 0.7, 0.1, 1) forwards;
}

/* Result Styles */
.result-container {
    margin-top: 1.5rem;
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 100%;
    gap: 2rem;
}

.winner-card {
    position: relative;
    padding: 2rem;
    border-radius: 1rem;
    background: linear-gradient(135deg, #fffdf0, #fff8e1);
    box-shadow: 0 10px 30px rgba(255, 215, 0, 0.2), 0 0 10px rgba(255, 215, 0, 0.1);
    border: 2px solid #ffd700;
    text-align: center;
    overflow: hidden;
    width: 80%;
    max-width: 600px;
}

.lucky-number-display {
    display: flex;
    flex-direction: column;
    align-items: center;
    margin-bottom: 1.5rem;
}

.lucky-number-label {
    font-size: 1.2rem;
    color: #ff6f00;
    text-transform: uppercase;
    margin-bottom: 0.5rem;
}

.lucky-number {
    font-size: 4rem;
    font-weight: bold;
    color: #ff6f00;
    background: linear-gradient(135deg, #ffd700, #ff6f00);
    background-clip: text;
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    display: inline-block;
    padding: 0.5rem 1.5rem;
    border-radius: 2rem;
    border: 2px dashed #ffd700;
    background-color: rgba(255, 255, 240, 0.8);
}

.winner-title {
    color: #ff6f00;
    font-size: 1.5rem;
    margin-bottom: 1rem;
}

.winner-email {
    font-size: 2rem;
    font-weight: bold;
    color: #222;
    margin-bottom: 0.5rem;
}

.no-winner {
    color: #f44336;
    font-style: italic;
}

.winner-name {
    font-size: 1.5rem;
    color: #444;
}

/* History styles */
.history-container {
    width: 80%;
    max-width: 600px;
    background: white;
    border-radius: 1rem;
    padding: 1.5rem;
    box-shadow: 0 5px 20px rgba(0, 0, 0, 0.1);
}

.history-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
    border-bottom: 1px solid #f0f0f0;
    padding-bottom: 0.5rem;
}

.history-title {
    color: #ff6f00;
    font-size: 1.3rem;
    margin: 0;
}

.history-list {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    max-height: 300px;
    overflow-y: auto;
}

.history-item {
    display: flex;
    align-items: center;
    padding: 0.75rem;
    background: #f9f9f9;
    border-radius: 0.5rem;
    transition: all 0.2s ease;
}

.history-item:hover {
    background: #f5f5f5;
    transform: translateY(-2px);
    box-shadow: 0 3px 10px rgba(0, 0, 0, 0.05);
}

.history-number {
    font-size: 1.5rem;
    font-weight: bold;
    color: #ff6f00;
    min-width: 50px;
    margin-right: 1rem;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(255, 215, 0, 0.1);
    border-radius: 50%;
    width: 50px;
    height: 50px;
}

.history-info {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    overflow: hidden;
}

.history-email {
    font-size: 1rem;
    font-weight: bold;
    color: #333;
    margin-bottom: 0.25rem;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 100%;
}

.history-name {
    font-size: 0.9rem;
    color: #666;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 100%;
}

/* Confetti animation */
.confetti-container {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    overflow: hidden;
    z-index: 0;
    pointer-events: none;
}

.confetti {
    position: absolute;
    animation: confetti-fall 3s linear forwards;
    z-index: 1;
}

/* Animations */
@keyframes spinning {
    from {
        transform: rotate(0deg);
    }
    to {
        transform: rotate(var(--random-rotation, 1800deg));
    }
}

@keyframes confetti-fall {
    0% { transform: translateY(-100%) rotate(0deg); opacity: 1; }
    100% { transform: translateY(1000%) rotate(360deg); opacity: 0; }
}
</style>