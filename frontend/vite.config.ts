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
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          'react-vendor': ['react', 'react-dom', 'react-router-dom'],
          'ui-vendor': ['lucide-react', 'class-variance-authority', 'clsx', 'tailwind-merge'],
          'hls': ['hls.js']
        }
      }
    },
    chunkSizeWarningLimit: 600,
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
