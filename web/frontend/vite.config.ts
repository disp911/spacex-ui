import { fileURLToPath, URL } from 'node:url'
import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'

// The built app is served by the Go panel from web/ui (see web/spa.go).
// Asset URLs stay relative; the Go side injects a <base href> that points
// at "<basePath>ui/", so the same build works under any panel base path.
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  // During `npm run dev`, API calls are proxied to a running panel.
  const panel = env.PANEL_URL || 'http://127.0.0.1:2053'

  return {
    base: './',
    plugins: [vue()],
    resolve: {
      alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) },
    },
    build: {
      outDir: '../ui',
      emptyOutDir: true,
      assetsDir: 'assets',
      chunkSizeWarningLimit: 600,
    },
    server: {
      proxy: {
        '^/(login|logout|getTwoFactorEnable|panel|ws)': { target: panel, changeOrigin: true, ws: true },
      },
    },
  }
})
