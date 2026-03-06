<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminAPI } from '@/api/admin'
import IdCell from '@/components/IdCell.vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Dialog, DialogScrollContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Textarea } from '@/components/ui/textarea'
import { confirmAction } from '@/utils/confirm'
import { notifyError, notifySuccess } from '@/utils/notify'

const { t } = useI18n()
const loading = ref(true)
const jobs = ref<any[]>([])
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0,
  total_page: 1,
})
const jumpPage = ref('')

const filters = reactive({
  status: '__all__',
  type: '__all__',
  connection_id: '__all__',
})

// 站点连接列表
const connections = ref<any[]>([])
const loadingConnections = ref(false)

const fetchConnections = async () => {
  loadingConnections.value = true
  try {
    const res = await adminAPI.getSiteConnections({ page: 1, page_size: 200 })
    connections.value = (res.data.data as any[]) || []
  } catch {
    connections.value = []
  } finally {
    loadingConnections.value = false
  }
}

// Job list
const statusOptions = [
  { value: '__all__', key: 'reconciliation.filters.allStatus' },
  { value: 'pending', key: 'reconciliation.status.pending' },
  { value: 'running', key: 'reconciliation.status.running' },
  { value: 'completed', key: 'reconciliation.status.completed' },
  { value: 'failed', key: 'reconciliation.status.failed' },
]

const typeOptions = [
  { value: '__all__', key: 'reconciliation.filters.allTypes' },
  { value: 'status', key: 'reconciliation.type.status' },
  { value: 'amount', key: 'reconciliation.type.amount' },
  { value: 'full', key: 'reconciliation.type.full' },
]

// New job form
const showNewJob = ref(false)
const newJobForm = reactive({
  connection_id: '',
  type: 'full',
  time_range_start: '',
  time_range_end: '',
})

const openNewJob = () => {
  if (connections.value.length === 0) fetchConnections()
  showNewJob.value = true
}
const submitting = ref(false)

// Detail dialog
const showDetail = ref(false)
const detailJob = ref<any>(null)
const detailItems = ref<any[]>([])
const detailItemsTotal = ref(0)
const detailItemsPage = ref(1)

// Resolve dialog
const showResolve = ref(false)
const resolveItemId = ref<number | null>(null)
const resolveRemark = ref('')
const resolving = ref(false)

const fetchJobs = async (page = 1) => {
  loading.value = true
  try {
    const params: any = { page, page_size: pagination.page_size }
    if (filters.status && filters.status !== '__all__') params.status = filters.status
    if (filters.type && filters.type !== '__all__') params.type = filters.type
    if (filters.connection_id && filters.connection_id !== '__all__') params.connection_id = filters.connection_id

    const res = await adminAPI.getReconciliationJobs(params)
    jobs.value = (res.data.data as any[]) || []
    const p = (res.data as any).pagination
    if (p) {
      pagination.page = p.page
      pagination.page_size = p.page_size
      pagination.total = p.total
      pagination.total_page = p.total_page
    }
  } catch {
    jobs.value = []
  } finally {
    loading.value = false
  }
}

const changePage = (page: number) => {
  if (page < 1 || page > pagination.total_page) return
  fetchJobs(page)
}

const jumpToPage = () => {
  if (!jumpPage.value) return
  const raw = Number(jumpPage.value)
  if (Number.isNaN(raw)) return
  const target = Math.min(Math.max(Math.floor(raw), 1), pagination.total_page)
  if (target === pagination.page) return
  changePage(target)
}

const handleSearch = () => {
  fetchJobs(1)
}

const connectionName = (id: number | string) => {
  const conn = connections.value.find((c: any) => String(c.id) === String(id))
  return conn ? conn.name : String(id)
}

