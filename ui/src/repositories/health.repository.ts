import type { HealthStatus } from '@/domain/health'
import type { HttpClient } from '@/lib/http/client'
import type { HealthRepository } from '@/repositories/types'

export function createHealthRepository(client: HttpClient): HealthRepository {
  return {
    healthz() {
      return client.get<HealthStatus>('/healthz')
    },
    readyz() {
      return client.get<HealthStatus>('/readyz')
    },
  }
}
