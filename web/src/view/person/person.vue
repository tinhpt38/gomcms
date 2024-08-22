<template>
  <div class="gva-form-box">
    <div class="grid grid-cols-12 w-full gap-2">
      <div class="col-span-3 h-full">
        <div class="w-full h-full bg-white dark:bg-slate-900 px-4 py-8 rounded-lg shadow-lg box-border">
          <div class="user-card px-6 text-center bg-white dark:bg-slate-900 shrink-0">
            <div class="flex justify-center">
              <SelectImage v-model="userStore.userInfo.headerImg" file-type="image" />
            </div>
            <div class="py-6 text-center">
              <p v-if="!editFlag" class="text-3xl flex justify-center items-center gap-4">
                {{ userStore.userInfo.nickName }}
                <el-icon class="cursor-pointer text-sm" color="#66b1ff" @click="openEdit">
                  <edit />
                </el-icon>
              </p>
              <p v-if="editFlag" class="flex justify-center items-center gap-4">
                <el-input v-model="nickName" />
                <el-icon class="cursor-pointer" color="#67c23a" @click="enterEdit">
                  <check />
                </el-icon>
                <el-icon class="cursor-pointer" color="#f23c3c" @click="closeEdit">
                  <close />
                </el-icon>
              </p>
              <p class="text-gray-500 mt-2 text-md">Người này lười biếng, không để lại gì cả</p>
            </div>
            <div class="w-full h-full text-left">
              <ul class="inline-block h-full w-full">
                <li class="info-list">
                  <el-icon>
                    <user />
                  </el-icon>
                  {{ userStore.userInfo.nickName }}
                </li>
                <el-tooltip class="item" effect="light" content="Công ty Công nghệ Phản chuyển Cực quang Bắc Kinh - Bộ phận Kỹ thuật - Nhóm Kỹ thuật Front-end" placement="top">
                  <li class="info-list">
                    <el-icon>
                      <data-analysis />
                    </el-icon>
                    Công ty Công nghệ Phản chuyển Cực quang Bắc Kinh - Bộ phận Kỹ thuật - Nhóm Kỹ thuật Front-end
                  </li>
                </el-tooltip>
                <li class="info-list">
                  <el-icon>
                    <video-camera />
                  </el-icon>
                  Trung Quốc · Thành phố Bắc Kinh · Quận Triều Tiên
                </li>
                <el-tooltip class="item" effect="light" content="GoLang/JavaScript/Vue/Gorm" placement="top">
                  <li class="info-list">
                    <el-icon>
                      <medal />
                    </el-icon>
                    GoLang/JavaScript/Vue/Gorm
                  </li>
                </el-tooltip>
              </ul>
            </div>
          </div>
        </div>
      </div>
      <div class="col-span-9 ">
        <div class="bg-white dark:bg-slate-900 h-full px-4 py-8 rounded-lg shadow-lg box-border">
          <el-tabs v-model="activeName" @tab-click="handleClick">
            <el-tab-pane label="Liên kết tài khoản" name="second">
              <ul>
                <li class="borderd">
                  <p class="pb-2.5 text-xl text-gray-600">Điện thoại bảo mật</p>
                  <p class="pb-2.5 text-lg text-gray-400">
                    Đã liên kết điện thoại: {{ userStore.userInfo.phone }}
                    <a href="javascript:void(0)" class="float-right text-blue-400"
                      @click="changePhoneFlag = true">Thay đổi ngay</a>
                  </p>
                </li>
                <li class="borderd pt-2.5">
                  <p class="pb-2.5 text-xl text-gray-600">Email bảo mật</p>
                  <p class="pb-2.5 text-lg text-gray-400">
                    Đã liên kết email: {{ userStore.userInfo.email }}
                    <a href="javascript:void(0)" class="float-right text-blue-400"
                      @click="changeEmailFlag = true">Thay đổi ngay</a>
                  </p>
                </li>
                <li class="borderd pt-2.5">
                  <p class="pb-2.5 text-xl text-gray-600">Câu hỏi bảo mật</p>
                  <p class="pb-2.5 text-lg text-gray-400">
                    Chưa thiết lập câu hỏi bảo mật
                    <a href="javascript:void(0)" class="float-right text-blue-400">Thiết lập ngay</a>
                  </p>
                </li>
                <li class="borderd pt-2.5">
                  <p class="pb-2.5 text-xl text-gray-600">Thay đổi mật khẩu</p>
                  <p class="pb-2.5 text-lg text-gray-400">
                    Thay đổi mật khẩu cá nhân
                    <a href="javascript:void(0)" class="float-right text-blue-400" @click="showPassword = true">Thay đổi mật khẩu</a>
                  </p>
                </li>
              </ul>
            </el-tab-pane>
          </el-tabs>
        </div>
      </div>
    </div>

    <el-dialog v-model="showPassword" title="Thay đổi mật khẩu" width="360px" @close="clearPassword">
      <el-form ref="modifyPwdForm" :model="pwdModify" :rules="rules" label-width="80px">
        <el-form-item :minlength="6" label="Mật khẩu cũ" prop="password">
          <el-input v-model="pwdModify.password" show-password />
        </el-form-item>
        <el-form-item :minlength="6" label="Mật khẩu mới" prop="newPassword">
          <el-input v-model="pwdModify.newPassword" show-password />
        </el-form-item>
        <el-form-item :minlength="6" label="Xác nhận mật khẩu" prop="confirmPassword">
          <el-input v-model="pwdModify.confirmPassword" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showPassword = false">Hủy</el-button>
          <el-button type="primary" @click="savePassword">Đồng ý</el-button>
        </div>
      </template>
    </el-dialog>

    <el-dialog v-model="changePhoneFlag" title="Liên kết điện thoại" width="600px">
      <el-form :model="phoneForm">
        <el-form-item label="Số điện thoại" label-width="120px">
          <el-input v-model="phoneForm.phone" placeholder="Vui lòng nhập số điện thoại" autocomplete="off" />
        </el-form-item>
        <el-form-item label="Mã xác nhận" label-width="120px">
          <div class="flex w-full gap-4">
            <el-input v-model="phoneForm.code" class="flex-1" autocomplete="off" placeholder="Vui lòng tự thiết kế dịch vụ tin nhắn, ở đây là giả lập viết bừa"
              style="width:300px" />
            <el-button type="primary" :disabled="time > 0" @click="getCode">{{ time > 0 ? `(${time}s) Sau khi lấy lại` : 'Lấy mã xác nhận'
              }}</el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="closeChangePhone">Hủy</el-button>
          <el-button type="primary" @click="changePhone">Thay đổi</el-button>
        </span>
      </template>
    </el-dialog>

    <el-dialog v-model="changeEmailFlag" title="Liên kết email" width="600px">
      <el-form :model="emailForm">
        <el-form-item label="Email" label-width="120px">
          <el-input v-model="emailForm.email" placeholder="Vui lòng nhập email" autocomplete="off" />
        </el-form-item>
        <el-form-item label="Mã xác nhận" label-width="120px">
          <div class="flex w-full gap-4">
            <el-input v-model="emailForm.code" class="flex-1" placeholder="Vui lòng tự thiết kế dịch vụ email, ở đây là giả lập viết bừa" autocomplete="off"
              style="width:300px" />
            <el-button type="primary" :disabled="emailTime > 0" @click="getEmailCode">{{
              emailTime > 0 ? `(${emailTime}s) Sau khi lấy lại` : 'Lấy mã xác nhận' }}</el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="closeChangeEmail">Hủy</el-button>
          <el-button type="primary" @click="changeEmail">Thay đổi</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { setSelfInfo, changePassword } from '@/api/user.js'
import { reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/pinia/modules/user'
import SelectImage from '@/components/selectImage/selectImage.vue'

defineOptions({
  name: 'Person',
})

const activeName = ref('second')
const rules = reactive({
  password: [
    { required: true, message: 'Vui lòng nhập mật khẩu', trigger: 'blur' },
    { min: 6, message: 'Ít nhất 6 ký tự', trigger: 'blur' },
  ],
  newPassword: [
    { required: true, message: 'Vui lòng nhập mật khẩu mới', trigger: 'blur' },
    { min: 6, message: 'Ít nhất 6 ký tự', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: 'Vui lòng nhập lại mật khẩu', trigger: 'blur' },
    { min: 6, message: 'Ít nhất 6 ký tự', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== pwdModify.value.newPassword) {
          callback(new Error('Mật khẩu không khớp'))
        } else {
          callback()
        }
      },
      trigger: 'blur',
    },
  ],
})

const userStore = useUserStore()
const modifyPwdForm = ref(null)
const showPassword = ref(false)
const pwdModify = ref({})
const nickName = ref('')
const editFlag = ref(false)
const savePassword = async () => {
  modifyPwdForm.value.validate((valid) => {
    if (valid) {
      changePassword({
        password: pwdModify.value.password,
        newPassword: pwdModify.value.newPassword,
      }).then((res) => {
        if (res.code === 0) {
          ElMessage.success('Đổi mật khẩu thành công!')
        }
        showPassword.value = false
      })
    } else {
      return false
    }
  })
}

