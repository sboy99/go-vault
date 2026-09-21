import type { Backup, BackupStatus } from '@/domain/backup'
import type { Job } from '@/domain/job'
import type { RestoreRequest } from '@/domain/restore'
import type { HealthStatus } from '@/domain/health'

export interface ListParams {
  limit?: number
  offset?: number
}

export interface BackupListParams extends ListParams {
  status?: BackupStatus | ''
}

export interface BackupRepository {
  list(params?: BackupListParams): Promise<Backup[]>
  getById(id: string): Promise<Backup>
  create(): Promise<Job>
  download(id: string): Promise<Blob>
}

export interface JobRepository {
  list(params?: ListParams): Promise<Job[]>
  getById(id: string): Promise<Job>
}

export interface RestoreRepository {
  create(request: RestoreRequest): Promise<Job>
}

export interface HealthRepository {
  healthz(): Promise<HealthStatus>
  readyz(): Promise<HealthStatus>
}

export interface Repositories {
  backups: BackupRepository
  jobs: JobRepository
  restores: RestoreRepository
  health: HealthRepository
}
