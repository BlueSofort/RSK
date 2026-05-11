<template>
  <div class="search-page min-h-screen theme-page pt-24 pb-16">
    <div class="container mx-auto px-4 max-w-4xl">
      <!-- Search Header -->
      <div class="mb-10">
        <div class="flex items-center gap-3 max-w-xl mx-auto">
          <!-- Dropdown -->
          <div class="relative">
            <select
              v-model="activeTab"
              class="appearance-none h-10 pl-4 pr-9 rounded-full border theme-border bg-background text-sm font-medium focus:outline-none focus:ring-2 focus:ring-primary/20 cursor-pointer"
            >
              <option value="all">{{ t('search.tabAll') }}</option>
              <option value="products">{{ t('search.tabProducts') }}</option>
              <option value="posts">{{ t('search.tabPosts') }}</option>
            </select>
            <svg class="absolute right-3 top-1/2 -translate-y-1/2 w-3.5 h-3.5 theme-text-muted pointer-events-none" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
          </div>

          <!-- Search Input -->
          <div class="relative flex-1">
            <svg class="absolute left-3.5 top-1/2 -translate-y-1/2 w-4 h-4 theme-text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
            <input
              ref="searchInput"
              v-model="query"
              type="text"
              :placeholder="t('search.placeholder')"
              class="w-full h-10 pl-10 pr-10 rounded-full border theme-border bg-background text-sm focus:outline-none focus:ring-2 focus:ring-primary/20"
              @keyup.enter="doSearch"
            />
            <button
              v-if="query"
              @click="query = ''; doSearch()"
              class="absolute right-2 top-1/2 -translate-y-1/2 p-1.5 rounded-full hover:bg-muted"
            >
              <svg class="w-3.5 h-3.5 theme-text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="flex justify-center py-20">
        <div class="w-8 h-8 border-2 border-primary border-t-transparent rounded-full animate-spin"></div>
      </div>

      <!-- All Results -->
      <template v-else-if="activeTab === 'all'">
        <!-- Products -->
        <section v-if="products.length > 0" class="mb-12">
          <h2 class="text-lg font-bold theme-text-primary mb-5 flex items-center gap-2">
            <span class="w-1 h-5 theme-accent-stick rounded-full"></span>
            {{ t('search.tabProducts') }} <span class="text-sm font-normal theme-text-muted">({{ products.length }})</span>
          </h2>
          <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
            <div
              v-for="product in products"
              :key="product.id"
              class="theme-panel border rounded-2xl overflow-hidden cursor-pointer hover:shadow-lg hover:-translate-y-1 transition-all duration-300"
              @click="goToProduct(product.slug)"
            >
              <div class="h-36 bg-muted relative overflow-hidden">
                <img
                  v-if="getFirstImageUrl(product.images)"
                  :src="getFirstImageUrl(product.images)"
                  class="w-full h-full object-cover"
                  loading="lazy"
                />
                <div v-else class="w-full h-full flex items-center justify-center theme-text-muted">
                  <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                  </svg>
                </div>
              </div>
              <div class="p-3">
                <div class="text-sm font-semibold theme-text-primary line-clamp-2">{{ getLocalizedText(product.title) }}</div>
                <div class="mt-1 text-sm font-bold theme-text-accent">{{ formatMoney(product.price_amount) }}</div>
              </div>
            </div>
          </div>
        </section>

        <!-- Posts -->
        <section v-if="posts.length > 0">
          <h2 class="text-lg font-bold theme-text-primary mb-5 flex items-center gap-2">
            <span class="w-1 h-5 theme-accent-stick rounded-full"></span>
            {{ t('search.tabPosts') }} <span class="text-sm font-normal theme-text-muted">({{ posts.length }})</span>
          </h2>
          <div class="space-y-3">
            <article
              v-for="post in posts"
              :key="post.id"
              class="theme-panel border rounded-2xl p-5 cursor-pointer hover:shadow-lg transition-all duration-300"
              @click="goToPost(post.slug)"
            >
              <div class="flex items-center gap-3 mb-2">
                <span class="theme-badge-meta text-xs" :class="post.type === 'blog' ? 'theme-badge-accent' : 'theme-badge-info'">
                  {{ post.type === 'blog' ? t('nav.blog') : t('nav.notice') }}
                </span>
                <time class="text-xs theme-text-muted">{{ formatDate(post.published_at || post.created_at) }}</time>
              </div>
              <h2 class="text-base font-bold theme-text-primary mb-1.5">{{ getLocalizedText(post.title) }}</h2>
              <p class="text-sm theme-text-secondary line-clamp-2">{{ getLocalizedText(post.summary) }}</p>
            </article>
          </div>
        </section>

        <div v-if="searched && products.length === 0 && posts.length === 0" class="text-center py-20">
          <svg class="w-16 h-16 mx-auto theme-text-muted mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <p class="theme-text-muted text-lg">{{ t('search.noResults') }}</p>
        </div>
      </template>

      <!-- Products Only -->
      <div v-else-if="activeTab === 'products'">
        <div v-if="products.length > 0" class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
          <div
            v-for="product in products"
            :key="product.id"
            class="theme-panel border rounded-2xl overflow-hidden cursor-pointer hover:shadow-lg hover:-translate-y-1 transition-all duration-300"
            @click="goToProduct(product.slug)"
          >
            <div class="h-36 bg-muted relative overflow-hidden">
              <img
                v-if="getFirstImageUrl(product.images)"
                :src="getFirstImageUrl(product.images)"
                class="w-full h-full object-cover"
                loading="lazy"
              />
              <div v-else class="w-full h-full flex items-center justify-center theme-text-muted">
                <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
              </div>
            </div>
            <div class="p-3">
              <div class="text-sm font-semibold theme-text-primary line-clamp-2">{{ getLocalizedText(product.title) }}</div>
              <div class="mt-1 text-sm font-bold theme-text-accent">{{ formatMoney(product.price_amount) }}</div>
            </div>
          </div>
        </div>
        <div v-else-if="searched" class="text-center py-20">
          <p class="theme-text-muted text-lg">{{ t('search.noProducts') }}</p>
        </div>
      </div>

      <!-- Posts Only -->
      <div v-else-if="activeTab === 'posts'">
        <div v-if="posts.length > 0" class="space-y-3">
          <article
            v-for="post in posts"
            :key="post.id"
            class="theme-panel border rounded-2xl p-5 cursor-pointer hover:shadow-lg transition-all duration-300"
            @click="goToPost(post.slug)"
          >
            <div class="flex items-center gap-3 mb-2">
              <span class="theme-badge-meta text-xs" :class="post.type === 'blog' ? 'theme-badge-accent' : 'theme-badge-info'">
                {{ post.type === 'blog' ? t('nav.blog') : t('nav.notice') }}
              </span>
              <time class="text-xs theme-text-muted">{{ formatDate(post.published_at || post.created_at) }}</time>
            </div>
            <h2 class="text-base font-bold theme-text-primary mb-1.5">{{ getLocalizedText(post.title) }}</h2>
            <p class="text-sm theme-text-secondary line-clamp-2">{{ getLocalizedText(post.summary) }}</p>
          </article>
        </div>
        <div v-else-if="searched" class="text-center py-20">
          <p class="theme-text-muted text-lg">{{ t('search.noPosts') }}</p>
        </div>
      </div>

      <!-- Initial State -->
      <div v-if="!searched" class="text-center py-20">
        <svg class="w-16 h-16 mx-auto theme-text-muted mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        <p class="theme-text-muted">{{ t('search.hint') }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '../stores/app'
