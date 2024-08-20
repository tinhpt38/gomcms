<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="addMenu(0)">
          Thêm menu gốc
        </el-button>
      </div>

      <!-- Vì menu ở đây tương ứng với danh sách bên trái nên không cần phân trang, pageSize mặc định là 999 -->
      <el-table :data="tableData" row-key="ID">
        <el-table-column align="left" label="ID" min-width="100" prop="ID" />
        <el-table-column align="left" label="Tên hiển thị" min-width="120" prop="authorityName">
          <template #default="scope">
            <span>{{ scope.row.meta.title }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="Biểu tượng" min-width="140" prop="authorityName">
          <template #default="scope">
            <div v-if="scope.row.meta.icon" class="icon-column">
              <el-icon>
                <component :is="scope.row.meta.icon" />
              </el-icon>
              <span>{{ scope.row.meta.icon }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="Tên đường dẫn" show-overflow-tooltip min-width="160" prop="name" />
        <el-table-column align="left" label="Đường dẫn" show-overflow-tooltip min-width="160" prop="path" />
        <el-table-column align="left" label="Ẩn/Hiện" min-width="100" prop="hidden">
          <template #default="scope">
            <span>{{ scope.row.hidden ? "Ẩn" : "Hiện" }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="Nút cha" min-width="90" prop="parentId" />
        <el-table-column align="left" label="Thứ tự" min-width="70" prop="sort" />
        <el-table-column align="left" label="Đường dẫn tệp" min-width="360" prop="component" />
        <el-table-column align="left" fixed="right" label="Thao tác" width="300">
          <template #default="scope">
            <el-button type="primary" link icon="plus" @click="addMenu(scope.row.ID)">
              Thêm menu con
            </el-button>
            <el-button type="primary" link icon="edit" @click="editMenu(scope.row.ID)">
              Chỉnh sửa
            </el-button>
            <el-button type="primary" link icon="delete" @click="deleteMenu(scope.row.ID)">
              Xóa
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
    <el-drawer v-model="dialogFormVisible" size="60%" :before-close="handleClose" :show-close="false">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ dialogTitle }}</span>
          <div>
            <el-button @click="closeDialog">
              Hủy bỏ
            </el-button>
            <el-button type="primary" @click="enterDialog">
              Xác nhận
            </el-button>
          </div>
        </div>
      </template>

      <warning-bar title="Thêm menu, cần cấu hình quyền trong quản lý vai trò để sử dụng" />
      <el-form v-if="dialogFormVisible" ref="menuForm" :inline="true" :model="form" :rules="rules" label-position="top">
        <el-row class="w-full">
          <el-col :span="16">
            <el-form-item label="Đường dẫn tệp" prop="component">
              <el-select v-model="form.component" filterable allow-create autocomplete="off" style="width: 100%"
                placeholder="Trang:view/xxx/xx.vue Plugin:plugin/xx/xx.vue" default-first-option @change="fmtComponent">
                <el-option v-for="(item, path) in pathOptions" :key="path" :label="path" :value="path" />
              </el-select>
              <span style="font-size: 12px; margin-right: 12px">Nếu menu có menu con, vui lòng tạo trang router-view cấp
                hai
                hoặc</span><el-button style="margin-top: 4px" @click="form.component = 'view/routerHolder.vue'">
                Nhấn vào đây để cài đặt
              </el-button>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="Tên hiển thị" prop="meta.title">
              <el-input v-model="form.meta.title" autocomplete="off" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row class="w-full">
          <el-col :span="8">
            <el-form-item label="Tên đường dẫn" prop="path">
              <el-input v-model="form.name" autocomplete="off" placeholder="Chuỗi duy nhất bằng tiếng Anh"
                @change="changeName" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item prop="path">
              <template #label>
                <span style="display: inline-flex; align-items: center">
                  <span>Đường dẫn</span>
                  <el-checkbox v-model="checkFlag" style="margin-left: 12px; height: auto">Thêm tham số</el-checkbox>
                </span>
              </template>

              <el-input v-model="form.path" :disabled="!checkFlag" autocomplete="off"
                placeholder="Nên chỉ thêm tham số ở phía sau" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="Ẩn/Hiện">
              <el-select v-model="form.hidden" style="width: 100%" placeholder="Ẩn/Hiện trong danh sách">
                <el-option :value="false" label="Không" />
                <el-option :value="true" label="Có" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row class="w-full">
          <el-col :span="8">
            <el-form-item label="ID nút cha">
              <el-cascader v-model="form.parentId" style="width: 100%" :disabled="!isEdit" :options="menuOption" :props="{
                checkStrictly: true,
                label: 'title',
                value: 'ID',
                disabled: 'disabled',
                emitPath: false,
              }" :show-all-levels="false" filterable />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="Biểu tượng" prop="meta.icon">
              <icon :meta="form.meta" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="Đánh dấu sắp xếp" prop="sort">
              <el-input v-model.number="form.sort" autocomplete="off" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row class="w-full">
          <el-col :span="8">
            <el-form-item prop="meta.activeName">
              <template #label>
                <div>
                  <span>Menu được tô sáng</span>
                  <el-tooltip
                    content="Chú ý: Khi đến đường dẫn này, menu bên trái được chỉ định sẽ ở trạng thái hoạt động (sáng lên), có thể để trống, nếu để trống thì sẽ là tên đường dẫn hiện tại."
                    placement="top" effect="light">
                    <el-icon>
                      <QuestionFilled />
                    </el-icon>
                  </el-tooltip>
                </div>
              </template>
              <el-input v-model="form.meta.activeName" :placeholder="form.name" autocomplete="off" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="KeepAlive" prop="meta.keepAlive">
              <el-select v-model="form.meta.keepAlive" style="width: 100%"
                placeholder="Có giữ trạng thái trang hay không">
                <el-option :value="false" label="Không" />
                <el-option :value="true" label="Có" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="CloseTab" prop="meta.closeTab">
              <el-select v-model="form.meta.closeTab" style="width: 100%" placeholder="Tự động đóng tab hay không">
                <el-option :value="false" label="Không" />
                <el-option :value="true" label="Có" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row class="w-full">
          <el-col :span="8">
            <el-form-item>
              <template #label>
                <div>
                  <span>Trang cơ bản</span>
                  <el-tooltip content="Nếu chọn là có, thì không hiển thị menu bên trái và thông tin phía trên."
                    placement="top" effect="light">
                    <el-icon>
                      <QuestionFilled />
                    </el-icon>
                  </el-tooltip>
                </div>
              </template>

              <el-select v-model="form.meta.defaultMenu" style="width: 100%" placeholder="Trang cơ bản hay không">
                <el-option :value="false" label="Không" />
                <el-option :value="true" label="Có" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <div>
        <div class="flex items-center gap-2">
          <el-button type="primary" icon="edit" @click="addParameter(form)">
            Thêm tham số menu
          </el-button>
        </div>
        <el-table :data="form.parameters" style="width: 100%; margin-top: 12px">
          <el-table-column align="left" prop="type" label="Loại tham số" width="180">
            <template #default="scope">
              <el-select v-model="scope.row.type" placeholder="Vui lòng chọn">
                <el-option key="query" value="query" label="query" />
                <el-option key="params" value="params" label="params" />
              </el-select>
            </template>
          </el-table-column>
          <el-table-column align="left" prop="key" label="Khóa tham số" width="180">
            <template #default="scope">
              <div>
                <el-input v-model="scope.row.key" />
              </div>
            </template>
          </el-table-column>
          <el-table-column align="left" prop="value" label="Giá trị tham số">
            <template #default="scope">
              <div>
                <el-input v-model="scope.row.value" />
              </div>
            </template>
          </el-table-column>
          <el-table-column align="left">
            <template #default="scope">
              <div>
                <el-button type="danger" icon="delete" @click="deleteParameter(form.parameters, scope.$index)">
                  Xóa
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>

        <div class="flex items-center gap-2 mt-3">
          <el-button type="primary" icon="edit" @click="addBtn(form)">
            Thêm nút điều khiển
          </el-button>
          <el-icon class="cursor-pointer" @click="
            toDoc('https://www.gin-vue-admin.com/guide/web/button-auth.html')
            ">
            <QuestionFilled />
          </el-icon>
        </div>

        <el-table :data="form.menuBtn" style="width: 100%; margin-top: 12px">
          <el-table-column align="left" prop="name" label="Tên nút" width="180">
            <template #default="scope">
              <div>
                <el-input v-model="scope.row.name" />
              </div>
            </template>
          </el-table-column>
          <el-table-column align="left" prop="name" label="Ghi chú" width="180">
            <template #default="scope">
              <div>
                <el-input v-model="scope.row.desc" />
              </div>
            </template>
          </el-table-column>
          <el-table-column align="left">
            <template #default="scope">
              <div>
                <el-button type="danger" icon="delete" @click="deleteBtn(form.menuBtn, scope.$index)">
                  Xóa
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  updateBaseMenu,
  getMenuList,
  addBaseMenu,
  deleteBaseMenu,
  getBaseMenuById,
} from '@/api/menu'
import icon from '@/view/superAdmin/menu/icon.vue'
import WarningBar from '@/components/warningBar/warningBar.vue'
import { canRemoveAuthorityBtnApi } from '@/api/authorityBtn'
import { reactive, ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { QuestionFilled } from '@element-plus/icons-vue'
import pathInfo from '@/pathInfo.json'

import { toDoc } from '@/utils/doc'
import { toLowerCase } from '@/utils/stringFun'

defineOptions({
  name: 'Menus',
})

const pathOptions = reactive({})

onMounted(() => {
  for (let pathInfoKey in pathInfo) {
    // Bỏ đi phần đầu tiên /src/
    pathOptions[pathInfoKey.replace(/^\/src\//, '')] = pathInfo[pathInfoKey]
  }
})

const rules = reactive({
  path: [{ required: true, message: 'Vui lòng nhập tên menu', trigger: 'blur' }],
  component: [{ required: true, message: 'Vui lòng nhập đường dẫn tệp tin', trigger: 'blur' }],
  'meta.title': [
    { required: true, message: 'Vui lòng nhập tên hiển thị menu', trigger: 'blur' },
  ],
})

const page = ref(1)
const total = ref(0)
const pageSize = ref(999)
const tableData = ref([])
const searchInfo = ref({})
// Tìm kiếm
const getTableData = async () => {
  const table = await getMenuList({
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo.value,
  })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// Thêm tham số
const addParameter = (form) => {
  if (!form.parameters) {
    form.parameters = []
  }
  form.parameters.push({
    type: 'query',
    key: '',
    value: '',
  })
}

const fmtComponent = () => {
  form.value.component = form.value.component.replace(/\\/g, '/')
  form.value.name = toLowerCase(pathOptions[form.value.component])
  form.value.path = form.value.name
}

// Xóa tham số
const deleteParameter = (parameters, index) => {
  parameters.splice(index, 1)
}

// Thêm nút điều khiển
const addBtn = (form) => {
  if (!form.menuBtn) {
    form.menuBtn = []
  }
  form.menuBtn.push({
    name: '',
    desc: '',
  })
}
// Xóa nút điều khiển
const deleteBtn = async (btns, index) => {
  const btn = btns[index]
  if (btn.ID === 0) {
    btns.splice(index, 1)
    return
  }
  const res = await canRemoveAuthorityBtnApi({ id: btn.ID })
  if (res.code === 0) {
    btns.splice(index, 1)
  }
}

const form = ref({
  ID: 0,
  path: '',
  name: '',
  hidden: false,
  parentId: 0,
  component: '',
  meta: {
    activeName: '',
    title: '',
    icon: '',
    defaultMenu: false,
    closeTab: false,
    keepAlive: false,
  },
  parameters: [],
  menuBtn: [],
})
const changeName = () => {
  form.value.path = form.value.name
}

const handleClose = (done) => {
  initForm()
  done()
}
// Xóa menu
const deleteMenu = (ID) => {
  ElMessageBox.confirm('Thao tác này sẽ xóa vĩnh viễn tất cả các menu dưới vai trò này, bạn có muốn tiếp tục?', 'Cảnh báo', {
    confirmButtonText: 'Đồng ý',
    cancelButtonText: 'Hủy',
    type: 'warning',
  })
    .then(async () => {
      const res = await deleteBaseMenu({ ID })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: 'Xóa thành công!',
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
        message: 'Đã hủy xóa',
      })
    })
}
// Khởi tạo các phương thức bảng trong hộp thoại
const menuForm = ref(null)
const checkFlag = ref(false)
const initForm = () => {
  checkFlag.value = false
  menuForm.value.resetFields()
  form.value = {
    ID: 0,
    path: '',
    name: '',
    hidden: false,
    parentId: 0,
    component: '',
    meta: {
      title: '',
      icon: '',
      defaultMenu: false,
      closeTab: false,
      keepAlive: false,
    },
  }
}
// Đóng hộp thoại

const dialogFormVisible = ref(false)
const closeDialog = () => {
  initForm()
  dialogFormVisible.value = false
}
// Thêm menu
const enterDialog = async () => {
  menuForm.value.validate(async (valid) => {
    if (valid) {
      let res
      if (isEdit.value) {
        res = await updateBaseMenu(form.value)
      } else {
        res = await addBaseMenu(form.value)
      }
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: isEdit.value ? 'Chỉnh sửa thành công' : 'Thêm thành công!',
        })
        getTableData()
      }
      initForm()
      dialogFormVisible.value = false
    }
  })
}

