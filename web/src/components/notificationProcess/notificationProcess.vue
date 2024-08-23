<template>
  <teleport to="body">
    <div v-if="displayFileProcess.length > 0" class="noti-container">
      <div class="noti-top">
        <div class="noti-title">
          <span>Quá trình tải dữ liệu ({{ displayFileProcess.length }})</span>
        </div>
        <div class="noti-action">
          <el-button style="outline: none; border: none; background-color: transparent;" @click="handleExpand">
            <el-icon v-if="isExpand" size="18">
              <ArrowDownBold />
            </el-icon>
            <el-icon v-else size="18">
              <ArrowUpBold />
            </el-icon>
          </el-button>
        </div>
      </div>
      <div v-if="isExpand" class="noti-bottom">
        <div v-if="displayFileProcess.length > 0" class="noti-list-file">
          <div v-for="item in displayFileProcess" :key="item.ID" class="noti-file">
            <div class="noti-file-left">
              <i class="el-icon-document" />
              <p class="noti-file-name truncate cursor-pointer" @click="handleFileProcessClick(item)">
                {{ item?.fileName || 'Đang cập nhật...' }}
              </p>
            </div>
            <div class="noti-file-right">
              <el-progress
                v-if="+item.percent < 100.00"
                class="mx-2"
                type="circle"
                :percentage="+item.percent"
                :width="25"
                :stroke-width="4"
                :show-text="false"
                :color="customColors"
              />

              <el-tag v-else-if="parseInt(item.percent) >= 100 && item.status === 'finish'" type="success">
                Hoàn thành
              </el-tag>

              <el-tag v-else-if="parseInt(item.percent) >= 100 && item.status === 'error'" type="danger">
                Lỗi
              </el-tag>

              <el-icon v-if="parseInt(item.percent) >= 100" class="ms-3 cursor-pointer" :size="18" @click="handleCloseFileProcess(item)">
                <Close />
              </el-icon>
            </div>
          </div>
        </div>
      </div>
    </div>
  </teleport>
</template>

<script>
export default {
  name: 'NotificationProcess',
  inheritAttrs: false
}
</script>

<script setup>
import { ref } from 'vue'

// Import strore
import { storeToRefs } from 'pinia'
import { useFileProcessStore } from '@/pinia/bl'
import { useRouter } from 'vue-router'

const router = useRouter()

const fileProcessStore = useFileProcessStore()
const { displayFileProcess } = storeToRefs(fileProcessStore)
const isExpand = ref(true)

const customColors = [
  { color: '#f56c6c', percentage: 20 },
  { color: '#e6a23c', percentage: 40 },
  { color: '#5cb87a', percentage: 60 },
  { color: '#1989fa', percentage: 80 },
  { color: '#67c23a', percentage: 100 },
]

const handleExpand = () => {
  isExpand.value = !isExpand.value
}

const handleCloseFileProcess = (fp) => {
  fileProcessStore.removeFileProcess(fp)
}

const handleFileProcessClick = (fp) => {
  console.info('Tính năng đang phát triển')
  // router.push({
  //   name: 'fileProcessManagement'
  // })
}
</script>

<style lang="scss" scoped>
.noti-container {
    width: 360px;
    overflow: hidden;

    position: fixed;
    bottom: 0;
    right: 24px;
    z-index: 9999;

    background-color: #fff;
    box-shadow: rgba(100, 100, 111, 0.2) 0px 7px 29px 0px;

    border-top-left-radius: 25px;
    border-top-right-radius: 25px;

    .noti-top {
        background-color: #f7fbff;
        height: 52px;
        padding: 0 20px;

        display: flex;
        align-items: center;
        justify-content: space-between;

        .noti-action {
            cursor: pointer;

            display: flex;
            align-self: center;
            & > button {
                padding: 3px;
            }
        }

        .noti-title {
            font-size: 15px;
            font-weight: bold;
        }
    }

    .noti-file {
        height: 48px;
        padding: 0 20px;
        background-color: #fff;

        display: flex;
        align-items: center;

        .noti-file-left {
            display: flex;
            align-items: center;
            line-height: 48px;

            .noti-file-name {
                font-size: 14px;
                margin-left: 4px;
                max-width: 200px;
            }
        }

        .noti-file-right {
            flex: 1;
            display: flex;
            align-items: center;
            justify-content: flex-end;
        }
    }
}

.truncate {
  text-overflow: ellipsis;
  /* Needed to make it work */
  overflow: hidden;
  white-space: nowrap;
}

</style>
