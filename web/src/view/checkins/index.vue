<template>
  <div class="bg-white relative">
    <main class="isolate">
      <!-- Header -->
      <header ref="header" class="sticky top-0 z-50">
        <div class="bg-white mx-auto px-4 sm:px-6 lg:px-52">
          <div class="container mx-auto py-3">
            <div class="flex items-center justify-between">
              <a href="/" class="flex items-center">
                <img src="/logo-dlu-full.webp" alt="Trường Đại học Đà Lạt" class="h-14 sm:h-16 lg:h-20 w-auto">
              </a>
              <div class="hidden md:flex items-center gap-4 lg:gap-7"></div>
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

      <!-- Hero Section -->
      <div class="relative isolate -z-10">
        <svg class="absolute inset-x-0 top-0 -z-10 h-[64rem] w-full stroke-gray-200 [mask-image:radial-gradient(32rem_32rem_at_center,white,transparent)]" aria-hidden="true">
          <defs>
            <pattern id="1f932ae7-37de-4c0a-a8b0-a6e3b4d44b84" width="200" height="200" x="50%" y="-1" patternUnits="userSpaceOnUse">
              <path d="M.5 200V.5H200" fill="none" />
            </pattern>
          </defs>
          <svg x="50%" y="-1" class="overflow-visible fill-gray-50">
            <path d="M-200 0h201v201h-201Z M600 0h201v201h-201Z M-400 600h201v201h-201Z M200 800h201v201h-201Z" stroke-width="0" />
          </svg>
          <rect width="100%" height="100%" stroke-width="0" fill="url(#1f932ae7-37de-4c0a-a8b0-a6e3b4d44b84)" />
        </svg>
        <div class="absolute left-1/2 right-0 top-0 -z-10 -ml-24 transform-gpu overflow-hidden blur-3xl lg:ml-24 xl:ml-48" aria-hidden="true">
          <div class="aspect-[801/1036] w-[50.0625rem] bg-gradient-to-tr from-[#ff80b5] to-[#9089fc] opacity-30"
            style="clip-path: polygon(63.1% 29.5%, 100% 17.1%, 76.6% 3%, 48.4% 0%, 44.6% 4.7%, 54.5% 25.3%, 59.8% 49%, 55.2% 57.8%, 44.4% 57.2%, 27.8% 47.9%, 35.1% 81.5%, 0% 97.7%, 39.2% 100%, 35.2% 81.4%, 97.2% 52.8%, 63.1% 29.5%)"></div>
        </div>
        <div class="overflow-hidden">
          <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 pb-[32px] pt-[32px]">
            <div class="mx-auto max-w-2xl gap-x-14 lg:mx-0 lg:flex lg:max-w-none lg:items-center">
              <div class="w-full max-w-xl lg:shrink-0 xl:max-w-2xl">
                <h1 class="text-3xl sm:text-4xl lg:text-6xl font-bold tracking-tight text-[#514C39]">
                  Hệ thống điểm danh<br><span class="text-[#7BA227]">Trường Đại học Đà Lạt</span>.
                </h1>
                <p class="relative mt-6 text-base sm:text-lg font-bold leading-8 text-[#E67F32]">
                  Bạn hãy đăng nhập bằng tài khoản Email của bạn với Google để đăng nhập.
                </p>
                <p class="mt-6 text-lg leading-8 text-gray-600 px-2 mb-6">
                  Bạn hãy sử dụng 1 trong 2 WiFi: DLU Student hoặc DLU Teacher để thực hiện điểm danh.
                </p>
                <p class="relative mt-1 text-sm sm:text-base italic leading-8 text-gray-500">
                  {{ isSupported ? "Trình duyệt hỗ trợ lấy vị trí: " : "Trình duyệt không hỗ trợ lấy vị trí" }}
                  {{ coords.latitude + ", " + coords.longitude }}
                </p>
                <div v-if="attendance.title != null" class="mt-4 p-2 rounded shadow-slate-400">
                  <h3>{{ attendance.title }}</h3>
                  <div v-if="conditionData.length > 0">
                    <div class="text-base text-gray-900">Danh sách điều kiện điểm danh</div>
                    <dl class="divide-y divide-gray-100 mt-2 mb-2">
                      <div v-for="(item, key) in conditionData" :key="key"
                        class="px-2 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-2 bg-slate-50 mt-2">
                        <dt class="text-sm font-medium text-base text-gray-900">
                          Lần thứ {{ key + 1 }}
                        </dt>
                        <dd class="my-1 text-sm text-base text-gray-700 sm:col-span-2 sm:mt-0">
                          {{ conditionString(item) }}
                        </dd>
                        <el-tag v-if="item.isPass" effect="dark" type="success">Bạn đã điểm danh</el-tag>
                        <el-tag v-else effect="dark" type="danger">Bạn chưa điểm danh</el-tag>
                        <el-tag v-if="!item.isPass" effect="dark" type="primary">
                          <strong>{{ item.msg }}</strong>
                        </el-tag>
                      </div>
                    </dl>
                  </div>
                </div>

                <!-- Google Login -->
                <GoogleLogin class="my-4" :callback="callback" :error="gError" prompt />
              
              </div>
              <div class="mt-14 flex justify-end gap-8 sm:-mt-44 sm:justify-start sm:pl-20 lg:mt-0 lg:pl-0">
                <!-- Phần ảnh trang trí -->
                <div class="ml-auto w-44 flex-none space-y-8 pt-32 sm:ml-0 sm:pt-80 lg:order-last lg:pt-36 xl:order-none xl:pt-80">
                  <div class="relative">
                    <img src="/dlu1.jpg" alt=""
                      class="aspect-[2/3] w-full rounded-xl bg-gray-900/5 object-cover shadow-lg">
                    <div class="pointer-events-none absolute inset-0 rounded-xl ring-1 ring-inset ring-gray-900/10"></div>
                  </div>
                </div>
                <div class="mr-auto w-44 flex-none space-y-8 sm:mr-0 sm:pt-52 lg:pt-36">
                  <div class="relative">
                    <img src="/dlu2.jpg" alt=""
                      class="aspect-[2/3] w-full rounded-xl bg-gray-900/5 object-cover shadow-lg">
                    <div class="pointer-events-none absolute inset-0 rounded-xl ring-1 ring-inset ring-gray-900/10"></div>
                  </div>
                  <div class="relative">
                    <img src="/dlu3.jpg" alt=""
                      class="aspect-[2/3] w-full rounded-xl bg-gray-900/5 object-cover shadow-lg">
                    <div class="pointer-events-none absolute inset-0 rounded-xl ring-1 ring-inset ring-gray-900/10"></div>
                  </div>
                </div>
                <div class="w-44 flex-none space-y-8 pt-32 sm:pt-0">
                  <div class="relative">
                    <img src="/dlu4.jpg" alt=""
                      class="aspect-[2/3] w-full rounded-xl bg-gray-900/5 object-cover shadow-lg">
                    <div class="pointer-events-none absolute inset-0 rounded-xl ring-1 ring-inset ring-gray-900/10"></div>
                  </div>
                  <div class="relative">
                    <img src="/dlu5.jpg" alt=""
                      class="aspect-[2/3] w-full rounded-xl bg-gray-900/5 object-cover shadow-lg">
                    <div class="pointer-events-none absolute inset-0 rounded-xl ring-1 ring-inset ring-gray-900/10"></div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- Modal Camera -->
    <el-dialog
      v-model="cameraOpen"
      title="Chụp ảnh trước khi điểm danh"
      width="60%"
      @opened="startCamera"
      @closed="stopCamera"
    >
      <div v-if="!capturedImage">
        <video ref="videoRef" autoplay playsinline style="width: 100%; max-height: 300px; border: 1px solid #ccc;"></video>
        <el-button type="primary" @click="capturePhoto" style="margin-top: 10px;">Chụp ảnh</el-button>
      </div>
      <div v-else>
        <p>Ảnh đã chụp:</p>
        <img :src="capturedImage" alt="Ảnh chụp" style="max-width: 100%; border: 1px solid #ccc;" />
        <div style="margin-top: 10px;">
          <el-button type="success" @click="confirmPhoto">Xác nhận</el-button>
          <el-button type="warning" @click="() => { capturedImage = null; startCamera() }" style="margin-left: 10px;">Chụp lại</el-button>
        </div>
      </div>
    </el-dialog>
    <!-- Modal điểm danh -->
    <transition name="fade">
  <div v-if="showAttendanceModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-60">
    <div class="bg-white p-8 rounded-xl shadow-2xl w-[90%] max-w-md">
      <button
        @click="cancelAttendance"
        class="absolute top-4 right-4 text-gray-500 hover:text-gray-700 text-2xl font-bold"
      >
        &times;
      </button>
      <h2 class="text-2xl font-semibold mb-6 text-center text-gray-800">Điểm danh bằng câu hỏi</h2>
      <!-- <div class="font-medium text-gray-800 mb-1">{{ attendance.question }}</div> -->
      <p class="text-gray-700 mb-4 text-lg"> {{ question }}</p>
      <input
        v-model="answer"
        class="w-full p-3 border border-gray-300 rounded-lg mb-5 focus:outline-none focus:ring-2 focus:ring-green-400"
        placeholder="Nhập câu trả lời của bạn..."
      />
      <div class="flex justify-between mt-4">
        <button 
          @click="submitAttendance" 
          class="px-6 py-3 bg-green-600 text-white rounded-lg hover:bg-green-600 transition duration-300 ease-in-out">
          Gửi
        </button>
        <button 
          @click="cancelAttendance" 
          class="px-6 py-3 bg-gray-300 text-gray-800 rounded-lg hover:bg-gray-400 transition duration-300 ease-in-out">
          Huỷ
        </button>
      </div>
      <p v-if="attendanceStatus" class="mt-5 font-semibold text-center text-lg"
        :class="{ 'text-green-600': attendanceStatus === 'Điểm danh thành công!', 'text-red-600': attendanceStatus !== 'Điểm danh thành công!' }">
        {{ attendanceStatus }}
      </p>
    </div>
  </div>