const handleNewJob = async () => {
  if (!newJobForm.connection_id || !newJobForm.time_range_start || !newJobForm.time_range_end) return
  submitting.value = true
  try {
    await adminAPI.runReconciliation({
      connection_id: Number(newJobForm.connection_id),
      type: newJobForm.type,
      time_range_start: new Date(newJobForm.time_range_start).toISOString(),
      time_range_end: new Date(newJobForm.time_range_end).toISOString(),
    })
    notifySuccess(t('reconciliation.form.submitSuccess'))
    showNewJob.value = false
    newJobForm.connection_id = ''
    newJobForm.time_range_start = ''
    newJobForm.time_range_end = ''
    fetchJobs(1)
  } catch (err: any) {
    notifyError(err?.response?.data?.message || err?.message)
  } finally {
    submitting.value = false
  }
}

const openDetail = async (job: any) => {
  try {
    const res = await adminAPI.getReconciliationJob(job.id, { items_page: 1, items_page_size: 20 })
    const data = res.data.data
    detailJob.value = data.job || job
    detailItems.value = data.items || []
    detailItemsTotal.value = data.items_total || 0
    detailItemsPage.value = 1
  } catch {
    detailJob.value = job
    detailItems.value = []
    detailItemsTotal.value = 0
  }
  showDetail.value = true
}

const loadDetailItems = async (page: number) => {
  if (!detailJob.value) return
  try {
    const res = await adminAPI.getReconciliationJob(detailJob.value.id, { items_page: page, items_page_size: 20 })
    const data = res.data.data
    detailItems.value = data.items || []
    detailItemsTotal.value = data.items_total || 0
    detailItemsPage.value = page
  } catch {
    // keep current items
  }
}

const closeDetail = () => {
  showDetail.value = false
  detailJob.value = null
  detailItems.value = []
}

const openResolve = (item: any) => {
  resolveItemId.value = item.id
  resolveRemark.value = ''
  showResolve.value = true
}

const handleResolve = async () => {
  if (!resolveItemId.value) return
  resolving.value = true
  try {
    await adminAPI.resolveReconciliationItem(resolveItemId.value, { remark: resolveRemark.value })
    notifySuccess(t('reconciliation.items.resolveSuccess'))
    showResolve.value = false
    // Refresh detail items
    if (detailJob.value) {
      await loadDetailItems(detailItemsPage.value)
    }
  } catch (err: any) {
    notifyError(err?.response?.data?.message || err?.message)
  } finally {
    resolving.value = false
  }
}

const statusBadgeClass = (status: string) => {
  switch (status) {
    case 'pending':
      return 'text-yellow-700 border-yellow-200 bg-yellow-50'
    case 'running':
      return 'text-blue-700 border-blue-200 bg-blue-50'
    case 'completed':
      return 'text-emerald-700 border-emerald-200 bg-emerald-50'
    case 'failed':
      return 'text-red-700 border-red-200 bg-red-50'
    default:
      return 'text-muted-foreground border-border bg-muted/30'
  }
}

const mismatchBadgeClass = (type: string) => {
  switch (type) {
    case 'status':
      return 'text-orange-700 border-orange-200 bg-orange-50'
    case 'amount':
      return 'text-purple-700 border-purple-200 bg-purple-50'
    case 'both':
      return 'text-red-700 border-red-200 bg-red-50'
    default:
      return 'text-muted-foreground border-border bg-muted/30'
  }
}

const formatTime = (raw?: string) => {
  if (!raw) return '-'
  const d = new Date(raw)
  if (Number.isNaN(d.getTime())) return '-'
  return d.toLocaleString()
}

