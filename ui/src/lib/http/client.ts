import { env } from '@/config/env'
import type { Envelope } from '@/lib/http/envelope'
import { ApiError } from '@/lib/http/errors'

export type TokenGetter = () => string | null | undefined
export type UnauthorizedHandler = () => void

export interface HttpClientOptions {
  baseUrl?: string
  getToken?: TokenGetter
  onUnauthorized?: UnauthorizedHandler
  fetchImpl?: typeof fetch
}

export interface RequestOptions {
  method?: string
  body?: unknown
  headers?: Record<string, string>
  signal?: AbortSignal
  raw?: boolean
}

export class HttpClient {
  private readonly baseUrl: string
  private readonly getToken: TokenGetter
  private readonly onUnauthorized?: UnauthorizedHandler
  private readonly fetchImpl: typeof fetch

  constructor(options: HttpClientOptions = {}) {
    this.baseUrl = (options.baseUrl ?? env.apiBaseUrl).replace(/\/$/, '')
    this.getToken = options.getToken ?? (() => null)
    this.onUnauthorized = options.onUnauthorized
    this.fetchImpl = options.fetchImpl ?? fetch.bind(globalThis)
  }

  async request<T>(path: string, options: RequestOptions = {}): Promise<T> {
    const url = path.startsWith('http') ? path : `${this.baseUrl}${path.startsWith('/') ? '' : '/'}${path}`
    const headers: Record<string, string> = {
      Accept: 'application/json',
      ...options.headers,
    }
    const token = this.getToken()
    if (token) {
      headers.Authorization = `Bearer ${token}`
    }
    let body: BodyInit | undefined
    if (options.body !== undefined) {
      headers['Content-Type'] = 'application/json'
      body = JSON.stringify(options.body)
    }

    const response = await this.fetchImpl(url, {
      method: options.method ?? (options.body !== undefined ? 'POST' : 'GET'),
      headers,
      body,
      signal: options.signal,
    })

    if (options.raw) {
      if (!response.ok) {
        const text = await response.text().catch(() => null)
        if (response.status === 401) this.onUnauthorized?.()
        throw new ApiError(text || response.statusText || 'request failed', response.status, text)
      }
      return response as unknown as T
    }

    const text = await response.text()
    let envelope: Envelope<T> | null = null
    if (text) {
      try {
        envelope = JSON.parse(text) as Envelope<T>
      } catch {
        if (!response.ok) {
          if (response.status === 401) this.onUnauthorized?.()
          throw new ApiError(text || response.statusText, response.status, text)
        }
        throw new ApiError('invalid JSON response', response.status, text)
      }
    }

    if (!response.ok || !envelope?.success) {
      const message = envelope?.error || response.statusText || 'request failed'
      if (response.status === 401) this.onUnauthorized?.()
      throw new ApiError(message, response.status, text)
    }

    return envelope.data as T
  }

  get<T>(path: string, options?: Omit<RequestOptions, 'method' | 'body'>): Promise<T> {
    return this.request<T>(path, { ...options, method: 'GET' })
  }

  post<T>(path: string, body?: unknown, options?: Omit<RequestOptions, 'method' | 'body'>): Promise<T> {
    return this.request<T>(path, { ...options, method: 'POST', body })
  }
}
