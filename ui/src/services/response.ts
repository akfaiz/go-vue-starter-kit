export interface ApiEnvelope<T> {
  status: number
  data: T
  message?: string
}

export interface ApiMessage {
  status: number
  message: string
}
