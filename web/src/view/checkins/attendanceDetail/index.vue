<template>
  <div class="p-2">
    <div class="p-1 my-1">
      <el-button type="primary" icon="upload">
        Nhập Excel
      </el-button>
      <el-button type="success" icon="download">
        Xuất Excel
      </el-button>
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
                  <el-input v-model="formData.clientUrl" type="text" clearable >
                    <template #prepend>https://</template>
                  </el-input>
                </el-form-item>
              </el-col>
              <el-col :span="12" class="grid-cell flex justify-end">
                <el-button type="primary" @click="saveAttendance">
                  Lưu
                </el-button>
              </el-col>
            </el-row>

            <el-row>
              <el-col :span="12" class="grid-cell">
                <el-form-item label="Ngày bắt đầu" label-width="150px" prop="startDate" class="required">
                  <el-date-picker v-model="formData.startDate" type="datetime" class="full-width-input" clearable />
                </el-form-item>
              </el-col>
              <el-col :span="12" class="grid-cell">
                <el-form-item label="Ngày kết thúc" label-width="150px" prop="endDate" class="required">
                  <el-date-picker v-model="formData.endDate" type="datetime" class="full-width-input" clearable />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row>
              <el-col :span="12" class="grid-cell">
                <el-form-item label="Cho thử nghiệm" label-width="150px" prop="isTrailer">
                  <el-switch v-model="formData.isTrailer" />
                </el-form-item>
              </el-col>
              <el-col :span="12" class="grid-cell">
                <el-form-item label="Khoá" label-width="150px" prop="isLocked">
                  <el-switch v-model="formData.isLocked" />
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
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
          <Partticipants :participants="partticipantsData" />
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
import { useRoute } from 'vue-router';

