<script lang="ts" setup>
import type { UpdateProfileRequest } from '@/services/profile'
import { deleteAccount, updateProfile } from '@/services/profile'
import { useAuthStore } from '@/stores/auth'
import type { FieldErrors } from '@/utils/errors'
import { AppError } from '@/utils/errors'

const auth = useAuthStore()
const router = useRouter()
const { user } = storeToRefs(auth)

const updateProfileForm = reactive<UpdateProfileRequest>({
  name: '',
  email: '',
})

watch(
  user,
  u => {
    updateProfileForm.name = u?.name ?? ''
    updateProfileForm.email = u?.email ?? ''
  },
  { immediate: true },
)

const deleteForm = reactive({
  password: '',
})

const isAccountDeleted = ref(false)
const dialogDeleteAccount = ref(false)
const loading = ref(false)
const loadingDelete = ref(false)
const isPasswordVisible = ref(false)
const validationErrors = ref<FieldErrors>({})

const resetForm = () => {
  updateProfileForm.name = user.value?.name || ''
  updateProfileForm.email = user.value?.email || ''
  validationErrors.value = {}
}

const handleSubmit = async () => {
  try {
    loading.value = true
    validationErrors.value = {}
    await updateProfile(updateProfileForm)
    loading.value = false
    auth.fetchMe(true)
  }
  catch (e) {
    loading.value = false
    if (e instanceof AppError && e.isValidation)
      validationErrors.value = e.fieldErrors || {}
    else
      throw e
  }
}

const handleDeleteAccount = async () => {
  try {
    loadingDelete.value = true
    await deleteAccount(deleteForm.password)
    auth.logout()
    router.replace({ path: '/login' })
  }
  catch (e) {
    loadingDelete.value = false
    if (e instanceof AppError && e.isValidation)
      validationErrors.value = e.fieldErrors || {}
    else
      throw e
  }
}
</script>

<template>
  <VRow>
    <VCol cols="12">
      <VCard
        title="Profile Information"
        subtitle="Update your account's profile information and email address."
      >
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
                  v-model="updateProfileForm.name"
                  placeholder="John Doe"
                  label="Full Name"
                  :error-messages="validationErrors.name"
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
                  v-model="updateProfileForm.email"
                  label="E-mail"
                  placeholder="johndoe@gmail.com"
                  type="email"
                  :error-messages="validationErrors.email"
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
                  Save
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
      <!-- 👉 Delete Account -->
      <VCard
        title="Delete Account"
        subtitle="Delete your account and all of its resources"
      >
        <VCardText>
          <div>
            <VCheckbox
              v-model="isAccountDeleted"
              label="I confirm my account deletion"
            />
          </div>

          <VBtn
            :disabled="!isAccountDeleted"
            color="error"
            class="mt-3"
            @click="dialogDeleteAccount = true"
          >
            Delete Account
          </VBtn>
        </VCardText>
      </VCard>
      <VDialog
        v-model="dialogDeleteAccount"
        max-width="500"
      >
        <VCard>
          <VCardTitle class="text-h5">
            Are you sure you want to delete your account?
          </VCardTitle>
          <VCardText>
            Once your account is deleted, all of its resources and data will be permanently deleted. Please
            enter your password to confirm you would like to permanently delete your account.
            <VTextField
              v-model="deleteForm.password"
              :type="isPasswordVisible ? 'text' : 'password'"
              :append-inner-icon="isPasswordVisible ? 'bx-hide' : 'bx-show'"
              label="Password"
              autocomplete="current-password"
              placeholder="············"
              :error-messages="validationErrors.password"
              class="mt-4"
              @click:append-inner="isPasswordVisible = !isPasswordVisible"
            />
          </VCardText>
          <VCardActions>
            <VSpacer />
            <VBtn
              variant="text"
              @click="dialogDeleteAccount = false"
            >
              Cancel
            </VBtn>
            <VBtn
              color="error"
              :loading="loadingDelete"
              @click="handleDeleteAccount"
            >
              Delete Account
            </VBtn>
          </VCardActions>
        </VCard>
      </VDialog>
    </VCol>
  </VRow>
</template>
