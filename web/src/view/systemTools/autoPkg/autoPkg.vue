<template>
  <div>
    <!-- <warning-bar
      href="https://www.bilibili.com/video/BV1kv4y1g7nT?p=3"
      title="This feature is for development environment only and not recommended for production. For specific usage, please refer to the video: https://www.bilibili.com/video/BV1kv4y1g7nT?p=3"
    /> -->
    <div class="gva-table-box">
      <div class="gva-btn-list gap-3 flex items-center">
        <el-button
          type="primary"
          icon="plus"
          @click="openDialog('addApi')"
        >
          Add
        </el-button>

      </div>
      <el-table :data="tableData">
        <el-table-column
          align="left"
          label="id"
          width="60"
          prop="ID"
        />
        <el-table-column
          align="left"
          label="Package Name"
          width="150"
          prop="packageName"
        />
        <el-table-column
            align="left"
            label="Template"
            width="150"
            prop="template"
        />
        <el-table-column
          align="left"
          label="Display Name"
          width="150"
          prop="label"
        />
        <el-table-column
          align="left"
          label="Description"
          min-width="150"
          prop="desc"
        />

        <el-table-column
          align="left"
          label="Action"
          width="200"
        >
          <template #default="scope">
            <el-button
              icon="delete"

              type="primary"
              link
              @click="deleteApiFunc(scope.row)"
            >
              Delete
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-drawer
      v-model="dialogFormVisible"
      size="40%"
      :show-close="false"
    >
      <warning-bar title="Template package will create a code package integrated into the project, while template plugin will create a plugin package" />
      <el-form
        ref="pkgForm"
        :model="form"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item
          label="Package Name"
          prop="packageName"
        >
          <el-input
            v-model="form.packageName"
            autocomplete="off"
          />
        </el-form-item>
        <el-form-item
          label="Template"
          prop="template"
        >
          <el-select
            v-model="form.template"
          >
            <el-option v-for="template in templatesOptions" :label="template" :value="template" :key="template"/>
          </el-select>
        </el-form-item>

        <el-form-item
          label="Display Name"
          prop="label"
        >
          <el-input
            v-model="form.label"
            autocomplete="off"
          />
        </el-form-item>
        <el-form-item
          label="Description"
          prop="desc"
        >
          <el-input
            v-model="form.desc"
            autocomplete="off"
          />
        </el-form-item>
      </el-form>
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">Create Package</span>
          <div>
            <el-button @click="closeDialog">
              Cancel
            </el-button>
            <el-button
              type="primary"
              @click="enterDialog"
            >
              Confirm
            </el-button>
          </div>
        </div>
      </template>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createPackageApi,
  getPackageApi,
  deletePackageApi,
  getTemplatesApi
} from '@/api/autoCode'
import { ref } from 'vue'
import WarningBar from '@/components/warningBar/warningBar.vue'
import { ElMessage, ElMessageBox } from 'element-plus'

defineOptions({
  name: 'AutoPkg',
})

const form = ref({
  packageName: '',
  template: '',
  label: '',
  desc: '',
})
const templatesOptions = ref([])

const getTemplates = async ()=>{
  const res = await getTemplatesApi()
  if (res.code === 0){
    templatesOptions.value = res.data
  }
}

getTemplates()

const validateNum = (rule, value, callback) => {
  if ((/^\d+$/.test(value[0]))) {
    callback(new Error('不能够以数字开头'))
  } else {
    callback()
  }
}

const rules = ref({
  packageName: [
    { required: true, message: '请输入包名', trigger: 'blur' },
    { validator: validateNum, trigger: 'blur' }
  ],
  template:[
    { required: true, message: '请选择模板', trigger: 'change' },
    { validator: validateNum, trigger: 'blur' }
  ]
})

const dialogFormVisible = ref(false)
const openDialog = () => {
  dialogFormVisible.value = true
}

const closeDialog = () => {
  dialogFormVisible.value = false
  form.value = {
    packageName: '',
    template: '',
    label: '',
    desc: '',
  }
}

const pkgForm = ref(null)
const enterDialog = async() => {
  pkgForm.value.validate(async valid => {
    if (valid) {
      const res = await createPackageApi(form.value)
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '添加成功',
          showClose: true
        })
      }
      getTableData()
      closeDialog()
    }
  })
}

const tableData = ref([])
const getTableData = async() => {
  const table = await getPackageApi()
  if (table.code === 0) {
    tableData.value = table.data.pkgs
  }
}

const deleteApiFunc = async(row) => {
  ElMessageBox.confirm('此操作仅删除数据库中的pkg存储，后端相应目录结构请自行删除与数据库保持一致！', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
    .then(async() => {
      const res = await deletePackageApi(row)
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功!'
        })
        getTableData()
      }
    })
}

getTableData()
</script>
