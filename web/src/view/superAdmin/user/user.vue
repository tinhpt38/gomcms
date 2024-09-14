<template>
  <div>
    <warning-bar title="Chú ý: Nhấp vào biểu tượng hình đại diện ở góc trên bên phải để chuyển đổi vai trò" />
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="addUser">Thêm người dùng</el-button>
      </div>
      <el-table :data="tableData" row-key="ID">
        <el-table-column align="left" label="#" min-width="75">
          <template #default="scope">
            <CustomPic style="margin-top:8px" :pic-src="scope.row.headerImg" />
          </template>
        </el-table-column>
        <el-table-column align="left" label="ID" min-width="50" prop="ID" />
        <el-table-column align="left" label="Tên người dùng" min-width="150" prop="userName" />
        <el-table-column align="left" label="Biệt danh" min-width="150" prop="nickName" />
        <!-- <el-table-column align="left" label="Số điện thoại" min-width="180" prop="phone" /> -->
        <el-table-column align="left" label="Email" min-width="200" prop="email" />
        <el-table-column align="left" label="Vai trò" min-width="200">
          <template #default="scope">
            <el-cascader v-model="scope.row.authorityIds" :options="authOptions" :show-all-levels="false" collapse-tags
              :props="{ multiple: true, checkStrictly: true, label: 'authorityName', value: 'authorityId', disabled: 'disabled', emitPath: false }"
              :clearable="false" @visible-change="(flag) => { changeAuthority(scope.row, flag, 0) }"
              @remove-tag="(removeAuth) => { changeAuthority(scope.row, false, removeAuth) }" />
          </template>
        </el-table-column>
        <el-table-column align="left" label="Kích hoạt" min-width="150">
          <template #default="scope">
            <el-switch v-model="scope.row.enable" inline-prompt :active-value="1" :inactive-value="2"
              @change="() => { switchEnable(scope.row) }" />
          </template>
        </el-table-column>

        <el-table-column label="Hành động" min-width="250" fixed="right">
          <template #default="scope">
            <el-button type="primary" link icon="delete" @click="deleteUserFunc(scope.row)">Xóa</el-button>
            <el-button type="primary" link icon="edit" @click="openEdit(scope.row)">Chỉnh sửa</el-button>
            <el-button type="primary" link icon="magic-stick" @click="resetPasswordFunc(scope.row)">Đặt lại mật
              khẩu</el-button>
          </template>
        </el-table-column>

      </el-table>
      <div class="gva-pagination">
        <el-pagination :current-page="page" :page-size="pageSize" :page-sizes="[10, 30, 50, 100]" :total="total"
          layout="total, sizes, prev, pager, next, jumper" @current-change="handleCurrentChange"
          @size-change="handleSizeChange" />
      </div>
    </div>
    <el-drawer v-model="addUserDialog" size="60%" :show-close="false" :close-on-press-escape="false"
      :close-on-click-modal="false">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">Người dùng</span>
          <div>
            <el-button @click="closeAddUserDialog">Hủy</el-button>
            <el-button type="primary" @click="enterAddUserDialog">Xác nhận</el-button>
          </div>
        </div>
      </template>

      <el-form ref="userForm" :rules="rules" :model="userInfo" label-width="80px">
        <el-form-item v-if="dialogFlag === 'add'" label="Tên người dùng" prop="userName">
          <el-input v-model="userInfo.userName" />
        </el-form-item>
        <el-form-item v-if="dialogFlag === 'add'" label="Mật khẩu" prop="password">
          <el-input v-model="userInfo.password" />
        </el-form-item>
        <el-form-item label="Biệt danh" prop="nickName">
          <el-input v-model="userInfo.nickName" />
        </el-form-item>
        <el-form-item label="Số điện thoại" prop="phone">
          <el-input v-model="userInfo.phone" />
        </el-form-item>
        <el-form-item label="Email" prop="email">
          <el-input v-model="userInfo.email" />
        </el-form-item>
        <el-form-item label="Vai trò người dùng" prop="authorityId">
          <el-cascader v-model="userInfo.authorityIds" style="width:100%" :options="authOptions"
            :show-all-levels="false"
            :props="{ multiple: true, checkStrictly: true, label: 'authorityName', value: 'authorityId', disabled: 'disabled', emitPath: false }"
            :clearable="false" />
        </el-form-item>
        <el-form-item label="Kích hoạt" prop="disabled">
          <el-switch v-model="userInfo.enable" inline-prompt :active-value="1" :inactive-value="2" />
        </el-form-item>
        <el-form-item label="Hình đại diện" label-width="80px">
          <SelectImage v-model="userInfo.headerImg" />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>

