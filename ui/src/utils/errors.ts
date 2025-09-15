// src/utils/errors.ts
import type { AxiosError } from 'axios'

/** Field-level error from backend */
export interface FieldError {
  field: string
  message: string
}

/** RFC 7807-ish problem details with extra `errors` */
export interface ProblemDetails {
  type?: string // e.g. "about:blank"
  title?: string // e.g. "Validation failed"
  status?: number // e.g. 422
  detail?: string // e.g. "One or more fields are invalid..."
  instance?: string // e.g. request id / correlation id
  errors?: FieldError[]
}

/** Normalized app error you can safely throw/catch in UI code */
export class AppError extends Error {
  name = 'AppError'
  status?: number
  type?: string
  instance?: string

  /** Full backend payload if available */
  problem?: ProblemDetails

  /** Field -> list of messages (for form binding) */
  fieldMap: Record<string, string>
  isValidation: boolean
  isNetwork: boolean
  isTimeout: boolean

  constructor(init: {
    message: string
    status?: number
    type?: string
    instance?: string
    problem?: ProblemDetails
    fieldMap?: Record<string, string>
    isValidation?: boolean
    isNetwork?: boolean
    isTimeout?: boolean
    cause?: unknown
  }) {
    super(init.message, { cause: init.cause })
    this.status = init.status
    this.type = init.type
    this.instance = init.instance
    this.problem = init.problem
    this.fieldMap = init.fieldMap ?? {}
    this.isValidation = !!init.isValidation
    this.isNetwork = !!init.isNetwork
    this.isTimeout = !!init.isTimeout
  }
}

/* ------------------------------- Utilities ------------------------------- */

function isObject(x: unknown): x is Record<string, unknown> {
  return !!x && typeof x === 'object'
}

function isProblemDetails(x: unknown): x is ProblemDetails {
  return (
    isObject(x)
    && ('title' in x || 'status' in x || 'detail' in x || 'type' in x || 'errors' in x)
  )
}

function buildFieldMap(errors?: FieldError[]): Record<string, string> {
  const map: Record<string, string> = {}
  if (!errors?.length)
    return map
  for (const e of errors) {
    if (!e?.field)
      continue
    if (!map[e.field])
      map[e.field] = ''
    if (e.message)
      map[e.field] = e.message
  }

  return map
}

/**
 * Convert any thrown value (AxiosError, plain Error, unknown) into AppError.
 * Understands:
 * - Network/timeout errors
 * - RFC7807-ish payloads with `errors: [{field, message}]`
 * - Falls back to a generic message
 */
export function toAppError(err: unknown): AppError {
  // AxiosError branch
  const ax = err as AxiosError
  if (ax && isObject(ax) && 'isAxiosError' in ax) {
    const status = ax.response?.status
    const code = ax.code ?? ''
    const net = !ax.response // no response => likely network
    const timeout = code === 'ECONNABORTED'

    if (isProblemDetails(ax.response?.data)) {
      const p = ax.response!.data as ProblemDetails
      const isValidation = (status === 422) || !!p.errors?.length
      const fieldMap = buildFieldMap(p.errors)

      const message
        = p.title
        || p.detail
        || (isValidation ? 'Validation failed' : `Request failed with status ${status ?? 'unknown'}`)

      return new AppError({
        message,
        status,
        type: p.type,
        instance: p.instance,
        problem: p,
        fieldMap,
        isValidation,
        isNetwork: net,
        isTimeout: timeout,
        cause: err,
      })
    }

    // Non-problem JSON or HTML/text response
    if (status) {
      const message
        = (isObject(ax.response?.data) && String((ax.response!.data as any).message || ''))
        || ax.message
        || `Request failed with status ${status}`

      return new AppError({
        message,
        status,
        isNetwork: net,
        isTimeout: timeout,
        cause: err,
      })
    }

    // Pure network/timeout
    return new AppError({
      message: timeout ? 'Request timed out' : 'Network error. Please check your connection.',
      isNetwork: net,
      isTimeout: timeout,
      cause: err,
    })
  }

  // Plain Error or unknown
  if (err instanceof Error)
    return new AppError({ message: err.message || 'Unexpected error', cause: err })

  return new AppError({ message: 'Unexpected error', cause: err })
}

/** Shorthand check */
export function isValidationError(e: unknown): e is AppError {
  return e instanceof AppError && e.isValidation
}

/** Get first error message for a given field (e.g., "email") */
export function getFieldError(e: unknown, field: string): string | undefined {
  const app = toAppError(e)
  const list = app.fieldMap[field]

  return list?.[0]
}

/** Convert to `{ field: message }` map (first message per field) for simple forms */
export function toFirstFieldErrorMap(e: unknown): Record<string, string> {
  const app = toAppError(e)
  const out: Record<string, string> = {}
  for (const [k, v] of Object.entries(app.fieldMap)) {
    if (v.length)
      out[k] = v[0]
  }

  return out
}

/** Collect human-friendly messages for toast/snackbar */
export function summarizeError(e: unknown): string {
  const app = toAppError(e)
  if (app.isValidation && Object.keys(app.fieldMap).length) {
    // e.g., "Email already registered"
    const first = Object.values(app.fieldMap)[0]?.[0]
    if (first)
      return first
  }
  if (app.problem?.detail)
    return app.problem.detail
  if (app.status)
    return `${app.message} (HTTP ${app.status})`

  return app.message
}
