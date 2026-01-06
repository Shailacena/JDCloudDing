import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'

// https://vite.dev/config/
export default defineConfig({
  base: process.env.NODE_ENV === 'production' ? '/web/' : '/', // 生产环境使用路径前缀
  plugins: [react()],
  server: {
    proxy: {
      '/web_api': {
        target: 'http://localhost:8000',
        changeOrigin: true
      },
    }
  }
})
