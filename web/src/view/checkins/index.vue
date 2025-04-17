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
                <p class="mt-6 text-lg leading-8 text-gray-600 px-2 mb-6">
                  Bạn hãy Sử dụng 1 trong 2 WiFi: DLU Student hoặc DLU Teacher để thực hiện điểm danh
                </p>
                <p class="relative mt-1 text-sm sm:text-base italic leading-8 text-gray-500 sm:max-w-md lg:max-w-none">
                  {{ isSupported ? "Trình duyệt hỗ trợ lấy vị trí: " : "Trình duyệt không hỗ trợ lấy vị trí" }}
                  {{ coords.latitude + ", " + coords.longitude }}
                </p>
                <div v-if="attendance.title != null" class="mt-4 p-2 rounded shadow-slate-400">
                  <h3>{{ attendance.title }}</h3>
                  <div v-if="conditionData.length > 0">
                    <div class="text-base text-gray-900">
                      Danh sách điều kiện điểm danh
                    </div>
                    <dl class="divide-y divide-gray-100 mt-2 mb-2">
                      <div v-for="(item, key) in conditionData" :key="key"
                        class="px-2 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-2  bg-slate-50 mt-2">
                        <dt class="text-sm font-medium text-base text-gray-900">
                          Lần thứ {{ key + 1 }}
                        </dt>
                        <dd class="my-1 text-sm text-base text-gray-700 sm:col-span-2 sm:mt-0">
                          {{ conditionString(item) }}
                        </dd>
                        <el-tag v-if="item.isPass" effect="dark" type="success">
                          Bạn đã điểm danh
                        </el-tag>
                        <el-tag v-if="!item.isPass" effect="dart" type="danger">
                          Bạn chưa điểm danh
                        </el-tag>
                        <el-tag v-if="!item.isPass" effect="dart" type="primary">
                          <strong>{{item.msg}}</strong>
                        </el-tag>
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

    <!-- Modal Camera: hiển thị camera để chụp ảnh điểm danh -->
    <el-dialog
      v-model="showCamera"
      title="Chụp ảnh điểm danh"
      width="50%"
      :append-to-body="true"
      :lock-scroll="false"
      @opened="handleDialogOpened"
      @close="stopCamera">
      <div class="relative flex flex-col items-center">
        <div v-if="flash" class="flash-overlay"></div>
        <video ref="video" autoplay playsinline width="400" height="300" class="rounded-lg border"></video>
        <canvas ref="canvas" width="400" height="300" style="display: none;"></canvas>
        <el-button type="primary" @click="capturePhoto" class="mt-4">
          Chụp Ảnh
        </el-button>
        <el-upload
          ref="uploadComponent"
          style="display: none"
          name="file"
          v-model:fileList="fileList"
          :action="`${getBaseUrl()}/fileUploadAndDownload/upload`"
          :auto-upload="false"
          :before-upload="checkFile"
          :on-success="uploadSuccess"
          :on-error="uploadError"
          :with-credentials="true"
          :headers="authHeaders">
        </el-upload>              
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { nextTick, onMounted, ref, watch, computed  } from 'vue'
import { useRoute } from 'vue-router'
import { useGeolocation } from '@vueuse/core'
import { decodeCredential, GoogleLogin } from 'vue3-google-login'
import { ElMessage, ElMessageBox, ElNotification } from 'element-plus'
import FingerprintJS from '@fingerprintjs/fingerprintjs'
import { publicAttendanceCheckIn } from '@/api/checkins/attendanceCheckIn'
import { formatDateTime, getBaseUrl } from '@/utils/format'
import { isVideoMime, isImageMime } from '@/utils/image'

defineOptions({
  name: "Checkins",
})

const route = useRoute()
const { coords, isSupported } = useGeolocation()
const uploadComponent = ref(null)
const tempPhotoUrl = ref(null)
const emit = defineEmits(['on-success'])
const token = ref(localStorage.getItem('token') || '')
const authHeaders = computed(() => ({
  'Authorization': `Bearer ${token.value}`
}))
//console.log("token", authHeaders)
const fileList = ref([])
const fullscreenLoading = ref(false)


