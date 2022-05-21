<template>
  <a-layout :style="{ minHeight: '100vh' }">
    <a-layout-header>欢迎来到存储服务系统</a-layout-header>
    <a-layout-content>
      <div id="components-form-demo-normal-login" class="login">
        <a-typography-title>登录</a-typography-title>
        <a-divider />
        <a-form
          :model="formState"
          name="normal_login"
          class="login-form"
          @finish="onFinish"
          @finishFailed="onFinishFailed"
        >
          <a-form-item
            label="邮箱"
            name="email"
            :rules="[
              { type: 'email', message: '邮箱格式无效' },
              { required: true },
            ]"
          >
            <a-input v-model:value="formState.email">
              <template #prefix>
                <MailOutlined class="site-form-item-icon" />
              </template>
            </a-input>
          </a-form-item>

          <a-form-item
            label="密码"
            name="password"
            :rules="[{ required: true }]"
          >
            <a-input-password v-model:value="formState.password">
              <template #prefix>
                <LockOutlined class="site-form-item-icon" />
              </template>
            </a-input-password>
          </a-form-item>

          <div class="login-form-wrap">
            <a-form-item name="remember" no-style>
              <a-checkbox v-model:checked="formState.remember"
                >记住我</a-checkbox
              >
            </a-form-item>
            <a class="login-form-forgot" href="">忘记密码</a>
          </div>

          <a-form-item>
            <a-button
              :disabled="disabled"
              :loading="loading"
              type="primary"
              html-type="submit"
              class="login-form-button"
            >
              提交
            </a-button>
            或者
            <router-link to="/register">立即注册！</router-link>
          </a-form-item>
        </a-form>
      </div>
    </a-layout-content>
    <Footer />
  </a-layout>
</template>

<script>
"use strict";
import gql from "graphql-tag";
import { useQuery, useMutation, useResult } from "@vue/apollo-composable";
const GQL_LOGIN = gql`
  mutation login($email: String!, $password: String!) {
    login(email: $email, password: $password) {
      token
      user {
        id
        name
        email
      }
    }
  }
`;
import { defineComponent, reactive, computed } from "vue";
import { useRouter } from "vue-router";
import { MailOutlined, LockOutlined } from "@ant-design/icons-vue";
import { message } from "ant-design-vue";
import Footer from "../components/Footer.vue";

export default defineComponent({
  components: {
    MailOutlined,
    LockOutlined,
    Footer,
  },

  setup() {
    const router = useRouter();
    const formState = reactive({
      email: "test@test.com",
      password: "test",
      remember: true,
    });
    const {
      mutate: doLogin,
      onDone,
      onError,
      loading,
    } = useMutation(GQL_LOGIN, () => ({
      variables: {
        email: formState.email,
        password: formState.password,
      },
    }));

    // const apollo = {}
    const onFinish = (values) => {
      doLogin();
      console.log("Success:", values);
    };
    onError((error) => {
      let msg = error.message;

      if (msg.includes("User not found")) message.error("登录失败: 用户不存在");
      else if (msg.includes("Password incorrect"))
        message.error("登录失败: 密码错误");
      else message.error("登录失败: 未知错误");
    });
    onDone(({ data }) => {
      if (data.login) {
        message.success("登录成功！", 1);
        localStorage.clear();
        localStorage.setItem("token", data.login.token);
        localStorage.setItem("user", JSON.stringify(data.login.user));

        router.replace("/home");
      }
    });

    const onFinishFailed = (errorInfo) => {
      console.log("Failed:", errorInfo);
    };

    const disabled = computed(() => {
      return !(formState.email && formState.password);
    });

    return {
      formState,
      onFinish,
      onFinishFailed,
      disabled,
      loading,
    };
  },
});
</script>

<style scoped>
.login {
  width: 380px;
  background: #f9f9f9;
  display: block;
  margin: 100px auto 0px auto;
  /*	clockwise margin setup */
  /*	border-bottom: 4px solid purple;*/
  border-radius: 4px;
  padding-top: 18px;
  padding-left: 25px;
  padding-right: 29px;
  padding-bottom: 60px;
  /*		box-shadow: 0px 0px 2px red; */
  box-shadow: 0px 0px 3px #dadada;
  border: 2px solid #dadada;
}

#components-form-demo-normal-login .login-form {
  max-width: 300px;
}
#components-form-demo-normal-login .login-form-wrap {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

#components-form-demo-normal-login .login-form-button {
  width: 100%;
}
.ant-layout-header {
  color: #fff;
}
.ant-layout-footer {
  color: rgb(0, 0, 0);
  text-align: center;
}
</style>

