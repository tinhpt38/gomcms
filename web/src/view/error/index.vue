<template>
  <div>
    <div class="w-full h-screen bg-gray-50 flex items-center justify-center">
      <div class="flex flex-col items-center text-2xl gap-4">
        <img class="w-1/3" src="../../assets/404.png">
        <p class="text-lg">Trang này đã bị một lực lượng bí ẩn hút đi</p>
        <p class="text-lg">Vấn đề thường gặp là vai trò hiện tại không có định tuyến hiện tại, nếu bạn muốn sử dụng định tuyến này, vui lòng phân công trong quản lý vai trò</p>
        <!-- <p>Địa chỉ dự án: <a href="https://github.com/flipped-aurora/gin-vue-admin" target="_blank" class="text-blue-600">https://github.com/flipped-aurora/gin-vue-admin</a></p> -->
        <el-button @click="toDashboard">Trở về trang chủ</el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useUserStore } from '@/pinia/modules/user'
import { useRouter } from 'vue-router'

defineOptions({
  name: 'Error'
})

const userStore = useUserStore()
const router = useRouter()
const toDashboard = () => {
  router.push({ name: userStore.userInfo.authority.defaultRouter })
}
</script>
