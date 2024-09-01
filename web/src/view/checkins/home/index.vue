<template>
    <div class="bg-white">
        <main class="isolate">
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
                    <div class="mx-auto max-w-7xl px-6 pb-32 pt-24 sm:pt-60 lg:px-8 lg:pt-32">
                        <div class="mx-auto max-w-2xl gap-x-14 lg:mx-0 lg:flex lg:max-w-none lg:items-center">
                            <div class="w-full max-w-xl lg:shrink-0 xl:max-w-2xl">
                                <h1 class="text-4xl font-bold tracking-tight text-gray-900 sm:text-6xl">
                                    Hệ thống điểm
                                    danh - Trường Đại học Đà Lạt.
                                </h1>
                                <h2>Danh sách điểm danh sắp tới</h2>
                                <div
                                    class="mx-auto mt-2 grid max-w-2xl grid-cols-1 gap-x-8 gap-y-16 border-t border-gray-200 pt-10 sm:mt-16 sm:pt-16 lg:mx-0 lg:max-w-none lg:grid-cols-3">
                                    <article v-for="post in postData" :key="post.ID"
                                        class="flex max-w-xl flex-col items-start justify-between">
                                        <div class="flex flex-col items-start gap-x-4 text-base">

                                            <a
                                                class="relative z-10 rounded-md py-2 my-1 bg-gray-50 font-medium text-gray-600 hover:bg-gray-100">{{
                                                    post.category.name }}</a>
                                            <a
                                                class="relative z-10 rounded-md py-2 my-1 bg-gray-50 font-medium text-gray-600 hover:bg-gray-100">{{
                                                    post.agency.name }}</a>
                                        </div>
                                        <div class="group relative">
                                            <h3
                                                class="mt-3 text-lg font-semibold leading-6 text-gray-900 group-hover:text-gray-600">
                                                <a :href="post.href">
                                                    <span class="absolute inset-0" />
                                                    {{ post.title }}
                                                </a>
                                            </h3>
                                            <p class="mt-5 line-clamp-3 text-sm leading-6 text-gray-600">{{
                                                post.description
                                                }}</p>
                                        </div>
                                        <div class="flex flex-col items-start gap-x-4 text-base">

                                            <span class="text-gray-500 py-2">Bắt đầu: {{ formatDate(post.startDate)
                                                }}</span>
                                            <span class="text-red-500 py-2">Kết thúc: {{ formatDate(post.endDate)
                                                }}</span>
                                        </div>
                                    </article>
                                </div>
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
        </main>
    </div>
</template>


<script setup>
import { formatDate } from '@/utils/format';
import { getAttendancePublic } from '@/api/checkins/attendance';
import { ref } from 'vue'
import moment from 'moment';


const postData = ref([])

const getPostData = async () => {
    var now = moment()
    // "2006-01-02T15:04:05Z07:00"
    // var searchInfo = {
    //     startDate: now.format('DD/MM/YYYY HH:mm:ss'),
    //     endDate: now.add(10, 'days').format('DD/MM/YYYY HH:mm:ss')
    // }
    var searchInfo = {
        startDate: now.format('YYYY-MM-DDTHH:mm:ssZ'),
        endDate: now.add(10, 'days').format('YYYY-MM-DDTHH:mm:ssZ')
    }
    console.log(searchInfo)
    const res = await getAttendancePublic({ page: 0, pageSize: -1, startDate: searchInfo.startDate, endDate: searchInfo.endDate })
    console.log(res)
    if (res.code == 0) {
        postData.value = res.data.list
    }
}

getPostData()


const posts = [
    {
        id: 1,
        title: 'Boost your conversion rate',
        href: '#',
        description:
            'Illo sint voluptas. Error voluptates culpa eligendi. Hic vel totam vitae illo. Non aliquid explicabo necessitatibus unde. Sed exercitationem placeat consectetur nulla deserunt vel. Iusto corrupti dicta.',
        startAt: 'Mar 16, 2020',
        endAt: '2020-03-16',
        agency: { name: 'Marketing', },
        condition: []
    },

]
</script>