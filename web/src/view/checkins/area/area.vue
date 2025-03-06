<template>
  <div>
    <!-- Tìm kiếm -->
    <div class="gva-search-box hidden">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule" @keyup.enter="onSubmit">
        <el-form-item label="Ngày tạo" prop="createdAt">
          <template #label>
            <span>
              Ngày tạo
              <el-tooltip content="Phạm vi tìm kiếm từ ngày bắt đầu (bao gồm) đến ngày kết thúc (không bao gồm)">
                <el-icon>
                  <QuestionFilled />
                </el-icon>
              </el-tooltip>
            </span>
          </template>
          <el-date-picker
            v-model="searchInfo.startCreatedAt"
            type="datetime"
            placeholder="Ngày bắt đầu"
            :disabled-date="time => searchInfo.endCreatedAt ? time.getTime() > searchInfo.endCreatedAt.getTime() : false"
          ></el-date-picker>
          —
          <el-date-picker
            v-model="searchInfo.endCreatedAt"
            type="datetime"
            placeholder="Ngày kết thúc"
            :disabled-date="time => searchInfo.startCreatedAt ? time.getTime() < searchInfo.startCreatedAt.getTime() : false"
          ></el-date-picker>
        </el-form-item>
        <template v-if="showAllQuery">
          <!-- Thêm các điều kiện tìm kiếm nếu cần -->
        </template>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">Tìm kiếm</el-button>
          <el-button icon="refresh" @click="onReset">Đặt lại</el-button>
          <el-button link type="primary" icon="arrow-down" @click="showAllQuery = true" v-if="!showAllQuery">Mở rộng</el-button>
          <el-button link type="primary" icon="arrow-up" @click="showAllQuery = false" v-else>Thu gọn</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- Hiển thị bảng dữ liệu -->
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog">Thêm mới</el-button>
        <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="onDelete">Xóa</el-button>
      </div>
      <el-table ref="multipleTable" style="width: 100%" tooltip-effect="dark" :data="tableData" row-key="ID" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="Tên khu vực" prop="name" width="300" />
        <el-table-column align="left" label="Vĩ độ" prop="latitude" width="180" />
        <el-table-column align="left" label="Kinh độ" prop="longitude" width="180" />
        <el-table-column align="left" label="Bán kính" prop="radius" width="120" />
        <el-table-column align="left" label="Hành động" fixed="right" min-width="240">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)">
              <el-icon style="margin-right: 5px">
                <InfoFilled />
              </el-icon>
              Xem chi tiết
            </el-button>
            <el-button type="primary" link icon="edit" class="table-button" @click="updateAreaFunc(scope.row)">Chỉnh sửa</el-button>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">Xóa</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination layout="total, sizes, prev, pager, next, jumper" :current-page="page" :page-size="pageSize" :page-sizes="[10, 30, 50, 100]" :total="total" @current-change="handleCurrentChange" @size-change="handleSizeChange" />
      </div>
    </div>

    <!-- Dialog thêm/sửa (có map vẽ khu vực) -->
    <el-drawer destroy-on-close size="800" v-model="dialogFormVisible" :show-close="false" :before-close="closeDialog">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? 'Thêm mới' : 'Chỉnh sửa' }}</span>
          <div>
            <el-button type="primary" @click="enterDialog">Đồng ý</el-button>
            <el-button @click="closeDialog">Hủy</el-button>
          </div>
        </div>
      </template>
      <el-form :model="formData" label-position="top" ref="elFormRef" :rules="rule" label-width="80px">
        <el-form-item label="Tên khu vực:" prop="name">
          <el-input v-model="formData.name" clearable placeholder="Nhập tên khu vực" />
        </el-form-item>
        <el-form-item label="Vĩ độ:" prop="latitude">
          <el-input-number v-model="formData.latitude" style="width:100%" :precision="15" clearable />
        </el-form-item>
        <el-form-item label="Kinh độ:" prop="longitude">
          <el-input-number v-model="formData.longitude" style="width:100%" :precision="15" clearable />
        </el-form-item>
        <el-form-item label="Bán kính:" prop="radius">
          <el-input-number v-model="formData.radius" style="width:100%" :precision="2" clearable />
        </el-form-item>
        <el-form-item label="Giới hạn địa chỉ IP">
          <span class="text-sm my-1 italic font-normal">
            Nhập các IP cho phép, cách nhau bởi dấu phẩy, không có khoảng trắng
          </span>
          <el-input v-model="formData.restrictIp" style="width:100%"></el-input>
        </el-form-item>
        <!-- Container map vẽ khu vực -->
        <el-form-item label="Vẽ khu vực:">
          <div id="draw-map" style="width: 100%; height: 400px;"></div>
        </el-form-item>
      </el-form>
    </el-drawer>

    <!-- Drawer xem chi tiết (có map hiển thị khu vực) -->
    <el-drawer destroy-on-close size="800" v-model="detailShow" :show-close="true" :before-close="closeDetailShow">
      <el-descriptions column="1" border>
        <el-descriptions-item label="Tên khu vực">{{ detailFrom.name }}</el-descriptions-item>
        <el-descriptions-item label="Vĩ độ">{{ detailFrom.latitude }}</el-descriptions-item>
        <el-descriptions-item label="Kinh độ">{{ detailFrom.longitude }}</el-descriptions-item>
        <el-descriptions-item label="Bán kính">{{ detailFrom.radius }}</el-descriptions-item>
        <el-descriptions-item label="IP giới hạn">{{ detailFrom.restrictIp }}</el-descriptions-item>
      </el-descriptions>
      <!-- Container map hiển thị khu vực -->
      <div id="openlayers-map" style="width: 100%; height: 400px; margin-top: 20px;"></div>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, nextTick, watch } from 'vue'
