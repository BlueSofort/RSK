<template>
  <footer
    class="relative theme-panel-strong theme-text-secondary border-t theme-border overflow-hidden">
    <div class="container mx-auto px-4 pt-10 pb-7 relative">
      <div class="grid grid-cols-2 md:grid-cols-4 gap-8 md:gap-12 mb-16">
        <!-- Brand -->
        <div class="col-span-2 space-y-6">
          <div class="flex items-center space-x-3">
            <div
              class="w-8 h-8 theme-btn-primary rounded-lg flex items-center justify-center">
              <span class="text-white font-black text-sm">{{ brandInitial }}</span>
            </div>
            <h3 class="theme-text-primary text-xl font-bold tracking-tight">{{ brandSiteName }}</h3>
          </div>
          <p class="text-sm leading-relaxed max-w-sm theme-text-muted">
            {{ brandDescription || '专业的在线发卡工作室，为您提供优质的产品和服务。' }}
          </p>
          <div class="flex space-x-4">
            <!-- Social Icons (Placeholder) -->
            <!--
            <a href="#" class="w-10 h-10 rounded-full bg-white/5 flex items-center justify-center hover:bg-white/10 hover:text-white transition-colors">
              <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24"><path d="M24 4.557c-.883.392-1.832.656-2.828.775 1.017-.609 1.798-1.574 2.165-2.724-.951.564-2.005.974-3.127 1.195-.897-.957-2.178-1.555-3.594-1.555-3.179 0-5.515 2.966-4.797 6.045-4.091-.205-7.719-2.165-10.148-5.144-1.29 2.213-.669 5.108 1.523 6.574-.806-.026-1.566-.247-2.229-.616-.054 2.281 1.581 4.415 3.949 4.89-.693.188-1.452.232-2.224.084.626 1.956 2.444 3.379 4.6 3.419-2.07 1.623-4.678 2.348-7.29 2.04 2.179 1.397 4.768 2.212 7.548 2.212 9.142 0 14.307-7.721 13.995-14.646.962-.695 1.797-1.562 2.457-2.549z"/></svg>
            </a>
            <a href="#" class="w-10 h-10 rounded-full bg-white/5 flex items-center justify-center hover:bg-white/10 hover:text-white transition-colors">
              <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>
            </a>
            -->
          </div>
        </div>

        <!-- Links -->
        <div>
          <h4 class="theme-text-primary font-bold mb-6 tracking-wide">{{ t('footer.quickLinks') }}</h4>
          <ul class="space-y-3 text-sm">
            <li v-for="item in quickLinks" :key="item.path">
              <router-link :to="item.path" class="theme-link-muted transition-colors flex items-center gap-2 group">
                <svg class="w-4 h-4 shrink-0 opacity-50 group-hover:opacity-80 transition-opacity" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.75" :d="item.icon" />
                </svg>
                {{ t(item.label) }}
              </router-link>
            </li>
          </ul>
        </div>

        <!-- Contact -->
        <div>
          <h4 class="theme-text-primary font-bold mb-6 tracking-wide">{{ t('footer.contact') }}</h4>
          <div class="space-y-4">
            <a v-if="config?.contact?.telegram" :href="config.contact.telegram" target="_blank"
              rel="noopener noreferrer"
              class="flex items-center space-x-3 text-sm hover:text-gray-900 dark:hover:text-white transition-colors p-3 rounded-lg bg-gray-50 dark:bg-white/5 hover:bg-gray-100 dark:hover:bg-white/10 border border-gray-100 dark:border-white/5 hover:border-gray-200 dark:hover:border-white/10">
              <svg class="w-5 h-5 text-blue-400" fill="currentColor" viewBox="0 0 24 24">
                <path
                  d="M12 0C5.373 0 0 5.373 0 12s5.373 12 12 12 12-5.373 12-12S18.627 0 12 0zm5.894 8.221l-1.97 9.28c-.145.658-.537.818-1.084.508l-3-2.21-1.446 1.394c-.14.18-.357.295-.6.295-.002 0-.003 0-.005 0l.213-3.054 5.56-5.022c.24-.213-.054-.334-.373-.121l-6.869 4.326-2.96-.924c-.64-.203-.658-.64.135-.954l11.566-4.458c.538-.196 1.006.128.832.941z" />
              </svg>
              <span>Telegram</span>
            </a>
            <a v-if="config?.contact?.whatsapp" :href="config.contact.whatsapp" target="_blank"
              rel="noopener noreferrer"
              class="flex items-center space-x-3 text-sm hover:text-gray-900 dark:hover:text-white transition-colors p-3 rounded-lg bg-gray-50 dark:bg-white/5 hover:bg-gray-100 dark:hover:bg-white/10 border border-gray-100 dark:border-white/5 hover:border-gray-200 dark:hover:border-white/10">
              <svg class="w-5 h-5 text-green-400" fill="currentColor" viewBox="0 0 24 24">
                <path
                  d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413Z" />
              </svg>
              <span>WhatsApp</span>
            </a>
          </div>
        </div>
      </div>

      <!-- Copyright -->
      <div
        class="border-t theme-border pt-8 flex flex-col items-center gap-6 text-xs theme-text-muted">
        <div class="flex flex-col items-center gap-4 w-full">
          <div class="flex flex-col md:flex-row items-center justify-between w-full gap-4">
            <div class="flex flex-wrap items-center justify-center md:justify-start gap-x-4 gap-y-1">
              <p>&copy; {{ currentYear }} {{ brandSiteName }}. {{ t('footer.rights') }}</p>
            </div>
            
            <div class="flex flex-wrap items-center gap-x-4 gap-y-1 justify-center md:justify-end">
              <router-link to="/privacy" class="hover:text-gray-900 dark:hover:text-gray-400">{{ t('footer.privacy') || 'Privacy Policy' }}</router-link>
              <router-link to="/terms" class="hover:text-gray-900 dark:hover:text-gray-400">{{ t('footer.terms') || 'Terms of Service' }}</router-link>
              <template v-if="footerLinks.length">
                <a
                  v-for="link in footerLinks"
                  :key="link.name"
                  :href="link.url || 'javascript:void(0)'"
                  :target="link.url ? '_blank' : undefined"
                  rel="noopener noreferrer"
                  class="hover:text-gray-900 dark:hover:text-gray-400"
                >{{ link.name }}</a>
              </template>
            </div>
          </div>

          <!-- Status Badges (Theme-aware style) -->
          <div class="flex items-center justify-center gap-1">
            <!-- Left Badge: Brand -->
            <div class="flex items-center bg-[#2d2d2d] dark:bg-[#f5f5f5] text-white dark:text-[#1a1a1a] rounded px-1.5 py-0.5 shadow-md border border-white/10 dark:border-black/5 transition-colors duration-300">
              <span class="mr-1 flex items-center justify-center text-[10px]">
                🍰
              </span>
              <span class="text-[10px] font-bold tracking-tight">RSK发卡小店</span>
            </div>
            
            <!-- Right Badge: Status Bubble -->
            <div class="relative flex items-center bg-[#2d2d2d] dark:bg-[#f5f5f5] text-white dark:text-[#1a1a1a] px-1.5 py-0.5 rounded shadow-md border border-white/10 dark:border-black/5 transition-colors duration-300">
              <!-- Bubble Arrow pointing left -->
              <div class="absolute -left-0.5 top-1/2 -translate-y-1/2 w-1 h-1 bg-[#2d2d2d] dark:bg-[#f5f5f5] rotate-45 border-l border-b border-white/10 dark:border-black/5"></div>
              <span class="relative z-10 text-[10px] font-bold">营业中</span>
            </div>
          </div>

          <!-- Transistor Uptime Clock (Digital style) -->
          <div class="inline-flex items-center px-3 py-1 bg-black dark:bg-[#1a1a1a] rounded-sm border border-cyan-500/40 dark:border-cyan-400/60 shadow-[0_0_8px_rgba(6,182,212,0.3)] dark:shadow-[0_0_12px_rgba(6,182,212,0.5)] transition-colors duration-300">
            <span class="font-digital text-[12px] font-bold tracking-[0.05em] text-[#00ffff] dark:text-[#00e5ff] drop-shadow-[0_0_4px_rgba(0,255,255,0.8)] dark:drop-shadow-[0_0_6px_rgba(0,229,255,0.9)] uppercase">
              {{ uptimeDisplay }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '../stores/app'

const { t } = useI18n()
const appStore = useAppStore()

const config = computed(() => appStore.config)

// Uptime Clock Logic
const launchDate = new Date('2026-04-18T00:00:00') // 运行起始日期
const uptimeDisplay = ref('')

const updateUptime = () => {
  const now = new Date()
  const diff = now.getTime() - launchDate.getTime()

  const years = Math.floor(diff / (1000 * 60 * 60 * 24 * 365))
  const days = Math.floor((diff % (1000 * 60 * 60 * 24 * 365)) / (1000 * 60 * 60 * 24))
  const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
  const seconds = Math.floor((diff % (1000 * 60)) / 1000)

  const pad = (n: number) => n.toString().padStart(2, '0')

  uptimeDisplay.value = `${years} YEAR ${days} DAYS ${pad(hours)}:${pad(minutes)}:${pad(seconds)}`
}

let timer: number | undefined

onMounted(() => {
  updateUptime()
  timer = window.setInterval(updateUptime, 1000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})

const brandSiteName = computed(() => {
  const siteName = config.value?.brand?.site_name
  return typeof siteName === 'string' && siteName.trim() ? siteName.trim() : 'Dujiao-Next'
})

const brandDescription = computed(() => {
  const desc = config.value?.brand?.site_description
  if (desc && typeof desc === 'object') {
    const val = desc[appStore.locale] || desc['zh-CN'] || ''
    return typeof val === 'string' ? val.trim() : ''
  }
  return ''
})

const brandInitial = computed(() => {
  return brandSiteName.value.charAt(0).toUpperCase()
})

const isListMode = computed(() => config.value?.template_mode === 'list')
const navConfig = computed(() => config.value?.nav_config as { builtin?: Record<string, boolean> } | undefined)

const quickLinks = computed(() => {
  const items = [
    { path: '/', label: 'nav.home', icon: 'M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-4 0a1 1 0 01-1-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 01-1 1' },
  ]
  if (!isListMode.value) {
    items.push({ path: '/products', label: 'nav.products', icon: 'M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z' })
  }
  const builtin = navConfig.value?.builtin
  if (!builtin || builtin.blog !== false) {
    items.push({ path: '/blog', label: 'nav.blog', icon: 'M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2' })
  }
  if (!builtin || builtin.about !== false) {
    items.push({ path: '/about', label: 'nav.about', icon: 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z' })
  }
  return items
})

const footerLinks = computed(() => {
  const links = config.value?.footer_links
  if (!Array.isArray(links)) return []
  return links.filter((item: any) => item && typeof item.name === 'string' && item.name.trim())
})

const currentYear = new Date().getFullYear()
</script>
