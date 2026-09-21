export interface Envelope<T> {
  success: boolean
  data: T | null
  error: string | null
}
