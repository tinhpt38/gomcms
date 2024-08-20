<template>
  <el-drawer v-model="drawer" title="Cấu hình hệ thống" direction="rtl" :size="width">
    <div class="flex flex-col">
      <div class="mb-8">
        <div class="text-gray-800 dark:text-gray-100">Chủ đề mặc định</div>
        <div class="mt-2 text-sm p-2 flex items-center gap-2">
          <el-segmented
            v-model="config.darkMode"
            :options="options"
            size="default"
            @change="appStore.toggleDarkMode"
          />
        </div>
      </div>
      <div class="mb-8">
        <div class="text-gray-800 dark:text-gray-100">Màu chủ đề</div>
        <div class="mt-2 text-sm p-2 flex items-center gap-2">
          <div
            v-for="item in colors"
            :key="item"
            class="w-5 h-5 rounded cursor-pointer flex items-center justify-center"
            :style="`background:${item}`"
            @click="appStore.togglePrimaryColor(item)"
          >
            <el-icon v-if="config.primaryColor === item">
              <Select />
            </el-icon>
          </div>
          <el-color-picker
            v-model="customColor"
            @change="appStore.togglePrimaryColor"
          />
        </div>
      </div>
      <div class="mb-8">
        <div class="text-gray-800 dark:text-gray-100">Giao diện hiển thị</div>
        <div class="mt-2 text-sm p-2">
          <div class="flex items-center justify-between">
            <div>Hiển thị watermark</div>
            <el-switch
              v-model="config.show_watermark"
              @change="appStore.toggleConfigWatermark"
            />
          </div>
          <div class="flex items-center justify-between">
            <div>Chế độ xám</div>
            <el-switch v-model="config.grey" @change="appStore.toggleGrey" />
          </div>
          <div class="flex items-center justify-between">
            <div>Chế độ yếu</div>
            <el-switch
              v-model="config.weakness"
              @change="appStore.toggleWeakness"
            />
          </div>
          <div class="flex items-center justify-between">
            <div>Chế độ menu</div>
            <el-segmented
              v-model="config.side_mode"
              :options="sideModes"
              size="default"
              @change="appStore.toggleSideModel"
            />
            <!-- <el-select
              v-model="config.side_mode"
              @change="handleSideModelChange"
            >
              <el-option value="normal" label="Chế độ tiêu chuẩn" />
              <el-option value="head" label="Thanh điều hướng trên cùng" />
              <el-option value="multilayer" disabled label="Chế độ đa cấp" />
            </el-select> -->
          </div>

          <div class="flex items-center justify-between">
            <div>Hiển thị tab</div>
            <el-switch
              v-model="config.showTabs"
              @change="appStore.toggleTabs"
            />
          </div>
        </div>
      </div>

      <div class="mb-8">
        <div class="text-gray-800 dark:text-gray-100">Cấu hình kích thước layout</div>
        <div class="mt-2 text-sm p-2">
          <div class="flex items-center justify-between mb-2">
            <div>Chiều rộng mở rộng thanh bên</div>
            <el-input-number
              v-model="config.layout_side_width"
              :min="150"
              :max="400"
              :step="10"
            ></el-input-number>
          </div>
          <div class="flex items-center justify-between mb-2">
            <div>Chiều rộng thu gọn thanh bên</div>
            <el-input-number
              v-model="config.layout_side_collapsed_width"
              :min="60"
              :max="100"
            ></el-input-number>
          </div>
          <div class="flex items-center justify-between mb-2">
            <div>Chiều cao mục con thanh bên</div>
            <el-input-number
              v-model="config.layout_side_item_height"
              :min="30"
              :max="50"
            ></el-input-number>
          </div>
        </div>
      </div>

      <el-alert type="warning" :closable="false">
        Vui lòng lưu tất cả cấu hình vào tệp
        <el-tag>config.json</el-tag> trên máy tính của bạn, nếu không cấu hình sẽ bị mất sau khi làm mới trang
      </el-alert>

      <el-button type="primary" class="mt-4" @click="copyConfig"
        >Sao chép cấu hình json</el-button
      >
    </div>
  </el-drawer>
</template>

<script setup>
import { useAppStore } from "@/pinia";
import { storeToRefs } from "pinia";
import { ref, computed } from "vue";
import { ElMessage } from "element-plus";
const appStore = useAppStore();
const { config, device } = storeToRefs(appStore);
defineOptions({
  name: "GvaSetting",
});

const width = computed(() => {
  return device.value === "mobile" ? "100%" : "500px";
});

const colors = [
  "#EB2F96",
  "#3b82f6",
  "#2FEB54",
  "#EBEB2F",
  "#EB2F2F",
  "#2FEBEB",
];

const drawer = defineModel("drawer", {
  default: true,
  type: Boolean,
});

const options = ["tối", "sáng", "tự động"];
const sideModes = [
  {
    label : "Chế độ tiêu chuẩn",
    value : "normal"
  },
  {
    label : "Thanh điều hướng trên cùng",
    value: "head"
  },
  {
    label : "Chế độ kết hợp",
    value: "combination"
  }
];

const copyConfig = () => {
  const input = document.createElement("textarea");
  input.value = JSON.stringify(config.value);
  // Thêm dấu xuống dòng
  input.value = input.value.replace(/,/g, ",\n");
  document.body.appendChild(input);
  input.select();
  document.execCommand("copy");
  document.body.removeChild(input);
  ElMessage.success("Sao chép thành công, vui lòng lưu vào tệp tin cục bộ");
};

const customColor = ref("");

const handleSideModelChange = (e) => {
  console.log(e);
};
</script>

<style lang="scss" scoped>
::v-deep(.el-drawer__header) {
  @apply border-gray-400 dark:border-gray-600;
}
</style>
