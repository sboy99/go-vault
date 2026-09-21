import type { Job } from '@/domain/job'
import type { RestoreRequest } from '@/domain/restore'
import type { HttpClient } from '@/lib/http/client'
import type { RestoreRepository } from '@/repositories/types'

export function createRestoreRepository(client: HttpClient): RestoreRepository {
  return {
    create(request: RestoreRequest) {
      return client.post<Job>('/v1/restores', request)
    },
  }
}
