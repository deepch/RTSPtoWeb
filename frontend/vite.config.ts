import path from "path"
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    proxy: {
      '/stream': {
        target: 'http://localhost:8083',
        changeOrigin: true,
        ws: true,
      },
      '/streams': {
        target: 'http://localhost:8083',
        changeOrigin: true,
      },
      '/login': {
        target: 'http://localhost:8083',
        changeOrigin: true,
      },
    }
  }
})
