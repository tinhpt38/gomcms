<template>
    <div>
        <div class="random-wheel-container my-8">
            <!-- Wheel with sectors - clickable -->
            <div class="wheel-container mb-8" @click="onLuckyClick"
                :class="{ 'cursor-not-allowed': isSpinning, 'cursor-pointer': !isSpinning }">
                <canvas id="wheel" width="400" height="400" ref="wheelCanvas"></canvas>
                <div class="wheel-pointer"></div>
                <div v-if="isSpinning" class="wheel-status">Đang quay...</div>
            </div>

            <!-- Lucky number popup that appears when a result is found -->
            <div v-if="showLuckyPopup" class="lucky-popup-overlay" @click="closeLuckyPopup">
                <div class="lucky-popup animate__animated animate__zoomIn" @click.stop>
                    <div class="confetti-container" v-if="showConfetti">
                        <div v-for="n in 100" :key="n" class="confetti" :style="{
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
                        <div v-for="n in 20" :key="`star-${n}`" class="star" :style="{
                            left: Math.random() * 100 + '%',
                            top: Math.random() * 100 + '%',
                            animationDelay: Math.random() * 2 + 's',
                            animationDuration: 1 + Math.random() * 2 + 's'
                        }">
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Lịch sử quay số -->
        <div class="mt-8">
            <h3 class="text-xl font-bold mb-4">Lịch sử quay số</h3>

            <el-table :data="paginatedHistory" style="width: 100%" border stripe v-loading="historyLoading"
                :empty-text="historyLoading ? 'Đang tải...' : 'Không có dữ liệu lịch sử'">
                <el-table-column prop="luckyNumber" label="Số may mắn" min-width="200" align="center" />
                <!-- <el-table-column label="Người trúng" min-width="200">
                    <template #default="scope">
                        <div>
                            <span class="font-bold">{{ scope.row.fullName || 'N/A' }}</span>
                        </div>
                        <div class="text-gray-500 text-sm">
                            {{ scope.row.email || 'N/A' }}
                        </div>
                    </template>
                </el-table-column> -->
                <el-table-column prop="createdAt" label="Thời gian" min-width="180" align="center">
                    <template #default="scope">
                        {{ new Date(scope.row.createdAt).toLocaleString() }}
                    </template>
                </el-table-column>
                <el-table-column label="Hành động" fixed="right" min-width="180" align="center">
                    <template #default="scope">
                        <el-button type="primary" size="small" @click="viewParticipantInfo(scope.row)">
                            Xem thông tin
                        </el-button>
                    </template>
                </el-table-column>

            </el-table>

            <!-- Pagination component -->
            <div class="pagination-container mt-4 flex justify-end" v-if="luckyHistory.length > 0">
                <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize"
                    :page-sizes="[5, 10, 20, 50]" layout="total, sizes, prev, pager, next, jumper"
                    :total="luckyHistory.length" @size-change="handleSizeChange" @current-change="handleCurrentChange"
                    background />
            </div>
        </div>

        <!-- Participant Info Dialog -->
        <el-dialog
            v-model="dialogVisible"
            title="Thông tin người trúng"
            width="520px"
            :show-close="true"
            :close-on-click-modal="true"
            :close-on-press-escape="true"
            class="participant-dialog"
            destroy-on-close
            top="15vh"
        >
            <div v-if="selectedParticipant" class="participant-info">
                <div class="info-section">
                    <div class="flex items-center justify-between mb-4">
                        <h3 class="text-lg font-bold">Thông tin số may mắn</h3>
                        <div class="golden-trophy">🏆</div>
                    </div>
                    
                    <div class="lucky-number-display mb-4 flex justify-center items-center">
                        <div class="lucky-number-tag">{{ selectedParticipant.luckyNumber }}</div>
                    </div>

                    <div class="info-row">
                        <span class="info-label">Thời gian:</span>
                        <span class="info-value time-value">
                            {{ new Date(selectedParticipant.createdAt).toLocaleString() }}
                        </span>
                    </div>
                </div>
                
                <div class="info-section">
                    <div class="flex items-center mb-4">
                        <h3 class="text-lg font-bold">Thông tin người trúng</h3>
                    </div>

                    <div class="info-row">
                        <span class="info-label">Họ tên:</span>
                        <span class="info-value font-medium">{{ selectedParticipant.fullName || '--' }}</span>
                    </div>
                    
                    <div class="info-row">
                        <span class="info-label">Email:</span>
                        <span class="info-value email-value">{{ selectedParticipant.email || '--' }}</span>
                    </div>
                    
                    <!-- Additional participant info section -->
                    <div v-if="selectedParticipant.participantId" class="mt-3 pt-3 border-t border-gray-200">
                        <div class="info-row">
                            <span class="info-label">ID:</span>
                            <span class="info-value id-value">{{ selectedParticipant.participantId }}</span>
                        </div>
                    </div>
                </div>
                
                <div class="mt-4 text-center text-sm text-gray-500">
                    Đã quay số vào {{ new Date(selectedParticipant.createdAt).toLocaleDateString() }}
                </div>
            </div>
            
            <div v-else class="no-data">
                <el-empty description="Không có thông tin" />
            </div>
            
            <template #footer>
                <span class="dialog-footer">
                    <el-button @click="dialogVisible = false">Đóng</el-button>
                </span>
            </template>
        </el-dialog>
    </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
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

