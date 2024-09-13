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
                <button
                  class="hidden md:block bg-[#79a227] text-white text-base lg:text-base py-2 px-6 lg:py-3 lg:px-8 rounded-xl outline-none border-none cursor-pointer"
                  @click="redirectToLogin">
                  Đăng nhập
                </button>
                <button
                  class="md:hidden bg-[#79a227] text-white text-base lg:text-base py-2 px-6 lg:py-3 lg:px-8 rounded-xl outline-none border-none cursor-pointer"
                  @click="redirectToLogin">
                  Đăng nhập
                </button>
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
                    Hãy đăng nhập bằng tài khoản Email của bạn với Google để đăng nhập.
                  </p>
                  
                  <!-- <div v-if="attendance.title != null" class="mt-4 bg-slate-50 p-2 rounded shadow-slate-400">
                    <h3>{{ attendance.title }}</h3>
                    <div v-if="conditionData.length > 0">
                      <div class="text-base text-gray-900">
                        Lịch sử điểm danh
                      </div>
                      <dl class="divide-y divide-gray-100">
                        <div v-for="(item, key) in conditionData" :key="key"
                          class="px-2 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                          <dt class="text-sm font-medium leading-6 text-gray-900">
                            Lần thứ {{ key + 1 }}
                          </dt>
                          <dd class="my-1 text-sm leading-6 text-gray-700 sm:col-span-2 sm:mt-0">
                            {{ conditionString(item) }}
                          </dd>
                          <el-tag v-if="item.isPass" effect="dark" type="success">
                            Bạn đã điểm danh
                          </el-tag>
                          <el-tag v-if="!item.isPass" effect="dart" type="danger">
                            Bạn chưa điểm danh
                          </el-tag>
                        </div>
                      </dl>
                    </div>
                  </div> -->
  
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
  import { decodeCredential } from 'vue3-google-login'
  
  import { useRoute } from 'vue-router';
  import { formatDateTime } from '@/utils/format';
  
  import FingerprintJS from '@fingerprintjs/fingerprintjs';
  
  
  const route = useRoute()
  
  
  defineOptions({
    name: "Checkins",
  })
  

  
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
    console.log("Handle the error", error)
  }
  
  
  async function getFingerprint() {
    const fpPromise = FingerprintJS.load();
    const fp = await fpPromise;
    const result = await fp.get();
    return result.visitorId;
  }
  
  const fingerPrint = () => {
    getFingerprint().then(visitorId => {
      data.value.visitorId = visitorId
    });
  }
  
  fingerPrint()
  
  
  const handleScroll = () => {
    console.log('header', header.value)
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
  </style>