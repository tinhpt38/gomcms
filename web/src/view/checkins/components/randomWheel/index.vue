<template>
    <div class="random-wheel-container my-8">
        <el-button type="success" @click="onLuckyClick" class="mt-2 mb-4" :disabled="isSpinning">
            <span v-if="!isSpinning">Tìm số may mắn</span>
            <span v-else>Đang quay...</span>
        </el-button>
        
        <!-- Wheel without pre-displayed numbers -->
        <div class="wheel-container mb-8" v-if="isSpinning">
            <div class="wheel" ref="wheelRef" :class="{ 'spinning': isSpinning }">
                <div class="wheel-inner">
                    <div v-for="index in 12" :key="index" class="wheel-item" 
                        :style="getWheelItemStyle(index - 1)">
                        <!-- No numbers displayed here -->
                    </div>
                    <div class="wheel-center"></div>
                </div>
                <div class="wheel-pointer"></div>
            </div>
        </div>
        
        <!-- Lucky number popup that appears when a result is found -->
        <div v-if="showLuckyPopup" class="lucky-popup-overlay" @click="closeLuckyPopup">
            <div class="lucky-popup animate__animated animate__zoomIn" @click.stop>
                <div class="confetti-container" v-if="showConfetti">
                    <div v-for="n in 100" :key="n" class="confetti" 
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
                
                <div class="spotlight-container">
                    <div class="spotlight spotlight-1"></div>
                    <div class="spotlight spotlight-2"></div>
                </div>
                
                <div class="lucky-popup-content">
                    <div class="lucky-number">{{ luckyNumber }}</div>
                </div>
                
                <div class="star-container">
                    <div v-for="n in 20" :key="`star-${n}`" class="star"
                        :style="{
                            left: Math.random() * 100 + '%',
                            top: Math.random() * 100 + '%',
                            animationDelay: Math.random() * 2 + 's',
                            animationDuration: 1 + Math.random() * 2 + 's'
                        }">
                    </div>
                </div>
            </div>
        </div>
        
        <!-- History is always visible when there's data -->
        <div class="history-container mt-6" v-if="luckyNumbersHistory.length > 0">
            <div class="history-header">
                <h3 class="history-title">Lịch sử số may mắn</h3>
                <el-button type="danger" size="small" @click="clearHistory" icon="Delete">Xóa lịch sử</el-button>
            </div>
            <div class="history-list">
                <div v-for="(item, idx) in luckyNumbersHistory" :key="idx" class="history-item">
                    <div class="history-number">{{ item.luckyNumber }}</div>
                    <div class="history-details">
                        <div class="history-info" v-if="!item.showDetails && item.email">
                            <span class="history-email">Đã có người trúng</span>
                        </div>
                        <div class="history-info" v-else-if="!item.showDetails && !item.email">
                            <span class="history-email no-winner">Chưa tìm được chủ nhân</span>
                        </div>
                        <div class="history-info expanded" v-if="item.showDetails && item.email">
                            <span class="history-email">{{ item.email }}</span>
                            <span class="history-name">{{ item.fullName || "--" }}</span>
                        </div>
                        <el-button 
                            v-if="item.email" 
                            type="primary" 
                            size="small" 
                            @click="toggleDetails(idx)" 
                            class="details-button">
                            {{ item.showDetails ? 'Ẩn thông tin' : 'Xem thông tin' }}
                        </el-button>
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
const showLuckyPopup = ref(false)
const wheelRef = ref(null)

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
    const angle = (360 / 12) * index;
    return {
        transform: `rotate(${angle}deg) translateY(-125px) rotate(-${angle}deg)`,
        backgroundColor: index % 2 === 0 ? '#ff9800' : '#ffeb3b',
    }
}

// Toggle details visibility for a history item
const toggleDetails = (index) => {
    const item = luckyNumbersHistory.value[index];
    if (item) {
        // Create a new array with the updated item to maintain reactivity
        luckyNumbersHistory.value = luckyNumbersHistory.value.map((historyItem, idx) => {
            if (idx === index) {
                return { ...historyItem, showDetails: !historyItem.showDetails };
            }
            return historyItem;
        });
    }
};

// Close the lucky popup
const closeLuckyPopup = () => {
    showLuckyPopup.value = false;
};

