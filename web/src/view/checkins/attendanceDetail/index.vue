<template>
  <div class="p-2">
    <div class="p-1 my-1 flex justify-end w-100">
      <ImportExcel :form-data="{ action: 'IMPORT_PARTICIPANT', attendanceId: currentId }" class="mb-2" />
      <el-button class="mx-4 hidden" type="danger" icon="download">
        Xuất Excel
      </el-button>
    </div>
    <el-tabs v-model="tabsActiveTab" type="border-card" @tab-click="tabHandleClick">
      <el-tab-pane name="attendanceInfoTab" label="Chi tiết">
        <div class="card-container">
          <el-form ref="elFormRef" :model="formData" label-position="top" :rules="rule" label-width="80px">
            <el-row>
              <el-col :span="12" class="grid-cell">
                <el-form-item label="Tiêu đề" prop="formData.title" class=" required">
                  <el-input v-model="formData.title" type="text" clearable />
                </el-form-item>
                <el-form-item label="URL Điểm danh" prop="formData.clientUrl" class="required">
                  <el-input v-model="formData.clientUrl" type="text" clearable disabled />
                </el-form-item>
                <div class=" flex justify-between">
                  <el-form-item label="Ngày bắt đầu" label-width="150px" prop="startDate" class="required">
                    <el-date-picker v-model="formData.startDate" type="datetime" class="full-width-input" clearable />
                  </el-form-item>
                  <el-form-item class=" required" label="Ngày kết thúc" label-width="150px" prop="endDate">
                    <el-date-picker v-model="formData.endDate" type="datetime" class="full-width-input" clearable />
                  </el-form-item>
                </div>

                <el-form-item label="Danh mục" prop="categoryId" class="w-full required">
                  <el-tree-select class="w-full" v-model="formData.categoryId" :data="categoryOptions"
                    check-on-click-node :render-after-expand="false" style="width: 240px" />
                </el-form-item>
                <el-form-item label="Đơn vị" prop="agencyId" class="w-full required">
                  <el-select v-model="formData.agencyId" placeholder="Chọn đơn vị" clearable filterable>
                    <el-option v-for="item in agencyOptions" :key="item.ID" :label="item.name" :value="item.ID" />
                  </el-select>
                </el-form-item>
                <el-form-item prop="description" label="Mô tả" class="w-full">
                  <el-input type="text" v-model="formData.description"></el-input>
                </el-form-item>
                <div class="flex justify-between">
                  <el-form-item label="Đóng điểm danh" label-width="150px" prop="isLocked">
                    <el-switch v-model="formData.isLocked" />
                  </el-form-item>
                  <el-form-item label="Cho phép khách" label-width="150px" prop="allowGuest">
                    <el-switch v-model="formData.allowGuest" />
                  </el-form-item>
                </div>
                <el-button link type="primary" icon="arrow-down" @click="showAllOptionConfig = true"
                  v-if="!showAllOptionConfig">Mở
                  rộng</el-button>
                <el-button link type="primary" icon="arrow-up" @click="showAllOptionConfig = false" v-else>Thu
                  gọn</el-button>
                <template v-if="showAllOptionConfig">
                  <div class="flex justify-between">
                    <el-form-item label="Hệ số" prop="formData.weight">
                      <el-input v-model="formData.weight" type="number" clearable />
                    </el-form-item>
                    <el-form-item label="Giới hạn số lần / thành viên" prop="formData.limitCount">
                      <el-input v-model="formData.limitCount" type="number" clearable />
                    </el-form-item>
                    <el-form-item label="Giới hạn số lần / máy" prop="formData.limitClientCount">
                      <el-input v-model="formData.limitClientCount" type="number" clearable />
                    </el-form-item>
                  </div>

                  <el-form-item label="Giới hạn IP truy cập" prop="formData.restrictIp">
                    <el-input v-model="formData.restrictIp" type="text" clearable
                      placeholder="172.0.0.1,196.0.0.1,10.0.0.0/32" />
                    <span class="text-sm my-1 italic font-normal">Để giới hạn các IP điểm danh, nhập các IP được cho
                      phép
                      vào ô dưới đây, cách nhau bởi dấu phẩy, không có khoảng trắng</span>
                  </el-form-item>
                  <el-form-item label="URL chuyển hướng" prop="formData.clientUrl">
                    <el-input v-model="formData.redirectUrl" type="text" clearable />
                    <span class="text-sm my-1 italic font-normal">Hệ thống sẽ chuyển hướng bạn đến địa chỉ được nhập vào
                      sau khi điểm danh</span>
                  </el-form-item>
                </template>

              </el-col>
              <el-col :span="12" class="grid-cell flex-column px-4">
                <div class="flex justify-end p-2">
                  <!-- <el-button type="success" @click="downloadQRCode">
                    Tải xuống QR Code
                  </el-button> -->
                  <el-button type="primary" @click="saveAttendance">
                    Lưu
                  </el-button>
                </div>
                <div class="flex flex-col items-center justify-center w-full">
                  <!-- <canvas ref="qrcodeCanvas" class="border-2 rounded border-gray-500" /> -->
                  <QRCodeVue3 :width="512" :height="512" :value="targetClientURL" :key="targetClientURL" :qr-options="{
                    errorCorrectionLevel: 'H'
                  }" :image-options="{ hideBackgroundDots: true, imageSize: 0.4, margin: 10 }"
                    :corners-square-options="{ type: 'dot', color: '#514C39' }" :corners-dot-options="{
                      type: undefined,
                      color: '#7BA227'
                    }" :dots-options="{
                      type: 'dots',
                      color: '#7BA227',
                      gradient: {
                        type: 'linear',
                        rotation: 0,
                        colorStops: [
                          { offset: 0, color: '#7BA227' },
                          { offset: 1, color: '#E67F32' }
                        ]
                      }
                    }" :download="true" image="/dlu.png" buttonName="Tải xuống" 
                    :downloadOptions="downloadQrOptions" downloadButton="download-btn" />
                  <!-- <div>Dùng QR Code này để quét điểm danh</div> -->
                </div>
              </el-col>
            </el-row>
          </el-form>
        </div>

        <el-divider />
      </el-tab-pane>
      <el-tab-pane name="partticipantsTab" label="Thành viên">
        <div class="table-container">
          <Partticipant :ac-id="currentId" :group-options="groupOptions" />
        </div>
      </el-tab-pane>
      <el-tab-pane name="groupTab" label="Nhóm">
        <div class="table-container">
          <Group :ac-id="currentId" />
        </div>
      </el-tab-pane>
      <el-tab-pane name="areaTab" label="Khu vực">
        <div class="table-container">
          <Area :ac-id="currentId" />
        </div>
      </el-tab-pane>
      <el-tab-pane name="conditionTab" label="Điều kiện">
        <div class="table-container">
          <Condition :ref="conditionTabRef" :ac-id="currentId" />
        </div>
      </el-tab-pane>
      <el-tab-pane name="histories" label="Lịch sử">
        <div class="table-container">
          <div class="text-xl">
            Danh sách điểm danh
          </div>
          <div class="flex flex-row justify-evenly py-2">
            <div class="bg-white shadow-md text-base p-4 rounded justify-between items-center">
              <p class="text-5xl font-bold text-center py-3">{{ formData.total }}</p>
              <p class="text-slate-700">Tổng số thành viên</p>
            </div>
            <div class="bg-white shadow-md text-base p-4 rounded justify-between items-center">
              <p class="text-5xl font-bold text-center py-3">{{ formData.totalCheckin }}</p>
              <p class="text-slate-700">Đã điểm danh</p>
            </div>
            <div class="bg-white shadow-md text-base p-4 rounded justify-between items-center">
              <p class="text-5xl font-bold text-center py-3">{{ formData.total - formData.totalCheckin }}</p>
              <p class="text-slate-700">Chưa điểm danh</p>
            </div>
            <div class="bg-white shadow-md text-base p-4 rounded justify-between items-center">
              <p class="text-5xl font-bold text-center py-3">{{ total }}</p>
              <p class="text-slate-700">Tổng lượt truy cập</p>
            </div>
          </div>
          <div class="my-4">
            <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline"
              :rules="searchRules" @keyup.enter="onSubmit">
              <el-form-item label="Ngày tạo" prop="createdAt">
                <el-date-picker v-model="searchInfo.startCreatedAt" type="date" placeholder="Ngày bắt đầu" />
              </el-form-item>
              <el-form-item label="Email" prop="email">
                <el-input v-model="searchInfo.email" type="text" placeholder="Email" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" icon="search" @click="onSubmit">
                  Tìm kiếm
                </el-button>
                <el-button icon="refresh" @click="onReset">
                  Đặt lại
                </el-button>
              </el-form-item>
            </el-form>
          </div>
          <div class="mt-4">
            <el-table style="width: 100%" tooltip-effect="dark" :data="tableData" row-key="ID" border>
              <el-table-column align="left" label="Ngày giờ" prop="checkinDate" width="180">
                <template #default="scope">
                  {{ formatDateTime(scope.row.checkinDate) }}
                </template>
              </el-table-column>
              <el-table-column align="left" label="Thành viên" prop="participant.fullName" width="200" />
              <el-table-column align="left" label="Khu vực" width="120">
                <template #default="scope">
                  <span>{{ scope.row.area?.name ?? '/' }}</span>
                </template>
              </el-table-column>
              <el-table-column align="left" label="Nhóm" width="120">
                <template #default="scope">
                  <span>{{ scope.row.group?.name ?? '/' }}</span>
                </template>
              </el-table-column>

              <el-table-column align="left" label="IP" prop="iP" width="90" />
              <el-table-column align="left" label="Vị trí" width="170">
                <template #default="scope">
                  <a target="_blank"
                    :href="'https://www.google.com/maps?q=' + scope.row.lattidue + ',' + scope.row.longtidue">{{
                      scope.row.lattidue }}, {{ scope.row.longtidue }}</a>
                </template>
              </el-table-column>
              <el-table-column align="left" label="Client ID" prop="visitorId" width="200" />
              <el-table-column class="overflow-hidden" label="Agent" prop="agent" />
            </el-table>
            <div class="gva-pagination">
              <el-pagination layout="total, sizes, prev, pager, next, jumper" :current-page="page" :page-size="pageSize"
                :page-sizes="[10, 30, 50, 100]" :total="total" @current-change="handleCurrentChange"
                @size-change="handleSizeChange" />
            </div>
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import {
  updateAttendance,
  findAttendance,
} from '@/api/checkins/attendance'
import {
  getAttendanceCheckInList
} from '@/api/checkins/attendanceCheckIn'