onMounted(() => {
  fetchConnections()
  fetchJobs()
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-semibold">{{ t('reconciliation.title') }}</h1>
      <Button size="sm" @click="openNewJob">{{ t('reconciliation.newJob') }}</Button>
    </div>

    <!-- Filters -->
    <div class="flex flex-wrap items-end gap-3">
      <div>
        <label class="mb-1.5 block text-xs font-medium text-muted-foreground">{{ t('reconciliation.columns.status') }}</label>
        <Select v-model="filters.status">
          <SelectTrigger class="h-9 w-40">
            <SelectValue :placeholder="t('reconciliation.filters.allStatus')" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem v-for="opt in statusOptions" :key="opt.value" :value="opt.value">
              {{ t(opt.key) }}
            </SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div>
        <label class="mb-1.5 block text-xs font-medium text-muted-foreground">{{ t('reconciliation.columns.type') }}</label>
        <Select v-model="filters.type">
          <SelectTrigger class="h-9 w-40">
            <SelectValue :placeholder="t('reconciliation.filters.allTypes')" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem v-for="opt in typeOptions" :key="opt.value" :value="opt.value">
              {{ t(opt.key) }}
            </SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div>
        <label class="mb-1.5 block text-xs font-medium text-muted-foreground">{{ t('reconciliation.filters.connectionId') }}</label>
        <Select v-model="filters.connection_id">
          <SelectTrigger class="h-9 w-48">
            <SelectValue :placeholder="t('reconciliation.filters.allConnections')" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="__all__">{{ t('reconciliation.filters.allConnections') }}</SelectItem>
            <SelectItem v-for="conn in connections" :key="conn.id" :value="String(conn.id)">
              {{ conn.name || `#${conn.id}` }}
            </SelectItem>
          </SelectContent>
        </Select>
      </div>
      <Button size="sm" class="h-9" @click="handleSearch">{{ t('admin.common.refresh') }}</Button>
    </div>

    <!-- Table -->
    <div class="rounded-xl border border-border bg-card">
      <Table>
        <TableHeader class="border-b border-border bg-muted/40 text-xs uppercase text-muted-foreground">
          <TableRow>
            <TableHead class="px-6 py-3">{{ t('reconciliation.columns.id') }}</TableHead>
            <TableHead class="px-6 py-3">{{ t('reconciliation.columns.connection') }}</TableHead>
            <TableHead class="px-6 py-3">{{ t('reconciliation.columns.type') }}</TableHead>
            <TableHead class="px-6 py-3">{{ t('reconciliation.columns.status') }}</TableHead>
            <TableHead class="px-6 py-3">{{ t('reconciliation.columns.timeRange') }}</TableHead>
            <TableHead class="px-6 py-3">{{ t('reconciliation.columns.total') }}</TableHead>
            <TableHead class="px-6 py-3">{{ t('reconciliation.columns.matched') }}</TableHead>
            <TableHead class="px-6 py-3">{{ t('reconciliation.columns.mismatched') }}</TableHead>
            <TableHead class="px-6 py-3">{{ t('reconciliation.columns.createdAt') }}</TableHead>
            <TableHead class="px-6 py-3 text-right">{{ t('reconciliation.columns.actions') }}</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody class="divide-y divide-border">
          <TableRow v-if="loading">
            <TableCell colspan="10" class="px-6 py-8 text-center text-muted-foreground">{{ t('admin.common.loading') }}</TableCell>
          </TableRow>
          <TableRow v-else-if="jobs.length === 0">
            <TableCell colspan="10" class="px-6 py-8 text-center text-muted-foreground">{{ t('reconciliation.empty') }}</TableCell>
          </TableRow>
          <TableRow
            v-for="job in jobs"
            :key="job.id"
            class="cursor-pointer hover:bg-muted/30"
            @click="openDetail(job)"
          >
            <TableCell class="px-6 py-4">
              <IdCell :value="job.id" />
            </TableCell>
            <TableCell class="px-6 py-4 font-medium text-foreground">
              {{ job.connection?.name || job.connection_id || '-' }}
            </TableCell>
            <TableCell class="px-6 py-4">
              <span class="text-xs">{{ t('reconciliation.type.' + job.type) }}</span>
            </TableCell>
            <TableCell class="px-6 py-4">
              <span
                class="inline-flex rounded-full border px-2.5 py-1 text-xs"
                :class="statusBadgeClass(job.status)"
              >
                {{ t('reconciliation.status.' + job.status) }}
              </span>
            </TableCell>
            <TableCell class="px-6 py-4 text-xs text-muted-foreground">
              {{ formatTime(job.time_range_start) }} ~ {{ formatTime(job.time_range_end) }}
            </TableCell>
            <TableCell class="px-6 py-4 text-xs text-muted-foreground">{{ job.total_count }}</TableCell>
            <TableCell class="px-6 py-4 text-xs text-emerald-600">{{ job.matched_count }}</TableCell>
            <TableCell class="px-6 py-4 text-xs" :class="job.mismatched_count > 0 ? 'text-red-600 font-semibold' : 'text-muted-foreground'">
              {{ job.mismatched_count }}
            </TableCell>
            <TableCell class="px-6 py-4 text-xs text-muted-foreground">
              {{ formatTime(job.created_at) }}
            </TableCell>
            <TableCell class="px-6 py-4 text-right" @click.stop>
              <Button size="sm" variant="outline" @click="openDetail(job)">
                {{ t('reconciliation.detail.title') }}
              </Button>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>

      <!-- Pagination -->
      <div v-if="pagination.total_page > 1" class="flex flex-wrap items-center justify-between gap-3 border-t border-border px-6 py-4">
        <span class="text-xs text-muted-foreground">
          {{ t('admin.common.pageInfo', { total: pagination.total, page: pagination.page, totalPage: pagination.total_page }) }}
        </span>
        <div class="flex flex-wrap items-center gap-2">
          <Input v-model="jumpPage" type="number" min="1" :max="pagination.total_page" class="h-8 w-20" :placeholder="t('admin.common.jumpPlaceholder')" />
          <Button variant="outline" size="sm" class="h-8" @click="jumpToPage">{{ t('admin.common.jumpTo') }}</Button>
          <Button variant="outline" size="sm" class="h-8" :disabled="pagination.page <= 1" @click="changePage(pagination.page - 1)">
            {{ t('admin.common.prevPage') }}
          </Button>
          <Button variant="outline" size="sm" class="h-8" :disabled="pagination.page >= pagination.total_page" @click="changePage(pagination.page + 1)">
            {{ t('admin.common.nextPage') }}
          </Button>
        </div>
      </div>
    </div>

    <!-- New Job Dialog -->
    <Dialog v-model:open="showNewJob">
      <DialogScrollContent class="w-full max-w-md" @interact-outside="(e: Event) => e.preventDefault()">
        <DialogHeader>
          <DialogTitle>{{ t('reconciliation.newJob') }}</DialogTitle>
        </DialogHeader>
        <div class="space-y-4 pt-2">
          <div>
            <label class="mb-1.5 block text-xs font-medium text-muted-foreground">{{ t('reconciliation.form.connectionId') }}</label>
            <Select v-model="newJobForm.connection_id">
              <SelectTrigger class="h-9">
                <SelectValue :placeholder="t('reconciliation.form.selectConnection')" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="conn in connections" :key="conn.id" :value="String(conn.id)">
                  {{ conn.name || `#${conn.id}` }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div>
            <label class="mb-1.5 block text-xs font-medium text-muted-foreground">{{ t('reconciliation.form.type') }}</label>
            <Select v-model="newJobForm.type">
              <SelectTrigger class="h-9">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="status">{{ t('reconciliation.type.status') }}</SelectItem>
                <SelectItem value="amount">{{ t('reconciliation.type.amount') }}</SelectItem>
                <SelectItem value="full">{{ t('reconciliation.type.full') }}</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div>
            <label class="mb-1.5 block text-xs font-medium text-muted-foreground">{{ t('reconciliation.form.timeRangeStart') }}</label>
            <Input v-model="newJobForm.time_range_start" type="datetime-local" class="h-9" />
          </div>
          <div>
            <label class="mb-1.5 block text-xs font-medium text-muted-foreground">{{ t('reconciliation.form.timeRangeEnd') }}</label>
            <Input v-model="newJobForm.time_range_end" type="datetime-local" class="h-9" />
          </div>
          <div class="flex justify-end gap-3 pt-2">
            <Button variant="outline" @click="showNewJob = false">{{ t('admin.common.cancel') }}</Button>
            <Button :disabled="submitting || !newJobForm.connection_id || !newJobForm.time_range_start || !newJobForm.time_range_end" @click="handleNewJob">
              {{ t('reconciliation.form.submit') }}
            </Button>
          </div>
        </div>
      </DialogScrollContent>
    </Dialog>

    <!-- Detail Dialog -->
    <Dialog v-model:open="showDetail" @update:open="(value: boolean) => { if (!value) closeDetail() }">
      <DialogScrollContent class="w-full max-w-4xl" @interact-outside="(e: Event) => e.preventDefault()">
        <DialogHeader>
          <DialogTitle>{{ t('reconciliation.detail.title') }}</DialogTitle>
        </DialogHeader>

        <div v-if="detailJob" class="space-y-4">
          <div class="grid grid-cols-1 gap-4 md:grid-cols-3">
            <div>
              <label class="mb-1 block text-xs font-medium text-muted-foreground">{{ t('reconciliation.columns.id') }}</label>
              <div class="text-sm">{{ detailJob.id }}</div>
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-muted-foreground">{{ t('reconciliation.columns.connection') }}</label>
              <div class="text-sm">{{ detailJob.connection?.name || detailJob.connection_id || '-' }}</div>
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-muted-foreground">{{ t('reconciliation.columns.type') }}</label>
              <div class="text-sm">{{ t('reconciliation.type.' + detailJob.type) }}</div>
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-muted-foreground">{{ t('reconciliation.columns.status') }}</label>
              <div>
                <span class="inline-flex rounded-full border px-2.5 py-1 text-xs" :class="statusBadgeClass(detailJob.status)">
                  {{ t('reconciliation.status.' + detailJob.status) }}
                </span>
              </div>
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-muted-foreground">{{ t('reconciliation.columns.total') }}</label>
              <div class="text-sm">{{ detailJob.total_count }}</div>
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-muted-foreground">{{ t('reconciliation.columns.matched') }} / {{ t('reconciliation.columns.mismatched') }}</label>
              <div class="text-sm">
                <span class="text-emerald-600">{{ detailJob.matched_count }}</span>
                /
                <span :class="detailJob.mismatched_count > 0 ? 'text-red-600 font-semibold' : ''">{{ detailJob.mismatched_count }}</span>
              </div>
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-muted-foreground">{{ t('reconciliation.columns.timeRange') }}</label>
              <div class="text-xs">{{ formatTime(detailJob.time_range_start) }} ~ {{ formatTime(detailJob.time_range_end) }}</div>
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-muted-foreground">{{ t('reconciliation.columns.startedAt') }}</label>
              <div class="text-sm">{{ formatTime(detailJob.started_at) }}</div>
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-muted-foreground">{{ t('reconciliation.columns.finishedAt') }}</label>
              <div class="text-sm">{{ formatTime(detailJob.finished_at) }}</div>
            </div>
          </div>

          <div v-if="detailJob.result_json" class="mt-4">
            <label class="mb-1 block text-xs font-medium text-muted-foreground">{{ t('reconciliation.detail.resultJson') }}</label>
            <div class="rounded-md border border-border bg-muted/30 p-3 text-sm font-mono whitespace-pre-wrap break-all">
              {{ detailJob.result_json }}
            </div>
          </div>

          <!-- Mismatch Items Table -->
          <div v-if="detailItems.length > 0" class="mt-6">
            <h3 class="mb-3 text-sm font-semibold">{{ t('reconciliation.items.title') }} ({{ detailItemsTotal }})</h3>
            <div class="rounded-lg border border-border">
              <Table>
                <TableHeader class="border-b border-border bg-muted/40 text-xs uppercase text-muted-foreground">
                  <TableRow>
                    <TableHead class="px-4 py-2">{{ t('reconciliation.items.localOrderNo') }}</TableHead>
                    <TableHead class="px-4 py-2">{{ t('reconciliation.items.upstreamOrderNo') }}</TableHead>
                    <TableHead class="px-4 py-2">{{ t('reconciliation.items.localStatus') }}</TableHead>
                    <TableHead class="px-4 py-2">{{ t('reconciliation.items.upstreamStatus') }}</TableHead>
                    <TableHead class="px-4 py-2">{{ t('reconciliation.items.mismatchType') }}</TableHead>
                    <TableHead class="px-4 py-2">{{ t('reconciliation.columns.status') }}</TableHead>
                    <TableHead class="px-4 py-2 text-right">{{ t('reconciliation.columns.actions') }}</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody class="divide-y divide-border">
                  <TableRow v-for="item in detailItems" :key="item.id">
                    <TableCell class="px-4 py-3 text-xs font-mono">{{ item.local_order_no || '-' }}</TableCell>
                    <TableCell class="px-4 py-3 text-xs font-mono">{{ item.upstream_order_no || '-' }}</TableCell>
                    <TableCell class="px-4 py-3 text-xs">{{ item.local_status || '-' }}</TableCell>
                    <TableCell class="px-4 py-3 text-xs">{{ item.upstream_status || '-' }}</TableCell>
                    <TableCell class="px-4 py-3">
                      <span
                        class="inline-flex rounded-full border px-2 py-0.5 text-xs"
                        :class="mismatchBadgeClass(item.mismatch_type)"
                      >
                        {{ t('reconciliation.mismatchType.' + item.mismatch_type) }}
                      </span>
                    </TableCell>
                    <TableCell class="px-4 py-3">
                      <span
                        class="inline-flex rounded-full border px-2 py-0.5 text-xs"
                        :class="item.resolved ? 'text-emerald-700 border-emerald-200 bg-emerald-50' : 'text-yellow-700 border-yellow-200 bg-yellow-50'"
                      >
                        {{ item.resolved ? t('reconciliation.items.resolved') : t('reconciliation.items.unresolved') }}
                      </span>
                    </TableCell>
                    <TableCell class="px-4 py-3 text-right">
                      <Button
                        v-if="!item.resolved"
                        size="sm"
                        variant="outline"
                        @click="openResolve(item)"
                      >
                        {{ t('reconciliation.items.resolve') }}
                      </Button>
                      <span v-else class="text-xs text-muted-foreground">{{ item.remark || '-' }}</span>
                    </TableCell>
                  </TableRow>
                </TableBody>
              </Table>
              <div v-if="detailItemsTotal > 20" class="flex items-center justify-end gap-2 border-t border-border px-4 py-3">
                <Button variant="outline" size="sm" class="h-7" :disabled="detailItemsPage <= 1" @click="loadDetailItems(detailItemsPage - 1)">
                  {{ t('admin.common.prevPage') }}
                </Button>
                <span class="text-xs text-muted-foreground">{{ detailItemsPage }}</span>
                <Button variant="outline" size="sm" class="h-7" :disabled="detailItemsPage * 20 >= detailItemsTotal" @click="loadDetailItems(detailItemsPage + 1)">
                  {{ t('admin.common.nextPage') }}
                </Button>
              </div>
            </div>
          </div>

          <div class="flex justify-end border-t border-border pt-4">
            <Button variant="outline" @click="closeDetail">{{ t('admin.common.cancel') }}</Button>
          </div>
        </div>
      </DialogScrollContent>
    </Dialog>

    <!-- Resolve Dialog -->
    <Dialog v-model:open="showResolve">
      <DialogScrollContent class="w-full max-w-md" @interact-outside="(e: Event) => e.preventDefault()">
        <DialogHeader>
          <DialogTitle>{{ t('reconciliation.items.resolve') }}</DialogTitle>
        </DialogHeader>
        <div class="space-y-4 pt-2">
          <div>
            <label class="mb-1.5 block text-xs font-medium text-muted-foreground">{{ t('reconciliation.items.remark') }}</label>
            <Textarea v-model="resolveRemark" :placeholder="t('reconciliation.items.remarkPlaceholder')" rows="3" />
          </div>
          <div class="flex justify-end gap-3 pt-2">
            <Button variant="outline" @click="showResolve = false">{{ t('admin.common.cancel') }}</Button>
            <Button :disabled="resolving" @click="handleResolve">{{ t('reconciliation.items.resolve') }}</Button>
          </div>
        </div>
      </DialogScrollContent>
    </Dialog>
  </div>
</template>
