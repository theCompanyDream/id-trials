import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    host: '0.0.0.0'
  },
  build: {
    outDir: 'dist/frontend'
  },
  resolve: {
    alias: {
      '@backpack': path.resolve(__dirname, './src/components/index'),
      '@backpack/*': path.resolve(__dirname, './src/components/*'),
      // If you need the wildcards to work for assets/pages too:
      '@components': path.resolve(__dirname, './src/components'),
      '@assets': path.resolve(__dirname, './src/assets'),
      '@pages': path.resolve(__dirname, './src/pages'),
    }
  }
})
