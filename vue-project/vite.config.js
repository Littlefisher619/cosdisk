import { fileURLToPath, URL } from 'url'
import pluginEnv from 'vite-plugin-vue-env';
import vue from '@vitejs/plugin-vue'
import { defineConfig, loadEnv } from 'vite';

// https://vitejs.dev/config/
export default ({ mode }) => {
  process.env = { ...process.env, ...loadEnv(mode, process.cwd()) };
  proxy = {
    "/query": {
      target: `http://127.0.0.1:8080`,
      changeOrigin: true,
      secure: false,
    },
  }
  if (mode === 'production') {
    proxy = {} 
  }
  console.log(process.env)
  return defineConfig({
    plugins: [vue(), pluginEnv()],
    server: {
      base: process.env.VITE_APP_BASE || '/',
      port: process.env.VITE_APP_PORT || 3000,
      host: process.env.VITE_APP_HOST || '127.0.0.1',
      proxy: proxy,
    },
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url))
      }
    },
  })
}
