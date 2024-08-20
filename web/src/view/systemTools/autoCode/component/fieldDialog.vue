<template>
  <div>
    <warning-bar title="id , created_at , updated_at , deleted_at sẽ tự động tạo, không tạo trùng lặp. Khi tìm kiếm, nếu điều kiện là LIKE chỉ hỗ trợ chuỗi" />
    <el-form
      ref="fieldDialogForm"
      :model="middleDate"
      label-width="120px"
      label-position="right"
      :rules="rules"
      class="grid grid-cols-2"
    >
      <el-form-item
        label="Tên trường"
        prop="fieldName"
      >
        <el-input
          v-model="middleDate.fieldName"
          autocomplete="off"
          style="width:80%"
        />
        <el-button
          style="width:18%;margin-left:2%"
          @click="autoFill"
        >
          <span style="font-size: 12px">Tự động điền</span>
        </el-button>
      </el-form-item>
      <el-form-item
        label="Tên tiếng Việt"
        prop="fieldDesc"
      >
        <el-input
          v-model="middleDate.fieldDesc"
          autocomplete="off"
        />
      </el-form-item>
      <el-form-item
        label="Trường JSON"
        prop="fieldJson"
      >
        <el-input
          v-model="middleDate.fieldJson"
          autocomplete="off"
        />
      </el-form-item>
      <el-form-item
        label="Tên trường cơ sở dữ liệu"
        prop="columnName"
      >
        <el-input
          v-model="middleDate.columnName"
          autocomplete="off"
        />
      </el-form-item>
      <el-form-item
        label="Mô tả trường cơ sở dữ liệu"
        prop="comment"
      >
        <el-input
          v-model="middleDate.comment"
          autocomplete="off"
        />
      </el-form-item>
      <el-form-item
        label="Loại trường"
        prop="fieldType"
      >
        <el-select
          v-model="middleDate.fieldType"
          style="width:100%"
          placeholder="Chọn loại trường"
          clearable
          @change="clearOther"
        >
          <el-option
            v-for="item in typeOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
            :disabled="item.disabled"
          />
        </el-select>
      </el-form-item>
      <el-form-item
        :label="middleDate.fieldType === 'enum' ? 'Giá trị enum' : 'Độ dài kiểu'"
        prop="dataTypeLong"
      >
        <el-input
          v-model="middleDate.dataTypeLong"
          :placeholder="middleDate.fieldType === 'enum'?`Ví dụ:'Hà Nội','TP.HCM'`:'Độ dài kiểu cơ sở dữ liệu'"
        />
      </el-form-item>
      <el-form-item
        label="Điều kiện tìm kiếm trường"
        prop="fieldSearchType"
      >
        <el-select
          v-model="middleDate.fieldSearchType"
          :disabled="middleDate.fieldType === 'json'"
          style="width:100%"
          placeholder="Chọn điều kiện tìm kiếm trường"
          clearable
        >
          <el-option
            v-for="item in typeSearchOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
            :disabled="canSelect(item.value)"
          />
        </el-select>
      </el-form-item>
      <el-form-item
        label="Liên kết từ điển"
        prop="dictType"
      >
        <el-select
          v-model="middleDate.dictType"
          style="width:100%"
          :disabled="middleDate.fieldType!=='string'"
          placeholder="Chọn từ điển"
          clearable
        >
          <el-option
            v-for="item in dictOptions"
            :key="item.type"
            :label="`${item.type}(${item.name})`"
            :value="item.type"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="Giá trị mặc định">
        <el-input
          v-model="middleDate.defaultValue"
          placeholder="Nhập giá trị mặc định"
        />
      </el-form-item>
      <el-form-item label="Khóa chính">
        <el-checkbox v-model="middleDate.primaryKey" />
      </el-form-item>
      <el-form-item
        label="Loại chỉ mục"
        prop="fieldIndexType"
      >
        <el-select
          v-model="middleDate.fieldIndexType"
          :disabled="middleDate.fieldType === 'json'"
          style="width:100%"
          placeholder="Chọn loại chỉ mục trường"
          clearable
        >
          <el-option
            v-for="item in typeIndexOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
            :disabled="canSelect(item.value)"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="Tạo/Chỉnh sửa giao diện trước">
        <el-switch v-model="middleDate.form" />
      </el-form-item>
      <el-form-item label="Cột bảng giao diện trước">
        <el-switch v-model="middleDate.table" />
      </el-form-item>
      <el-form-item label="Chi tiết giao diện trước">
        <el-switch v-model="middleDate.desc" />
      </el-form-item>
      <el-form-item label="Có sắp xếp">
        <el-switch v-model="middleDate.sort" />
      </el-form-item>
      <el-form-item label="Bắt buộc">
        <el-switch v-model="middleDate.require" />
      </el-form-item>
      <el-form-item label="Có thể xóa">
        <el-switch v-model="middleDate.clearable" />
      </el-form-item>
      <el-form-item label="Ẩn điều kiện tìm kiếm">
        <el-switch v-model="middleDate.fieldSearchHide" :disabled="!middleDate.fieldSearchType" />
      </el-form-item>
      <el-form-item label="Văn bản thông báo lỗi">
        <el-input v-model="middleDate.errorText" />
      </el-form-item>
    </el-form>
    <el-collapse v-model="activeNames">
      <el-collapse-item
        title="Cấu hình nguồn dữ liệu (Cấu hình nâng cao, nếu không vững chắc về lập trình có thể dẫn đến mã tự động không sử dụng được)"
        name="1"
      >
        <el-row :gutter="8">
          <el-col
            :span="3"
          >
            <el-select
              v-model="middleDate.dataSource.association"
              placeholder="Chế độ liên kết"
              @change="associationChange"
            >
              <el-option
                label="Một một"
                :value="1"
              />
              <el-option
                label="Một nhiều"
                :value="2"
              />
            </el-select>
          </el-col>

          <el-col :span="7">
            <el-select
              v-model="middleDate.dataSource.table" placeholder="Chọn bảng nguồn dữ liệu"
              filterable allow-create @focus="getDBTableList" @change="selectDB"
            >
              <el-option
                v-for="item in dbTableList" :key="item.tableName" :label="item.tableName"
                :value="item.tableName"
              />
            </el-select>
            <!-- <el-input v-model="middleDate.dataSource.table" placeholder="Bảng nguồn dữ liệu" /> -->
          </el-col>
          <el-col :span="7">
            <el-select v-model="middleDate.dataSource.value" placeholder="Hãy chọn dữ liệu cần lưu trữ">
              <template #label="{ value }">
                <span>Lưu trữ: </span>
                <span style="font-weight: bold">{{ value }}</span>
              </template>
              <el-option v-for="item in dbColumnList" :key="item.columnName" :value="item.columnName">
                <span style="float: left"> <el-tag :type="item.isPrimary ? 'primary' : 'info'">
                  {{ item.isPrimary ? "Khóa chính" : "Không phải khóa chính" }}
                </el-tag> {{ item.columnName }}</span>
                <span
                  style="
          float: right;
          margin-left:5px;
          color: var(--el-text-color-secondary);
          font-size: 13px;
        "
                >
                  Kiểu dữ liệu: {{ item.type }} <block v-if="item.comment != ''">, Mô tả trường: {{ item.comment }}</block>
                </span>
              </el-option>
            </el-select>
            <!-- <el-input v-model="middleDate.dataSource.value" placeholder="Trường lưu trữ" /> -->
          </el-col>
          <el-col :span="7">
            <el-select v-model="middleDate.dataSource.label" placeholder="Hãy chọn dữ liệu cần hiển thị">
              <template #label="{ value }">
                <span>Hiển thị: </span>
                <span style="font-weight: bold">{{ value }}</span>
              </template>
              <el-option v-for="item in dbColumnList" :key="item.columnName" :value="item.columnName">
                <span style="float: left"> <el-tag :type="item.isPrimary ? 'primary' : 'info'">
                  {{ item.isPrimary ? "Khóa chính" : "Không phải khóa chính" }}
                </el-tag> {{ item.columnName }}</span>
                <span
                    style="
          float: right;
          margin-left:5px;
          color: var(--el-text-color-secondary);
          font-size: 13px;
        "
                >
                  Kiểu dữ liệu: {{ item.type }} <span v-if="item.comment != ''">, Mô tả trường: {{ item.comment }}</span>
                </span>
              </el-option>
            </el-select>
            <!-- <el-input v-model="middleDate.dataSource.label" placeholder="Trường hiển thị" /> -->
          </el-col>
        </el-row>
      </el-collapse-item>
    </el-collapse>
  </div>
