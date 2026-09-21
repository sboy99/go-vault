import type { Backup, BackupStatus } from '@/domain/backup'
import { isApiError } from '@/lib/http/errors'
import type { Repositories } from '@/repositories/types'
import {
  errorAsync,
  idleAsync,
  loadingAsync,
  successAsync,
  type AsyncState,
} from '@/store/async-state'
import { logAction } from '@/store/create-store'
import type { SetStoreFunction } from 'solid-js/store'

export interface BackupsState {
  list: AsyncState<Backup[]>
  selected: AsyncState<Backup>
  statusFilter: BackupStatus | ''
  creating: boolean
}

export const initialBackupsState: BackupsState = {
  list: idleAsync<Backup[]>([]),
  selected: idleAsync<Backup>(),
  statusFilter: '',
  creating: false,
}

export function createBackupsActions(
  get: () => BackupsState,
  set: SetStoreFunction<BackupsState>,
  repos: Repositories,
  notify: (kind: 'success' | 'error' | 'info', title: string, description?: string) => void,
  onUnauthorized: () => void,
) {
  return {
    async fetchList() {
      logAction('backups', 'fetchList')
      const prev = get().list.data
      set('list', loadingAsync(prev ?? []))
      try {
        const status = get().statusFilter
        const data = await repos.backups.list({
          limit: 100,
          status: status || undefined,
        })
        set('list', successAsync(data))
      } catch (e) {
        if (isApiError(e) && e.isUnauthorized) onUnauthorized()
        const message = e instanceof Error ? e.message : 'Failed to load backups'
        set('list', errorAsync(message, prev ?? []))
        notify('error', 'Backups', message)
      }
    },

    setStatusFilter(status: BackupStatus | '') {
      set('statusFilter', status)
    },

    async fetchById(id: string) {
      logAction('backups', 'fetchById', id)
      set('selected', loadingAsync())
      try {
        const data = await repos.backups.getById(id)
        set('selected', successAsync(data))
      } catch (e) {
        if (isApiError(e) && e.isUnauthorized) onUnauthorized()
        const message = e instanceof Error ? e.message : 'Failed to load backup'
        set('selected', errorAsync(message))
        notify('error', 'Backup', message)
      }
    },

    clearSelected() {
      set('selected', idleAsync())
    },

    async createBackup(): Promise<string | null> {
      logAction('backups', 'createBackup')
      set('creating', true)
      try {
        const job = await repos.backups.create()
        notify('success', 'Backup started', `Job ${job.id}`)
        return job.id
      } catch (e) {
        if (isApiError(e) && e.isUnauthorized) onUnauthorized()
        if (isApiError(e) && e.isConflict) {
          notify('info', 'Job already running', e.message)
          return null
        }
        const message = e instanceof Error ? e.message : 'Failed to start backup'
        notify('error', 'Backup', message)
        return null
      } finally {
        set('creating', false)
      }
    },

    async download(id: string, filename: string) {
      try {
        const blob = await repos.backups.download(id)
        const url = URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.href = url
        a.download = filename
        a.click()
        URL.revokeObjectURL(url)
        notify('success', 'Download started', filename)
      } catch (e) {
        if (isApiError(e) && e.isUnauthorized) onUnauthorized()
        const message = e instanceof Error ? e.message : 'Download failed'
        notify('error', 'Download', message)
      }
    },
  }
}

export type BackupsActions = ReturnType<typeof createBackupsActions>
