<template>
  <div class="bg-white relative">
    <main class="isolate">
      <header ref="header" class="sticky top-0 z-50">
        <div class="bg-white  mx-auto px-4 sm:px-6 lg:px-52">
          <div class="container mx-auto py-3">
            <div class="flex items-center justify-between">
              <a href="/" class="flex items-center">
                <img src="/logo-dlu-full.webp" alt="Trường Đại học Đà Lạt" class="h-14 sm:h-16 lg:h-20 w-auto">
              </a>
              <div class="hidden md:flex items-center gap-4 lg:gap-7" />
              <!-- <button
                class="hidden md:block bg-[#79a227] text-white text-base lg:text-base py-2 px-6 lg:py-3 lg:px-8 rounded-xl outline-none border-none cursor-pointer"
                @click="redirectToHistory">
                Lịch sử
              </button> -->
              <!-- <button
                class="md:hidden bg-[#79a227] text-white text-base lg:text-base py-2 px-6 lg:py-3 lg:px-8 rounded-xl outline-none border-none cursor-pointer"
                @click="redirectToHistory">
                Lịch sử
              </button> -->
            </div>
          </div>
        </div>
        <div class="flex w-full h-[3px]">
          <div class="w-[10%] bg-[#E67F32]"></div>
          <div class="w-[20%] bg-[#607F23]"></div>
          <div class="w-[30%] bg-[#98BE3B]"></div>
          <div class="w-[50%] bg-[#C2D88B]"></div>
        </div>
      </header>

      <!-- Hero section -->
      <div class="relative isolate -z-10">
        <svg
          class="absolute inset-x-0 top-0 -z-10 h-[64rem] w-full stroke-gray-200 [mask-image:radial-gradient(32rem_32rem_at_center,white,transparent)]"
          aria-hidden="true">
          <defs>
            <pattern id="1f932ae7-37de-4c0a-a8b0-a6e3b4d44b84" width="200" height="200" x="50%" y="-1"
              patternUnits="userSpaceOnUse">
              <path d="M.5 200V.5H200" fill="none" />
            </pattern>
          </defs>
          <svg x="50%" y="-1" class="overflow-visible fill-gray-50">
            <path d="M-200 0h201v201h-201Z M600 0h201v201h-201Z M-400 600h201v201h-201Z M200 800h201v201h-201Z"
              stroke-width="0" />
          </svg>
          <rect width="100%" height="100%" stroke-width="0" fill="url(#1f932ae7-37de-4c0a-a8b0-a6e3b4d44b84)" />
        </svg>
        <div
          class="absolute left-1/2 right-0 top-0 -z-10 -ml-24 transform-gpu overflow-hidden blur-3xl lg:ml-24 xl:ml-48"
          aria-hidden="true">
          <div class="aspect-[801/1036] w-[50.0625rem] bg-gradient-to-tr from-[#ff80b5] to-[#9089fc] opacity-30"
            style="clip-path: polygon(63.1% 29.5%, 100% 17.1%, 76.6% 3%, 48.4% 0%, 44.6% 4.7%, 54.5% 25.3%, 59.8% 49%, 55.2% 57.8%, 44.4% 57.2%, 27.8% 47.9%, 35.1% 81.5%, 0% 97.7%, 39.2% 100%, 35.2% 81.4%, 97.2% 52.8%, 63.1% 29.5%)" />
        </div>
        <div class="overflow-hidden">
          <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 pb-[32px] pt-[32px]">
            <div class="mx-auto max-w-2xl gap-x-14 lg:mx-0 lg:flex lg:max-w-none lg:items-center">
              <div class="w-full max-w-xl lg:shrink-0 xl:max-w-2xl">
                <h1 class="text-3xl sm:text-4xl lg:text-6xl font-bold tracking-tight text-[#514C39]">
                  Hệ thống điểm danh<br><span class="text-[#7BA227]">Trường Đại học Đà Lạt</span>.
                </h1>
                <p
                  class="relative mt-6 text-base sm:text-lg font-bold leading-8 text-[#E67F32] sm:max-w-md lg:max-w-none">
                  Bạn hãy đăng nhập bằng tài khoản Email của bạn với Google để đăng nhập.
                </p>
                <!-- <p class="mt-6 text-lg leading-8 text-gray-600 px-2 mb-6">
                  Bạn hãy Sử dụng 1 trong 2 WiFi: DLU Student hoặc DLU Teacher để thực hiện điểm danh
                </p> -->
                <p class="relative mt-1 text-sm sm:text-base italic leading-8 text-gray-500 sm:max-w-md lg:max-w-none">
                  {{ isSupported ? "Trình duyệt hỗ trợ lấy vị trí: " : "Trình duyệt không hỗ trợ lấy vị trí" }}
                  {{ coords.latitude + ", " + coords.longitude }}
                </p>
                <div v-if="attendance.title != null" class="mt-4 p-2 rounded shadow-slate-400">
                  <h3>{{ attendance.title }}</h3>
                  <!-- Lucky Number feature -->
                <LuckyNumberDisplay 
                  :lucky-number="apiResponse?.data?.luckyNumber" 
                  :current-checkins="checkinCount"
                  :required-checkins="attendance.luckyShowAfterMinCount || 0"
                  :is-new-lucky-number="isNewLuckyNumber"
                  :show-progress="attendance.useLuckyNumber"
                />
                  <div class="mt-4" v-if="conditionData.length > 0">
                    <div class="text-base text-gray-900">
                      Danh sách điều kiện điểm danh
                    </div>
                    <dl class="divide-y divide-gray-100 mt-2 mb-2">
                      <div v-for="(item, key) in conditionData" :key="key"
                        class="px-2 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-2  bg-slate-50 mt-2">
                        <dt class="text-sm font-medium text-base text-gray-900">
                          Điều kiện {{ key + 1 }}
                        </dt>
                        <dd class="my-1 text-sm text-base text-gray-700 sm:col-span-2 sm:mt-0">
                          {{ conditionString(item) }}
                        </dd>
                        <el-tag v-if="item.IsPass" effect="dark" type="success" class="flex items-center gap-1">
                          <span>Bạn đã điểm danh thành công {{item.counter + 1 }} lần</span>
                        </el-tag>
                        <el-tag v-if="item.IsPass && item.ShowLuckyNumber" class="ml-2 mt-2" type="warning">
                          <i class="el-icon-star-on mr-1"></i>Điều kiện may mắn
                        </el-tag>
                        <el-tag v-if="!item.IsPass && item.ShowLuckyNumber" class="ml-2 mt-2" type="info">
                          <i class="el-icon-star-on mr-1"></i>Điều kiện nhận số may mắn
                        </el-tag>
                        <el-tag v-if="!item.IsPass" effect="dart" type="danger">
                          Bạn chưa điểm danh
                        </el-tag>
                        <el-text v-if="!item.IsPass" effect="dart" type="primary">
                          <strong>{{ item.Message }}</strong>
                        </el-text>
                      </div>
                    </dl>
                  </div>
                </div>

                

                <GoogleLogin class="my-4" :callback="callback" :error="gError" prompt />
              </div>
              <div class="mt-14 flex justify-end gap-8 sm:-mt-44 sm:justify-start sm:pl-20 lg:mt-0 lg:pl-0">
                <div
                  class="ml-auto w-44 flex-none space-y-8 pt-32 sm:ml-0 sm:pt-80 lg:order-last lg:pt-36 xl:order-none xl:pt-80">
                  <div class="relative">
                    <img src="/dlu1.jpg" alt=""
                      class="aspect-[2/3] w-full rounded-xl bg-gray-900/5 object-cover shadow-lg">
                    <div class="pointer-events-none absolute inset-0 rounded-xl ring-1 ring-inset ring-gray-900/10" />
                  </div>
                </div>
                <div class="mr-auto w-44 flex-none space-y-8 sm:mr-0 sm:pt-52 lg:pt-36">
                  <div class="relative">
                    <img src="/dlu2.jpg" alt=""
                      class="aspect-[2/3] w-full rounded-xl bg-gray-900/5 object-cover shadow-lg">
                    <div class="pointer-events-none absolute inset-0 rounded-xl ring-1 ring-inset ring-gray-900/10" />
                  </div>
                  <div class="relative">
                    <img src="/dlu3.jpg" alt=""
                      class="aspect-[2/3] w-full rounded-xl bg-gray-900/5 object-cover shadow-lg">
                    <div class="pointer-events-none absolute inset-0 rounded-xl ring-1 ring-inset ring-gray-900/10" />
                  </div>
                </div>
                <div class="w-44 flex-none space-y-8 pt-32 sm:pt-0">
                  <div class="relative">
                    <img src="/dlu4.jpg" alt=""
                      class="aspect-[2/3] w-full rounded-xl bg-gray-900/5 object-cover shadow-lg">
                    <div class="pointer-events-none absolute inset-0 rounded-xl ring-1 ring-inset ring-gray-900/10" />
                  </div>
                  <div class="relative">
                    <img src="/dlu5.jpg" alt=""
                      class="aspect-[2/3] w-full rounded-xl bg-gray-900/5 object-cover shadow-lg">
                    <div class="pointer-events-none absolute inset-0 rounded-xl ring-1 ring-inset ring-gray-900/10" />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { nextTick, onMounted, ref, onUnmounted, watch } from 'vue'
