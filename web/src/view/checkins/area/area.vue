<template>
  <div>
    <!-- TÌM KIẾM -->
    <div class="gva-search-box">
      <el-form 
        ref="elSearchFormRef" 
        :inline="true" 
        :model="searchInfo" 
        class="demo-form-inline"
        :rules="searchRule"
        @keyup.enter="onSubmit"
      >
        <el-form-item label="Khu vực" prop="areaId">
          <el-select v-model="searchInfo.areaId" placeholder="Chọn khu vực" clearable filterable>
            <el-option v-for="item in tableData" :key="item.ID" :label="item.name" :value="item.ID" />
          </el-select>
        </el-form-item>
        <template v-if="showAllQuery">
          <el-form-item label="Ngày tạo" prop="createdAt">
            <template #label>
              <span>
                Ngày tạo
                <el-tooltip content="Phạm vi tìm kiếm từ ngày bắt đầu (bao gồm) đến ngày kết thúc (không bao gồm)">
                  <el-icon><QuestionFilled /></el-icon>
                </el-tooltip>
              </span>
            </template>
            <el-date-picker 
              v-model="searchInfo.startCreatedAt" 
              type="datetime" 
              placeholder="Ngày bắt đầu"
              :disabled-date="time => searchInfo.endCreatedAt ? time.getTime() > searchInfo.endCreatedAt.getTime() : false" 
            />
            —
            <el-date-picker 
              v-model="searchInfo.endCreatedAt" 
              type="datetime" 
              placeholder="Ngày kết thúc"
              :disabled-date="time => searchInfo.startCreatedAt ? time.getTime() < searchInfo.startCreatedAt.getTime() : false" 
            />
          </el-form-item>        
        </template>
      </el-form>
    </div>
    <!-- DỮ LIỆU -->
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog">Thêm mới</el-button>
        <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="onDelete"/>
      </div>
      <el-table 
        ref="multipleTable" 
        style="width: 100%" 
        tooltip-effect="dark" 
        :data="filteredData" 
        row-key="ID"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="Tên khu vực" prop="name" width="300" />
        <el-table-column align="left" label="Vĩ độ" prop="latitude" width="180" />
        <el-table-column align="left" label="Kinh độ" prop="longitude" width="180" />
        <el-table-column align="left" label="Bán kính" prop="radius" width="120" />
        <el-table-column align="left" label="Hành động" fixed="right" min-width="240">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)">
              <el-icon style="margin-right: 5px"><InfoFilled /></el-icon>
              Xem chi tiết
            </el-button>
            <el-button type="primary" link icon="edit" class="table-button" @click="updateAreaFunc(scope.row)">Chỉnh sửa</el-button>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">Xóa</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination 
          layout="total, sizes, prev, pager, next, jumper" 
          :current-page="page" 
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]" 
          :total="total" 
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>
    <!-- Cửa sổ pop-up Thêm/Sửa -->
    <el-drawer 
      v-model="dialogFormVisible" 
      destroy-on-close size="800"  
      :show-close="false" 
      :before-close="closeDialog"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? 'Thêm mới' : 'Chỉnh sửa' }}</span>
          <div>
            <el-button type="primary" @click="enterDialog">Đồng ý</el-button>
            <el-button @click="closeDialog">Hủy</el-button>
          </div>
        </div>
      </template>
      <el-form ref="elFormRef" :model="formData" label-position="top" :rules="rule" label-width="80px">
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
            Để giới hạn các IP điểm danh, nhập các IP được cho phép vào ô dưới đây, cách nhau bởi dấu phẩy, không có khoảng trắng
          </span>
          <el-input v-model="formData.restrictIp" style="width:100%" />
        </el-form-item>
        <el-form-item label="Vẽ khu vực:">
          <el-button type="danger" @click="clearGeometry">Xóa toàn bộ vùng vẽ</el-button>
          <!-- <el-button type="primary" @click="enableModifyMode" style="margin-left: 10px;">Chỉnh sửa vùng vẽ</el-button> -->
          <div id="draw-map" style="width: 100%; height: 400px;"></div>
        </el-form-item>
      </el-form>
    </el-drawer>
    <!-- Cửa sổ chi tiết -->
    <el-drawer v-model="detailShow" destroy-on-close size="800" :show-close="true" :before-close="closeDetailShow">
      <el-descriptions column="1" border>
        <el-descriptions-item label="Tên khu vực">{{ detailFrom.name }}</el-descriptions-item>
        <el-descriptions-item label="Vĩ độ">{{ detailFrom.latitude }}</el-descriptions-item>
        <el-descriptions-item label="Kinh độ">{{ detailFrom.longitude }}</el-descriptions-item>
        <el-descriptions-item label="Bán kính" v-if="!detailFrom.geometry">{{ detailFrom.radius }}</el-descriptions-item>
        <el-descriptions-item label="Vùng">{{ detailFrom.geometry }}</el-descriptions-item>
        <el-descriptions-item label="IP giới hạn">{{ detailFrom.restrictIp }}</el-descriptions-item>
      </el-descriptions>
      <div id="openlayers-map" style="width: 100%; height: 400px; margin-top: 20px;"></div>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createArea,
  deleteArea,
  deleteAreaByIds,
  updateArea,
  findArea,
  getAreaList
} from '@/api/checkins/area'
import { getDictFunc, formatDate, formatBoolean, filterDict, filterDataSource, returnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, nextTick, computed, watch } from 'vue'
import 'ol/ol.css'
import Map from 'ol/Map'
import View from 'ol/View'
import TileLayer from 'ol/layer/Tile'
import OSM from 'ol/source/OSM'
import { fromLonLat, toLonLat } from 'ol/proj'
import { Circle as CircleGeom } from 'ol/geom'
import { Feature } from 'ol'
import VectorLayer from 'ol/layer/Vector'
import VectorSource from 'ol/source/Vector'
import { Style, Stroke, Fill } from 'ol/style'
import CircleStyle from 'ol/style/Circle'
import { Draw, Modify } from 'ol/interaction'
import { GeoJSON } from 'ol/format'
import Snap from 'ol/interaction/Snap'


