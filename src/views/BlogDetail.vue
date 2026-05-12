<template>
  <div
    class="blog-detail-page min-h-screen theme-page pt-24 pb-16 relative overflow-hidden">
    <div class="container mx-auto px-4 max-w-4xl relative z-10">
      <!-- Loading State -->
      <div v-if="loading" class="animate-pulse space-y-8">
        <div class="h-8 theme-surface-muted rounded w-1/3"></div>
        <div class="space-y-4">
          <div class="h-12 theme-surface-muted rounded w-3/4"></div>
          <div class="h-6 theme-surface-muted rounded w-1/2"></div>
        </div>
        <div class="h-96 theme-surface-muted rounded-3xl"></div>
      </div>

      <!-- Post Content -->
      <article v-else-if="post">
        <!-- Breadcrumb -->
        <nav class="mb-8 flex items-center space-x-2 text-sm theme-text-muted font-medium">
          <router-link to="/" class="theme-link-muted">{{ t('nav.home')
          }}</router-link>
          <span>/</span>
          <router-link :to="backLink" class="theme-link-muted">{{ backText
          }}</router-link>
          <span>/</span>
          <span class="theme-text-primary truncate max-w-[200px]">{{ getLocalizedText(post.title) }}</span>
        </nav>

        <div
          class="theme-panel backdrop-blur-xl border rounded-3xl p-8 md:p-12 shadow-2xl relative overflow-hidden">
          <!-- Featured Image -->
          <div v-if="post.thumbnail" class="mb-12 relative h-64 md:h-96 rounded-2xl overflow-hidden group">
            <img :src="getImageUrl(post.thumbnail)" :alt="getLocalizedText(post.title)"
              loading="lazy" class="w-full h-full object-cover">
            <div class="absolute inset-0 bg-black/20 dark:bg-black/35"></div>
          </div>

          <!-- Post Header -->
          <header class="mb-12 border-b theme-border pb-12">
            <div class="flex flex-wrap items-center gap-4 mb-6">
              <span class="theme-badge-meta" :class="post.type === 'blog'
                ? 'theme-badge-accent'
                : 'theme-badge-info'">
                {{ post.type === 'blog' ? t('nav.blog') : t('nav.notice') }}
              </span>
              <time class="text-sm theme-text-muted font-mono">
                {{ formatDate(post.published_at || post.created_at) }}
              </time>
            </div>

            <h1 class="text-3xl md:text-5xl font-black theme-text-primary mb-6 leading-tight tracking-tight">
              {{ getLocalizedText(post.title) }}
            </h1>

            <p v-if="post.summary" class="text-xl theme-text-secondary leading-relaxed font-light">
              {{ getLocalizedText(post.summary) }}
            </p>
          </header>

          <!-- Post Content -->
          <div v-html="processHtmlForDisplay(getLocalizedText(post.content))"
            class="prose prose-lg max-w-none dark:prose-invert theme-prose">
          </div>

          <!-- Footer -->
          <footer class="mt-16 pt-12 border-t theme-border flex justify-center">
            <router-link :to="backLink"
              class="group inline-flex items-center space-x-3 theme-link-muted px-6 py-3 border theme-btn-secondary rounded-full">
              <svg class="w-5 h-5 transition-transform group-hover:-translate-x-1" fill="none" stroke="currentColor"
                viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
              </svg>
              <span class="font-medium">{{ backText }}</span>
            </router-link>
          </footer>

          <!-- Comments Section -->
          <section class="mt-16 pt-12 border-t theme-border">
            <h2 class="text-2xl font-bold theme-text-primary mb-8">
              {{ t('blogDetail.comments') }} <span v-if="commentsTotal > 0" class="text-lg theme-text-muted font-normal">({{ commentsTotal }})</span>
            </h2>

            <!-- Comment Form -->
            <div v-if="isLoggedIn" class="mb-10">
              <div v-if="replyTo" class="mb-3 flex items-center gap-2 text-sm theme-text-muted">
                <span>{{ t('blogDetail.replyingTo') }} <strong>{{ replyTo.user_name }}</strong></span>
                <button @click="cancelReply" class="theme-link-muted text-xs">&times; {{ t('blogDetail.cancelReply') }}</button>
              </div>
              <div class="flex gap-3">
                <img
                  :src="getImageUrl(currentUserAvatar)"
                  class="w-10 h-10 rounded-full object-cover flex-shrink-0 bg-muted"
                  @error="($event.target as HTMLImageElement).src = defaultAvatar"
                />
                <div class="flex-1 space-y-3">
                  <textarea
                    v-model="commentForm.content"
                    :placeholder="t('blogDetail.commentPlaceholder')"
                    rows="3"
                    maxlength="100"
                    class="w-full rounded-xl border border-border bg-background px-4 py-3 text-sm resize-none focus:outline-none focus:ring-2 focus:ring-primary/20"
                  ></textarea>
                  <div class="flex items-center justify-between">
                    <span class="text-xs theme-text-muted">{{ commentForm.content.length }}/100</span>
                    <button
                      @click="submitComment"
                      :disabled="!commentForm.content.trim() || submittingComment"
                      class="px-5 py-2 rounded-full theme-btn-primary text-sm font-semibold disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      {{ submittingComment ? t('blogDetail.submitting') : t('blogDetail.submitComment') }}
                    </button>
                  </div>
                </div>
              </div>
            </div>
            <div v-else class="mb-10 p-6 rounded-2xl theme-surface-muted text-center">
              <p class="theme-text-muted text-sm">
                <router-link to="/login" class="theme-link font-semibold">{{ t('blogDetail.loginToComment') }}</router-link>
              </p>
            </div>

            <!-- Comment List -->
            <div v-if="commentsLoading" class="text-center py-8">
              <div class="inline-block w-6 h-6 border-2 border-primary border-t-transparent rounded-full animate-spin"></div>
            </div>
            <div v-else-if="comments.length === 0" class="text-center py-8 theme-text-muted text-sm">
              {{ t('blogDetail.noComments') }}
            </div>
            <div v-else class="space-y-6">
              <div v-for="comment in comments" :key="comment.id" class="group">
                <!-- Parent Comment -->
                <div class="flex gap-3">
                  <img
                    :src="getImageUrl(comment.user_avatar)"
                    class="w-10 h-10 rounded-full object-cover flex-shrink-0 bg-muted"
                    @error="($event.target as HTMLImageElement).src = defaultAvatar"
                  />
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center gap-2 mb-1">
                      <span class="font-semibold text-sm theme-text-primary">{{ comment.user_name }}</span>
                      <time class="text-xs theme-text-muted">{{ formatDate(comment.created_at) }}</time>
                    </div>
                    <p class="text-sm theme-text-secondary leading-relaxed break-words">{{ comment.content }}</p>
                    <div class="flex items-center gap-3 mt-2">
                      <button
                        v-if="isLoggedIn"
                        @click="replyToComment(comment)"
                        class="text-xs theme-link-muted hover:underline"
                      >{{ t('blogDetail.reply') }}</button>
                      <button
                        v-if="isLoggedIn && currentUserId === comment.user_id"
                        @click="deleteComment(comment.id)"
                        class="text-xs text-red-400 hover:text-red-600 hover:underline"
                      >{{ t('blogDetail.delete') }}</button>
                    </div>

                    <!-- Replies -->
                    <div v-if="repliesMap[comment.id]?.length" class="mt-4 ml-2 pl-4 border-l-2 theme-border space-y-4">
                      <div v-for="reply in repliesMap[comment.id]" :key="reply.id" class="flex gap-3">
                        <img
                          :src="getImageUrl(reply.user_avatar)"
                          class="w-8 h-8 rounded-full object-cover flex-shrink-0 bg-muted"
                          @error="($event.target as HTMLImageElement).src = defaultAvatar"
                        />
                        <div class="flex-1 min-w-0">
                          <div class="flex items-center gap-2 mb-1">
                            <span class="font-semibold text-xs theme-text-primary">{{ reply.user_name }}</span>
                            <time class="text-xs theme-text-muted">{{ formatDate(reply.created_at) }}</time>
                          </div>
                          <p class="text-sm theme-text-secondary leading-relaxed break-words">{{ reply.content }}</p>
                          <button
                            v-if="isLoggedIn && currentUserId === reply.user_id"
                            @click="deleteComment(reply.id)"
                            class="text-xs text-red-400 hover:text-red-600 hover:underline mt-1"
                          >{{ t('blogDetail.delete') }}</button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Load More -->
              <div v-if="commentsHasMore" class="text-center pt-4">
                <button
                  @click="loadMoreComments"
                  :disabled="commentsLoadingMore"
                  class="px-6 py-2 rounded-full border theme-btn-secondary text-sm"
                >
                  {{ commentsLoadingMore ? t('blogDetail.loading') : t('blogDetail.loadMore') }}
                </button>
              </div>
            </div>
          </section>
        </div>
      </article>

      <!-- Error State -->
      <div v-else
        class="text-center py-24 theme-panel rounded-3xl border backdrop-blur-sm">
        <svg class="w-20 h-20 mx-auto theme-text-muted mb-6" fill="none" stroke="currentColor"
          viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
            d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <p class="theme-text-muted text-xl mb-8">
          {{ t('blogDetail.notFound') }}
        </p>
        <router-link to="/blog"
          class="inline-block theme-btn-primary px-8 py-3 rounded-full font-bold hover:scale-105 transition-transform">
          {{ t('blogDetail.backToBlog') }}
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, reactive } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '../stores/app'
import { useUserAuthStore } from '../stores/userAuth'
import { postAPI, commentAPI } from '../api'
import { getImageUrl } from '../utils/image'
import { processHtmlForDisplay } from '../utils/content'
import { toast } from '../composables/useToast'

