<template>
  <div v-loading.fullscreen.lock="fullscreenLoading">
    <div class="gva-table-box">
      <warning-bar
        title="Nhấp vào 'Tên tệp/Ghi chú' để chỉnh sửa tên tệp hoặc nội dung ghi chú."
      />
      <div class="gva-btn-list gap-3">
        <upload-common
          :image-common="imageCommon"
          @on-success="getTableData"
        />
        <upload-image
          :image-url="imageUrl"
          :file-size="512"
          :max-w-h="1080"
          @on-success="getTableData"
        />
        <el-input
          v-model="search.keyword"
          class="keyword"
          placeholder="Nhập tên tệp hoặc ghi chú"
        />
        <el-button
          type="primary"
          icon="search"
          @click="getTableData"
        >Tìm kiếm</el-button>
      </div>

      <el-table :data="tableData">
        <el-table-column
          align="left"
          label="Xem trước"
          width="100"
        >
          <template #default="scope">
            <CustomPic
              pic-type="file"
              :pic-src="scope.row.url"
              preview
            />
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="Ngày"
          prop="UpdatedAt"
          width="180"
        >
          <template #default="scope">
            <div>{{ formatDate(scope.row.UpdatedAt) }}</div>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="Tên tệp/Ghi chú"
          prop="name"
          width="180"
        >
          <template #default="scope">
            <div
              class="name"
              @click="editFileNameFunc(scope.row)"
            >{{ scope.row.name }}</div>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="Liên kết"
          prop="url"
          min-width="300"
        />
        <el-table-column
          align="left"
          label="Nhãn"
          prop="tag"
          width="100"
        >
          <template #default="scope">
            <el-tag
              :type="scope.row.tag === 'jpg' ? 'info' : 'success'"
              disable-transitions
            >{{ scope.row.tag }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="Thao tác"
          width="160"
        >
          <template #default="scope">
            <el-button
              icon="download"
              type="primary"
              link
              @click="downloadFile(scope.row)"
            >Tải xuống</el-button>
            <el-button
              icon="delete"
              type="primary"
              link
              @click="deleteFileFunc(scope.row)"
            >Xóa</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :style="{ float: 'right', padding: '20px' }"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { getFileList, deleteFile, editFileName } from '@/api/fileUploadAndDownload'
import { downloadImage } from '@/utils/downloadImg'
import CustomPic from '@/components/customPic/index.vue'
import UploadImage from '@/components/upload/image.vue'
import UploadCommon from '@/components/upload/common.vue'
import { formatDate } from '@/utils/format'
import WarningBar from '@/components/warningBar/warningBar.vue'

import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

defineOptions({
  name: 'Upload',
})

const path = ref(import.meta.env.VITE_BASE_API)

const imageUrl = ref('')
const imageCommon = ref('')

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const search = ref({})
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

// Tìm kiếm
const getTableData = async() => {
  const table = await getFileList({ page: page.value, pageSize: pageSize.value, ...search.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
  // console.log("file path: ", path.value)
}
getTableData()

const deleteFileFunc = async(row) => {
  ElMessageBox.confirm('Hành động này sẽ xóa tệp tin vĩnh viễn, bạn có muốn tiếp tục?', 'Cảnh báo', {
    confirmButtonText: 'Đồng ý',
    cancelButtonText: 'Hủy',
    type: 'warning',
  })
    .then(async() => {
      const res = await deleteFile(row)
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

const downloadFile = (row) => {
  if (row.url.indexOf('http://') > -1 || row.url.indexOf('https://') > -1) {
    downloadImage(row.url, row.name)
  } else {
    debugger
    downloadImage(path.value + '/' + row.url, row.name)
  }
}

/**
 * Chỉnh sửa tên tệp tin hoặc ghi chú
 * @param row
 * @returns {Promise<void>}
 */
const editFileNameFunc = async(row) => {
  ElMessageBox.prompt('Vui lòng nhập tên tệp tin hoặc ghi chú', 'Chỉnh sửa', {
    confirmButtonText: 'Đồng ý',
    cancelButtonText: 'Hủy',
    inputPattern: /\S/,
    inputErrorMessage: 'Không được để trống',
    inputValue: row.name
  }).then(async({ value }) => {
    row.name = value
    // console.log(row)
    const res = await editFileName(row)
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: 'Chỉnh sửa thành công!',
      })
      getTableData()
    }
  }).catch(() => {
    ElMessage({
      type: 'info',
      message: 'Hủy chỉnh sửa'
    })
  })
}
</script>

<style scoped>
.name {
  cursor: pointer;
}

</style>