</transition>
  </div>
</template> 

<script setup>
import { nextTick, onMounted, ref, onUnmounted } from 'vue'
import axios from 'axios'
import { publicAttendanceCheckIn,
  createAttendanceCheckIn,
  getAttendanceCheckInLogList,
  submitAttendanceAnswer
} from '@/api/checkins/attendanceCheckIn'
import {
  updateAttendance,
  findAttendance,
  createAttendanceQuestion
} from '@/api/checkins/attendance'
import { ElMessage, ElMessageBox, ElNotification } from 'element-plus'
import { useGeolocation } from '@vueuse/core'
import { decodeCredential } from 'vue3-google-login'
import { useRoute } from 'vue-router'
import { formatDateTime } from '@/utils/format'
import FingerprintJS from '@fingerprintjs/fingerprintjs'
defineOptions({ name: "Checkins" })


// Câu hỏi và đáp án có thể sinh động hoặc lấy từ API
const question = 'Trường bạn đang học là trường nào?'
//const question = ref('')
const showAttendanceModal = ref(false);
const answer = ref('')
const attendanceStatus = ref('')
const props = defineProps({
  acId: Number
})
const attendanceId = ref(props.acId)
const currentUserId = ref(1) // Giả định bạn có userId (nên lấy từ Vuex/pinia/session)
const route = useRoute()
const { coords } = useGeolocation()
const currentUser = ref(null)
const data = ref({
  email: null,
  code: null,
  lat: null,
  lng: null,
  visitorId: null,
  fullName: null,
})
const submitAttendance = async () => {
  if (!answer.value) {
    ElMessage.warning('Vui lòng nhập câu trả lời!')
    return
  }

  try {
    const payload = {
      attendanceId: attendance?.id ?? 0,        // chữ thường 
      partpaticipantId: currentUserId.value,   // chữ thường
      answer: answer.value,
      checkinDate: new Date().toISOString()    // chữ thường
    }

    console.log('payload gửi lên:', payload) // kiểm tra thử trước khi gửi

    const res = await submitAttendanceAnswer(payload)
    if (res.code === 0) {
      ElMessage.success('Gửi câu trả lời thành công!')
      //publicAttendanceCheckIn({ attendanceId: props.acId }) // Cập nhật log
    } else {
      ElMessage.error(res.msg || 'Lỗi khi gửi câu trả lời')
    }
  } catch (err) {
    console.error(err)
    ElMessage.error('Gửi thất bại. Vui lòng thử lại.')
  }
}
const loadAttendance = async () => {
  try {
    const res = await findAttendance({ id: attendanceId }) // giả sử đây là API đã import từ chỗ khác
    // xử lý dữ liệu ở đây
  } catch (error) {
    console.error(error)
  }
}
// onMounted(async () => {
//   const response = await findAttendance()
//   question.value = response.data.question
// })