const clearPassword = () => {
  pwdModify.value = {
    password: '',
    newPassword: '',
    confirmPassword: '',
  }
  modifyPwdForm.value.clearValidate()
}

watch(() => userStore.userInfo.headerImg, async (val) => {
  const res = await setSelfInfo({ headerImg: val })
  if (res.code === 0) {
    userStore.ResetUserInfo({ headerImg: val })
    ElMessage({
      type: 'success',
      message: 'Cài đặt thành công',
    })
  }
})

const openEdit = () => {
  nickName.value = userStore.userInfo.nickName
  editFlag.value = true
}

const closeEdit = () => {
  nickName.value = ''
  editFlag.value = false
}

const enterEdit = async () => {
  const res = await setSelfInfo({
    nickName: nickName.value
  })
  if (res.code === 0) {
    userStore.ResetUserInfo({ nickName: nickName.value })
    ElMessage({
      type: 'success',
      message: 'Cài đặt thành công',
    })
  }
  nickName.value = ''
  editFlag.value = false
}

const handleClick = (tab, event) => {
  console.log(tab, event)
}

const changePhoneFlag = ref(false)
const time = ref(0)
const phoneForm = reactive({
  phone: '',
  code: ''
})

const getCode = async () => {
  time.value = 60
  let timer = setInterval(() => {
    time.value--
    if (time.value <= 0) {
      clearInterval(timer)
      timer = null
    }
  }, 1000)
}

const closeChangePhone = () => {
  changePhoneFlag.value = false
  phoneForm.phone = ''
  phoneForm.code = ''
}

const changePhone = async () => {
  const res = await setSelfInfo({ phone: phoneForm.phone })
  if (res.code === 0) {
    ElMessage.success('Thay đổi thành công')
    userStore.ResetUserInfo({ phone: phoneForm.phone })
    closeChangePhone()
  }
}

const changeEmailFlag = ref(false)
const emailTime = ref(0)
const emailForm = reactive({
  email: '',
  code: ''
})

const getEmailCode = async () => {
  emailTime.value = 60
  let timer = setInterval(() => {
    emailTime.value--
    if (emailTime.value <= 0) {
      clearInterval(timer)
      timer = null
    }
  }, 1000)
}

const closeChangeEmail = () => {
  changeEmailFlag.value = false
  emailForm.email = ''
  emailForm.code = ''
}

const changeEmail = async () => {
  const res = await setSelfInfo({ email: emailForm.email })
  if (res.code === 0) {
    ElMessage.success('Thay đổi thành công')
    userStore.ResetUserInfo({ email: emailForm.email })
    closeChangeEmail()
  }
}

</script>

<style lang="scss">
.borderd {
  @apply border-b-2 border-solid border-gray-100 dark:border-gray-500 border-t-0 border-r-0 border-l-0;

  &:last-child {
    @apply border-b-0;
  }
}

.info-list {
  @apply w-full whitespace-nowrap overflow-hidden text-ellipsis py-3 text-lg text-gray-700
}
</style>
