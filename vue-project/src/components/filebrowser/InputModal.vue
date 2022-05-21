<template>
  <a-modal
    v-model:visible="modalVisible"
    :title="title"
    @ok="$emit('ok', modalState.ret)"
    :confirmLoading="confirmLoading"
  >
    <a-input
      v-model:value="modalState.ret"
      :placeholder="description"
      @pressEnter="$emit('ok', modalState.ret)"
    />
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
      default: "请输入...",
    },
    description: {
      type: String,
      default: "",
    },
    visible: {
      type: Boolean,
      default: false,
    },
    confirmLoading: {
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
    ok: null,
    "update:visible": null,
  },
  components: {
    InboxOutlined,
  },
  watch: {
    visible(newVal, oldVal) {
      this.modalState.ret = "";
    },
  },
  setup() {
    const modalState = reactive({
      ret: "",
    });
    return {
      modalState,
    };
  },
});
</script>
