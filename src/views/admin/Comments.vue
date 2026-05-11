<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminAPI } from '@/api/admin'
import type { AdminComment } from '@/api/types'
import IdCell from '@/components/IdCell.vue'
import { formatDate } from '@/utils/format'
import { getImageUrl } from '@/utils/image'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import TableSkeleton from '@/components/TableSkeleton.vue'
import { confirmAction } from '@/utils/confirm'
import { notifyError } from '@/utils/notify'

const { t } = useI18n()

const loading = ref(false)
const comments = ref<AdminComment[]>([])
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0,
  total_page: 0,
})
const jumpPage = ref('')
const filterPostId = ref('')

const fetchComments = async () => {
  loading.value = true
  try {
    const params: Record<string, unknown> = {
      page: pagination.page,
      page_size: pagination.page_size,
    }
    const pid = Number(filterPostId.value.trim())
    if (Number.isFinite(pid) && pid > 0) {
      params.post_id = pid
    }
    const res = await adminAPI.getComments(params)
    comments.value = (res.data.data as AdminComment[]) || []
    if (res.data.pagination) {
      Object.assign(pagination, res.data.pagination)
    }
  } catch {
    comments.value = []
  } finally {
    loading.value = false
  }
}

const changePage = (page: number) => {
  pagination.page = page
  fetchComments()
}

const jumpToPage = () => {
  if (!jumpPage.value) return
  const target = Math.min(Math.max(Math.floor(Number(jumpPage.value)), 1), pagination.total_page || 1)
  if (target === pagination.page) return
  changePage(target)
}

const handleDelete = async (comment: AdminComment) => {
  const confirmed = await confirmAction({
    description: t('admin.comments.confirmDelete', { id: comment.id }),
    confirmText: t('admin.common.delete'),
    variant: 'destructive',
  })
  if (!confirmed) return
  try {
    await adminAPI.deleteComment(comment.id)
    fetchComments()
  } catch (err) {
    notifyError(t('admin.comments.deleteFailed', { message: (err as Error).message || '' }))
  }
}

const statusBadge = (status: string) => {
  if (status === 'approved') return 'text-emerald-700 border-emerald-200 bg-emerald-50'
  if (status === 'rejected') return 'text-red-700 border-red-200 bg-red-50'
  return 'text-amber-700 border-amber-200 bg-amber-50'
}

const statusLabel = (status: string) => {
  if (status === 'approved') return t('admin.comments.status.approved')
  if (status === 'rejected') return t('admin.comments.status.rejected')
  return status || '-'
}

const truncate = (text: string, max = 60) => {
  if (!text) return '-'
  return text.length > max ? text.slice(0, max) + '...' : text
}