const menuOption = ref([
  {
    ID: '0',
    title: 'Menu gốc',
  },
])
const setOptions = () => {
  menuOption.value = [
    {
      ID: 0,
      title: 'Thư mục gốc',
    },
  ]
  setMenuOptions(tableData.value, menuOption.value, false)
}
const setMenuOptions = (menuData, optionsData, disabled) => {
  menuData &&
    menuData.forEach((item) => {
      if (item.children && item.children.length) {
        const option = {
          title: item.meta.title,
          ID: item.ID,
          disabled: disabled || item.ID === form.value.ID,
          children: [],
        }
        setMenuOptions(
          item.children,
          option.children,
          disabled || item.ID === form.value.ID
        )
        optionsData.push(option)
      } else {
        const option = {
          title: item.meta.title,
          ID: item.ID,
          disabled: disabled || item.ID === form.value.ID,
        }
        optionsData.push(option)
      }
    })
}

// Phương thức thêm menu, id = 0 là thêm menu gốc
const isEdit = ref(false)
const dialogTitle = ref('Thêm menu')
const addMenu = (id) => {
  dialogTitle.value = 'Thêm menu'
  form.value.parentId = id
  isEdit.value = false
  setOptions()
  dialogFormVisible.value = true
}
// Phương thức chỉnh sửa menu
const editMenu = async (id) => {
  dialogTitle.value = 'Chỉnh sửa menu'
  const res = await getBaseMenuById({ id })
  form.value = res.data.menu
  isEdit.value = true
  setOptions()
  dialogFormVisible.value = true
}
</script>

<style scoped lang="scss">
.warning {
  color: #dc143c;
}

.icon-column {
  display: flex;
  align-items: center;

  .el-icon {
    margin-right: 8px;
  }
}
</style>
