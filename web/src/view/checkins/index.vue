<template>
  <div class="bg-white">
    <main class="isolate">
      <!-- Hero section -->
      <div class="relative isolate -z-10">
        <svg
          class="absolute inset-x-0 top-0 -z-10 h-[64rem] w-full stroke-gray-200 [mask-image:radial-gradient(32rem_32rem_at_center,white,transparent)]"
          aria-hidden="true"
        >
          <defs>
            <pattern
              id="1f932ae7-37de-4c0a-a8b0-a6e3b4d44b84" width="200" height="200" x="50%" y="-1"
              patternUnits="userSpaceOnUse"
            >
              <path d="M.5 200V.5H200" fill="none" />
            </pattern>
          </defs>
          <svg x="50%" y="-1" class="overflow-visible fill-gray-50">
            <path
              d="M-200 0h201v201h-201Z M600 0h201v201h-201Z M-400 600h201v201h-201Z M200 800h201v201h-201Z"
              stroke-width="0"
            />
          </svg>
          <rect
            width="100%" height="100%" stroke-width="0"
            fill="url(#1f932ae7-37de-4c0a-a8b0-a6e3b4d44b84)"
          />
        </svg>
        <div
          class="absolute left-1/2 right-0 top-0 -z-10 -ml-24 transform-gpu overflow-hidden blur-3xl lg:ml-24 xl:ml-48"
          aria-hidden="true"
        >
          <div
            class="aspect-[801/1036] w-[50.0625rem] bg-gradient-to-tr from-[#ff80b5] to-[#9089fc] opacity-30"
            style="clip-path: polygon(63.1% 29.5%, 100% 17.1%, 76.6% 3%, 48.4% 0%, 44.6% 4.7%, 54.5% 25.3%, 59.8% 49%, 55.2% 57.8%, 44.4% 57.2%, 27.8% 47.9%, 35.1% 81.5%, 0% 97.7%, 39.2% 100%, 35.2% 81.4%, 97.2% 52.8%, 63.1% 29.5%)"
          />
        </div>
        <div class="overflow-hidden">
          <div class="mx-auto max-w-7xl px-6 pb-32 pt-24 sm:pt-60 lg:px-8 lg:pt-32">
            <div class="mx-auto max-w-2xl gap-x-14 lg:mx-0 lg:flex lg:max-w-none lg:items-center">
              <div class="w-full max-w-xl lg:shrink-0 xl:max-w-2xl">
                <h1 class="text-4xl font-bold tracking-tight text-gray-900 sm:text-6xl">
                  Hệ thống điểm
                  danh - Trường Đại học Đà Lạt.
                </h1>
                <p class="relative mt-6 text-lg leading-8 text-gray-600 sm:max-w-md lg:max-w-none">
                  Hãy đăng nhập bằng tài khoảng Email của bạn với Google để đăng nhập.
                </p>
                <div v-if="conditionData.length > 0" class="px-4 sm:px-0">
                  <h3 class="text-base font-semibold leading-7 text-gray-900">
                    Danh sách cần điểm danh
                  </h3>
                  <dl class="divide-y divide-gray-100">
                    <div
                      v-for="(item, key) in conditionData"
                      :key="key"
                      class="px-2 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0"
                    >
                      <dt class="text-sm font-medium leading-6 text-gray-900">
                        Lần thứ {{ key + 1
                        }}
                      </dt>
                      <dd class="my-1 text-sm leading-6 text-gray-700 sm:col-span-2 sm:mt-0">
                        {{ conditionString(item) }}
                      </dd>
                      <el-tag v-if="item.isPass" effect="dark" type="success">
                        Đã điểm
                        danh
                      </el-tag>
                      <el-tag v-if="!item.isPass" effect="dart" type="danger">
                        Chưa điểm
                        danh
                      </el-tag>
                    </div>
                  </dl>
                </div>
                <GoogleLogin class="my-4" :callback="callback" :error="gError" prompt />
              </div>
              <div
                class="mt-14 flex justify-end gap-8 sm:-mt-44 sm:justify-start sm:pl-20 lg:mt-0 lg:pl-0"
              >
                <div
                  class="ml-auto w-44 flex-none space-y-8 pt-32 sm:ml-0 sm:pt-80 lg:order-last lg:pt-36 xl:order-none xl:pt-80"
                >
                  <div class="relative">
                    <img
                      src="https://images.unsplash.com/photo-1557804506-669a67965ba0?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&h=528&q=80"
                      alt=""
                      class="aspect-[2/3] w-full rounded-xl bg-gray-900/5 object-cover shadow-lg"
                    >
                    <div
                      class="pointer-events-none absolute inset-0 rounded-xl ring-1 ring-inset ring-gray-900/10"
                    />
                  </div>
                </div>
                <div class="mr-auto w-44 flex-none space-y-8 sm:mr-0 sm:pt-52 lg:pt-36">
                  <div class="relative">
                    <img
                      src="https://images.unsplash.com/photo-1485217988980-11786ced9454?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&h=528&q=80"
                      alt=""
                      class="aspect-[2/3] w-full rounded-xl bg-gray-900/5 object-cover shadow-lg"
                    >
                    <div
                      class="pointer-events-none absolute inset-0 rounded-xl ring-1 ring-inset ring-gray-900/10"
                    />
                  </div>
                  <div class="relative">
                    <img
                      src="https://images.unsplash.com/photo-1559136555-9303baea8ebd?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&crop=focalpoint&fp-x=.4&w=396&h=528&q=80"
                      alt=""
                      class="aspect-[2/3] w-full rounded-xl bg-gray-900/5 object-cover shadow-lg"
                    >
                    <div
                      class="pointer-events-none absolute inset-0 rounded-xl ring-1 ring-inset ring-gray-900/10"
                    />
                  </div>
                </div>
                <div class="w-44 flex-none space-y-8 pt-32 sm:pt-0">
                  <div class="relative">
                    <img
                      src="https://images.unsplash.com/photo-1670272504528-790c24957dda?ixlib=rb-4.0.3&ixid=MnwxMjA3fDF8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&crop=left&w=400&h=528&q=80"
                      alt=""
                      class="aspect-[2/3] w-full rounded-xl bg-gray-900/5 object-cover shadow-lg"
                    >
                    <div
                      class="pointer-events-none absolute inset-0 rounded-xl ring-1 ring-inset ring-gray-900/10"
                    />
                  </div>
                  <div class="relative">
                    <img
                      src="https://images.unsplash.com/photo-1670272505284-8faba1c31f7d?ixlib=rb-4.0.3&ixid=MnwxMjA3fDF8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&h=528&q=80"
                      alt=""
                      class="aspect-[2/3] w-full rounded-xl bg-gray-900/5 object-cover shadow-lg"
                    >
                    <div
                      class="pointer-events-none absolute inset-0 rounded-xl ring-1 ring-inset ring-gray-900/10"
                    />
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
import { onMounted, ref } from 'vue'
import { publicAttendanceCheckIn } from '@/api/checkins/attendanceCheckIn'
import { ElMessage, ElMessageBox, ElNotification } from 'element-plus';
import { formatDate, useGeolocation } from '@vueuse/core'
import { decodeCredential } from 'vue3-google-login'

