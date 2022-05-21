<template>
  <a-modal
    v-model:visible="modalVisible"
    :title="title"
    @ok="$emit('doUpload', uploadState)"
    :confirmLoading="uploading"
    ok-text="上传"
  >
    <a-upload-dragger
      v-model:fileList="uploadState.uploadFileList"
      name="file"
      :multiple="true"
      :before-upload="beforeUpload"
      @remove="handleRemoveUploadFile"
    >
      <p class="ant-upload-drag-icon">
        <inbox-outlined></inbox-outlined>
      </p>
      <p class="ant-upload-text">点击或将文件拖拽到此区域以选择文件</p>
      <p class="ant-upload-hint">
        一次支持上传一个文件，点击上传按钮将开始上传
      </p>
    </a-upload-dragger>
  </a-modal>
</template>

<script>
"use strict";
import { defineComponent, reactive } from "vue";
import { InboxOutlined } from "@ant-design/icons-vue";
export default defineComponent({
  name: "Upload",
  props: {
    title: {
      type: String,
      default: "文件上传",
    },
    visible: {
      type: Boolean,
      default: false,
    },
    uploading: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    modalVisible: {
      get() {
        return this.visible;
      },
      set(value) {
        this.$emit("update:visible", value);
      },
    },
  },
  emits: {
    doUpload: null,
    "update:visible": null,
  },
  components: {
    InboxOutlined,
  },
  watch: {
    visible(newVal, oldVal) {
      this.uploadState.uploadFileList = [];
    },
  },
  setup() {
    const uploadState = reactive({
      uploadFileList: [],
      uploading: false,
    });
    const handleRemoveUploadFile = (file) => {
      const index = uploadState.uploadFileList.indexOf(file);
      const newFileList = uploadState.uploadFileList.slice();
      newFileList.splice(index, 1);
      uploadState.uploadFileList = newFileList;
    };

    const beforeUpload = (file) => {
      uploadState.uploadFileList = [...uploadState.uploadFileList, file];
      return false;
    };

    return {
      uploadState,
      handleRemoveUploadFile,
      beforeUpload,
    };
  },
});
</script>
