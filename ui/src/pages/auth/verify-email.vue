<script setup lang="ts">
import { sendEmailVerification, verifyEmail } from '@/services/auth'
import { useAuthStore } from '@/stores/auth'
import { AppError } from '@/utils/errors'
import logo from '@images/logo.svg?raw'

const auth = useAuthStore()
const router = useRouter()

const alert = ref({
  type: 'success' as 'success' | 'error',
  message: '',
  show: false,
})

const loading = ref(false)

const handleSendEmailVerification = async () => {
  try {
    alert.value.show = false
    loading.value = true

    await sendEmailVerification()

    alert.value = {
      type: 'success',
      message: 'A new verification link has been sent to the email address you provided during registration.',
      show: true,
    }
    loading.value = false
  }
  catch (e) {
    loading.value = false
    if (e instanceof AppError) {
      alert.value = {
        type: 'error',
        message: e.message,
        show: true,
      }
    }
    else {
      throw e
    }
  }
}

const handleLogout = async () => {
  await auth.logout()
  router.replace({ path: '/login' })
}

onMounted(async () => {
  if (auth.user?.email_verified_at)
    router.replace({ path: '/dashboard' })
  const url = new URL(window.location.href)
  const token = url.searchParams.get('token') || ''
  const userId = url.searchParams.get('user_id') || ''
  if (!token || !userId) {
    return
  }
  const userIDNum = Number(userId)
  if (isNaN(userIDNum)) {
    return
  }
  try {
    loading.value = true
    await verifyEmail(token, userIDNum)
    alert.value = {
      type: 'success',
      message: 'Your email has been successfully verified. You can now access the dashboard.',
      show: true,
    }
    // Refresh user data
    await auth.fetchMe(true)
    loading.value = false
    router.replace({ path: '/dashboard' })
  }
  catch (e) {
    loading.value = false
    if (e instanceof AppError) {
      alert.value = {
        type: 'error',
        message: e.message,
        show: true,
      }
    }
    else {
      throw e
    }
  }
})
</script>

<template>
  <div class="auth-wrapper d-flex align-center justify-center pa-4">
    <div class="position-relative my-sm-16">
      <!-- 👉 Auth Card -->
      <VCard
        class="auth-card"
        max-width="460"
        :class="$vuetify.display.smAndUp ? 'pa-6' : 'pa-0'"
      >
        <VCardItem class="justify-center">
          <RouterLink
            to="/"
            class="app-logo"
          >
            <!-- eslint-disable vue/no-v-html -->
            <div
              class="d-flex"
              v-html="logo"
            />
            <h1 class="app-logo-title">
              govue
            </h1>
          </RouterLink>
        </VCardItem>

        <VCardText class="text-center">
          <h4 class="text-h5 mb-1">
            Verify Email
          </h4>
          <p class="mb-0">
            Please verify your email address by clicking on the link we just emailed to you.
          </p>
        </VCardText>

        <VCardText
          v-if="alert.show"
          class="pt-0"
        >
          <VAlert
            :color="alert.type === 'success' ? 'success' : 'error'"
            variant="text"
            :text="alert.message"
          />
        </VCardText>

        <VCardText>
          <VRow>
              <VCol cols="12">
                <VBtn
                  block
                  :loading="loading"
                  @click="handleSendEmailVerification"  
                >
                  Resend Verification Email
                </VBtn>
              </VCol>

              <VCol
                cols="12"
                class="text-body-1 text-center"
                @click="handleLogout"
              >
                <VBtn variant="text">
                  Log out
                </VBtn>
              </VCol>
            </VRow>
        </VCardText>
      </VCard>
    </div>
  </div>
</template>

<style lang="scss">
@use "@core/scss/template/pages/page-auth";
</style>