</template>

<script setup>
import { toLowerCase, toSQLLine } from '@/utils/stringFun'
import { getSysDictionaryList } from '@/api/sysDictionary'
import WarningBar from '@/components/warningBar/warningBar.vue'
import { ref } from 'vue'
import { ElMessageBox } from 'element-plus'
import {getColumn, getTable} from "@/api/autoCode";

defineOptions({
  name: 'FieldDialog'
})

const props = defineProps({
  dialogMiddle: {
    type: Object,
    default: function() {
      return {}
    }
  },
  typeOptions: {
    type: Array,
    default: function() {
      return []
    }
  },
  typeSearchOptions: {
    type: Array,
    default: function() {
      return []
    }
  },
  typeIndexOptions: {
    type: Array,
    default: function() {
      return []
    }
  },
})

const activeNames = ref([])

const middleDate = ref({})
const dictOptions = ref([])

const validateDataTypeLong = (rule, value, callback) => {
  const regex = /^('([^']*)'(?:,'([^']+)'*)*)$/;
  if (middleDate.value.fieldType == "enum" && !regex.test(value)) {
    callback(new Error("Lỗi kiểm tra giá trị enum"));
  } else {
    callback();
  }
};

const rules = ref({
  fieldName: [
    { required: true, message: 'Vui lòng nhập tên trường tiếng Anh', trigger: 'blur' }
  ],
  fieldDesc: [
    { required: true, message: 'Vui lòng nhập tên trường tiếng Việt', trigger: 'blur' }
  ],
  fieldJson: [
    { required: true, message: 'Vui lòng nhập định dạng trường json', trigger: 'blur' }
  ],
  columnName: [
    { required: true, message: 'Vui lòng nhập trường cơ sở dữ liệu', trigger: 'blur' }
  ],
  fieldType: [
    { required: true, message: 'Vui lòng chọn loại trường', trigger: 'blur' }
  ],
  dataTypeLong: [
    { validator: validateDataTypeLong, trigger: "blur" }
  ],
})