import { publicAttendanceCheckIn } from '@/api/checkins/attendanceCheckIn'
import { ElMessage, ElMessageBox, ElNotification } from 'element-plus';
import { useGeolocation } from '@vueuse/core'
import { decodeCredential } from 'vue3-google-login'
import LuckyNumberDisplay from '@/components/LuckyNumberDisplay.vue'

import { useRoute } from 'vue-router';
import { formatDateTime } from '@/utils/format';

import FingerprintJS from '@fingerprintjs/fingerprintjs';


const route = useRoute()


defineOptions({
  name: "Checkins",
})

const { coords, locatedAt, error, isSupported, isLoading } = useGeolocation()

const data = ref({
  email: null,
  code: null,
  lat: null,
  lng: null,
})

const header = ref(null)


const callback = async (response) => {
  const userData = decodeCredential(response.credential)
  data.value.email = userData.email
  data.value.fullName = (userData?.given_name || "") + ' ' + (userData?.family_name || "")
  await requestCheckin()
}

const gError = (error) => {
  //console.log("Handle the error", error)
}

function toBinaryStr(str) {
  const encoder = new TextEncoder();
  // 1: split the UTF-16 string into an array of bytes
  const charCodes = encoder.encode(str);
  // 2: concatenate byte data to create a binary string
  return String.fromCharCode(...charCodes);
}