const checkFile = (file) => {
  fullscreenLoading.value = true
  const isLt500K = file.size / 1024 / 1024 < 1.5 // 500K, @todo should be configurable in the project
  const isLt5M = file.size / 1024 / 1024 < 5 // 5MB, @todo should be configurable in the project
  const isVideo = isVideoMime(file.type)
  const isImage = isImageMime(file.type)
  let pass = true
  if (!isVideo && !isImage) {
    ElMessage.error('Chỉ có thể tải lên hình ảnh định dạng jpg, png, svg, webp hoặc video định dạng mp4, webm!')
    fullscreenLoading.value = false
    pass = false
  }
  if (!isLt5M && isVideo) {
    ElMessage.error('Kích thước video tải lên không được vượt quá 5MB')
    fullscreenLoading.value = false
    pass = false
  }
  if (!isLt500K && isImage) {
    ElMessage.error('Kích thước hình ảnh tải lên chưa được nén không được vượt quá 500KB, vui lòng sử dụng chức năng nén trước khi tải lên')
    fullscreenLoading.value = false
    pass = false
  }
  //console.log('Kết quả kiểm tra tệp tải lên: ', pass)
  return pass
}

const uploadSuccess = async (res) => {
  //console.log(res)
  const { data: resData } = res
  //console.log({ resData })
  if (resData.file && resData.file.url) {
    //console.log(resData.file.url)
    data.value.photoURL = resData.file.url
    data.value.requirePhoto = false
  }
  fileList.value = []
  fullscreenLoading.value = false
  await requestCheckin()
  showCamera.value = false
  stopCamera()
  
  await ElMessageBox.alert('Bạn đã điểm danh thành công', 'Thông báo', {
    confirmButtonText: 'OK',
    type: 'success'
  });
}


const uploadError = () => {
  ElMessage({
    type: 'error',
    message: 'Tải lên thất bại'
  })
}

//console.log("Geolocation initialized:", { coords, isSupported });

const data = ref({
  email: null,
  fullName: null,
  code: null,
  lat: null,
  lng: null,
  accuracy: null,
  visitorId: null,
  photoURL: null,
  requirePhoto: false
})
//console.log("Initial data:", data.value);

const showCamera = ref(false)
const video = ref(null)
const canvas = ref(null)
const stream = ref(null)
const flash = ref(false)

// Các biến hiển thị thông tin kết quả và điều kiện điểm danh
const conditionData = ref([])
const attendance = ref({})

// Đồng bộ vị trí lấy từ trình duyệt
watch(coords, (newCoords) => {
  if (newCoords && newCoords.latitude && newCoords.longitude) {
    data.value.lat = newCoords.latitude;
    data.value.lng = newCoords.longitude;
    data.value.accuracy = newCoords.accuracy;
    //console.log("Updated coordinates:", newCoords);
  }
}, { immediate: true });

// Lấy Fingerprint của thiết bị và gán vào data.visitorId
async function getFingerprint() {
  const fp = await FingerprintJS.load().then(fp => fp.get());
  return fp.visitorId;
}
function visitorTemplate(val) {
  return "dlu_activities_20422_5BS:W`A8nF<J6Y{V4Nv.r!Je_" + val;
}
getFingerprint().then(visitorId => {
  data.value.visitorId = visitorTemplate(visitorId);
  //console.log("Fingerprint obtained:", data.value.visitorId);
});

// Hàm chuyển chuỗi thành dạng nhị phân và mã hóa dữ liệu
function toBinaryStr(str) {
  const encoder = new TextEncoder();
  return String.fromCharCode(...encoder.encode(str));
}
const keyRandom = 'E;>YIws8_DdsSMG£sL£@lq8E<(O?Sc5';
function encodeVal(payload) {
  const jsonString = JSON.stringify(payload);
  const encodedData = btoa(toBinaryStr(jsonString));
  const result = keyRandom + "_" + encodedData;
  //console.log("Encoded data:", result);
  return result;
}

// Hàm định dạng chuỗi điều kiện điểm danh
const conditionString = (item) => {
  const condition = {
    group: item?.group?.name,
    area: item?.area?.area?.name,
    startAt: item.startAt ? formatDateTime(item.startAt) : null,
    endAt: item.endAt ? formatDateTime(item.endAt) : null
  };
  if (condition.group && condition.area && condition.startAt && condition.endAt) {
    return `Nhóm ${condition.group}, tại ${condition.area}, từ ${condition.startAt} đến ${condition.endAt}`;
  } else if (condition.group && condition.area) {
    return `Nhóm ${condition.group}, tại ${condition.area}`;
  } else if (condition.group) {
    return `Nhóm ${condition.group}`;
  } else if (condition.area) {
    return `Tại ${condition.area}`;
  } else if (condition.startAt && condition.endAt) {
    return `Từ ${condition.startAt} đến ${condition.endAt}`;
  } else if (condition.startAt) {
    return `Từ ${condition.startAt}`;
  } else if (condition.endAt) {
    return `Hạn cuối ${condition.endAt}`;
  }
  return "";
};

