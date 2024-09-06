<template>
    <div class="bg-white">
        <main class="isolate">
            <header ref="header" class="sticky top-0 bg-white z-50 mx-auto px-4 sm:px-6 lg:px-52">
                <div class="container mx-auto py-3">
                    <div class="flex items-center justify-between">
                        <a href="/" class="flex items-center">
                            <img src="/logo-hoat-dong.webp" alt="Trường Đại học Đà Lạt"
                                class="h-16 sm:h-20 lg:h-28 w-auto">
                        </a>
                        <div class="hidden md:flex items-center gap-4 lg:gap-7">

                        </div>
                        <button
                            class="hidden md:block bg-[#79a227] text-white text-base lg:text-lg py-2 px-6 lg:py-3 lg:px-10 rounded-xl outline-none border-none cursor-pointer"
                            @click="redirectToLogin">
                            Đăng nhập
                        </button>
                        <button class="md:hidden text-gray-500 hover:text-gray-700" @click="toggleMobileMenu">
                            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"
                                xmlns="http://www.w3.org/2000/svg">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M4 6h16M4 12h16M4 18h16" />
                            </svg>
                        </button>
                    </div>
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
                <div class="overflow-scroll">
                    <div class="mx-auto max-w-7xl px-6 pb-32 pt-24 sm:pt-60 lg:px-8 lg:pt-10">
                        <div class="mx-auto max-w-2xl gap-x-14 lg:mx-0 lg:flex lg:max-w-none lg:items-center">
                            <div class="w-full max-w-xl lg:shrink-0 xl:max-w-4xl">
                                <h1 class="text-4xl font-bold tracking-tight text-[#514C39] sm:text-6xl">
                                    Hệ thống điểm
                                    danh<br><span class="text-[#7BA227]">Trường Đại học Đà Lạt</span>.
                                </h1>
                            </div>

                            <div
                                class="mt-14 flex justify-end gap-8 sm:-mt-44 sm:justify-start sm:pl-20 lg:mt-0 lg:pl-0">
                                <div
                                    class="ml-auto w-44 flex-none space-y-8 pt-32 sm:ml-0 sm:pt-80 lg:order-last lg:pt-36 xl:order-none xl:pt-80">
                                    <div class="relative">
                                        <img src="https://tansinhvien.dlu.edu.vn/wp-content/uploads/2024/08/IMG_3589-scaled-1.jpg"
                                            alt=""
                                            class="aspect-[2/3] w-full rounded-xl bg-gray-900/5 object-cover shadow-lg">
                                        <div
                                            class="pointer-events-none absolute inset-0 rounded-xl ring-1 ring-inset ring-gray-900/10" />
                                    </div>
                                </div>
                                <div class="mr-auto w-44 flex-none space-y-8 sm:mr-0 sm:pt-52 lg:pt-36">
                                    <div class="relative">
                                        <img src="https://tansinhvien.dlu.edu.vn/wp-content/uploads/2024/08/IMG_4601-scaled-1.jpg"
                                            alt=""
                                            class="aspect-[2/3] w-full rounded-xl bg-gray-900/5 object-cover shadow-lg">
                                        <div
                                            class="pointer-events-none absolute inset-0 rounded-xl ring-1 ring-inset ring-gray-900/10" />
                                    </div>
                                    <div class="relative">
                                        <img src="https://tansinhvien.dlu.edu.vn/wp-content/uploads/2024/08/IMG_9950-scaled-1.jpg"
                                            alt=""
                                            class="aspect-[2/3] w-full rounded-xl bg-gray-900/5 object-cover shadow-lg">
                                        <div
                                            class="pointer-events-none absolute inset-0 rounded-xl ring-1 ring-inset ring-gray-900/10" />
                                    </div>
                                </div>
                                <div class="w-44 flex-none space-y-8 pt-32 sm:pt-0">
                                    <div class="relative">
                                        <img src="https://tansinhvien.dlu.edu.vn/wp-content/uploads/2024/08/Thung-Lung-1.jpg"
                                            alt=""
                                            class="aspect-[2/3] w-full rounded-xl bg-gray-900/5 object-cover shadow-lg">
                                        <div
                                            class="pointer-events-none absolute inset-0 rounded-xl ring-1 ring-inset ring-gray-900/10" />
                                    </div>
                                    <div class="relative">
                                        <img src="https://tansinhvien.dlu.edu.vn/wp-content/uploads/2024/08/IMG_8452-1-scaled-1.jpg"
                                            alt=""
                                            class="aspect-[2/3] w-full rounded-xl bg-gray-900/5 object-cover shadow-lg">
                                        <div
                                            class="pointer-events-none absolute inset-0 rounded-xl ring-1 ring-inset ring-gray-900/10" />
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            <section class="py-16 bg-white">
                <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                    <h2 class="text-2xl sm:text-3xl font-bold text-center mb-8 sm:mb-12 text-gray-800">
                        Các hoạt động
                    </h2>
                    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 sm:gap-8">
                        <div v-for="(activity, index) in postData" :key="index"
                            class="bg-white rounded-lg shadow-md overflow-hidden flex flex-col">
                            <div class="relative overflow-hidden">
                                <!-- <img :src="activity.image" :alt="activity.title" class="w-full h-48 object-cover"> -->
                                <canvas :ref="canvasRefs[index]" class="w-full h-48"></canvas>
                                <!-- <div class="absolute top-2 left-2 bg-[#79a227] text-white px-3 py-1 rounded">
                                    Mã QR 
                                </div> -->
                                <!-- <canvas :ref="(el) => (canvasRefs.value[index] = el)" class="w-full h-48"></canvas> -->
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
                                    <p class="p-1"><span class="font-medium">Danh mục:</span>
                                        <span class="bg-blue-100 text-blue-800 px-2 py-1 rounded">{{
                                            activity.category.name }}</span>
                                    </p>
                                    <p class="p-1"><span class="font-medium">Hệ số:</span> <span
                                            class="bg-green-100 text-green-800 px-2 py-1 rounded">{{
                                                activity.weight }}</span></p>
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
                    <div class="flex justify-center mt-8">
                        <nav class="inline-flex rounded-md shadow">
                            <a href="#"
                                class="px-3 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50">
                                Trước
                            </a>
                            <a href="#"
                                class="px-3 py-2 border-t border-b border-gray-300 bg-white text-sm font-medium text-gray-700 hover:bg-gray-50">
                                1
                            </a>
                            <a href="#"
                                class="px-3 py-2 border-t border-b border-gray-300 bg-white text-sm font-medium text-gray-700 hover:bg-gray-50">
                                2
                            </a>
                            <a href="#"
                                class="px-3 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50">
                                Sau
                            </a>
                        </nav>
                    </div>
                </div>
            </section>
        </main>
    </div>
