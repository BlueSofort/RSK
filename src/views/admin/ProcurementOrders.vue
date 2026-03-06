<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminAPI } from '@/api/admin'
import { getLocalizedText, formatMoney } from '@/utils/format'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Dialog, DialogScrollContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { confirmAction } from '@/utils/confirm'
import { notifyError, notifySuccess } from '@/utils/notify'

const { t } = useI18n()
const loading = ref(true)
const orders = ref<any[]>([])
const connections = ref<any[]>([])
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0,
  total_page: 1,
})
const jumpPage = ref('')

const filters = reactive({
  status: '__all__',
  connection_id: '__all__',
  order_no: '',
  upstream_order_no: '',
  created_from: '',
  created_to: '',
})

const showDetail = ref(false)
const detailOrder = ref<any>(null)
const detailLoading = ref(false)
const retryingId = ref<number | null>(null)
const cancelingId = ref<number | null>(null)

const statusOptions = [
  { value: '__all__', key: 'procurement.filters.allStatus' },
  { value: 'pending', key: 'procurement.status.pending' },
  { value: 'accepted', key: 'procurement.status.accepted' },
  { value: 'rejected', key: 'procurement.status.rejected' },
  { value: 'failed', key: 'procurement.status.failed' },
  { value: 'fulfilled', key: 'procurement.status.fulfilled' },
  { value: 'canceled', key: 'procurement.status.canceled' },
]

// Stats
const stats = computed(() => {
  const all = orders.value
  return {
    total: pagination.total,
    pending: all.filter((o: any) => o.status === 'pending').length,
    failed: all.filter((o: any) => o.status === 'failed' || o.status === 'rejected').length,
    fulfilled: all.filter((o: any) => o.status === 'fulfilled').length,
  }
})

const fetchConnections = async () => {
  try {
    const res = await adminAPI.getSiteConnections({ page: 1, page_size: 100 })
    connections.value = (res.data.data as any[]) || []
  } catch { /* ignore */ }
}

const fetchOrders = async (page = 1) => {
  loading.value = true
  try {
    const params: any = {
      page,
      page_size: pagination.page_size,
    }
    if (filters.status && filters.status !== '__all__') params.status = filters.status
    if (filters.connection_id && filters.connection_id !== '__all__') params.connection_id = filters.connection_id
    if (filters.order_no) params.order_no = filters.order_no
    if (filters.upstream_order_no) params.upstream_order_no = filters.upstream_order_no
    if (filters.created_from) params.created_from = filters.created_from
    if (filters.created_to) params.created_to = filters.created_to

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
  detailLoading.value = true
  showDetail.value = true
  detailOrder.value = order
  try {
    const res = await adminAPI.getProcurementOrder(order.id)
    detailOrder.value = res.data.data || order
  } catch { /* keep original */ }
  detailLoading.value = false
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
    if (showDetail.value && detailOrder.value?.id === order.id) {
      openDetail(order)
    }
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
    if (showDetail.value && detailOrder.value?.id === order.id) {
      openDetail(order)
    }
  } catch (err: any) {
    notifyError(err?.response?.data?.message || err?.message)
  } finally {
    cancelingId.value = null
  }
}

const statusBadgeClass = (status: string) => {
  switch (status) {
    case 'pending': return 'text-amber-700 border-amber-200 bg-amber-50'
    case 'submitted': return 'text-blue-700 border-blue-200 bg-blue-50'
    case 'accepted': return 'text-sky-700 border-sky-200 bg-sky-50'
    case 'rejected': return 'text-red-700 border-red-200 bg-red-50'
    case 'failed': return 'text-red-700 border-red-200 bg-red-50'
    case 'fulfilled': return 'text-emerald-700 border-emerald-200 bg-emerald-50'
    case 'completed': return 'text-emerald-700 border-emerald-200 bg-emerald-50'
    case 'canceled': return 'text-muted-foreground border-border bg-muted/30'
    default: return 'text-muted-foreground border-border bg-muted/30'
  }
}

const statusIcon = (status: string) => {
  switch (status) {
    case 'pending': return 'M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z' // clock
    case 'accepted': return 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z' // check circle
    case 'rejected': case 'failed': return 'M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z' // exclamation
    case 'fulfilled': case 'completed': return 'M5 13l4 4L19 7' // check
    case 'canceled': return 'M6 18L18 6M6 6l12 12' // x
    default: return 'M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01'
  }
}

const formatTime = (raw?: string) => {
  if (!raw) return '-'
  const d = new Date(raw)
  if (Number.isNaN(d.getTime())) return '-'
  return d.toLocaleString()
}

