<template>
  <a-row>
    <a-col :span="6">
      <a-statistic
        title="Username"
        :value="userInfo.name"
        style="margin-right: 0px"
      />
    </a-col>
    <a-col :span="6">
      <a-statistic title="邮箱" :value="userInfo.email" />
    </a-col>
    <a-col :span="6">
      <a-statistic title="注册日期" value="2022/4/23" />
    </a-col>
    <a-col :span="6">
      <a-statistic title="上次登录" value="2022/4/23" />
    </a-col>

    <a-col :span="6">
      <a-statistic title="总文件数" :value="5" style="margin-right: 50px" />
    </a-col>
    <a-col :span="6">
      <a-statistic title="配额" value="2048MB" />
    </a-col>
    <a-col :span="6">
      <a-statistic title="已使用" value="2048MB" />
    </a-col>
    <a-col :span="6">
      <a-statistic title="剩余" value="0MB" />
    </a-col>
  </a-row>
</template>
<script>
import { defineComponent, reactive, computed, ref } from "vue";
import { CheckOutlined, EditOutlined } from "@ant-design/icons-vue";
export default defineComponent({
  components: {
    CheckOutlined,
    EditOutlined,
  },
  props: {
    userInfo: {
      type: Object,
      default: () => ({
        name: "",
        email: "",
      }),
    },
  },
  setup() {
  },
});
</script>
