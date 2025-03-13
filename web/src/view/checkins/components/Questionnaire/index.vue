<template>
  <div class="p-4 relative">
    <div class="absolute top-4 right-4">
      <el-switch v-model="enableAnswering" active-text="Bật chức năng trả lời" inactive-text="Tắt chức năng trả lời" />
    </div>
    
    <el-button type="primary" @click="showTable = true" class="mb-4" v-if="!showTable">
      Tạo bảng câu hỏi và câu trả lời
    </el-button>

    <div v-if="showTable">
      <el-table :data="questions" border stripe class="w-full max-w-2xl mx-auto" :disabled="!enableAnswering">
        <el-table-column prop="text" label="Câu hỏi" width="300">
          <template #default="{ row }">
            <el-input v-model="row.text" placeholder="Nhập câu hỏi..." :disabled="!enableAnswering" />
          </template>
        </el-table-column>
        <el-table-column label="Câu trả lời">
          <template #default="{ row }">
            <el-input v-model="row.answer" placeholder="Nhập câu trả lời..." :disabled="!enableAnswering" />
          </template>
        </el-table-column>
        <el-table-column label="Hành động" width="100">
          <template #default="{ row, $index }">
            <el-button type="danger" @click="removeQuestion($index)" size="small" :disabled="!enableAnswering">Xóa</el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="mt-4 w-full max-w-2xl mx-auto flex gap-2">
        <el-input v-model="newQuestion" placeholder="Nhập câu hỏi mới..." class="flex-grow" :disabled="!enableAnswering" />
        <el-button type="success" @click="addQuestion" :disabled="!enableAnswering">Thêm câu hỏi</el-button>
      </div>
    </div>
  </div>
</template>

<script setup>

import{
  createQuestion

}from '@/api/checkins/question'

import { ref, watch } from 'vue';
import { ElMessage } from 'element-plus';

const props = defineProps({
  acId: {
    type: String,
    required: true,
  },
});

const showTable = ref(false);
const questions = ref([]);
const newQuestion = ref("");
const enableAnswering = ref(true);

const addQuestion = () => {
  if (newQuestion.value.trim()) {
    questions.value.push({ text: newQuestion.value, answer: '' });
    newQuestion.value = "";
  } else {
    ElMessage.warning('Vui lòng nhập câu hỏi trước khi thêm');
  }
};

const removeQuestion = (index) => {
  questions.value.splice(index, 1);
};
</script>
