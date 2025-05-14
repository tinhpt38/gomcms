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

            <!-- Hiển thị thông báo cập nhật lịch sử -->
            <div v-if="isSpinning" class="updating-history-message">
                <!-- <span class="loading-icon">🔄</span> Đang quay số, lịch sử sẽ được cập nhật sau khi vòng quay dừng lại... -->
            </div>

            <!-- Thêm ô tìm kiếm số -->
            <div class="search-container mb-4">
                <div class="search-title mb-2 text-center">
                    <span class="search-icon">🔍</span>
                    <span class="ml-2">Tìm kiếm số may mắn</span>
                </div>
                <el-input v-model="searchQuery" placeholder="Nhập số may mắn để tìm kiếm" clearable @clear="resetSearch"
                    @input="handleSearch">
                    <template #prefix>
                        <el-icon>
                            <Search />
                        </el-icon>
                    </template>
                </el-input>
                <div v-if="searchQuery && filteredHistory.length > 0" class="search-results-info">
                    Tìm thấy <strong>{{ filteredHistory.length }}</strong> kết quả
                </div>
                <div v-if="searchQuery && filteredHistory.length === 0" class="search-no-results">
                    Không tìm thấy kết quả nào cho "<strong>{{ searchQuery }}</strong>"
                </div>
            </div>

            <el-table :data="paginatedHistory" style="width: 100%" border stripe v-loading="historyLoading"
                :empty-text="historyLoading ? 'Đang tải...' : 'Vui lòng quay số để xem kết quả và lịch sử quay số!'"
                :class="['history-table', { 'updated': historyJustUpdated }]">
                <el-table-column prop="luckyNumber" label="Số may mắn" min-width="200" align="center">
                    <template #default="scope">
                        <div class="lucky-number-container">
                            <div class="lucky-number-tag">
                                {{ scope.row.luckyNumber }}
                            </div>
                        </div>
                    </template>
                </el-table-column>
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
            <div class="pagination-container mt-4 flex justify-end" v-if="filteredHistory.length > 0">
                <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize"
                    :page-sizes="[5, 10, 20, 50]" layout="total, sizes, prev, pager, next, jumper"
                    :total="filteredHistory.length" @size-change="handleSizeChange"
                    @current-change="handleCurrentChange" background />
            </div>
        </div>

        <!-- Participant Info Dialog -->
        <el-dialog v-model="dialogVisible" title="Thông tin người trúng" width="520px" :show-close="true"
            :close-on-click-modal="true" :close-on-press-escape="true" class="participant-dialog" destroy-on-close
            top="15vh">
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
                    <!-- <div v-if="selectedParticipant.participantId" class="mt-3 pt-3 border-t border-gray-200">
                        <div class="info-row">
                            <span class="info-label">ID:</span>
                            <span class="info-value id-value">{{ selectedParticipant.participantId }}</span>
                        </div>
                    </div> -->
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
import { Search } from '@element-plus/icons-vue'
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

// Lịch sử quay số - khởi tạo rỗng
const luckyHistory = ref([])

// Pagination variables
const currentPage = ref(1)
const pageSize = ref(10)

// Search functionality
const searchQuery = ref('')

// Computed property for filtered history data
const filteredHistory = computed(() => {
    if (!searchQuery.value) {
        return luckyHistory.value
    }
    const query = searchQuery.value.trim().toLowerCase()
    return luckyHistory.value.filter(item =>
        item.luckyNumber.toString().includes(query)
    )
})

