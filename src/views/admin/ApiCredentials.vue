<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminAPI } from '@/api/admin'
import IdCell from '@/components/IdCell.vue'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Dialog, DialogScrollContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { confirmAction } from '@/utils/confirm'
import { notifyError, notifySuccess } from '@/utils/notify'

const { t } = useI18n()
const loading = ref(true)
const credentials = ref<any[]>([])
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0,
  total_page: 1,
})
const jumpPage = ref('')

const filterStatus = ref('__all__')

const statusOptions = [
  { value: '__all__', key: 'apiCredentials.filters.allStatus' },
  { value: 'pending_review', key: 'apiCredentials.status.pending_review' },
  { value: 'approved', key: 'apiCredentials.status.approved' },
  { value: 'rejected', key: 'apiCredentials.status.rejected' },
]

// Detail dialog
const showDetail = ref(false)
const detailCred = ref<any>(null)

// Reject dialog
const showReject = ref(false)
const rejectId = ref<number>(0)
const rejectReason = ref('')

// Secret dialog (after approve)
const showSecret = ref(false)
const secretData = ref<{ api_key: string; api_secret: string }>({ api_key: '', api_secret: '' })

const fetchCredentials = async (page = 1) => {
  loading.value = true
  try {
    const params: any = { page, page_size: pagination.page_size }
    if (filterStatus.value && filterStatus.value !== '__all__') params.status = filterStatus.value
    const res = await adminAPI.getApiCredentials(params)
    credentials.value = res.data?.data || []
    const pg = res.data?.pagination
    if (pg) {
      pagination.page = pg.page
      pagination.page_size = pg.page_size
      pagination.total = pg.total
      pagination.total_page = pg.total_page
    }
  } catch (e: any) {
    notifyError(e?.response?.data?.message || e.message)
  } finally {
    loading.value = false
  }
}

const statusVariant = (status: string) => {
  switch (status) {
    case 'approved': return 'default'
    case 'pending_review': return 'secondary'
    case 'rejected': return 'destructive'
    default: return 'outline'
  }
}

const formatDate = (d: string) => d ? new Date(d).toLocaleString() : '-'

const openDetail = (cred: any) => {
  detailCred.value = cred
  showDetail.value = true
}

const handleApprove = async (id: number) => {
  const ok = await confirmAction(t('apiCredentials.approve.confirm'))
  if (!ok) return
  try {
    const res = await adminAPI.approveApiCredential(id)
    const data = res.data?.data
    if (data?.api_key && data?.api_secret) {
      secretData.value = { api_key: data.api_key, api_secret: data.api_secret }
      showSecret.value = true
    }
    notifySuccess(t('apiCredentials.approve.success'))
    fetchCredentials(pagination.page)
  } catch (e: any) {
    notifyError(e?.response?.data?.message || e.message)
  }
}

const openReject = (id: number) => {
  rejectId.value = id
  rejectReason.value = ''
  showReject.value = true
}

const handleReject = async () => {
  if (!rejectReason.value.trim()) return
  try {
    await adminAPI.rejectApiCredential(rejectId.value, { reason: rejectReason.value.trim() })
    notifySuccess(t('apiCredentials.reject.success'))
    showReject.value = false
    fetchCredentials(pagination.page)
  } catch (e: any) {
    notifyError(e?.response?.data?.message || e.message)
  }
}

const handleToggle = async (cred: any) => {
  const msg = cred.is_active
    ? t('apiCredentials.toggle.disableConfirm')
    : t('apiCredentials.toggle.enableConfirm')
  const ok = await confirmAction(msg)
  if (!ok) return
  try {
    await adminAPI.updateApiCredentialStatus(cred.id, { is_active: !cred.is_active })
    notifySuccess(t('apiCredentials.toggle.success'))
    fetchCredentials(pagination.page)
  } catch (e: any) {
    notifyError(e?.response?.data?.message || e.message)
  }
}

const handleDelete = async (id: number) => {
  const ok = await confirmAction(t('apiCredentials.delete.confirm'))
  if (!ok) return
  try {
    await adminAPI.deleteApiCredential(id)
    notifySuccess(t('apiCredentials.delete.success'))
    fetchCredentials(pagination.page)
  } catch (e: any) {
    notifyError(e?.response?.data?.message || e.message)
  }
}

const goPage = (page: number) => {
  if (page < 1 || page > pagination.total_page) return
  fetchCredentials(page)
}