async function getFingerprint() {
  const fpPromise = FingerprintJS.load();
  const fp = await fpPromise;
  const result = await fp.get();
  return result.visitorId;
}

const fingerPrint = () => {
  getFingerprint().then(visitorId => {
    data.value.visitorId = visitorTemplate(visitorId)
  });
}

fingerPrint()

const visitorTemplate = (val) => {
  return "dlu_activities_20422_5BS:W`A8nF<J6Y{V4Nv.r!Je_" + val
}

const keyRandom = 'E;>YIws8_DdsSMG£sL£@lq8E<(O?Sc5'
const encodeVal = (data) => {
  const jsonString = JSON.stringify(data);
  // //console.log("jsonString: ", jsonString)
  const encodedData = btoa(toBinaryStr(jsonString));
  // //console.log("encodedData: ", encodedData)
  return keyRandom + "_" + encodedData
}

const redirectToHistory = () => {
  window.location.href = '/history'
}



const conditionString = (item) => {
  var condition = {
    group: item?.group?.name,
    area: item?.area?.area?.name,
    startAt: item.startAt ? formatDateTime(item.startAt) : null,
    endAt: item.endAt ? formatDateTime(item.endAt) : null,
    isPass: item.isPass
  }
  if (condition.group != null && condition.area != null && condition.startAt != null && condition.endAt != null) {
    return `Nhóm ${condition.group}, tại ${condition.area}, từ ${condition.startAt} đến ${condition.endAt}`
  } else if (condition.group != null && condition.area == null && condition.startAt != null) {
    return `Nhóm ${condition.group}, từ ${condition.startAt} đến ${condition.endAt}`
  } else if (condition.group != null && condition.area != null) {
    return `Nhóm ${condition.group}, tại ${condition.area}`
  } else if (condition.group != null) {
    return `Nhóm ${condition.group}`
  } else if (condition.area != null) {
    return `Tại ${condition.area}`
  } else if (condition.area != null && condition.startAt != null && condition.endAt != null) {
    return `Tại ${condition.area}, từ ${condition.startAt} đến ${condition.endAt}`
  } else if (condition.area != null && condition.startAt != null) {
    return `Tại ${condition.area}, từ ${condition.startAt}`
  } else if (condition.startAt != null && condition.endAt != null) {
    return `Từ ${condition.startAt} đến ${condition.endAt}`
  } else if (condition.startAt != null) {
    return `Từ ${condition.startAt}`
  } else if (condition.endAt != null) {
    return `Hạn cuối ${condition.endAt}`
  }
}