const route = useRoute()
const { t } = useI18n()
const appStore = useAppStore()
const authStore = useUserAuthStore()

const loading = ref(true)
const post = ref<any>(null)

// --- Comments ---
const comments = ref<any[]>([])
const repliesMap = ref<Record<number, any[]>>({})
const commentsLoading = ref(false)
const commentsLoadingMore = ref(false)
const commentsPage = ref(1)
const commentsTotal = ref(0)
const commentsHasMore = ref(false)
const submittingComment = ref(false)
const replyTo = ref<any>(null)
const commentForm = reactive({ content: '' })

const defaultAvatar = 'data:image/svg+xml,' + encodeURIComponent('<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 40 40"><rect width="40" height="40" fill="#e5e7eb"/><text x="20" y="26" text-anchor="middle" font-size="18" fill="#9ca3af">?</text></svg>')

const isLoggedIn = computed(() => authStore.isAuthenticated)
const currentUserId = computed(() => authStore.user?.id)
const currentUserAvatar = computed(() => authStore.user?.avatar || '')

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

// Flatten nested user object to flat fields expected by template
const flattenUser = (item: any) => ({
  ...item,
  user_name: item.user?.display_name || '',
  user_avatar: item.user?.avatar || '',
})

const backLink = computed(() => {
  if (!post.value) return '/blog'
  return post.value.type === 'notice' ? '/notice' : '/blog'
})