import {
  findAttendanceArea
} from '@/api/checkins/attendance'

import {
  getGroupList
} from '@/api/checkins/group'

import { useRoute } from 'vue-router';
import QRCode from 'qrcode'
import { ElForm, ElMessage } from 'element-plus'
import { ref, reactive, nextTick } from 'vue'
import Partticipant from '@/view/checkins/components/participant/index.vue'
import Group from '@/view/checkins/components/group/index.vue'
import Area from '@/view/checkins/components/area/index.vue'
import Condition from '@/view/checkins/components/condition/index.vue'
import ImportExcel from '@/components/importExcel/index.vue'
import { formatDateTime, formatDate } from '@/utils/format'
import base32 from 'hi-base32'
import * as echarts from 'echarts'

import {
  getAttendanceCategoryList
} from '@/api/checkins/attendanceCategory'

import {
  getAttendanceAgencyList
} from '@/api/checkins/attendanceAgency'

import QRCodeVue3 from 'qrcode-vue3'

defineOptions({
  name: 'AttendanceDetail'
})

const $route = useRoute()
const tabsActiveTab = ref('attendanceInfoTab')
const currentId = ref($route.params.id)

const clientURL = ref(import.meta.env.VITE_CLIENT_URL)

const formData = ref({
  isLocked: false,
  allowGuest: false,
  title: '',
  categoryId: null,
  agencyId: null,

})

