<template>
    <div class="bg-white">
        <main class="isolate bg-gradient-to-b from-white to-gray-100">
            <header ref="header" class="sticky top-0 z-50">
                <div class="bg-white  mx-auto px-4 sm:px-6 lg:px-52">
                    <div class="container mx-auto py-3">
                        <div class="flex items-center justify-between">
                            <a href="/" class="flex items-center">
                                <img src="/logo-dlu-full.webp" alt="Trường Đại học Đà Lạt"
                                    class="h-14 sm:h-16 lg:h-20 w-auto">
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
                <svg class="absolute inset-x-0 top-0 -z-10 h-[64rem] w-full stroke-gray-200 [mask-image:radial-gradient(32rem_32rem_at_center,white,transparent)]"
                    aria-hidden="true">
                    <defs>
                        <pattern id="1f932ae7-37de-4c0a-a8b0-a6e3b4d44b84" width="200" height="200" x="50%" y="-1"
                            patternUnits="userSpaceOnUse">
                            <path d="M.5 200V.5H200" fill="none" />
                        </pattern>
                    </defs>
                    <svg x="50%" y="-1" class="overflow-visible fill-gray-50">
                        <path
                            d="M-200 0h201v201h-201Z M600 0h201v201h-201Z M-400 600h201v201h-201Z M200 800h201v201h-201Z"
                            stroke-width="0" />
                    </svg>
                    <rect width="100%" height="100%" stroke-width="0"
                        fill="url(#1f932ae7-37de-4c0a-a8b0-a6e3b4d44b84)" />
                </svg>
                <div class="absolute left-1/2 right-0 top-0 -z-10 -ml-24 transform-gpu overflow-hidden blur-3xl lg:ml-24 xl:ml-48"
                    aria-hidden="true">
                    <div class="aspect-[801/1036] w-[50.0625rem] bg-gradient-to-tr from-[#ff80b5] to-[#9089fc] opacity-30"
                        style="clip-path: polygon(63.1% 29.5%, 100% 17.1%, 76.6% 3%, 48.4% 0%, 44.6% 4.7%, 54.5% 25.3%, 59.8% 49%, 55.2% 57.8%, 44.4% 57.2%, 27.8% 47.9%, 35.1% 81.5%, 0% 97.7%, 39.2% 100%, 35.2% 81.4%, 97.2% 52.8%, 63.1% 29.5%)" />
                </div>
                <div class="overflow-hidden">
                    <!-- <div class="mx-auto max-w-7xl px-6 pb-[32px] pt-[32px]">
                        <div class="mx-auto max-w-2xl gap-x-14 lg:mx-0 lg:flex lg:max-w-none lg:items-center">
                            <div class="w-full max-w-xl lg:shrink-0 xl:max-w-2xl">
                                <h1 class="text-4xl font-bold tracking-tight text-[#514C39] sm:text-6xl mb-6">
                                    Hệ thống điểm danh
                                    <br><span class="text-[#7BA227]">Trường Đại học Đà Lạt</span>
                                </h1>
                                <p class="mt-6 text-lg leading-8 text-gray-600">
                                    Chào mừng bạn đến với hệ thống điểm danh hiện đại của Trường Đại học Đà Lạt.
                                    Hãy đăng nhập để bắt đầu quá trình điểm danh của bạn.
                                </p>
                            </div>

                            <div
                                class="mt-14 flex justify-end gap-8 sm:-mt-44 sm:justify-start sm:pl-20 lg:mt-0 lg:pl-0">
                                <div
                                    class="ml-auto w-44 flex-none space-y-8 pt-32 sm:ml-0 sm:pt-80 lg:order-last lg:pt-36 xl:order-none xl:pt-80">
                                    <div class="relative">
                                        <img src="/dlu1.jpg" alt=""
                                            class="aspect-[2/3] w-full rounded-xl bg-gray-900/5 object-cover shadow-lg">
                                        <div
                                            class="pointer-events-none absolute inset-0 rounded-xl ring-1 ring-inset ring-gray-900/10" />
                                    </div>
                                </div>
                                <div class="mr-auto w-44 flex-none space-y-8 sm:mr-0 sm:pt-52 lg:pt-36">
                                    <div class="relative">
                                        <img src="/dlu2.jpg" alt=""
                                            class="aspect-[2/3] w-full rounded-xl bg-gray-900/5 object-cover shadow-lg">
                                        <div
                                            class="pointer-events-none absolute inset-0 rounded-xl ring-1 ring-inset ring-gray-900/10" />
                                    </div>
                                    <div class="relative">
                                        <img src="/dlu3.jpg" alt=""
                                            class="aspect-[2/3] w-full rounded-xl bg-gray-900/5 object-cover shadow-lg">
                                        <div
                                            class="pointer-events-none absolute inset-0 rounded-xl ring-1 ring-inset ring-gray-900/10" />
                                    </div>
                                </div>
                                <div class="w-44 flex-none space-y-8 pt-32 sm:pt-0">
                                    <div class="relative">
                                        <img src="/dlu4.jpg" alt=""
                                            class="aspect-[2/3] w-full rounded-xl bg-gray-900/5 object-cover shadow-lg">
                                        <div
                                            class="pointer-events-none absolute inset-0 rounded-xl ring-1 ring-inset ring-gray-900/10" />
                                    </div>
                                    <div class="relative">
                                        <img src="/dlu5.jpg" alt=""
                                            class="aspect-[2/3] w-full rounded-xl bg-gray-900/5 object-cover shadow-lg">
                                        <div
                                            class="pointer-events-none absolute inset-0 rounded-xl ring-1 ring-inset ring-gray-900/10" />
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div> -->
                    <div class="w-full max-w-7xl mx-auto text-center mt-8">
                        <div class="transition-all duration-120 delay-350 ease-in-out px-4">
                            <div class="relative flex flex-col">
                                <template :key="index" v-for="(item, index) in slider.images">
                                    <TransitionRoot :show="active === index"
                                        enter="transition ease-in-out duration-500 delay-200 order-first"
                                        enterFrom="opacity-0 scale-105" enterTo="opacity-100 scale-100"
                                        leave="transition ease-in-out duration-300 absolute"
                                        leaveFrom="opacity-100 scale-100" leaveTo="opacity-0 scale-95">
                                        <img class="rounded-xl lg:h[440px] sm:h-auto w-[100%]" :src="getUrl(item)">
                                    </TransitionRoot>
                                </template>

                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <section class="py-8">
                <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                    <h2 class="text-2xl sm:text-3xl font-bold text-center mb-8 sm:mb-12 text-gray-800">
                        <span class="text-3xl font-bold text-gray-800">
                            <span class="text-[#7BA227] font-bold relative">
                                CÁC HOẠT ĐỘNG
                                <span class="absolute -bottom-2 left-0 w-full h-[2px] bg-[#DDD3B3]" />
                            </span>
                        </span>
                    </h2>
                    <div>

                        <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline"
                            :rules="searchRule">
                            <el-form-item label="Đơn vị" prop="agencyId">
                                <el-select class="w-[300px]" v-model="searchInfo.agencyId" placeholder="Chọn đơn vị"
                                    clearable filterable @change="onSubmit">
                                    <el-option v-for="item in agencyOptions" :key="item.ID" :label="item.name"
                                        :value="item.ID" />
                                </el-select>
                            </el-form-item>
                        </el-form>
                    </div>
                    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 sm:gap-8">
                        <div v-for="(activity, index) in postData" :key="index"
                            class="bg-white rounded-lg shadow-md overflow-hidden flex flex-col">
                            <div class="relative overflow-hidden rounded-t-lg shadow-md">
                                <!-- <div class="relative w-full h-80 bg-gradient-to-r from-blue-100 to-green-100"> -->
                                <div class="relative w-full h-96 bg-white">
                                    <div class="absolute inset-0 flex items-center justify-center">
                                        <div class="bg-transparent p-4 rounded-lg">
                                            <!-- <canvas :id="`qr-canvas-${index}`" class="w-full h-max aspect-square" /> -->
                                            <div class="w-full h-max aspect-square">
                                                <QRCodeVue3 :width="512" :height="512" :value="activity.clientUrl"
                                                    :key="index" :qr-options="{
                                                        errorCorrectionLevel: 'H'
                                                    }"
                                                    :image-options="{ hideBackgroundDots: true, imageSize: 0.4, margin: 10 }"
                                                    :corners-square-options="{ type: 'extra-rounded', color: '#514C39' }"
                                                    :corners-dot-options="{
                                                        type: undefined,
                                                        color: '#7BA227'
                                                    }" :dots-options="{
                                                        type: 'dots',
                                                        color: '#7BA227',
                                                        gradient: {
                                                            type: 'linear',
                                                            rotation: 45,
                                                            colorStops: [
                                                                { offset: 0, color: '#7BA227' },
                                                                { offset: 1, color: '#E67F32' }
                                                            ]
                                                        }
                                                    }" :backgroundOptions="backgroundQROptions" :download="false"
                                                    image="/dlu.svg" imgclass="qrcode-img" myclass="home-qrcode" />
                                            </div>
                                        </div>
                                    </div>
                                    <div
                                        class="absolute bottom-0 left-0 right-0 h-16 bg-gradient-to-t from-white opacity-70" />
                                </div>
                            </div>
                            <div
                                class="p-6 flex-grow flex flex-col justify-between bg-gradient-to-br from-white to-gray-50">
                                <div>
                                    <h3 class="text-2xl font-bold mb-3 mt-0 text-gray-800">
                                        {{ activity.title }}
                                    </h3>
                                    <p class="text-gray-600 mb-4 leading-relaxed">
                                        {{ activity.description }}
                                    </p>
                                </div>
                                <div class="text-sm text-gray-500 space-y-2">
                                    <p class="p-1">
                                        <span class="font-medium">Đơn vị:
                                        </span> <span class="text-[#E67F32]"> {{ activity.agency.name }}</span>
                                    </p>
                                    <p class="p-1">
                                        <span class="font-medium">Danh mục:</span>
                                        <span class="bg-blue-100 text-blue-800 px-2 py-1 rounded">{{
                                            activity.category.name }}</span>
                                    </p>
                                    <p class="p-1">
                                        <span class="font-medium">Hệ số:</span> <span
                                            class="bg-green-100 text-green-800 px-2 py-1 rounded">{{
                                                activity.weight }}</span>
                                    </p>
                                    <p class="p-1">
                                        <span class="font-medium">Ngày bắt đầu:</span> <span
                                            class="bg-green-100 text-green-800 px-2 py-1 rounded">{{
                                                formatDateTime(activity.startDate) }}</span>
                                    </p>
                                    <p class="p-1">
                                        <span class="font-medium">Ngày kết thúc:</span> <span
                                            class="bg-red-100 text-red-800 px-2 py-1 rounded">{{
                                                formatDateTime(activity.endDate) }}</span>
                                    </p>
                                </div>
                            </div>
                            <div class="px-6 py-4 bg-gradient-to-r from-green-50 to-blue-50">
                                <div class="flex justify-between items-center mb-2">
                                    <span class="text-gray-700 font-medium">Đã điểm danh:</span>
                                    <div class="bg-white rounded-full px-4 py-2 shadow-md">
                                        <span class="text-green-600 font-bold">{{ activity.totalCheckin }}</span>
                                        <span class="text-gray-400 mx-1">/</span>
                                        <span class="text-gray-600 font-bold">{{ activity.total }}</span>
                                    </div>
                                </div>
                                <div class="w-full bg-gray-200 rounded-full h-2">
                                    <div :style="{ width: calProgress(activity.totalCheckin, activity.total) }"
                                        class="bg-green-600 h-2 rounded-full" />
                                </div>
                            </div>
                        </div>
                    </div>
                    <div class="flex justify-center mt-8 text-[#7BA227]">
                        <el-pagination layout="prev, pager, next" :total="total" :current-page="page"
                            :page-size="pageSize" @current-change="handlePageChange" />
                    </div>
                </div>
            </section>
        </main>
    </div>
