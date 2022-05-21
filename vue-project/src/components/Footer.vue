<template>
  <a-layout-footer style="text-align: center">
    <github-outlined style="margin: 0 10px" />
    <a href="https://www.github.com/Littlefisher619/cosdisk"
      >Github Project CosDisk</a
    >
  </a-layout-footer>
</template>

<script>
import {
  GithubOutlined,
} from "@ant-design/icons-vue";
export default {
  name: "AppFooter",
  components: {
    GithubOutlined,
  },
  setup() {
    return {};
  },
};
</script>>