// Computed property for paginated history data
const paginatedHistory = computed(() => {
    const startIndex = (currentPage.value - 1) * pageSize.value
    const endIndex = startIndex + pageSize.value
    return filteredHistory.value.slice(startIndex, endIndex)
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

// Store history update state
const historyJustUpdated = ref(false);

// Start animations when result is shown
const startAnimations = () => {
    // Hiển thị các hiệu ứng
    showAnimation.value = true;
    showConfetti.value = true;
    showLuckyPopup.value = true;

    // Cập nhật lịch sử quay số sau khi wheel đã dừng
    if (tempApiResult.value) {
        // Lưu lịch sử quay số nếu có từ API
        if (tempApiResult.value.luckyHistory && Array.isArray(tempApiResult.value.luckyHistory)) {
            luckyHistory.value = tempApiResult.value.luckyHistory;
            console.log('Đã cập nhật lịch sử quay số:', luckyHistory.value);
        }

        // Thêm kết quả mới vào lịch sử
        const newResult = {
            id: tempApiResult.value.id || Date.now(),
            luckyNumber: tempApiResult.value.luckyNumber,
            participantId: tempApiResult.value.participantId,
            email: tempApiResult.value.email || '',
            fullName: tempApiResult.value.fullName || '',
            createdAt: new Date().toISOString()
        };

        // Kiểm tra xem kết quả này đã tồn tại trong lịch sử chưa
        const existingIndex = luckyHistory.value.findIndex(item =>
            item.id === newResult.id ||
            (item.luckyNumber === newResult.luckyNumber &&
                new Date(item.createdAt).toDateString() === new Date(newResult.createdAt).toDateString())
        );

        if (existingIndex === -1) {
            luckyHistory.value.unshift(newResult);
            console.log('Đã thêm kết quả mới vào lịch sử:', newResult);

            // Đánh dấu là lịch sử vừa được cập nhật để hiển thị animation
            historyJustUpdated.value = true;
            setTimeout(() => {
                historyJustUpdated.value = false;
            }, 2000);
        }
    }

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

        // Đợi cho vòng quay dừng hoàn toàn rồi mới hiển thị kết quả
        hasResult.value = true;

        // Nếu đã có kết quả từ API, hiển thị kết quả và cập nhật lịch sử
        if (luckyNumber.value) {
            // Thêm dữ liệu vào lịch sử chỉ khi vòng quay dừng lại
            if (luckyHistory.value.length > 0) {
                // Đảm bảo kết quả mới nhất được hiển thị ở đầu tiên
                currentPage.value = 1;
            }

            // Bắt đầu hiệu ứng
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

// Store temporary API result
const tempApiResult = ref(null);

const onLuckyClick = () => {
    if (isSpinning.value) return;

    isSpinning.value = true;
    showLuckyPopup.value = false;
    hasResult.value = false;

    // Reset lucky number and member before making the API call
    luckyNumber.value = null;
    luckyMember.value = { email: null, fullName: null };
    // Reset temporary result
    tempApiResult.value = null;

    // Gọi API để tìm số may mắn khi ấn vào vòng quay
    searchInfo.value.attendanceId = props.acId;
    findLuckyParticipant(searchInfo.value).then(res => {
        // Lưu kết quả số may mắn và thông tin người chơi
        if (res.data && res.data.luckyNumber) {
            // Lưu kết quả vào biến tạm, chưa cập nhật vào luckyHistory
            luckyNumber.value = res.data.luckyNumber;
            luckyMember.value = res.data.participant || {};
            tempApiResult.value = res.data;

            // Không cập nhật history ở đây, mà sẽ cập nhật sau khi vòng quay dừng lại
            console.log('Đã nhận kết quả từ API, đợi vòng quay dừng lại để hiển thị');
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
        // Note: In a real implementation, you'd fetch history from a dedicated API endpoint
        // For now, we'll just set loading to false without adding any sample data
        await new Promise(resolve => setTimeout(resolve, 300))

        // Not adding any sample data anymore - history will be populated only after spinning
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

// Search functionality handlers
const handleSearch = () => {
    // Reset to first page when searching
    currentPage.value = 1
    // The actual filtering happens in the filteredHistory computed property
}

const resetSearch = () => {
    searchQuery.value = ''
    currentPage.value = 1
}

onMounted(() => {
    // Initialize the wheel
    initWheel();

    // Fetch history data only (without calling the findLuckyParticipant API)
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

    /* Thêm style cho lucky number trong bảng */
    .el-table .lucky-number-tag {
        font-size: 20px;
        font-weight: 800;
        padding: 8px 16px;
        background: linear-gradient(135deg, #fff1c1, #ffb700);
        color: #7d4e00;
        border-radius: 50px;
        box-shadow: 0 4px 12px rgba(255, 183, 0, 0.4);
        text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
        min-width: 50px;
        text-align: center;
        transition: all 0.3s ease;
        position: relative;
        animation: lucky-number-glow 2s infinite;
    }

    .el-table .lucky-number-tag:hover {
        transform: scale(1.1);
        box-shadow: 0 6px 15px rgba(255, 183, 0, 0.6);
    }

    /* Thêm hiệu ứng ánh sáng trước và sau số */
    .el-table .lucky-number-tag::before,
    .el-table .lucky-number-tag::after {
        content: '✨';
        position: absolute;
        top: 50%;
        transform: translateY(-50%);
        font-size: 14px;
        color: #ffb700;
        opacity: 0.8;
    }

    .el-table .lucky-number-tag::before {
        left: -5px;
        animation: sparkle 1.5s infinite ease-in-out;
    }

    .el-table .lucky-number-tag::after {
        right: -5px;
        animation: sparkle 1.5s infinite ease-in-out 0.5s;
    }

    /* Keyframe cho hiệu ứng phát sáng */
    @keyframes lucky-number-glow {

        0%,
        100% {
            box-shadow: 0 4px 12px rgba(255, 183, 0, 0.4);
        }

        50% {
            box-shadow: 0 4px 20px rgba(255, 183, 0, 0.7);
        }
    }

    @keyframes sparkle {

        0%,
        100% {
            opacity: 0.5;
            transform: translateY(-50%) scale(0.8);
        }

        50% {
            opacity: 1;
            transform: translateY(-50%) scale(1.2);
        }
    }

    .lucky-number-display {
        padding: 10px 0;
    }

    .golden-trophy {
        font-size: 24px;
        animation: trophy-glow 2s ease-in-out infinite;
    }

    .time-value,
    .email-value,
    .id-value {
        color: #666;
    }
}

.no-data {
    padding: 30px 0;
}

/* Container cho số may mắn trong bảng */
.lucky-number-container {
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 5px 0;
    position: relative;
}

/* Hiệu ứng ánh sáng xung quanh container */
.lucky-number-container::before {
    content: '';
    position: absolute;
    top: 50%;
    left: 50%;
    width: 80%;
    height: 80%;
    background: radial-gradient(circle, rgba(255, 215, 0, 0.2) 0%, rgba(255, 215, 0, 0) 70%);
    transform: translate(-50%, -50%);
    border-radius: 50%;
    z-index: 0;
    animation: pulse-bg 3s infinite;
}

@keyframes pulse-bg {

    0%,
    100% {
        opacity: 0.2;
        transform: translate(-50%, -50%) scale(1);
    }

    50% {
        opacity: 0.5;
        transform: translate(-50%, -50%) scale(1.2);
    }
}

/* Tùy chỉnh bảng lịch sử */
.el-table {
    border-radius: 8px;
    overflow: hidden;
    box-shadow: 0 5px 15px rgba(0, 0, 0, 0.05);
    border: 1px solid rgba(255, 215, 0, 0.2);
}

/* Tùy chỉnh header của bảng */
:deep(.el-table__header) {
    background: linear-gradient(180deg, #fff8e6, #fff5d6);
}

:deep(.el-table__header .cell) {
    font-weight: bold;
    font-size: 15px;
}

/* Tùy chỉnh tiêu đề của phần lịch sử quay số */
.mt-8>h3 {
    position: relative;
    display: inline-block;
    padding-bottom: 8px;
    margin-bottom: 16px;
}

.mt-8>h3::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 0;
    width: 100%;
    height: 3px;
    background: linear-gradient(90deg, #ffb700, transparent);
    border-radius: 3px;
}

/* Tùy chỉnh tiêu đề cột số may mắn */
:deep(.el-table__header-wrapper .el-table__cell[data-label="Số may mắn"]) {
    background: linear-gradient(90deg, rgba(255, 215, 0, 0.1), rgba(255, 215, 0, 0.3), rgba(255, 215, 0, 0.1));
    font-weight: bold;
    color: #d48806;
}

/* Tùy chỉnh hàng trong bảng khi hover */
:deep(.el-table__row:hover) {
    background-color: rgba(255, 248, 230, 0.7) !important;
}

/* Hiệu ứng hover cho số may mắn */
.el-table__row:hover .lucky-number-tag {
    transform: scale(1.1) rotate(5deg);
    box-shadow: 0 8px 20px rgba(255, 183, 0, 0.5);
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

    0%,
    100% {
        text-shadow: 0 0 0px gold;
        transform: scale(1);
    }

    50% {
        text-shadow: 0 0 10px gold, 0 0 20px gold;
        transform: scale(1.1);
    }
}

/* Search container styling */
.search-container {
    position: relative;
    margin-bottom: 20px;
    max-width: 500px;
    margin-left: auto;
    margin-right: auto;
}

.search-title {
    font-size: 16px;
    font-weight: 600;
    color: #666;
    margin-bottom: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
}

.search-icon {
    font-size: 18px;
    animation: bounce 1s infinite alternate;
}

.search-results-info {
    margin-top: 10px;
    text-align: center;
    font-size: 14px;
    color: #666;
    background-color: #f8f9fa;
    padding: 8px;
    border-radius: 4px;
    animation: fadeIn 0.5s ease-in;
}

.search-no-results {
    margin-top: 10px;
    text-align: center;
    font-size: 14px;
    color: #ff9800;
    background-color: #fff8e1;
    padding: 8px;
    border-radius: 4px;
    border-left: 3px solid #ff9800;
    animation: fadeIn 0.5s ease-in;
}

.search-container :deep(.el-input__wrapper) {
    border-radius: 50px;
    padding: 2px 15px;
    box-shadow: 0 3px 10px rgba(255, 184, 0, 0.1);
    border: 1px solid rgba(255, 183, 0, 0.3);
    transition: all 0.3s ease;
}

.search-container :deep(.el-input__wrapper:hover),
.search-container :deep(.el-input__wrapper.is-focus) {
    box-shadow: 0 5px 15px rgba(255, 184, 0, 0.2);
    border-color: rgba(255, 183, 0, 0.5);
}

.search-container :deep(.el-input__prefix) {
    color: #ffb700;
}

.search-container :deep(.el-input__inner::placeholder) {
    color: #999;
    font-style: italic;
}

.search-container :deep(.el-input__inner) {
    font-size: 15px;
    color: #333;
}

/* Add animation for the clear button */
.search-container :deep(.el-input__suffix) {
    cursor: pointer;
    transition: transform 0.2s ease;
}

.search-container :deep(.el-input__suffix:hover) {
    transform: scale(1.2);
}

@keyframes bounce {
    0% {
        transform: translateY(0);
    }

    100% {
        transform: translateY(-3px);
    }
}

@keyframes fadeIn {
    from {
        opacity: 0;
        transform: translateY(-5px);
    }

    to {
        opacity: 1;
        transform: translateY(0);
    }
}

/* Thêm style cho thông báo đang cập nhật lịch sử */
.updating-history-message {
    text-align: center;
    padding: 10px;
    margin-bottom: 15px;
    background: #fff9e6;
    border-radius: 8px;
    border-left: 4px solid #ffb700;
    color: #d48806;
    animation: fadeIn 0.5s ease-in;
    font-size: 15px;
    display: flex;
    align-items: center;
    justify-content: center;
}

.loading-icon {
    animation: spin 1.5s linear infinite;
    display: inline-block;
    margin-right: 8px;
    font-size: 18px;
}

@keyframes spin {
    0% {
        transform: rotate(0deg);
    }

    100% {
        transform: rotate(360deg);
    }
}

/* Thêm style để highlight bảng lịch sử khi có cập nhật mới */
.history-table.updated {
    animation: highlight-update 1.5s ease;
}

@keyframes highlight-update {
    0% {
        box-shadow: 0 0 0px rgba(255, 183, 0, 0);
    }

    50% {
        box-shadow: 0 0 20px rgba(255, 183, 0, 0.8);
    }

    100% {
        box-shadow: 0 0 0px rgba(255, 183, 0, 0);
    }
}
</style>