// Lịch sử quay số
const luckyHistory = ref([
    // Dữ liệu mẫu - sẽ được ghi đè bởi dữ liệu từ API
    {
        "id": 62,
        "luckyNumber": 13,
        "participantId": null,
        "email": "",
        "fullName": "",
        "createdAt": "2025-05-14T12:14:30.918+07:00"
    }
])

// Pagination variables
const currentPage = ref(1)
const pageSize = ref(10)

// Computed property for paginated history data
const paginatedHistory = computed(() => {
    const startIndex = (currentPage.value - 1) * pageSize.value
    const endIndex = startIndex + pageSize.value
    return luckyHistory.value.slice(startIndex, endIndex)
})

// Pagination event handlers
const handleSizeChange = (val) => {
    pageSize.value = val
    currentPage.value = 1 // Reset to first page when changing page size
}

const handleCurrentChange = (val) => {
    currentPage.value = val
}

// Wheel related states
const isSpinning = ref(false)
const hasResult = ref(false)
const showAnimation = ref(false)
const showConfetti = ref(false)
const showLuckyPopup = ref(false)
const wheelCanvas = ref(null)

// Wheel configuration
const sectors = ref([])
const ctx = ref(null)
const canvasWidth = 400
const canvasHeight = 400
const rad = canvasWidth / 2
const PI = Math.PI
const TAU = 2 * PI
let arc = TAU / 12 // Will be updated based on sectors.length
const friction = 0.992 // Increased from 0.985 for smoother deceleration
let angVel = 0 // Angular velocity
let ang = 0 // Angle in radians
let requestId = null

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

// Initialize wheel sectors with colors
const initSectors = () => {
    const total = 12; // 12 sectors for the wheel
    sectors.value = [];

    // Array of vibrant colors for the wheel
    const colors = [
        '#FF5252', // Red
        '#FF4081', // Pink
        '#E040FB', // Purple
        '#7C4DFF', // Deep Purple
        '#536DFE', // Indigo
        '#448AFF', // Blue
        '#40C4FF', // Light Blue
        '#18FFFF', // Cyan
        '#64FFDA', // Teal
        '#69F0AE', // Green
        '#B2FF59', // Light Green
        '#EEFF41', // Lime
        '#FFFF00', // Yellow
        '#FFD740', // Amber
        '#FFAB40', // Orange
        '#FF6E40'  // Deep Orange
    ];

    for (let i = 1; i <= total; i++) {
        sectors.value.push({
            number: i, // Still keep the number for reference
            color: colors[(i - 1) % colors.length], // Cycle through colors
            textColor: '#FFFFFF' // Not used anymore but kept for compatibility
        });
    }

    arc = TAU / sectors.value.length;
}

