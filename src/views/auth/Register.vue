<template>
  <div class="relative flex min-h-screen items-center justify-center theme-page theme-auth-page px-4 py-16 sm:px-6">
    <div class="relative z-10 w-full max-w-lg">
      <div class="mb-4 flex items-center justify-between px-1">
        <router-link
          to="/"
          class="theme-nav-link rounded-full gap-1"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          {{ t('auth.login.backHome') }}
        </router-link>
        <span class="rounded-full border theme-pill-neutral px-3 py-1 text-xs font-semibold">
          {{ t('auth.register.title') }}
        </span>
      </div>

      <div class="theme-auth-card">
        <div class="mb-8 text-center">
          <p class="text-xs font-semibold uppercase tracking-[0.22em] theme-text-accent">{{ brandSiteName }}</p>
          <h1 class="mt-3 text-3xl font-black theme-text-primary">{{ t('auth.register.title') }}</h1>
          <p class="mt-2 text-sm theme-text-muted">{{ t('auth.register.subtitle') }}</p>
        </div>

        <form class="theme-auth-form" @submit.prevent="handleRegister">
          <div>
            <label class="mb-2 flex items-center gap-1.5 text-xs font-semibold uppercase tracking-[0.18em] theme-text-muted">
              <svg class="w-3.5 h-3.5 opacity-60" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
              </svg>
              {{ t('auth.register.emailLabel') }}
            </label>
            <input
              v-model="email"
              type="email"
              required
              class="w-full form-input-lg"
              :placeholder="t('auth.register.emailPlaceholder')"
            />
          </div>

          <div>
            <label class="mb-2 flex items-center gap-1.5 text-xs font-semibold uppercase tracking-[0.18em] theme-text-muted">
              <svg class="w-3.5 h-3.5 opacity-60" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
              </svg>
              {{ t('auth.register.passwordLabel') }}
            </label>
            <input
              v-model="password"
              type="password"
              required
              class="w-full form-input-lg"
              :placeholder="t('auth.register.passwordPlaceholder')"
            />
          </div>

          <div v-if="sendCodeCaptchaEnabled">
            <label class="mb-2 flex items-center gap-1.5 text-xs font-semibold uppercase tracking-[0.18em] theme-text-muted">
              <svg class="w-3.5 h-3.5 opacity-60" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
              </svg>
              {{ t('auth.common.captchaLabel') }}
            </label>
            <ImageCaptcha
              v-if="captchaProvider === 'image'"
              ref="imageCaptchaRef"
              v-model="captchaPayload"
              :disabled="sending || countdown > 0"
              @config-stale="handleCaptchaConfigStale"
            />
            <TurnstileCaptcha
              v-else-if="captchaProvider === 'turnstile'"
              ref="turnstileRef"
              v-model="turnstileToken"
              :site-key="turnstileSiteKey"
            />
          </div>

          <div>
            <label class="mb-2 flex items-center gap-1.5 text-xs font-semibold uppercase tracking-[0.18em] theme-text-muted">
              <svg class="w-3.5 h-3.5 opacity-60" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
              </svg>
              {{ t('auth.register.codeLabel') }}
            </label>
            <div class="flex flex-col gap-2 sm:flex-row">
              <input
                v-model="code"
                type="text"
                required
                class="min-w-0 flex-1 form-input-lg"
                :placeholder="t('auth.register.codePlaceholder')"
              />
              <button
                type="button"
                @click="handleSendCode"
                :disabled="sending || countdown > 0"
                class="whitespace-nowrap rounded-xl border theme-btn-secondary px-4 py-3 text-xs font-semibold disabled:cursor-not-allowed disabled:opacity-50"
              >
                {{ countdown > 0 ? t('auth.common.countdown', { seconds: countdown }) : t('auth.common.sendCode') }}
              </button>
            </div>
          </div>

          <label class="flex items-start gap-3 rounded-xl border theme-surface-soft px-4 py-3 text-sm theme-text-secondary transition-colors">
            <input
              v-model="agreed"
              type="checkbox"
              class="mt-0.5 h-4 w-4 rounded border-gray-300 theme-accent-checkbox dark:border-white/20 dark:bg-black/20"
            />
            <span class="leading-6">
              {{ t('auth.register.agreementPrefix') }}
              <router-link to="/privacy" target="_blank" rel="noopener noreferrer" class="font-semibold theme-link">
                {{ t('footer.privacy') }}
              </router-link>
              {{ t('auth.register.agreementAnd') }}
              <router-link to="/terms" target="_blank" rel="noopener noreferrer" class="font-semibold theme-link">
                {{ t('footer.terms') }}
              </router-link>
            </span>
          </label>

          <div
            v-if="error"
            class="rounded-xl border theme-alert-danger px-4 py-3 text-center text-sm"
          >
            {{ error }}
          </div>

          <button
            type="submit"
            :disabled="userAuthStore.loading || !agreed"
            class="inline-flex w-full items-center justify-center rounded-xl theme-btn-primary px-4 py-3 text-sm font-bold transition-colors disabled:cursor-not-allowed disabled:opacity-50"
          >
            <svg v-if="!userAuthStore.loading" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
            </svg>
            {{ userAuthStore.loading ? t('auth.register.creating') : t('auth.register.create') }}
          </button>
        </form>
      </div>

      <div class="mt-4 text-center">
        <router-link
          to="/auth/login"
          class="theme-link-muted text-sm"
        >
          {{ t('auth.register.hasAccount') }}
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserAuthStore } from '../../stores/userAuth'
import { useI18n } from 'vue-i18n'
import { debounceAsync } from '../../utils/debounce'
import { useAppStore } from '../../stores/app'
import type { CaptchaPayload } from '../../api'
import ImageCaptcha from '../../components/captcha/ImageCaptcha.vue'
import TurnstileCaptcha from '../../components/captcha/TurnstileCaptcha.vue'