onMounted(() => {
  fetchComments()
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <h1 class="text-2xl font-semibold">{{ t('admin.comments.title') }}</h1>
    </div>

    <!-- Filter -->
    <div class="flex items-center gap-3">
      <Input
        v-model="filterPostId"
        type="number"
        :placeholder="t('admin.comments.filterPostId')"
        class="h-9 w-full sm:w-48"
      />
      <Button variant="outline" size="sm" class="h-9" @click="fetchComments">
        {{ t('admin.common.filter') }}
      </Button>
    </div>

    <div class="rounded-xl border border-border bg-card overflow-x-auto">
      <Table class="min-w-[900px]">
        <TableHeader class="border-b border-border bg-muted/40 text-xs uppercase text-muted-foreground">
          <TableRow>
            <TableHead class="px-6 py-3">{{ t('admin.comments.table.id') }}</TableHead>
            <TableHead class="px-6 py-3">{{ t('admin.comments.table.post') }}</TableHead>
            <TableHead class="px-6 py-3">{{ t('admin.comments.table.user') }}</TableHead>
            <TableHead class="min-w-[200px] px-6 py-3">{{ t('admin.comments.table.content') }}</TableHead>
            <TableHead class="px-6 py-3">{{ t('admin.comments.table.status') }}</TableHead>
            <TableHead class="px-6 py-3">{{ t('admin.comments.table.createdAt') }}</TableHead>
            <TableHead class="px-6 py-3 text-right">{{ t('admin.comments.table.action') }}</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody class="divide-y divide-border">
          <TableRow v-if="loading">
            <TableCell :colspan="7" class="p-0">
              <TableSkeleton :columns="7" :rows="5" />
            </TableCell>
          </TableRow>
          <TableRow v-else-if="comments.length === 0">
            <TableCell colspan="7" class="px-6 py-8 text-center text-muted-foreground">
              {{ t('admin.comments.empty') }}
            </TableCell>
          </TableRow>
          <TableRow v-for="comment in comments" :key="comment.id" class="hover:bg-muted/30 group">
            <TableCell class="px-6 py-4">
              <IdCell :value="comment.id" />
            </TableCell>
            <TableCell class="px-6 py-4">
              <span class="font-mono text-xs text-muted-foreground">#{{ comment.post_id }}</span>
              <span v-if="comment.post_title" class="ml-1 text-sm">{{ comment.post_title }}</span>
            </TableCell>
            <TableCell class="px-6 py-4">
              <div class="flex items-center gap-2">
                <img
                  v-if="comment.user_avatar"
                  :src="getImageUrl(comment.user_avatar)"
                  class="w-6 h-6 rounded-full object-cover bg-muted"
                />
                <div class="text-sm">
                  <div class="font-medium">{{ comment.user_name }}</div>
                  <div class="text-xs text-muted-foreground">{{ comment.user_email }}</div>
                </div>
              </div>
            </TableCell>
            <TableCell class="min-w-[200px] px-6 py-4">
              <div class="text-sm break-words">{{ truncate(comment.content, 100) }}</div>
              <div v-if="comment.parent_id > 0" class="text-xs text-muted-foreground mt-1">
                {{ t('admin.comments.replyTo') }} #{{ comment.parent_id }}
              </div>
            </TableCell>
            <TableCell class="px-6 py-4">
              <span class="inline-flex rounded-full border px-2.5 py-1 text-xs" :class="statusBadge(comment.status)">
                {{ statusLabel(comment.status) }}
              </span>
            </TableCell>
            <TableCell class="px-6 py-4 text-xs text-muted-foreground">{{ formatDate(comment.created_at) }}</TableCell>
            <TableCell class="px-6 py-4 text-right">
              <div class="flex items-center justify-end gap-2 md:opacity-0 md:group-hover:opacity-100 transition-opacity">
                <Button size="icon-sm" variant="destructive" @click="handleDelete(comment)">
                  <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </Button>
              </div>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>

      <div
        v-if="pagination.total_page > 1"
        class="flex flex-col gap-3 border-t border-border px-6 py-4 sm:flex-row sm:items-center sm:justify-between"
      >
        <div class="flex items-center gap-3">
          <span class="text-xs text-muted-foreground">
            {{ t('admin.common.pageInfo', { total: pagination.total, page: pagination.page, totalPage: pagination.total_page }) }}
          </span>
        </div>
        <div class="flex items-center gap-3">
          <Input v-model="jumpPage" type="number" min="1" :max="pagination.total_page" class="h-8 w-20" />
          <Button variant="outline" size="sm" class="h-8" @click="jumpToPage">{{ t('admin.common.jumpTo') }}</Button>
          <Button variant="outline" size="sm" class="h-8" :disabled="pagination.page <= 1" @click="changePage(pagination.page - 1)">{{ t('admin.common.prevPage') }}</Button>
          <Button variant="outline" size="sm" class="h-8" :disabled="pagination.page >= pagination.total_page" @click="changePage(pagination.page + 1)">{{ t('admin.common.nextPage') }}</Button>
        </div>
      </div>
    </div>
  </div>
</template>