const relativeTime = (raw?: string) => {
  if (!raw) return ''
  const d = new Date(raw)
  if (Number.isNaN(d.getTime())) return ''
  const diff = Date.now() - d.getTime()
  const minutes = Math.floor(diff / 60000)
  if (minutes < 1) return t('procurement.time.justNow')
  if (minutes < 60) return t('procurement.time.minutesAgo', { n: minutes })
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return t('procurement.time.hoursAgo', { n: hours })
  const days = Math.floor(hours / 24)
  return t('procurement.time.daysAgo', { n: days })
}

const getOrderTitle = (order: any) => {
  const localOrder = order.local_order
  if (!localOrder) return order.local_order_no || '-'
  const items = localOrder.items || []
  if (items.length > 0 && items[0].title) {
    return getLocalizedText(items[0].title)
  }
  return order.local_order_no || '-'
}

const profitAmount = (order: any) => {
  const sell = parseFloat(order.local_sell_amount || '0')
  const cost = parseFloat(order.upstream_amount || '0')
  if (!cost || !sell) return null
  return (sell - cost).toFixed(2)
}

const profitClass = (order: any) => {
  const p = profitAmount(order)
  if (p === null) return ''
  return parseFloat(p) >= 0 ? 'text-emerald-600' : 'text-red-600'
}

const canRetry = (status: string) => ['failed', 'rejected'].includes(status)
const canCancel = (status: string) => ['pending', 'submitted', 'accepted', 'failed'].includes(status)

