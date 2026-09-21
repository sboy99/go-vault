import type { Job } from '@/domain/job'
import { isActiveJob } from '@/domain/job'
import { JOB_POLL_INTERVAL_MS } from '@/config/constants'
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

export interface JobsState {
  list: AsyncState<Job[]>
  selected: AsyncState<Job>
  polling: boolean
}

export const initialJobsState: JobsState = {
  list: idleAsync<Job[]>([]),
  selected: idleAsync<Job>(),
  polling: false,
}

let pollTimer: ReturnType<typeof setInterval> | null = null

export function createJobsActions(
  get: () => JobsState,
  set: SetStoreFunction<JobsState>,
  repos: Repositories,
  notify: (kind: 'success' | 'error' | 'info', title: string, description?: string) => void,
  onUnauthorized: () => void,
  onJobSettled?: (job: Job) => void,
) {
  const actions = {
    async fetchList() {
      logAction('jobs', 'fetchList')
      const prev = get().list.data
      set('list', loadingAsync(prev ?? []))
      try {
        const data = await repos.jobs.list({ limit: 100 })
        set('list', successAsync(data))
        actions.syncPolling()
      } catch (e) {
        if (isApiError(e) && e.isUnauthorized) onUnauthorized()
        const message = e instanceof Error ? e.message : 'Failed to load jobs'
        set('list', errorAsync(message, prev ?? []))
        notify('error', 'Jobs', message)
      }
    },

    async fetchById(id: string) {
      logAction('jobs', 'fetchById', id)
      set('selected', loadingAsync())
      try {
        const data = await repos.jobs.getById(id)
        set('selected', successAsync(data))
        // merge into list
        const list = get().list.data ?? []
        const idx = list.findIndex((j) => j.id === id)
        if (idx >= 0) {
          set('list', 'data', (prev) => {
            if (!prev) return prev
            return prev.map((j, i) => (i === idx ? data : j))
          })
        } else {
          set('list', 'data', (prev) => [...(prev ?? []), data])
        }
        return data
      } catch (e) {
        if (isApiError(e) && e.isUnauthorized) onUnauthorized()
        const message = e instanceof Error ? e.message : 'Failed to load job'
        set('selected', errorAsync(message))
        notify('error', 'Job', message)
        return null
      }
    },

    upsertJob(job: Job) {
      set('list', 'data', (prev) => {
        const list = prev ?? []
        const idx = list.findIndex((j) => j.id === job.id)
        if (idx >= 0) {
          return list.map((j, i) => (i === idx ? job : j))
        }
        return [job, ...list]
      })
      actions.syncPolling()
    },

    syncPolling() {
      const list = get().list.data ?? []
      const hasActive = list.some(isActiveJob)
      if (hasActive && !pollTimer) {
        set('polling', true)
        pollTimer = setInterval(() => {
          void actions.pollActive()
        }, JOB_POLL_INTERVAL_MS)
      } else if (!hasActive && pollTimer) {
        clearInterval(pollTimer)
        pollTimer = null
        set('polling', false)
      }
    },

    async pollActive() {
      const list = get().list.data ?? []
      const active = list.filter(isActiveJob)
      if (active.length === 0) {
        actions.syncPolling()
        return
      }
      for (const job of active) {
        try {
          const updated = await repos.jobs.getById(job.id)
          set('list', 'data', (prev) =>
            (prev ?? []).map((j) => (j.id === updated.id ? updated : j)),
          )
          if (!isActiveJob(updated)) {
            onJobSettled?.(updated)
            if (updated.status === 'succeeded') {
              notify('success', `${updated.type} finished`, `Job ${updated.id}`)
            } else if (updated.status === 'failed') {
              notify('error', `${updated.type} failed`, updated.error || `Job ${updated.id}`)
            }
          }
        } catch {
          // ignore transient poll errors
        }
      }
      actions.syncPolling()
    },

    stopPolling() {
      if (pollTimer) {
        clearInterval(pollTimer)
        pollTimer = null
      }
      set('polling', false)
    },
  }
  return actions
}

export type JobsActions = ReturnType<typeof createJobsActions>