// Draw a sector on the wheel
const drawSector = (sector, i) => {
    const angle = arc * i;
    ctx.value.save();

    // Draw sector color
    ctx.value.beginPath();
    ctx.value.fillStyle = sector.color;
    ctx.value.moveTo(rad, rad);
    ctx.value.arc(rad, rad, rad, angle, angle + arc);
    ctx.value.lineTo(rad, rad);
    ctx.value.fill();

    // Add a border between sectors
    ctx.value.beginPath();
    ctx.value.strokeStyle = '#FFFFFF';
    ctx.value.lineWidth = 2;
    ctx.value.moveTo(rad, rad);
    ctx.value.arc(rad, rad, rad, angle, angle + arc);
    ctx.value.lineTo(rad, rad);
    ctx.value.stroke();

    ctx.value.restore();
}

// Get current sector index
const getIndex = () => Math.floor(sectors.value.length - (ang / TAU) * sectors.value.length) % sectors.value.length;

// Rotate the wheel
const rotate = () => {
    if (!ctx.value || !wheelCanvas.value) return;
    // Sử dụng transform với hardware acceleration 
    wheelCanvas.value.style.transform = `rotate(${ang - PI / 2}rad)`;
    wheelCanvas.value.style.backfaceVisibility = 'hidden'; // Cải thiện hiệu suất
}

// Animation frame function
const frame = () => {
    // If stopped and we had been spinning
    if (!angVel && isSpinning.value) {
        isSpinning.value = false;
        const finalSector = sectors.value[getIndex()];

        // Không cần gọi API vì đã gọi ở bước click
        hasResult.value = true;

        // Nếu đã có kết quả từ API, hiển thị kết quả
        if (luckyNumber.value) {
            startAnimations();
        } else {
            ElMessage.error('Không tìm thấy kết quả số may mắn');
        }

        return;
    }

    angVel *= friction; // Decrease velocity by friction
    if (angVel < 0.001) angVel = 0; // Stop when very slow (giảm ngưỡng để vòng quay nhẹ nhàng hơn)
    ang += angVel; // Update angle
    ang %= TAU; // Normalize angle to 0-2π
    rotate();
}

// Animation engine
const engine = () => {
    frame();
    requestId = requestAnimationFrame(engine);
}

// Initialize the wheel
const initWheel = () => {
    if (!wheelCanvas.value) return;

    ctx.value = wheelCanvas.value.getContext('2d');
    initSectors();

    // Draw all sectors
    sectors.value.forEach((sector, i) => drawSector(sector, i));

    // Start the animation engine
    engine();
}

const onLuckyClick = () => {
    if (isSpinning.value) return;

    isSpinning.value = true;
    showLuckyPopup.value = false;
    hasResult.value = false;

    // Gọi API để tìm số may mắn ngay khi ấn vào vòng quay
    searchInfo.value.attendanceId = props.acId;
    findLuckyParticipant(searchInfo.value).then(res => {
        // Lưu kết quả số may mắn và thông tin người chơi
        if (res.data && res.data.luckyNumber) {
            luckyNumber.value = res.data.luckyNumber;
            luckyMember.value = res.data.participant || {};

            // Lưu lịch sử quay số nếu có
            if (res.data.luckyHistory && Array.isArray(res.data.luckyHistory)) {
                luckyHistory.value = res.data.luckyHistory;
                console.log('Lịch sử quay số:', luckyHistory.value);
            }

            // Add the new result to the history
            luckyHistory.value.unshift({
                id: res.data.id || Date.now(),
                luckyNumber: res.data.luckyNumber,
                participantId: res.data.participantId,
                email: res.data.email || '',
                fullName: res.data.fullName || '',
                createdAt: new Date().toISOString()
            });

            // Reset to first page to show the latest result
            currentPage.value = 1;
        }
    }).catch(err => {
        ElMessage.error('Có lỗi xảy ra khi tìm người trúng');
        console.error(err);
        // Dừng vòng quay nếu có lỗi
        isSpinning.value = false;
        return;
    });

    // Calculate random spin - increasing velocity for better spin effect
    const minVelocity = 0.5;
    const maxVelocity = 0.9;
    angVel = minVelocity + Math.random() * (maxVelocity - minVelocity);
}

// Loading state for history table
const historyLoading = ref(false)