</template>


<script setup>
import { getAttendancePublic } from '@/api/checkins/attendance';
import { onMounted, ref, watch, onUnmounted } from 'vue'
import moment from 'moment';
import QRCodeVue3 from 'qrcode-vue3'
import {
    getAttendanceAgencyPublic
} from '@/api/checkins/attendanceAgency'
import { getSliderBuilderPublic } from '@/api/uibuilder/sliderBuilder'
import { getUrl } from '@/utils/image'
import { TransitionRoot } from '@headlessui/vue'
import { formatDateTime } from '@/utils/format';


const postData = ref([])
const page = ref(1)
const total = ref(0)
const pageSize = ref(6)
const backgroundQROptions = ref({
    color: '#FFFFFF00'
})


const searchInfo = ref({
    agencyId: null,
    startDate: null,
    endDate: null
})


const getPostData = async () => {
    var now = moment()
    searchInfo.value.startDate = now.format('YYYY-MM-DDTHH:mm:ssZ')
    searchInfo.value.endDate = now.add(30, 'days').format('YYYY-MM-DDTHH:mm:ssZ')
    // const res = await getAttendancePublic({ page: page.value, pageSize: pageSize.value, startDate: searchInfo.value.startDate, endDate: searchInfo.value.endDate })
    const res = await getAttendancePublic({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
    if (res.code == 0) {
        postData.value = res.data.list
        total.value = res.data.total
        page.value = res.data.page
        pageSize.value = res.data.pageSize
    }

}

const onSubmit = async () => {
    // console.log("onSubmit")
    await getPostData()
}

const agencyOptions = ref([])
const getAgencyOptions = async () => {
    // console.log("agencyOptions")
    const table = await getAttendanceAgencyPublic({ page: 0, pageSize: -1 })
    // console.log("after call await")
    if (table.code === 0) {
        agencyOptions.value = table.data.list
    }
    // console.log("agencyOptions")
    // console.log(agencyOptions.value)
}
getAgencyOptions()

const redirectToLogin = () => {
    window.location.href = '/login'
}

const calProgress = (totalCheckin, total) => {
    total = total == 0 ? 1 : total
    return `${(totalCheckin / total * 100).toFixed(2)}%`;
}

const handlePageChange = async (newPage) => {
    page.value = newPage
    await getPostData()
}

//Slider
const slider = ref({})

const getSlider = async () => {
    const res = await getSliderBuilderPublic()
    if (res.code === 0) {
        slider.value = res.data
    }
    console.log("getSliderBuilderPublic", res)
}

getSlider()

const duration = 5000
const frame = ref(0)
const firstFrameTime = ref(performance.now())
const active = ref(0)


const startAnimation = () => {
    firstFrameTime.value = performance.now()
    frame.value = requestAnimationFrame(animate)
}

const animate = (now) => {
    let timeDifference = now - firstFrameTime.value
    let timeFraction = Math.max(0, timeDifference) / duration
    if (timeFraction <= 1) {
        frame.value = requestAnimationFrame(animate)
    } else {
        timeFraction = 1
        active.value = (active.value + 1) % slider.value.images.length
    }
}

//System 
onMounted(async () => {
    startAnimation()
    await getPostData()
})


onUnmounted(() => cancelAnimationFrame(frame.value))

watch(active, () => {
    cancelAnimationFrame(frame.value)
    startAnimation()
})

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

img.qrcode-img {
    width: 320px;
    height: 320px;
}
</style>