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
import { confirmAction } from '@/utils/confirm'
import { notifyError, notifySuccess } from '@/utils/notify'

const { t } = useI18n()
const loading = ref(true)
const orders = ref<any[]>([])
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0,
  total_page: 1,
})
const jumpPage = ref('')

const filters = reactive({
  status: '__all__',
  connection_id: '',
  order_no: '',
})

const showDetail = ref(false)
const detailOrder = ref<any>(null)
const retryingId = ref<number | null>(null)
const cancelingId = ref<number | null>(null)

const statusOptions = [
  { value: '__all__', key: 'procurement.filters.allStatus' },
  { value: 'pending', key: 'procurement.status.pending' },
  { value: 'submitted', key: 'procurement.status.submitted' },
  { value: 'accepted', key: 'procurement.status.accepted' },
  { value: 'rejected', key: 'procurement.status.rejected' },
  { value: 'failed', key: 'procurement.status.failed' },
  { value: 'fulfilled', key: 'procurement.status.fulfilled' },
  { value: 'completed', key: 'procurement.status.completed' },
  { value: 'canceled', key: 'procurement.status.canceled' },
]

const fetchOrders = async (page = 1) => {
  loading.value = true
  try {
    const params: any = {
      page,
      page_size: pagination.page_size,
    }
    if (filters.status && filters.status !== '__all__') params.status = filters.status
    if (filters.connection_id) params.connection_id = filters.connection_id
    if (filters.order_no) params.order_no = filters.order_no

    const res = await adminAPI.getProcurementOrders(params)
    orders.value = (res.data.data as any[]) || []
    const p = (res.data as any).pagination
    if (p) {
      pagination.page = p.page
      pagination.page_size = p.page_size
      pagination.total = p.total
      pagination.total_page = p.total_page
    }
  } catch {
    orders.value = []
  } finally {
    loading.value = false
  }
}

const changePage = (page: number) => {
  if (page < 1 || page > pagination.total_page) return
  fetchOrders(page)
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
  fetchOrders(1)
}

const openDetail = async (order: any) => {
  try {
    const res = await adminAPI.getProcurementOrder(order.id)
    detailOrder.value = res.data.data || order
  } catch {
    detailOrder.value = order
  }
  showDetail.value = true
}

const closeDetail = () => {
  showDetail.value = false
  detailOrder.value = null
}

const handleRetry = async (order: any) => {
  const confirmed = await confirmAction({
    description: t('procurement.actions.retryConfirm', { id: order.id }),
    confirmText: t('procurement.actions.retry'),
  })
  if (!confirmed) return
  retryingId.value = order.id
  try {
    await adminAPI.retryProcurementOrder(order.id)
    notifySuccess(t('procurement.actions.retrySuccess'))
    fetchOrders(pagination.page)
  } catch (err: any) {
    notifyError(err?.response?.data?.message || err?.message)
  } finally {
    retryingId.value = null
  }
}

const handleCancel = async (order: any) => {
  const confirmed = await confirmAction({
    description: t('procurement.actions.cancelConfirm', { id: order.id }),
    confirmText: t('procurement.actions.cancel'),
    variant: 'destructive',
  })
  if (!confirmed) return
  cancelingId.value = order.id
  try {
    await adminAPI.cancelProcurementOrder(order.id)
    notifySuccess(t('procurement.actions.cancelSuccess'))
    fetchOrders(pagination.page)
  } catch (err: any) {
    notifyError(err?.response?.data?.message || err?.message)
  } finally {
    cancelingId.value = null
  }
}

