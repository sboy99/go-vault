import type { Job } from '@/domain/job'
import type { HttpClient } from '@/lib/http/client'
import type { JobRepository, ListParams } from '@/repositories/types'

function toQuery(params?: ListParams): string {
  if (!params) return ''
  const q = new URLSearchParams()
  if (params.limit != null) q.set('limit', String(params.limit))
  if (params.offset != null) q.set('offset', String(params.offset))
  const s = q.toString()
  return s ? `?${s}` : ''
}

export function createJobRepository(client: HttpClient): JobRepository {
  return {
    list(params) {
      return client.get<Job[]>(`/v1/jobs${toQuery(params)}`)
    },
    getById(id) {
      return client.get<Job>(`/v1/jobs/${encodeURIComponent(id)}`)
    },
  }
}
