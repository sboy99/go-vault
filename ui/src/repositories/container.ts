import { HttpClient, type TokenGetter, type UnauthorizedHandler } from '@/lib/http/client'
import { createBackupRepository } from '@/repositories/backup.repository'
import { createJobRepository } from '@/repositories/job.repository'
import { createRestoreRepository } from '@/repositories/restore.repository'
import { createHealthRepository } from '@/repositories/health.repository'
import type { Repositories } from '@/repositories/types'

export interface CreateRepositoriesOptions {
  getToken?: TokenGetter
  onUnauthorized?: UnauthorizedHandler
  baseUrl?: string
  fetchImpl?: typeof fetch
}

export function createRepositories(options: CreateRepositoriesOptions = {}): Repositories {
  const client = new HttpClient({
    baseUrl: options.baseUrl,
    getToken: options.getToken,
    onUnauthorized: options.onUnauthorized,
    fetchImpl: options.fetchImpl,
  })
  return {
    backups: createBackupRepository(client),
    jobs: createJobRepository(client),
    restores: createRestoreRepository(client),
    health: createHealthRepository(client),
  }
}

export type { Repositories } from '@/repositories/types'
