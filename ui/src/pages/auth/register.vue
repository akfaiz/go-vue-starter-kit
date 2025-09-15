<script setup lang="ts">
import { register } from '@/services/auth'
import { AppError } from '@/utils/errors'
import logo from '@images/logo.svg?raw'

const form = ref({
  name: '',
  email: '',
  password: '',
  password_confirmation: '',
})

const isPasswordVisible = ref(false)
const loading = ref(false)
const validationErrors = ref<Record<string, string>>({})

const handleSubmit = async () => {
  try {
    loading.value = true
    validationErrors.value = {}
    await register({
      name: form.value.name,
      email: form.value.email,
      password: form.value.password,
      password_confirmation: form.value.password_confirmation,
    })
    loading.value = false

    window.location.href = '/login'
  }
  catch (e) {
    loading.value = false
    if (e instanceof AppError && e.isValidation)
      validationErrors.value = e.fieldMap || {}
  }
}
</script>

<template>
  <div class="auth-wrapper d-flex align-center justify-center pa-4">
    <div class="position-relative my-sm-16">
      <!-- 👉 Auth card -->
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
            Create an account
          </h4>
          <p class="mb-0">
            Enter your details below to create your account
          </p>
        </VCardText>

        <VCardText>
          <VForm @submit.prevent="handleSubmit">
            <VRow>
              <!-- Name -->
              <VCol cols="12">
                <VTextField
                  v-model="form.name"
                  autofocus
                  label="Name"
                  placeholder="John Doe"
                  :error-messages="validationErrors.name"
                />
              </VCol>
              <!-- email -->
              <VCol cols="12">
                <VTextField
                  v-model="form.email"
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
                  autocomplete="password"
                  placeholder="············"
                  :type="isPasswordVisible ? 'text' : 'password'"
                  :append-inner-icon="isPasswordVisible ? 'bx-hide' : 'bx-show'"
                  :error-messages="validationErrors.password"
                  @click:append-inner="isPasswordVisible = !isPasswordVisible"
                />
              </VCol>

              <!-- confirm password -->
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

                <VBtn
                  block
                  type="submit"
                  class="mt-6"
                  :loading="loading"
                >
                  Sign up
                </VBtn>
              </VCol>

              <!-- login instead -->
              <VCol
                cols="12"
                class="text-center text-base"
              >
                <span>Already have an account?</span>
                <RouterLink
                  class="text-primary ms-1"
                  to="/login"
                >
                  Sign in instead
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
