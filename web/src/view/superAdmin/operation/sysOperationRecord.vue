<template>
  <div>
    <div class="gva-search-box">
      <el-form
        :inline="true"
        :model="searchInfo"
      >
        <el-form-item label="Phương thức yêu cầu">
          <el-input
            v-model="searchInfo.method"
            placeholder="Điều kiện tìm kiếm"
          />
        </el-form-item>
        <el-form-item label="Đường dẫn yêu cầu">
          <el-input
            v-model="searchInfo.path"
            placeholder="Điều kiện tìm kiếm"
          />
        </el-form-item>
        <el-form-item label="Mã trạng thái kết quả">
          <el-input
            v-model="searchInfo.status"
            placeholder="Điều kiện tìm kiếm"
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            icon="search"
            @click="onSubmit"
          >Tìm kiếm</el-button>
          <el-button
            icon="refresh"
            @click="onReset"
          >Đặt lại</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button
          icon="delete"
          :disabled="!multipleSelection.length"
          @click="onDelete"
        >Xóa</el-button>
      </div>
      <el-table
        ref="multipleTable"
        :data="tableData"
        style="width: 100%"
        tooltip-effect="dark"
        row-key="ID"
        @selection-change="handleSelectionChange"
      >
        <el-table-column
          align="left"
          type="selection"
          width="55"
        />
        <el-table-column
          align="left"
          label="Người thực hiện"
          width="140"
        >
          <template #default="scope">
            <div>{{ scope.row.user.userName }}({{ scope.row.user.nickName }})</div>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="Ngày"
          width="180"
        >
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column
          align="left"
          label="Mã trạng thái"
          prop="status"
          width="120"
        >
          <template #default="scope">
            <div>
              <el-tag type="success">{{ scope.row.status }}</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="Địa chỉ IP yêu cầu"
          prop="ip"
          width="120"
        />
        <el-table-column
          align="left"
          label="Phương thức yêu cầu"
          prop="method"
          width="120"
        />
        <el-table-column
          align="left"
          label="Đường dẫn yêu cầu"
          prop="path"
          width="240"
        />
        <el-table-column
          align="left"
          label="Yêu cầu"
          prop="path"
          width="80"
        >
          <template #default="scope">
            <div>
              <el-popover
                v-if="scope.row.body"
                placement="left-start"
                :width="444"
              >
                <div class="popover-box">
                  <pre>{{ fmtBody(scope.row.body) }}</pre>
                </div>
                <template #reference>
                  <el-icon style="cursor: pointer;"><warning /></el-icon>
                </template>
              </el-popover>

              <span v-else>Không có</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="Phản hồi"
          prop="path"
          width="80"
        >
          <template #default="scope">
            <div>
              <el-popover
                v-if="scope.row.resp"
                placement="left-start"
                :width="444"
              >
                <div class="popover-box">
                  <pre>{{ fmtBody(scope.row.resp) }}</pre>
                </div>
                <template #reference>
                  <el-icon style="cursor: pointer;"><warning /></el-icon>
                </template>
              </el-popover>
              <span v-else>Không có</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="Thao tác"
        >
          <template #default="scope">
            <el-button
              icon="delete"
              type="primary"
              link
              @click="deleteSysOperationRecordFunc(scope.row)"
            >Xóa</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
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
import {
  deleteSysOperationRecord,
  getSysOperationRecordList,
  deleteSysOperationRecordByIds
} from '@/api/sysOperationRecord' // Replace this with your own address
import { formatDate } from '@/utils/format'
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

defineOptions({
  name: 'SysOperationRecord'
})

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const onReset = () => {
  searchInfo.value = {}
}
// Method for frontend conditional search
const onSubmit = () => {
  page.value = 1
  pageSize.value = 10
  if (searchInfo.value.status === '') {
    searchInfo.value.status = null
  }
  getTableData()
}

// Pagination
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// Query
const getTableData = async() => {
  const table = await getSysOperationRecordList({
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

const multipleSelection = ref([])
const handleSelectionChange = (val) => {
  multipleSelection.value = val
}
const onDelete = async() => {
  ElMessageBox.confirm('Are you sure you want to delete?', 'Prompt', {
    confirmButtonText: 'Confirm',
    cancelButtonText: 'Cancel',
    type: 'warning'
  }).then(async() => {
    const ids = []
    multipleSelection.value &&
        multipleSelection.value.forEach(item => {
          ids.push(item.ID)
        })
    const res = await deleteSysOperationRecordByIds({ ids })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: 'Deleted successfully'
      })
      if (tableData.value.length === ids.length && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  })
}
const deleteSysOperationRecordFunc = async(row) => {
  ElMessageBox.confirm('Are you sure you want to delete?', 'Prompt', {
    confirmButtonText: 'Confirm',
    cancelButtonText: 'Cancel',
    type: 'warning'
  }).then(async() => {
    const res = await deleteSysOperationRecord({ ID: row.ID })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: 'Deleted successfully'
      })
      if (tableData.value.length === 1 && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  })
}
const fmtBody = (value) => {
  try {
    return JSON.parse(value)
  } catch (err) {
    return value
  }
}

</script>

<style lang="scss">
.table-expand {
  padding-left: 60px;
  font-size: 0;
  label {
    width: 90px;
    color: #99a9bf;
    .el-form-item {
      margin-right: 0;
      margin-bottom: 0;
      width: 50%;
    }
  }
}
.popover-box {
  background: #112435;
  color: #f08047;
  height: 600px;
  width: 420px;
  overflow: auto;
}
.popover-box::-webkit-scrollbar {
  display: none; /* Chrome Safari */
}
</style>
