<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="goAutoCode(null)">
          Thêm mới
        </el-button>
      </div>
      <el-table :data="tableData">
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="id" width="60" prop="ID" />
        <el-table-column align="left" label="Ngày" width="180">
          <template #default="scope">
            {{
              formatDate(scope.row.CreatedAt)
            }}
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="Tên cấu trúc"
          min-width="150"
          prop="structName"
        />
        <el-table-column
          align="left"
          label="Mô tả cấu trúc"
          min-width="150"
          prop="structCNName"
        />
        <el-table-column
          align="left"
          label="Tên bảng"
          min-width="150"
          prop="tableName"
        />
        <el-table-column
          align="left"
          label="Cờ rollback"
          min-width="150"
          prop="flag"
        >
          <template #default="scope">
            <el-tag v-if="scope.row.flag" type="danger" effect="dark">
              Đã rollback
            </el-tag>
            <el-tag v-else type="success" effect="dark">
              Chưa rollback
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="Hành động" min-width="240">
          <template #default="scope">
            <div>
              <el-button
                type="primary"
                link
                :disabled="scope.row.flag === 1"
                @click="addFuncBtn(scope.row)"
              >
                Thêm phương thức
              </el-button>
              <el-button
                type="primary"
                link
                :disabled="scope.row.flag === 1"
                @click="openDialog(scope.row)"
              >
                Rollback
              </el-button>
              <el-button type="primary" link @click="goAutoCode(scope.row)">
                Sử dụng lại
              </el-button>
              <el-button type="primary" link @click="deleteRow(scope.row)">
                Xóa
              </el-button>
            </div>
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
    <el-dialog
      v-model="dialogFormVisible"
      :title="dialogFormTitle"
      :before-close="closeDialog"
      width="600px"
    >
      <el-form :inline="true" :model="formData" label-width="80px">
        <el-form-item label="Tùy chọn:">
          <el-checkbox
            v-model="formData.deleteApi"
            label="Xóa API"
          />
          <el-checkbox
            v-model="formData.deleteMenu"
            label="Xóa menu"
          />
          <el-checkbox
            v-model="formData.deleteTable"
            label="Xóa bảng"
            @change="deleteTableCheck"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">
            Hủy
          </el-button>
          <el-popconfirm
            title="Thao tác này sẽ rollback các file đã tạo và các mục đã chọn, bạn có muốn tiếp tục?"
            @confirm="enterDialog"
          >
            <template #reference>
              <el-button type="primary">
                Xác nhận
              </el-button>
            </template>
          </el-popconfirm>
        </div>
      </template>
    </el-dialog>


    <el-drawer
      v-model="funcFlag"
      size="60%"
      :show-close="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">Thanh công cụ</span>
          <div>
            <el-button
              type="primary"
              @click="runFunc"
            >
              Tạo
            </el-button>
            <el-button
              type="primary"
              @click="closeFunc"
            >
              Hủy
            </el-button>
          </div>
        </div>
      </template>
      <div class="">
        <el-form label-position="top" :model="autoFunc" label-width="80px">
          <el-form-item label="Tên gói:">
            <el-input v-model="autoFunc.package" placeholder="Nhập tên gói" disabled />
          </el-form-item>
          <el-form-item label="Tên cấu trúc:">
            <el-input v-model="autoFunc.structName" placeholder="Nhập tên cấu trúc" disabled />
          </el-form-item>
          <el-form-item label="Tên file frontend:">
            <el-input v-model="autoFunc.packageName" placeholder="Nhập tên file" disabled />
          </el-form-item>
          <el-form-item label="Tên file backend:">
            <el-input v-model="autoFunc.humpPackageName" placeholder="Nhập tên file" disabled />
          </el-form-item>
          <el-form-item label="Mô tả:">
            <el-input v-model="autoFunc.description" placeholder="Nhập mô tả" disabled />
          </el-form-item>
          <el-form-item label="Viết tắt:">
            <el-input v-model="autoFunc.abbreviation" placeholder="Nhập viết tắt" disabled />
          </el-form-item>
          <el-form-item label="Tên phương thức:">
            <el-input v-model="autoFunc.funcName" placeholder="Nhập tên phương thức" />
          </el-form-item>
          <el-form-item label="Phương thức:">
            <el-select v-model="autoFunc.method" placeholder="Chọn phương thức">
              <el-option
                v-for="item in ['GET', 'POST', 'PUT', 'DELETE']"
                :key="item"
                :label="item"
                :value="item"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="Đường dẫn router:">
            <el-input v-model="autoFunc.router" placeholder="Đường dẫn router" />
            <div>Đường dẫn API: [{{ autoFunc.method }}]  /{{ autoFunc.abbreviation }}/{{ autoFunc.router }}</div>
          </el-form-item>
        </el-form>
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
import { getSysHistory, rollback, delSysHistory,addFunc } from "@/api/autoCode.js";
import { useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import { ref } from "vue";
import { formatDate } from "@/utils/format";

defineOptions({
  name: "AutoCodeAdmin",
});

const formData = ref({
  id: undefined,
  deleteApi: true,
  deleteMenu: true,
  deleteTable: false,
});

const router = useRouter();
const dialogFormVisible = ref(false);
const dialogFormTitle = ref("");

const page = ref(1);
const total = ref(0);
const pageSize = ref(10);
const tableData = ref([]);

const autoFunc = ref({
  package:"",
  funcName:"",
  structName:"",
  packageName:"",
  description:"",
  abbreviation:"",
  humpPackageName:"",
  businessDB:"",
  method:"",
})

const addFuncBtn =  (row) => {
  const req = JSON.parse(row.request)
  autoFunc.value.package = req.package
  autoFunc.value.structName = req.structName
  autoFunc.value.packageName = req.packageName
  autoFunc.value.description = req.description
  autoFunc.value.abbreviation = req.abbreviation
  autoFunc.value.humpPackageName = req.humpPackageName
  autoFunc.value.businessDB = req.businessDB
  autoFunc.value.method = ""
  autoFunc.value.funcName = ""
  autoFunc.value.router = ""
  funcFlag.value = true;
};

const funcFlag = ref(false);

const closeFunc = () => {
  funcFlag.value = false;
};

const runFunc = async () =>{
  const res = await addFunc(autoFunc.value)
  if (res.code === 0) {
    ElMessage.success("增加方法成功");
    closeFunc()
  }
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val;
  getTableData();
};

const handleCurrentChange = (val) => {
  page.value = val;
  getTableData();
};

// 查询
const getTableData = async () => {
  const table = await getSysHistory({
    page: page.value,
    pageSize: pageSize.value,
  });
  if (table.code === 0) {
    tableData.value = table.data.list;
    total.value = table.data.total;
    page.value = table.data.page;
    pageSize.value = table.data.pageSize;
  }
};

getTableData();

const deleteRow = async (row) => {
  ElMessageBox.confirm("此操作将删除本历史, 是否继续?", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(async () => {
    const res = await delSysHistory({ id: Number(row.ID) });
    if (res.code === 0) {
      ElMessage.success("删除成功");
      getTableData();
    }
  });
};

// 打开弹窗
const openDialog = (row) => {
  dialogFormTitle.value = "回滚：" + row.structName;
  formData.value.id = row.ID;
  dialogFormVisible.value = true;
};

// 关闭弹窗
const closeDialog = () => {
  dialogFormVisible.value = false;
  formData.value = {
    id: undefined,
    deleteApi: true,
    deleteMenu: true,
    deleteTable: false,
  };
};

// 确认删除表
const deleteTableCheck = (flag) => {
  if (flag) {
    ElMessageBox.confirm(
      `此操作将删除自动创建的文件和api（会删除表！！！）, 是否继续?`,
      "提示",
      {
        closeOnClickModal: false,
        distinguishCancelAndClose: true,
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      }
    )
      .then(() => {
        ElMessageBox.confirm(
          `此操作将删除自动创建的文件和api（会删除表！！！）, 请继续确认！！！`,
          "会删除表",
          {
            closeOnClickModal: false,
            distinguishCancelAndClose: true,
            confirmButtonText: "确定",
            cancelButtonText: "取消",
            type: "warning",
          }
        ).catch(() => {
          formData.value.deleteTable = false;
        });
      })
      .catch(() => {
        formData.value.deleteTable = false;
      });
  }
};

const enterDialog = async () => {
  const res = await rollback(formData.value);
  if (res.code === 0) {
    ElMessage.success("回滚成功");
    getTableData();
  }
};

const goAutoCode = (row) => {
  if (row) {
    router.push({
      name: "autoCodeEdit",
      params: {
        id: row.ID,
      },
    });
  } else {
    router.push({ name: "autoCode" });
  }
};
</script>
