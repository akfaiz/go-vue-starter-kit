<script lang="ts" setup>
import type { ChangePasswordRequest } from '@/services/profile'
import { changePassword } from '@/services/profile'
import type { FieldErrors } from '@/utils/errors'
import { AppError } from '@/utils/errors'

const isPasswordVisible = ref(false)
const loading = ref(false)
const validationErrors = ref<FieldErrors>({})

const form = reactive<ChangePasswordRequest>({
  current_password: '',
  new_password: '',
  new_password_confirmation: '',
})

const handleSubmit = async () => {
  try {
    loading.value = true
    validationErrors.value = {}
    await changePassword(form)
    loading.value = false

    // Optionally, you can add a success message or redirect the user
    form.current_password = ''
    form.new_password = ''
    form.new_password_confirmation = ''
  }
  catch (e) {
    loading.value = false
    if (e instanceof AppError && e.isValidation)
      validationErrors.value = e.fieldErrors || {}
    else
      throw e
  }
}

const handleReset = () => {
  form.current_password = ''
  form.new_password = ''
  form.new_password_confirmation = ''
  validationErrors.value = {}
}
</script>

<template>
  <VRow>
    <!-- SECTION: Change Password -->
    <VCol cols="12">
      <VCard
        title="Update Password"
        subtitle="Ensure your account is using a long, random password to stay secure."
      >
        <VForm @submit.prevent="handleSubmit">
          <VCardText>
            <!-- 👉 Current Password -->
            <VRow>
              <VCol
                cols="12"
                md="6"
              >
                <!-- 👉 current password -->
                <VTextField
                  v-model="form.current_password"
                  :type="isPasswordVisible ? 'text' : 'password'"
                  :append-inner-icon="isPasswordVisible ? 'bx-hide' : 'bx-show'"
                  label="Current Password"
                  placeholder="············"
                  :error-messages="validationErrors.current_password"
                  @click:append-inner="isPasswordVisible = !isPasswordVisible"
                />
              </VCol>
            </VRow>

            <!-- 👉 New Password -->
            <VRow>
              <VCol
                cols="12"
                md="6"
              >
                <!-- 👉 new password -->
                <VTextField
                  v-model="form.new_password"
                  :type="isPasswordVisible ? 'text' : 'password'"
                  :append-inner-icon="isPasswordVisible ? 'bx-hide' : 'bx-show'"
                  label="New Password"
                  autocomplete="on"
                  placeholder="············"
                  :error-messages="validationErrors.new_password"
                  @click:append-inner="isPasswordVisible = !isPasswordVisible"
                />
              </VCol>
            </VRow>
            <VRow>
              <VCol
                cols="12"
                md="6"
              >
                <!-- 👉 confirm password -->
                <VTextField
                  v-model="form.new_password_confirmation"
                  :type="isPasswordVisible ? 'text' : 'password'"
                  :append-inner-icon="isPasswordVisible ? 'bx-hide' : 'bx-show'"
                  label="Confirm New Password"
                  placeholder="············"
                  :error-messages="validationErrors.new_password_confirmation"
                  @click:append-inner="isPasswordVisible = !isPasswordVisible"
                />
              </VCol>
            </VRow>
          </VCardText>

          <!-- 👉 Action Buttons -->
          <VCardText class="d-flex flex-wrap gap-4">
            <VBtn
              type="submit"
              :loading="loading"
            >
              Save
            </VBtn>

            <VBtn
              color="secondary"
              variant="tonal"
              @click="handleReset"
            >
              Reset
            </VBtn>
          </VCardText>
        </VForm>
      </VCard>
    </VCol>
    <!-- !SECTION -->
  </VRow>
</template>
