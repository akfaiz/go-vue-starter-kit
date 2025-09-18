<script setup lang="ts">
import type { ResetPasswordRequest } from '@/services/auth'
import { resetPassword, validatePasswordResetToken } from '@/services/auth'
import type { FieldErrors } from '@/utils/errors'
import { AppError } from '@/utils/errors'
import logo from '@images/logo.svg?raw'

const router = useRouter()

const form = reactive<ResetPasswordRequest>({
  token: '',
  email: '',
  password: '',
  password_confirmation: '',
})

const alert = ref({
  type: 'success' as 'success' | 'error',
  message: 'Lorem ipsum dolor sit amet consectetur adipisicing elit.',
  show: false,
})

const loading = ref(false)
const isPasswordVisible = ref(false)
const validationErrors = ref<FieldErrors>({})
const tokenValidated = ref(false)

onMounted(async () => {
  const url = new URL(window.location.href)
  const email = url.searchParams.get('email') || ''
  const token = url.searchParams.get('token') || ''
  if (!email || !token) {
    alert.value = {
      type: 'error',
      message: 'Invalid password reset link. Please request a new link.',
      show: true,
    }

    return
  }

  form.email = email
  form.token = token

  try {
    await validatePasswordResetToken(token, email)
    tokenValidated.value = true
  }
  catch (e) {
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

const handleSubmit = async () => {
  try {
    alert.value.show = false
    validationErrors.value = {}
    loading.value = true

    const msg = await resetPassword(form)

    alert.value = {
      type: 'success',
      message: msg.message,
      show: true,
    }
    loading.value = false
    router.replace({ path: '/login' })
  }
  catch (e) {
    loading.value = false
    if (e instanceof AppError) {
      if (e.isValidation) {
        validationErrors.value = e.fieldErrors || {}
      }
      else {
        alert.value = {
          type: 'error',
          message: e.message,
          show: true,
        }
      }
    }
    else {
      throw e
    }
  }
}
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

        <VCardText>
          <h4 class="text-h5 mb-1">
            Reset password
          </h4>
          <p class="mb-0">
            Please enter your new password below
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
          <VForm @submit.prevent="handleSubmit">
            <VRow>
              <!-- email -->
              <VCol cols="12">
                <VTextField
                  v-model="form.email"
                  autofocus
                  label="Email"
                  type="email"
                  disabled
                />
              </VCol>

              <VCol cols="12">
                <VTextField
                  v-model="form.password"
                  label="Password"
                  autocomplete="password"
                  placeholder="············"
                  :type="isPasswordVisible ? 'text' : 'password'"
                  :append-inner-icon="isPasswordVisible ? 'bx-hide' : 'bx-show'"
                  :error-messages="validationErrors.password"
                  @click:append-inner="isPasswordVisible = !isPasswordVisible"
                />
              </VCol>

              <VCol cols="12">
                <VTextField
                  v-model="form.password_confirmation"
                  label="Confirm Password"
                  autocomplete="password"
                  placeholder="············"
                  :type="isPasswordVisible ? 'text' : 'password'"
                  :append-inner-icon="isPasswordVisible ? 'bx-hide' : 'bx-show'"
                  :error-messages="validationErrors.password_confirmation"
                  @click:append-inner="isPasswordVisible = !isPasswordVisible"
                />
              </VCol>

              <VCol cols="12">
                <VBtn
                  block
                  type="submit"
                  :loading="loading"
                  :disabled="!tokenValidated"
                >
                  Reset Password
                </VBtn>
              </VCol>
            </VRow>
          </VForm>
        </VCardText>
      </VCard>
    </div>
  </div>
</template>

<style lang="scss">
@use "@core/scss/template/pages/page-auth";
</style>
