<template>
  <div id="userLayout" class="w-full h-full relative">
    <div class="rounded-lg flex items-center justify-evenly w-full h-full md:w-screen md:h-screen md:bg-white">
      <div class="md:w-3/5 w-10/12 h-full flex items-center justify-evenly">
        <div class="oblique h-[130%] w-3/5 bg-white dark:bg-slate-900 transform -rotate-12 absolute -ml-52" />
        <!-- Phân khối chéo -->
        <div class="z-[999] pt-12 pb-10 md:w-96 w-full  rounded-lg flex flex-col justify-between box-border">
          <div>
            <div class="flex items-center justify-center">
              <img class="w-24" :src="$GIN_VUE_ADMIN.appLogo" alt>
            </div>
            <div class="mb-9">
              <p class="text-center text-4xl font-bold">
                {{ $GIN_VUE_ADMIN.appName }}
              </p>
              <p class="text-center text-sm font-normal text-gray-500 mt-2.5">
                Hệ thống điểm danh sự kiện - Trường Đại học Đà Lạt
              </p>
            </div>
            <el-form
              ref="loginForm" :model="loginFormData" :rules="rules" :validate-on-rule-change="false"
              @keyup.enter="submitForm"
            >
              <el-form-item prop="username" class="mb-6">
                <el-input
                  v-model="loginFormData.username" size="large" placeholder="Nhập tên người dùng"
                  suffix-icon="user"
                />
              </el-form-item>
              <el-form-item prop="password" class="mb-6">
                <el-input
                  v-model="loginFormData.password" show-password size="large" type="password"
                  placeholder="Nhập mật khẩu"
                />
              </el-form-item>
              <el-form-item v-if="loginFormData.openCaptcha" prop="captcha" class="mb-6">
                <div class="flex w-full justify-between">
                  <el-input
                    v-model="loginFormData.captcha" placeholder="Nhập mã xác nhận" size="large"
                    class="flex-1 mr-5"
                  />
                  <div class="w-1/3 h-11 bg-[#c3d4f2] rounded">
                    <img
                      v-if="picPath" class="w-full h-full" :src="picPath" alt="Nhập mã xác nhận"
                      @click="loginVerify()"
                    >
                  </div>
                </div>
              </el-form-item>
              <el-form-item class="mb-6">
                <el-button class="shadow shadow-active h-11 w-full" type="primary" size="large" @click="submitForm">
                  Đăng
                  nhập
                </el-button>
              </el-form-item>
              <!-- <el-form-item class="mb-6">
                <el-button class="shadow shadow-active h-11 w-full" type="primary" size="large" @click="checkInit">Đi
                  tới khởi tạo</el-button>
              </el-form-item> -->
            </el-form>
          </div>
        </div>
      </div>
      <!-- <div class="hidden md:block w-1/2 h-full float-right bg-[#194bfb]">
        <img class="h-full" src="@/assets/login_right_banner.jpg" alt="banner">
      </div> -->
      <div class="hidden md:block w-1/2 h-full float-right bg-white">
        <img class="h-full object-contain	w-full" src="@/assets/checkin@3x.webp" alt="banner">
      </div>
    </div>

    <!-- <BottomInfo class="left-0 right-0 absolute bottom-3 mx-auto  w-full z-20">
      <div class="links items-center justify-center gap-2 hidden md:flex">
        <a href="https://www.gin-vue-admin.com/" target="_blank">
          <img src="@/assets/docs.png" class="w-8 h-8" alt="Tài liệu">
        </a>
        <a href="https://support.qq.com/product/371961" target="_blank">
          <img src="@/assets/kefu.png" class="w-8 h-8" alt="Hỗ trợ">
        </a>
        <a href="https://github.com/flipped-aurora/gin-vue-admin" target="_blank">
          <img src="@/assets/github.png" class="w-8 h-8" alt="github">
        </a>
        <a href="https://space.bilibili.com/322210472" target="_blank">
          <img src="@/assets/video.png" class="w-8 h-8" alt="Trang video">
        </a>
      </div>
    </BottomInfo> -->
  </div>
</template>

<script setup>
import { captcha } from '@/api/user'
import { checkDB } from '@/api/initdb'
import BottomInfo from '@/components/bottomInfo/bottomInfo.vue'
import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/pinia/modules/user'

defineOptions({
  name: "Login",
})

const router = useRouter()
// Hàm xác nhận
const checkUsername = (rule, value, callback) => {
  if (value.length < 5) {
    return callback(new Error('Vui lòng nhập tên người dùng đúng'))
  } else {
    callback()
  }
}
const checkPassword = (rule, value, callback) => {
  if (value.length < 6) {
    return callback(new Error('Vui lòng nhập mật khẩu đúng'))
  } else {
    callback()
  }
}

// Lấy mã xác nhận
const loginVerify = async () => {
  const ele = await captcha()
  rules.captcha.push({
    max: ele.data.captchaLength,
    min: ele.data.captchaLength,
    message: `Vui lòng nhập mã xác nhận ${ele.data.captchaLength} ký tự`,
    trigger: 'blur',
  })
  picPath.value = ele.data.picPath
  loginFormData.captchaId = ele.data.captchaId
  loginFormData.openCaptcha = ele.data.openCaptcha
}
loginVerify()

// Các hoạt động liên quan đến đăng nhập
const loginForm = ref(null)
const picPath = ref('')
const loginFormData = reactive({
  username: 'admin',
  password: '',
  captcha: '',
  captchaId: '',
  openCaptcha: false,
})
const rules = reactive({
  username: [{ validator: checkUsername, trigger: 'blur' }],
  password: [{ validator: checkPassword, trigger: 'blur' }],
  captcha: [
    {
      message: 'Mã xác nhận không đúng định dạng',
      trigger: 'blur',
    },
  ],
})

const userStore = useUserStore()
const login = async () => {
  return await userStore.LoginIn(loginFormData)
}
const submitForm = () => {

  loginForm.value.validate(async (v) => {
    if (!v) {
      // Không vượt qua kiểm tra tĩnh phía trước
      ElMessage({
        type: 'error',
        message: 'Vui lòng điền thông tin đăng nhập đúng',
        showClose: true,
      })
      loginVerify()
      return false
    }

    // Vượt qua kiểm tra, yêu cầu đăng nhập
    const flag = await login()

    // Đăng nhập thất bại, làm mới mã xác nhận
    if (!flag) {
      loginVerify()
      return false
    }

    // Đăng nhập thành công
    return true
  })
}

// Chuyển hướng khởi tạo
const checkInit = async () => {
  const res = await checkDB()
  if (res.code === 0) {
    if (res.data?.needInit) {
      userStore.NeedInit()
      router.push({ name: 'Init' })
    } else {
      ElMessage({
        type: 'info',
        message: 'Đã cấu hình thông tin cơ sở dữ liệu, không thể khởi tạo',
      })
    }
  }
}

</script>
