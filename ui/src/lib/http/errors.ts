export class ApiError extends Error {
  readonly status: number
  readonly body: string | null

  constructor(message: string, status: number, body: string | null = null) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.body = body
  }

  get isUnauthorized(): boolean {
    return this.status === 401
  }

  get isConflict(): boolean {
    return this.status === 409
  }
}

export function isApiError(error: unknown): error is ApiError {
  return error instanceof ApiError
}
