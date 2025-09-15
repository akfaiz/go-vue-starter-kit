<script lang="ts" setup>
import { updateProfile } from '@/services/profile'
import { useAuthStore } from '@/stores/auth'
import { AppError } from '@/utils/errors'

const auth = useAuthStore()
const { user } = storeToRefs(auth)

const accountData = {
  name: user.value?.name ?? '',
  email: user.value?.email ?? '',
}

const accountDataLocal = ref(structuredClone(accountData))
const isAccountDeactivated = ref(false)
const loading = ref(false)

const resetForm = () => {
  accountDataLocal.value = structuredClone(accountData)
}

const validationErrors = ref<Record<string, string>>({})

const handleSubmit = async () => {
  try {
    loading.value = true
    validationErrors.value = {}
    await updateProfile({
      name: accountDataLocal.value.name,
      email: accountDataLocal.value.email,
    })
    loading.value = false
    auth.fetchMe(true)
  }
  catch (e) {
    loading.value = false
    if (e instanceof AppError && e.isValidation)
      validationErrors.value = e.fieldMap || {}
    else
      throw e
  }
}
</script>

<template>
  <VRow>
    <VCol cols="12">
      <VCard title="Account Details">
        <VCardText>
          <!-- 👉 Form -->
          <VForm
            class="mt-6"
            @submit.prevent="handleSubmit"
          >
            <VRow>
              <!-- 👉 Full Name -->
              <VCol
                md="6"
                cols="12"
              >
                <VTextField
                  v-model="accountDataLocal.name"
                  placeholder="John Doe"
                  label="Full Name"
                />
              </VCol>
            </VRow>

            <VRow>
              <!-- 👉 Email -->
              <VCol
                cols="12"
                md="6"
              >
                <VTextField
                  v-model="accountDataLocal.email"
                  label="E-mail"
                  placeholder="johndoe@gmail.com"
                  type="email"
                />
              </VCol>

              <!-- 👉 Form Actions -->
              <VCol
                cols="12"
                class="d-flex flex-wrap gap-4"
              >
                <VBtn
                  type="submit"
                  :loading="loading"
                >
                  Save changes
                </VBtn>

                <VBtn
                  color="secondary"
                  variant="tonal"
                  type="reset"
                  @click.prevent="resetForm"
                >
                  Reset
                </VBtn>
              </VCol>
            </VRow>
          </VForm>
        </VCardText>
      </VCard>
    </VCol>

    <VCol cols="12">
      <!-- 👉 Deactivate Account -->
      <VCard title="Deactivate Account">
        <VCardText>
          <div>
            <VCheckbox
              v-model="isAccountDeactivated"
              label="I confirm my account deactivation"
            />
          </div>

          <VBtn
            :disabled="!isAccountDeactivated"
            color="error"
            class="mt-3"
          >
            Deactivate Account
          </VBtn>
        </VCardText>
      </VCard>
    </VCol>
  </VRow>
</template>
