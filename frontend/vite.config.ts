import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

// The build lands in web/dist, which the Go binary embeds. Asset URLs are
// relative: the Go server writes a <base href> with the panel's (secret,
// configurable) base path into index.html, so the same build works under any
// base path. Everything goes under static/ to stay clear of the old /assets.
export default defineConfig(({ command }) => ({
  base: command === 'build' ? './' : '/',
  plugins: [vue()],
  resolve: {
    alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) },
  },
  build: {
    outDir: '../web/dist',
    emptyOutDir: true,
    assetsDir: 'static',
    target: 'es2022',
    chunkSizeWarningLimit: 600,
  },
}))