const elFormRef = ref();
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const elSearchFormRef = ref()
const categoryOptions = ref([])
const agencyOptions = ref([])

const showAllOptionConfig = ref(false)
const conditionTabRef = ref()

const searchRules = reactive({
  createdAt: [
    {
      validator: (rule, value, callback) => {
        if (!searchInfo.value.startCreatedAt) {
          callback(new Error('Vui lòng nhập ngày kết thúc'))
        } else {
          callback()
        }
      }, trigger: 'change'
    }
  ],
  email: [
    {
      validator: (rule, value, callback) => {
        if (!searchInfo.value.email) {
          callback(new Error('Vui lòng nhập email'))
        } else {
          callback()
        }
      }, trigger: 'change'
    }
  ],
})

const getDetailData = async () => {
  var id = $route.params.id
  console.log('client URL', clientURL.value)
  const res = await findAttendance({ id: $route.params.id })
  if (res.code == 0) {
    formData.value = res.data
  }
  generateQRCode();
  // console.log("getDetailData: ", formData.value)

}
getDetailData();
//TODO: thay đổi cấu trúc nhận props của các component 
/**
 * Đưa các API lấy dữ liệu Option từ bên index
 * Truyền thông qua props của component
 * Emit event ra ngoài để bên index lấy dữ liệu
 */

const saveAttendance = async () => {
  if (formData.value.limitCount) {
    formData.value.limitCount = Number(formData.value.limitCount)
  }
  if (formData.value.limitClientCount) {
    formData.value.limitClientCount = Number(formData.value.limitClientCount)
  }

  if (formData.value.weight) {
    formData.value.weight = Number(formData.value.weight)
  }

  elFormRef.value?.validate(async (valid) => {
    if (!valid) return
    var res = await updateAttendance(formData.value)
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: 'Tạo/cập nhật thành công'
      })
      getDetailData()
    }
  })
}