</template>


<script setup>
import { getAttendancePublic } from '@/api/checkins/attendance';
import { onMounted, ref, watch, nextTick } from 'vue'
import moment from 'moment';
import QRCode from 'qrcode'

const postData = ref([])
const page = ref(1)
const total = ref(0)
const pageSize = ref(6)


const getPostData = async () => {
    var now = moment()
    var searchInfo = {
        startDate: now.format('YYYY-MM-DDTHH:mm:ssZ'),
        endDate: now.add(30, 'days').format('YYYY-MM-DDTHH:mm:ssZ')
    }
    const res = await getAttendancePublic({ page: page.value, pageSize: pageSize.value, startDate: searchInfo.startDate, endDate: searchInfo.endDate })
    console.log(res)
    if (res.code == 0) {
        postData.value = res.data.list
        total.value = res.data.total
        page.value = res.data.page
        pageSize.value = res.data.pageSize
        // Update canvasRefs to reference canvas elements
        // await nextTick()  // Ensure the DOM is updated
        
        postData.value.forEach((activity, index) => {
            generateQRCode(index, activity.clientUrl)
        })
    }
    console.log(postData.value)
}




const canvasRefs = ref( Array(6).fill(null) )
console.log("canvasRefs: ", canvasRefs)
const generateQRCode = async(index, url) => {
    const canvas = canvasRefs.value[index]
    if (canvas) {
        try {
            await QRCode.toCanvas(canvas, url, { width: 300 })
        } catch (error) {
            console.error("QR code canvas error: ", error)
        }
    }
}


watch(postData, (newPostData) => {
    newPostData.forEach((activity, index) => {
        generateQRCode(index, activity.clientUrl)
    })
})

onMounted(() => {
    getPostData()
})

const redirectToLogin = () => {
    window.location.href = '/login'
}

const calProgress = (totalCheckin, total) => {
    total = total == 0 ? 1 : total
    return `${(totalCheckin / total * 100).toFixed(2)}%`;
}

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