<template>
  <UploadModal
    :title="'上传文件到 ' + toolBarState.path"
    v-model:visible="toolBarState.uploadVisible"
    :uploading="uploading"
    @doUpload="handleUpload"
  />
  <InputModal
    :title="'在' + toolBarState.path + ' 创建目录'"
    description="输入目录名"
    v-model:visible="toolBarState.mkDirVisble"
    :confirmLoading="mkdiring"
    @ok="handleMkDir"
  />
  <InputModal
    title="获取分享"
    description="输入分享码并提交，文件的文件将会放入到根目录"
    v-model:visible="toolBarState.getShareVisble"
    :confirmLoading="getshareing"
    @ok="(input) => mutateGetShare({ shareID: input })"
  />
  <a-input-group compact>
    <a-tooltip title="返回根目录">
      <a-button
        @click="
          emitChDir('/');
          toolBarState.inputpath = toolBarState.path;
        "
      >
        <template #icon><HomeOutlined /></template>
      </a-button>
    </a-tooltip>
    <a-input
      v-model:value="toolBarState.inputpath"
      @pressEnter="emitChDir(toolBarState.inputpath)"
      :loading="loading"
      style="width: calc(100% - 242px)"
    />
    <a-tooltip title="刷新">
      <a-button
        class="address-toolbar-btn"
        @click="
          emitRefresh();
          toolBarState.inputpath = toolBarState.path;
        "
      >
        <template #icon><SyncOutlined /></template>
      </a-button>
    </a-tooltip>
    <a-tooltip title="回退">
      <a-button class="address-toolbar-btn" @click="emitChDir('../')">
        <template #icon><ArrowLeftOutlined /></template>
      </a-button>
    </a-tooltip>
    <a-tooltip title="创建目录">
      <a-button
        class="address-toolbar-btn"
        @click="toolBarState.mkDirVisble = true"
      >
        <template #icon><FolderAddOutlined /></template>
      </a-button>
    </a-tooltip>
    <a-tooltip title="上传">
      <a-button
        class="address-toolbar-btn"
        @click="toolBarState.uploadVisible = true"
      >
        <template #icon><UploadOutlined /></template>
      </a-button>
    </a-tooltip>
    <a-tooltip title="获取分享">
      <a-button
        class="address-toolbar-btn"
        @click="toolBarState.getShareVisble = true"
      >
        <template #icon><PaperClipOutlined /></template>
      </a-button>
    </a-tooltip>
    <a-tooltip title="粘贴">
      <a-button
        type="primary"
        :disabled="true"
        class="address-toolbar-btn"
        @click="toolBarState.uploadVisible = true"
      >
        <template #icon><CopyFilled /></template>
      </a-button>
    </a-tooltip>
  </a-input-group>
</template>
<style scoped>
.address-toolbar-btn {
  margin-left: 4px;
}
</style>
<script>
"use strict";
import {
  ArrowLeftOutlined,
  HomeOutlined,
  UploadOutlined,
  FolderAddOutlined,
  CopyFilled,
  SyncOutlined,
  PaperClipOutlined,
} from "@ant-design/icons-vue";
import gql from "graphql-tag";
const GQL_UPLOAD_FILE = gql`
  mutation singleUpload($path: String!, $file: Upload!) {
    singleUpload(input: { path: $path, file: $file }) {
      name
      contentType
    }
  }
`;
const GQL_MKDIR = gql`
  mutation createDir($path: String!) {
    createDir(path: $path)
  }
`;
const GQL_GET_SHARE = gql`
  mutation getSharedFile($shareID: String!) {
    getSharedFile(shareID: $shareID)
  }
`;
import UploadModal from "./UploadModal.vue";
import InputModal from "./InputModal.vue";
import { defineComponent, reactive, ref } from "vue";
import { useQuery, useMutation, useResult } from "@vue/apollo-composable";
import { message } from "ant-design-vue";
const icons = {
  ArrowLeftOutlined,
  HomeOutlined,
  UploadOutlined,
  CopyFilled,
  FolderAddOutlined,
  SyncOutlined,
  PaperClipOutlined,
};

export default defineComponent({
  components: {
    ...icons,
    UploadModal,
    InputModal,
  },
  props: {
    path: {
      type: String,
      default: "/",
    },
    inputpath: {
      type: String,
      default: "/",
    },
    loading: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    "toolBarState.path": {
      get() {
        return this.path;
      },
      set(value) {
        this.$emit("update:path", value);
      },
    },
    "toolBarState.mkdiring": {
      get() {
        return this.mkdiring;
      },
      set(value) {
        this.$emit("update:mkdiring", value);
      },
    },
    "toolBarState.uploading": {
      get() {
        return this.uploading;
      },
      set(value) {
        this.$emit("update:uploading", value);
      },
    },
  },
  emits: {
    chDir: null,
    "update:path": null,
    "upadte:uploading": null,
    "update:mkdiring": null,
  },
  setup(props, { emit }) {
    const toolBarState = reactive({
      path: "/",
      inputpath: "/",
      mkDirVisble: false,
      uploadVisible: false,
      getShareVisble: false,
    });

    const {
      mutate: doMkDir,
      onDone: onMkDirDone,
      loading: mkdiring,
    } = useMutation(GQL_MKDIR, () => ({
      fetchPolicy: "no-cache",
    }));

    const {
      mutate: doFileUpload,
      onDone: onFileUploadDone,
      onError: onFileUploadError,
      loading: uploading,
    } = useMutation(GQL_UPLOAD_FILE, {
      fetchPolicy: "no-cache",
    });

    const {
      mutate: mutateGetShare,
      onDone: onGetShareDone,
      onError: onGetShareError,
      loading: getshareing,
    } = useMutation(GQL_GET_SHARE, {
      fetchPolicy: "no-cache",
    });

    const handleUpload = ({ uploadFileList }) => {
      if (uploadFileList.length > 1) {
        message.error("一次只能传一个文件");
        return;
      }
      if (uploadFileList.length == 0) {
        message.error("你没有选择任何文件");
        return;
      }

      let fileToUpload = uploadFileList[0].originFileObj;
      doFileUpload({
        path: toolBarState.path + fileToUpload.name,
        file: fileToUpload,
      });
    };

    const handleMkDir = (name) => {
      console.log("mkdir", toolBarState.path + name);
      doMkDir({
        path: toolBarState.path + name,
      });
    };

    const emitChDir = (path) => {
      emit("chDir", path);
    };
    const emitRefresh = () => {
      emit("chDir", toolBarState.path);
    };

    onMkDirDone(({ data }) => {
      if (data.createDir) {
        message.success("目录已创建");
        emitRefresh();
        toolBarState.mkDirVisble = false;
      }
    });

    onFileUploadDone(({ data }) => {
      message.success("上传成功");
      emitRefresh();
      toolBarState.uploadVisible = false;
    });

    onGetShareDone(({ data }) => {
      message.success("分享转存成功");
      emitRefresh();
      toolBarState.getShareVisble = false;
    });

    return {
      toolBarState,
      handleUpload,
      handleMkDir,
      mutateGetShare,

      emitChDir,
      emitRefresh,
      uploading,
      mkdiring,
      getshareing,
    };
  },
  watch: {
    path(newVal, oldVal) {
      this.toolBarState.inputpath = newVal;
      this.toolBarState.path = newVal;
    },
  },
});
</script>
