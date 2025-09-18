<script setup lang="ts">
import type { LoginRequest } from '@/services/auth'
import { useAuthStore } from '@/stores/auth'
import type { FieldErrors } from '@/utils/errors'
import { AppError } from '@/utils/errors'
import logo from '@images/logo.svg?raw'

const auth = useAuthStore()
const router = useRouter()

const form = reactive<LoginRequest>({
  email: '',
  password: '',
})

const isPasswordVisible = ref(false)
const validationErrors = ref<FieldErrors>({})
const loading = ref(false)

const handleSubmit = async () => {
  try {
    loading.value = true
    validationErrors.value = {}
    await auth.login(form)
    loading.value = false

    router.replace({ path: '/dashboard' })
  }
  catch (e) {
    loading.value = false
    if (e instanceof AppError && e.isValidation)
      validationErrors.value = e.fieldErrors || {}
    else
      throw e
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
            Log in to your account
          </h4>
          <p class="mb-0">
            Enter your email and password below to log in
          </p>
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
                <VTextField
                  v-model="form.password"
                  label="Password"
                  placeholder="············"
                  :type="isPasswordVisible ? 'text' : 'password'"
                  autocomplete="password"
                  :append-inner-icon="isPasswordVisible ? 'bx-hide' : 'bx-show'"
                  :error-messages="validationErrors.password"
                  @click:append-inner="isPasswordVisible = !isPasswordVisible"
                />

                <!-- remember me checkbox -->
                <div class="d-flex align-center justify-space-between flex-wrap my-6">
                  <VCheckbox label="Remember me" />
                  <RouterLink
                    class="text-primary ms-1 d-inline-block text-body-1"
                    to="/forgot-password"
                  >
                    Forgot Password?
                  </RouterLink>
                </div>

                <!-- login button -->
                <VBtn
                  block
                  type="submit"
                  :loading="loading"
                >
                  Login
                </VBtn>
              </VCol>

              <!-- create account -->
              <VCol
                cols="12"
                class="text-body-1 text-center"
              >
                <span class="d-inline-block">
                  New on our platform?
                </span>
                <RouterLink
                  class="text-primary ms-1 d-inline-block text-body-1"
                  to="/register"
                >
                  Create an account
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