const targetClientURL = ref("")
const downloadQrOptions = ref({})
const generateQRCode = async () => {
  var params = base32.encode($route.params.id)
  console.log('params-endcode' + params)
  console.log(clientURL.value)
  var url = clientURL.value + '/?c=' + params
  formData.value.clientUrl = url
  console.log('url', url)
  targetClientURL.value = url
  downloadQrOptions.value = {
    name: 'QR Điểm danh - ' + formData.value.title + " ngày " + formatDate(formData.value.startDate),
    extension: 'png'
  }

}

const downloadQRCode = async () => {
  const canvas = qrcodeCanvas.value
  const url = canvas.toDataURL('image/png')
  const link = document.createElement('a')
  link.href = url
  link.download = 'QR Điểm danh - ' + formData.value.title + '.png'
  link.click()
  await saveAttendance()
}

const rule = reactive({
  title: [{
    required: true,
    message: 'Tiêu đề không được để trống',
    trigger: ['input', 'blur'],
  },
  {
    whitespace: true,
    message: 'Bạn không thể chỉ nhập dấu cách',
    trigger: ['input', 'blur'],
  }
  ],
  startDate: [{
    required: true,
    message: 'Ngày bắt đầu là bắt buộc',
    trigger: ['input', 'blur'],
  },
  ],
  endDate: [{
    required: true,
    message: 'Ngày kết thúc là bắt buộc',
    trigger: ['input', 'blur'],
  }],
  categoryId: [{
    required: true,
    message: 'Danh mục không được để trống',
    trigger: ['input', 'blur'],
  }],
  agencyId: [{
    required: true,
    message: 'Đơn vị không được để trống',
    trigger: ['input', 'blur'],
  }],

})


const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}


const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}


const getTableData = async () => {
  searchInfo.value.attendanceId = $route.params.id
  const table = await getAttendanceCheckInList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
  // console.log("Danh sách điểm danh")
  // console.log(table)
  generateQRCode()
}

getTableData()

const onSubmit = () => {
  page.value = 1
  pageSize.value = 10
  getTableData()
  // elSearchFormRef.value?.validate(async (valid) => {
  //   if (!valid) return
  //   page.value = 1
  //   pageSize.value = 10
  //   getTableData()
  // })

}

const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

const getCategoryOptions = async () => {
  const table = await getAttendanceCategoryList({ page: 0, pageSize: -1 })
  if (table.code === 0) {
    categoryOptions.value = convertToTree(table.data.list)
  }
  // console.log("parent Options", categoryOptions.value)
}

getCategoryOptions()

const getAgencyOptions = async () => {
  const table = await getAttendanceAgencyList({ page: 0, pageSize: -1 })
  if (table.code === 0) {
    agencyOptions.value = table.data.list
  }
}
getAgencyOptions()


const convertToTree = (data) => {
  const map = {}
  const roots = []

  // Create a map of nodes using their ID as the key
  data.forEach((node) => {
    map[node.ID] = { ...node, value: node.ID, label: node.name, children: [] }
  })

  // Iterate over the nodes and assign children to their parent
  data.forEach((node) => {
    const parent = map[node.parentId]

    if (parent) {
      parent.children.push(map[node.ID])
    } else {
      roots.push(map[node.ID])
    }
  })
  console.log('roots', roots)
  return roots
}

const tabHandleClick = async (tab, event) => {

  if (tab.props.name === 'conditionTab') {
    await conditionTabRef.value?.getAgencyOptions()
    await conditionTabRef.value?.getGroupOptions()
  }
}

const areaOptions = ref([])
const getAreaListData = async () => {
  const table = await findAttendanceArea({ id: searchInfo.value.attendanceId })
  if (table.code === 0) {
    areaOptions.value = table.data.map(item => {
      return {
        ID: item.ID,
        name: item.area?.name
      }
    })
  }

}
getAreaListData();

const groupOptions = ref([])
const getGroupOptions = async () => {
  const table = await getGroupList({ page: 1, pageSize: -1, ...searchInfo.value })
  if (table.code === 0) {
    groupOptions.value = table.data.list
  }
  console.log('groupOptions', groupOptions.value)
}
getGroupOptions();



</script>

<style lang="scss">
button.download-btn {
    padding: 12px 32px;
    border: none;
    justify-content: center;
    display: block;
    width: 100%;
    border-radius: 4px;
    background-color: #67C239;
    color: #fff;
    font-weight: 400;
    font-size: 14px;
}
</style>