import {
  createArea,
  deleteArea,
  deleteAreaByIds,
  updateArea,
  findArea,
  getAreaList
} from '@/api/checkins/area'
import { ElMessage, ElMessageBox } from 'element-plus'

// --- Các biến và cấu hình chung ---
defineOptions({ name: 'Area' })

// Biến cho tìm kiếm và phân trang
const showAllQuery = ref(false)
const searchInfo = ref({})
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])

// Dữ liệu form
const formData = ref({
  name: '',
  latitude: null,
  longitude: null,
  radius: 50,
  restrictIp: null
})

// Validation rules
const rule = reactive({
  name: [
    { required: true, message: 'Tên khu vực là bắt buộc', trigger: ['input', 'blur'] },
    { whitespace: true, message: 'Không được chỉ nhập khoảng trắng', trigger: ['input', 'blur'] }
  ]
})
const searchRule = reactive({
  createdAt: [
    {
      validator: (rule, value, callback) => {
        if (searchInfo.value.startCreatedAt && !searchInfo.value.endCreatedAt) {
          callback(new Error('Vui lòng nhập ngày kết thúc'))
        } else if (!searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt) {
          callback(new Error('Vui lòng nhập ngày bắt đầu'))
        } else if (
          searchInfo.value.startCreatedAt &&
          searchInfo.value.endCreatedAt &&
          (searchInfo.value.startCreatedAt.getTime() === searchInfo.value.endCreatedAt.getTime() ||
            searchInfo.value.startCreatedAt.getTime() > searchInfo.value.endCreatedAt.getTime())
        ) {
          callback(new Error('Ngày bắt đầu phải trước ngày kết thúc'))
        } else {
          callback()
        }
      },
      trigger: 'change'
    }
  ]
})
const elFormRef = ref()
const elSearchFormRef = ref()

// --- Hàm điều khiển bảng ---
const getTableData = async () => {
  const table = await getAreaList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}
getTableData()

