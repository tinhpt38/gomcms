<template>
  <div class="p-4 relative">
    <div class="absolute top-4 right-4">
      <el-switch v-model="enableAnswering" active-text="Bật chức năng trả lời" inactive-text="Tắt chức năng trả lời" />
    </div>
    
    <el-button type="primary" @click="showTable = true" class="mb-4" v-if="!showTable">
      Tạo bảng câu hỏi
    </el-button>

    <div v-if="showTable">
      <el-table :data="questions" border stripe class="w-full max-w-2xl mx-auto" :disabled="!enableAnswering">
        <el-table-column prop="text" label="Câu hỏi" width="400">
          <template #default="{ row }">
            <el-input v-model="row.text" placeholder="Nhập câu hỏi..." :disabled="!enableAnswering" />
          </template>
        </el-table-column>  
        <el-table-column label="Hành động" width="100">
          <template #default="{ $index }">
            <el-button type="danger" @click="removeQuestion($index)" size="small" :disabled="!enableAnswering">Xóa</el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="mt-4 w-full max-w-2xl mx-auto flex gap-2">
        <el-input v-model="newQuestion" placeholder="Nhập câu hỏi mới..." class="flex-grow" :disabled="!enableAnswering" />
        <el-button type="success" @click="addQuestion" :disabled="!enableAnswering">Thêm câu hỏi</el-button>
      </div>
    </div>

    <el-button 
      type="primary" 
      class="mt-4" 
      @click="saveQuestions" 
      :disabled="!enableAnswering || questions.length === 0"
    >
      Lưu lên hệ thống
    </el-button>
  </div>
</template>


<script setup>
import { ref, watch } from 'vue';
import { ElMessage } from 'element-plus';
import { createAttendanceQuestion } from '@/api/checkins/attendance';

const props = defineProps({
  acId: {
    type: Number,
    required: true
  }
});

const emit = defineEmits(['onSuccess'])

const showTable = ref(true);
const questions = ref([]);
const newQuestion = ref("");
const enableAnswering = ref(true);

watch(enableAnswering, (newValue) => {
  if (newValue) {
    showTable.value = true;
  }
});

const addQuestion = () => {
  if (newQuestion.value.trim()) {
    questions.value.push({ text: newQuestion.value });
    newQuestion.value = "";
  } else {
    ElMessage.warning('Vui lòng nhập câu hỏi trước khi thêm');
  }
};

const removeQuestion = (index) => {
  questions.value.splice(index, 1);
};

const saveQuestions = async () => {
  try {
    for (const q of questions.value) {
      const payload = {
        attendance_id: props.acId,
        question: q.text
      };

      const res = await createAttendanceQuestion(payload);
      if (res.code !== 0) {
        throw new Error(res.message || 'Gửi thất bại');
      }
    }
    ElMessage.success('Đã lưu tất cả câu hỏi!');
    emit('onSuccess');
  } 
  catch (error) {
    ElMessage.error('Lỗi khi lưu: ' + error.message);
  }
};
</script>