defineOptions({ name: 'AreaArea' })

// Hàm fix GeoJSON như cũ
function fixGeoJsonStr(geoStr) {
  geoStr = geoStr.replace(/'/g, '"').trim();
  geoStr = geoStr.replace(/,\s*([\]\}])/g, '$1');
  if (!geoStr.endsWith('}')) {
    geoStr += '}';
  }
  if (geoStr.indexOf('"coordinates"') !== -1) {
    while (!geoStr.trim().endsWith(']]}')) {
      geoStr += ']';
    }
  }
  return geoStr;
}

// ==================== Dữ liệu form và validation ====================
const formData = ref({
  name: '',
  latitude: null,
  longitude: null,
  radius: 50,
  restrictIp: null,
  geometry: null 
})

const rule = reactive({
  name: [
    { required: true, message: 'Tên khu vực là bắt buộc ', trigger: ['input', 'blur'] },
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
        } else if (searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt &&
          (searchInfo.value.startCreatedAt.getTime() === searchInfo.value.endCreatedAt.getTime() ||
           searchInfo.value.startCreatedAt.getTime() > searchInfo.value.endCreatedAt.getTime())) {
          callback(new Error('Ngày bắt đầu phải trước ngày kết thúc'))
        } else {
          callback()
        }
      }, trigger: 'change'
    }
  ]
})

const elFormRef = ref()
const elSearchFormRef = ref()

// ==================== Bảng điều khiển ====================
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({ areaId: null, startCreatedAt: null, endCreatedAt: null })
const showAllQuery = ref(false)

watch(searchInfo, () => { getTableData() }, { deep: true })

const onSubmit = () => { getTableData() }
const handleSizeChange = (val) => { pageSize.value = val; getTableData() }
const handleCurrentChange = (val) => { page.value = val; getTableData() }