import { useRoute } from 'vue-router';
import { formatDateTime } from '@/utils/format';

import FingerprintJS from '@fingerprintjs/fingerprintjs';


const route = useRoute()


defineOptions({
    name: "Checkins",
})
const { coords, locatedAt, error, resume, pause } = useGeolocation()

const data = ref({
    email: null,
    code: null,
    lat: null,
    lng: null,
})

const callback = async (response) => {
    const userData = decodeCredential(response.credential)
    data.value.email = userData.email
    data.value.fullName = userData?.given_name + ' ' + userData?.family_name

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

const fingerPrint = () =>{
  getFingerprint().then(visitorId => {
  console.log(visitorId);
});
}

fingerPrint()



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
const requestCheckin = async () => {
    if (route.query?.c == null) {
        ElNotification("Không có điểm danh nào đang hiện hành")
        return
    }
    data.value.lat = coords.value.latitude
    data.value.lng = coords.value.longitude
    data.value.code = route.query.c

    var res = await publicAttendanceCheckIn(data.value)

    console.log('res', res)
    if (res.code == 0) {
        conditionData.value = res.data.conditions
        attendance.value = res.data.attendance
        ElMessageBox.alert('Điểm danh thành công', 'Thông báo', {
            confirmButtonText: 'Xem chi tiết',
            type: 'success'
        }).then(() => {

        });
    } else {
        ElMessage(res.msg)
    }
}


</script>

<style scope>
body, html {
    height: auto;
    overflow-y: auto;
}
main.isolate {
    overflow-y: auto;
}

#app {
    height: 100%;
    overflow-y: auto;
}

</style>