const backText = computed(() => {
  if (!post.value) return t('blogDetail.backToBlog')
  return post.value.type === 'notice' ? t('blogDetail.backToNotice') : t('blogDetail.backToBlog')
})

const loadPost = async () => {
  loading.value = true
  try {
    const slug = route.params.slug as string
    const response = await postAPI.detail(slug)
    post.value = response.data.data || null
    if (post.value) {
      loadComments()
    }
  } catch (error) {
    console.error('Failed to load post:', error)
    post.value = null
  } finally {
    loading.value = false
  }
}

// --- Comment functions ---
const loadComments = async (page = 1) => {
  if (!post.value) return
  const isFirstPage = page === 1
  if (isFirstPage) {
    commentsLoading.value = true
    comments.value = []
    repliesMap.value = {}
  } else {
    commentsLoadingMore.value = true
  }
  try {
    const res = await commentAPI.list(post.value.id, { page, page_size: 10 })
    const data = res.data
    const list = (data.data?.list || []).map(flattenUser)
    const rawReplies = (data.data?.replies || []).map(flattenUser)
    const pagination = data.pagination || {}

    // Group flat replies array by parent_id
    const replies: Record<number, any[]> = {}
    if (Array.isArray(rawReplies)) {
      for (const r of rawReplies) {
        const pid = r.parent_id
        if (!replies[pid]) replies[pid] = []
        replies[pid].push(r)
      }
    }

    if (isFirstPage) {
      comments.value = list
      repliesMap.value = replies
    } else {
      comments.value = [...comments.value, ...list]
      Object.assign(repliesMap.value, replies)
    }
    commentsTotal.value = pagination.total || 0
    commentsPage.value = page
    commentsHasMore.value = page < (pagination.total_page || 1)
  } catch (err) {
    console.error('Failed to load comments:', err)
  } finally {
    commentsLoading.value = false
    commentsLoadingMore.value = false
  }
}

const loadMoreComments = () => {
  loadComments(commentsPage.value + 1)
}

const submitComment = async () => {
  if (!commentForm.content.trim() || !post.value) return
  submittingComment.value = true
  try {
    await commentAPI.create({
      post_id: post.value.id,
      parent_id: replyTo.value?.id || 0,
      content: commentForm.content.trim(),
    })
    commentForm.content = ''
    replyTo.value = null
    loadComments() // reload
  } catch (err: any) {
    const sensitiveWord = err?.responseData?.sensitive_word || ''
    if (sensitiveWord) {
      toast.error(t('blogDetail.commentSensitive', { word: sensitiveWord }))
    } else {
      toast.error(err?.message || t('blogDetail.commentFailed'))
    }
  } finally {
    submittingComment.value = false
  }
}

const replyToComment = (comment: any) => {
  replyTo.value = comment
  commentForm.content = ''
  // scroll to form
  const formEl = document.querySelector('textarea')
  formEl?.focus()
}

const cancelReply = () => {
  replyTo.value = null
  commentForm.content = ''
}

const deleteComment = async (commentId: number) => {
  if (!confirm(t('blogDetail.confirmDelete'))) return
  try {
    await commentAPI.delete(commentId)
    loadComments()
  } catch (err) {
    console.error('Failed to delete comment:', err)
  }
}

onMounted(() => {
  loadPost()
})
</script>
