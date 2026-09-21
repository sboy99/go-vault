export type AsyncStatus = 'idle' | 'loading' | 'success' | 'error'

export interface AsyncState<T> {
  status: AsyncStatus
  data: T | null
  error: string | null
}

export function idleAsync<T>(data: T | null = null): AsyncState<T> {
  return { status: 'idle', data, error: null }
}

export function loadingAsync<T>(data: T | null = null): AsyncState<T> {
  return { status: 'loading', data, error: null }
}

export function successAsync<T>(data: T): AsyncState<T> {
  return { status: 'success', data, error: null }
}

export function errorAsync<T>(error: string, data: T | null = null): AsyncState<T> {
  return { status: 'error', data, error }
}
