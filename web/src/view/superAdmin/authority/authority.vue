<template>
  <div class="authority">
    <warning-bar title="Chú ý: Bạn có thể chuyển đổi vai trò trong menu thả xuống ở góc trên bên phải" />
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button
          type="primary"
          icon="plus"
          @click="addAuthority(0)"
        >Thêm vai trò</el-button>
      </div>
      <el-table
        :data="tableData"
        :tree-props="{children: 'children', hasChildren: 'hasChildren'}"
        row-key="authorityId"
        style="width: 100%"
      >
        <el-table-column
          label="ID vai trò"
          min-width="180"
          prop="authorityId"
        />
        <el-table-column
          align="left"
          label="Tên vai trò"
          min-width="180"
          prop="authorityName"
        />
        <el-table-column
          align="left"
          label="Hành động"
          width="460"
        >
          <template #default="scope">
            <el-button
              icon="setting"
              type="primary"
              link
              @click="openDrawer(scope.row)"
            >Cấu hình quyền</el-button>
            <el-button
              icon="plus"
              type="primary"
              link
              @click="addAuthority(scope.row.authorityId)"
            >Thêm vai trò con</el-button>
            <el-button
              icon="copy-document"
              type="primary"
              link
              @click="copyAuthorityFunc(scope.row)"
            >Sao chép</el-button>
            <el-button
              icon="edit"
              type="primary"
              link
              @click="editAuthority(scope.row)"
            >Chỉnh sửa</el-button>
            <el-button
              icon="delete"
              type="primary"
              link
              @click="deleteAuth(scope.row)"
            >Xóa</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
    <!-- Dialog thêm vai trò -->
    <el-drawer
      v-model="authorityFormVisible"
      :show-close="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ authorityTitleForm }}</span>
          <div>
            <el-button @click="closeAuthorityForm">Hủy</el-button>
            <el-button
              type="primary"
              @click="submitAuthorityForm"
            >Xác nhận</el-button>
          </div>
        </div>
      </template>
      <el-form
        ref="authorityForm"
        :model="form"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item
          label="Vai trò cha"
          prop="parentId"
        >
          <el-cascader
            v-model="form.parentId"
            style="width:100%"
            :disabled="dialogType==='add'"
            :options="AuthorityOption"
            :props="{ checkStrictly: true,label:'authorityName',value:'authorityId',disabled:'disabled',emitPath:false}"
            :show-all-levels="false"
            filterable
          />
        </el-form-item>
        <el-form-item
          label="ID vai trò"
          prop="authorityId"
        >
          <el-input
            v-model="form.authorityId"
            :disabled="dialogType==='edit'"
            autocomplete="off"
            maxlength="15"
          />
        </el-form-item>
        <el-form-item
          label="Tên vai trò"
          prop="authorityName"
        >
          <el-input
            v-model="form.authorityName"
            autocomplete="off"
          />
        </el-form-item>
      </el-form>
    </el-drawer>

    <el-drawer
      v-if="drawer"
      v-model="drawer"
      :with-header="false"
      size="40%"
      title="Cấu hình vai trò"
    >
      <el-tabs
        :before-leave="autoEnter"
        type="border-card"
      >
        <el-tab-pane label="Menu vai trò">
          <Menus
            ref="menus"
            :row="activeRow"
            @changeRow="changeRow"
          />
        </el-tab-pane>
        <el-tab-pane label="API vai trò">
          <Apis
            ref="apis"
            :row="activeRow"
            @changeRow="changeRow"
          />
        </el-tab-pane>
        <el-tab-pane label="Quyền tài nguyên">
          <Datas
            ref="datas"
            :authority="tableData"
            :row="activeRow"
            @changeRow="changeRow"
          />
        </el-tab-pane>
      </el-tabs>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  getAuthorityList,
  deleteAuthority,
  createAuthority,
  updateAuthority,
  copyAuthority
} from '@/api/authority'

import Menus from '@/view/superAdmin/authority/components/menus.vue'
import Apis from '@/view/superAdmin/authority/components/apis.vue'
import Datas from '@/view/superAdmin/authority/components/datas.vue'
import WarningBar from '@/components/warningBar/warningBar.vue'

import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'


defineOptions({
  name: 'Authority'
})

const mustUint = (rule, value, callback) => {
  if (!/^[0-9]*[1-9][0-9]*$/.test(value)) {
    return callback(new Error('Vui lòng nhập số nguyên dương'))
  }
  return callback()
}

const AuthorityOption = ref([
  {
    authorityId: 0,
    authorityName: 'Vai trò gốc'
  }
])
const drawer = ref(false)
const dialogType = ref('add')
const activeRow = ref({})

const authorityTitleForm = ref('Thêm vai trò')
const authorityFormVisible = ref(false)
const apiDialogFlag = ref(false)
const copyForm = ref({})

const form = ref({
  authorityId: 0,
  authorityName: '',
  parentId: 0
})
const rules = ref({
  authorityId: [
    { required: true, message: 'Vui lòng nhập ID vai trò', trigger: 'blur' },
    { validator: mustUint, trigger: 'blur', message: 'Phải là số nguyên dương' }
  ],
  authorityName: [
    { required: true, message: 'Vui lòng nhập tên vai trò', trigger: 'blur' }
  ],
  parentId: [
    { required: true, message: 'Vui lòng chọn vai trò cha', trigger: 'blur' },
  ]
})

