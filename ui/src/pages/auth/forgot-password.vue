<script setup lang="ts">
import { sendPasswordResetEmail } from '@/services/auth'
import type { FieldErrors } from '@/utils/errors'
import { AppError } from '@/utils/errors'
import logo from '@images/logo.svg?raw'

const form = ref({
  email: '',
})

const alert = ref({
  type: 'success' as 'success' | 'error',
  message: '',
  show: false,
})

const loading = ref(false)
const validationErrors = ref<FieldErrors>({})

const handleSubmit = async () => {
  try {
    alert.value.show = false
    validationErrors.value = {}
    loading.value = true

    const msg = await sendPasswordResetEmail(form.value.email)

    alert.value = {
      type: 'success',
      message: msg.message,
      show: true,
    }
    loading.value = false
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

        <VCardText class="text-center">
          <h4 class="text-h5 mb-1">
            Forgot password
          </h4>
          <p class="mb-0">
            Enter your email to receive a password reset link
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
                  placeholder="johndoe@email.com"
                  :error-messages="validationErrors.email"
                />
              </VCol>

              <!-- password -->

              <VCol cols="12">
                <VBtn
                  block
                  type="submit"
                  :loading="loading"
                >
                  Email password reset link
                </VBtn>
              </VCol>

              <VCol
                cols="12"
                class="text-body-1 text-center"
              >
                <span class="d-inline-block">
                  Or, return to
                </span>
                <RouterLink
                  class="text-primary ms-1 d-inline-block text-body-1"
                  to="/login"
                >
                  Log in
                </RouterLink>
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
