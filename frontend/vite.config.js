// frontend/vite.config.js
import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

export default defineConfig({
  plugins: [vue()],
  server: {
    port: 5173,
    strictPort: true,
    proxy: {
      // API routes -> Go backend
      "/v1": {
        target: "http://localhost:3000",
        changeOrigin: true
      },
      "/health": {
        target: "http://localhost:3000",
        changeOrigin: true
      }
    }
  }
});