const init = async() => {
  middleDate.value = props.dialogMiddle
  const dictRes = await getSysDictionaryList({
    page: 1,
    pageSize: 999999
  })

  dictOptions.value = dictRes.data
}
init()

const autoFill = () => {
  middleDate.value.fieldJson = toLowerCase(middleDate.value.fieldName)
  middleDate.value.columnName = toSQLLine(middleDate.value.fieldJson)
}

const canSelect = (item) => {
  const fieldType = middleDate.value.fieldType
  if (fieldType !== 'string' && item === 'LIKE') {
    return true
  }

  if ((fieldType !== 'int' && fieldType !== 'time.Time' && fieldType !== 'float64') && (item === 'BETWEEN' || item === 'NOT BETWEEN')) {
    return true
  }
  return false
}

const clearOther = () => {
  middleDate.value.fieldSearchType = ''
  middleDate.value.dictType = ''
}

const associationChange = (val) => {
  if (val === 2) {
    ElMessageBox.confirm(
      'Trong chế độ liên kết một-nhiều, kiểu dữ liệu sẽ thay đổi thành mảng, được biểu diễn bởi json, cụ thể là chế độ mảng, bạn có muốn tiếp tục?',
      'Cảnh báo',
      {
        confirmButtonText: 'Tiếp tục',
        cancelButtonText: 'Hủy',
        type: 'warning'
      }
    ).then(() => {
      middleDate.value.fieldType = 'array'
    }).catch(() => {
      middleDate.value.dataSource.association = 1
    })
  }
}


const dbTableList = ref([])

const getDBTableList = async () => {
  const res = await getTable()
  console.log(res);
  if (res.code === 0) {
    let list = res.data.tables; // Đảm bảo đây là mảng tables
    dbTableList.value = list.map(item => ({
      tableName: item.tableName,
      value: item.tableName // Giả sử value cũng là tableName, nếu khác vui lòng điều chỉnh
    }));
  }
}

const dbColumnList = ref([])
const selectDB = async (val) => {
  middleDate.value.dataSource.table = val
  const res = await getColumn({
    tableName: val
  })
  console.log(res)
  if (res.code === 0) {
    let list = res.data.columns; // Đảm bảo đây là mảng columns
    dbColumnList.value = list.map(item => ({
      columnName: item.columnName,
      value: item.columnName,
      type: item.dataType,
      isPrimary: item.primaryKey,
      comment: item.columnComment
    }));
    if (dbColumnList.value.length > 0) {
      middleDate.value.dataSource.label = dbColumnList.value[0].columnName
      middleDate.value.dataSource.value = dbColumnList.value[0].columnName
    }
  }
}


const fieldDialogForm = ref(null)
defineExpose({ fieldDialogForm })
</script>