const cancelAttendance = () => {
  showAttendanceModal.value = false;
  attendanceStatus.value = '';
  answer.value = '';
};

const formData = ref({
  requirePhoto: ref(true),

})
// Quản lý camera
const cameraOpen = ref(false)
const capturedImage = ref(null)
const videoRef = ref(null)
const videoStream = ref(null)
let canvas = null

// Google login callback
// const callback = async (response) => {
//   const userData = decodeCredential(response.credential)
//   data.value.email = userData.email
//   data.value.fullName = (userData?.given_name || "") + ' ' + (userData?.family_name || "")

//   if (formData.value.requirePhoto) {
//     // Mở giao diện chụp ảnh
//     cameraOpen.value = true
//   } else {
//     // Không cần ảnh thì hiện form câu hỏi luôn
//     await openQuestionForm()
//   }
// }
// const openQuestionForm = async () => {
//   showAttendanceModal.value = true
//   await getQuestionAndCheckAttendance()
// }

const callback = async (response) => {
  // Giải mã thông tin từ Google
  const userData = decodeCredential(response.credential)
  data.value.email = userData.email
  data.value.fullName = (userData?.given_name || '') + ' ' + (userData?.family_name || '')
  // Mở bảng câu hỏi
  showAttendanceModal.value = true
  // Gọi API để lấy câu hỏi
  //await getQuestionAndCheckAttendance()
}
const requestCheckin = async () => {
  console.log("Requesting checkin...");
  if (route.query?.c == null) {
    ElNotification({
      title: 'Thông báo',
      message: 'Không có điểm danh nào đang hiện hành',
      type: 'warning'
    });
    console.warn("No checkin code found in route query.");
    return;
  }

  if (!coords.value.latitude || !coords.value.longitude ||
      coords.value.latitude === 0 || coords.value.longitude === 0) {
    ElMessage.error("Không thể lấy vị trí của bạn");
    console.error("Invalid coordinates:", coords.value);
    return;
  }

  // Gán dữ liệu vị trí và mã điểm danh
  data.value.lat = coords.value.latitude;
  data.value.lng = coords.value.longitude;
  data.value.accuracy = coords.value.accuracy;
  data.value.code = route.query.c;
  console.log("Data for checkin:", data.value);

  const encodedData = encodeVal(data.value);
  console.log("Sending encoded data to API...");

  try {
    const res = await publicAttendanceCheckIn({ data: encodedData });
    console.log("API response:", res);

    if (res.code === 7 && res.data.attendance?.requirePhoto) {
      attendance.value = res.data.attendance;
      showCamera.value = true;
      return;
    }

    if (res.code === 0) {
      attendance.value = res.data.attendance;
      conditionData.value = res.data.conditions?.filter((condition, index, self) =>
        index === self.findIndex((c) => c.ID === condition.ID)
      ) || [];

      // ❗Nếu yêu cầu chụp ảnh
      if (attendance.value.requirePhoto) {
        showCamera.value = true;
        return;
      }

      // ❗Nếu yêu cầu trả lời câu hỏi (điểm danh bằng câu hỏi)
      if (attendance.value.requireQuestion) {
        question.value = res.data.questions || [];
        showQuestionModal.value = true;
        return;
      }

      // ✅ Nếu không có điều kiện đặc biệt: xử lý thông báo điểm danh
      let msg = "Bạn điểm danh không thành công. Vui lòng thao tác lại.";
      const passcount = conditionData.value.reduce((count, item) => item.isPass ? count + 1 : count, 0);
      if (passcount > 0) {
        msg = `Điểm danh thành công ${passcount}/${conditionData.value.length} lần`;
      }
      if (conditionData.value.length === 0) {
        msg = "Bạn đã điểm danh thành công";
      }

      console.log("Checkin message to display:", msg);
      await ElMessageBox.alert(msg, 'Thông báo', {
        confirmButtonText: 'OK',
        type: 'success'
      }).then(() => {
        if (attendance.value.redirectUrl) {
          window.location.href = attendance.value.redirectUrl;
        }
      });
    } else {
      ElMessage(res.data?.msg ?? res.msg);
      ElMessage.error(res.data?.msg || res.msg || "Có lỗi xảy ra");
      console.error("API responded with error:", res);
    }
  } catch (error) {
    ElMessage.error("Lỗi điểm danh: " + error);
    console.error("Exception during checkin:", error);
  }
};