const handleJump = () => {
  const p = parseInt(jumpPage.value)
  if (p >= 1 && p <= pagination.total_page) {
    fetchCredentials(p)
  }
  jumpPage.value = ''
}

onMounted(() => fetchCredentials())
</script>

<template>
  <div class="space-y-4 p-4 md:p-6">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <h1 class="text-2xl font-bold">{{ t('apiCredentials.title') }}</h1>
    </div>

    <!-- Filters -->
    <div class="flex flex-wrap gap-3 items-center">
      <Select v-model="filterStatus" @update:model-value="fetchCredentials(1)">
        <SelectTrigger class="w-[160px]">
          <SelectValue :placeholder="t('apiCredentials.filters.statusPlaceholder')" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem v-for="opt in statusOptions" :key="opt.value" :value="opt.value">
            {{ t(opt.key) }}
          </SelectItem>
        </SelectContent>
      </Select>
    </div>

    <!-- Table -->
    <div class="rounded-md border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>{{ t('apiCredentials.columns.id') }}</TableHead>
            <TableHead>{{ t('apiCredentials.columns.userId') }}</TableHead>
            <TableHead>{{ t('apiCredentials.columns.apiKey') }}</TableHead>
            <TableHead>{{ t('apiCredentials.columns.status') }}</TableHead>
            <TableHead>{{ t('apiCredentials.columns.isActive') }}</TableHead>
            <TableHead>{{ t('apiCredentials.columns.lastUsedAt') }}</TableHead>
            <TableHead>{{ t('apiCredentials.columns.createdAt') }}</TableHead>
            <TableHead>{{ t('apiCredentials.columns.actions') }}</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-if="loading">
            <TableCell :colspan="8" class="text-center py-8 text-muted-foreground">Loading...</TableCell>
          </TableRow>
          <TableRow v-else-if="credentials.length === 0">
            <TableCell :colspan="8" class="text-center py-8 text-muted-foreground">{{ t('apiCredentials.empty') }}</TableCell>
          </TableRow>
          <TableRow v-for="cred in credentials" :key="cred.id">
            <TableCell><IdCell :id="cred.id" /></TableCell>
            <TableCell>{{ cred.user_id }}</TableCell>
            <TableCell class="font-mono text-xs max-w-[200px] truncate">{{ cred.api_key || '-' }}</TableCell>
            <TableCell>
              <Badge :variant="statusVariant(cred.status)">
                {{ t(`apiCredentials.status.${cred.status}`, cred.status) }}
              </Badge>
            </TableCell>
            <TableCell>
              <Badge :variant="cred.is_active ? 'default' : 'outline'">
                {{ cred.is_active ? 'ON' : 'OFF' }}
              </Badge>
            </TableCell>
            <TableCell class="text-xs">{{ formatDate(cred.last_used_at) }}</TableCell>
            <TableCell class="text-xs">{{ formatDate(cred.created_at) }}</TableCell>
            <TableCell>
              <div class="flex flex-wrap gap-1">
                <Button size="xs" variant="outline" @click="openDetail(cred)">
                  {{ t('apiCredentials.actions.detail') }}
                </Button>
                <Button
                  v-if="cred.status === 'pending_review'"
                  size="xs"
                  variant="default"
                  @click="handleApprove(cred.id)"
                >
                  {{ t('apiCredentials.actions.approve') }}
                </Button>
                <Button
                  v-if="cred.status === 'pending_review'"
                  size="xs"
                  variant="destructive"
                  @click="openReject(cred.id)"
                >
                  {{ t('apiCredentials.actions.reject') }}
                </Button>
                <Button
                  v-if="cred.status === 'approved'"
                  size="xs"
                  :variant="cred.is_active ? 'outline' : 'default'"
                  @click="handleToggle(cred)"
                >
                  {{ cred.is_active ? t('apiCredentials.actions.disable') : t('apiCredentials.actions.enable') }}
                </Button>
                <Button size="xs" variant="destructive" @click="handleDelete(cred.id)">
                  {{ t('apiCredentials.actions.delete') }}
                </Button>
              </div>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>

    <!-- Pagination -->
    <div v-if="pagination.total_page > 1" class="flex items-center justify-between gap-2 text-sm">
      <span class="text-muted-foreground">
        {{ pagination.page }} / {{ pagination.total_page }} ({{ pagination.total }})
      </span>
      <div class="flex items-center gap-1">
        <Button size="sm" variant="outline" :disabled="pagination.page <= 1" @click="goPage(pagination.page - 1)">
          &laquo;
        </Button>
        <Button size="sm" variant="outline" :disabled="pagination.page >= pagination.total_page" @click="goPage(pagination.page + 1)">
          &raquo;
        </Button>
        <Input
          v-model="jumpPage"
          class="w-16 h-8 text-xs"
          placeholder="Go"
          @keyup.enter="handleJump"
        />
      </div>
    </div>

    <!-- Detail Dialog -->
    <Dialog v-model:open="showDetail">
      <DialogScrollContent class="max-w-lg">
        <DialogHeader>
          <DialogTitle>{{ t('apiCredentials.detail.title') }}</DialogTitle>
        </DialogHeader>
        <div v-if="detailCred" class="space-y-3 text-sm">
          <div class="grid grid-cols-[120px_1fr] gap-y-2">
            <span class="text-muted-foreground">{{ t('apiCredentials.columns.id') }}</span>
            <span>{{ detailCred.id }}</span>
            <span class="text-muted-foreground">{{ t('apiCredentials.columns.userId') }}</span>
            <span>{{ detailCred.user_id }}</span>
            <span class="text-muted-foreground">{{ t('apiCredentials.columns.apiKey') }}</span>
            <span class="font-mono text-xs break-all">{{ detailCred.api_key || '-' }}</span>
            <span class="text-muted-foreground">{{ t('apiCredentials.columns.status') }}</span>
            <Badge :variant="statusVariant(detailCred.status)" class="w-fit">
              {{ t(`apiCredentials.status.${detailCred.status}`, detailCred.status) }}
            </Badge>
            <span class="text-muted-foreground">{{ t('apiCredentials.columns.isActive') }}</span>
            <span>{{ detailCred.is_active ? 'ON' : 'OFF' }}</span>
            <template v-if="detailCred.reject_reason">
              <span class="text-muted-foreground">{{ t('apiCredentials.detail.rejectReason') }}</span>
              <span class="text-destructive">{{ detailCred.reject_reason }}</span>
            </template>
            <span class="text-muted-foreground">{{ t('apiCredentials.columns.approvedAt') }}</span>
            <span>{{ formatDate(detailCred.approved_at) }}</span>
            <span class="text-muted-foreground">{{ t('apiCredentials.columns.lastUsedAt') }}</span>
            <span>{{ formatDate(detailCred.last_used_at) }}</span>
            <span class="text-muted-foreground">{{ t('apiCredentials.columns.createdAt') }}</span>
            <span>{{ formatDate(detailCred.created_at) }}</span>
          </div>
        </div>
      </DialogScrollContent>
    </Dialog>

    <!-- Reject Dialog -->
    <Dialog v-model:open="showReject">
      <DialogScrollContent class="max-w-md">
        <DialogHeader>
          <DialogTitle>{{ t('apiCredentials.reject.title') }}</DialogTitle>
        </DialogHeader>
        <div class="space-y-4">
          <div class="space-y-2">
            <Label>{{ t('apiCredentials.reject.reasonLabel') }}</Label>
            <Input v-model="rejectReason" :placeholder="t('apiCredentials.reject.reasonPlaceholder')" />
          </div>
          <div class="flex justify-end gap-2">
            <Button variant="outline" @click="showReject = false">Cancel</Button>
            <Button variant="destructive" :disabled="!rejectReason.trim()" @click="handleReject">
              {{ t('apiCredentials.actions.reject') }}
            </Button>
          </div>
        </div>
      </DialogScrollContent>
    </Dialog>

    <!-- Secret Dialog (after approve) -->
    <Dialog v-model:open="showSecret">
      <DialogScrollContent class="max-w-md">
        <DialogHeader>
          <DialogTitle>{{ t('apiCredentials.approve.secretTitle') }}</DialogTitle>
        </DialogHeader>
        <div class="space-y-3">
          <p class="text-sm text-amber-600 dark:text-amber-400 font-medium">
            {{ t('apiCredentials.approve.secretWarning') }}
          </p>
          <div class="space-y-2 text-sm">
            <div>
              <span class="text-muted-foreground">API Key: </span>
              <code class="bg-muted px-2 py-1 rounded text-xs break-all">{{ secretData.api_key }}</code>
            </div>
            <div>
              <span class="text-muted-foreground">API Secret: </span>
              <code class="bg-muted px-2 py-1 rounded text-xs break-all">{{ secretData.api_secret }}</code>
            </div>
          </div>
          <div class="flex justify-end">
            <Button @click="showSecret = false">OK</Button>
          </div>
        </div>
      </DialogScrollContent>
    </Dialog>
  </div>
</template>