const statusBadgeClass = (status: string) => {
  switch (status) {
    case 'pending':
      return 'text-yellow-700 border-yellow-200 bg-yellow-50'
    case 'submitted':
      return 'text-blue-700 border-blue-200 bg-blue-50'
    case 'accepted':
      return 'text-blue-700 border-blue-200 bg-blue-50'
    case 'rejected':
      return 'text-red-700 border-red-200 bg-red-50'
    case 'failed':
      return 'text-red-700 border-red-200 bg-red-50'
    case 'fulfilled':
      return 'text-emerald-700 border-emerald-200 bg-emerald-50'
    case 'completed':
      return 'text-emerald-700 border-emerald-200 bg-emerald-50'
    case 'canceled':
      return 'text-muted-foreground border-border bg-muted/30'
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

const truncate = (text: string, maxLen = 40) => {
  if (!text) return '-'
  return text.length > maxLen ? text.slice(0, maxLen) + '...' : text
}

const canRetry = (status: string) => ['failed', 'rejected'].includes(status)
const canCancel = (status: string) => ['pending', 'submitted', 'accepted'].includes(status)

onMounted(() => {
  fetchOrders()
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-semibold">{{ t('procurement.title') }}</h1>
    </div>

    <!-- Filters -->
    <div class="flex flex-wrap items-end gap-3">
      <div>
        <label class="mb-1.5 block text-xs font-medium text-muted-foreground">{{ t('procurement.columns.status') }}</label>
        <Select v-model="filters.status">
          <SelectTrigger class="h-9 w-40">
            <SelectValue :placeholder="t('procurement.filters.allStatus')" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem v-for="opt in statusOptions" :key="opt.value" :value="opt.value">
              {{ t(opt.key) }}
            </SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div>
        <label class="mb-1.5 block text-xs font-medium text-muted-foreground">{{ t('procurement.filters.connectionId') }}</label>
        <Input v-model="filters.connection_id" class="h-9 w-32" placeholder="ID" />
      </div>
      <div>
        <label class="mb-1.5 block text-xs font-medium text-muted-foreground">{{ t('procurement.filters.orderNo') }}</label>
        <Input v-model="filters.order_no" class="h-9 w-48" :placeholder="t('procurement.filters.orderNo')" />
      </div>
      <Button size="sm" class="h-9" @click="handleSearch">{{ t('admin.common.refresh') }}</Button>
    </div>

    <!-- Table -->
    <div class="rounded-xl border border-border bg-card">
      <Table>
        <TableHeader class="border-b border-border bg-muted/40 text-xs uppercase text-muted-foreground">
          <TableRow>
            <TableHead class="px-6 py-3">{{ t('procurement.columns.id') }}</TableHead>
            <TableHead class="px-6 py-3">{{ t('procurement.columns.connection') }}</TableHead>
            <TableHead class="px-6 py-3">{{ t('procurement.columns.localOrderNo') }}</TableHead>
            <TableHead class="px-6 py-3">{{ t('procurement.columns.upstreamOrderNo') }}</TableHead>
            <TableHead class="px-6 py-3">{{ t('procurement.columns.status') }}</TableHead>
            <TableHead class="px-6 py-3">{{ t('procurement.columns.upstreamAmount') }}</TableHead>
            <TableHead class="px-6 py-3">{{ t('procurement.columns.errorMessage') }}</TableHead>
            <TableHead class="px-6 py-3">{{ t('procurement.columns.retryCount') }}</TableHead>
            <TableHead class="px-6 py-3">{{ t('procurement.columns.createdAt') }}</TableHead>
            <TableHead class="px-6 py-3 text-right">{{ t('procurement.columns.actions') }}</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody class="divide-y divide-border">
          <TableRow v-if="loading">
            <TableCell colspan="10" class="px-6 py-8 text-center text-muted-foreground">{{ t('admin.common.loading') }}</TableCell>
          </TableRow>
          <TableRow v-else-if="orders.length === 0">
            <TableCell colspan="10" class="px-6 py-8 text-center text-muted-foreground">{{ t('procurement.empty') }}</TableCell>
          </TableRow>
          <TableRow
            v-for="order in orders"
            :key="order.id"
            class="cursor-pointer hover:bg-muted/30"
            @click="openDetail(order)"
          >
            <TableCell class="px-6 py-4">
              <IdCell :value="order.id" />
            </TableCell>
            <TableCell class="px-6 py-4 font-medium text-foreground">
              {{ order.connection?.name || order.connection_id || '-' }}
            </TableCell>
            <TableCell class="px-6 py-4 text-xs font-mono text-muted-foreground">
              {{ order.local_order_no || '-' }}
            </TableCell>
            <TableCell class="px-6 py-4 text-xs font-mono text-muted-foreground">
              {{ order.upstream_order_no || '-' }}
            </TableCell>
            <TableCell class="px-6 py-4">
              <span
                class="inline-flex rounded-full border px-2.5 py-1 text-xs"
                :class="statusBadgeClass(order.status)"
              >
                {{ t('procurement.status.' + order.status) }}
              </span>
            </TableCell>
            <TableCell class="px-6 py-4 text-xs text-muted-foreground">
              {{ order.upstream_amount || '-' }}
            </TableCell>
            <TableCell class="px-6 py-4 text-xs text-muted-foreground" :title="order.error_message">
              {{ truncate(order.error_message) }}
            </TableCell>
            <TableCell class="px-6 py-4 text-xs text-muted-foreground">
              {{ order.retry_count ?? 0 }}
            </TableCell>
            <TableCell class="px-6 py-4 text-xs text-muted-foreground">
              {{ formatTime(order.created_at) }}
            </TableCell>
            <TableCell class="px-6 py-4 text-right" @click.stop>
              <div class="flex items-center justify-end gap-2">
                <Button
                  v-if="canRetry(order.status)"
                  size="sm"
                  variant="outline"
                  :disabled="retryingId === order.id"
                  @click="handleRetry(order)"
                >
                  {{ t('procurement.actions.retry') }}
                </Button>
                <Button
                  v-if="canCancel(order.status)"
                  size="sm"
                  variant="destructive"
                  :disabled="cancelingId === order.id"
                  @click="handleCancel(order)"
                >
                  {{ t('procurement.actions.cancel') }}
                </Button>
              </div>
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

    <!-- Detail Dialog -->
    <Dialog v-model:open="showDetail" @update:open="(value: boolean) => { if (!value) closeDetail() }">
      <DialogScrollContent class="w-full max-w-2xl" @interact-outside="(e: Event) => e.preventDefault()">
        <DialogHeader>
          <DialogTitle>{{ t('procurement.detail.title') }}</DialogTitle>
        </DialogHeader>

        <div v-if="detailOrder" class="space-y-4">
          <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
            <div>
              <label class="mb-1 block text-xs font-medium text-muted-foreground">{{ t('procurement.columns.id') }}</label>
              <div class="text-sm">{{ detailOrder.id }}</div>
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-muted-foreground">{{ t('procurement.columns.connection') }}</label>
              <div class="text-sm">{{ detailOrder.connection?.name || detailOrder.connection_id || '-' }}</div>
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-muted-foreground">{{ t('procurement.columns.localOrderNo') }}</label>
              <div class="text-sm font-mono">{{ detailOrder.local_order_no || '-' }}</div>
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-muted-foreground">{{ t('procurement.columns.upstreamOrderNo') }}</label>
              <div class="text-sm font-mono">{{ detailOrder.upstream_order_no || '-' }}</div>
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-muted-foreground">{{ t('procurement.columns.status') }}</label>
              <div>
                <span
                  class="inline-flex rounded-full border px-2.5 py-1 text-xs"
                  :class="statusBadgeClass(detailOrder.status)"
                >
                  {{ t('procurement.status.' + detailOrder.status) }}
                </span>
              </div>
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-muted-foreground">{{ t('procurement.columns.upstreamAmount') }}</label>
              <div class="text-sm">{{ detailOrder.upstream_amount || '-' }}</div>
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-muted-foreground">{{ t('procurement.detail.currency') }}</label>
              <div class="text-sm">{{ detailOrder.currency || '-' }}</div>
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-muted-foreground">{{ t('procurement.columns.retryCount') }}</label>
              <div class="text-sm">{{ detailOrder.retry_count ?? 0 }}</div>
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-muted-foreground">{{ t('procurement.detail.traceId') }}</label>
              <div class="text-sm font-mono">{{ detailOrder.trace_id || '-' }}</div>
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-muted-foreground">{{ t('procurement.detail.upstreamOrderId') }}</label>
              <div class="text-sm font-mono">{{ detailOrder.upstream_order_id || '-' }}</div>
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-muted-foreground">{{ t('procurement.columns.createdAt') }}</label>
              <div class="text-sm">{{ formatTime(detailOrder.created_at) }}</div>
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-muted-foreground">{{ t('procurement.detail.updatedAt') }}</label>
              <div class="text-sm">{{ formatTime(detailOrder.updated_at) }}</div>
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-muted-foreground">{{ t('procurement.detail.nextRetryAt') }}</label>
              <div class="text-sm">{{ formatTime(detailOrder.next_retry_at) }}</div>
            </div>
          </div>

          <div v-if="detailOrder.error_message">
            <label class="mb-1 block text-xs font-medium text-muted-foreground">{{ t('procurement.columns.errorMessage') }}</label>
            <div class="rounded-md border border-border bg-muted/30 p-3 text-sm font-mono whitespace-pre-wrap break-all">
              {{ detailOrder.error_message }}
            </div>
          </div>

          <div v-if="detailOrder.upstream_response">
            <label class="mb-1 block text-xs font-medium text-muted-foreground">Upstream Response</label>
            <div class="rounded-md border border-border bg-muted/30 p-3 text-sm font-mono whitespace-pre-wrap break-all max-h-60 overflow-y-auto">
              {{ typeof detailOrder.upstream_response === 'string' ? detailOrder.upstream_response : JSON.stringify(detailOrder.upstream_response, null, 2) }}
            </div>
          </div>

          <div class="flex justify-end gap-3 border-t border-border pt-4">
            <Button
              v-if="canRetry(detailOrder.status)"
              variant="outline"
              :disabled="retryingId === detailOrder.id"
              @click="handleRetry(detailOrder)"
            >
              {{ t('procurement.actions.retry') }}
            </Button>
            <Button
              v-if="canCancel(detailOrder.status)"
              variant="destructive"
              :disabled="cancelingId === detailOrder.id"
              @click="handleCancel(detailOrder)"
            >
              {{ t('procurement.actions.cancel') }}
            </Button>
            <Button variant="outline" @click="closeDetail">{{ t('admin.common.cancel') }}</Button>
          </div>
        </div>
      </DialogScrollContent>
    </Dialog>
  </div>
</template>