const router = useRouter()
const userAuthStore = useUserAuthStore()
const appStore = useAppStore()
const { t } = useI18n()

const brandSiteName = computed(() => {
  const siteName = String(appStore.config?.brand?.site_name || '').trim()
  return siteName !== '' ? siteName : 'Dujiao-Next'
})

const email = ref('')
const password = ref('')
const code = ref('')
const agreed = ref(false)
const error = ref('')
const sending = ref(false)
const countdown = ref(0)
const captchaPayload = ref<CaptchaPayload>({})
const turnstileToken = ref('')
const imageCaptchaRef = ref<InstanceType<typeof ImageCaptcha> | null>(null)
const turnstileRef = ref<InstanceType<typeof TurnstileCaptcha> | null>(null)
let timer: number | undefined

const captchaConfig = computed(() => appStore.config?.captcha || null)
const captchaProvider = computed(() => String(captchaConfig.value?.provider || 'none'))
const sendCodeCaptchaEnabled = computed(() => !!captchaConfig.value?.scenes?.register_send_code && captchaProvider.value !== 'none')
const turnstileSiteKey = computed(() => String(captchaConfig.value?.turnstile?.site_key || ''))

const startCountdown = () => {
  countdown.value = 60
  timer = window.setInterval(() => {
    countdown.value -= 1
    if (countdown.value <= 0 && timer) {
      clearInterval(timer)
      timer = undefined
    }
  }, 1000)
}

const getCaptchaPayload = (): CaptchaPayload | undefined => {
  if (!sendCodeCaptchaEnabled.value) return undefined
  if (captchaProvider.value === 'image') {
    return {
      captcha_id: captchaPayload.value.captcha_id || '',
      captcha_code: captchaPayload.value.captcha_code || '',
    }
  }
  if (captchaProvider.value === 'turnstile') {
    return {
      turnstile_token: turnstileToken.value,
    }
  }
  return undefined
}

const handleCaptchaConfigStale = async () => {
  await appStore.loadConfig(true)
  captchaPayload.value = {}
  turnstileToken.value = ''
}

const performSendCode = async () => {
  error.value = ''
  if (!email.value) {
    error.value = t('auth.register.errors.emailRequired')
    return
  }
  if (countdown.value > 0) return

  if (sendCodeCaptchaEnabled.value && captchaProvider.value === 'image') {
    if (!captchaPayload.value.captcha_id || !captchaPayload.value.captcha_code) {
      error.value = t('auth.common.captchaRequired')
      return
    }
  }
  if (sendCodeCaptchaEnabled.value && captchaProvider.value === 'turnstile') {
    if (!turnstileToken.value) {
      error.value = t('auth.common.captchaRequired')
      return
    }
  }

  sending.value = true
  try {
    await userAuthStore.sendVerifyCode({
      email: email.value,
      purpose: 'register',
      captcha_payload: getCaptchaPayload(),
    })
    startCountdown()
  } catch (err: any) {
    error.value = err.message || t('auth.register.errors.sendCodeFailed')
    if (captchaProvider.value === 'image') {
      imageCaptchaRef.value?.refresh()
    }
    if (captchaProvider.value === 'turnstile') {
      turnstileRef.value?.reset()
      turnstileToken.value = ''
    }
  } finally {
    sending.value = false
  }
}

const performRegister = async () => {
  error.value = ''
  if (!email.value || !password.value || !code.value) return
  if (!agreed.value) {
    error.value = t('auth.register.errors.agreementRequired')
    return
  }
  try {
    await userAuthStore.register({
      email: email.value,
      password: password.value,
      code: code.value,
      agreement_accepted: agreed.value,
    })
    router.push('/me/orders')
  } catch (err: any) {
    error.value = err.message || t('auth.register.errors.registerFailed')
  }
}

const handleSendCode = debounceAsync(performSendCode, 200)
const handleRegister = debounceAsync(performRegister, 200)

onMounted(async () => {
  await appStore.loadConfig(true)
})
</script>