const conditionData = ref([])
const attendance = ref({})
const apiResponse = ref({data: {}}) // Store the entire response for reference
const showBounceAnimation = ref(false) // Control bounce animation
const showConfetti = ref(false) // Control confetti animation
const checkinsRemaining = ref(0) // Remaining check-ins needed for lucky number
const checkinsProgress = ref(0) // Progress percentage towards lucky number
const isNewLuckyNumber = ref(false) // Whether the lucky number was just received

// Function to generate random colors for confetti
const getRandomColor = () => {
  const colors = ['#ffd700', '#ffcc33', '#ffc107', '#ffb300', '#ffa000', '#ff8f00', '#ff6f00']
  return colors[Math.floor(Math.random() * colors.length)]
}

const requestCheckin = async () => {
  if (route.query?.c == null) {
    ElNotification("Không có điểm danh nào đang hiện hành")
    return
  }

  if (typeof (coords.value.latitude) == "Infinity" || typeof (coords.value.longitude) == "Infinity") {
    ElMessage.error("Không thể lấy vị trí của bạn")
    return
  }

  if (!coords.value.latitude || !coords.value.longitude) {
    ElMessage.error("Không thể lấy vị trí của bạn")
    return
  }

  if (coords.value.latitude === 0 || coords.value.longitude === 0) {
    ElMessage.error("Không thể lấy vị trí của bạn")
    return
  }


  if (coords.value.latitude === null || coords.value.longitude === null) {
    ElMessage.error("Không thể lấy vị trí của bạn")
    return
  }



  data.value.lat = coords.value.latitude
  data.value.lng = coords.value.longitude
  data.value.accuracy = coords.value.accuracy
  data.value.code = route.query.c

  var encodedData = encodeVal(data.value)
  // var res = await publicAttendanceCheckIn({ ...data.value })

  var apiResult = await publicAttendanceCheckIn({ data: encodedData })
  console.log("res: ", apiResult)
  apiResponse.value = apiResult // Store the entire response
  
  if (apiResult.code == 0) {
    if (apiResult.data.conditions != null) {
      conditionData.value = apiResult.data.conditions.filter((condition, index, self) =>
        index === self.findIndex((c) => c.ID === condition.ID)
      )

      // Ensure counter information is correctly processed
      // Conditions may already have counter data from backend
      conditionData.value.forEach(condition => {
        // Ensure counter property exists and has a valid value
        if (typeof condition.counter === 'undefined') {
          condition.counter = 0
        }
        
        // Normalize ShowLuckyNumber/showLuckyNumber property name
        if (condition.ShowLuckyNumber !== undefined && condition.showLuckyNumber === undefined) {
          condition.showLuckyNumber = condition.ShowLuckyNumber
        } else if (condition.showLuckyNumber !== undefined && condition.ShowLuckyNumber === undefined) {
          condition.ShowLuckyNumber = condition.showLuckyNumber
        }
        
        // Normalize message property name
        if (condition.Message !== undefined && condition.msg === undefined) {
          condition.msg = condition.Message
        } else if (condition.msg !== undefined && condition.Message === undefined) {
          condition.Message = condition.msg
        }
      })
      console.log("conditionData: ", conditionData.value)
    }
    attendance.value = apiResult.data.attendance
    // Retrieve total check-in count and lucky number from response
    var checkinCount = 0;
    var luckyNumber = apiResponse.value.data.luckyNumber ?? null
    
    // Get check-in counter from the first check-in if available
    if (apiResponse.value.data.checkins && apiResponse.value.data.checkins.length > 0) {
      checkinCount = apiResponse.value.data.checkins.reduce((total, checkin) => total + (checkin.counter || 0), 0);
    }
    
    // Calculate condition stats if conditions exist
    var passcount = 0;
    var totalConditions = 0;
    var totalCheckIns = 0;
    
    if (conditionData.value && conditionData.value.length > 0) {
      passcount = conditionData.value.reduce((count, item) => {
        return item.isPass ? count + 1 : count;
      }, 0);
      
      totalConditions = conditionData.value.length;
      
      // Get total check-in count including repetitions
      totalCheckIns = conditionData.value.reduce((total, item) => {
        return item.isPass ? total + (item.counter || 1) : total;
      }, 0);
    }
    
    // Update the lucky number progress data
    if (attendance.value.useLuckyNumber && !luckyNumber) {
      const requiredCheckins = attendance.value.luckyShowAfterMinCount || 0;
      if (requiredCheckins > 0 && checkinCount < requiredCheckins) {
        checkinsRemaining.value = requiredCheckins - checkinCount;
        checkinsProgress.value = (checkinCount / requiredCheckins) * 100;
      }
    }
    
    // Build a structured message for better readability
    let msgComponents = [];
    
    // 1. Basic success message
    if (apiResponse?.data?.message) {
      // If server provided a message, use it
      msgComponents.push(apiResponse.data.msg);
    } else {
      // Otherwise build our own message
      if (conditionData.value.length === 0) {
        msgComponents.push("Bạn đã điểm danh thành công");
      } else if (passcount > 0) {
        msgComponents.push(`Điểm danh thành công ${passcount}/${totalConditions} điều kiện`);
      } else {
        msgComponents.push("Điểm danh thành công");
      }
    }
    
    // 2. Check-in count info if not in the basic message
    if (!apiResponse.data?.message && checkinCount > 1 && !msgComponents[0]?.includes("lần")) {
      msgComponents.push(`Tổng số lần điểm danh: ${checkinCount}`);
    }
    
    // 3. Lucky number info if not in the basic message
    if (!apiResponse.data?.message && luckyNumber != null && !msgComponents[0]?.includes("may mắn")) {
      msgComponents.push(`<span class="text-yellow-600 font-bold"><i class="el-icon-star-on"></i> Chúc mừng! Số may mắn của bạn: ${luckyNumber}</span>`);
    }
    
    // 4. Additional info about attendance
    if (attendance.value.useLuckyNumber && !luckyNumber && checkinCount < attendance.value.luckyShowAfterMinCount) {
      const remaining = attendance.value.luckyShowAfterMinCount - checkinCount;
      msgComponents.push(`Bạn cần điểm danh thêm ${remaining} lần nữa và điểm danh thành công ít nhất một điều kiện để nhận số may mắn`);
      
      // Update the progress bar values
      checkinsRemaining.value = remaining;
      checkinsProgress.value = (checkinCount / attendance.value.luckyShowAfterMinCount) * 100;
    }
    
    // Combine all message components
    const msg = msgComponents.join("<br>");

    // Show message with HTML support
    ElMessageBox.alert(msg, 'Thông báo', {
      confirmButtonText: 'OK',
      type: 'success',
      dangerouslyUseHTMLString: true,
      center: true
    }).then(() => {
      if (luckyNumber) {
        // Set flag to indicate this is a new lucky number
        isNewLuckyNumber.value = true
        
        // Reset the flag after a delay
        setTimeout(() => { 
          isNewLuckyNumber.value = false 
        }, 8000)
      }
      
      if (attendance.value.redirectUrl) {
        window.location.href = attendance.value.redirectUrl
      }
    });
  } else {
    ElMessage(apiResult.data?.msg ?? apiResult.msg)
  }
}


