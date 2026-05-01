import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import Sitemap from 'vite-plugin-sitemap' // 1. 引入插件

const cfAsyncModuleScriptPlugin = () => ({
  name: 'cfasync-module-script',
  transformIndexHtml(html: string) {
    return html.replace(
      /<script\s+type="module"(?![^>]*data-cfasync)/g,
      '<script data-cfasync="false" type="module"',
    )
  },
})

// https://vite.dev/config/
export default defineConfig(({ mode }) => ({
  plugins: [
    vue(), cfAsyncModuleScriptPlugin(),
    Sitemap({ 
      hostname: 'https://www.rsk.cc.cd', // 替换成你真实的域名
      // 根据项目路由补全公开页面
      generateRobotsTxt: false, 
      // 明确指定你要索引的路径
      dynamicRoutes: [
        '/',
        '/products',
        '/cart',
        '/checkout',
        '/pay',
        '/guest/orders',
        '/blog',
        '/notice',
        '/about',
        '/terms',
        '/privacy',
        '/auth/login',
        '/auth/register',
        '/auth/forgot'
      ],
      changefreq: 'weekly', // 更新频率
      priority: 1.0,        // 权重
    }),
  ],
  esbuild: mode === 'production' ? { drop: ['console', 'debugger'] } : {},
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          'vendor-qrcode': ['qrcode'],
          'vendor-vue-i18n': ['vue-i18n'],
        },
      },
    },
  },
  server: {
    host: '0.0.0.0', // 监听所有网络接口
    port: 5173,
    strictPort: true,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/uploads': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      }
    }
  },
}))
