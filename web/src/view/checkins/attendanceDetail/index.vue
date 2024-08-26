<template>
  <div class="p-2">
    <div class="p-1 my-1 flex">
      <ImportExcel :form-data="{ action: 'IMPORT_PARTICIPANT' }" />
      <el-button class="mx-4 hidden" type="danger" icon="download"> Xuất Excel</el-button>
    </div>
    <el-tabs v-model="tabsActiveTab" type="border-card">
      <el-tab-pane name="attendanceInfoTab" label="Chi tiết">
        <div class="card-container">
          <el-form ref="elFormRef" :model="formData" label-position="top" :rules="rule" label-width="80px">
            <el-row>
              <el-col :span="12" class="grid-cell">
                <el-form-item label="Tiêu đề" prop="formData.title" class="required">
                  <el-input v-model="formData.title" type="text" clearable />
                </el-form-item>
                <el-form-item label="Client URL" prop="formData.clientUrl" class="required">
                  <el-input v-model="formData.clientUrl" type="text" clearable>
                    <template #prepend>https://</template>
                  </el-input>
                </el-form-item>
                <div class=" flex justify-between">
                  <el-form-item label="Ngày bắt đầu" label-width="150px" prop="startDate" class="required">
                    <el-date-picker v-model="formData.startDate" type="datetime" class="full-width-input" clearable />
                  </el-form-item>
                  <el-form-item label="Ngày kết thúc" label-width="150px" prop="endDate" class="required">
                    <el-date-picker v-model="formData.endDate" type="datetime" class="full-width-input" clearable />
                  </el-form-item>
                </div>
                <div class="flex justify-between">
                  <el-form-item label="Cho thử nghiệm" label-width="150px" prop="isTrial">
                    <el-switch v-model="formData.isTrial" />
                  </el-form-item>
                  <el-form-item label="Khoá" label-width="150px" prop="isLocked">
                    <el-switch v-model="formData.isLocked" />
                  </el-form-item>
                </div>

              </el-col>
              <el-col :span="12" class="grid-cell flex-column px-4">
                <div class="flex justify-end p-2">
                  <el-button type="success" @click="downloadQRCode">Tải xuống QR Code</el-button>
                  <el-button type="primary" @click="saveAttendance">Lưu</el-button>
                </div>
                <div class="flex flex-col items-center justify-center w-full">
                  <canvas class="border-2 rounded border-gray-500" ref="qrcodeCanvas"></canvas>
                  <div>Dùng QR Code này để quét điểm danh</div>
                </div>
              </el-col>
            </el-row>

          </el-form>
        </div>

        <el-divider />
        <div class="text-xl">Danh sách điểm danh</div>
        <div class="mt-4">
          <el-table style="width: 100%" tooltip-effect="dark" :data="tableData" row-key="ID">
            <el-table-column align="left" label="Ngày giờ" prop="checkinDate" width="180">
              <template #default="scope">{{ formatDateTime(scope.row.checkinDate) }}</template>
            </el-table-column>
            <el-table-column align="left" label="Thành viên" prop="participant.fullName" width="200" />
            <el-table-column align="left" label="Khu vực" width="120">
              <template #default="scope">
                <span>{{ scope.row.area?.name ?? '/' }}</span>
              </template>
            </el-table-column>
            <el-table-column align="left" label="Nhóm" width="120">
              <template #default="scope">
                <span>{{ scope.row.group.name ?? '/' }}</span>
              </template>
            </el-table-column>
            <el-table-column align="left" label="IP" prop="iP" width="90" />
            <el-table-column align="left" label="Kinh độ" prop="lattidue" width="170">
              <template #default="scope">
                <span>{{ '11.953687161569116' }}</span>
              </template>
            </el-table-column>
            <el-table-column align="left" label="Vĩ độ" prop="longtidue" width="170">
              <template #default="scope">
                <span>{{ '11.953687161569116' }}</span>
              </template>
            </el-table-column>
            <el-table-column class="overflow-hidden" label="Agent" prop="agent"></el-table-column>
          </el-table>
          <div class="gva-pagination">
            <el-pagination layout="total, sizes, prev, pager, next, jumper" :current-page="page" :page-size="pageSize"
              :page-sizes="[10, 30, 50, 100]" :total="total" @current-change="handleCurrentChange"
              @size-change="handleSizeChange" />
          </div>
        </div>
        <div class="table-container">
          <table class="table-layout">
            <tbody>
              <tr>
                <td class="table-cell" />
              </tr>
            </tbody>
          </table>
        </div>
      </el-tab-pane>
      <el-tab-pane name="partticipantsTab" label="Thành viên">
        <div class="table-container">
          <Partticipant :ac-id="currentId" />
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
          <Condition :ac-id="currentId" />
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

import { useRoute } from 'vue-router';
import QRCode from 'qrcode'
import { ElForm, ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'
import Partticipant from '@/view/checkins/components/participant/index.vue'
import Group from '@/view/checkins/components/group/index.vue'
import Area from '@/view/checkins/components/area/index.vue'
import Condition from '@/view/checkins/components/condition/index.vue'
import ImportExcel from '@/components/importExcel/index.vue'
import { formatDateTime } from '@/utils/format'
import base32 from 'hi-base32'

defineOptions({
  name: 'AttendanceDetail'
})

const $route = useRoute()
const tabsActiveTab = ref('attendanceInfoTab')
const currentId = ref($route.params.id)
const formData = ref({})
const elFormRef = ref();

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

const getDetailData = async () => {
  var id = $route.params.id
  const res = await findAttendance({ id: $route.params.id })
  if (res.code == 0) {
    formData.value = res.data
  }
  generateQRCode();
}
getDetailData();


const saveAttendance = async () => {
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

const qrcodeCanvas = ref(null)

const generateQRCode = async () => {
  var params = base32.encode($route.params.id)
  console.log('params-endcode', params)
  var url = formData.value.clientUrl + '/checkin/?c=' + params
  QRCode.toCanvas(qrcodeCanvas.value, url, { width: 300 }, (error) => {
    if (error) console.error(error)
  })
}

const downloadQRCode = () => {
  const canvas = qrcodeCanvas.value
  const url = canvas.toDataURL('image/png')
  const link = document.createElement('a')
  link.href = url
  link.download = 'QR Điểm danh số - ' + $route.params.id + '.png'
  link.click()
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
  },
  ],
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
}

getTableData()
</script>

<style lang="scss"></style>