<template>
  <div style="margin: 16px 0">
    <Toolbar
      v-model:path="tableState.path"
      :loading="listdiring"
      @chDir="chDir"
    />
  </div>
  <a-table
    :columns="columns"
    :data-source="tableState.files"
    :loading="listdiring"
  >
    <template #bodyCell="{ column, record }">
      <template v-if="column.key === 'name'">
        <span>
          <template v-if="record.contentType === 'file'">
            <file-outlined style="margin: 0 10px" />
            <span>{{ record.name }}</span>
          </template>
          <template v-else>
            <folder-outlined style="margin: 0 10px" />
            <a @click="chDir(tableState.path + record.name)">
              {{ record.name }}
            </a>
          </template>
        </span>
      </template>
      <template v-if="column.key === 'action'">
        <template v-if="record.contentType === 'file'">
          <a-tooltip title="下载">
            <a-button
              type="primary"
              class="tblop"
              @click="handleDownloadFile(record.name)"
            >
              <template #icon><DownloadOutlined /></template>
            </a-button>
          </a-tooltip>
          <a-tooltip title="重命名">
            <a-button class="tblop">
              <template #icon><EditOutlined /></template>
            </a-button>
          </a-tooltip>

          <a-tooltip title="剪切">
            <a-button class="tblop">
              <template #icon><ScissorOutlined /></template>
            </a-button>
          </a-tooltip>

          <a-tooltip title="复制">
            <a-button class="tblop">
              <template #icon><CopyOutlined /></template>
            </a-button>
          </a-tooltip>

          <a-tooltip title="分享">
            <a-button
              @click="
                mutateCreateShare({
                  path: tableState.path + record.name,
                  expireDays: 7,
                })
              "
              class="tblop"
            >
              <template #icon><SendOutlined /></template>
            </a-button>
          </a-tooltip>

          <a-popconfirm
            title="你确定要删除吗，该操作不可逆！"
            ok-text="确定"
            cancel-text="取消"
            @confirm="mutateRmfile({ path: tableState.path + record.name })"
          >
            <a-tooltip title="Remove File">
              <a-button danger class="tblop" href="#">
                <template #icon><DeleteOutlined /></template>
              </a-button>
            </a-tooltip>
          </a-popconfirm>
        </template>

        <template v-else>
          <a-tooltip title="打开">
            <a-button
              type="primary"
              class="tblop"
              @click="chDir(tableState.path + record.name)"
            >
              <template #icon><FolderOpenFilled /></template>
            </a-button>
          </a-tooltip>
          <a-tooltip title="重命名">
            <a-button class="tblop">
              <template #icon><EditOutlined /></template>
            </a-button>
          </a-tooltip>
          <a-tooltip title="剪切">
            <a-button class="tblop">
              <template #icon><ScissorOutlined /></template>
            </a-button>
          </a-tooltip>
          <a-popconfirm
            title="你确定要删除吗，该操作不可逆！"
            ok-text="确定"
            cancel-text="取消"
            @confirm="mutateRmdir({ path: tableState.path + record.name })"
          >
            <a-tooltip title="删除目录">
              <a-button danger class="tblop" href="#">
                <template #icon><DeleteOutlined /></template>
              </a-button>
            </a-tooltip>
          </a-popconfirm>
        </template>
      </template>
    </template>
  </a-table>
</template>
<style scoped>
.tblop {
  margin-right: 8px;
}
</style>
<script>
"use strict";
import {
  ScissorOutlined,
  FileOutlined,
  DownloadOutlined,
  DeleteOutlined,
  DragOutlined,
  EditOutlined,
  FolderOpenFilled,
  FolderOutlined,
  CopyOutlined,
  SendOutlined,
} from "@ant-design/icons-vue";
import Toolbar from "./Toolbar.vue";
import { defineComponent, reactive, ref } from "vue";
import { useRoute } from "vue-router";
import { useQuery, useMutation, useResult } from "@vue/apollo-composable";
import { message, Modal } from "ant-design-vue";

const columns = [
  {
    name: "Name",
    title: "文件名",
    dataIndex: "name",
    key: "name",
    sorter: (a, b) =>
      a.contentType == b.contentType
        ? -a.name.localeCompare(b)
        : a.contentType == "file"
        ? -1
        : 1,
    sortDirections: ["descend"],
  },
  {
    name: "Type",
    title: "类型",
    dataIndex: "contentType",
    key: "type",
  },
  {
    name: "Action",
    title: "操作",
    key: "action",
  },
];