onMounted(() => {
  fetchConnections()
  fetchOrders()
})
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-semibold">{{ t('procurement.title') }}</h1>
        <p class="mt-1 text-sm text-muted-foreground">{{ t('procurement.subtitle') }}</p>
      </div>
      <Button size="sm" variant="outline" @click="fetchOrders(pagination.page)">
        {{ t('admin.common.refresh') }}
      </Button>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-2 gap-3 md:grid-cols-4">
      <div class="rounded-xl border border-border bg-card p-4">
        <div class="text-xs font-medium text-muted-foreground">{{ t('procurement.stats.total') }}</div>
        <div class="mt-1 text-2xl font-bold">{{ stats.total }}</div>
      </div>
      <div class="rounded-xl border border-amber-200 bg-amber-50/50 p-4">
        <div class="text-xs font-medium text-amber-700">{{ t('procurement.stats.pending') }}</div>
        <div class="mt-1 text-2xl font-bold text-amber-700">{{ stats.pending }}</div>
      </div>
      <div class="rounded-xl border border-red-200 bg-red-50/50 p-4">
        <div class="text-xs font-medium text-red-700">{{ t('procurement.stats.failed') }}</div>
        <div class="mt-1 text-2xl font-bold text-red-700">{{ stats.failed }}</div>
      </div>
      <div class="rounded-xl border border-emerald-200 bg-emerald-50/50 p-4">
        <div class="text-xs font-medium text-emerald-700">{{ t('procurement.stats.fulfilled') }}</div>
        <div class="mt-1 text-2xl font-bold text-emerald-700">{{ stats.fulfilled }}</div>
      </div>
    </div>

    <!-- Filters -->
    <div class="flex flex-wrap items-end gap-3">
      <div>
        <label class="mb-1.5 block text-xs font-medium text-muted-foreground">{{ t('procurement.columns.status') }}</label>
        <Select v-model="filters.status" @update:model-value="handleSearch">
          <SelectTrigger class="h-9 w-36">
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
        <label class="mb-1.5 block text-xs font-medium text-muted-foreground">{{ t('procurement.columns.connection') }}</label>
        <Select v-model="filters.connection_id" @update:model-value="handleSearch">
          <SelectTrigger class="h-9 w-44">
            <SelectValue :placeholder="t('procurement.filters.allConnections')" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="__all__">{{ t('procurement.filters.allConnections') }}</SelectItem>
            <SelectItem v-for="conn in connections" :key="conn.id" :value="String(conn.id)">
              {{ conn.name || `#${conn.id}` }}
            </SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div>
        <label class="mb-1.5 block text-xs font-medium text-muted-foreground">{{ t('procurement.filters.orderNo') }}</label>
        <Input v-model="filters.order_no" class="h-9 w-44" :placeholder="t('procurement.filters.orderNoPlaceholder')" @keyup.enter="handleSearch" />
      </div>
      <div>
        <label class="mb-1.5 block text-xs font-medium text-muted-foreground">{{ t('procurement.filters.upstreamOrderNo') }}</label>
        <Input v-model="filters.upstream_order_no" class="h-9 w-44" :placeholder="t('procurement.filters.upstreamOrderNoPlaceholder')" @keyup.enter="handleSearch" />
      </div>
      <div>
        <label class="mb-1.5 block text-xs font-medium text-muted-foreground">{{ t('procurement.filters.dateRange') }}</label>
        <div class="flex items-center gap-1.5">
          <Input v-model="filters.created_from" type="date" class="h-9 w-36" @change="handleSearch" />
          <span class="text-xs text-muted-foreground">—</span>
          <Input v-model="filters.created_to" type="date" class="h-9 w-36" @change="handleSearch" />
        </div>
      </div>
      <Button size="sm" class="h-9" @click="handleSearch">{{ t('procurement.filters.search') }}</Button>
    </div>

    <!-- Order Cards -->
    <div class="space-y-3">
      <div v-if="loading" class="rounded-xl border border-border bg-card p-12 text-center text-muted-foreground">
        {{ t('admin.common.loading') }}
      </div>
      <div v-else-if="orders.length === 0" class="rounded-xl border border-border bg-card p-12 text-center text-muted-foreground">
        {{ t('procurement.empty') }}
      </div>

      <div
        v-for="order in orders"
        :key="order.id"
        class="group cursor-pointer rounded-xl border border-border bg-card transition-all hover:shadow-md"
        :class="{ 'border-red-200 bg-red-50/30': order.status === 'failed' || order.status === 'rejected' }"
        @click="openDetail(order)"
      >
        <div class="p-4">
          <!-- Top row: status + ID + time -->
          <div class="flex items-start justify-between gap-3">
            <div class="flex items-center gap-3">
              <div
                class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full border"
                :class="statusBadgeClass(order.status)"
              >
                <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path :d="statusIcon(order.status)" />
                </svg>
              </div>
              <div>
                <div class="flex items-center gap-2">
                  <span class="text-sm font-semibold text-foreground">{{ getOrderTitle(order) }}</span>
                  <span
                    class="inline-flex rounded-full border px-2 py-0.5 text-[11px] font-medium"
                    :class="statusBadgeClass(order.status)"
                  >
                    {{ t('procurement.status.' + order.status) }}
                  </span>
                </div>
                <div class="mt-0.5 flex items-center gap-3 text-xs text-muted-foreground">
                  <span class="font-mono">#{{ order.id }}</span>
                  <span>{{ order.connection?.name || '-' }}</span>
                  <span>{{ relativeTime(order.created_at) }}</span>
                </div>
              </div>
            </div>
            <div class="flex shrink-0 items-center gap-2" @click.stop>
              <Button
                v-if="canRetry(order.status)"
                size="sm"
                variant="outline"
                class="h-7 text-xs"
                :disabled="retryingId === order.id"
                @click="handleRetry(order)"
              >
                {{ t('procurement.actions.retry') }}
              </Button>
              <Button
                v-if="canCancel(order.status)"
                size="sm"
                variant="ghost"
                class="h-7 text-xs text-muted-foreground hover:text-red-600"
                :disabled="cancelingId === order.id"
                @click="handleCancel(order)"
              >
                {{ t('procurement.actions.cancel') }}
              </Button>
            </div>
          </div>

          <!-- Info row -->
          <div class="mt-3 grid grid-cols-2 gap-x-6 gap-y-1.5 text-xs md:grid-cols-4 lg:grid-cols-6">
            <div>
              <span class="text-muted-foreground">{{ t('procurement.columns.localOrderNo') }}</span>
              <div class="mt-0.5 font-mono text-foreground">{{ order.local_order_no || '-' }}</div>
            </div>
            <div>
              <span class="text-muted-foreground">{{ t('procurement.columns.upstreamOrderNo') }}</span>
              <div class="mt-0.5 font-mono text-foreground">{{ order.upstream_order_no || '-' }}</div>
            </div>
            <div>
              <span class="text-muted-foreground">{{ t('procurement.columns.localSellAmount') }}</span>
              <div class="mt-0.5 font-medium text-foreground">{{ formatMoney(order.local_sell_amount, order.currency) }}</div>
            </div>
            <div>
              <span class="text-muted-foreground">{{ t('procurement.columns.upstreamAmount') }}</span>
              <div class="mt-0.5 font-medium text-foreground">{{ order.upstream_amount && order.upstream_amount !== '0.00' ? formatMoney(order.upstream_amount, order.currency) : '-' }}</div>
            </div>
            <div>
              <span class="text-muted-foreground">{{ t('procurement.detail.profit') }}</span>
              <div class="mt-0.5 font-medium" :class="profitClass(order)">
                {{ profitAmount(order) !== null ? formatMoney(profitAmount(order)!, order.currency) : '-' }}
              </div>
            </div>
            <div>
              <span class="text-muted-foreground">{{ t('procurement.columns.retryCount') }}</span>
              <div class="mt-0.5 text-foreground">{{ order.retry_count ?? 0 }}</div>
            </div>
          </div>

          <!-- Error message -->
          <div
            v-if="order.error_message"
            class="mt-3 flex items-start gap-2 rounded-lg border border-red-200 bg-red-50 px-3 py-2"
          >
            <svg class="mt-0.5 h-3.5 w-3.5 shrink-0 text-red-500" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            <span class="text-xs text-red-700">{{ order.error_message }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Pagination -->
    <div v-if="pagination.total_page > 1" class="flex flex-wrap items-center justify-between gap-3">
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

    <!-- Detail Dialog -->
    <Dialog v-model:open="showDetail" @update:open="(value: boolean) => { if (!value) closeDetail() }">
      <DialogScrollContent class="w-full max-w-3xl" @interact-outside="(e: Event) => e.preventDefault()">
        <DialogHeader>
          <DialogTitle class="flex items-center gap-2">
            {{ t('procurement.detail.title') }}
            <span v-if="detailOrder" class="font-mono text-sm text-muted-foreground">#{{ detailOrder.id }}</span>
          </DialogTitle>
        </DialogHeader>

        <div v-if="detailOrder" class="space-y-5">
          <!-- Status banner -->
          <div
            class="flex items-center gap-3 rounded-lg border px-4 py-3"
            :class="statusBadgeClass(detailOrder.status)"
          >
            <svg class="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path :d="statusIcon(detailOrder.status)" />
            </svg>
            <div>
              <div class="text-sm font-semibold">{{ t('procurement.status.' + detailOrder.status) }}</div>
              <div v-if="detailOrder.error_message" class="mt-0.5 text-xs opacity-80">{{ detailOrder.error_message }}</div>
            </div>
          </div>

          <!-- Order info -->
          <div class="rounded-lg border border-border">
            <div class="border-b border-border bg-muted/30 px-4 py-2 text-xs font-semibold text-muted-foreground uppercase">
              {{ t('procurement.detail.orderInfo') }}
            </div>
            <div class="grid grid-cols-2 gap-4 p-4 md:grid-cols-3">
              <div>
                <div class="text-xs text-muted-foreground">{{ t('procurement.columns.localOrderNo') }}</div>
                <div class="mt-1 text-sm font-mono font-medium">{{ detailOrder.local_order_no || '-' }}</div>
              </div>
              <div>
                <div class="text-xs text-muted-foreground">{{ t('procurement.columns.upstreamOrderNo') }}</div>
                <div class="mt-1 text-sm font-mono font-medium">{{ detailOrder.upstream_order_no || '-' }}</div>
              </div>
              <div>
                <div class="text-xs text-muted-foreground">{{ t('procurement.detail.upstreamOrderId') }}</div>
                <div class="mt-1 text-sm font-mono">{{ detailOrder.upstream_order_id || '-' }}</div>
              </div>
              <div>
                <div class="text-xs text-muted-foreground">{{ t('procurement.columns.connection') }}</div>
                <div class="mt-1 text-sm font-medium">{{ detailOrder.connection?.name || detailOrder.connection_id || '-' }}</div>
              </div>
              <div>
                <div class="text-xs text-muted-foreground">{{ t('procurement.detail.traceId') }}</div>
                <div class="mt-1 text-sm font-mono text-muted-foreground break-all">{{ detailOrder.trace_id || '-' }}</div>
              </div>
              <div>
                <div class="text-xs text-muted-foreground">{{ t('procurement.detail.currency') }}</div>
                <div class="mt-1 text-sm">{{ detailOrder.currency || '-' }}</div>
              </div>
            </div>
          </div>

          <!-- Financial -->
          <div class="rounded-lg border border-border">
            <div class="border-b border-border bg-muted/30 px-4 py-2 text-xs font-semibold text-muted-foreground uppercase">
              {{ t('procurement.detail.financial') }}
            </div>
            <div class="grid grid-cols-3 divide-x divide-border">
              <div class="p-4 text-center">
                <div class="text-xs text-muted-foreground">{{ t('procurement.columns.localSellAmount') }}</div>
                <div class="mt-1 text-lg font-bold">{{ formatMoney(detailOrder.local_sell_amount, detailOrder.currency) }}</div>
              </div>
              <div class="p-4 text-center">
                <div class="text-xs text-muted-foreground">{{ t('procurement.columns.upstreamAmount') }}</div>
                <div class="mt-1 text-lg font-bold">
                  {{ detailOrder.upstream_amount && detailOrder.upstream_amount !== '0.00' ? formatMoney(detailOrder.upstream_amount, detailOrder.currency) : '-' }}
                </div>
              </div>
              <div class="p-4 text-center">
                <div class="text-xs text-muted-foreground">{{ t('procurement.detail.profit') }}</div>
                <div class="mt-1 text-lg font-bold" :class="profitClass(detailOrder)">
                  {{ profitAmount(detailOrder) !== null ? formatMoney(profitAmount(detailOrder)!, detailOrder.currency) : '-' }}
                </div>
              </div>
            </div>
          </div>

          <!-- Local Order Info -->
          <div v-if="detailOrder.local_order" class="rounded-lg border border-border">
            <div class="border-b border-border bg-muted/30 px-4 py-2 text-xs font-semibold text-muted-foreground uppercase">
              {{ t('procurement.detail.localOrder') }}
            </div>
            <div class="p-4 space-y-2">
              <div v-if="detailOrder.local_order.items?.length" class="space-y-2">
                <div v-for="(item, idx) in detailOrder.local_order.items" :key="idx" class="flex items-center gap-3 text-sm">
                  <div class="flex-1">
                    <span class="font-medium">{{ getLocalizedText(item.title) }}</span>
                    <span v-if="item.sku_snapshot?.sku_code" class="ml-2 text-xs text-muted-foreground">
                      SKU: {{ item.sku_snapshot.sku_code }}
                    </span>
                  </div>
                  <span class="text-muted-foreground">x{{ item.quantity }}</span>
                  <span class="font-mono">{{ formatMoney(item.total_amount, detailOrder.currency) }}</span>
                </div>
              </div>
              <div class="flex items-center gap-4 text-xs text-muted-foreground pt-1 border-t border-border">
                <span>{{ t('procurement.detail.orderStatus') }}: {{ detailOrder.local_order.status }}</span>
                <span v-if="detailOrder.local_order.user_email">{{ detailOrder.local_order.user_email }}</span>
              </div>
            </div>
          </div>

          <!-- Timeline -->
          <div class="rounded-lg border border-border">
            <div class="border-b border-border bg-muted/30 px-4 py-2 text-xs font-semibold text-muted-foreground uppercase">
              {{ t('procurement.detail.timeline') }}
            </div>
            <div class="p-4 space-y-3">
              <div class="flex items-start gap-3 text-xs">
                <div class="w-32 shrink-0 text-muted-foreground">{{ t('procurement.columns.createdAt') }}</div>
                <div>{{ formatTime(detailOrder.created_at) }}</div>
              </div>
              <div class="flex items-start gap-3 text-xs">
                <div class="w-32 shrink-0 text-muted-foreground">{{ t('procurement.detail.updatedAt') }}</div>
                <div>{{ formatTime(detailOrder.updated_at) }}</div>
              </div>
              <div v-if="detailOrder.next_retry_at" class="flex items-start gap-3 text-xs">
                <div class="w-32 shrink-0 text-muted-foreground">{{ t('procurement.detail.nextRetryAt') }}</div>
                <div>{{ formatTime(detailOrder.next_retry_at) }}</div>
              </div>
              <div class="flex items-start gap-3 text-xs">
                <div class="w-32 shrink-0 text-muted-foreground">{{ t('procurement.columns.retryCount') }}</div>
                <div>{{ detailOrder.retry_count ?? 0 }}</div>
              </div>
            </div>
          </div>

          <!-- Error detail -->
          <div v-if="detailOrder.error_message" class="rounded-lg border border-red-200 bg-red-50/50">
            <div class="border-b border-red-200 bg-red-50 px-4 py-2 text-xs font-semibold text-red-700 uppercase">
              {{ t('procurement.columns.errorMessage') }}
            </div>
            <div class="p-4 text-sm font-mono whitespace-pre-wrap break-all text-red-700">
              {{ detailOrder.error_message }}
            </div>
          </div>

          <!-- Upstream payload -->
          <div v-if="detailOrder.upstream_payload" class="rounded-lg border border-border">
            <div class="border-b border-border bg-muted/30 px-4 py-2 text-xs font-semibold text-muted-foreground uppercase">
              {{ t('procurement.detail.upstreamPayload') }}
            </div>
            <div class="max-h-48 overflow-y-auto p-4 text-xs font-mono whitespace-pre-wrap break-all text-muted-foreground">
              {{ detailOrder.upstream_payload }}
            </div>
          </div>

          <!-- Actions -->
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
