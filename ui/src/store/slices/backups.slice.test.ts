import { describe, expect, it, vi, beforeEach } from 'vitest'
import { createStore } from 'solid-js/store'
import { createBackupsActions, initialBackupsState } from '@/store/slices/backups.slice'
import { createMemoryRepositories } from '@/repositories/memory'
import type { Backup } from '@/domain/backup'

const sample: Backup = {
  id: 'b1',
  name: 'db.dump',
  storage_key: 'db.dump',
  database_type: 'POSTGRESQL',
  storage_type: 'LOCAL',
  status: 'success',
  started_at: '2026-01-01T00:00:00Z',
  created_at: '2026-01-01T00:00:00Z',
  size_bytes: 1024,
  format: 'custom',
  verified: true,
}

describe('backups slice', () => {
  beforeEach(() => {
    vi.useRealTimers()
  })

  it('loads backups from repository', async () => {
    const repos = createMemoryRepositories({ backups: [sample] })
    const [state, setState] = createStore(initialBackupsState)
    const notify = vi.fn()
    const actions = createBackupsActions(
      () => state,
      setState,
      repos,
      notify,
      vi.fn(),
    )
    await actions.fetchList()
    expect(state.list.status).toBe('success')
    expect(state.list.data).toHaveLength(1)
    expect(state.list.data?.[0].id).toBe('b1')
  })

  it('creates a backup job', async () => {
    const repos = createMemoryRepositories({ backups: [] })
    const [state, setState] = createStore(initialBackupsState)
    const notify = vi.fn()
    const actions = createBackupsActions(
      () => state,
      setState,
      repos,
      notify,
      vi.fn(),
    )
    const id = await actions.createBackup()
    expect(id).toMatch(/^job_/)
    expect(notify).toHaveBeenCalledWith('success', 'Backup started', expect.any(String))
  })
})
