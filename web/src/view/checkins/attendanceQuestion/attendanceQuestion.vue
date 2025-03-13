<template>
    <div>
      <div class="gva-search-box">
        <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule"
          @keyup.enter="onSubmit">
          <el-form-item label="Câu hỏi" prop="question">
            <el-input v-model="searchInfo.question" placeholder="Nhập câu hỏi" clearable />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" icon="search" @click="onSubmit">Tìm kiếm</el-button>
            <el-button icon="refresh" @click="onReset">Làm mới</el-button>
          </el-form-item>
        </el-form>
      </div>
      
      <div class="gva-table-box">
        <div class="gva-btn-list">
          <el-button type="primary" icon="plus" @click="openDialog">Thêm mới</el-button>
          <el-button icon="delete" :disabled="!multipleSelection.length" @click="onDelete">Xóa</el-button>
        </div>
        
        <el-table ref="multipleTable" :data="tableData" row-key="ID" @selection-change="handleSelectionChange">
          <el-table-column align="left" label="Câu hỏi" prop="question" width="400" />
          <el-table-column align="left" label="Hành động" fixed="right" min-width="240">
            <template #default="scope">
              <el-button type="primary" link icon="edit" @click="updateQuestion(scope.row)">Chỉnh sửa</el-button>
              <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">Xóa</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>
  </template>
  
  <script setup>
  import { ref, reactive } from 'vue';
  import { ElMessage, ElMessageBox } from 'element-plus';
  import { getAttendanceQuestionList, createAttendanceQuestion, updateAttendanceQuestion, deleteAttendanceQuestion } from '@/api/attendanceQuestion';
  
  const tableData = ref([]);
  const searchInfo = ref({ question: '' });
  const multipleSelection = ref([]);
  
  const getTableData = async () => {
    const response = await getAttendanceQuestionList({ ...searchInfo.value });
    if (response.code === 0) {
      tableData.value = response.data.list;
    }
  };
  
  const onSubmit = () => {
    getTableData();
  };
  
  const onReset = () => {
    searchInfo.value.question = '';
    getTableData();
  };
  
  const handleSelectionChange = (val) => {
    multipleSelection.value = val;
  };
  
  const deleteRow = async (row) => {
    await deleteAttendanceQuestion({ ID: row.ID });
    ElMessage.success('Xóa thành công');
    getTableData();
  };
  
  const onDelete = async () => {
    const IDs = multipleSelection.value.map(item => item.ID);
    await deleteAttendanceQuestion({ IDs });
    ElMessage.success('Xóa thành công');
    getTableData();
  };
  
  const openDialog = () => {
    // Open the form dialog for creating a new question
  };
  </script>
  
  <style></style>