// Callback khi đăng nhập bằng Google thành công
const callback = async (response) => {
  const userData = decodeCredential(response.credential);
  data.value.email = userData.email;
  data.value.fullName = (userData?.given_name || "") + " " + (userData?.family_name || "");
  //console.log("Google login successful:", data.value);

  localStorage.setItem('token', response.credential)
  token.value = response.credential

  if (data.value.requirePhoto) {
    showCamera.value = true;
  } else {
    await requestCheckin();
  }
};

// Hàm xử lý lỗi khi đăng nhập bằng Google
const gError = (error) => {
  console.error("Lỗi Google Login:", error);
};

// Khi modal được mở, yêu cầu truy cập camera
const handleDialogOpened = async () => {
  try {
    stream.value = await navigator.mediaDevices.getUserMedia({ video: true })
    if (video.value) {
      video.value.srcObject = stream.value
    }
  } catch (err) {
    ElMessage.error("Không thể truy cập camera: " + err.message)
  }
}
const stopCamera = () => {
  if (stream.value) {
    stream.value.getTracks().forEach(track => track.stop());
    stream.value = null;
  }
  showCamera.value = false;
}
const triggerFlash = () => {
  flash.value = true
  setTimeout(() => { flash.value = false }, 150)
}

const downloadFile = (file, filename) => {
  // Tạo URL tạm cho blob từ file
  const url = URL.createObjectURL(file);
  // Tạo phần tử <a> ẩn để tải file xuống
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  // Giải phóng URL sau khi tải về
  URL.revokeObjectURL(url);
};

const capturePhoto = async () => {
  //console.log("CapturePhoto triggered");
  const context = canvas.value.getContext('2d');
  context.drawImage(video.value, 0, 0, canvas.value.width, canvas.value.height);
  triggerFlash();  // Hiệu ứng flash

  // Chuyển canvas thành Blob (ảnh JPEG nén)
  canvas.value.toBlob(async (blob) => {
    if (!blob) {
      ElMessage.error('Không thể chuyển ảnh thành blob');
      return;
    }
    if (blob.size > 500 * 1024) {
      ElMessage.error('Ảnh vượt quá 500KB. Vui lòng thử lại.');
      return;
    }
    // Lấy email và id phiên điểm danh (ví dụ sử dụng data.value.code làm id phiên)
  const email = data.value.email || "unknown";
  const name = data.value.fullName || "unknown";
  const attendanceId = data.value.code || "0";
  const newFileName = `${attendanceId}-${email}-${name}.jpg`;
    // Tạo File từ Blob
    const file = new File([blob], newFileName, { type: 'image/jpeg' });
    file.uid = Date.now().toString();
    file.status = 'ready';
    file.raw = file;
    fileList.value = [file];
    console.log("File created:", file);  
    // downloadFile(file, file.name);
    // Tạo URL tạm thời để xem trước ảnh
    // tempPhotoUrl.value = URL.createObjectURL(file);
    // console.log("tempPhotoUrl.value:", tempPhotoUrl.value);
    
    // Đợi Vue cập nhật DOM (nextTick) và thêm delay nếu cần
    await nextTick();
    setTimeout(() => {
      if (uploadComponent.value && typeof uploadComponent.value.submit === 'function') {
        uploadComponent.value.submit();
        //console.log("Upload triggered via submit() method.");
      } else {
        //console.error("uploadComponent không khả dụng hoặc không có phương thức submit()");
      }
    }, 100);
    stopCamera();
  }, 'image/jpeg', 0.9);
  stopCamera();
};

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
        // hiện modal camera
        showCamera.value = true;
        return;
      }
    if (res.code === 0) {
      if (res.data.conditions != null) {
        conditionData.value = res.data.conditions.filter((condition, index, self) =>
          index === self.findIndex((c) => c.ID === condition.ID)
        );
        console.log("Filtered condition data:", conditionData.value);
      }
      attendance.value = res.data.attendance;
      conditionData.value = res.data.conditions || [];
      if (attendance.value.requirePhoto) {
        showCamera.value = true;
        return;
      }

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

// Xử lý hiệu ứng cuộn cho header
const header = ref(null);
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
  await nextTick()
  window.addEventListener('scroll', handleScroll)
})
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

.flash-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 400px;
  height: 300px;
  background: white;
  opacity: 0.8;
  z-index: 10;
  pointer-events: none;
  animation: flashFade 0.15s ease-out forwards;
}
@keyframes flashFade {
  from { opacity: 0.8; }
  to { opacity: 0; }
}
</style>
