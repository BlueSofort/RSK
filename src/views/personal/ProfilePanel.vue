<template>
  <div class="theme-personal-card">
    <div class="mb-6 flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
      <div>
        <h2 class="text-xl font-bold theme-text-primary">{{ t('personalCenter.profile.title') }}</h2>
        <p class="mt-1 text-sm theme-text-muted">{{ t('personalCenter.profile.subtitle') }}</p>
      </div>
      <span class="theme-badge theme-badge-accent px-3 py-1 text-xs font-semibold">
        {{ t('personalCenter.tabs.profile') }}
      </span>
    </div>

    <form class="space-y-6" @submit.prevent="handleSaveProfile">
      <!-- Avatar Section -->
      <div class="flex items-center gap-5 pb-5 border-b border-gray-200/70 dark:border-white/10">
        <div class="relative">
          <img
            :src="avatarPreview || getAvatarUrl()"
            class="w-16 h-16 rounded-full object-cover border-2 border-gray-200 dark:border-white/10"
            @error="($event.target as HTMLImageElement).src = 'data:image/svg+xml,' + encodeURIComponent('<svg xmlns=%22http://www.w3.org/2000/svg%22 viewBox=%220 0 40 40%22><rect width=%2240%22 height=%2240%22 fill=%22%23e5e7eb%22/><text x=%2220%22 y=%2226%22 text-anchor=%22middle%22 font-size=%2218%22 fill=%22%239ca3af%22>?</text></svg>')"
          />
          <div v-if="avatarUploading" class="absolute inset-0 rounded-full bg-black/40 flex items-center justify-center">
            <div class="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
          </div>
        </div>
        <div class="space-y-2">
          <label class="inline-flex items-center gap-2 px-4 py-2 rounded-xl border cursor-pointer text-sm font-medium theme-btn-secondary hover:scale-[1.02] transition-transform">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
            {{ avatarUploading ? t('personalCenter.profile.uploading') : t('personalCenter.profile.changeAvatar') }}
            <input type="file" accept="image/*" class="hidden" @change="handleAvatarChange" :disabled="avatarUploading" />
          </label>
          <p class="text-xs theme-text-muted">{{ t('personalCenter.profile.avatarHint') }}</p>
        </div>
      </div>

      <div class="grid grid-cols-1 gap-5 md:grid-cols-2">
        <div class="md:col-span-2">
          <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-200">{{ t('personalCenter.profile.emailLabel') }}</label>
          <input
            :value="userProfileStore.profile?.email || ''"
            disabled
            class="w-full rounded-xl border border-gray-200 bg-gray-100 px-4 py-3 text-gray-500 dark:border-white/10 dark:bg-white/5"
          />
        </div>

        <div>
          <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-200">{{ t('personalCenter.profile.nicknameLabel') }}</label>
          <input
            v-model="profileForm.nickname"
            :placeholder="t('personalCenter.profile.nicknamePlaceholder')"
            class="w-full form-input-lg"
          />
        </div>

        <div>
          <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-200">{{ t('personalCenter.profile.localeLabel') }}</label>
          <select
            v-model="profileForm.locale"
            class="w-full form-input-lg"
          >
            <option value="zh-CN">简体中文</option>
            <option value="zh-TW">繁體中文</option>
            <option value="en-US">English</option>
          </select>
        </div>
      </div>

      <div class="flex flex-col gap-3 border-t border-gray-200/70 pt-5 dark:border-white/10 sm:flex-row sm:items-center sm:justify-between">
        <p class="text-xs theme-text-muted">{{ t('personalCenter.profile.subtitle') }}</p>
        <button
          type="submit"
          :disabled="userProfileStore.savingProfile"
          class="inline-flex items-center justify-center rounded-xl theme-btn-primary px-6 py-3 text-sm font-bold transition-colors disabled:cursor-not-allowed disabled:opacity-60"
        >
          {{ userProfileStore.savingProfile ? t('personalCenter.profile.saving') : t('personalCenter.profile.save') }}
        </button>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useUserProfileStore } from '../../stores/userProfile'
import { userProfileAPI } from '../../api'
import { toast } from '../../composables/useToast'

const { t } = useI18n()
const userProfileStore = useUserProfileStore()

const profileForm = reactive({
  nickname: '',
  locale: 'zh-CN',
})

const avatarPreview = ref('')
const avatarUploading = ref(false)

const getAvatarUrl = () => {
  const avatar = userProfileStore.profile?.avatar
  if (!avatar) {
    return 'data:image/svg+xml,' + encodeURIComponent('<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 40 40"><rect width="40" height="40" fill="#e5e7eb"/><text x="20" y="26" text-anchor="middle" font-size="18" fill="#9ca3af">?</text></svg>')
  }
  if (avatar.startsWith('http')) return avatar
  const apiBase = import.meta.env.VITE_API_BASE_URL || ''
  return `${apiBase}${avatar.startsWith('/') ? '' : '/'}${avatar}`
}

const handleAvatarChange = async (e: Event) => {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return

  // Preview
  const reader = new FileReader()
  reader.onload = () => {
    avatarPreview.value = reader.result as string
  }
  reader.readAsDataURL(file)

  // Upload
  avatarUploading.value = true
  const formData = new FormData()
  formData.append('avatar', file)
  try {
    await userProfileAPI.uploadAvatar(formData)
    // Refresh profile to get updated avatar
    await userProfileStore.loadProfile()
    toast.success(t('personalCenter.profile.avatarSuccess'))
  } catch (err: any) {
    toast.error(err?.message || t('personalCenter.profile.avatarFailed'))
    avatarPreview.value = ''
  } finally {
    avatarUploading.value = false
    input.value = ''
  }
}

const handleSaveProfile = async () => {
  const payload = {
    nickname: profileForm.nickname.trim(),
    locale: profileForm.locale,
  }
  const ok = await userProfileStore.saveProfile(payload)
  if (!ok) {
    toast.error(userProfileStore.profileError || t('personalCenter.common.saveFailed'))
    return
  }
  toast.success(t('personalCenter.profile.saveSuccess'))
}

watch(
  () => userProfileStore.profile,
  (profile) => {
    if (!profile) return
    profileForm.nickname = profile.nickname || ''
    profileForm.locale = profile.locale || 'zh-CN'
    avatarPreview.value = '' // reset preview on profile load
  },
  { immediate: true }
)
</script>
