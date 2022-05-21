<template>
  <a-layout>
    <UserProfile
      @save="handleSaveProfile"
      :initialState="homeState.userInfo"
      v-model:visible="homeState.userProfileVisible"
    />
    <a-layout-sider
      id="components-layout-demo-responsive"
      :style="{
        overflow: 'auto',
        position: 'fixed',
        left: 0,
        top: 0,
        bottom: 0,
      }"
    >
      <div class="logo" />
      <a-menu
        v-model:selectedKeys="homeState.selectedKeys"
        theme="dark"
        mode="inline"
      >
        <a-menu-item key="1">
          <user-outlined />
          <span class="nav-text">账户信息</span>
        </a-menu-item>
        <a-menu-item key="2">
          <file-outlined />
          <span class="nav-text">文件管理</span>
        </a-menu-item>
        <a-menu-item key="3">
          <send-outlined />
          <span class="nav-text">分享管理</span>
        </a-menu-item>
        <a-menu-item key="4">
          <question-circle-outlined />
          <span class="nav-text">用户手册</span>
        </a-menu-item>
      </a-menu>
    </a-layout-sider>
    <a-layout :style="{ minHeight: '100vh', marginLeft: '200px' }">
      <a-layout-header
        :style="{ background: '#fff', marginLeft: '16px', paddingLeft: '24px' }"
      >
        <span style="margin-right: 10px"><database-outlined /></span>
        <span>存储服务系统</span>
        <span style="margin-left: 10px; float: right">
          <span
            >{{
              homeState.userInfo
                ? "你好，" + homeState.userInfo.name
                : "请先登录"
            }}
            !</span
          >
          <a-tooltip title="Edit Profile">
            <form-outlined
              class="user-op-icon"
              @click="homeState.userProfileVisible = true"
            />
          </a-tooltip>
          <a-tooltip title="Logout">
            <logout-outlined class="user-op-icon" @click="handleLogout" />
          </a-tooltip>
        </span>
      </a-layout-header>
      <a-layout-content :style="{ margin: '24px 16px 0' }">
        <div :style="{ padding: '24px', background: '#fff' }">
          <Account
            :userInfo="homeState.userInfo"
            v-if="homeState.selectedKeys.includes('1')"
          />
          <Table v-if="homeState.selectedKeys.includes('2')" />
        </div>
      </a-layout-content>
      <Footer />
    </a-layout>
  </a-layout>
</template>

<style scoped>
#components-layout-demo-responsive .logo {
  height: 32px;
  background: rgba(255, 255, 255, 0.2);
  margin: 16px;
}

.site-layout-sub-header-background {
  background: #fff;
}

.site-layout-background {
  background: #fff;
}

.user-op-icon {
  width: 20px;
  margin-left: 12px;
  cursor: pointer;
}
.user-op-icon:hover {
  color: #108ee9;
}

[data-theme="dark"] .site-layout-sub-header-background {
  background: #141414;
}
</style>

<script>
import {
  UserOutlined,
  SendOutlined,
  FileOutlined,
  FormOutlined,
  LogoutOutlined,
  DatabaseOutlined,
  QuestionCircleOutlined,
} from "@ant-design/icons-vue";
import { defineComponent, reactive, ref, h } from "vue";

import Footer from "../components/Footer.vue";
import Table from "../components/filebrowser/Table.vue";
import Account from "../components/Account.vue";
import UserProfile from "../components/UserProfile.vue";
import { useRouter } from "vue-router";
import { notification, Button } from "ant-design-vue";
export default defineComponent({
  components: {
    Table,
    Account,
    Footer,
    UserProfile,
    UserOutlined,
    QuestionCircleOutlined,
    FileOutlined,
    LogoutOutlined,
    FormOutlined,
    SendOutlined,
    DatabaseOutlined,
  },
  beforeRouteEnter(to, from, next) {
    console.log("home beforeRouteEnter");
    if (!localStorage.getItem("user")) {
      return next("/");
    }
    next(vm => {
      vm.homeState.userInfo = JSON.parse(localStorage.getItem("user"));
    });
  },
  setup() {
    console.log("home setup")
    const router = useRouter();
    const homeState = reactive({
      userInfo: {
        name: "",
        email: "",
      },
      userProfileVisible: false,
      selectedKeys: ["1"],
    });
    const onCollapse = (collapsed, type) => {
      console.log(collapsed, type);
    };

    const onBreakpoint = (broken) => {
      console.log(broken);
    };


    const handleLogout = () => {
      localStorage.clear();
      router.replace("/login");
    };
    const handleSaveProfile = (state) => {
      console.log("save:", state);
      homeState.userProfileVisible = false;
      homeState.userInfo.name = state.name;
      homeState.userInfo.email = state.email;
      if (state.newPass) {
        notification["warning"]({
          message: "正在登出",
          description:
            "你刚才更新了密码，稍后需要重新登录",
          onClose: handleLogout,
        });
      }
    };
    return {
      homeState,
      onCollapse,
      onBreakpoint,
      handleLogout,
      handleSaveProfile,
      router,
    };
  },
});
</script>

<script setup>
</script>