const getTableData = async () => {
  const response = await getAreaList({ page: page.value, pageSize: pageSize.value })
  if (response.code === 0) {
    tableData.value = response.data.list
    total.value = response.data.total
  }
}
getTableData()

const filteredData = computed(() => {
  return tableData.value.filter(item => {
    return (
      (!searchInfo.value.areaId || item.ID === searchInfo.value.areaId) &&
      (!searchInfo.value.startCreatedAt || new Date(item.createdAt) >= new Date(searchInfo.value.startCreatedAt)) &&
      (!searchInfo.value.endCreatedAt || new Date(item.createdAt) < new Date(searchInfo.value.endCreatedAt))
    )
  })
})

// ==================== Đa chọn và các thao tác bảng ====================
const multipleSelection = ref([])
const handleSelectionChange = (val) => { multipleSelection.value = val }

const deleteRow = (row) => {
  ElMessageBox.confirm('Bạn có chắc muốn xóa không?', 'Cảnh báo', {
    confirmButtonText: 'Đồng ý',
    cancelButtonText: 'Hủy',
    type: 'warning'
  }).then(() => { deleteAreaFunc(row) })
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
    multipleSelection.value.forEach(item => { IDs.push(item.ID) })
    const res = await deleteAreaByIds({ IDs })
    if (res.code === 0) {
      ElMessage({ type: 'success', message: 'Xóa thành công' })
      if (tableData.value.length === IDs.length && page.value > 1) { page.value-- }
      getTableData()
    }
  })
}

// ==================== Cửa sổ pop-up Thêm/Sửa ====================
const type = ref('')
const dialogFormVisible = ref(false)