const onReset = () => {
  searchInfo.value = {}
  getTableData()
}
const onSubmit = () => {
  elSearchFormRef.value?.validate(async (valid) => {
    if (!valid) return
    page.value = 1
    pageSize.value = 10
    getTableData()
  })
}
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// --- Các hàm thao tác dữ liệu ---
const multipleSelection = ref([])
const handleSelectionChange = (val) => {
  multipleSelection.value = val
}
const deleteRow = (row) => {
  ElMessageBox.confirm('Bạn có chắc muốn xóa không?', 'Cảnh báo', {
    confirmButtonText: 'Đồng ý',
    cancelButtonText: 'Hủy',
    type: 'warning'
  }).then(() => {
    deleteAreaFunc(row)
  })
}
const onDelete = async () => {
  ElMessageBox.confirm('Bạn có chắc muốn xóa không?', 'Cảnh báo', {
    confirmButtonText: 'Đồng ý',
    cancelButtonText: 'Hủy',
    type: 'warning'
  }).then(async () => {
    const IDs = []
    if (multipleSelection.value.length === 0) {
      ElMessage({ type: 'warning', message: 'Vui lòng chọn dữ liệu để xóa' })
      return
    }
    multipleSelection.value.forEach(item => {
      IDs.push(item.ID)
    })
    const res = await deleteAreaByIds({ IDs })
    if (res.code === 0) {
      ElMessage({ type: 'success', message: 'Xóa thành công' })
      if (tableData.value.length === IDs.length && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  })
}

// --- Dialog thêm/sửa ---
const type = ref('')
const dialogFormVisible = ref(false)
const openDialog = () => {
  type.value = 'create'
  dialogFormVisible.value = true
}
const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = { name: '', latitude: 0, longitude: 0, radius: 0, restrictIp: null }
}
const enterDialog = async () => {
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return
    let res
    if (type.value === 'create') {
      res = await createArea(formData.value)
    } else if (type.value === 'update') {
      res = await updateArea(formData.value)
    }
    if (res.code === 0) {
      ElMessage({ type: 'success', message: 'Tạo/cập nhật thành công' })
      closeDialog()
      getTableData()
    }
  })
}
const updateAreaFunc = async (row) => {
  const res = await findArea({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data
    dialogFormVisible.value = true
  }
}
const deleteAreaFunc = async (row) => {
  const res = await deleteArea({ ID: row.ID })
  if (res.code === 0) {
    ElMessage({ type: 'success', message: 'Xóa thành công' })
    if (tableData.value.length === 1 && page.value > 1) {
      page.value--
    }
    getTableData()
  }
}

// --- Drawer xem chi tiết ---
const detailFrom = ref({})
const detailShow = ref(false)
const openDetailShow = () => {
  detailShow.value = true
}
const getDetails = async (row) => {
  const res = await findArea({ ID: row.ID })
  if (res.code === 0) {
    detailFrom.value = res.data
    openDetailShow()
  }
}
const closeDetailShow = () => {
  detailShow.value = false
  detailFrom.value = {}
}

// --- OpenLayers tích hợp cho Drawer xem chi tiết ---
import 'ol/ol.css'
import Map from 'ol/Map'
import View from 'ol/View'
import TileLayer from 'ol/layer/Tile'
import OSM from 'ol/source/OSM'
import { fromLonLat } from 'ol/proj'
import { Circle as CircleGeom } from 'ol/geom'
import { Feature } from 'ol'
import VectorLayer from 'ol/layer/Vector'
import VectorSource from 'ol/source/Vector'
import { Style, Stroke, Fill } from 'ol/style'

let detailMapInstance = null
const initOpenLayersMap = () => {
  const mapContainer = document.getElementById('openlayers-map')
  if (!mapContainer) return
  const lat = Number(detailFrom.value.latitude)
  const lng = Number(detailFrom.value.longitude)
  const radius = Number(detailFrom.value.radius)
  if (isNaN(lat) || isNaN(lng)) return
  const center = fromLonLat([lng, lat])
  const circle = new CircleGeom(center, radius)
  const circleFeature = new Feature(circle)
  const vectorSource = new VectorSource({ features: [circleFeature] })
  const vectorLayer = new VectorLayer({
    source: vectorSource,
    style: new Style({
      stroke: new Stroke({ color: 'red', width: 2 }),
      fill: new Fill({ color: 'rgba(255, 0, 0, 0.3)' })
    })
  })
  if (detailMapInstance) {
    detailMapInstance.setTarget(null)
  }
  detailMapInstance = new Map({
    target: 'openlayers-map',
    layers: [new TileLayer({ source: new OSM() }), vectorLayer],
    view: new View({ center: center, zoom: 15 })
  })
  detailMapInstance.getView().fit(circle.getExtent(), { duration: 1000, padding: [20, 20, 20, 20] })
}
watch(detailShow, (newVal) => {
  if (newVal) {
    nextTick(() => {
      initOpenLayersMap()
    })
  }
})

// --- OpenLayers tích hợp cho map vẽ trong Dialog chỉnh sửa ---
import { Draw } from 'ol/interaction'
import { toLonLat } from 'ol/proj'

let drawingMapInstance = null
const initDrawingMap = () => {
  const drawMapContainer = document.getElementById('draw-map')
  if (!drawMapContainer) return
  // Tạo vector source và layer để chứa đối tượng vẽ
  const vectorSource = new VectorSource()
  const vectorLayer = new VectorLayer({
    source: vectorSource,
    style: new Style({
      stroke: new Stroke({ color: 'blue', width: 2 }),
      fill: new Fill({ color: 'rgba(0, 0, 255, 0.3)' })
    })
  })
  drawingMapInstance = new Map({
    target: 'draw-map',
    layers: [new TileLayer({ source: new OSM() }), vectorLayer],
    view: new View({
      center: formData.value.longitude && formData.value.latitude
        ? fromLonLat([formData.value.longitude, formData.value.latitude])
        : fromLonLat([105.83416, 21.027764]),
      zoom: 12
    })
  })
  const draw = new Draw({ source: vectorSource, type: 'Circle' })
  drawingMapInstance.addInteraction(draw)
  // Khi bắt đầu vẽ, xoá feature cũ để không tồn tại nhiều hình cùng lúc
  draw.on('drawstart', function() {
    vectorSource.clear();
  })
  // Khi vẽ xong, cập nhật formData với tọa độ và bán kính
  draw.on('drawend', function(event) {
    const circleGeom = event.feature.getGeometry()
    const center = circleGeom.getCenter()
    const radius = circleGeom.getRadius()
    const centerLonLat = toLonLat(center)
    formData.value.latitude = centerLonLat[1]
    formData.value.longitude = centerLonLat[0]
    formData.value.radius = radius
  })
}
watch(dialogFormVisible, (newVal) => {
  if (newVal) {
    nextTick(() => {
      initDrawingMap()
    })
  }
})
</script>

<style>
/* Thêm style tùy chỉnh nếu cần */
</style>