// Lỗi Google login
const gError = (error) => {
  console.error("Google login error:", error)
}

// Lấy giá trị `requirePhoto` từ API khi component mount
onMounted(async () => {
  try {
    const response = await axios.get('/api/attendance')
    console.log("Dữ liệu điểm danh:", response.data)  // Kiểm tra dữ liệu ở đây
    formData.value = response.data
  } catch (error) {
    console.error("Lỗi lấy dữ liệu điểm danh:", error)
  }
})

// Hàm mở camera
const startCamera = () => {
  navigator.mediaDevices.getUserMedia({ video: true })
    .then(stream => {
      if (videoRef.value) {
        videoRef.value.srcObject = stream
      }
    })
    .catch(err => {
      console.error("Lỗi truy cập camera:", err)
    })
}

// Hàm dừng camera
const stopCamera = () => {
  if (videoRef.value && videoRef.value.srcObject) {
    videoRef.value.srcObject.getTracks().forEach(track => track.stop())
  }
}

// Hàm chụp ảnh
const capturePhoto = () => {
  if (!canvas) {
    canvas = document.createElement('canvas')
  }
  const video = videoRef.value
  if (video) {
    canvas.width = video.videoWidth
    canvas.height = video.videoHeight
    const ctx = canvas.getContext('2d')
    ctx.drawImage(video, 0, 0, canvas.width, canvas.height)
    capturedImage.value = canvas.toDataURL('image/jpeg')
    stopCamera()
  }
}

