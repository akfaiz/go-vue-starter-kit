import type { ApiMessage } from './response'
import { clearAuthTokens } from '@/utils/token'

export interface ChangePasswordRequest {
  current_password: string
  new_password: string
  new_password_confirmation: string
}

export interface UpdateProfileRequest {
  name: string
  email: string
}

export async function updateProfile(payload: UpdateProfileRequest): Promise<ApiMessage> {
  const { data } = await $api.put<ApiMessage>('/v1/profile', payload)

  return data
}

export async function changePassword(payload: ChangePasswordRequest): Promise<ApiMessage> {
  const { data } = await $api.put<ApiMessage>('/v1/profile/password', payload)

  return data
}

export async function deleteAccount(password: string): Promise<ApiMessage> {
  const { data } = await $api.delete<ApiMessage>('/v1/profile', {
    data: { password },
  })

  // On successful account deletion, clear tokens to force re-login
  clearAuthTokens()

  return data
}
