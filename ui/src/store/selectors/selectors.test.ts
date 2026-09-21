import { describe, expect, it } from 'vitest'
import { selectBackupCounts, selectLatestBackup } from '@/store/selectors/backups.selectors'
import { selectHasActiveJob } from '@/store/selectors/jobs.selectors'
import { idleAsync, successAsync } from '@/store/async-state'
import type { Backup } from '@/domain/backup'
import type { Job } from '@/domain/job'
import type { BackupsState } from '@/store/slices/backups.slice'
import type { JobsState } from '@/store/slices/jobs.slice'

const backup = (partial: Partial<Backup> & Pick<Backup, 'id' | 'status' | 'created_at'>): Backup => ({
  name: partial.name ?? partial.id,
  storage_key: partial.storage_key ?? partial.id,
  database_type: 'POSTGRESQL',
  storage_type: 'LOCAL',
  size_bytes: partial.size_bytes ?? 100,
  format: 'custom',
  verified: false,
  started_at: partial.created_at,
  ...partial,
})

describe('backup selectors', () => {
  it('counts and latest', () => {
    const state: BackupsState = {
      list: successAsync([
        backup({ id: 'a', status: 'success', created_at: '2026-01-01T00:00:00Z', size_bytes: 10 }),
        backup({ id: 'b', status: 'failed', created_at: '2026-01-02T00:00:00Z', size_bytes: 20 }),
        backup({ id: 'c', status: 'running', created_at: '2026-01-03T00:00:00Z', size_bytes: 0 }),
      ]),
      selected: idleAsync(),
      statusFilter: '',
      creating: false,
    }
    const counts = selectBackupCounts(state)
    expect(counts.total).toBe(3)
    expect(counts.success).toBe(1)
    expect(counts.failed).toBe(1)
    expect(counts.running).toBe(1)
    expect(counts.storageBytes).toBe(30)
    expect(selectLatestBackup(state)?.id).toBe('c')
  })
})

describe('job selectors', () => {
  it('detects active jobs', () => {
    const jobs: Job[] = [
      {
        id: '1',
        type: 'backup',
        status: 'running',
        created_at: '2026-01-01T00:00:00Z',
      },
      {
        id: '2',
        type: 'restore',
        status: 'succeeded',
        created_at: '2026-01-01T00:00:00Z',
      },
    ]
    const state: JobsState = {
      list: successAsync(jobs),
      selected: idleAsync(),
      polling: true,
    }
    expect(selectHasActiveJob(state)).toBe(true)
  })
})
