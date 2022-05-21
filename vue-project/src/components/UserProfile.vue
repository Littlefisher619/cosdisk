<template>
  <div id="components-form-demo-normal-register">
    <a-modal
      v-model:visible="modalVisible"
      :title="title"
      :rules="rules"
      @ok="onOk"
      :confirmLoading="confirmLoading"
    >
      <a-form
        :model="formState"
        :rules="rules"
        ref="formRef"
        name="update_profile"
        class="register-form"
        v-bind="{
          labelCol: {
            span: 8,
          },
          labelAlign: 'left',
        }"
      >
       <a-form-item label="邮箱" name="email" required>
          <a-tooltip title="邮箱不允许更改">
            <span>{{ formState.email }}</span>
          </a-tooltip>
        </a-form-item>
        <a-form-item label="用户名" name="name" required>
          <a-input v-model:value="formState.name" />
        </a-form-item>
        <a-form-item has-feedback label="新密码" name="newPass">
          <a-input
            v-model:value="formState.newPass"
            type="password"
            placeholder="留空则不更改密码"
            autocomplete="off"
          />
        </a-form-item>

        <a-form-item
          has-feedback
          v-if="formState.newPass != ''"
          label="Confirm"
          name="checkPass"
        >
          <a-input
            v-model:value="formState.checkPass"
            type="password"
            autocomplete="off"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>
<script>
"use strict";
import { defineComponent, reactive, ref } from "vue";
import { message } from "ant-design-vue";
export default defineComponent({
  props: {
    title: {
      type: String,
      default: "更新用户资料",
    },
    visible: {
      type: Boolean,
      default: false,
    },
    confirmLoading: {
      type: Boolean,
      default: false,
    },
    initialState: {
      type: Object,
      default: () => ({
        name: "",
        email: "",
      }),
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
    save: null,
    "update:visible": null,
  },

  setup(props, { emit }) {
    const formRef = ref();
    const formState = reactive({
      checkPass: "",
      newPass: "",
      email: "youremail@email.com",
      name: "username",
    });
    let validatePass = async (_rule, value) => {
      if (value === "") {
        formState.checkPass = "";
      } else {
        if (formState.checkPass !== "") {
          formRef.value.validateFields("checkPass");
        }

        return Promise.resolve();
      }
    };

    let validatePass2 = async (_rule, value) => {
      if (value === "") {
        return Promise.reject("请再输入一遍密码");
      } else if (value !== formState.newPass) {
        return Promise.reject("重复输入的密码不匹配");
      } else {
        return Promise.resolve();
      }
    };

    const onOk = () => {
      formRef.value
        .validate()
        .then(async (valid) => {
          if (valid) {
            emit("save", formState);
          }
        })
        .catch(() => {
          message.error("请检查你的输入");
        });
    };

    const rules = {
      newPass: [
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
          trigger: "change",
        },
      ],
    };
    return {
      formRef,
      formState,

      rules,
      onOk,
    };
  },
  watch: {
    visible(newVal, oldVal) {
      this.formState.newPass = "";
      this.formState.checkPass = "";
      this.formState.email = this.initialState.email;
      this.formState.name = this.initialState.name;
    },
  },
});
</script>

<style scoped>
#components-form-demo-normal-register {
  margin: 0 auto;
  display: flex;

  justify-content: space-around;
}
#components-form-demo-normal-register .register-form {
  max-width: 300px;
}
#components-form-demo-normal-register .register-form-wrap {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
</style>