import { ElForm, ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import Partticipants from '@/view/checkins/components/participant/index.vue'
import Group from '@/view/checkins/components/group/index.vue'
import Area from '@/view/checkins/components/area/index.vue'
import Condition from '@/view/checkins/components/condition/index.vue'
import { getParticipantList } from '@/api/checkins/participant'

defineOptions({
  name: 'AttendanceDetail'
})

const $route = useRoute()
const tabsActiveTab = ref('attendanceInfoTab')
const currentId = ref($route.params.id)
const formData = ref({})
const elFormRef = ref();

const getDetailData = async () => {
  var id = $route.params.id
  console.log('id', id)
  const res = await findAttendance({ id: $route.params.id })
  if (res.code == 0) {
    formData.value = res.data
  }
  console.log('res', res)
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


const partticipantsData = ref([])
const groupData = ref([])
const areaData = ref([])
const conditionsData = ref([])

const initExampleData = () => {
  partticipantsData.value = [
    {
      id: 1,
      email: "ex1@example.com",
      group: "Group 1",
      totalCheckin: 10,
      totalPass: 5, // past condition
      fullname: "Example 1"
    },
    {
      id: 2,
      email: "ex2@example.com",
      group: "Group 2",
      totalCheckin: 8,
      totalPass: 3,
      fullname: "Example 2"
    },
    {
      id: 3,
      email: "ex3@example.com",
      group: "Group 3",
      totalCheckin: 12,
      totalPass: 7,
      fullname: "Example 3"
    },
    {
      id: 4,
      email: "ex4@example.com",
      group: "Group 4",
      totalCheckin: 6,
      totalPass: 2,
      fullname: "Example 4"
    },
    {
      id: 5,
      email: "ex5@example.com",
      group: "Group 5",
      totalCheckin: 15,
      totalPass: 10,
      fullname: "Example 5"
    },
    {
      id: 6,
      email: "ex6@example.com",
      group: "Group 6",
      totalCheckin: 9,
      totalPass: 4,
      fullname: "Example 6"
    },
    {
      id: 7,
      email: "ex7@example.com",
      group: "Group 7",
      totalCheckin: 11,
      totalPass: 6,
      fullname: "Example 7"
    },
    {
      id: 8,
      email: "ex8@example.com",
      group: "Group 8",
      totalCheckin: 7,
      totalPass: 3,
      fullname: "Example 8"
    },
    {
      id: 9,
      email: "ex9@example.com",
      group: "Group 9",
      totalCheckin: 13,
      totalPass: 8,
      fullname: "Example 9"
    },
    {
      id: 10,
      email: "ex10@example.com",
      group: "Group 10",
      totalCheckin: 5,
      totalPass: 2,
      fullname: "Example 10"
    },
    {
      id: 11,
      email: "ex11@example.com",
      group: "Group 11",
      totalCheckin: 14,
      totalPass: 9,
      fullname: "Example 11"
    },
    {
      id: 12,
      email: "ex12@example.com",
      group: "Group 12",
      totalCheckin: 7,
      totalPass: 4,
      fullname: "Example 12"
    },
    {
      id: 13,
      email: "ex13@example.com",
      group: "Group 13",
      totalCheckin: 10,
      totalPass: 6,
      fullname: "Example 13"
    },
    {
      id: 14,
      email: "ex14@example.com",
      group: "Group 14",
      totalCheckin: 9,
      totalPass: 5,
      fullname: "Example 14"
    },
    {
      id: 15,
      email: "ex15@example.com",
      group: "Group 15",
      totalCheckin: 12,
      totalPass: 8,
      fullname: "Example 15"
    },
    {
      id: 16,
      email: "ex16@example.com",
      group: "Group 16",
      totalCheckin: 6,
      totalPass: 3,
      fullname: "Example 16"
    },
    {
      id: 17,
      email: "ex17@example.com",
      group: "Group 17",
      totalCheckin: 11,
      totalPass: 7,
      fullname: "Example 17"
    },
    {
      id: 18,
      email: "ex18@example.com",
      group: "Group 18",
      totalCheckin: 8,
      totalPass: 4,
      fullname: "Example 18"
    },
    {
      id: 19,
      email: "ex19@example.com",
      group: "Group 19",
      totalCheckin: 13,
      totalPass: 9,
      fullname: "Example 19"
    },
    {
      id: 20,
      email: "ex20@example.com",
      group: "Group 20",
      totalCheckin: 7,
      totalPass: 5,
      fullname: "Example 20"
    }
  ]


  conditionsData.value = [
    {
      id: 1,
      areaId: 1,
      groupId: 1,
      startAt: "2021-09-01 00:00:00",
      endAt: "2021-09-30 23:59:59",
    },
    {
      id: 2,
      areaId: 2,
      groupId: 2,
      startAt: "2021-10-01 00:00:00",
      endAt: "2021-10-31 23:59:59",
    },
    {
      id: 3,
      areaId: 3,
      groupId: 3,
      startAt: "2021-11-01 00:00:00",
      endAt: "2021-11-30 23:59:59",
    },
    {
      id: 4,
      areaId: 4,
      groupId: 4,
      startAt: "2021-12-01 00:00:00",
      endAt: "2021-12-31 23:59:59",
    },
    {
      id: 5,
      areaId: 5,
      groupId: 5,
      startAt: "2022-01-01 00:00:00",
      endAt: "2022-01-31 23:59:59",
    },
    {
      id: 6,
      areaId: 6,
      groupId: 6,
      startAt: "2022-02-01 00:00:00",
      endAt: "2022-02-28 23:59:59",
    }

  ]

}
initExampleData();

//endregion

//region Group
const groupPage = ref(4)
const groupsSize = ref(20)
const groupSize = ref(20)


const groupHandleCurrentChange = (val) => {
  console.log('groupHandleCurrentChange', val)
}
const groupHandleSizeChange = (val) => {
  console.log('groupHandleSizeChange', val)
}

const getParticipantListData = async () => {
    const res = await getParticipantList({ page: groupPage.value, size: groupSize.value })
    if (res.code == 0) {
        partticipantsData.value = res.data.list
    }
    console.log('res', res)
}

//endregion
getParticipantListData()

//region Area
const areaPage = ref(4)
const areasSize = ref(20)
const areaSize = ref(20)


const areaHandleCurrentChange = (val) => {
  console.log('groupHandleCurrentChange', val)
}
const areaHandleSizeChange = (val) => {
  console.log('groupHandleSizeChange', val)
}
//endregion
</script>

<style lang="scss"></style>