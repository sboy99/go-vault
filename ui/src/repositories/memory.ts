import type { Backup } from '@/domain/backup'
import type { Job } from '@/domain/job'
import type { HealthStatus } from '@/domain/health'
import type {
  BackupRepository,
  HealthRepository,
  JobRepository,
  Repositories,
  RestoreRepository,
} from '@/repositories/types'

export function createMemoryRepositories(seed?: {
  backups?: Backup[]
  jobs?: Job[]
}): Repositories & {
  _backups: Backup[]
  _jobs: Job[]
} {
  const backups = [...(seed?.backups ?? [])]
  const jobs = [...(seed?.jobs ?? [])]

  const backupRepo: BackupRepository = {
    async list(params) {
      let list = [...backups]
      if (params?.status) {
        list = list.filter((b) => b.status === params.status)
      }
      const offset = params?.offset ?? 0
      const limit = params?.limit ?? list.length
      return list.slice(offset, offset + limit)
    },
    async getById(id) {
      const found = backups.find((b) => b.id === id)
      if (!found) throw new Error('backup not found')
      return found
    },
    async create() {
      const job: Job = {
        id: `job_${jobs.length + 1}`,
        type: 'backup',
        status: 'running',
        created_at: new Date().toISOString(),
        started_at: new Date().toISOString(),
      }
      jobs.push(job)
      return job
    },
    async download() {
      return new Blob(['fake'], { type: 'application/octet-stream' })
    },
  }

  const jobRepo: JobRepository = {
    async list(params) {
      const offset = params?.offset ?? 0
      const limit = params?.limit ?? jobs.length
      return jobs.slice(offset, offset + limit)
    },
    async getById(id) {
      const found = jobs.find((j) => j.id === id)
      if (!found) throw new Error('job not found')
      return found
    },
  }

  const restoreRepo: RestoreRepository = {
    async create(request) {
      if (request.confirm !== 'testdb') {
        throw Object.assign(new Error('restore confirm mismatch'), { status: 400 })
      }
      const job: Job = {
        id: `job_${jobs.length + 1}`,
        type: 'restore',
        status: 'running',
        backup_id: request.backup_id,
        created_at: new Date().toISOString(),
        started_at: new Date().toISOString(),
      }
      jobs.push(job)
      return job
    },
  }

  const healthRepo: HealthRepository = {
    async healthz(): Promise<HealthStatus> {
      return { status: 'ok' }
    },
    async readyz(): Promise<HealthStatus> {
      return { status: 'ready' }
    },
  }

  return {
    backups: backupRepo,
    jobs: jobRepo,
    restores: restoreRepo,
    health: healthRepo,
    _backups: backups,
    _jobs: jobs,
  }
}