// Khi người dùng nhấn "Xác nhận", gọi API điểm danh và đóng modal
// const confirmPhoto = async () => {
//   if (!capturedImage.value) {
//     ElMessage.error("Vui lòng chụp ảnh trước khi xác nhận điểm danh!");
//     return;
//   }
//   try {
//     const blob = dataURLtoBlob(capturedImage.value);
//     const formDataUpload = new FormData();
//     formDataUpload.append("file", blob, "checkin.jpg");
//     const uploadRes = await axios.post('/api/upload', formDataUpload, {
//       headers: { "Content-Type": "multipart/form-data" }
//     });
//     console.log("uploadRes:", uploadRes);
//     const photoUrl = uploadRes.data.url;
    
//     const checkinRes = await axios.post('/api/checkin', { photoUrl /*, email: data.value.email, lat: coords.latitude, lng: coords.longitude, ...*/ });
//     console.log("checkinRes:", checkinRes);
    
//     if (checkinRes.data && checkinRes.data.success) {
//       cameraOpen.value = false;
//       ElMessage.success("Bạn đã điểm danh thành công!");
//     } else {
//       ElMessage.error("Điểm danh thất bại!");
//     }
//   } catch (err) {
//     console.error("Lỗi điểm danh:", err);
//     ElMessage.error("Điểm danh thất bại!");
//   }
// }
const confirmPhoto = async () => {
  if (!capturedImage.value) {
    ElMessage.error("Vui lòng chụp ảnh trước khi xác nhận điểm danh!")
    return
  }
  cameraOpen.value = false
  ElMessage.success("Bạn đã điểm danh thành công!")
}

// Nếu muốn chụp lại
const retakePhoto = () => {
  capturedImage.value = null
  startCamera()
}

// Hàm chuyển đổi dataURL sang Blob
function dataURLtoBlob(dataurl) {
  const arr = dataurl.split(',');
  const mime = arr[0].match(/:(.*?);/)[1];
  const bstr = atob(arr[1]);
  let n = bstr.length;
  const u8arr = new Uint8Array(n);
  while(n--) {
    u8arr[n] = bstr.charCodeAt(n);
  }
  return new Blob([u8arr], { type: mime });
}

// Hàm gọi API điểm danh, gửi kèm URL ảnh nếu có
async function requestCheckin(photoUrl) {
  return axios.post('/api/checkin', { photoUrl });
}

// Hàm chuyển đổi data thành chuỗi nhị phân và mã hóa
function toBinaryStr(str) {
  const encoder = new TextEncoder();
  const charCodes = encoder.encode(str);
  return String.fromCharCode(...charCodes);
}
const keyRandom = 'E;>YIws8_DdsSMG£sL£@lq8E<(O?Sc5'
const encodeVal = (data) => {
  const jsonString = JSON.stringify(data);
  const encodedData = btoa(toBinaryStr(jsonString));
  return keyRandom + "_" + encodedData;
}

// Hàm lấy fingerprint
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
fingerPrint();
const visitorTemplate = (val) => {
  return "dlu_activities_20422_5BS:W`A8nF<J6Y{V4Nv.r!Je_" + val;
}

// Các biến và hàm khác của phần điểm danh, QR code, tìm kiếm, etc.
const conditionData = ref([]);
const attendance = ref({});
const isSupported = ref(true); // Giả sử browser hỗ trợ

onMounted(async () => {
  await nextTick();
  window.addEventListener('scroll', handleScroll);
});

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll);
});

function handleScroll() {
  if (!header.value) return;
  const value = window.scrollY;
  if (value > 0) {
    header.value.classList.add('bg-white');
    header.value.classList.remove('bg-transparent');
  } else {
    header.value.classList.add('bg-transparent');
    header.value.classList.remove('bg-white');
  }
}
</script>

<style scoped>
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
</style>
