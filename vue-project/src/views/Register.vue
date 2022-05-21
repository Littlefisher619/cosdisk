<template>
  <a-layout :style="{ minHeight: '100vh' }">
    <a-layout-header>欢迎来到存储服务系统</a-layout-header>
    <a-layout-content>
      <div id="components-form-demo-normal-register" class="register">
        <a-typography-title>注册</a-typography-title>
        <a-divider />
        <a-form
          :model="formState"
          :rules="rules"
          ref="formRef"
          name="normal_register"
          class="register-form"
          v-bind="layout"
          @finish="onFinish"
          @finishFailed="onFinishFailed"
        >
          <a-form-item label="邮箱" name="email" has-feedback>
            <a-input v-model:value="formState.email" />
          </a-form-item>

          <a-form-item has-feedback label="用户名" name="username">
            <a-input v-model:value="formState.username" autocomplete="off" />
          </a-form-item>

          <a-form-item has-feedback label="密码" name="pass">
            <a-input
              v-model:value="formState.pass"
              type="password"
              autocomplete="off"
            />
          </a-form-item>

          <a-form-item has-feedback label="确认密码" name="checkPass">
            <a-input
              v-model:value="formState.checkPass"
              type="password"
              autocomplete="off"
            />
          </a-form-item>

          <div class="register-form-wrap">
            <a-form-item style="width: 100%">
              <a-button
                type="primary"
                html-type="submit"
                style="width: 100%"
                :loading="loading"
                >提交</a-button
              >
              或者
              <router-link to="/login">去登录！</router-link>
            </a-form-item>
          </div>
        </a-form>
      </div>
    </a-layout-content>
    <Footer />
  </a-layout>
</template>
<script>
"use strict";
import { defineComponent, reactive, ref } from "vue";
import gql from "graphql-tag";
import { useRouter } from "vue-router";
import { useQuery, useMutation, useResult } from "@vue/apollo-composable";
import { message } from "ant-design-vue";
import Footer from "../components/Footer.vue";
const GQL_REGISTER = gql`
  mutation register($name: String!, $email: String!, $password: String!) {
    register(input: { name: $name, email: $email, password: $password }) {
      token
      user {
        id
        name
        email
      }
    }
  }
`;

export default defineComponent({
  components: {
    Footer,
  },
  setup() {
    const formRef = ref();
    const formState = reactive({
      pass: "password",
      checkPass: "password",
      email: "youremail@email.com",
      username: "username",
    });
    const userRegex = /^[a-zA-Z0-9_]+$/;
    const router = useRouter();

    let validatePass = async (_rule, value) => {
      if (value === "") {
        return Promise.reject("请输入密码");
      } else {
        if (formState.username !== "" && value.includes(formState.username)) {
          return Promise.reject("密码不能包含用户名");
        }
        if (formState.checkPass !== "") {
          formRef.value.validateFields("checkPass");
        }

        return Promise.resolve();
      }
    };

    let validatePass2 = async (_rule, value) => {
      if (value === "") {
        return Promise.reject("请再输入一遍密码");
      } else if (value !== formState.pass) {
        return Promise.reject("重复输入的密码不匹配");
      } else {
        return Promise.resolve();
      }
    };

    let validateUserName = async (_rule, value) => {
      if (value === "") {
        return Promise.reject("请输入用户名");
      } else if (value.length < 5) {
        return Promise.reject("用户名至少要有五个字符");
      } else if (!userRegex.test(value)) {
        return Promise.reject("用户名不符合规范");
      } else {
        return Promise.resolve();
      }
    };

    const rules = {
      pass: [
        {
          required: true,
          validator: validatePass,
          trigger: "change",
        },
      ],
      checkPass: [
        {
          required: true,
          validator: validatePass2,
          trigger: "change",
        },
      ],
      username: [
        {
          required: true,
          validator: validateUserName,
          trigger: "change",
        },
      ],
      email: [
        {
          required: true,
          trigger: "change",
        },
        {
          type: "email",
          trigger: "change",
        },
      ],
    };
    const layout = {
      labelCol: {
        span: 8,
      },
      labelAlign: "left",
    };

    const {
      mutate: doRegister,
      onDone,
      onError,
      loading,
    } = useMutation(GQL_REGISTER, () => ({
      variables: {
        name: formState.username,
        email: formState.email,
        password: formState.pass,
      },
    }));

    const onFinish = (values) => {
      doRegister();
      console.log("Success:", values);
    };
    onDone(({ data }) => {
      console.log(data);
      if (data.register) {
        message.success("注册成功，正在跳转...");
        localStorage.clear();
        localStorage.setItem("token", data.register.token);
        localStorage.setItem("user", JSON.stringify(data.register.user));
        router.replace("/home");
      }
    });
    onError((error) => {
      let msg = error.message;

      if (msg.includes("User already register"))
        message.error("注册失败：用户名或邮箱已经被注册");
      else message.error("注册失败：未知错误");
    });
    const onFinishFailed = (errors) => {
      message.error("请检查输入");
    };

    return {
      formState,
      formRef,
      rules,
      layout,
      onFinishFailed,
      onFinish,
      loading,
    };
  },
});
</script>


<script setup>
</script>

<style scoped>
.register {
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

#components-form-demo-normal-register .register-form {
  max-width: 300px;
}
#components-form-demo-normal-register .register-form-wrap {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.ant-layout-header {
  color: #fff;
}
.ant-layout-footer {
  color: rgb(0, 0, 0);
  text-align: center;
}
</style>