// Function to fetch history data from the API
const fetchLuckyHistory = async () => {
    historyLoading.value = true
    try {
        // Replace this with your actual API call
        // Example: const response = await yourApiService.getLuckyHistory()
        // luckyHistory.value = response.data

        // Simulate API delay
        await new Promise(resolve => setTimeout(resolve, 800))

        // For now, we'll use the sample data
        // In a real implementation, you'd fetch this from your API
        luckyHistory.value = [
            {
                "id": 62,
                "luckyNumber": 13,
                "participantId": null,
                "email": "user1@example.com",
                "fullName": "User One",
                "createdAt": "2023-05-14T12:14:30.918+07:00"
            },
            {
                "id": 61,
                "luckyNumber": 7,
                "participantId": null,
                "email": "user2@example.com",
                "fullName": "User Two",
                "createdAt": "2023-05-14T11:45:20.112+07:00"
            },
            {
                "id": 60,
                "luckyNumber": 21,
                "participantId": null,
                "email": "user3@example.com",
                "fullName": "User Three",
                "createdAt": "2023-05-14T10:30:15.654+07:00"
            },
            // Add more sample entries to demonstrate pagination
            // ... more sample data entries ...
        ]
    } catch (error) {
        console.error('Failed to fetch lucky history:', error)
    } finally {
        historyLoading.value = false
    }
}

// Dialog state and selected participant
const dialogVisible = ref(false)
const selectedParticipant = ref(null)

// Function to view participant info
const viewParticipantInfo = (participant) => {
    console.log('Viewing participant:', participant)
    selectedParticipant.value = { ...participant }
    dialogVisible.value = true
}

onMounted(() => {
    // Initialize the wheel
    initWheel();

    // Load previous lucky number if exists
    findLuckyParticipant({ attendanceId: props.acId }).then(res => {
        if (res.data && res.data.luckyNumber) {
            luckyNumber.value = res.data.luckyNumber;
            luckyMember.value = res.data.participant || {};
            hasResult.value = true;

            // Lưu lịch sử quay số nếu có
            if (res.data.luckyHistory && Array.isArray(res.data.luckyHistory)) {
                luckyHistory.value = res.data.luckyHistory;
                console.log('Lịch sử quay số khi mount:', luckyHistory.value);
            }
        }
    }).catch(err => {
        console.error("Không thể tải thông tin số may mắn", err);
    });

    // Fetch history data
    fetchLuckyHistory();
})

onUnmounted(() => {
    // Clean up animation when component is destroyed
    if (requestId) {
        cancelAnimationFrame(requestId);
    }
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
    width: 400px;
    height: 400px;
    margin: 0 auto;
    border-radius: 50%;
    cursor: pointer;
    transition: transform 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
    transform-origin: center center;
}

.wheel-container:hover {
    transform: scale(1.03);
}

.wheel-container.cursor-not-allowed {
    cursor: not-allowed;
}

canvas#wheel {
    position: relative;
    z-index: 1;
    border-radius: 50%;
    box-shadow: 0 0 20px rgba(0, 0, 0, 0.3);
    will-change: transform;
    transform: translateZ(0);
    backface-visibility: hidden;
    perspective: 1000px;
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
    z-index: 2;
}

