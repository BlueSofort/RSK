<template>
  <div class="blog-page min-h-screen theme-page pt-20 pb-16">
    <div class="container mx-auto px-4">
      <!-- Page Header -->
      <div class="mb-12 mt-12 text-center">
        <h1 class="text-4xl md:text-5xl font-black mb-4 tracking-tight theme-text-primary">{{ t('nav.blog') }}</h1>
        <p class="theme-text-secondary max-w-2xl mx-auto text-lg border-b theme-border pb-8">
          {{ t('blog.subtitle') }}
        </p>
      </div>

      <div class="flex flex-col lg:flex-row gap-8">
        <!-- Sidebar -->
        <CategorySidebar
          :categories="categoryGroups"
          :selected-category="selectedCategory"
          :expanded-parent-ids="expandedParentIds"
          :show-drawer="showFilterDrawer"
          :show-search="true"
          :search-query="searchQuery"
          :search-label="t('blog.searchLabel')"
          :search-placeholder="t('blog.searchPlaceholder')"
          :all-label="t('blog.allPosts')"
          :title-label="t('blog.categories')"
          @select-category="selectCategory"
          @toggle-parent="toggleParentCategory"
          @update:show-drawer="showFilterDrawer = $event"
          @update:search-query="onSearchChange"
          @clear-search="clearSearch"
        />

        <!-- Main Content - Posts List -->
        <main class="flex-1">
          <!-- Loading Skeleton -->
          <div v-if="loading" class="space-y-6">
            <div v-for="i in 5" :key="i" class="theme-panel rounded-2xl border p-6 animate-pulse">
              <div class="h-4 w-20 rounded theme-skeleton mb-3"></div>
              <div class="h-6 w-3/4 rounded theme-skeleton mb-2"></div>
              <div class="h-4 w-full rounded theme-skeleton mb-1"></div>
              <div class="h-4 w-2/3 rounded theme-skeleton"></div>
            </div>
          </div>

          <!-- Posts List -->
          <div v-else-if="posts.length > 0" class="space-y-6">
            <article
              v-for="post in posts"
              :key="post.id"
              class="group theme-panel backdrop-blur-xl border rounded-2xl overflow-hidden hover:shadow-xl hover:-translate-y-1 transition-all duration-300 cursor-pointer"
              @click="goToPost(post.slug)"
            >
              <div class="flex flex-col md:flex-row">
                <div v-if="post.thumbnail" class="md:w-48 h-48 md:h-auto shrink-0 overflow-hidden">
                  <img :src="getImageUrl(post.thumbnail)" :alt="getLocalizedText(post.title)"
                    loading="lazy" class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105" />
                </div>
                <div class="p-6 flex flex-col flex-1">
                  <div class="flex items-center justify-between mb-3">
                    <span class="theme-badge-meta text-xs"
                      :class="post.type === 'blog' ? 'theme-badge-accent' : 'theme-badge-info'">
                      {{ post.type === 'blog' ? t('nav.blog') : t('nav.notice') }}
                    </span>
                    <time class="text-xs theme-text-muted font-mono">
                      {{ formatDate(post.published_at || post.created_at) }}
                    </time>
                  </div>
                  <h2 class="text-xl font-bold mb-2 theme-text-primary leading-tight group-hover:text-primary transition-colors">
                    {{ getLocalizedText(post.title) }}
                  </h2>
                  <p class="text-sm theme-text-secondary line-clamp-2 leading-relaxed flex-1">
                    {{ getLocalizedText(post.summary) }}
                  </p>
                  <div class="flex items-center text-sm font-medium theme-text-muted group-hover:text-foreground transition-colors mt-4 pt-4 border-t theme-border">
                    {{ t('blog.readMore') }}
                    <svg class="w-4 h-4 ml-2 transition-transform group-hover:translate-x-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
                    </svg>
                  </div>
                </div>
              </div>
            </article>

            <PaginationNav
              :current-page="currentPage"
              :total-pages="totalPages"
              @change-page="changePage"
            />
          </div>

          <!-- Empty State -->
          <div v-else class="text-center py-20 border theme-panel-soft rounded-2xl backdrop-blur-sm">
            <svg class="w-20 h-20 mx-auto theme-text-muted mb-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z" />
            </svg>
            <p class="theme-text-muted text-lg">
              {{ (searchQuery || selectedCategory) ? t('blog.emptyFiltered') : t('blog.empty') }}
            </p>
            <button
              v-if="searchQuery || selectedCategory"
              class="mt-4 theme-btn-inline-md border theme-btn-secondary font-semibold rounded-full"
              @click="clearSearch(); selectCategory(null)"
            >
              {{ t('blog.clearFilters') }}
            </button>
          </div>
        </main>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '../stores/app'
import { postAPI, categoryAPI } from '../api'
import { getImageUrl } from '../utils/image'
import CategorySidebar from '../components/CategorySidebar.vue'
import PaginationNav from '../components/PaginationNav.vue'

const router = useRouter()
const { t } = useI18n()
const appStore = useAppStore()

const loading = ref(true)
const posts = ref<any[]>([])
const currentPage = ref(1)
const totalPages = ref(0)
const pageSize = 12

const categories = ref<any[]>([])
const categoryGroups = ref<any[]>([])
const selectedCategory = ref<any>(null)
const expandedParentIds = ref<number[]>([])
const showFilterDrawer = ref(false)
const searchQuery = ref('')

let searchTimer: ReturnType<typeof setTimeout> | null = null

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

const fetchCategories = async () => {
  try {
    const res = await categoryAPI.list({ type: 'post' })
    const all = res.data.data || []
    categories.value = all
    // Build flat group for sidebar
    categoryGroups.value = all
      .filter((c: any) => !c.parent_id || c.parent_id === 0)
      .map((parent: any) => ({
        ...parent,
        children: all.filter((c: any) => c.parent_id === parent.id),
      }))
  } catch {
    categories.value = []
    categoryGroups.value = []
  }
}

const selectCategory = (cat: any) => {
  selectedCategory.value = cat
  currentPage.value = 1
  loadPosts()
}

const toggleParentCategory = (id: number) => {
  const idx = expandedParentIds.value.indexOf(id)
  if (idx >= 0) {
    expandedParentIds.value.splice(idx, 1)
  } else {
    expandedParentIds.value.push(id)
  }
}

const onSearchChange = (q: string) => {
  searchQuery.value = q
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    currentPage.value = 1
    loadPosts()
  }, 400)
}

const clearSearch = () => {
  searchQuery.value = ''
  loadPosts()
}

const loadPosts = async () => {
  loading.value = true
  try {
    const params: any = {
      type: 'blog',
      page: currentPage.value,
      page_size: pageSize,
    }
    if (selectedCategory.value?.id) {
      params.category_id = selectedCategory.value.id
    }
    if (searchQuery.value.trim()) {
      params.search = searchQuery.value.trim()
    }
    const response = await postAPI.list(params)
    posts.value = response.data.data || []
    if (response.data.pagination) {
      totalPages.value = response.data.pagination.total_page || 0
    }
  } catch (error) {
    console.error('Failed to load posts:', error)
  } finally {
    loading.value = false
  }
}

const goToPost = (slug: string) => {
  router.push(`/blog/${slug}`)
}

const changePage = (page: number) => {
  currentPage.value = page
  loadPosts()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => {
  fetchCategories()
  loadPosts()
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