import {
  getUserList,
  setUserAuthorities,
  register,
  deleteUser
} from '@/api/user'

import { getAuthorityList } from '@/api/authority'
import CustomPic from '@/components/customPic/index.vue'
import WarningBar from '@/components/warningBar/warningBar.vue'
import { setUserInfo, resetPassword } from '@/api/user.js'

import { nextTick, ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import SelectImage from '@/components/selectImage/selectImage.vue'

defineOptions({
  name: 'User',
})

const path = ref(import.meta.env.VITE_BASE_API + '/')
// Khởi tạo
const setAuthorityOptions = (AuthorityData, optionsData) => {
  AuthorityData &&
    AuthorityData.forEach(item => {
      if (item.children && item.children.length) {
        const option = {
          authorityId: item.authorityId,
          authorityName: item.authorityName,
          children: []
        }
        setAuthorityOptions(item.children, option.children)
        optionsData.push(option)
      } else {
        const option = {
          authorityId: item.authorityId,
          authorityName: item.authorityName
        }
        optionsData.push(option)
      }
    })
}

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
// Phân trang
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// Truy vấn
const getTableData = async () => {
  const table = await getUserList({ page: page.value, pageSize: pageSize.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

watch(() => tableData.value, () => {
  setAuthorityIds()
})

const initPage = async () => {
  getTableData()
  const res = await getAuthorityList({ page: 1, pageSize: 999 })
  setOptions(res.data.list)
}

initPage()

const resetPasswordFunc = (row) => {
  ElMessageBox.confirm(
    'Bạn có chắc muốn đặt lại mật khẩu của người dùng này thành 123456?',
    'Cảnh báo',
    {
      confirmButtonText: 'Đồng ý',
      cancelButtonText: 'Hủy',
      type: 'warning',
    }
  ).then(async () => {
    const res = await resetPassword({
      ID: row.ID,
    })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: res.msg,
      })
    } else {
      ElMessage({
        type: 'error',
        message: res.msg,
      })
    }
  })
}
const setAuthorityIds = () => {
  tableData.value && tableData.value.forEach((user) => {
    user.authorityIds = user.authorities && user.authorities.map(i => {
      return i.authorityId
    })
  })
}

const chooseImg = ref(null)
const openHeaderChange = () => {
  chooseImg.value.open()
}

const authOptions = ref([])
const setOptions = (authData) => {
  authOptions.value = []
  setAuthorityOptions(authData, authOptions.value)
}

const deleteUserFunc = async (row) => {
  ElMessageBox.confirm('Bạn có chắc muốn xóa không?', 'Thông báo', {
    confirmButtonText: 'Đồng ý',
    cancelButtonText: 'Hủy',
    type: 'warning'
  }).then(async () => {
    const res = await deleteUser({ id: row.ID })
    if (res.code === 0) {
      ElMessage.success('Xóa thành công')
      await getTableData()
    }
  })
}

// Liên quan đến hộp thoại
const userInfo = ref({
  userName: '',
  password: '',
  nickName: '',
  headerImg: '',
  authorityId: '',
  authorityIds: [],
  enable: 1,
})

const rules = ref({
  userName: [
    { required: true, message: 'Vui lòng nhập tên người dùng', trigger: 'blur' },
    { min: 5, message: 'Ít nhất 5 ký tự', trigger: 'blur' }
  ],
  password: [
    { required: true, message: 'Vui lòng nhập mật khẩu người dùng', trigger: 'blur' },
    { min: 6, message: 'Ít nhất 6 ký tự', trigger: 'blur' }
  ],
  nickName: [
    { required: true, message: 'Vui lòng nhập biệt danh người dùng', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1([38][0-9]|4[014-9]|[59][0-35-9]|6[2567]|7[0-8])\d{8}$/, message: 'Vui lòng nhập số điện thoại hợp lệ', trigger: 'blur' },
  ],
  email: [
    { pattern: /^([0-9A-Za-z\-_.]+)@([0-9a-z]+\.[a-z]{2,3}(\.[a-z]{2})?)$/g, message: 'Vui lòng nhập địa chỉ email hợp lệ', trigger: 'blur' },
  ],
  authorityId: [
    { required: true, message: 'Vui lòng chọn vai trò người dùng', trigger: 'blur' }
  ]
})
const userForm = ref(null)
const enterAddUserDialog = async () => {
  userInfo.value.authorityId = userInfo.value.authorityIds[0]
  userForm.value.validate(async valid => {
    if (valid) {
      const req = {
        ...userInfo.value
      }
      if (dialogFlag.value === 'add') {
        const res = await register(req)
        if (res.code === 0) {
          ElMessage({ type: 'success', message: 'Tạo thành công' })
          await getTableData()
          closeAddUserDialog()
        }
      }
      if (dialogFlag.value === 'edit') {
        const res = await setUserInfo(req)
        if (res.code === 0) {
          ElMessage({ type: 'success', message: 'Chỉnh sửa thành công' })
          await getTableData()
          closeAddUserDialog()
        }
      }
    }
  })
}

const addUserDialog = ref(false)
const closeAddUserDialog = () => {
  userForm.value.resetFields()
  userInfo.value.headerImg = ''
  userInfo.value.authorityIds = []
  addUserDialog.value = false
}

const dialogFlag = ref('add')

const addUser = () => {
  dialogFlag.value = 'add'
  addUserDialog.value = true
}

const tempAuth = {}
const changeAuthority = async (row, flag, removeAuth) => {
  if (flag) {
    if (!removeAuth) {
      tempAuth[row.ID] = [...row.authorityIds]
    }
    return
  }
  await nextTick()
  const res = await setUserAuthorities({
    ID: row.ID,
    authorityIds: row.authorityIds
  })
  if (res.code === 0) {
    ElMessage({ type: 'success', message: 'Cài đặt vai trò thành công' })
  } else {
    if (!removeAuth) {
      row.authorityIds = [...tempAuth[row.ID]]
      delete tempAuth[row.ID]
    } else {
      row.authorityIds = [removeAuth, ...row.authorityIds]
    }
  }
}

const openEdit = (row) => {
  dialogFlag.value = 'edit'
  userInfo.value = JSON.parse(JSON.stringify(row))
  addUserDialog.value = true
}

const switchEnable = async (row) => {
  userInfo.value = JSON.parse(JSON.stringify(row))
  await nextTick()
  const req = {
    ...userInfo.value
  }
  const res = await setUserInfo(req)
  if (res.code === 0) {
    ElMessage({ type: 'success', message: `${req.enable === 2 ? 'Vô hiệu hóa' : 'Kích hoạt'} thành công` })
    await getTableData()
    userInfo.value.headerImg = ''
    userInfo.value.authorityIds = []
  }
}

</script>

<style lang="scss">
.header-img-box {
  @apply w-52 h-52 border border-solid border-gray-300 rounded-xl flex justify-center items-center cursor-pointer;
}
</style>
