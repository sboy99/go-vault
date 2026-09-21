import type { Backup } from '@/domain/backup'
import type { Job } from '@/domain/job'
import type { HttpClient } from '@/lib/http/client'
import type { BackupListParams, BackupRepository } from '@/repositories/types'

function toQuery(params?: BackupListParams): string {
  if (!params) return ''
  const q = new URLSearchParams()
  if (params.limit != null) q.set('limit', String(params.limit))
  if (params.offset != null) q.set('offset', String(params.offset))
  if (params.status) q.set('status', params.status)
  const s = q.toString()
  return s ? `?${s}` : ''
}

export function createBackupRepository(client: HttpClient): BackupRepository {
  return {
    list(params) {
      return client.get<Backup[]>(`/v1/backups${toQuery(params)}`)
    },
    getById(id) {
      return client.get<Backup>(`/v1/backups/${encodeURIComponent(id)}`)
    },
    create() {
      return client.post<Job>('/v1/backups')
    },
    async download(id) {
      const response = await client.request<Response>(`/v1/backups/${encodeURIComponent(id)}/download`, {
        method: 'GET',
        raw: true,
      })
      return response.blob()
    },
  }
}
