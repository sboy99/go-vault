import { describe, expect, it, vi } from 'vitest'
import { HttpClient } from '@/lib/http/client'
import { ApiError } from '@/lib/http/errors'

describe('HttpClient', () => {
  it('unwraps success envelope', async () => {
    const fetchImpl = vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ success: true, data: { id: '1' }, error: null }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      }),
    )
    const client = new HttpClient({ baseUrl: '/api', fetchImpl })
    const data = await client.get<{ id: string }>('/v1/jobs/1')
    expect(data).toEqual({ id: '1' })
    expect(fetchImpl).toHaveBeenCalledWith(
      '/api/v1/jobs/1',
      expect.objectContaining({ method: 'GET' }),
    )
  })

  it('attaches bearer token', async () => {
    const fetchImpl = vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ success: true, data: [], error: null }), { status: 200 }),
    )
    const client = new HttpClient({
      baseUrl: '/api',
      getToken: () => 'secret',
      fetchImpl,
    })
    await client.get('/v1/backups')
    const init = fetchImpl.mock.calls[0][1] as RequestInit
    expect((init.headers as Record<string, string>).Authorization).toBe('Bearer secret')
  })

  it('throws ApiError on failure and calls onUnauthorized for 401', async () => {
    const onUnauthorized = vi.fn()
    const fetchImpl = vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ success: false, data: null, error: 'unauthorized' }), {
        status: 401,
      }),
    )
    const client = new HttpClient({ baseUrl: '/api', fetchImpl, onUnauthorized })
    await expect(client.get('/v1/backups')).rejects.toBeInstanceOf(ApiError)
    expect(onUnauthorized).toHaveBeenCalled()
  })

  it('marks 409 as conflict', async () => {
    const fetchImpl = vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ success: false, data: null, error: 'job already running' }), {
        status: 409,
      }),
    )
    const client = new HttpClient({ baseUrl: '/api', fetchImpl })
    try {
      await client.post('/v1/backups')
      expect.unreachable()
    } catch (e) {
      expect(e).toBeInstanceOf(ApiError)
      expect((e as ApiError).isConflict).toBe(true)
    }
  })
})