const handleScroll = () => {
  if (!header.value) return;
  const value = window.scrollY;

  if (value > 0) {
    header.value.classList.add('bg-white');
    header.value.classList.remove('bg-transparent');
  } else {
    header.value.classList.add('bg-transparent');
    header.value.classList.remove('bg-white');
  }
};

onMounted(async () => {
  await nextTick();
  window.addEventListener('scroll', handleScroll);
});

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll);
});

// watch([isLoading, error], ([loading, err]) => {
//   if (!loading && !err) {
//     ElMessage.success('Lấy vị trí thành công!');
//   } else if (!loading && err) {
//     ElMessage.error(errorMessage.value || 'Không thể lấy vị trí');
//   }
// });

</script>

<style scope>
body,
html {
  height: 100%;
  overflow-y: auto;
}

main.isolate {
  overflow-y: unset;
}

#app {
  height: 100%;
  overflow-y: auto;
}

/* Import styles for the legacy UI parts */
/* @import '../../style/lucky-animation.css'; */

/* Confetti animation styles */
.confetti {
  position: absolute;
  width: 10px;
  height: 10px;
  background-color: #ffd700;
  opacity: 0.7;
  animation: confetti-fall 3s ease-in-out infinite, confetti-shake 2s ease-in-out infinite;
}

@keyframes confetti-fall {
  0% {
    top: -10%;
    opacity: 1;
  }
  100% {
    top: 100%;
    opacity: 0;
  }
}

@keyframes confetti-shake {
  0% {
    transform: skew(0deg, 0deg) rotate(0deg);
  }
  25% {
    transform: skew(5deg, 5deg) rotate(5deg);
  }
  50% {
    transform: skew(-5deg, -5deg) rotate(10deg);
  }
  75% {
    transform: skew(5deg, 5deg) rotate(5deg);
  }
  100% {
    transform: skew(0deg, 0deg) rotate(0deg);
  }
}
</style>