import gql from "graphql-tag";
import axios from "axios";
const GQL_LIST_DIR = gql`
  mutation listDir($path: String!) {
    listDir(path: $path) {
      name
      contentType
    }
  }
`;
const GQL_DELETE_FILE = gql`
  mutation deleteFile($path: String!) {
    deleteFile(path: $path)
  }
`;

const GQL_RMDIR = gql`
  mutation deleteDir($path: String!) {
    deleteDir(path: $path)
  }
`;

const GQL_GET_DOWNLOAD_URL = gql`
  mutation getDownloadURL($path: String!) {
    getDownloadURL(path: $path)
  }
`;

const GQL_CREATE_SHARE = gql`
  mutation createShareFile($path: String!, $expireDays: Int!) {
    createShareFile(input: { path: $path, expireDays: $expireDays })
  }
`;

const icons = {
  DownloadOutlined,
  DeleteOutlined,
  EditOutlined,
  DragOutlined,
  FileOutlined,
  FolderOutlined,
  CopyOutlined,
  SendOutlined,
  ScissorOutlined,
  FolderOpenFilled,
};

const filterPath = (path) => {
  let paths = path.split("/");
  for (let i = 0; i < paths.length; i++) {
    if (paths[i].includes("..")) {
      paths.splice(i - 1, 2);
      i -= 2;
    }
  }
  let filtered =
    "/" +
    paths
      .filter((value, index, arr) => {
        return value != "";
      })
      .join("/");
  return filtered;
};
export default defineComponent({
  components: {
    ...icons,
    Toolbar,
  },

  setup() {
    const route = useRoute();
    const router = useRoute();
    const tableState = reactive({
      path: "/",
      files: [],
    });

    const {
      mutate: mutateListDir,
      onDone: onListDirDone,
      loading: listdiring,
    } = useMutation(GQL_LIST_DIR, () => ({
      variables: {
        path: tableState.path,
      },
      fetchPolicy: "no-cache",
    }));
    const { mutate: mutateCreateShare, onDone: onCreateShareDone } =
      useMutation(GQL_CREATE_SHARE, {
        fetchPolicy: "no-cache",
      });
    const { mutate: mutateRmdir, onDone: onRmdirDone } = useMutation(
      GQL_RMDIR,
      {
        fetchPolicy: "no-cache",
      }
    );

    const { mutate: mutateRmfile, onDone: onRmfileDone } = useMutation(
      GQL_DELETE_FILE,
      {
        fetchPolicy: "no-cache",
      }
    );

    const { mutate: mutateGetDownloadURL } = useMutation(GQL_GET_DOWNLOAD_URL, {
      fetchPolicy: "no-cache",
    });
    const chDir = (path) => {
      console.log("chdir to: ", path);
      let newpath = path.startsWith("/") ? path : tableState.path + path;
      let filtered = filterPath(newpath);
      filtered = filtered.endsWith("/") ? filtered : filtered + "/";

      mutateListDir({
        path: filtered,
      }).then(() => {
        tableState.path = filtered;
        console.log("set:", filtered);
      });
    };
    const handleDownloadFile = async (fileName) => {
      mutateGetDownloadURL({
        path: tableState.path + fileName,
      }).then(async ({ data }) => {
        message.info("下载已开始");
        if (data.getDownloadURL) {
          const response = await axios.get(data.getDownloadURL, {
            responseType: "blob",
          });
          const blob = new Blob([response.data], {
            type: "application/octet-stream",
          });
          const link = document.createElement("a");
          link.href = URL.createObjectURL(blob);
          link.download = fileName;
          link.click();
          URL.revokeObjectURL(link.href);
        }
      });
    };

    onListDirDone(({ data }) => {
      if (data.listDir) {
        tableState.files = data.listDir;
      }
    });
    onRmdirDone(({ data }) => {
      message.success("目录已删除");
      mutateListDir();
    });
    onRmfileDone(({ data }) => {
      message.success("文件已删除");
      mutateListDir();
    });
    onCreateShareDone(({ data }) => {
      Modal.info({
        title: "分享",
        content: "你的分享码是：" + data.createShareFile,
      });
    });
    return {
      tableState,
      columns,
      listdiring,
      mutateListDir,
      mutateRmdir,
      mutateRmfile,
      mutateCreateShare,
      handleDownloadFile,
      chDir,
      route,
      router,
    };
  },
  created: function () {
    if (!this.route.hash) {
      this.chDir("/");
    } else {
      console.log(this.route.hash);
      let hash = this.route.hash.substr(1);
      this.chDir(hash);
    }
  },
});
</script>

