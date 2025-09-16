// src/services/auth.ts
/* eslint-disable camelcase */

import type { ApiEnvelope, ApiMessage } from './response'
import { clearAuthTokens, setAuthTokens } from '@/utils/token'

/* ------------------------------- Types ----------------------------------- */

export interface User {
  id: number | string
  name: string
  email: string
  created_at?: string
  updated_at?: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  name: string
  email: string
  password: string
  password_confirmation: string
}

export interface ResetPasswordRequest {
  token: string
  email: string
  password: string
  password_confirmation: string
}

export interface TokenResponse {
  access_token: string
  refresh_token: string | null
}

/* ----------------------------- Services ---------------------------------- */

/** /login -> { status, message, data: { access_token, refresh_token } } */
export async function login(payload: LoginRequest): Promise<TokenResponse> {
  const res = await $api.post<ApiEnvelope<TokenResponse>>('/v1/auth/login', payload)

  const { access_token, refresh_token } = res.data.data ?? {}

  if (!access_token)
    throw new Error('Login failed: access_token missing in response')

  setAuthTokens(access_token, refresh_token ?? '')

  return { access_token, refresh_token }
}

/** /register -> { status, message } (no tokens) */
export async function register(payload: RegisterRequest): Promise<TokenResponse> {
  const { data } = await $api.post<ApiEnvelope<TokenResponse>>('/v1/auth/register', payload)

  const { access_token, refresh_token } = data.data ?? {}

  if (!access_token)
    throw new Error('Registration failed: access_token missing in response')

  setAuthTokens(access_token, refresh_token ?? '')

  return { access_token, refresh_token }
}

/** /me -> { status, data: User } */
export async function me(): Promise<User> {
  const res = await $api.get<ApiEnvelope<User>>('/v1/profile')
  const user = res.data.data
  if (!user)
    throw new Error('Invalid /me response: data is missing')

  return user
}

export async function logout(): Promise<void> {
  // Optional: await $api.post("/v1/auth/logout").catch(() => {});
  clearAuthTokens()
}

export async function sendPasswordResetEmail(email: string): Promise<ApiMessage> {
  const { data } = await $api.post<ApiMessage>('/v1/auth/forgot-password', { email })

  return data
}

export async function validatePasswordResetToken(token: string, email: string): Promise<ApiMessage> {
  const { data } = await $api.post<ApiMessage>('/v1/auth/validate-reset-password', { token, email })

  return data
}

export async function resetPassword(payload: ResetPasswordRequest): Promise<ApiMessage> {
  const { data } = await $api.post<ApiMessage>('/v1/auth/reset-password', payload)

  return data
}

export async function sendEmailVerification(): Promise<ApiMessage> {
  const { data } = await $api.post<ApiMessage>('/v1/auth/email/send-verification')

  return data
}

export async function verifyEmail(token: string, user_id: number): Promise<ApiMessage> {
  const { data } = await $api.post<ApiMessage>('/v1/auth/email/verify', { token, user_id })

  return data
}