const updateAreaFunc = async (row) => {
  const res = await findArea({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = { ...res.data }
    dialogFormVisible.value = true
  }
}

const deleteAreaFunc = async (row) => {
  const res = await deleteArea({ ID: row.ID })
  if (res.code === 0) {
    ElMessage({ type: 'success', message: 'Xóa thành công' })
    if (tableData.value.length === 1 && page.value > 1) { page.value-- }
    getTableData()
  }
}

const openDialog = () => {
  type.value = 'create'
  dialogFormVisible.value = true
}

const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = { name: '', latitude: 0, longitude: 0, radius: 0, restrictIp: null, geometry: null }
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

// ==================== Cửa sổ chi tiết ====================
const detailFrom = ref({})
const detailShow = ref(false)

const openDetailShow = () => { detailShow.value = true }

const getDetails = async (row) => {
  const res = await findArea({ ID: row.ID })
  if (res.code === 0) {
    detailFrom.value = res.data
    openDetailShow()
    nextTick(() => { initOpenLayersMap() })
  }
}
// Đóng cửa sổ chi tiết
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
let currentInteraction = null
let drawingMapInstance = null
let drawingVectorSource = null

// Hàm khởi tạo map "Xem chi tiết"
const initOpenLayersMap = () => {
  if (detailMapInstance) {
    detailMapInstance.setTarget(null);
    detailMapInstance = null;
  }
  const mapContainer = document.getElementById('openlayers-map');
  if (!mapContainer) return;

  const vectorSource = new VectorSource();
  const vectorLayer = new VectorLayer({
    source: vectorSource,
    style: new Style({
      stroke: new Stroke({ color: 'red', width: 2 }),
      fill: new Fill({ color: 'rgba(255, 0, 0, 0.1)' })
    })
  });

  detailMapInstance = new Map({
    target: 'openlayers-map',
    layers: [new TileLayer({ source: new OSM() }), vectorLayer],
    view: new View({
      center: fromLonLat([108.4428, 11.9556]),
      zoom: 15
    })
  });

  if (detailFrom.value.geometry) {
    try {
      let geoStr = fixGeoJsonStr(detailFrom.value.geometry);
      let geoObj = JSON.parse(geoStr);
      if (geoObj.type === 'Polygon' && Array.isArray(geoObj.coordinates) && geoObj.coordinates.length > 0) {
        let ring = geoObj.coordinates[0];
        if (ring.length > 0) {
          const firstPoint = ring[0];
          const lastPoint = ring[ring.length - 1];
          if (firstPoint[0] !== lastPoint[0] || firstPoint[1] !== lastPoint[1]) {
            ring.push(firstPoint);
          }
          geoObj.coordinates[0] = ring;
        }
      }
      const validGeoStr = JSON.stringify(geoObj);
      const geoJsonFormat = new GeoJSON();
      // Chỉ định cả dataProjection và featureProjection vì dữ liệu của bạn đã ở EPSG:3857
      const feature = geoJsonFormat.readFeature(validGeoStr, { 
        dataProjection: 'EPSG:3857',
        featureProjection: 'EPSG:3857'
      });
      if (!feature.getGeometry() || !feature.getGeometry().getCoordinates()) {
        throw new Error("Feature không có geometry hợp lệ");
      }
      vectorSource.addFeature(feature);
      const extent = vectorSource.getExtent();
      if (extent && !extent.every(val => val === Infinity)) {
        detailMapInstance.getView().fit(extent, { padding: [20,20,20,20], duration: 200 });
      } else {
        console.error("Extent rỗng, không thể fit view");
      }
    } catch (error) {
      console.error("Lỗi khi đọc GeoJSON ở phần xem chi tiết:", error);
    }
  } else {
    const lat = Number(detailFrom.value.latitude);
    const lng = Number(detailFrom.value.longitude);
    const radius = Number(detailFrom.value.radius);
    if (!isNaN(lat) && !isNaN(lng) && radius > 0) {
      const center = fromLonLat([lng, lat]);
      const circle = new CircleGeom(center, radius);
      const circleFeature = new Feature(circle);
      vectorSource.addFeature(circleFeature);
      detailMapInstance.getView().fit(circle.getExtent(), { padding: [20,20,20,20], duration: 200 });
    }
  }
};
// Hàm xóa toàn bộ vùng vẽ, cho phép người dùng vẽ lại từ đầu
const clearGeometry = () => {
  // Xóa dữ liệu geometry cũ trong form
  formData.value.geometry = null;
  // Xóa hết các feature trong vector source
  if (drawingVectorSource) {
    drawingVectorSource.clear();
  }
  // Nếu có interaction cũ, loại bỏ nó
  if (currentInteraction) {
    drawingMapInstance.removeInteraction(currentInteraction);
    currentInteraction = null;
  }
  // Tạo lại Draw interaction trên vector source đã có
  currentInteraction = new Draw({
    source: drawingVectorSource,
    type: 'Polygon',
    freehand: false  // sử dụng click-to-draw cho tính chính xác
  });
  drawingMapInstance.addInteraction(currentInteraction);
  
  // Gắn sự kiện drawend để lưu dữ liệu mới khi vẽ xong
  currentInteraction.on('drawend', function(event) {
    const geometry = event.feature.getGeometry();
    let coords = geometry.getCoordinates();
    if (coords.length > 0 && coords[0].length > 0) {
      const firstPoint = coords[0][0];
      const lastPoint = coords[0][coords[0].length - 1];
      // Đảm bảo điểm đầu trùng với điểm cuối
      if (firstPoint[0] !== lastPoint[0] || firstPoint[1] !== lastPoint[1]) {
        coords[0].push(firstPoint);
        geometry.setCoordinates(coords);
      }
    }
    const geoJson = new GeoJSON().writeGeometry(geometry);
    formData.value.geometry = geoJson;
    if (geometry.getType() === 'Polygon') {
      const firstPoint = toLonLat(coords[0][0]);
      formData.value.latitude = firstPoint[1];
      formData.value.longitude = firstPoint[0];
      formData.value.radius = 0;
    }
    console.log("Saved GeoJSON:", geoJson);
  });
};

// Hàm chuyển sang chế độ chỉnh sửa (modify)
// const enableModifyMode = () => {
//   if (!formData.value.geometry) {
//     ElMessage({ type: 'warning', message: 'Chưa có vùng vẽ để chỉnh sửa, vui lòng vẽ trước.' });
//     return;
//   }
//   // Khởi tạo lại bản đồ vẽ để load geometry hiện có
//   initDrawingMap(); // initDrawingMap sẽ load geometry từ formData nếu có
//   if (currentInteraction) {
//     drawingMapInstance.removeInteraction(currentInteraction);
//     currentInteraction = null;
//   }
//   // Sử dụng drawingVectorSource đã được tạo trong initDrawingMap
//   currentInteraction = new Modify({
//     source: drawingVectorSource,
//     style: new Style({
//       image: new CircleStyle({
//         radius: 5,
//         fill: new Fill({ color: 'rgba(255, 255, 255, 0.8)' }),
//         stroke: new Stroke({ color: '#ffcc33', width: 2 })
//       })
//     })
//   });
//   drawingMapInstance.addInteraction(currentInteraction);
//   currentInteraction.on('modifyend', (event) => {
//     const modifiedFeatures = event.features.getArray();
//     if (modifiedFeatures.length > 0) {
//       const modifiedFeature = modifiedFeatures[0];
//       let updatedGeom = closeRing(modifiedFeature.getGeometry());
//       modifiedFeature.setGeometry(updatedGeom);
//       const newGeoJson = new GeoJSON().writeGeometry(updatedGeom);
//       formData.value.geometry = newGeoJson;
//       if (updatedGeom.getType() === 'Polygon') {
//         const firstPoint = toLonLat(updatedGeom.getCoordinates()[0][0]);
//         formData.value.latitude = firstPoint[1];
//         formData.value.longitude = firstPoint[0];
//         formData.value.radius = 0;
//       }
//       console.log("Modified GeoJSON:", newGeoJson);
//     }
//   });
// };
// Đặt hàm closeRing ở đây để đảm bảo có sẵn cho toàn bộ file
const closeRing = (geometry) => {
  if (geometry.getType() === 'Polygon') {
    let rings = geometry.getCoordinates();
    if (rings && rings.length > 0) {
      let ring = rings[0];
      if (ring.length > 0) {
        const first = ring[0];
        const last = ring[ring.length - 1];
        if (first[0] !== last[0] || first[1] !== last[1]) {
          ring.push(first);
        }
      }
      geometry.setCoordinates([ring]);
    }
  }
  return geometry;
};

const initDrawingMap = () => {
  // Nếu đã có instance, hủy bỏ nó trước khi khởi tạo mới
  if (drawingMapInstance) {
    drawingMapInstance.setTarget(null);
    drawingMapInstance = null;
  }
  const drawMapContainer = document.getElementById('draw-map');
  if (!drawMapContainer) return;
  
  // Khởi tạo vectorSource và lưu vào biến toàn cục
  drawingVectorSource = new VectorSource();
  const vectorLayer = new VectorLayer({
    source: drawingVectorSource,
    style: new Style({
      stroke: new Stroke({ color: 'blue', width: 2 }),
      fill: new Fill({ color: 'rgba(0, 0, 255, 0.1)' })
    })
  });
  
  // Khởi tạo bản đồ
  drawingMapInstance = new Map({
    target: 'draw-map',
    layers: [new TileLayer({ source: new OSM() }), vectorLayer],
    view: new View({
      center: (formData.value.longitude && formData.value.latitude)
        ? fromLonLat([formData.value.longitude, formData.value.latitude])
        : fromLonLat([108.4428, 11.9556]),
      zoom: 16
    })
  });
  
  // Nếu có interaction cũ, loại bỏ nó
  if (currentInteraction) {
    drawingMapInstance.removeInteraction(currentInteraction);
    currentInteraction = null;
  }
  
  if (formData.value.geometry) {
    try {
      let validGeoJsonStr = fixGeoJsonStr(formData.value.geometry);
      const geoJsonFormat = new GeoJSON();
      let feature = geoJsonFormat.readFeature(validGeoJsonStr, {
        dataProjection: 'EPSG:3857',
        featureProjection: 'EPSG:3857'
      });
      let geom = feature.getGeometry();
      geom = closeRing(geom);
      feature.setGeometry(geom);
      
      // Thêm feature vào vector source
      drawingVectorSource.addFeature(feature);
      
      const extent = drawingVectorSource.getExtent();
      if (extent && !extent.every(val => val === Infinity)) {
        drawingMapInstance.getView().fit(extent, { padding: [20,20,20,20], duration: 200 });
      }
      
      // Thêm Modify interaction để cho phép chỉnh sửa feature
      currentInteraction = new Modify({
        source: drawingVectorSource,
        style: new Style({
          image: new CircleStyle({
            radius: 5,
            fill: new Fill({ color: 'rgba(255, 255, 255, 0.8)' }),
            stroke: new Stroke({ color: '#ffcc33', width: 2 })
          })
        })
      });
      drawingMapInstance.addInteraction(currentInteraction);
      
      currentInteraction.on('modifyend', (event) => {
        const modifiedFeatures = event.features.getArray();
        if (modifiedFeatures.length > 0) {
          const modifiedFeature = modifiedFeatures[0];
          let updatedGeom = closeRing(modifiedFeature.getGeometry());
          modifiedFeature.setGeometry(updatedGeom);
          const newGeoJson = new GeoJSON().writeGeometry(updatedGeom);
          formData.value.geometry = newGeoJson;
          if (updatedGeom.getType() === 'Polygon') {
            const firstPoint = toLonLat(updatedGeom.getCoordinates()[0][0]);
            formData.value.latitude = firstPoint[1];
            formData.value.longitude = firstPoint[0];
            formData.value.radius = 0;
          }
        }
      });
    } catch (error) {
      console.error("Lỗi khi đọc GeoJSON:", error);
    }
  } else {
    // Nếu chưa có geometry, sử dụng Draw interaction kiểu "Polygon" (click-to-draw)
    currentInteraction = new Draw({
      source: drawingVectorSource,
      type: 'Polygon',
      freehand: false  // Không dùng freehand để đảm bảo tính chính xác
    });
    drawingMapInstance.addInteraction(currentInteraction);
    
    currentInteraction.on('drawend', function(event) {
      const geometry = event.feature.getGeometry();
      let coords = geometry.getCoordinates();
      if (coords.length > 0 && coords[0].length > 0) {
        const firstPoint = coords[0][0];
        const lastPoint = coords[0][coords[0].length - 1];
        // Nếu điểm đầu không trùng với điểm cuối, thêm điểm đầu vào cuối mảng
        if (firstPoint[0] !== lastPoint[0] || firstPoint[1] !== lastPoint[1]) {
          coords[0].push(firstPoint);
          geometry.setCoordinates(coords);
        }
      }
      const geoJson = new GeoJSON().writeGeometry(geometry);
      formData.value.geometry = geoJson;
      if (geometry.getType() === 'Polygon') {
        const firstPoint = toLonLat(coords[0][0]);
        formData.value.latitude = firstPoint[1];
        formData.value.longitude = firstPoint[0];
        formData.value.radius = 0;
      }
      console.log("Saved GeoJSON:", geoJson);
    });
  }
};

// ==================== Watcher Debounce ====================
let detailMapTimeout = null
watch(detailShow, (newVal) => {
  if (newVal) {
    if (detailMapTimeout) clearTimeout(detailMapTimeout)
    detailMapTimeout = setTimeout(() => { nextTick(() => { initOpenLayersMap() }) }, 300)
  }
})

let drawingMapTimeout = null
watch(dialogFormVisible, (newVal) => {
  if (newVal) {
    if (drawingMapTimeout) clearTimeout(drawingMapTimeout)
    drawingMapTimeout = setTimeout(() => { nextTick(() => { initDrawingMap() }) }, 300)
  }
})
</script>

<style></style>