import { productAPI, postAPI } from '../api'
import { getFirstImageUrl } from '../utils/image'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const appStore = useAppStore()

const query = ref('')
const activeTab = ref<'all' | 'products' | 'posts'>('all')
const loading = ref(false)
const searched = ref(false)
const products = ref<any[]>([])
const posts = ref<any[]>([])
const searchInput = ref<HTMLInputElement>()

const getLocalizedText = (jsonData: any) => {
  if (!jsonData) return ''
  const locale = appStore.locale
  return jsonData[locale] || jsonData['zh-CN'] || jsonData['en-US'] || ''
}

const formatDate = (dateString: any) => {
  if (!dateString) return '-'
  const s = String(dateString)
  return s.length >= 10 ? s.substring(0, 10) : s
}

const formatMoney = (amount: any) => {
  if (amount === null || amount === undefined) return '-'
  const num = typeof amount === 'object' ? (amount.decimal || amount.amount || 0) : Number(amount)
  if (isNaN(num)) return '-'
  return '¥' + Number(num).toFixed(2)
}

const doSearch = async () => {
  const q = query.value.trim()
  if (!q) {
    searched.value = false
    products.value = []
    posts.value = []
    return
  }
  searched.value = true
  loading.value = true
  try {
    if (activeTab.value === 'all') {
      const [prodRes, postRes] = await Promise.all([
        productAPI.list({ search: q, page: 1, page_size: 8 }),
        postAPI.list({ search: q, page: 1, page_size: 10, type: 'blog' }),
      ])
      products.value = prodRes.data.data || []
      posts.value = postRes.data.data || []
    } else if (activeTab.value === 'products') {
      const res = await productAPI.list({ search: q, page: 1, page_size: 12 })
      products.value = res.data.data || []
      posts.value = []
    } else {
      const res = await postAPI.list({ search: q, page: 1, page_size: 20, type: 'blog' })
      posts.value = res.data.data || []
      products.value = []
    }
  } catch {
    products.value = []
    posts.value = []
  } finally {
    loading.value = false
  }
}

const goToProduct = (slug: string) => {
  router.push(`/products/${slug}`)
}

const goToPost = (slug: string) => {
  router.push(`/blog/${slug}`)
}

onMounted(() => {
  const q = route.query.q as string
  if (q) {
    query.value = q
    doSearch()
  }
  searchInput.value?.focus()
})
</script>

<style scoped>
.line-clamp-2 {
  overflow: hidden;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  line-clamp: 2;
}
</style>
