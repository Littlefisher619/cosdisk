
<template>
    <a-config-provider :locale="locale">
    <router-view v-slot="{ Component }">
      <transition name="fade" >
        <keep-alive>
          
            <component :is="Component"/>
          
        </keep-alive>
      </transition>
    </router-view>
    </a-config-provider>
    <!--    <Login/>-->
</template>
<script>
  import zhCN from 'ant-design-vue/es/locale/zh_CN';
  import dayjs from 'dayjs';
  import 'dayjs/locale/zh-cn';
  dayjs.locale('zh-cn');
  export default {
    data() {
      return {
        locale: zhCN,
      };
    },
  };
</script>

<style>
html, body, #app {
  /* overflow: hidden; */
  height:100%;
}
</style>