.wheel-status {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    background: rgba(0, 0, 0, 0.7);
    color: white;
    padding: 10px 20px;
    border-radius: 20px;
    font-weight: bold;
    z-index: 3;
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

/* History Container */
.history-container {
    width: 100%;
    max-width: 800px;
    margin-top: 2rem;
    padding: 1rem;
    background-color: #ffffff;
    border-radius: 10px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
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

/* Pagination styles */
.pagination-container {
    margin-top: 15px;
}

/* Optional: Custom styling for the pagination component */
:deep(.el-pagination) {
    --el-pagination-button-bg-color: var(--el-color-primary-light-9);
    --el-pagination-hover-color: var(--el-color-primary);
}

:deep(.el-pagination .el-pager li.is-active) {
    background-color: var(--el-color-primary);
    color: white;
}

/* Participant info dialog styles */
.participant-dialog {
    :deep(.el-dialog__header) {
        border-bottom: 1px solid #eee;
        padding-bottom: 15px;
        margin-bottom: 0;
        background-color: #f8f9fa;
    }

    :deep(.el-dialog__title) {
        font-weight: bold;
        color: var(--el-color-primary);
        position: relative;
        padding-left: 10px;
        
        &::before {
            content: '';
            position: absolute;
            left: 0;
            top: 50%;
            transform: translateY(-50%);
            width: 4px;
            height: 18px;
            background-color: var(--el-color-primary);
            border-radius: 2px;
        }
    }

    :deep(.el-dialog__body) {
        padding: 25px;
    }

    :deep(.el-dialog__footer) {
        border-top: 1px solid #eee;
        padding-top: 15px;
    }
}

.participant-info {
    .info-section {
        margin-bottom: 25px;
        padding: 18px;
        background-color: #f8f9fa;
        border-radius: 8px;
        border-left: 4px solid var(--el-color-primary-light-5);
        box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
    }

    .info-row {
        display: flex;
        margin-bottom: 14px;
        align-items: center;

        &:last-child {
            margin-bottom: 0;
        }
    }

    .info-label {
        font-weight: 600;
        width: 120px;
        color: #606266;
        position: relative;
        padding-left: 15px;
        
        &::before {
            content: '•';
            position: absolute;
            left: 0;
            color: var(--el-color-primary);
        }
    }

    .info-value {
        flex: 1;
    }

    .lucky-number-tag {
        display: inline-block;
        padding: 6px 15px;
        background: linear-gradient(135deg, #ffd580, #ffb700);
        color: #7d4e00;
        font-weight: bold;
        border-radius: 20px;
        font-size: 18px;
        box-shadow: 0 2px 6px rgba(255, 183, 0, 0.3);
        text-shadow: 0 1px 1px rgba(255, 255, 255, 0.7);
    }

    .lucky-number-display {
        padding: 10px 0;
    }

    .golden-trophy {
        font-size: 24px;
        animation: trophy-glow 2s ease-in-out infinite;
    }

    .time-value, .email-value, .id-value {
        color: #666;
    }
}

.no-data {
    padding: 30px 0;
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
    0% {
        transform: translateY(-100%) rotate(0deg);
        opacity: 1;
    }

    100% {
        transform: translateY(1000%) rotate(720deg);
        opacity: 0;
    }
}

@keyframes pulse {
    0% {
        transform: scale(1);
        text-shadow: 0 0 20px rgba(255, 183, 0, 0.9), 0 0 40px rgba(255, 183, 0, 0.7), 0 0 80px rgba(255, 183, 0, 0.5);
    }

    50% {
        transform: scale(1.1);
        text-shadow: 0 0 30px rgba(255, 183, 0, 1), 0 0 60px rgba(255, 183, 0, 0.8), 0 0 100px rgba(255, 183, 0, 0.6);
    }

    100% {
        transform: scale(1);
        text-shadow: 0 0 20px rgba(255, 183, 0, 0.9), 0 0 40px rgba(255, 183, 0, 0.7), 0 0 80px rgba(255, 183, 0, 0.5);
    }
}

@keyframes shake {

    0%,
    100% {
        transform: translateX(0);
    }

    10%,
    30%,
    50%,
    70%,
    90% {
        transform: translateX(-5px);
    }

    20%,
    40%,
    60%,
    80% {
        transform: translateX(5px);
    }
}

@keyframes spotlight {
    0% {
        transform: translate(0, 0) scale(1);
    }

    25% {
        transform: translate(50px, 50px) scale(1.2);
    }

    50% {
        transform: translate(0, 100px) scale(1);
    }

    75% {
        transform: translate(-50px, 50px) scale(0.8);
    }

    100% {
        transform: translate(0, 0) scale(1);
    }
}

@keyframes twinkle {

    0%,
    100% {
        opacity: 0.2;
        transform: scale(0.5);
    }

    50% {
        opacity: 1;
        transform: scale(1.2);
    }
}

@keyframes trophy-glow {
    0%, 100% {
        text-shadow: 0 0 0px gold;
        transform: scale(1);
    }
    50% {
        text-shadow: 0 0 10px gold, 0 0 20px gold;
        transform: scale(1.1);
    }
}
</style>