// Start animations when result is shown
const startAnimations = () => {
    showAnimation.value = true;
    showConfetti.value = true;
    showLuckyPopup.value = true;
    
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
        showLuckyPopup.value = false;
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
                // Add the showDetails property to each history item
                luckyNumbersHistory.value = res.data.luckyHistory.map(item => ({
                    ...item,
                    showDetails: false
                }));
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
    // Tải danh sách lịch sử từ server khi component được mount
    findLuckyParticipant({ attendanceId: props.acId }).then(res => {
        if (res.data && res.data.luckyHistory) {
            // Add the showDetails property to each history item
            luckyNumbersHistory.value = res.data.luckyHistory.map(item => ({
                ...item,
                showDetails: false
            }));
            
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

/* Lucky Popup Styles */
.lucky-popup-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(255, 236, 179, 0.85);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 9999;
    cursor: pointer;
}

.lucky-popup {
    position: relative;
    width: 90%;
    max-width: 600px;
    padding: 5rem;
    border-radius: 3rem;
    background: radial-gradient(circle, #ffffff, #fff5e0);
    box-shadow: 0 0 100px rgba(255, 184, 0, 0.6), 0 0 60px rgba(255, 184, 0, 0.4);
    border: 6px solid #ffb700;
    text-align: center;
    overflow: hidden;
    cursor: default;
}

.lucky-popup-content {
    position: relative;
    z-index: 5;
}

.lucky-number {
    font-size: 12rem;
    font-weight: 900;
    color: #ff8c00;
    background: linear-gradient(135deg, #ffb700, #ff8c00, #ff7800);
    background-clip: text;
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    text-shadow: 
        0 0 20px rgba(255, 183, 0, 0.9),
        0 0 40px rgba(255, 183, 0, 0.7),
        0 0 80px rgba(255, 183, 0, 0.5);
    display: inline-block;
    animation: pulse 1.5s infinite, shake 5s infinite;
}

/* Spotlight effects */
.spotlight-container {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    overflow: hidden;
    z-index: 1;
}

.spotlight {
    position: absolute;
    width: 300px;
    height: 300px;
    border-radius: 50%;
    filter: blur(40px);
    opacity: 0.7;
    z-index: 1;
    animation: spotlight 8s infinite linear;
}

.spotlight-1 {
    background: radial-gradient(circle, rgba(255, 183, 0, 0.8), transparent 70%);
    top: -100px;
    left: -100px;
    animation-delay: 0s;
}

.spotlight-2 {
    background: radial-gradient(circle, rgba(255, 140, 0, 0.6), transparent 70%);
    bottom: -100px;
    right: -100px;
    animation-delay: 4s;
}

/* Star effects */
.star-container {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    z-index: 2;
    pointer-events: none;
}

.star {
    position: absolute;
    width: 20px;
    height: 20px;
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='%23FF8C00'%3E%3Cpath d='M12 17.27L18.18 21l-1.64-7.03L22 9.24l-7.19-.61L12 2 9.19 8.63 2 9.24l5.46 4.73L5.82 21z'/%3E%3C/svg%3E");
    background-size: contain;
    opacity: 0.8;
    animation: twinkle 2s infinite;
}

/* Confetti animation */
.confetti-container {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    overflow: hidden;
    z-index: 3;
    pointer-events: none;
}

.confetti {
    position: absolute;
    animation: confetti-fall 5s linear forwards;
    z-index: 1;
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

.history-details {
    display: flex;
    flex: 1;
    flex-direction: column;
    align-items: flex-start;
    overflow: hidden;
}

.history-info {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    margin-bottom: 0.5rem;
    width: 100%;
}

.history-info.expanded {
    background-color: rgba(255, 248, 225, 0.5);
    padding: 0.5rem;
    border-radius: 0.5rem;
    border-left: 3px solid #ff6f00;
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

.details-button {
    align-self: flex-end;
    margin-top: 0.25rem;
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
    100% { transform: translateY(1000%) rotate(720deg); opacity: 0; }
}

@keyframes pulse {
    0% { transform: scale(1); text-shadow: 0 0 20px rgba(255, 183, 0, 0.9), 0 0 40px rgba(255, 183, 0, 0.7), 0 0 80px rgba(255, 183, 0, 0.5); }
    50% { transform: scale(1.1); text-shadow: 0 0 30px rgba(255, 183, 0, 1), 0 0 60px rgba(255, 183, 0, 0.8), 0 0 100px rgba(255, 183, 0, 0.6); }
    100% { transform: scale(1); text-shadow: 0 0 20px rgba(255, 183, 0, 0.9), 0 0 40px rgba(255, 183, 0, 0.7), 0 0 80px rgba(255, 183, 0, 0.5); }
}

@keyframes shake {
    0%, 100% { transform: translateX(0); }
    10%, 30%, 50%, 70%, 90% { transform: translateX(-5px); }
    20%, 40%, 60%, 80% { transform: translateX(5px); }
}

@keyframes spotlight {
    0% { transform: translate(0, 0) scale(1); }
    25% { transform: translate(50px, 50px) scale(1.2); }
    50% { transform: translate(0, 100px) scale(1); }
    75% { transform: translate(-50px, 50px) scale(0.8); }
    100% { transform: translate(0, 0) scale(1); }
}

@keyframes twinkle {
    0%, 100% { opacity: 0.2; transform: scale(0.5); }
    50% { opacity: 1; transform: scale(1.2); }
}
</style>