const page = ref(1)
const total = ref(0)
const pageSize = ref(999)
const tableData = ref([])
const searchInfo = ref({})

// Lấy dữ liệu bảng
const getTableData = async() => {
  const table = await getAuthorityList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

const changeRow = (key, value) => {
  activeRow.value[key] = value
}
const menus = ref(null)
const apis = ref(null)
const datas = ref(null)
const autoEnter = (activeName, oldActiveName) => {
  const paneArr = [menus, apis, datas]
  if (oldActiveName) {
    if (paneArr[oldActiveName].value.needConfirm) {
      paneArr[oldActiveName].value.enterAndNext()
      paneArr[oldActiveName].value.needConfirm = false
    }
  }
}
// Sao chép vai trò
const copyAuthorityFunc = (row) => {
  setOptions()
  authorityTitleForm.value = 'Sao chép vai trò'
  dialogType.value = 'copy'
  for (const k in form.value) {
    form.value[k] = row[k]
  }
  copyForm.value = row
  authorityFormVisible.value = true
}
const openDrawer = (row) => {
  drawer.value = true
  activeRow.value = row
}
// Xóa vai trò
const deleteAuth = (row) => {
  ElMessageBox.confirm('Thao tác này sẽ xóa vai trò vĩnh viễn, bạn có muốn tiếp tục?', 'Cảnh báo', {
    confirmButtonText: 'Đồng ý',
    cancelButtonText: 'Hủy',
    type: 'warning'
  })
    .then(async() => {
      const res = await deleteAuthority({ authorityId: row.authorityId })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: 'Xóa thành công!'
        })
        if (tableData.value.length === 1 && page.value > 1) {
          page.value--
        }
        getTableData()
      }
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: 'Đã hủy xóa'
      })
    })
}
// Khởi tạo form
const authorityForm = ref(null)
const initForm = () => {
  if (authorityForm.value) {
    authorityForm.value.resetFields()
  }
  form.value = {
    authorityId: 0,
    authorityName: '',
    parentId: 0
  }
}
// Đóng cửa sổ
const closeAuthorityForm = () => {
  initForm()
  authorityFormVisible.value = false
  apiDialogFlag.value = false
}
// Xác nhận form
const submitAuthorityForm = () => {
  authorityForm.value.validate(async valid => {
    if (valid) {
      form.value.authorityId = Number(form.value.authorityId)
      switch (dialogType.value) {
        case 'add':
          {
            const res = await createAuthority(form.value)
            if (res.code === 0) {
              ElMessage({
                type: 'success',
                message: 'Thêm thành công!'
              })
              getTableData()
              closeAuthorityForm()
            }
          }
          break
        case 'edit':
          {
            const res = await updateAuthority(form.value)
            if (res.code === 0) {
              ElMessage({
                type: 'success',
                message: 'Cập nhật thành công!'
              })
              getTableData()
              closeAuthorityForm()
            }
          }
          break
        case 'copy': {
          const data = {
            authority: {
              authorityId: 0,
              authorityName: '',
              datauthorityId: [],
              parentId: 0
            },
            oldAuthorityId: 0
          }
          data.authority.authorityId = form.value.authorityId
          data.authority.authorityName = form.value.authorityName
          data.authority.parentId = form.value.parentId
          data.authority.dataAuthorityId = copyForm.value.dataAuthorityId
          data.oldAuthorityId = copyForm.value.authorityId
          const res = await copyAuthority(data)
          if (res.code === 0) {
            ElMessage({
              type: 'success',
              message: 'Sao chép thành công!'
            })
            getTableData()
          }
        }
      }

      initForm()
      authorityFormVisible.value = false
    }
  })
}
const setOptions = () => {
  AuthorityOption.value = [
    {
      authorityId: 0,
      authorityName: 'Vai trò gốc'
    }
  ]
  setAuthorityOptions(tableData.value, AuthorityOption.value, false)
}
const setAuthorityOptions = (AuthorityData, optionsData, disabled) => {
  AuthorityData &&
        AuthorityData.forEach(item => {
          if (item.children && item.children.length) {
            const option = {
              authorityId: item.authorityId,
              authorityName: item.authorityName,
              disabled: disabled || item.authorityId === form.value.authorityId,
              children: []
            }
            setAuthorityOptions(
              item.children,
              option.children,
              disabled || item.authorityId === form.value.authorityId
            )
            optionsData.push(option)
          } else {
            const option = {
              authorityId: item.authorityId,
              authorityName: item.authorityName,
              disabled: disabled || item.authorityId === form.value.authorityId
            }
            optionsData.push(option)
          }
        })
}
// Thêm vai trò
const addAuthority = (parentId) => {
  initForm()
  authorityTitleForm.value = 'Thêm vai trò'
  dialogType.value = 'add'
  form.value.parentId = parentId
  setOptions()
  authorityFormVisible.value = true
}
// Chỉnh sửa vai trò
const editAuthority = (row) => {
  setOptions()
  authorityTitleForm.value = 'Chỉnh sửa vai trò'
  dialogType.value = 'edit'
  for (const key in form.value) {
    form.value[key] = row[key]
  }
  setOptions()
  authorityForm.value && authorityForm.value.clearValidate()
  authorityFormVisible.value = true
}

</script>

<style lang="scss">
.authority {
  .el-input-number {
    margin-left: 15px;
    span {
      display: none;
    }
  }
}
.tree-content{
  margin-top: 10px;
  height: calc(100vh - 158px);
  overflow: auto;
}

</style>
